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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DevSpaceSpec defines the desired state of DevSpace
type DevSpaceSpec struct {
	// CPU is the CPU limit
	// +kubebuilder:default:cpu="2"
	CPU string `json:"cpu,omitempty"`
	// Memory is the memory limit
	// +kubebuilder:default:memory="4Gi"
	Memory string `json:"memory,omitempty"`
	// Storage is the storage size
	// +kubebuilder:default:storage="50Gi"
	Storage string `json:"storage,omitempty"`
	Image   string `json:"image,omitempty"`
	// Replicas is the number of replicas
	// +kubebuilder:default:replicas=1
	Replicas int32 `json:"replicas,omitempty"`
	// Host is the hostname
	Host        string            `json:"host,omitempty"`
	Repository  *GitRepository    `json:"repository,omitempty"`
	Auth        DevSpaceAuth      `json:"auth,omitempty"`
	Environment map[string]string `json:"env,omitempty"`
	HostAliases []v1.HostAlias    `json:"hostAliases,omitempty"`
	Windows     []Window          `json:"windows,omitempty"`
	InitScript  string            `json:"initScript,omitempty"`
	Services    Services          `json:"services,omitempty"`
}

type Services struct {
	Docker *Docker `json:"docker,omitempty"`
	MySQL  *MySQL  `json:"mysql,omitempty"`
	Redis  *Redis  `json:"redis,omitempty"`
}

type Docker struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type MySQL struct {
	Enabled  bool   `json:"enabled,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	Image    string `json:"image,omitempty"`
}

type Redis struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Window struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type DevSpaceAuth struct {
	BasicAuth *BasicAuth `json:"basicAuth,omitempty"`
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type GitRepository struct {
	URL string `json:"url"`
	// Branch is the branch
	// +kubebuilder:default:branch="master"
	Branch string `json:"branch,omitempty"`
	// Username is the username
	Username string `json:"username"`
	// Password is the password
	Password string `json:"password,omitempty"`
}

// DevSpaceStatus defines the observed state of DevSpace
type DevSpaceStatus struct {
	Link         string                    `json:"link,omitempty"`
	ExposeLinks  []ExposeLink              `json:"exposeLinks,omitempty"`
	DeployStatus string                    `json:"deployStatus,omitempty"`
	Pods         []v1.LocalObjectReference `json:"pods,omitempty"`
	Phase        DevSpacePhase             `json:"phase,omitempty"`
}

type DevSpacePhase string

const (
	DevSpacePhaseReady DevSpacePhase = "Ready"
	DevSpacePhaseOff   DevSpacePhase = "Off"
)

type ExposeLink struct {
	Link string `json:"link,omitempty"`
	Port int    `json:"port,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Link",type=string,JSONPath=`.status.link`
// +kubebuilder:printcolumn:name="DeployStatus",type=string,JSONPath=`.status.deployStatus`
// +kubebuilder:printcolumn:name="Replicas",type=integer,JSONPath=`.spec.replicas`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:subresource:status

// DevSpace is the Schema for the devspaces API
type DevSpace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevSpaceSpec   `json:"spec,omitempty"`
	Status DevSpaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DevSpaceList contains a list of DevSpace
type DevSpaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevSpace `json:"items"`
}

const (
	AnnoKeyExposePorts      = "linuxsuren.github.io/expose-ports"
	AnnoKeyWebhookToken     = "linuxsuren.github.io/webhook-token"
	AnnoKeyImagePullPolicy  = "linuxsuren.github.io/image-pull-policy"
	AnnoKeyStorageClassName = "linuxsuren.github.io/storage-class-name"
	AnnoKeyVolumeAccessMode = "linuxsuren.github.io/volume-access-mode"
	AnnoKeyVolumeMode       = "linuxsuren.github.io/volume-mode"
	AnnoKeyIngressMode      = "linuxsuren.github.io/ingress-mode"
	AnnoKeyBasicAuth        = "linuxsuren.github.io/basic-auth"
	AnnoKeyMaintainMode     = "linuxsuren.github.io/maintain-mode"
)

func init() {
	SchemeBuilder.Register(&DevSpace{}, &DevSpaceList{})
}
