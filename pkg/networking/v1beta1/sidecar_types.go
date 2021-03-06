// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Sidecar describes the configuration of the sidecar proxy that mediates
// inbound and outbound communication to the workload instance it is attached to. By
// default, Istio will program all sidecar proxies in the mesh with the
// necessary configuration required to reach every workload instance in the mesh, as
// well as accept traffic on all the ports associated with the
// workload. The `SidecarSpec` configuration provides a way to fine tune the set of
// ports, protocols that the proxy will accept when forwarding traffic to
// and from the workload. In addition, it is possible to restrict the set
// of services that the proxy can reach when forwarding outbound traffic
// from workload instances.
//
// Services and configuration in a mesh are organized into one or more
// namespaces (e.g., a Kubernetes namespace or a CF org/space). A `SidecarSpec`
// configuration in a namespace will apply to one or more workload instances in the same
// namespace, selected using the `workloadSelector` field. In the absence of a
// `workloadSelector`, it will apply to all workload instances in the same
// namespace. When determining the `SidecarSpec` configuration to be applied to a
// workload instance, preference will be given to the resource with a
// `workloadSelector` that selects this workload instance, over a `SidecarSpec` configuration
// without any `workloadSelector`.
//
// NOTE 1: *_Each namespace can have only one `SidecarSpec` configuration without any
// `workloadSelector`_*. The behavior of the system is undefined if more
// than one selector-less `SidecarSpec` configurations exist in a given namespace. The
// behavior of the system is undefined if two or more `SidecarSpec` configurations
// with a `workloadSelector` select the same workload instance.
//
// NOTE 2: *_A `SidecarSpec` configuration in the `MeshConfig`
// [root namespace](https://istio.io/docs/reference/config/istio.mesh.v1alpha1/#MeshConfig)
// will be applied by default to all namespaces without a `SidecarSpec`
// configuration_*. This global default `SidecarSpec` configuration should not have
// any `workloadSelector`.
//
// The example below declares a global default `SidecarSpec` configuration in the
// root namespace called `istio-config`, that configures sidecars in
// all namespaces to allow egress traffic only to other workloads in
// the same namespace, and to services in the `istio-system` namespace.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: SidecarSpec
// metadata:
//   name: default
//   namespace: istio-config
// spec:
//   egress:
//   - hosts:
//     - "./*"
//     - "istio-system/*"
//```
//
// The example below declares a `SidecarSpec` configuration in the `prod-us1`
// namespace that overrides the global default defined above, and
// configures the sidecars in the namespace to allow egress traffic to
// public services in the `prod-us1`, `prod-apis`, and the `istio-system`
// namespaces.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: SidecarSpec
// metadata:
//   name: default
//   namespace: prod-us1
// spec:
//   egress:
//   - hosts:
//     - "prod-us1/*"
//     - "prod-apis/*"
//     - "istio-system/*"
// ```
//
// The example below declares a `SidecarSpec` configuration in the `prod-us1` namespace
// that accepts inbound HTTP traffic on port 9080 and forwards
// it to the attached workload instance listening on a Unix domain socket. In the
// egress direction, in addition to the `istio-system` namespace, the sidecar
// proxies only HTTP traffic bound for port 9080 for services in the
// `prod-us1` namespace.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: SidecarSpec
// metadata:
//   name: default
//   namespace: prod-us1
// spec:
//   ingress:
//   - port:
//       number: 9080
//       protocol: HTTP
//       name: somename
//     defaultEndpoint: unix:///var/run/someuds.sock
//   egress:
//   - port:
//       number: 9080
//       protocol: HTTP
//       name: egresshttp
//     hosts:
//     - "prod-us1/*"
//   - hosts:
//     - "istio-system/*"
// ```
//
// If the workload is deployed without IPTables-based traffic capture, the
// `SidecarSpec` configuration is the only way to configure the ports on the proxy
// attached to the workload instance. The following example declares a `SidecarSpec`
// configuration in the `prod-us1` namespace for all pods with labels
// `app: productpage` belonging to the `productpage.prod-us1` service. Assuming
// that these pods are deployed without IPtable rules (i.e. the `istio-init`
// container) and the proxy metadata `ISTIO_META_INTERCEPTION_MODE` is set to
// `NONE`, the specification, below, allows such pods to receive HTTP traffic
// on port 9080 and forward it to the application listening on
// `127.0.0.1:8080`. It also allows the application to communicate with a
// backing MySQL database on `127.0.0.1:3306`, that then gets proxied to the
// externally hosted MySQL service at `mysql.foo.com:3306`.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: SidecarSpec
// metadata:
//   name: no-ip-tables
//   namespace: prod-us1
// spec:
//   workloadSelector:
//     labels:
//       app: productpage
//   ingress:
//   - port:
//       number: 9080 # binds to proxy_instance_ip:9080 (0.0.0.0:9080, if no unicast IP is available for the instance)
//       protocol: HTTP
//       name: somename
//     defaultEndpoint: 127.0.0.1:8080
//     captureMode: NONE # not needed if metadata is set for entire proxy
//   egress:
//   - port:
//       number: 3306
//       protocol: MYSQL
//       name: egressmysql
//     captureMode: NONE # not needed if metadata is set for entire proxy
//     bind: 127.0.0.1
//     hosts:
//     - "*/mysql.foo.com"
// ```
//
// And the associated service entry for routing to `mysql.foo.com:3306`
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: ServiceEntry
// metadata:
//   name: external-svc-mysql
//   namespace: ns1
// spec:
//   hosts:
//   - mysql.foo.com
//   ports:
//   - number: 3306
//     name: mysql
//     protocol: MYSQL
//   location: MESH_EXTERNAL
//   resolution: DNS
// ```
//
// It is also possible to mix and match traffic capture modes in a single
// proxy. For example, consider a setup where internal services are on the
// `192.168.0.0/16` subnet. So, IP tables are setup on the VM to capture all
// outbound traffic on `192.168.0.0/16` subnet. Assume that the VM has an
// additional network interface on `172.16.0.0/16` subnet for inbound
// traffic. The following `SidecarSpec` configuration allows the VM to expose a
// listener on `172.16.1.32:80` (the VM's IP) for traffic arriving from the
// `172.16.0.0/16` subnet. Note that in this scenario, the
// `ISTIO_META_INTERCEPTION_MODE` metadata on the proxy in the VM should
// contain `REDIRECT` or `TPROXY` as its value, implying that IP tables
// based traffic capture is active.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: SidecarSpec
// metadata:
//   name: partial-ip-tables
//   namespace: prod-us1
// spec:
//   workloadSelector:
//     labels:
//       app: productpage
//   ingress:
//   - bind: 172.16.1.32
//     port:
//       number: 80 # binds to 172.16.1.32:80
//       protocol: HTTP
//       name: somename
//     defaultEndpoint: 127.0.0.1:8080
//     captureMode: NONE
//   egress:
//     # use the system detected defaults
//     # sets up configuration to handle outbound traffic to services
//     # in 192.168.0.0/16 subnet, based on information provided by the
//     # service registry
//   - captureMode: IPTABLES
//     hosts:
//     - "*/*"
// ```
type Sidecar struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SidecarSpec `json:"spec"`
}

