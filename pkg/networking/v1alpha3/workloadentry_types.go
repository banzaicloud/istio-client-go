// Copyright © 2020 Banzai Cloud
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

package v1alpha3

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	istioApi "github.com/banzaicloud/istio-client-go/pkg/networking/v1alpha3/istioapi"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// WorkloadEntry
type WorkloadEntry struct {
	v1.TypeMeta   `json:",inline"`
	v1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the implementation of this definition.
	Spec   WorkloadEntrySpec    `json:"spec"`
	Status istioApi.IstioStatus `json:"status"`
}

// `WorkloadEntry` enables operators to describe the properties of a
// single non-Kubernetes workload such as a VM or a bare metal server
// as it is are onboarded into the mesh. A `WorkloadEntry` must be
// accompanied by an Istio `ServiceEntry` that selects the workload
// through the appropriate labels and provides the service definition
// for a `MESH_INTERNAL` service (hostnames, port properties, etc.). A
// `ServiceEntry` object can select multiple workload entries as well
// as Kubernetes pods based on the label selector specified in the
// service entry.
//
// When a workload connects to `istiod`, the status field in the
// custom resource will be updated to indicate the health of the
// workload along with other details, similar to how Kubernetes
// updates the status of a pod.
//
// The following example declares a workload entry representing a
// VM for the `details.bookinfo.com` service. This VM has
// sidecar installed and bootstrapped using the `details-legacy`
// service account. The sidecar receives HTTP traffic on port 80
// (wrapped in istio mutual TLS) and forwards it to the application on
// the localhost on the same port.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: WorkloadEntry
// metadata:
//   name: details-svc
// spec:
//   # use of the service account indicates that the workload has a
//   # sidecar proxy bootstrapped with this service account. Pods with
//   # sidecars will automatically communicate with the workload using
//   # istio mutual TLS.
//   serviceAccount: details-legacy
//   address: 2.2.2.2
//   labels:
//     app: details-legacy
//     instance-id: vm1
//   # ports if not specified will be the same as service ports
// ```
//
// and the associated service entry
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: ServiceEntry
// metadata:
//   name: details-svc
// spec:
//   hosts:
//   - details.bookinfo.com
//   location: MESH_INTERNAL
//   ports:
//   - number: 80
//     name: http
//     protocol: HTTP
//   resolution: STATIC
//   workloadSelector:
//     labels:
//       app: details-legacy
// ```
//
// The following example declares the same VM workload using
// its fully qualified DNS name. The service entry's resolution
// mode should be changed to DNS to indicate that the client-side
// sidecars should dynamically resolve the DNS name at runtime before
// forwarding the request.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: WorkloadEntry
// metadata:
//   name: details-svc
// spec:
//   # use of the service account indicates that the workload has a
//   # sidecar proxy bootstrapped with this service account. Pods with
//   # sidecars will automatically communicate with the workload using
//   # istio mutual TLS.
//   serviceAccount: details-legacy
//   address: vm1.vpc01.corp.net
//   labels:
//     app: details-legacy
//     instance-id: vm1
//   # ports if not specified will be the same as service ports
// ```
//
// and the associated service entry
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: ServiceEntry
// metadata:
//   name: details-svc
// spec:
//   hosts:
//   - details.bookinfo.com
//   location: MESH_INTERNAL
//   ports:
//   - number: 80
//     name: http
//     protocol: HTTP
//   resolution: DNS
//   workloadSelector:
//     labels:
//       app: details-legacy
// ```
type WorkloadEntrySpec struct {
	// Address associated with the network endpoint without the
	// port.  Domain names can be used if and only if the resolution is set
	// to DNS, and must be fully-qualified without wildcards. Use the form
	// unix:///absolute/path/to/socket for Unix domain socket endpoints.
	Address string `json:"address"`
	// Set of ports associated with the endpoint. The ports must be
	// associated with a port name that was declared as part of the
	// service. Do not use for `unix://` addresses.
	Ports map[string]uint32 `json:"ports,omitempty"`
	// One or more labels associated with the endpoint.
	Labels map[string]string `json:"labels,omitempty"`
	// Network enables Istio to group endpoints resident in the same L3
	// domain/network. All endpoints in the same network are assumed to be
	// directly reachable from one another. When endpoints in different
	// networks cannot reach each other directly, an Istio Gateway can be
	// used to establish connectivity (usually using the
	// `AUTO_PASSTHROUGH` mode in a Gateway Server). This is
	// an advanced configuration used typically for spanning an Istio mesh
	// over multiple clusters.
	Network string `json:"network,omitempty"`
	// The locality associated with the endpoint. A locality corresponds
	// to a failure domain (e.g., country/region/zone). Arbitrary failure
	// domain hierarchies can be represented by separating each
	// encapsulating failure domain by /. For example, the locality of an
	// an endpoint in US, in US-East-1 region, within availability zone
	// az-1, in data center rack r11 can be represented as
	// us/us-east-1/az-1/r11. Istio will configure the sidecar to route to
	// endpoints within the same locality as the sidecar. If none of the
	// endpoints in the locality are available, endpoints parent locality
	// (but within the same network ID) will be chosen. For example, if
	// there are two endpoints in same network (networkID "n1"), say e1
	// with locality us/us-east-1/az-1/r11 and e2 with locality
	// us/us-east-1/az-2/r12, a sidecar from us/us-east-1/az-1/r11 locality
	// will prefer e1 from the same locality over e2 from a different
	// locality. Endpoint e2 could be the IP associated with a gateway
	// (that bridges networks n1 and n2), or the IP associated with a
	// standard service endpoint.
	Locality string `json:"locality,omitempty"`
	// The load balancing weight associated with the endpoint. Endpoints
	// with higher weights will receive proportionally higher traffic.
	Weight uint32 `json:"weight,omitempty"`
	// The service account associated with the workload if a sidecar
	// is present in the workload. The service account must be present
	// in the same namespace as the configuration ( WorkloadEntry or a
	// ServiceEntry)
	ServiceAccount string `json:"serviceAccount,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// WorkloadEntryList is a collection of EnvoyFilters.
type WorkloadEntryList struct {
	v1.TypeMeta `json:",inline"`
	v1.ListMeta `json:"metadata"`
	Items       []WorkloadEntry `json:"items"`
}
