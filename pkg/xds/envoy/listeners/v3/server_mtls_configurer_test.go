package v3_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	test_model "github.com/kumahq/kuma/pkg/test/resources/model"
	util_proto "github.com/kumahq/kuma/pkg/util/proto"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	envoy_common "github.com/kumahq/kuma/pkg/xds/envoy"
	. "github.com/kumahq/kuma/pkg/xds/envoy/listeners"
)

var _ = Describe("ServerMtlsConfigurer", func() {

	type testCase struct {
		listenerName     string
		listenerProtocol core_xds.SocketAddressProtocol
		listenerAddress  string
		listenerPort     uint32
		statsName        string
		clusters         []envoy_common.Cluster
		ctx              xds_context.Context
		expected         string
	}

	DescribeTable("should generate proper Envoy config",
		func(given testCase) {
			// when
			listener, err := NewListenerBuilder(envoy_common.APIV3).
				Configure(InboundListener(given.listenerName, given.listenerAddress, given.listenerPort, given.listenerProtocol)).
				Configure(FilterChain(NewFilterChainBuilder(envoy_common.APIV3).
					Configure(ServerSideMTLS(given.ctx)).
					Configure(TcpProxy(given.statsName, given.clusters...)))).
				Build()
			// then
			Expect(err).ToNot(HaveOccurred())

			// when
			actual, err := util_proto.ToYAML(listener)
			Expect(err).ToNot(HaveOccurred())
			// and
			Expect(actual).To(MatchYAML(given.expected))
		},
		Entry("basic tcp_proxy with mTLS", testCase{
			listenerName:    "inbound:192.168.0.1:8080",
			listenerAddress: "192.168.0.1",
			listenerPort:    8080,
			statsName:       "localhost:8080",
			clusters: []envoy_common.Cluster{envoy_common.NewCluster(
				envoy_common.WithService("localhost:8080"),
				envoy_common.WithWeight(200),
			)},
			ctx: xds_context.Context{
				ControlPlane: &xds_context.ControlPlaneContext{},
				Mesh: xds_context.MeshContext{
					Resource: &core_mesh.MeshResource{
						Meta: &test_model.ResourceMeta{
							Name: "default",
						},
						Spec: &mesh_proto.Mesh{
							Mtls: &mesh_proto.Mesh_Mtls{
								EnabledBackend: "builtin",
								Backends: []*mesh_proto.CertificateAuthorityBackend{
									{
										Name: "builtin",
										Type: "builtin",
									},
								},
							},
						},
					},
				},
			},
			expected: `
            address:
              socketAddress:
                address: 192.168.0.1
                portValue: 8080
            filterChains:
            - filters:
              - name: envoy.filters.network.tcp_proxy
                typedConfig:
                  '@type': type.googleapis.com/envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy
                  cluster: localhost:8080
                  statPrefix: localhost_8080
              transportSocket:
                name: envoy.transport_sockets.tls
                typedConfig:
                  '@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext
                  commonTlsContext:
                    combinedValidationContext:
                      defaultValidationContext:
                        matchSubjectAltNames:
                        - prefix: spiffe://default/
                      validationContextSdsSecretConfig:
                        name: mesh_ca:secret:default
                        sdsConfig:
                          ads: {}
                          resourceApiVersion: V3
                    tlsCertificateSdsSecretConfigs:
                    - name: identity_cert:secret:default
                      sdsConfig:
                        ads: {}
                        resourceApiVersion: V3
                  requireClientCertificate: true
            name: inbound:192.168.0.1:8080
            trafficDirection: INBOUND
`,
		}),
		Entry("basic tcp_proxy with mTLS and Dataplane credentials", testCase{
			listenerName:    "inbound:192.168.0.1:8080",
			listenerAddress: "192.168.0.1",
			listenerPort:    8080,
			statsName:       "localhost:8080",
			clusters: []envoy_common.Cluster{envoy_common.NewCluster(
				envoy_common.WithService("localhost:8080"),
				envoy_common.WithWeight(200),
			)},
			ctx: xds_context.Context{
				ControlPlane: &xds_context.ControlPlaneContext{},
				Mesh: xds_context.MeshContext{
					Resource: &core_mesh.MeshResource{
						Meta: &test_model.ResourceMeta{
							Name: "default",
						},
						Spec: &mesh_proto.Mesh{
							Mtls: &mesh_proto.Mesh_Mtls{
								EnabledBackend: "builtin",
								Backends: []*mesh_proto.CertificateAuthorityBackend{
									{
										Name: "builtin",
										Type: "builtin",
									},
								},
							},
						},
					},
				},
			},
			expected: `
            address:
              socketAddress:
                address: 192.168.0.1
                portValue: 8080
            filterChains:
            - filters:
              - name: envoy.filters.network.tcp_proxy
                typedConfig:
                  '@type': type.googleapis.com/envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy
                  cluster: localhost:8080
                  statPrefix: localhost_8080
              transportSocket:
                name: envoy.transport_sockets.tls
                typedConfig:
                  '@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext
                  commonTlsContext:
                    combinedValidationContext:
                      defaultValidationContext:
                        matchSubjectAltNames:
                        - prefix: spiffe://default/
                      validationContextSdsSecretConfig:
                        name: mesh_ca:secret:default
                        sdsConfig:
                          ads: {}
                          resourceApiVersion: V3
                    tlsCertificateSdsSecretConfigs:
                    - name: identity_cert:secret:default
                      sdsConfig:
                        ads: {}
                        resourceApiVersion: V3
                  requireClientCertificate: true
            name: inbound:192.168.0.1:8080
            trafficDirection: INBOUND`,
		}),
	)
})
