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

	"github.com/banzaicloud/istio-client-go/pkg/common/v1alpha1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Configuration affecting traffic routing. Here are a few terms useful to define
// in the context of traffic routing.
//
// `Service` a unit of application behavior bound to a unique name in a
// service registry. Services consist of multiple network *endpoints*
// implemented by workload instances running on pods, containers, VMs etc.
//
// `Service versions (a.k.a. subsets)` - In a continuous deployment
// scenario, for a given service, there can be distinct subsets of
// instances running different variants of the application binary. These
// variants are not necessarily different API versions. They could be
// iterative changes to the same service, deployed in different
// environments (prod, staging, dev, etc.). Common scenarios where this
// occurs include A/B testing, canary rollouts, etc. The choice of a
// particular version can be decided based on various criterion (headers,
// url, etc.) and/or by weights assigned to each version. Each service has
// a default version consisting of all its instances.
//
// `Source` - A downstream client calling a service.
//
// `Host` - The address used by a client when attempting to connect to a
// service.
//
// `Access model` - Applications address only the destination service
// (Host) without knowledge of individual service versions (subsets). The
// actual choice of the version is determined by the proxy/sidecar, enabling the
// application code to decouple itself from the evolution of dependent
// services.
//
// A `VirtualService` defines a set of traffic routing rules to apply when a host is
// addressed. Each routing rule defines matching criteria for traffic of a specific
// protocol. If the traffic is matched, then it is sent to a named destination service
// (or subset/version of it) defined in the registry.
//
// The source of traffic can also be matched in a routing rule. This allows routing
// to be customized for specific client contexts.
//
// The following example on Kubernetes, routes all HTTP traffic by default to
// pods of the reviews service with label "version: v1". In addition,
// HTTP requests with path starting with /wpcatalog/ or /consumercatalog/ will
// be rewritten to /newcatalog and sent to pods with label "version: v2".
//
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: reviews-route
// spec:
//   hosts:
//   - reviews.prod.svc.cluster.local
//   http:
//   - name: "reviews-v2-routes"
//     match:
//     - uri:
//         prefix: "/wpcatalog"
//     - uri:
//         prefix: "/consumercatalog"
//     rewrite:
//       uri: "/newcatalog"
//     route:
//     - destination:
//         host: reviews.prod.svc.cluster.local
//         subset: v2
//   - name: "reviews-v1-route"
//     route:
//     - destination:
//         host: reviews.prod.svc.cluster.local
//         subset: v1
// ```
//
// A subset/version of a route destination is identified with a reference
// to a named service subset which must be declared in a corresponding
// `DestinationRule`.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: DestinationRule
// metadata:
//   name: reviews-destination
// spec:
//   host: reviews.prod.svc.cluster.local
//   subsets:
//   - name: v1
//     labels:
//       version: v1
//   - name: v2
//     labels:
//       version: v2
// ```
type VirtualService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VirtualServiceSpec `json:"spec"`
}

