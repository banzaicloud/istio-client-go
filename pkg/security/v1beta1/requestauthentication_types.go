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

// RequestAuthentication defines what request authentication methods are supported by a workload.
// If will reject a request if the request contains invalid authentication information, based on the
// configured authentication rules. A request that does not contain any authentication credentials
// will be accepted but will not have any authenticated identity. To restrict access to authenticated
// requests only, this should be accompanied by an authorization rule.
// Examples:
//
// - Require JWT for all request for workloads that have label `app:httpbin`
//
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: RequestAuthentication
// metadata:
//   name: httpbin
//   namespace: foo
// spec:
//   selector:
//     matchLabels:
//       app: httpbin
//   jwtRules:
//   - issuer: "issuer-foo"
//     jwksUri: https://example.com/.well-known/jwks.json
// ---
// apiVersion: security.istio.io/v1beta1
// kind: AuthorizationPolicy
// metadata:
//   name: httpbin
//   namespace: foo
// spec:
//   selector:
//     matchLabels:
//       app: httpbin
//   rules:
//   - from:
//     - source:
//         requestPrincipals: ["*"]
// ```
//
// - The next example shows how to set a different JWT requirement for a different `host`. The `RequestAuthentication`
// declares it can accpet JWTs issuer by either `issuer-foo` or `issuer-bar` (the public key set is implicitly
// set from the OpenID Connect spec).
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: RequestAuthentication
// metadata:
//   name: httpbin
//   namespace: foo
// spec:
//   selector:
//     matchLabels:
//       app: httpbin
//   jwtRules:
//   - issuer: "issuer-foo"
//   - issuer: "issuer-bar"
// ---
// apiVersion: security.istio.io/v1beta1
// kind: AuthorizationPolicy
// metadata:
//   name: httpbin
//   namespace: foo
// spec:
//   selector:
//     matchLabels:
//       app: httpbin
//  rules:
//  - from:
//    - source:
//        requestPrincipals: ["issuer-foo/*"]
//    to:
//      hosts: ["example.com"]
//  - from:
//    - source:
//        requestPrincipals: ["issuer-bar/*"]
//    to:
//      hosts: ["another-host.com"]
// ```
//
// - You can fine tune the authorization policy to set different requirement per path. For example,
// to require JWT on all paths, except /healthz, the same `RequestAuthentication` can be used, but the
// authorization policy could be:
//
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: AuthorizationPolicy
// metadata:
//   name: httpbin
//   namespace: foo
// spec:
//   selector:
//     matchLabels:
//       app: httpbin
//  rules:
//  - from:
//    - source:
//        requestPrincipals: ["*"]
//  - to:
//    - operation:
//        paths: ["/healthz]
// ```

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// RequestAuthentication
type RequestAuthentication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RequestAuthenticationSpec `json:"spec"`
}

type RequestAuthenticationSpec struct {
	// The selector determines the workloads to apply the RequestAuthentication on.
	// If not set, the policy will be applied to all workloads in the same namespace as the policy.
	Selector *selector.WorkloadSelector `json:"selector,omitempty"`
	// Define the list of JWTs that can be validated at the selected workloads' proxy. A valid token
	// will be used to extract the authenticated identity.
	// Each rule will be activated only when a token is presented at the location recorgnized by the
	// rule. The token will be validated based on the JWT rule config. If validation fails, the request will
	// be rejected.
	// Note: if more than one token is presented (at different locations), the output principal is nondeterministic.
	JwtRules []*JWTRule `json:"jwtRules,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// RequestAuthenticationList is a list of RequestAuthentication resources
type RequestAuthenticationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []RequestAuthentication `json:"items"`
}