// SidecarSpec describes the configuration of the sidecar proxy that mediates
// inbound and outbound communication of the workload instance to which it is
// attached.
type SidecarSpec struct {
	// Criteria used to select the specific set of pods/VMs on which this
	// `SidecarSpec` configuration should be applied. If omitted, the `SidecarSpec`
	// configuration will be applied to all workload instances in the same namespace.
	WorkloadSelector *WorkloadSelector `json:"workloadSelector,omitempty"`
	// Ingress specifies the configuration of the sidecar for processing
	// inbound traffic to the attached workload instance. If omitted, Istio will
	// automatically configure the sidecar based on the information about the workload
	// obtained from the orchestration platform (e.g., exposed ports, services,
	// etc.). If specified, inbound ports are configured if and only if the
	// workload instance is associated with a service.
	Ingress []*IstioIngressListener `json:"ingress,omitempty"`
	// Egress specifies the configuration of the sidecar for processing
	// outbound traffic from the attached workload instance to other services in the
	// mesh.
	Egress []*IstioEgressListener `json:"egress"`
	// This allows to configure the outbound traffic policy.
	// If your application uses one or more external
	// services that are not known apriori, setting the policy to `ALLOW_ANY`
	// will cause the sidecars to route any unknown traffic originating from
	// the application to its requested destination.
	OutboundTrafficPolicy *OutboundTrafficPolicy `json:"outboundTrafficPolicy,omitempty"`
}