// Configuration affecting traffic routing.
type VirtualServiceSpec struct {
	// REQUIRED. The destination hosts to which traffic is being sent. Could
	// be a DNS name with wildcard prefix or an IP address.  Depending on the
	// platform, short-names can also be used instead of a FQDN (i.e. has no
	// dots in the name). In such a scenario, the FQDN of the host would be
	// derived based on the underlying platform.
	//
	// A single VirtualService can be used to describe all the traffic
	// properties of the corresponding hosts, including those for multiple
	// HTTP and TCP ports. Alternatively, the traffic properties of a host
	// can be defined using more than one VirtualService, with certain
	// caveats. Refer to the
	// [Operations Guide](https://istio.io/docs/ops/traffic-management/deploy-guidelines/#multiple-virtual-services-and-destination-rules-for-the-same-host)
	// for details.
	//
	// *Note for Kubernetes users*: When short names are used (e.g. "reviews"
	// instead of "reviews.default.svc.cluster.local"), Istio will interpret
	// the short name based on the namespace of the rule, not the service. A
	// rule in the "default" namespace containing a host "reviews" will be
	// interpreted as "reviews.default.svc.cluster.local", irrespective of
	// the actual namespace associated with the reviews service. _To avoid
	// potential misconfigurations, it is recommended to always use fully
	// qualified domain names over short names._
	//
	// The hosts field applies to both HTTP and TCP services. Service inside
	// the mesh, i.e., those found in the service registry, must always be
	// referred to using their alphanumeric names. IP addresses are allowed
	// only for services defined via the Gateway.
	Hosts []string `json:"hosts"`

	// The names of gateways and sidecars that should apply these routes. A
	// single VirtualService is used for sidecars inside the mesh as well as
	// for one or more gateways. The selection condition imposed by this
	// field can be overridden using the source field in the match conditions
	// of protocol-specific routes. The reserved word `mesh` is used to imply
	// all the sidecars in the mesh. When this field is omitted, the default
	// gateway (`mesh`) will be used, which would apply the rule to all
	// sidecars in the mesh. If a list of gateway names is provided, the
	// rules will apply only to the gateways. To apply the rules to both
	// gateways and sidecars, specify `mesh` as one of the gateway names.
	Gateways []string `json:"gateways,omitempty"`

	// An ordered list of route rules for HTTP traffic. HTTP routes will be
	// applied to platform service ports named 'http-*'/'http2-*'/'grpc-*', gateway
	// ports with protocol HTTP/HTTP2/GRPC/ TLS-terminated-HTTPS and service
	// entry ports using HTTP/HTTP2/GRPC protocols.  The first rule matching
	// an incoming request is used.
	HTTP []HTTPRoute `json:"http,omitempty"`

	// An ordered list of route rule for non-terminated TLS & HTTPS
	// traffic. Routing is typically performed using the SNI value presented
	// by the ClientHello message. TLS routes will be applied to platform
	// service ports named 'https-*', 'tls-*', unterminated gateway ports using
	// HTTPS/TLS protocols (i.e. with "passthrough" TLS mode) and service
	// entry ports using HTTPS/TLS protocols.  The first rule matching an
	// incoming request is used.  NOTE: Traffic 'https-*' or 'tls-*' ports
	// without associated virtual service will be treated as opaque TCP
	// traffic.
	TLS []TLSRoute `json:"tls,omitempty"`

	// An ordered list of route rules for opaque TCP traffic. TCP routes will
	// be applied to any port that is not a HTTP or TLS port. The first rule
	// matching an incoming request is used.
	TCP []TCPRoute `json:"tcp,omitempty"`

	// A list of namespaces to which this virtual service is exported. Exporting a
	// virtual service allows it to be used by sidecars and gateways defined in
	// other namespaces. This feature provides a mechanism for service owners
	// and mesh administrators to control the visibility of virtual services
	// across namespace boundaries.
	//
	// If no namespaces are specified then the virtual service is exported to all
	// namespaces by default.
	//
	// The value "." is reserved and defines an export to the same namespace that
	// the virtual service is declared in. Similarly the value "*" is reserved and
	// defines an export to all namespaces.
	//
	// NOTE: in the current release, the `exportTo` value is restricted to
	// "." or "*" (i.e., the current namespace or all namespaces).
	ExportTo []string `json:"exportTo,omitempty"`
}

// Describes match conditions and actions for routing HTTP/1.1, HTTP2, and
// gRPC traffic. See VirtualService for usage examples.
type HTTPRoute struct {
	// The name assigned to the route for debugging purposes. The
	// route's name will be concatenated with the match's name and will
	// be logged in the access logs for requests matching this
	// route/match.
	Name *string `json:"name,omitempty"`

	// Match conditions to be satisfied for the rule to be
	// activated. All conditions inside a single match block have AND
	// semantics, while the list of match blocks have OR semantics. The rule
	// is matched if any one of the match blocks succeed.
	Match []*HTTPMatchRequest `json:"match,omitempty"`

	// A http rule can either redirect or forward (default) traffic. The
	// forwarding target can be one of several versions of a service (see
	// glossary in beginning of document). Weights associated with the
	// service version determine the proportion of traffic it receives.
	Route []*HTTPRouteDestination `json:"route,omitempty"`

	// A http rule can either redirect or forward (default) traffic. If
	// traffic passthrough option is specified in the rule,
	// route/redirect will be ignored. The redirect primitive can be used to
	// send a HTTP 301 redirect to a different URI or Authority.
	Redirect *HTTPRedirect `json:"redirect,omitempty"`

	// Rewrite HTTP URIs and Authority headers. Rewrite cannot be used with
	// Redirect primitive. Rewrite will be performed before forwarding.
	Rewrite *HTTPRewrite `json:"rewrite,omitempty"`

	// Timeout for HTTP requests.
	Timeout *string `json:"timeout,omitempty"`

	// Retry policy for HTTP requests.
	Retries *HTTPRetry `json:"retries,omitempty"`

	// Fault injection policy to apply on HTTP traffic at the client side.
	// Note that timeouts or retries will not be enabled when faults are
	// enabled on the client side.
	Fault *HTTPFaultInjection `json:"fault,omitempty"`

	// Mirror HTTP traffic to a another destination in addition to forwarding
	// the requests to the intended destination. Mirrored traffic is on a
	// best effort basis where the sidecar/gateway will not wait for the
	// mirrored cluster to respond before returning the response from the
	// original destination.  Statistics will be generated for the mirrored
	// destination.
	Mirror *Destination `json:"mirror,omitempty"`

	// Percentage of the traffic to be mirrored by the `mirror` field.
	// Use of integer `mirror_percent` value is deprecated. Use the
	// double `mirror_percentage` field instead
	MirrorPercent *uint32 `json:"mirrorPercent,omitempty"`

	// Percentage of the traffic to be mirrored by the `mirror` field.
	// If this field is absent, all the traffic (100%) will be mirrored.
	// Max value is 100.
	MirrorPercentage *Percentage `json:"mirrorPercentage,omitempty"`

	// Cross-Origin Resource Sharing policy (CORS). Refer to
	// [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
	// for further details about cross origin resource sharing.
	CorsPolicy *CorsPolicy `json:"corsPolicy,omitempty"`

	// Header manipulation rules
	Headers *Headers `json:"headers,omitempty"`
}

