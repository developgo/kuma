package api_server_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	config "github.com/kumahq/kuma/pkg/config/api-server"
	"github.com/kumahq/kuma/pkg/core"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/manager"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/resources/store"
	"github.com/kumahq/kuma/pkg/metrics"
	"github.com/kumahq/kuma/pkg/plugins/resources/memory"
	"github.com/kumahq/kuma/pkg/test/kds/samples"
	"github.com/kumahq/kuma/pkg/test/matchers"
	test_model "github.com/kumahq/kuma/pkg/test/resources/model"
	util_proto "github.com/kumahq/kuma/pkg/util/proto"
)

type dataplaneBuilder core_mesh.DataplaneResource

func newMesh(name string) *core_mesh.MeshResource {
	return &core_mesh.MeshResource{
		Meta: &test_model.ResourceMeta{Name: name},
		Spec: &mesh_proto.Mesh{},
	}
}

func newDataplane() *dataplaneBuilder {
	return &dataplaneBuilder{
		Spec: &mesh_proto.Dataplane{
			Networking: &mesh_proto.Dataplane_Networking{
				Address: "1.1.1.1",
			},
		},
	}
}

func (b *dataplaneBuilder) build() *core_mesh.DataplaneResource {
	return (*core_mesh.DataplaneResource)(b)
}

func (b *dataplaneBuilder) meta(name, mesh string) *dataplaneBuilder {
	b.Meta = &test_model.ResourceMeta{Name: name, Mesh: mesh}
	return b
}

func (b *dataplaneBuilder) inbound(service, ip string, dpPort, workloadPort uint32) *dataplaneBuilder {
	b.Spec.Networking.Inbound = append(b.Spec.Networking.Inbound, &mesh_proto.Dataplane_Networking_Inbound{
		Address:     ip,
		Port:        dpPort,
		ServicePort: workloadPort,
		Tags: map[string]string{
			mesh_proto.ServiceTag:  service,
			mesh_proto.ProtocolTag: "http",
		},
	})
	return b
}

func (b *dataplaneBuilder) outbound(service, ip string, port uint32) *dataplaneBuilder {
	b.Spec.Networking.Outbound = append(b.Spec.Networking.Outbound, &mesh_proto.Dataplane_Networking_Outbound{
		Address: ip,
		Port:    port,
		Tags: map[string]string{
			mesh_proto.ServiceTag: service,
		},
	})
	return b
}

type selectors []*mesh_proto.Selector

func anyService() []*mesh_proto.Selector {
	return []*mesh_proto.Selector{
		{
			Match: map[string]string{
				mesh_proto.ServiceTag: "*",
			},
		},
	}
}

func serviceSelector(name, protocol string) *mesh_proto.Selector {
	if protocol == "" {
		return &mesh_proto.Selector{
			Match: map[string]string{
				mesh_proto.ServiceTag: name,
			},
		}
	} else {
		return &mesh_proto.Selector{
			Match: map[string]string{
				mesh_proto.ServiceTag:  name,
				mesh_proto.ProtocolTag: protocol,
			},
		}
	}
}

