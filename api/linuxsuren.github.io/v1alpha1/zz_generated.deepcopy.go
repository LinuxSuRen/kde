//go:build !ignore_autogenerated

/*
Copyright 2024 kde authors.

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
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BasicAuth) DeepCopyInto(out *BasicAuth) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BasicAuth.
func (in *BasicAuth) DeepCopy() *BasicAuth {
	if in == nil {
		return nil
	}
	out := new(BasicAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevSpace) DeepCopyInto(out *DevSpace) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevSpace.
func (in *DevSpace) DeepCopy() *DevSpace {
	if in == nil {
		return nil
	}
	out := new(DevSpace)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DevSpace) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevSpaceAuth) DeepCopyInto(out *DevSpaceAuth) {
	*out = *in
	if in.BasicAuth != nil {
		in, out := &in.BasicAuth, &out.BasicAuth
		*out = new(BasicAuth)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevSpaceAuth.
func (in *DevSpaceAuth) DeepCopy() *DevSpaceAuth {
	if in == nil {
		return nil
	}
	out := new(DevSpaceAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevSpaceList) DeepCopyInto(out *DevSpaceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DevSpace, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevSpaceList.
func (in *DevSpaceList) DeepCopy() *DevSpaceList {
	if in == nil {
		return nil
	}
	out := new(DevSpaceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DevSpaceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevSpaceSpec) DeepCopyInto(out *DevSpaceSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Repository != nil {
		in, out := &in.Repository, &out.Repository
		*out = new(GitRepository)
		**out = **in
	}
	in.Auth.DeepCopyInto(&out.Auth)
	if in.Environment != nil {
		in, out := &in.Environment, &out.Environment
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.HostAliases != nil {
		in, out := &in.HostAliases, &out.HostAliases
		*out = make([]v1.HostAlias, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Windows != nil {
		in, out := &in.Windows, &out.Windows
		*out = make([]Window, len(*in))
		copy(*out, *in)
	}
	in.Services.DeepCopyInto(&out.Services)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevSpaceSpec.
func (in *DevSpaceSpec) DeepCopy() *DevSpaceSpec {
	if in == nil {
		return nil
	}
	out := new(DevSpaceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevSpaceStatus) DeepCopyInto(out *DevSpaceStatus) {
	*out = *in
	if in.ExposeLinks != nil {
		in, out := &in.ExposeLinks, &out.ExposeLinks
		*out = make([]ExposeLink, len(*in))
		copy(*out, *in)
	}
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevSpaceStatus.
func (in *DevSpaceStatus) DeepCopy() *DevSpaceStatus {
	if in == nil {
		return nil
	}
	out := new(DevSpaceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Docker) DeepCopyInto(out *Docker) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Docker.
func (in *Docker) DeepCopy() *Docker {
	if in == nil {
		return nil
	}
	out := new(Docker)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExposeLink) DeepCopyInto(out *ExposeLink) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExposeLink.
func (in *ExposeLink) DeepCopy() *ExposeLink {
	if in == nil {
		return nil
	}
	out := new(ExposeLink)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitRepository) DeepCopyInto(out *GitRepository) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitRepository.
func (in *GitRepository) DeepCopy() *GitRepository {
	if in == nil {
		return nil
	}
	out := new(GitRepository)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MySQL) DeepCopyInto(out *MySQL) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MySQL.
func (in *MySQL) DeepCopy() *MySQL {
	if in == nil {
		return nil
	}
	out := new(MySQL)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Postgres) DeepCopyInto(out *Postgres) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Postgres.
func (in *Postgres) DeepCopy() *Postgres {
	if in == nil {
		return nil
	}
	out := new(Postgres)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RabbitMQ) DeepCopyInto(out *RabbitMQ) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RabbitMQ.
func (in *RabbitMQ) DeepCopy() *RabbitMQ {
	if in == nil {
		return nil
	}
	out := new(RabbitMQ)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Redis) DeepCopyInto(out *Redis) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Redis.
func (in *Redis) DeepCopy() *Redis {
	if in == nil {
		return nil
	}
	out := new(Redis)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Services) DeepCopyInto(out *Services) {
	*out = *in
	if in.Docker != nil {
		in, out := &in.Docker, &out.Docker
		*out = new(Docker)
		**out = **in
	}
	if in.MySQL != nil {
		in, out := &in.MySQL, &out.MySQL
		*out = new(MySQL)
		**out = **in
	}
	if in.Postgres != nil {
		in, out := &in.Postgres, &out.Postgres
		*out = new(Postgres)
		**out = **in
	}
	if in.TDEngine != nil {
		in, out := &in.TDEngine, &out.TDEngine
		*out = new(TDEngine)
		**out = **in
	}
	if in.RabbitMQ != nil {
		in, out := &in.RabbitMQ, &out.RabbitMQ
		*out = new(RabbitMQ)
		**out = **in
	}
	if in.Redis != nil {
		in, out := &in.Redis, &out.Redis
		*out = new(Redis)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Services.
func (in *Services) DeepCopy() *Services {
	if in == nil {
		return nil
	}
	out := new(Services)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TDEngine) DeepCopyInto(out *TDEngine) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TDEngine.
func (in *TDEngine) DeepCopy() *TDEngine {
	if in == nil {
		return nil
	}
	out := new(TDEngine)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *User) DeepCopyInto(out *User) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new User.
func (in *User) DeepCopy() *User {
	if in == nil {
		return nil
	}
	out := new(User)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *User) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UserList) DeepCopyInto(out *UserList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]User, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UserList.
func (in *UserList) DeepCopy() *UserList {
	if in == nil {
		return nil
	}
	out := new(UserList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *UserList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UserSpec) DeepCopyInto(out *UserSpec) {
	*out = *in
	if in.Groups != nil {
		in, out := &in.Groups, &out.Groups
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.ResourceQuota.DeepCopyInto(&out.ResourceQuota)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UserSpec.
func (in *UserSpec) DeepCopy() *UserSpec {
	if in == nil {
		return nil
	}
	out := new(UserSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UserStatus) DeepCopyInto(out *UserStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UserStatus.
func (in *UserStatus) DeepCopy() *UserStatus {
	if in == nil {
		return nil
	}
	out := new(UserStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Window) DeepCopyInto(out *Window) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Window.
func (in *Window) DeepCopy() *Window {
	if in == nil {
		return nil
	}
	out := new(Window)
	in.DeepCopyInto(out)
	return out
}
