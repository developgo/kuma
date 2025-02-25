package gateway

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/model"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	"github.com/kumahq/kuma/pkg/plugins/runtime/gateway/match"
	"github.com/kumahq/kuma/pkg/plugins/runtime/gateway/merge"
	"github.com/kumahq/kuma/pkg/plugins/runtime/gateway/route"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	envoy_listeners "github.com/kumahq/kuma/pkg/xds/envoy/listeners"
	envoy_names "github.com/kumahq/kuma/pkg/xds/envoy/names"
	envoy_routes "github.com/kumahq/kuma/pkg/xds/envoy/routes"
)

const WildcardHostname = "*"

// RoutePolicyTypes specifies the resource types the gateway will bind
// for routes.
var RoutePolicyTypes = []model.ResourceType{
	core_mesh.GatewayRouteType,
}

// ConnectionPolicyTypes specifies the resource types the gateway will
// bind for connection policies.
var ConnectionPolicyTypes = []model.ResourceType{
	core_mesh.CircuitBreakerType,
	core_mesh.FaultInjectionType,
	core_mesh.HealthCheckType,
	core_mesh.RateLimitType,
	core_mesh.RetryType,
	core_mesh.TimeoutType,
}

type GatewayHost struct {
	Hostname string
	Routes   []model.Resource
	Policies map[model.ResourceType][]match.RankedPolicy
	TLS      *mesh_proto.Gateway_TLS_Conf
}

// Resources tracks partially-built xDS resources that can be updated
// by multiple gateway generators.
type Resources struct {
	Listener           *envoy_listeners.ListenerBuilder
	FilterChain        *envoy_listeners.FilterChainBuilder
	RouteConfiguration *envoy_routes.RouteConfigurationBuilder
}

type GatewayListener struct {
	Port         uint32
	Protocol     mesh_proto.Gateway_Listener_Protocol
	ResourceName string
}

type GatewayResourceInfo struct {
	Proxy            *core_xds.Proxy
	Dataplane        *core_mesh.DataplaneResource
	Gateway          *core_mesh.GatewayResource
	ExternalServices *core_mesh.ExternalServiceResourceList

	Listener   GatewayListener
	Host       GatewayHost
	Resources  Resources
	RouteTable route.Table
}

// GatewayHostGenerator is responsible for generating xDS resources for a single GatewayHost.
type GatewayHostGenerator interface {
	GenerateHost(xds_context.Context, *GatewayResourceInfo) (*core_xds.ResourceSet, error)
	SupportsProtocol(mesh_proto.Gateway_Listener_Protocol) bool
}

// Generator generates xDS resources for an entire Gateway.
type Generator struct {
	Generators []GatewayHostGenerator
}

func (g Generator) Generate(ctx xds_context.Context, proxy *core_xds.Proxy) (*core_xds.ResourceSet, error) {
	gateway := match.Gateway(ctx.Mesh.Resources.Gateways(), proxy.Dataplane.Spec.Matches)

	if gateway == nil {
		log.V(1).Info("no matching gateway for dataplane",
			"name", proxy.Dataplane.Meta.GetName(),
			"mesh", proxy.Dataplane.Meta.GetMesh(),
			"service", proxy.Dataplane.Spec.GetIdentifyingService(),
		)

		return nil, nil
	}

	log.V(1).Info(fmt.Sprintf("matched gateway %q to dataplane %q",
		gateway.Meta.GetName(), proxy.Dataplane.Meta.GetName()))

	// Canonicalize the tags on each listener to be the merged resources
	// of dataplane, gateway and listener tags.
	for _, listener := range gateway.Spec.GetConf().GetListeners() {
		listener.Tags = match.MergeSelectors(
			proxy.Dataplane.Spec.GetNetworking().GetGateway().GetTags(),
			gateway.Spec.GetTags(),
			listener.GetTags(),
		)
	}

	// Multiple listener specifications can have the same port. If
	// they are compatible, then we can collapse those specifications
	// down to a single listener.
	collapsed := map[uint32][]*mesh_proto.Gateway_Listener{}
	for _, ep := range gateway.Spec.GetConf().GetListeners() {
		collapsed[ep.GetPort()] = append(collapsed[ep.GetPort()], ep)
	}

	resources := ResourceAggregator{core_xds.NewResourceSet()}

	externalServices := ctx.Mesh.Resources.ExternalServices()

	for port, listeners := range collapsed {
		// Force all listeners on the same port to have the same protocol.
		for i := range listeners {
			if listeners[i].GetProtocol() != listeners[0].GetProtocol() {
				return nil, errors.Errorf(
					"cannot collapse listener protocols %s and %s on port %d",
					listeners[i].GetProtocol(), listeners[0].GetProtocol(), port,
				)
			}
		}

		listener, hosts, err := MakeGatewayListener(ctx.Mesh, gateway, listeners)
		if err != nil {
			return nil, err
		}

		hosts = RedistributeWildcardRoutes(hosts)

		// Sort by reverse hostname, so that fully qualified hostnames sort
		// before wildcard domains, and "*" is last.
		sort.Slice(hosts, func(i, j int) bool {
			return hosts[i].Hostname > hosts[j].Hostname
		})

		info := GatewayResourceInfo{
			Proxy:            proxy,
			Dataplane:        proxy.Dataplane,
			Gateway:          gateway,
			ExternalServices: externalServices,
			Listener:         listener,
		}

		// Make a pass over the generators for each virtual host.
		for _, host := range hosts {
			info.RouteTable = route.Table{}
			info.Host = host

			// Ensure that generators don't get duplicate routes,
			// which could happen after redistributing wildcards.
			info.Host.Routes = merge.UniqueResources(info.Host.Routes)

			for _, generator := range g.Generators {
				if !generator.SupportsProtocol(listener.Protocol) {
					continue
				}

				if err := resources.AddSet(generator.GenerateHost(ctx, &info)); err != nil {
					return nil, errors.Wrapf(err, "%T failed to generate resources for dataplane %q",
						generator, proxy.Id)
				}
			}
		}

		info.Resources.Listener.Configure(envoy_listeners.FilterChain(info.Resources.FilterChain))

		if err := resources.AddSet(BuildResourceSet(info.Resources.Listener)); err != nil {
			return nil, errors.Wrapf(err, "failed to build listener resource")
		}

		if err := resources.AddSet(BuildResourceSet(info.Resources.RouteConfiguration)); err != nil {
			return nil, errors.Wrapf(err, "failed to build route configuration resource")
		}
	}

	return resources.Get(), nil
}