// `OutboundTrafficPolicy` sets the default behavior of the sidecar for
// handling outbound traffic from the application.
// If your application uses one or more external
// services that are not known apriori, setting the policy to `ALLOW_ANY`
// will cause the sidecars to route any unknown traffic originating from
// the application to its requested destination.  Users are strongly
// encouraged to use `ServiceEntry` configurations to explicitly declare any external
// dependencies, instead of using `ALLOW_ANY`, so that traffic to these
// services can be monitored.
type OutboundTrafficPolicy struct {
	Mode *OutboundTrafficPolicyMode `json:"mode,omitempty"`
}

type OutboundTrafficPolicyMode string

const (
	// Outbound traffic will be restricted to services defined in the
	// service registry as well as those defined through `ServiceEntry` configurations.
	OutboundTrafficPolicyRegistryOnly OutboundTrafficPolicyMode = "REGISTRY_ONLY"
	// Outbound traffic to unknown destinations will be allowed, in case
	// there are no services or `ServiceEntry` configurations for the destination port.
	OutboundTrafficPolicyAllowAny OutboundTrafficPolicyMode = "ALLOW_ANY"
)

// IstioIngressListener specifies the properties of an inbound
// traffic listener on the sidecar proxy attached to a workload instance.
type IstioIngressListener struct {
	// The port associated with the listener.
	Port *Port `json:"port"`
	// The IP to which the listener should be bound. Must be in the
	// format `x.x.x.x`. Unix domain socket addresses are not allowed in
	// the bind field for ingress listeners. If omitted, Istio will
	// automatically configure the defaults based on imported services
	// and the workload instances to which this configuration is applied
	// to.
	Bind string `json:"bind,omitempty"`
	// The captureMode option dictates how traffic to the listener is
	// expected to be captured (or not).
	CaptureMode CaptureMode `json:"captureMode,omitempty"`
	// The loopback IP endpoint or Unix domain socket to which
	// traffic should be forwarded to. This configuration can be used to
	// redirect traffic arriving at the bind `IP:Port` on the sidecar to a `localhost:port`
	// or Unix domain socket where the application workload instance is listening for
	// connections. Format should be `127.0.0.1:PORT` or `unix:///path/to/socket`
	DefaultEndpoint string `json:"defaultEndpoint"`
}

