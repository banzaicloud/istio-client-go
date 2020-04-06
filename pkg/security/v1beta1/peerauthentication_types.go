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

type MTLSMode string

const (
	// Inherit from parent, if has one. Otherwise treated as PERMISSIVE.
	MTLSModeUnset MTLSMode = "UNSET"
	// Connection is not tunneled.
	MTLSModeDisable MTLSMode = "DISABLE"
	// Connection can be either plaintext or mTLS tunnel.
	MTLSModePermissive MTLSMode = "PERMISSIVE"
	// Connection is an mTLS tunnel (TLS with client cert must be presented).
	MTLSModeStrict MTLSMode = "STRICT"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// PeerAuthentication
type PeerAuthentication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PeerAuthenticationSpec `json:"spec"`
}

// PeerAuthentication defines how traffic will be tunneled (or not) to the sidecar.
//
// Examples:
//
// Policy to allow mTLS traffic for all workloads under namespace `foo`:
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: PeerAuthentication
// metadata:
//   name: default
//   namespace: foo
// spec:
//   mtls:
//     mode: STRICT
// ```
// For mesh level, put the policy in root-namespace according to your Istio installation.
//
// Policies to allow both mTLS & plaintext traffic for all workloads under namespace `foo`, but
// require mTLS for workload `finance`.
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: PeerAuthentication
// metadata:
//   name: default
//   namespace: foo
// spec:
//   mtls:
//     mode: PERMISSIVE
// ---
// apiVersion: security.istio.io/v1beta1
// kind: PeerAuthentication
// metadata:
//   name: default
//   namespace: foo
// spec:
//   selector:
//     matchLabels:
//       app: finance
//   mtls:
//     mode: STRICT
// ```
// Policy to allow mTLS strict for all workloads, but leave port 8080 to
// plaintext:
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: PeerAuthentication
// metadata:
//   name: default
//   namespace: foo
// spec:
//   selector:
//     matchLabels:
//       app: finance
//   mtls:
//     mode: STRICT
//   portLevelMtls:
//     8080:
//       mode: DISABLE
// ```
// Policy to inherite mTLS mode from namespace (or mesh) settings, and overwrite
// settings for port 8080
// ```yaml
// apiVersion: security.istio.io/v1beta1
// kind: PeerAuthentication
// metadata:
//   name: default
//   namespace: foo
// spec:
//   selector:
//     matchLabels:
//       app: finance
//   mtls:
//     mode: UNSET
//   portLevelMtls:
//     8080:
//       mode: DISABLE
// ```
type PeerAuthenticationSpec struct {
	// The selector determines the workloads to apply the ChannelAuthentication on.
	// If not set, the policy will be applied to all workloads in the same namespace as the policy.
	Selector *selector.WorkloadSelector `json:"selector,omitempty"`
	// Mutual TLS settings for workload. If not defined, inherit from parent.
	Mtls *PeerAuthenticationMTLS `json:"mtls,omitempty"`
	// Port specific mutual TLS settings.
	PortLevelMtls map[uint32]*PeerAuthenticationMTLS `json:"portLevelMtls,omitempty"`
}

// Mutual TLS settings.
type PeerAuthenticationMTLS struct {
	// Defines the mTLS mode used for peer authentication.
	Mode MTLSMode `json:"mode,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// PeerAuthenticationList is a list of PeerAuthentication resources
type PeerAuthenticationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []PeerAuthentication `json:"items"`
}
