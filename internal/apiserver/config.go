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

package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	"github.com/linuxsuren/kde/pkg/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Server) GetConfig(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := getNamespaceFromQuery(c)

	cm := getConfigMap("config.yaml")
	cm.SetNamespace(namespace)

	if config, err := core.GetConfigFromConfigMap(ctx, s.Client.CoreV1().ConfigMaps(namespace), cm.GetName()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, config)
	}
}

func (s *Server) UpdateConfig(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := getNamespaceFromQuery(c)

	var err error
	config := &core.Config{}
	if err = c.BindJSON(config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    // clean invalid languages
    config.Languages = core.CleanInvalidLanguages(config.Languages)

	cm := getConfigMap("config.yaml")
	cm.SetNamespace(namespace)

	cm, err = s.Client.CoreV1().ConfigMaps(namespace).Get(ctx, cm.GetName(), metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var data []byte
	if data, err = config.ToJSON(); err == nil {
		cm.Data[core.ConfigFileName] = string(data)
		_, err = s.Client.CoreV1().ConfigMaps(namespace).Update(ctx, cm, metav1.UpdateOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

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