// Message headers can be manipulated when Envoy forwards requests to,
// or responses from, a destination service. Header manipulation rules can
// be specified for a specific route destination or for all destinations.
// The following VirtualService adds a `test` header with the value `true`
// to requests that are routed to any `reviews` service destination.
// It also romoves the `foo` response header, but only from responses
// coming from the `v1` subset (version) of the `reviews` service.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: reviews-route
// spec:
//   hosts:
//   - reviews.prod.svc.cluster.local
//   http:
//   - headers:
//       request:
//         set:
//           test: true
//     route:
//     - destination:
//         host: reviews.prod.svc.cluster.local
//         subset: v2
//       weight: 25
//     - destination:
//         host: reviews.prod.svc.cluster.local
//         subset: v1
//       headers:
//         response:
//           remove:
//           - foo
//       weight: 75
// ```
type Headers struct {
	// Header manipulation rules to apply before forwarding a request
	// to the destination service
	Request *HeaderOperations `json:"request,omitempty"`

	// Header manipulation rules to apply before returning a response
	// to the caller
	Response *HeaderOperations `json:"response,omitempty"`
}

// HeaderOperations Describes the header manipulations to apply
type HeaderOperations struct {
	// Overwrite the headers specified by key with the given values
	Set map[string]string `json:"set,omitempty"`

	// Append the given values to the headers specified by keys
	// (will create a comma-separated list of values)
	Add map[string]string `json:"add,omitempty"`

	// Remove a the specified headers
	Remove []string `json:"remove,omitempty"`
}

