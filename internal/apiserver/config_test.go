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
	"testing"

	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	"github.com/linuxsuren/kde/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaultConfig(t *testing.T) {
	verify := func(t *testing.T, devspace *v1alpha1.DevSpace, config *core.Config) {
		setDefaultConfig(devspace, config)

		assert.Equal(t, "Always", devspace.Annotations[v1alpha1.AnnoKeyImagePullPolicy])
		assert.Equal(t, "standard", devspace.Annotations[v1alpha1.AnnoKeyStorageClassName])
		assert.Equal(t, "ReadWriteOnce", devspace.Annotations[v1alpha1.AnnoKeyVolumeAccessMode])
		assert.Equal(t, "Filesystem", devspace.Annotations[v1alpha1.AnnoKeyVolumeMode])
		assert.Equal(t, "nginx", devspace.Annotations[v1alpha1.AnnoKeyIngressMode])
	}

	t.Run("annotations are empty", func(t *testing.T) {
		devspace := &v1alpha1.DevSpace{}
		config := &core.Config{
			ImagePullPolicy:  "Always",
			StorageClassName: "standard",
			VolumeAccessMode: "ReadWriteOnce",
			VolumeMode:       "Filesystem",
			IngressMode:      "nginx",
		}
		verify(t, devspace, config)
	})

	t.Run("annotations are not empty", func(t *testing.T) {
		devspace := &v1alpha1.DevSpace{}
		devspace.Annotations = map[string]string{
			v1alpha1.AnnoKeyImagePullPolicy:  "Always",
			v1alpha1.AnnoKeyStorageClassName: "standard",
			v1alpha1.AnnoKeyVolumeAccessMode: "ReadWriteOnce",
			v1alpha1.AnnoKeyVolumeMode:       "Filesystem",
			v1alpha1.AnnoKeyIngressMode:      "nginx",
		}
		config := &core.Config{
			ImagePullPolicy:  "Never",
			StorageClassName: "fast",
			VolumeAccessMode: "ReadWriteMany",
			VolumeMode:       "Block",
			IngressMode:      "traefik",
		}
		verify(t, devspace, config)
	})
}
