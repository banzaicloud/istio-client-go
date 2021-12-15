// Copyright © 2021 Banzai Cloud
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// WorkloadGroup
type WorkloadGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WorkloadGroupSpec `json:"spec"`
}

// `WorkloadGroup` describes a collection of workload instances.
// It provides a specification that the workload instances can use to bootstrap
// their proxies, including the metadata and identity. It is only intended to
// be used with non-k8s workloads like Virtual Machines, and is meant to mimic
// the existing sidecar injection and deployment specification model used for
// Kubernetes workloads to bootstrap Istio proxies.
//
// The following example declares a workload group representing a collection
// of workloads that will be registered under `reviews` in namespace
// `bookinfo`. The set of labels will be associated with each workload
// instance during the bootstrap process, and the ports 3550 and 8080
// will be associated with the workload group and use service account `default`.
// `app.kubernetes.io/version` is just an arbitrary example of a label.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: WorkloadGroup
// metadata:
//   name: reviews
//   namespace: bookinfo
// spec:
//   metadata:
//     labels:
//       app.kubernetes.io/name: reviews
//       app.kubernetes.io/version: "1.3.4"
//   template:
//     ports:
//       grpc: 3550
//       http: 8080
//     serviceAccount: default
//   probe:
//     initialDelaySeconds: 5
//     timeoutSeconds: 3
//     periodSeconds: 4
//     successThreshold: 3
//     failureThreshold: 3
//     httpGet:
//      path: /foo/bar
//      host: 127.0.0.1
//      port: 3100
//      scheme: HTTPS
//      httpHeaders:
//      - name: Lit-Header
//        value: Im-The-Best
// ```

// NOTE: Use this page https://istio.io/latest/docs/reference/config/networking/workload-group as a reference to verify
// if fields are required.

type WorkloadGroupSpec struct {
	// Metadata that will be used for all corresponding `WorkloadEntries`.
	// User labels for a workload group should be set here in `metadata` rather than in `template`.
	Metadata *WorkloadGroupObjectMeta `json:"metadata,omitempty"`
	// REQUIRED. Template to be used for the generation of `WorkloadEntry` resources that belong to this `WorkloadGroup`.
	// Please note that `address` and `labels` fields should not be set in the template, and an empty `serviceAccount`
	// should default to `default`. The workload identities (mTLS certificates) will be bootstrapped using the
	// specified service account's token. Workload entries in this group will be in the same namespace as the
	// workload group, and inherit the labels and annotations from the above `metadata` field.
	Template *WorkloadEntry `json:"template"`
	// `ReadinessProbe` describes the configuration the user must provide for health-checking on their workload.
	// This configuration mirrors K8S in both syntax and logic for the most part.
	Probe *ReadinessProbe `json:"probe,omitempty"`
}

// WorkloadGroupObjectMeta describes metadata that will be attached to a `WorkloadEntry`.
// It is a subset of the supported Kubernetes metadata.
type WorkloadGroupObjectMeta struct {
	// Labels to attach.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations to attach.
	Annotations map[string]string `json:"annotations,omitempty"`
}

type ReadinessProbe struct {
	// Number of seconds after the container has started before readiness probes are initiated.
	InitialDelaySeconds int32 `json:"initial_delay_seconds,omitempty"`
	// Number of seconds after which the probe times out.
	// Defaults to 1 second. Minimum value is 1 second.
	TimeoutSeconds int32 `json:"timeout_seconds,omitempty"`
	// How often (in seconds) to perform the probe.
	// Default to 10 seconds. Minimum value is 1 second.
	PeriodSeconds int32 `json:"period_seconds,omitempty"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// Defaults to 1 second.
	SuccessThreshold int32 `json:"success_threshold,omitempty"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// Defaults to 3 seconds.
	FailureThreshold int32 `json:"failure_threshold,omitempty"`

	// Users can only provide one configuration for health-checks (tcp, http, exec),
	// and this is expressed as a one-of. All the other configuration values
	// hold true for any of the health-check methods.
	HTTPGet   *HTTPHealthCheckConfig `json:"http_get,omitempty"`
	TCPSocket *TCPHealthCheckConfig  `json:"tcp_socket,omitempty"`
	Exec      *ExecHealthCheckConfig `json:"exec,omitempty"`
}

type HTTPHealthCheckConfig struct {
	// Path to access on the HTTP server.
	Path string `json:"path,omitempty"`
	// REQUIRED. Port on which the endpoint lives.
	Port uint32 `json:"port"`
	// Host to connect to, defaults to the pod IP. You probably want to set
	// "Host" in httpHeaders instead.
	Host string `json:"host,omitempty"`
	// HTTP or HTTPS, defaults to HTTP.
	Scheme string `json:"scheme,omitempty"`
	// Headers the proxy will pass on to make the request.
	// Allows repeated headers.
	HTTPHeaders []*HTTPHeader `json:"http_headers,omitempty"`
}

type HTTPHeader struct {
	// The header field name
	Name string `json:"name,omitempty"`
	// The header field value
	Value string `json:"value,omitempty"`
}

type TCPHealthCheckConfig struct {
	// Host to connect to, defaults to localhost.
	Host string `json:"host,omitempty"`
	// REQUIRED. Port of host.
	Port uint32 `json:"port"`
}

type ExecHealthCheckConfig struct {
	// Command to run. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
	Command []string `json:"command,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// WorkloadGroupList is a list of WorkloadGroup resources
type WorkloadGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []WorkloadGroup `json:"items"`
}
