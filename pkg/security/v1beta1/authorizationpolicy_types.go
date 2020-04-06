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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	selector "github.com/banzaicloud/istio-client-go/pkg/type/v1beta1"
)

// Istio Authorization Policy enables access control on workloads in the mesh.
//
// Authorization policy supports both allow and deny policies. When allow and
// deny policies are used for a workload at the same time, the deny policies are
// evaluated first. The evaluation is determined by the following rules:
//
// 1. If there are any DENY policies that match the request, deny the request.
// 2. If there are no ALLOW policies for the workload, allow the request.
// 3. If any of the ALLOW policies match the request, allow the request.
// 4. Deny the request.
//
// For example, the following authorization policy sets the `action` to "ALLOW"
// to create an allow policy. The default action is "ALLOW" but it is useful
// to be explicit in the policy.
//
// It allows requests from:
//
// - service account "cluster.local/ns/default/sa/sleep" or
// - namespace "test"
//
// to access the workload with:
//
// - "GET" method at paths of prefix "/info" or,
// - "POST" method at path "/data".
//
// when the request has a valid JWT token issued by "https://accounts.google.com".
//
// Any other requests will be denied.
//
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: AuthorizationPolicy
// metadata:
//  name: httpbin
//  namespace: foo
// spec:
//  action: ALLOW
//  rules:
//  - from:
//    - source:
//        principals: ["cluster.local/ns/default/sa/sleep"]
//    - source:
//        namespaces: ["test"]
//    to:
//    - operation:
//        methods: ["GET"]
//        paths: ["/info*"]
//    - operation:
//        methods: ["POST"]
//        paths: ["/data"]
//    when:
//    - key: request.auth.claims[iss]
//      values: ["https://accounts.google.com"]
// ```
//
// The following is another example that sets `action` to "DENY" to create a deny policy.
// It denies requests from the "dev" namespace to the "POST" method on all workloads
// in the "foo" namespace.
//
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: AuthorizationPolicy
// metadata:
//  name: httpbin
//  namespace: foo
// spec:
//  action: DENY
//  rules:
//  - from:
//    - source:
//        namespaces: ["dev"]
//    to:
//    - operation:
//        methods: ["POST"]
// ```
//
// Authorization Policy scope (target) is determined by "metadata/namespace" and
// an optional "selector".
//
// - "metadata/namespace" tells which namespace the policy applies. If set to root
// namespace, the policy applies to all namespaces in a mesh.
// - workload "selector" can be used to further restrict where a policy applies.
//
// For example,
//
// The following authorization policy applies to workloads containing label
// "app: httpbin" in namespace bar.
//
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: AuthorizationPolicy
// metadata:
//  name: policy
//  namespace: bar
// spec:
//  selector:
//    matchLabels:
//      app: httpbin
// ```
//
// The following authorization policy applies to all workloads in namespace foo.
//
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: AuthorizationPolicy
// metadata:
//  name: policy
//  namespace: foo
// spec:
//   {}
// ```
//
// The following authorization policy applies to workloads containing label
// "version: v1" in all namespaces in the mesh. (Assuming the root namespace is
// configured to "istio-config").
//
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: AuthorizationPolicy
// metadata:
//  name: policy
//  namespace: istio-config
// spec:
//  selector:
//    matchLabels:
//      version: v1
// ```

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// AuthorizationPolicy
type AuthorizationPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AuthorizationPolicySpec `json:"spec"`
}

// AuthorizationPolicy enables access control on workloads.
//
// For example, the following authorization policy denies all requests to workloads
// in namespace foo.
//
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: AuthorizationPolicy
// metadata:
//  name: deny-all
//  namespace: foo
// spec:
//   {}
// ```
//
// The following authorization policy allows all requests to workloads in namespace
// foo.
//
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: AuthorizationPolicy
// metadata:
//  name: allow-all
//  namespace: foo
// spec:
//  rules:
//  - {}
// ```
type AuthorizationPolicySpec struct {
	// Optional. Workload selector decides where to apply the authorization policy.
	// If not set, the authorization policy will be applied to all workloads in the
	// same namespace as the authorization policy.
	Selector *selector.WorkloadSelector `json:"selector,omitempty"`
	// Optional. A list of rules to match the request. A match occurs when at least
	// one rule matches the request.
	//
	// If not set, the match will never occur. This is equivalent to setting a
	// default of deny for the target workloads.
	Rules []*Rule `json:"rules,omitempty"`
	// Optional. The action to take if the request is matched with the rules.
	Action AuthorizationPolicyAction `json:"action,omitempty"`
}

// Action specifies the operation to take.
type AuthorizationPolicyAction string

const (
	// Allow a request only if it matches the rules. This is the default type.
	AuthorizationPolicyActionAllow AuthorizationPolicyAction = "ALLOW"
	// Deny a request if it matches any of the rules.
	AuthorizationPolicyActionDeny AuthorizationPolicyAction = "DENY"
)

