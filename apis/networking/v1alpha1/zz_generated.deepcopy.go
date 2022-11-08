//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2020 The Crossplane Authors.

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

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancer) DeepCopyInto(out *LoadBalancer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancer.
func (in *LoadBalancer) DeepCopy() *LoadBalancer {
	if in == nil {
		return nil
	}
	out := new(LoadBalancer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LoadBalancer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerAlgorithm) DeepCopyInto(out *LoadBalancerAlgorithm) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerAlgorithm.
func (in *LoadBalancerAlgorithm) DeepCopy() *LoadBalancerAlgorithm {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerAlgorithm)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerList) DeepCopyInto(out *LoadBalancerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LoadBalancer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerList.
func (in *LoadBalancerList) DeepCopy() *LoadBalancerList {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LoadBalancerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerObservation) DeepCopyInto(out *LoadBalancerObservation) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerObservation.
func (in *LoadBalancerObservation) DeepCopy() *LoadBalancerObservation {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerParameters) DeepCopyInto(out *LoadBalancerParameters) {
	*out = *in
	out.Algorithm = in.Algorithm
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = new(map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
	if in.Location != nil {
		in, out := &in.Location, &out.Location
		*out = new(string)
		**out = **in
	}
	if in.Network != nil {
		in, out := &in.Network, &out.Network
		*out = new(uint64)
		**out = **in
	}
	if in.NetworkZone != nil {
		in, out := &in.NetworkZone, &out.NetworkZone
		*out = new(string)
		**out = **in
	}
	if in.PublicInterface != nil {
		in, out := &in.PublicInterface, &out.PublicInterface
		*out = new(bool)
		**out = **in
	}
	if in.Services != nil {
		in, out := &in.Services, &out.Services
		*out = new([]LoadBalancerService)
		if **in != nil {
			in, out := *in, *out
			*out = make([]LoadBalancerService, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
	if in.Targets != nil {
		in, out := &in.Targets, &out.Targets
		*out = new([]LoadBalancerTarget)
		if **in != nil {
			in, out := *in, *out
			*out = make([]LoadBalancerTarget, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerParameters.
func (in *LoadBalancerParameters) DeepCopy() *LoadBalancerParameters {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerService) DeepCopyInto(out *LoadBalancerService) {
	*out = *in
	in.HTTP.DeepCopyInto(&out.HTTP)
	in.HealthCheck.DeepCopyInto(&out.HealthCheck)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerService.
func (in *LoadBalancerService) DeepCopy() *LoadBalancerService {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerServiceHTTP) DeepCopyInto(out *LoadBalancerServiceHTTP) {
	*out = *in
	if in.CookieName != nil {
		in, out := &in.CookieName, &out.CookieName
		*out = new(string)
		**out = **in
	}
	if in.CookieLifetime != nil {
		in, out := &in.CookieLifetime, &out.CookieLifetime
		*out = new(int)
		**out = **in
	}
	if in.Certificates != nil {
		in, out := &in.Certificates, &out.Certificates
		*out = new([]int)
		if **in != nil {
			in, out := *in, *out
			*out = make([]int, len(*in))
			copy(*out, *in)
		}
	}
	if in.RedirectHTTP != nil {
		in, out := &in.RedirectHTTP, &out.RedirectHTTP
		*out = new(bool)
		**out = **in
	}
	if in.StickySessions != nil {
		in, out := &in.StickySessions, &out.StickySessions
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerServiceHTTP.
func (in *LoadBalancerServiceHTTP) DeepCopy() *LoadBalancerServiceHTTP {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerServiceHTTP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerServiceHealthCheck) DeepCopyInto(out *LoadBalancerServiceHealthCheck) {
	*out = *in
	if in.HTTP != nil {
		in, out := &in.HTTP, &out.HTTP
		*out = new(LoadBalancerServiceHealthCheckHTTP)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerServiceHealthCheck.
func (in *LoadBalancerServiceHealthCheck) DeepCopy() *LoadBalancerServiceHealthCheck {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerServiceHealthCheck)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerServiceHealthCheckHTTP) DeepCopyInto(out *LoadBalancerServiceHealthCheckHTTP) {
	*out = *in
	if in.Response != nil {
		in, out := &in.Response, &out.Response
		*out = new(string)
		**out = **in
	}
	if in.StatusCodes != nil {
		in, out := &in.StatusCodes, &out.StatusCodes
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerServiceHealthCheckHTTP.
func (in *LoadBalancerServiceHealthCheckHTTP) DeepCopy() *LoadBalancerServiceHealthCheckHTTP {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerServiceHealthCheckHTTP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerSpec) DeepCopyInto(out *LoadBalancerSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerSpec.
func (in *LoadBalancerSpec) DeepCopy() *LoadBalancerSpec {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerStatus) DeepCopyInto(out *LoadBalancerStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	out.AtProvider = in.AtProvider
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerStatus.
func (in *LoadBalancerStatus) DeepCopy() *LoadBalancerStatus {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerTarget) DeepCopyInto(out *LoadBalancerTarget) {
	*out = *in
	if in.Server != nil {
		in, out := &in.Server, &out.Server
		*out = new(LoadBalancerTargetServer)
		**out = **in
	}
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(LoadBalancerTargetLabelSelector)
		**out = **in
	}
	if in.IP != nil {
		in, out := &in.IP, &out.IP
		*out = new(LoadBalancerTargetIP)
		**out = **in
	}
	if in.HealthStatus != nil {
		in, out := &in.HealthStatus, &out.HealthStatus
		*out = new([]LoadBalancerTargetHealthStatus)
		if **in != nil {
			in, out := *in, *out
			*out = make([]LoadBalancerTargetHealthStatus, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
	if in.UsePrivateIP != nil {
		in, out := &in.UsePrivateIP, &out.UsePrivateIP
		*out = new(bool)
		**out = **in
	}
	if in.Targets != nil {
		in, out := &in.Targets, &out.Targets
		*out = new([]LoadBalancerTarget)
		if **in != nil {
			in, out := *in, *out
			*out = make([]LoadBalancerTarget, len(*in))
			for i := range *in {
				(*in)[i].DeepCopyInto(&(*out)[i])
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerTarget.
func (in *LoadBalancerTarget) DeepCopy() *LoadBalancerTarget {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerTargetHealthStatus) DeepCopyInto(out *LoadBalancerTargetHealthStatus) {
	*out = *in
	if in.ListenPort != nil {
		in, out := &in.ListenPort, &out.ListenPort
		*out = new(int)
		**out = **in
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerTargetHealthStatus.
func (in *LoadBalancerTargetHealthStatus) DeepCopy() *LoadBalancerTargetHealthStatus {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerTargetHealthStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerTargetIP) DeepCopyInto(out *LoadBalancerTargetIP) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerTargetIP.
func (in *LoadBalancerTargetIP) DeepCopy() *LoadBalancerTargetIP {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerTargetIP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerTargetLabelSelector) DeepCopyInto(out *LoadBalancerTargetLabelSelector) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerTargetLabelSelector.
func (in *LoadBalancerTargetLabelSelector) DeepCopy() *LoadBalancerTargetLabelSelector {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerTargetLabelSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerTargetServer) DeepCopyInto(out *LoadBalancerTargetServer) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerTargetServer.
func (in *LoadBalancerTargetServer) DeepCopy() *LoadBalancerTargetServer {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerTargetServer)
	in.DeepCopyInto(out)
	return out
}