// HttpMatchRequest specifies a set of criterion to be met in order for the
// rule to be applied to the HTTP request. For example, the following
// restricts the rule to match only requests where the URL path
// starts with /ratings/v2/ and the request contains a custom `end-user` header
// with value `jason`.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: ratings-route
// spec:
//   hosts:
//   - ratings.prod.svc.cluster.local
//   http:
//   - match:
//     - headers:
//         end-user:
//           exact: jason
//       uri:
//         prefix: "/ratings/v2/"
//       ignoreUriCase: true
//     route:
//     - destination:
//         host: ratings.prod.svc.cluster.local
// ```
//
// HTTPMatchRequest CANNOT be empty.
type HTTPMatchRequest struct {
	// The name assigned to a match. The match's name will be
	// concatenated with the parent route's name and will be logged in
	// the access logs for requests matching this route.
	Name *string `json:"name,omitempty"`

	// URI to match
	// values are case-sensitive and formatted as follows:
	//
	// - `exact: "value"` for exact string match
	//
	// - `prefix: "value"` for prefix-based match
	//
	// - `regex: "value"` for ECMAscript style regex-based match
	//
	// **Note:** Case-insensitive matching could be enabled via the
	// `ignore_uri_case` flag.
	URI *v1alpha1.StringMatch `json:"uri,omitempty"`

	// URI Scheme
	// values are case-sensitive and formatted as follows:
	//
	// - `exact: "value"` for exact string match
	//
	// - `prefix: "value"` for prefix-based match
	//
	// - `regex: "value"` for ECMAscript style regex-based match
	//
	Scheme *v1alpha1.StringMatch `json:"scheme,omitempty"`

	// HTTP Method
	// values are case-sensitive and formatted as follows:
	//
	// - `exact: "value"` for exact string match
	//
	// - `prefix: "value"` for prefix-based match
	//
	// - `regex: "value"` for ECMAscript style regex-based match
	//
	Method *v1alpha1.StringMatch `json:"method,omitempty"`

	// HTTP Authority
	// values are case-sensitive and formatted as follows:
	//
	// - `exact: "value"` for exact string match
	//
	// - `prefix: "value"` for prefix-based match
	//
	// - `regex: "value"` for ECMAscript style regex-based match
	//
	Authority *v1alpha1.StringMatch `json:"authority,omitempty"`

	// The header keys must be lowercase and use hyphen as the separator,
	// e.g. _x-request-id_.
	//
	// Header values are case-sensitive and formatted as follows:
	//
	// - `exact: "value"` for exact string match
	//
	// - `prefix: "value"` for prefix-based match
	//
	// - `regex: "value"` for ECMAscript style regex-based match
	//
	// **Note:** The keys `uri`, `scheme`, `method`, and `authority` will be ignored.
	Headers map[string]v1alpha1.StringMatch `json:"headers,omitempty"`

	// Specifies the ports on the host that is being addressed. Many services
	// only expose a single port or label ports with the protocols they support,
	// in these cases it is not required to explicitly select the port.
	Port *uint32 `json:"port,omitempty"`

	// One or more labels that constrain the applicability of a rule to
	// workloads with the given labels. If the VirtualService has a list of
	// gateways specified at the top, it must include the reserved gateway
	// `mesh` for this field to be applicable.
	SourceLabels map[string]string `json:"sourceLabels,omitempty"`

	// Query parameters for matching.
	//
	// Ex:
	// - For a query parameter like "?key=true", the map key would be "key" and
	//   the string match could be defined as `exact: "true"`.
	// - For a query parameter like "?key", the map key would be "key" and the
	//   string match could be defined as `exact: ""`.
	// - For a query parameter like "?key=123", the map key would be "key" and the
	//   string match could be defined as `regex: "\d+$"`. Note that this
	//   configuration will only match values like "123" but not "a123" or "123a".
	//
	// **Note:** `prefix` matching is currently not supported.
	QueryParams map[string]*v1alpha1.StringMatch `json:"queryParams,omitempty"`

	// Flag to specify whether the URI matching should be case-insensitive.
	//
	// **Note:** The case will be ignored only in the case of `exact` and `prefix`
	// URI matches.
	IgnoreURICase *bool `json:"ignoreUriCase,omitempty"`
}

// Each routing rule is associated with one or more service versions (see
// glossary in beginning of document). Weights associated with the version
// determine the proportion of traffic it receives. For example, the
// following rule will route 25% of traffic for the "reviews" service to
// instances with the "v2" tag and the remaining traffic (i.e., 75%) to
// "v1".
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: reviews-route
// spec:
//   hosts:
//   - reviews.prod.svc.cluster.local
//   http:
//   - route:
//     - destination:
//         host: reviews.prod.svc.cluster.local
//         subset: v2
//       weight: 25
//     - destination:
//         host: reviews.prod.svc.cluster.local
//         subset: v1
//       weight: 75
// ```
//
// And the associated DestinationRule
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: DestinationRule
// metadata:
//   name: reviews-destination
// spec:
//   host: reviews.prod.svc.cluster.local
//   subsets:
//   - name: v1
//     labels:
//       version: v1
//   - name: v2
//     labels:
//       version: v2
// ```
//
// Traffic can also be split across two entirely different services without
// having to define new subsets. For example, the following rule forwards 25% of
// traffic to reviews.com to dev.reviews.com
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: reviews-route-two-domains
// spec:
//   hosts:
//   - reviews.com
//   http:
//   - route:
//     - destination:
//         host: dev.reviews.com
//       weight: 25
//     - destination:
//         host: reviews.com
//       weight: 75
// ```
type HTTPRouteDestination struct {
	// REQUIRED. Destination uniquely identifies the instances of a service
	// to which the request/connection should be forwarded to.
	Destination *Destination `json:"destination"`

	// REQUIRED. The proportion of traffic to be forwarded to the service
	// version. (0-100). Sum of weights across destinations SHOULD BE == 100.
	// If there is only one destination in a rule, the weight value is assumed to
	// be 100.
	Weight *int `json:"weight,omitempty"`

	// Header manipulation rules
	Headers *Headers `json:"headers,omitempty"`
}

// L4 routing rule weighted destination.
type RouteDestination struct {
	// REQUIRED. Destination uniquely identifies the instances of a service
	// to which the request/connection should be forwarded to.
	Destination *Destination `json:"destination"`

	// REQUIRED. The proportion of traffic to be forwarded to the service
	// version. (0-100). Sum of weights across destinations SHOULD BE == 100.
	// If there is only one destination in a rule, the weight value is assumed to
	// be 100.
	Weight *int `json:"weight,omitempty"`
}

