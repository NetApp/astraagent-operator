//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Astra) DeepCopyInto(out *Astra) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Astra.
func (in *Astra) DeepCopy() *Astra {
	if in == nil {
		return nil
	}
	out := new(Astra)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AstraAgent) DeepCopyInto(out *AstraAgent) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AstraAgent.
func (in *AstraAgent) DeepCopy() *AstraAgent {
	if in == nil {
		return nil
	}
	out := new(AstraAgent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AstraAgent) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AstraAgentList) DeepCopyInto(out *AstraAgentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AstraAgent, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AstraAgentList.
func (in *AstraAgentList) DeepCopy() *AstraAgentList {
	if in == nil {
		return nil
	}
	out := new(AstraAgentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AstraAgentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AstraAgentSpec) DeepCopyInto(out *AstraAgentSpec) {
	*out = *in
	out.NatssyncClient = in.NatssyncClient
	out.HttpProxyClient = in.HttpProxyClient
	out.EchoClient = in.EchoClient
	out.Nats = in.Nats
	out.ConfigMap = in.ConfigMap
	out.Astra = in.Astra
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AstraAgentSpec.
func (in *AstraAgentSpec) DeepCopy() *AstraAgentSpec {
	if in == nil {
		return nil
	}
	out := new(AstraAgentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AstraAgentStatus) DeepCopyInto(out *AstraAgentStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.NatssyncClient = in.NatssyncClient
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AstraAgentStatus.
func (in *AstraAgentStatus) DeepCopy() *AstraAgentStatus {
	if in == nil {
		return nil
	}
	out := new(AstraAgentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigMap) DeepCopyInto(out *ConfigMap) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigMap.
func (in *ConfigMap) DeepCopy() *ConfigMap {
	if in == nil {
		return nil
	}
	out := new(ConfigMap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EchoClient) DeepCopyInto(out *EchoClient) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EchoClient.
func (in *EchoClient) DeepCopy() *EchoClient {
	if in == nil {
		return nil
	}
	out := new(EchoClient)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HttpProxyClient) DeepCopyInto(out *HttpProxyClient) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HttpProxyClient.
func (in *HttpProxyClient) DeepCopy() *HttpProxyClient {
	if in == nil {
		return nil
	}
	out := new(HttpProxyClient)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Nats) DeepCopyInto(out *Nats) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Nats.
func (in *Nats) DeepCopy() *Nats {
	if in == nil {
		return nil
	}
	out := new(Nats)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NatssyncClient) DeepCopyInto(out *NatssyncClient) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NatssyncClient.
func (in *NatssyncClient) DeepCopy() *NatssyncClient {
	if in == nil {
		return nil
	}
	out := new(NatssyncClient)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NatssyncClientStatus) DeepCopyInto(out *NatssyncClientStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NatssyncClientStatus.
func (in *NatssyncClientStatus) DeepCopy() *NatssyncClientStatus {
	if in == nil {
		return nil
	}
	out := new(NatssyncClientStatus)
	in.DeepCopyInto(out)
	return out
}
