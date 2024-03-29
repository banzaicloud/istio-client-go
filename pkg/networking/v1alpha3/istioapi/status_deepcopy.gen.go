// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: meta/v1alpha1/status.proto

package v1alpha1

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// DeepCopyInto supports using IstioStatus within kubernetes types, where deepcopy-gen is used.
func (in *IstioStatus) DeepCopyInto(out *IstioStatus) {
	p := proto.Clone(in).(*IstioStatus)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IstioStatus. Required by controller-gen.
func (in *IstioStatus) DeepCopy() *IstioStatus {
	if in == nil {
		return nil
	}
	out := new(IstioStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new IstioStatus. Required by controller-gen.
func (in *IstioStatus) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using IstioCondition within kubernetes types, where deepcopy-gen is used.
func (in *IstioCondition) DeepCopyInto(out *IstioCondition) {
	p := proto.Clone(in).(*IstioCondition)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IstioCondition. Required by controller-gen.
func (in *IstioCondition) DeepCopy() *IstioCondition {
	if in == nil {
		return nil
	}
	out := new(IstioCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new IstioCondition. Required by controller-gen.
func (in *IstioCondition) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