// Destination indicates the network addressable service to which the
// request/connection will be sent after processing a routing rule. The
// destination.host should unambiguously refer to a service in the service
// registry. Istio's service registry is composed of all the services found
// in the platform's service registry (e.g., Kubernetes services, Consul
// services), as well as services declared through the
// [ServiceEntry](https://istio.io/docs/reference/config/networking/v1beta1/service-entry/#ServiceEntry) resource.
//
// *Note for Kubernetes users*: When short names are used (e.g. "reviews"
// instead of "reviews.default.svc.cluster.local"), Istio will interpret
// the short name based on the namespace of the rule, not the service. A
// rule in the "default" namespace containing a host "reviews will be
// interpreted as "reviews.default.svc.cluster.local", irrespective of the
// actual namespace associated with the reviews service. _To avoid potential
// misconfigurations, it is recommended to always use fully qualified
// domain names over short names._
//
// The following Kubernetes example routes all traffic by default to pods
// of the reviews service with label "version: v1" (i.e., subset v1), and
// some to subset v2, in a Kubernetes environment.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: reviews-route
//   namespace: foo
// spec:
//   hosts:
//   - reviews # interpreted as reviews.foo.svc.cluster.local
//   http:
//   - match:
//     - uri:
//         prefix: "/wpcatalog"
//     - uri:
//         prefix: "/consumercatalog"
//     rewrite:
//       uri: "/newcatalog"
//     route:
//     - destination:
//         host: reviews # interpreted as reviews.foo.svc.cluster.local
//         subset: v2
//   - route:
//     - destination:
//         host: reviews # interpreted as reviews.foo.svc.cluster.local
//         subset: v1
// ```
//
// And the associated DestinationRule
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: DestinationRule
// metadata:
//   name: reviews-destination
//   namespace: foo
// spec:
//   host: reviews # interpreted as reviews.foo.svc.cluster.local
//   subsets:
//   - name: v1
//     labels:
//       version: v1
//   - name: v2
//     labels:
//       version: v2
// ```
//
// The following VirtualService sets a timeout of 5s for all calls to
// productpage.prod.svc.cluster.local service in Kubernetes. Notice that
// there are no subsets defined in this rule. Istio will fetch all
// instances of productpage.prod.svc.cluster.local service from the service
// registry and populate the sidecar's load balancing pool. Also, notice
// that this rule is set in the istio-system namespace but uses the fully
// qualified domain name of the productpage service,
// productpage.prod.svc.cluster.local. Therefore the rule's namespace does
// not have an impact in resolving the name of the productpage service.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: my-productpage-rule
//   namespace: istio-system
// spec:
//   hosts:
//   - productpage.prod.svc.cluster.local # ignores rule namespace
//   http:
//   - timeout: 5s
//     route:
//     - destination:
//         host: productpage.prod.svc.cluster.local
// ```
//
// To control routing for traffic bound to services outside the mesh, external
// services must first be added to Istio's internal service registry using the
// ServiceEntry resource. VirtualServices can then be defined to control traffic
// bound to these external services. For example, the following rules define a
// Service for wikipedia.org and set a timeout of 5s for http requests.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: ServiceEntry
// metadata:
//   name: external-svc-wikipedia
// spec:
//   hosts:
//   - wikipedia.org
//   location: MESH_EXTERNAL
//   ports:
//   - number: 80
//     name: example-http
//     protocol: HTTP
//   resolution: DNS
//
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: my-wiki-rule
// spec:
//   hosts:
//   - wikipedia.org
//   http:
//   - timeout: 5s
//     route:
//     - destination:
//         host: wikipedia.org
// ```
type Destination struct {
	// REQUIRED. The name of a service from the service registry. Service
	// names are looked up from the platform's service registry (e.g.,
	// Kubernetes services, Consul services, etc.) and from the hosts
	// declared by [ServiceEntry](https://istio.io/docs/reference/config/networking/v1beta1/service-entry/#ServiceEntry). Traffic forwarded to
	// destinations that are not found in either of the two, will be dropped.
	//
	// *Note for Kubernetes users*: When short names are used (e.g. "reviews"
	// instead of "reviews.default.svc.cluster.local"), Istio will interpret
	// the short name based on the namespace of the rule, not the service. A
	// rule in the "default" namespace containing a host "reviews will be
	// interpreted as "reviews.default.svc.cluster.local", irrespective of
	// the actual namespace associated with the reviews service. _To avoid
	// potential misconfigurations, it is recommended to always use fully
	// qualified domain names over short names._
	Host string `json:"host"`

	// The name of a subset within the service. Applicable only to services
	// within the mesh. The subset must be defined in a corresponding
	// DestinationRule.
	Subset *string `json:"subset,omitempty"`

	// Specifies the port on the host that is being addressed. If a service
	// exposes only a single port it is not required to explicitly select the
	// port.
	Port *PortSelector `json:"port,omitempty"`
}

