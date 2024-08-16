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
	"context"
	"embed"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	kdeClient "github.com/linuxsuren/kde/pkg/client/clientset/versioned"
	"github.com/linuxsuren/kde/pkg/core"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type Server struct {
	Client          *kubernetes.Clientset
	KClient         kdeClient.Interface
	DClient         dynamic.Interface
	ExtClient       *apiextensionsclientset.Clientset
	SystemNamespace string
}

func (s *Server) CreateDevSpace(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := getNamespaceFromQuery(c)
	// clone git repo

	// find the config file

	// create the devspace resource
	devSpace := &v1alpha1.DevSpace{}
	if err := c.BindJSON(devSpace); err != nil {
		c.Error(err)
	} else {
		config, err := core.GetConfigFromConfigMap(ctx, s.Client.CoreV1().ConfigMaps(namespace), "config")
		if err != nil {
			c.Error(err)
		}
		setDefaultConfig(devSpace, config)

		result, err := s.KClient.LinuxsurenV1alpha1().DevSpaces(namespace).Create(ctx, devSpace, metav1.CreateOptions{})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, result)
		}
	}

	// query the status of the devspace resource

	// return the space address
}

func (s *Server) ListDevSpace(c *gin.Context) {
	namespace := getNamespaceFromQuery(c)
	result, err := s.KClient.LinuxsurenV1alpha1().DevSpaces(namespace).List(c.Request.Context(), metav1.ListOptions{})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (s *Server) DeleteDevSpace(c *gin.Context) {
	name := c.Params.ByName("devspace")
	namespace := getNamespaceFromQuery(c)
	err := s.KClient.LinuxsurenV1alpha1().DevSpaces(namespace).Delete(c.Request.Context(), name, metav1.DeleteOptions{})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func (s *Server) UpdateDevSpace(c *gin.Context) {
	name := c.Params.ByName("devspace")
	namespace := getNamespaceFromQuery(c)
	devSpace := &v1alpha1.DevSpace{}
	devSpace.Name = name
	if err := c.BindJSON(devSpace); err != nil {
		c.Error(err)
	} else {
		result, err := s.KClient.LinuxsurenV1alpha1().DevSpaces(namespace).Update(c.Request.Context(), devSpace, metav1.UpdateOptions{})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func (s *Server) RestartDevSpace(c *gin.Context) {
	name := c.Params.ByName("devspace")
	namespace := getNamespaceFromQuery(c)

	err := s.updateReplicas(c.Request.Context(), namespace, name, 0)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = s.updateReplicas(c.Request.Context(), namespace, name, 1)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func (s *Server) SetDevSpaceReplicas(c *gin.Context) {
	name := c.Params.ByName("devspace")
	namespace := getNamespaceFromQuery(c)
	replicas := c.Query("replicas")

	var replicaNum int
	var err error
	if replicaNum, err = strconv.Atoi(replicas); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = s.updateReplicas(c.Request.Context(), namespace, name, int32(replicaNum))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func (s *Server) updateReplicas(ctx context.Context, namespace, name string, replicas int32) (err error) {
	var devSpace *v1alpha1.DevSpace
	devSpace, err = s.KClient.LinuxsurenV1alpha1().DevSpaces(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return
	}

	devSpace.Spec.Replicas = &replicas
	_, err = s.KClient.LinuxsurenV1alpha1().DevSpaces(namespace).Update(ctx, devSpace, metav1.UpdateOptions{})
	return
}

func (s *Server) GetDevSpace(c *gin.Context) {
	name := c.Params.ByName("devspace")
	namespace := getNamespaceFromQuery(c)
	result, err := s.KClient.LinuxsurenV1alpha1().DevSpaces(namespace).Get(c.Request.Context(), name, metav1.GetOptions{})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (s *Server) GetDevSpaceLanguages(c *gin.Context) {
	languagesData, err := embedFS.ReadFile("data/languages.json")
	sliceData := []interface{}{}
	if err == nil {
		err = json.Unmarshal(languagesData, &sliceData)
	}

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
	} else {
		cm := getConfigMap("config.yaml")
		cm.SetNamespace("default")
		ctx := c.Request.Context()

		if config, err := core.GetConfigFromConfigMap(ctx, s.Client.CoreV1().ConfigMaps(cm.GetNamespace()), cm.GetName()); err == nil {
			for _, lan := range config.Languages {
				lan.Name = strings.TrimSpace(lan.Name)
				lan.Image = strings.TrimSpace(lan.Image)
				if lan.Name == "" || lan.Image == "" {
					continue
				}

				sliceData = append(sliceData, map[string]interface{}{
					"name":  lan.Name,
					"image": lan.Image,
				})
			}
		}

		c.JSON(http.StatusOK, sliceData)
	}
}

//go:embed data/*
var embedFS embed.FS