// MakeGatewayListener converts a collapsed set of listener configurations
// in to a single configuration with a matched set of route resources. The
// given listeners must have a consistent protocol and port.
func MakeGatewayListener(
	meshContext xds_context.MeshContext,
	gateway *core_mesh.GatewayResource,
	listeners []*mesh_proto.Gateway_Listener,
) (GatewayListener, []GatewayHost, error) {
	hostsByName := map[string]GatewayHost{}

	listener := GatewayListener{
		Port:     listeners[0].GetPort(),
		Protocol: listeners[0].GetProtocol(),
		ResourceName: envoy_names.GetGatewayListenerName(
			gateway.Meta.GetName(),
			listeners[0].GetProtocol().String(),
			listeners[0].GetPort(),
		),
	}

	// Hostnames must be unique to a listener to remove ambiguity
	// in policy selection and TLS configuration.
	for _, l := range listeners {
		// An empty hostname is the same as "*", i.e. matches all hosts.
		hostname := l.GetHostname()
		if hostname == "" {
			hostname = WildcardHostname
		}

		if _, ok := hostsByName[hostname]; ok {
			return listener, nil, errors.Errorf("duplicate hostname %q", hostname)
		}

		host := GatewayHost{
			Hostname: hostname,
			Policies: map[model.ResourceType][]match.RankedPolicy{},
			TLS:      l.GetTls(),
		}

		switch listener.Protocol {
		case mesh_proto.Gateway_Listener_HTTP,
			mesh_proto.Gateway_Listener_HTTPS:
			host.Routes = append(host.Routes,
				match.Routes(meshContext.Resources.GatewayRoutes(), l.GetTags())...)
		default:
			// TODO(jpeach) match other route types that are appropriate to the protocol.
		}

		for _, t := range ConnectionPolicyTypes {
			matches := match.ConnectionPoliciesBySource(
				l.GetTags(),
				match.ToConnectionPolicies(meshContext.Resources[t]))
			host.Policies[t] = matches
		}

		hostsByName[hostname] = host
	}

	var hosts []GatewayHost
	for _, host := range hostsByName {
		hosts = append(hosts, host)
	}

	return listener, hosts, nil
}

// RedistributeWildcardRoutes takes the routes from the wildcard host
// and redistributes them to hosts with matching names, creating new
// hosts if necessary.
//
// This process is necessary because:
//
// 1. We might have a listener with hostname A and some routes, but also
//    a wildcard listener with routes for hostname A. We want all the routes
//    for hostname A in the same virtual host.
// 2. Routes with hostnames that are attached to a wildcard listener
//    should implicitly create virtual hosts so that we can generate a
//    consistent config. For example, if a wildcard listener has a route for
//    hostname A and a route for hostname B, that doesn't mean that the routes
//    are for hostnames A or B. We still want the routes to match the hostname
//    that they were specified with.
func RedistributeWildcardRoutes(
	hosts []GatewayHost,
) []GatewayHost {
	hostsByName := map[string]GatewayHost{}

	for _, h := range hosts {
		hostsByName[h.Hostname] = h
	}

	wild, ok := hostsByName[WildcardHostname]
	if !ok {
		return hosts
	}

	wildcardRoutes := wild.Routes
	wild.Routes = nil // We are rebuilding this.
	for _, r := range wildcardRoutes {
		gw, ok := r.(*core_mesh.GatewayRouteResource)
		if !ok {
			continue
		}

		names := gw.Spec.GetConf().GetHttp().GetHostnames()

		// No hostnames on this route, it stays as a wildcard route.
		if len(names) == 0 {
			wild.Routes = append(wild.Routes, r)
			continue
		}

		for _, n := range names {
			host, ok := hostsByName[n]

			// When generating a new implicit virtualhost,
			// initialize it by shallow copying from the
			// wildcard source.
			if !ok {
				host = wild
				host.Routes = nil
			}

			// Note that if we already have a virtualhost for this
			// name, and add the route to it, it might be a duplicate.
			host.Routes = append(host.Routes, r)
			host.Hostname = n

			hostsByName[n] = host
		}
	}

	hostsByName[WildcardHostname] = wild

	var flattened []GatewayHost
	for _, host := range hostsByName {
		flattened = append(flattened, host)
	}

	return flattened
}