// PortSelector specifies the number of a port to be used for
// matching or selection for final routing.
type PortSelector struct {
	// Valid port number
	Number uint32 `json:"number"`
}

// Describes match conditions and actions for routing TCP traffic. The
// following routing rule forwards traffic arriving at port 27017 for
// mongo.prod.svc.cluster.local to another Mongo server on port 5555.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: bookinfo-Mongo
// spec:
//   hosts:
//   - mongo.prod.svc.cluster.local
//   tcp:
//   - match:
//     - port: 27017
//     route:
//     - destination:
//         host: mongo.backup.svc.cluster.local
//         port:
//           number: 5555
// ```
type TCPRoute struct {
	// Match conditions to be satisfied for the rule to be
	// activated. All conditions inside a single match block have AND
	// semantics, while the list of match blocks have OR semantics. The rule
	// is matched if any one of the match blocks succeed.
	Match []L4MatchAttributes `json:"match"`

	// The destination to which the connection should be forwarded to.
	Route []*RouteDestination `json:"route"`
}

// Describes match conditions and actions for routing unterminated TLS
// traffic (TLS/HTTPS) The following routing rule forwards unterminated TLS
// traffic arriving at port 443 of gateway called mygateway to internal
// services in the mesh based on the SNI value.
//
// ```yaml
// kind: VirtualService
// metadata:
//   name: bookinfo-sni
// spec:
//   hosts:
//   - '*.bookinfo.com'
//   gateways:
//   - mygateway
//   tls:
//   - match:
//     - port: 443
//       sniHosts:
//       - login.bookinfo.com
//     route:
//     - destination:
//         host: login.prod.svc.cluster.local
//   - match:
//     - port: 443
//       sniHosts:
//       - reviews.bookinfo.com
//     route:
//     - destination:
//         host: reviews.prod.svc.cluster.local
// ```
type TLSRoute struct {
	// REQUIRED. Match conditions to be satisfied for the rule to be
	// activated. All conditions inside a single match block have AND
	// semantics, while the list of match blocks have OR semantics. The rule
	// is matched if any one of the match blocks succeed.
	Match []TLSMatchAttributes `json:"match"`

	// The destination to which the connection should be forwarded to.
	Route []*RouteDestination `json:"route"`
}

// L4 connection match attributes. Note that L4 connection matching support
// is incomplete.
type L4MatchAttributes struct {
	// IPv4 or IPv6 ip addresses of destination with optional subnet.  E.g.,
	// a.b.c.d/xx form or just a.b.c.d.
	DestinationSubnets []string `json:"destinationSubnets,omitempty"`

	// Specifies the port on the host that is being addressed. Many services
	// only expose a single port or label ports with the protocols they support,
	// in these cases it is not required to explicitly select the port.
	Port *int `json:"port,omitempty"`

	// One or more labels that constrain the applicability of a rule to
	// workloads with the given labels. If the VirtualService has a list of
	// gateways specified at the top, it should include the reserved gateway
	// `mesh` in order for this field to be applicable.
	SourceLabels map[string]string `json:"sourceLabels,omitempty"`

	// Names of gateways where the rule should be applied to. Gateway names
	// at the top of the VirtualService (if any) are overridden. The gateway
	// match is independent of sourceLabels.
	Gateways []string `json:"gateways,omitempty"`
}

// TLS connection match attributes.
type TLSMatchAttributes struct {
	// REQUIRED. SNI (server name indicator) to match on. Wildcard prefixes
	// can be used in the SNI value, e.g., *.com will match foo.example.com
	// as well as example.com. An SNI value must be a subset (i.e., fall
	// within the domain) of the corresponding virtual serivce's hosts.
	SniHosts []string `json:"sniHosts"`

	// IPv4 or IPv6 ip addresses of destination with optional subnet.  E.g.,
	// a.b.c.d/xx form or just a.b.c.d.
	DestinationSubnets []string `json:"destinationSubnets,omitempty"`

	// Specifies the port on the host that is being addressed. Many services
	// only expose a single port or label ports with the protocols they support,
	// in these cases it is not required to explicitly select the port.
	Port *int `json:"port,omitempty"`

	// One or more labels that constrain the applicability of a rule to
	// workloads with the given labels. If the VirtualService has a list of
	// gateways specified at the top, it should include the reserved gateway
	// `mesh` in order for this field to be applicable.
	SourceLabels map[string]string `json:"sourceLabels,omitempty"`

	// Names of gateways where the rule should be applied to. Gateway names
	// at the top of the VirtualService (if any) are overridden. The gateway
	// match is independent of sourceLabels.
	Gateways []string `json:"gateways,omitempty"`
}