var _ = Describe("Inspect WS", func() {

	type testCase struct {
		path       string
		goldenFile string
		resources  []core_model.Resource
	}

	DescribeTable("should return valid response",
		func(given testCase) {
			// setup
			resourceStore := memory.NewStore()
			metrics, err := metrics.NewMetrics("Standalone")
			Expect(err).ToNot(HaveOccurred())

			core.Now = func() time.Time { return time.Time{} }

			rm := manager.NewResourceManager(resourceStore)
			for _, resource := range given.resources {
				err = rm.Create(context.Background(), resource,
					store.CreateBy(core_model.MetaToResourceKey(resource.GetMeta())))
				Expect(err).ToNot(HaveOccurred())
			}

			apiServer := createTestApiServer(resourceStore, config.DefaultApiServerConfig(), true, metrics)

			stop := make(chan struct{})
			defer close(stop)
			go func() {
				defer GinkgoRecover()
				err := apiServer.Start(stop)
				Expect(err).ToNot(HaveOccurred())
			}()

			// when
			var resp *http.Response
			Eventually(func() error {
				r, err := http.Get((&url.URL{
					Scheme: "http",
					Host:   apiServer.Address(),
					Path:   given.path,
				}).String())
				resp = r
				return err
			}, "3s").ShouldNot(HaveOccurred())

			// then
			bytes, err := io.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytes).To(matchers.MatchGoldenJSON(path.Join("testdata", given.goldenFile)))
		},
		Entry("inspect dataplane", testCase{
			path:       "/meshes/default/dataplanes/backend-1/policies",
			goldenFile: "inspect_dataplane.json",
			resources: []core_model.Resource{
				newMesh("default"),
				newDataplane().
					meta("backend-1", "default").
					inbound("backend", "192.168.0.1", 80, 81).
					outbound("redis", "192.168.0.2", 8080).
					outbound("gateway", "192.168.0.3", 8080).
					outbound("postgres", "192.168.0.4", 8080).
					outbound("web", "192.168.0.2", 8080).
					build(),
				&core_mesh.TrafficPermissionResource{
					Meta: &test_model.ResourceMeta{Name: "tp-1", Mesh: "default"},
					Spec: &mesh_proto.TrafficPermission{
						Sources:      anyService(),
						Destinations: anyService(),
					},
				},
				&core_mesh.FaultInjectionResource{
					Meta: &test_model.ResourceMeta{Name: "fi-1", Mesh: "default"},
					Spec: &mesh_proto.FaultInjection{
						Sources: anyService(),
						Destinations: selectors{
							serviceSelector("backend", "http"),
						},
						Conf: &mesh_proto.FaultInjection_Conf{
							Delay: &mesh_proto.FaultInjection_Conf_Delay{
								Value:      durationpb.New(5 * time.Second),
								Percentage: util_proto.Double(90),
							},
						},
					},
				},
				&core_mesh.FaultInjectionResource{
					Meta: &test_model.ResourceMeta{Name: "fi-2", Mesh: "default"},
					Spec: &mesh_proto.FaultInjection{
						Sources: anyService(),
						Destinations: selectors{
							serviceSelector("backend", "http"),
						},
						Conf: &mesh_proto.FaultInjection_Conf{
							Abort: &mesh_proto.FaultInjection_Conf_Abort{
								HttpStatus: util_proto.UInt32(500),
								Percentage: util_proto.Double(80),
							},
						},
					},
				},
				&core_mesh.TimeoutResource{
					Meta: &test_model.ResourceMeta{Name: "t-1", Mesh: "default"},
					Spec: &mesh_proto.Timeout{
						Sources: anyService(),
						Destinations: selectors{
							serviceSelector("redis", ""),
						},
						Conf: samples.Timeout.Conf,
					},
				},
				&core_mesh.HealthCheckResource{
					Meta: &test_model.ResourceMeta{Name: "hc-1", Mesh: "default"},
					Spec: &mesh_proto.HealthCheck{
						Sources: selectors{
							serviceSelector("backend", ""),
						},
						Destinations: anyService(),
						Conf:         samples.HealthCheck.Conf,
					},
				},
			},
		}),
		Entry("inspect traffic permission", testCase{
			path:       "/meshes/default/traffic-permissions/tp-1/dataplanes",
			goldenFile: "inspect_traffic-permission.json",
			resources: []core_model.Resource{
				newMesh("default"),
				&core_mesh.TrafficPermissionResource{
					Meta: &test_model.ResourceMeta{Name: "tp-1", Mesh: "default"},
					Spec: &mesh_proto.TrafficPermission{
						Sources: anyService(),
						Destinations: selectors{
							serviceSelector("backend", "http"),
							serviceSelector("redis", "http"),
							serviceSelector("gateway", "http"),
						},
					},
				},
				newDataplane().
					meta("backend-1", "default").
					inbound("backend", "192.168.0.1", 80, 81).
					outbound("redis", "192.168.0.2", 8080).
					outbound("gateway", "192.168.0.3", 8080).
					build(),
				newDataplane().
					meta("redis-1", "default").
					inbound("redis", "192.168.0.1", 80, 81).
					outbound("backend", "192.168.0.2", 8080).
					outbound("gateway", "192.168.0.3", 8080).
					build(),
				newDataplane().
					meta("gateway-1", "default").
					inbound("gateway", "192.168.0.1", 80, 81).
					outbound("backend", "192.168.0.2", 8080).
					outbound("redis", "192.168.0.3", 8080).
					build(),
				newDataplane(). // not matched by TrafficPermission
						meta("web-1", "default").
						inbound("web", "192.168.0.1", 80, 81).
						build(),
			},
		}),
		Entry("inspect fault injection", testCase{
			path:       "/meshes/mesh-1/fault-injections/fi-1/dataplanes",
			goldenFile: "inspect_fault-injection.json",
			resources: []core_model.Resource{
				newMesh("mesh-1"),
				&core_mesh.FaultInjectionResource{
					Meta: &test_model.ResourceMeta{Name: "fi-1", Mesh: "mesh-1"},
					Spec: &mesh_proto.FaultInjection{
						Sources: anyService(),
						Destinations: selectors{
							serviceSelector("backend", "http"),
							serviceSelector("redis", "http"),
							serviceSelector("gateway", "http"),
						},
						Conf: samples.FaultInjection.Conf,
					},
				},
				newDataplane().
					meta("backend-redis-1", "mesh-1").
					inbound("backend", "192.168.0.1", 80, 81).
					inbound("redis", "192.168.0.2", 80, 81).
					build(),
				newDataplane().
					meta("gateway-1", "mesh-1").
					inbound("gateway", "192.168.0.1", 80, 81).
					outbound("backend", "192.168.0.2", 8080).
					outbound("redis", "192.168.0.3", 8080).
					build(),
				newDataplane(). // not matched by FaultInjection
						meta("web-1", "mesh-1").
						inbound("web", "192.168.0.1", 80, 81).
						build(),
			},
		}),
		Entry("inspect rate limit", testCase{
			path:       "/meshes/mesh-1/rate-limits/rl-1/dataplanes",
			goldenFile: "inspect_rate-limit.json",
			resources: []core_model.Resource{
				newMesh("mesh-1"),
				&core_mesh.RateLimitResource{
					Meta: &test_model.ResourceMeta{Name: "rl-1", Mesh: "mesh-1"},
					Spec: &mesh_proto.RateLimit{
						Sources: anyService(),
						Destinations: selectors{
							serviceSelector("backend", "http"),
							serviceSelector("redis", "http"),
							serviceSelector("gateway", "http"),
							serviceSelector("es", ""),
						},
						Conf: samples.RateLimit.Conf,
					},
				},
				newDataplane().
					meta("gateway-1", "mesh-1").
					inbound("gateway", "192.168.0.1", 80, 81).
					outbound("backend", "192.168.0.2", 8080).
					outbound("redis", "192.168.0.3", 8080).
					outbound("es", "192.168.0.4", 8080).
					build(),
				newDataplane(). // not matched by RateLimit
						meta("web-1", "mesh-1").
						inbound("web", "192.168.0.1", 80, 81).
						build(),
				&core_mesh.ExternalServiceResource{
					Meta: &test_model.ResourceMeta{Name: "es-1", Mesh: "mesh-1"},
					Spec: &mesh_proto.ExternalService{
						Networking: &mesh_proto.ExternalService_Networking{Address: "2.2.2.2:80"},
						Tags: map[string]string{
							mesh_proto.ServiceTag: "es",
						},
					},
				},
			},
		}),
		Entry("inspect traffic log", testCase{
			path:       "/meshes/mesh-1/traffic-logs/tl-1/dataplanes",
			goldenFile: "inspect_traffic-log.json",
			resources: []core_model.Resource{
				newMesh("mesh-1"),
				&core_mesh.TrafficLogResource{
					Meta: &test_model.ResourceMeta{Name: "tl-1", Mesh: "mesh-1"},
					Spec: &mesh_proto.TrafficLog{
						Sources: anyService(),
						Destinations: selectors{
							serviceSelector("redis", ""),
							serviceSelector("gateway", ""),
						},
						Conf: samples.TrafficLog.Conf,
					},
				},
				newDataplane().
					meta("backend-1", "mesh-1").
					inbound("backend", "192.168.0.1", 80, 81).
					outbound("redis", "192.168.0.2", 8080).
					outbound("gateway", "192.168.0.3", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
				newDataplane().
					meta("redis-1", "mesh-1").
					inbound("redis", "192.168.0.1", 80, 81).
					outbound("gateway", "192.168.0.2", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
			},
		}),
		Entry("inspect health check", testCase{
			path:       "/meshes/mesh-1/health-checks/hc-1/dataplanes",
			goldenFile: "inspect_health-check.json",
			resources: []core_model.Resource{
				newMesh("mesh-1"),
				&core_mesh.HealthCheckResource{
					Meta: &test_model.ResourceMeta{Name: "hc-1", Mesh: "mesh-1"},
					Spec: &mesh_proto.HealthCheck{
						Sources: anyService(),
						Destinations: selectors{
							serviceSelector("redis", ""),
							serviceSelector("gateway", ""),
						},
						Conf: samples.HealthCheck.Conf,
					},
				},
				newDataplane().
					meta("backend-1", "mesh-1").
					inbound("backend", "192.168.0.1", 80, 81).
					outbound("redis", "192.168.0.2", 8080).
					outbound("gateway", "192.168.0.3", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
				newDataplane().
					meta("redis-1", "mesh-1").
					inbound("redis", "192.168.0.1", 80, 81).
					outbound("gateway", "192.168.0.2", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
			},
		}),
		Entry("inspect circuit breaker", testCase{
			path:       "/meshes/mesh-1/circuit-breakers/cb-1/dataplanes",
			goldenFile: "inspect_circuit-breaker.json",
			resources: []core_model.Resource{
				newMesh("mesh-1"),
				&core_mesh.CircuitBreakerResource{
					Meta: &test_model.ResourceMeta{Name: "cb-1", Mesh: "mesh-1"},
					Spec: &mesh_proto.CircuitBreaker{
						Sources: anyService(),
						Destinations: selectors{
							serviceSelector("redis", ""),
							serviceSelector("gateway", ""),
						},
						Conf: samples.CircuitBreaker.Conf,
					},
				},
				newDataplane().
					meta("backend-1", "mesh-1").
					inbound("backend", "192.168.0.1", 80, 81).
					outbound("redis", "192.168.0.2", 8080).
					outbound("gateway", "192.168.0.3", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
				newDataplane().
					meta("redis-1", "mesh-1").
					inbound("redis", "192.168.0.1", 80, 81).
					outbound("gateway", "192.168.0.2", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
			},
		}),
		Entry("inspect retry", testCase{
			path:       "/meshes/mesh-1/retries/r-1/dataplanes",
			goldenFile: "inspect_retry.json",
			resources: []core_model.Resource{
				newMesh("mesh-1"),
				&core_mesh.RetryResource{
					Meta: &test_model.ResourceMeta{Name: "r-1", Mesh: "mesh-1"},
					Spec: &mesh_proto.Retry{
						Sources: anyService(),
						Destinations: selectors{
							serviceSelector("redis", ""),
							serviceSelector("gateway", ""),
						},
						Conf: samples.Retry.Conf,
					},
				},
				newDataplane().
					meta("backend-1", "mesh-1").
					inbound("backend", "192.168.0.1", 80, 81).
					outbound("redis", "192.168.0.2", 8080).
					outbound("gateway", "192.168.0.3", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
				newDataplane().
					meta("redis-1", "mesh-1").
					inbound("redis", "192.168.0.1", 80, 81).
					outbound("gateway", "192.168.0.2", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
			},
		}),
		Entry("inspect timeout", testCase{
			path:       "/meshes/mesh-1/timeouts/t-1/dataplanes",
			goldenFile: "inspect_timeout.json",
			resources: []core_model.Resource{
				newMesh("mesh-1"),
				&core_mesh.TimeoutResource{
					Meta: &test_model.ResourceMeta{Name: "t-1", Mesh: "mesh-1"},
					Spec: &mesh_proto.Timeout{
						Sources: anyService(),
						Destinations: selectors{
							serviceSelector("redis", ""),
							serviceSelector("gateway", ""),
						},
						Conf: samples.Timeout.Conf,
					},
				},
				newDataplane().
					meta("backend-1", "mesh-1").
					inbound("backend", "192.168.0.1", 80, 81).
					outbound("redis", "192.168.0.2", 8080).
					outbound("gateway", "192.168.0.3", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
				newDataplane().
					meta("redis-1", "mesh-1").
					inbound("redis", "192.168.0.1", 80, 81).
					outbound("gateway", "192.168.0.2", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
			},
		}),
		Entry("inspect traffic route", testCase{
			path:       "/meshes/mesh-1/traffic-routes/t-1/dataplanes",
			goldenFile: "inspect_traffic-route.json",
			resources: []core_model.Resource{
				newMesh("mesh-1"),
				&core_mesh.TrafficRouteResource{
					Meta: &test_model.ResourceMeta{Name: "t-1", Mesh: "mesh-1"},
					Spec: &mesh_proto.TrafficRoute{
						Sources:      anyService(),
						Destinations: anyService(),
						Conf:         samples.TrafficRoute.Conf,
					},
				},
				newDataplane().
					meta("backend-1", "mesh-1").
					inbound("backend", "192.168.0.1", 80, 81).
					outbound("redis", "192.168.0.2", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
				newDataplane().
					meta("redis-1", "mesh-1").
					inbound("redis", "192.168.0.1", 80, 81).
					outbound("gateway", "192.168.0.2", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
			},
		}),
		Entry("inspect traffic trace", testCase{
			path:       "/meshes/mesh-1/traffic-traces/tt-1/dataplanes",
			goldenFile: "inspect_traffic-trace.json",
			resources: []core_model.Resource{
				newMesh("mesh-1"),
				&core_mesh.TrafficTraceResource{
					Meta: &test_model.ResourceMeta{Name: "tt-1", Mesh: "mesh-1"},
					Spec: &mesh_proto.TrafficTrace{
						Selectors: anyService(),
						Conf:      samples.TrafficTrace.Conf,
					},
				},
				newDataplane().
					meta("backend-1", "mesh-1").
					inbound("backend", "192.168.0.1", 80, 81).
					outbound("redis", "192.168.0.2", 8080).
					outbound("gateway", "192.168.0.3", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
				newDataplane().
					meta("redis-1", "mesh-1").
					inbound("redis", "192.168.0.1", 80, 81).
					outbound("gateway", "192.168.0.2", 8080).
					outbound("web", "192.168.0.4", 8080).
					build(),
				newDataplane().
					meta("web-1", "mesh-1").
					inbound("web", "192.168.0.1", 80, 81).
					outbound("gateway", "192.168.0.2", 8080).
					outbound("backend", "192.168.0.4", 8080).
					build(),
			},
		}),
	)

	It("should change response if state changed", func() {
		// setup
		resourceStore := memory.NewStore()
		metrics, err := metrics.NewMetrics("Standalone")
		Expect(err).ToNot(HaveOccurred())
		rm := manager.NewResourceManager(resourceStore)
		apiServer := createTestApiServer(resourceStore, config.DefaultApiServerConfig(), true, metrics)

		stop := make(chan struct{})
		defer close(stop)
		go func() {
			defer GinkgoRecover()
			err := apiServer.Start(stop)
			Expect(err).ToNot(HaveOccurred())
		}()

		// when init the state
		// TrafficPermission that selects 2 DPPs
		initState := []core_model.Resource{
			newMesh("default"),
			&core_mesh.TrafficPermissionResource{
				Meta: &test_model.ResourceMeta{Name: "tp-1", Mesh: "default"},
				Spec: &mesh_proto.TrafficPermission{
					Sources: anyService(),
					Destinations: selectors{
						serviceSelector("backend", "http"),
						serviceSelector("redis", "http"),
					},
				},
			},
			newDataplane().
				meta("backend-1", "default").
				inbound("backend", "192.168.0.1", 80, 81).
				outbound("redis", "192.168.0.2", 8080).
				outbound("gateway", "192.168.0.3", 8080).
				build(),
			newDataplane().
				meta("redis-1", "default").
				inbound("redis", "192.168.0.1", 80, 81).
				outbound("backend", "192.168.0.2", 8080).
				outbound("gateway", "192.168.0.3", 8080).
				build(),
		}
		for _, resource := range initState {
			err = rm.Create(context.Background(), resource,
				store.CreateBy(core_model.MetaToResourceKey(resource.GetMeta())))
			Expect(err).ToNot(HaveOccurred())
		}

		// then
		var resp *http.Response
		Eventually(func() error {
			r, err := http.Get((&url.URL{
				Scheme: "http",
				Host:   apiServer.Address(),
				Path:   "/meshes/default/traffic-permissions/tp-1/dataplanes",
			}).String())
			resp = r
			return err
		}, "3s").ShouldNot(HaveOccurred())
		bytes, err := io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).To(matchers.MatchGoldenJSON(path.Join("testdata", "inspect_changed_state_before.json")))

		// when change the state
		err = rm.Delete(context.Background(), core_mesh.NewDataplaneResource(), store.DeleteByKey("backend-1", "default"))
		Expect(err).ToNot(HaveOccurred())

		// then
		Eventually(func() error {
			r, err := http.Get((&url.URL{
				Scheme: "http",
				Host:   apiServer.Address(),
				Path:   "/meshes/default/traffic-permissions/tp-1/dataplanes",
			}).String())
			resp = r
			return err
		}, "3s").ShouldNot(HaveOccurred())
		bytes, err = io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).To(matchers.MatchGoldenJSON(path.Join("testdata", "inspect_changed_state_after.json")))
	})
})
