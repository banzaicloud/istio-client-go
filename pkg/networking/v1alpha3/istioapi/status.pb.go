// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: meta/v1alpha1/status.proto

package v1alpha1

import (
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
)

type IstioStatus struct {
	// Current service state of pod.
	// More info: https://istio.io/docs/reference/config/config-status/
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []*IstioCondition `protobuf:"bytes,1,rep,name=conditions,proto3" json:"conditions,omitempty"`
	// Resource Generation to which the Reconciled Condition refers.
	// When this value is not equal to the object's metadata generation, reconciled condition  calculation for the current
	// generation is still in progress.  See https://istio.io/latest/docs/reference/config/config-status/ for more info.
	// +optional
	ObservedGeneration   int64    `protobuf:"varint,2,opt,name=observed_generation,json=observedGeneration,proto3" json:"observed_generation,omitempty"`
}

func (m *IstioStatus) Reset()         { *m = IstioStatus{} }
func (m *IstioStatus) String() string { return proto.CompactTextString(m) }
func (*IstioStatus) ProtoMessage()    {}

type IstioCondition struct {
	// Type is the type of the condition.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	// Last time we probed the condition.
	// +optional
	LastProbeTime *types.Timestamp `protobuf:"bytes,3,opt,name=last_probe_time,json=lastProbeTime,proto3" json:"last_probe_time,omitempty"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime *types.Timestamp `protobuf:"bytes,4,opt,name=last_transition_time,json=lastTransitionTime,proto3" json:"last_transition_time,omitempty"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `protobuf:"bytes,5,opt,name=reason,proto3" json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message              string   `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *IstioCondition) Reset()         { *m = IstioCondition{} }
func (m *IstioCondition) String() string { return proto.CompactTextString(m) }
func (*IstioCondition) ProtoMessage()    {}