// HTTPRedirect can be used to send a 301 redirect response to the caller,
// where the Authority/Host and the URI in the response can be swapped with
// the specified values. For example, the following rule redirects
// requests for /v1/getProductRatings API on the ratings service to
// /v1/bookRatings provided by the bookratings service.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: ratings-route
// spec:
//   hosts:
//   - ratings.prod.svc.cluster.local
//   http:
//   - match:
//     - uri:
//         exact: /v1/getProductRatings
//     redirect:
//       uri: /v1/bookRatings
//       authority: newratings.default.svc.cluster.local
//   ...
// ```
type HTTPRedirect struct {
	// On a redirect, overwrite the Path portion of the URL with this
	// value. Note that the entire path will be replaced, irrespective of the
	// request URI being matched as an exact path or prefix.
	URI *string `json:"uri,omitempty"`

	// On a redirect, overwrite the Authority/Host portion of the URL with
	// this value.
	Authority *string `json:"authority,omitempty"`

	// On a redirect, Specifies the HTTP status code to use in the redirect
	// response. The default response code is MOVED_PERMANENTLY (301).
	RedirectCode *uint32 `json:"redirectCode,omitempty"`
}

// HTTPRewrite can be used to rewrite specific parts of a HTTP request
// before forwarding the request to the destination. Rewrite primitive can
// be used only with HTTPRouteDestination. The following example
// demonstrates how to rewrite the URL prefix for api call (/ratings) to
// ratings service before making the actual API call.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: ratings-route
// spec:
//   hosts:
//   - ratings.prod.svc.cluster.local
//   http:
//   - match:
//     - uri:
//         prefix: /ratings
//     rewrite:
//       uri: /v1/bookRatings
//     route:
//     - destination:
//         host: ratings.prod.svc.cluster.local
//         subset: v1
// ```
type HTTPRewrite struct {
	// rewrite the path (or the prefix) portion of the URI with this
	// value. If the original URI was matched based on prefix, the value
	// provided in this field will replace the corresponding matched prefix.
	URI *string `json:"uri,omitempty"`

	// rewrite the Authority/Host header with this value.
	Authority *string `json:"authority,omitempty"`
}

// Describes the retry policy to use when a HTTP request fails. For
// example, the following rule sets the maximum number of retries to 3 when
// calling ratings:v1 service, with a 2s timeout per retry attempt.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: ratings-route
// spec:
//   hosts:
//   - ratings.prod.svc.cluster.local
//   http:
//   - route:
//     - destination:
//         host: ratings.prod.svc.cluster.local
//         subset: v1
//     retries:
//       attempts: 3
//       perTryTimeout: 2s
//       retryOn: gateway-error,connect-failure,refused-stream
// ```
type HTTPRetry struct {
	// REQUIRED. Number of retries for a given request. The interval
	// between retries will be determined automatically (25ms+). Actual
	// number of retries attempted depends on the httpReqTimeout.
	Attempts int `json:"attempts"`

	// Timeout per retry attempt for a given request. format: 1h/1m/1s/1ms. MUST BE >=1ms.
	PerTryTimeout string `json:"perTryTimeout"`

	// Specifies the conditions under which retry takes place.
	// One or more policies can be specified using a ‘,’ delimited list.
	// See the [retry policies](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/router_filter#x-envoy-retry-on)
	// and [gRPC retry policies](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/router_filter#x-envoy-retry-grpc-on) for more details.
	RetryOn *string `json:"retryOn,omitempty"`
}