// IstioEgressListener specifies the properties of an outbound traffic
// listener on the sidecar proxy attached to a workload instance.
type IstioEgressListener struct {
	// The port associated with the listener. If using Unix domain socket,
	// use 0 as the port number, with a valid protocol. The port if
	// specified, will be used as the default destination port associated
	// with the imported hosts. If the port is omitted, Istio will infer the
	// listener ports based on the imported hosts. Note that when multiple
	// egress listeners are specified, where one or more listeners have
	// specific ports while others have no port, the hosts exposed on a
	// listener port will be based on the listener with the most specific
	// port.
	Port *Port `json:"port,omitempty"`
	// The IP or the Unix domain socket to which the listener should be bound
	// to. Port MUST be specified if bind is not empty. Format: `x.x.x.x` or
	// `unix:///path/to/uds` or `unix://@foobar` (Linux abstract namespace). If
	// omitted, Istio will automatically configure the defaults based on imported
	// services, the workload instances to which this configuration is applied to and
	// the captureMode. If captureMode is `NONE`, bind will default to
	// 127.0.0.1.
	Bind string `json:"bind,omitempty"`
	// When the bind address is an IP, the captureMode option dictates
	// how traffic to the listener is expected to be captured (or not).
	// captureMode must be DEFAULT or `NONE` for Unix domain socket binds.
	CaptureMode CaptureMode `json:"captureMode,omitempty"`
	// One or more service hosts exposed by the listener
	// in `namespace/dnsName` format. Services in the specified namespace
	// matching `dnsName` will be exposed.
	// The corresponding service can be a service in the service registry
	// (e.g., a Kubernetes or cloud foundry service) or a service specified
	// using a `ServiceEntry` or `VirtualService` configuration. Any
	// associated `DestinationRule` in the same namespace will also be used.
	//
	// The `dnsName` should be specified using FQDN format, optionally including
	// a wildcard character in the left-most component (e.g., `prod/*.example.com`).
	// Set the `dnsName` to `*` to select all services from the specified namespace
	// (e.g., `prod/*`).
	//
	// The `namespace` can be set to `*`, `.`, or `~`, representing any, the current,
	// or no namespace, respectively. For example, `*/foo.example.com` selects the
	// service from any available namespace while `./foo.example.com` only selects
	// the service from the namespace of the sidecar. If a host is set to `*/*`,
	// Istio will configure the sidecar to be able to reach every service in the
	// mesh that is exported to the sidecar's namespace. The value `~/*` can be used
	// to completely trim the configuration for sidecars that simply receive traffic
	// and respond, but make no outbound connections of their own.
	//
	// NOTE: Only services and configuration artifacts exported to the sidecar's
	// namespace (e.g., `exportTo` value of `*`) can be referenced.
	// Private configurations (e.g., `exportTo` set to `.`) will
	// not be available. Refer to the `exportTo` setting in `VirtualService`,
	// `DestinationRule`, and `ServiceEntry` configurations for details.
	//
	// **WARNING:** The list of egress hosts in a `SidecarSpec` must also include
	// the Mixer control plane services if they are enabled. Envoy will not
	// be able to reach them otherwise. For example, add host
	// `istio-system/istio-telemetry.istio-system.svc.cluster.local` if telemetry
	// is enabled, `istio-system/istio-policy.istio-system.svc.cluster.local` if
	// policy is enabled, or add `istio-system/*` to allow all services in the
	// `istio-system` namespace. This requirement is temporary and will be removed
	// in a future Istio release.
	Hosts []string `json:"hosts"`
}

// WorkloadSelector specifies the criteria used to determine if the `Gateway`,
// `SidecarSpec`, or `EnvoyFilter` configuration can be applied to a proxy. The matching criteria
// includes the metadata associated with a proxy, workload instance info such as
// labels attached to the pod/VM, or any other info that the proxy provides
// to Istio during the initial handshake. If multiple conditions are
// specified, all conditions need to match in order for the workload instance to be
// selected. Currently, only label based selection mechanism is supported.
type WorkloadSelector struct {
	// One or more labels that indicate a specific set of pods/VMs
	// on which this `SidecarSpec` configuration should be applied. The scope of
	// label search is restricted to the configuration namespace in which the
	// the resource is present.
	Labels map[string]string `json:"labels"`
}

// CaptureMode describes how traffic to a listener is expected to be
// captured. Applicable only when the listener is bound to an IP.
type CaptureMode string

const (
	// The default capture mode defined by the environment.
	CaptureModeDefault CaptureMode = "DEFAULT"
	// Capture traffic using IPtables redirection.
	CaptureModeIPTables CaptureMode = "IPTABLES"
	// No traffic capture. When used in an egress listener, the application is
	// expected to explicitly communicate with the listener port or Unix
	// domain socket. When used in an ingress listener, care needs to be taken
	// to ensure that the listener port is not in use by other processes on
	// the host.
	CaptureModeNone CaptureMode = "NONE"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SidecarList is a list of Sidecar resources
type SidecarList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Sidecar `json:"items"`
}