// Rule matches requests from a list of sources that perform a list of operations subject to a
// list of conditions. A match occurs when at least one source, operation and condition
// matches the request. An empty rule is always matched.
//
// Any string field in the rule supports Exact, Prefix, Suffix and Presence match:
//
// - Exact match: "abc" will match on value "abc".
// - Prefix match: "abc*" will match on value "abc" and "abcd".
// - Suffix match: "*abc" will match on value "abc" and "xabc".
// - Presence match: "*" will match when value is not empty.
type Rule struct {
	// Optional. from specifies the source of a request.
	//
	// If not set, any source is allowed.
	From []*RuleFrom `json:"from,omitempty"`
	// Optional. to specifies the operation of a request.
	//
	// If not set, any operation is allowed.
	To []*RuleTo `json:"to,omitempty"`
	// Optional. when specifies a list of additional conditions of a request.
	//
	// If not set, any condition is allowed.
	When []*Condition `json:"when,omitempty"`
}

// From includes a list or sources.
type RuleFrom struct {
	// Source specifies the source of a request.
	Source *Source `json:"source,omitempty"`
}

// To includes a list or operations.
type RuleTo struct {
	// Operation specifies the operation of a request.
	Operation *Operation `json:"operation,omitempty"`
}

// Source specifies the source identities of a request. Fields in the source are
// ANDed together.
//
// For example, the following source matches if the principal is "admin" or "dev"
// and the namespace is "prod" or "test" and the ip is not "1.2.3.4".
//
// ```yaml
// principals: ["admin", "dev"]
// namespaces: ["prod", "test"]
// not_ipblocks: ["1.2.3.4"]
// ```
type Source struct {
	// Optional. A list of source peer identities (i.e. service account), which
	// matches to the "source.principal" attribute. This field requires mTLS enabled.
	//
	// If not set, any principal is allowed.
	Principals []string `json:"principals,omitempty"`
	// Optional. A list of negative match of source peer identities.
	NotPrincipals []string `json:"notPrincipals,omitempty"`
	// Optional. A list of request identities (i.e. "iss/sub" claims), which
	// matches to the "request.auth.principal" attribute.
	//
	// If not set, any request principal is allowed.
	RequestPrincipals []string `json:"requestPrincipals,omitempty"`
	// Optional. A list of negative match of request identities.
	NotRequestPrincipals []string `json:"notRequestPrincipals,omitempty"`
	// Optional. A list of namespaces, which matches to the "source.namespace"
	// attribute. This field requires mTLS enabled.
	//
	// If not set, any namespace is allowed.
	Namespaces []string `json:"namespaces,omitempty"`
	// Optional. A list of negative match of namespaces.
	NotNamespaces []string `json:"notNamespaces,omitempty"`
	// Optional. A list of IP blocks, which matches to the "source.ip" attribute.
	// Single IP (e.g. "1.2.3.4") and CIDR (e.g. "1.2.3.0/24") are supported.
	//
	// If not set, any IP is allowed.
	IPBlocks []string `json:"ipBlocks,omitempty"`
	// Optional. A list of negative match of IP blocks.
	NotIPBlocks []string `json:"notIpBlocks,omitempty"`
}

// Operation specifies the operations of a request. Fields in the operation are
// ANDed together.
//
// For example, the following operation matches if the host has suffix ".example.com"
// and the method is "GET" or "HEAD" and the path doesn't have prefix "/admin".
//
// ```yaml
// hosts: ["*.example.com"]
// methods: ["GET", "HEAD"]
// not_paths: ["/admin*"]
// ```
type Operation struct {
	// Optional. A list of hosts, which matches to the "request.host" attribute.
	//
	// If not set, any host is allowed. Must be used only with HTTP.
	Hosts []string `json:"hosts,omitempty"`
	// Optional. A list of negative match of hosts.
	NotHosts []string `json:"notHosts,omitempty"`
	// Optional. A list of ports, which matches to the "destination.port" attribute.
	//
	// If not set, any port is allowed.
	Ports []string `json:"ports,omitempty"`
	// Optional. A list of negative match of ports.
	NotPorts []string `json:"notPorts,omitempty"`
	// Optional. A list of methods, which matches to the "request.method" attribute.
	// For gRPC service, this will always be "POST".
	//
	// If not set, any method is allowed. Must be used only with HTTP.
	Methods []string `json:"methods,omitempty"`
	// Optional. A list of negative match of methods.
	NotMethods []string `json:"notMethods,omitempty"`
	// Optional. A list of paths, which matches to the "request.url_path" attribute.
	// For gRPC service, this will be the fully-qualified name in the form of
	// "/package.service/method".
	//
	// If not set, any path is allowed. Must be used only with HTTP.
	Paths []string `json:"paths,omitempty"`
	// Optional. A list of negative match of paths.
	NotPaths []string `json:"notPaths,omitempty"`
}

// Condition specifies additional required attributes.
type Condition struct {
	// The name of an Istio attribute.
	// See the [full list of supported attributes](https://istio.io/docs/reference/config/security/conditions/).
	Key string `json:"key,omitempty"`
	// Optional. A list of allowed values for the attribute.
	// Note: at least one of values or not_values must be set.
	Values []string `json:"values,omitempty"`
	// Optional. A list of negative match of values for the attribute.
	// Note: at least one of values or not_values must be set.
	NotValues []string `json:"notValues,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// AuthorizationPolicyList is a list of AuthorizationPolicy resources
type AuthorizationPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AuthorizationPolicy `json:"items"`
}
