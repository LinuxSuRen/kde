/*
Copyright 2024 kde authrors.

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

package apiserver

import (
	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	"github.com/linuxsuren/kde/pkg/core"
)

func setDefaultConfig(devspace *v1alpha1.DevSpace, config *core.Config) {
	if devspace.Annotations == nil {
		devspace.Annotations = make(map[string]string)
	}
	if _, ok := devspace.Annotations[v1alpha1.AnnoKeyImagePullPolicy]; !ok {
		devspace.Annotations[v1alpha1.AnnoKeyImagePullPolicy] = config.ImagePullPolicy
	}
	if _, ok := devspace.Annotations[v1alpha1.AnnoKeyStorageClassName]; !ok {
		devspace.Annotations[v1alpha1.AnnoKeyStorageClassName] = config.StorageClassName
	}
	if _, ok := devspace.Annotations[v1alpha1.AnnoKeyVolumeAccessMode]; !ok {
		devspace.Annotations[v1alpha1.AnnoKeyVolumeAccessMode] = config.VolumeAccessMode
	}
	if _, ok := devspace.Annotations[v1alpha1.AnnoKeyVolumeMode]; !ok {
		devspace.Annotations[v1alpha1.AnnoKeyVolumeMode] = config.VolumeMode
	}
	if _, ok := devspace.Annotations[v1alpha1.AnnoKeyIngressMode]; !ok {
		devspace.Annotations[v1alpha1.AnnoKeyIngressMode] = config.IngressMode
	}
}