// Describes the Cross-Origin Resource Sharing (CORS) policy, for a given
// service. Refer to [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS)
// for further details about cross origin resource sharing. For example,
// the following rule restricts cross origin requests to those originating
// from example.com domain using HTTP POST/GET, and sets the
// `Access-Control-Allow-Credentials` header to false. In addition, it only
// exposes `X-Foo-bar` header and sets an expiry period of 1 day.
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: ratings-route
// spec:
//   hosts:
//   - ratings.prod.svc.cluster.local
//   http:
//   - route:
//     - destination:
//         host: ratings.prod.svc.cluster.local
//         subset: v1
//     corsPolicy:
//       allowOrigin:
//       - example.com
//       allowMethods:
//       - POST
//       - GET
//       allowCredentials: false
//       allowHeaders:
//       - X-Foo-Bar
//       maxAge: "24h"
// ```
type CorsPolicy struct {
	// The list of origins that are allowed to perform CORS requests. The
	// content will be serialized into the Access-Control-Allow-Origin
	// header. Wildcard * will allow all origins.
	AllowOrigin []string `json:"allowOrigin,omitempty"`

	// List of HTTP methods allowed to access the resource. The content will
	// be serialized into the Access-Control-Allow-Methods header.
	AllowMethods []string `json:"allowMethods,omitempty"`

	// List of HTTP headers that can be used when requesting the
	// resource. Serialized to Access-Control-Allow-Methods header.
	AllowHeaders []string `json:"allowHeaders,omitempty"`

	// A white list of HTTP headers that the browsers are allowed to
	// access. Serialized into Access-Control-Expose-Headers header.
	ExposeHeaders []string `json:"exposeHeaders,omitempty"`

	// Specifies how long the results of a preflight request can be
	// cached. Translates to the `Access-Control-Max-Age` header.
	MaxAge *string `json:"maxAge,omitempty"`

	// Indicates whether the caller is allowed to send the actual request
	// (not the preflight) using credentials. Translates to
	// `Access-Control-Allow-Credentials` header.
	AllowCredentials *bool `json:"allowCredentials,omitempty"`
}

// HTTPFaultInjection can be used to specify one or more faults to inject
// while forwarding http requests to the destination specified in a route.
// Fault specification is part of a VirtualService rule. Faults include
// aborting the Http request from downstream service, and/or delaying
// proxying of requests. A fault rule MUST HAVE delay or abort or both.
//
// *Note:* Delay and abort faults are independent of one another, even if
// both are specified simultaneously.
type HTTPFaultInjection struct {
	// Delay requests before forwarding, emulating various failures such as
	// network issues, overloaded upstream service, etc.
	Delay *Delay `json:"delay,omitempty"`

	// Abort Http request attempts and return error codes back to downstream
	// service, giving the impression that the upstream service is faulty.
	Abort *Abort `json:"abort,omitempty"`
}

// Delay specification is used to inject latency into the request
// forwarding path. The following example will introduce a 5 second delay
// in 1 out of every 1000 requests to the "v1" version of the "reviews"
// service from all pods with label env: prod
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: reviews-route
// spec:
//   hosts:
//   - reviews.prod.svc.cluster.local
//   http:
//   - match:
//     - sourceLabels:
//         env: prod
//     route:
//     - destination:
//         host: reviews.prod.svc.cluster.local
//         subset: v1
//     fault:
//       delay:
//         percentage:
//           value: 0.1
//         fixedDelay: 5s
// ```
//
// The _fixedDelay_ field is used to indicate the amount of delay in seconds.
// The optional _percentage_ field can be used to only delay a certain
// percentage of requests. If left unspecified, all request will be delayed.
type Delay struct {
	// REQUIRED. Add a fixed delay before forwarding the request. Format:
	// 1h/1m/1s/1ms. MUST be >=1ms.
	FixedDelay string `json:"fixedDelay"`

	// Percentage of requests on which the delay will be injected.
	Percentage *Percentage `json:"percentage,omitempty"`
}

// Abort specification is used to prematurely abort a request with a
// pre-specified error code. The following example will return an HTTP 400
// error code for 1 out of every 1000 requests to the "ratings" service "v1".
//
// ```yaml
// apiVersion: networking.istio.io/v1beta1
// kind: VirtualService
// metadata:
//   name: ratings-route
// spec:
//   hosts:
//   - ratings.prod.svc.cluster.local
//   http:
//   - route:
//     - destination:
//         host: ratings.prod.svc.cluster.local
//         subset: v1
//     fault:
//       abort:
//         percentage:
//           value: 0.1
//         httpStatus: 400
// ```
//
// The _httpStatus_ field is used to indicate the HTTP status code to
// return to the caller. The optional _percentage_ field can be used to only
// abort a certain percentage of requests. If not specified, all requests are
// aborted.
type Abort struct {
	// REQUIRED. HTTP status code to use to abort the Http request.
	HTTPStatus int `json:"httpStatus"`

	// Percentage of requests on which the delay will be injected.
	Percentage *Percentage `json:"percentage,omitempty"`
}

// Percent specifies a percentage in the range of [0.0, 100.0].
type Percentage struct {
	Value float32 `json:"value"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualServiceList is a list of VirtualService resources
type VirtualServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []VirtualService `json:"items"`
}
