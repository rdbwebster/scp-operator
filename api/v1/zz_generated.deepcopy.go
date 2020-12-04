// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedOperator) DeepCopyInto(out *ManagedOperator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedOperator.
func (in *ManagedOperator) DeepCopy() *ManagedOperator {
	if in == nil {
		return nil
	}
	out := new(ManagedOperator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedOperator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedOperatorList) DeepCopyInto(out *ManagedOperatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ManagedOperator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedOperatorList.
func (in *ManagedOperatorList) DeepCopy() *ManagedOperatorList {
	if in == nil {
		return nil
	}
	out := new(ManagedOperatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedOperatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedOperatorSpec) DeepCopyInto(out *ManagedOperatorSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedOperatorSpec.
func (in *ManagedOperatorSpec) DeepCopy() *ManagedOperatorSpec {
	if in == nil {
		return nil
	}
	out := new(ManagedOperatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedOperatorStatus) DeepCopyInto(out *ManagedOperatorStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedOperatorStatus.
func (in *ManagedOperatorStatus) DeepCopy() *ManagedOperatorStatus {
	if in == nil {
		return nil
	}
	out := new(ManagedOperatorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SCPcluster) DeepCopyInto(out *SCPcluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SCPcluster.
func (in *SCPcluster) DeepCopy() *SCPcluster {
	if in == nil {
		return nil
	}
	out := new(SCPcluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SCPcluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SCPclusterList) DeepCopyInto(out *SCPclusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SCPcluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SCPclusterList.
func (in *SCPclusterList) DeepCopy() *SCPclusterList {
	if in == nil {
		return nil
	}
	out := new(SCPclusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SCPclusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SCPclusterSpec) DeepCopyInto(out *SCPclusterSpec) {
	*out = *in
	in.Connected.DeepCopyInto(&out.Connected)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SCPclusterSpec.
func (in *SCPclusterSpec) DeepCopy() *SCPclusterSpec {
	if in == nil {
		return nil
	}
	out := new(SCPclusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SCPclusterStatus) DeepCopyInto(out *SCPclusterStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SCPclusterStatus.
func (in *SCPclusterStatus) DeepCopy() *SCPclusterStatus {
	if in == nil {
		return nil
	}
	out := new(SCPclusterStatus)
	in.DeepCopyInto(out)
	return out
}