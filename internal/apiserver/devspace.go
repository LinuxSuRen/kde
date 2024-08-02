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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	"github.com/linuxsuren/kde/config"
	kdeClient "github.com/linuxsuren/kde/pkg/client/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type Server struct {
	Client    *kubernetes.Clientset
	KClient   *kdeClient.Clientset
	DClient   dynamic.Interface
	ExtClient *apiextensionsclientset.Clientset
}

func (s *Server) CreateDevSpace(c *gin.Context) {
	fmt.Println("CreateDevSpace")
	// clone git repo

	// find the config file

	// create the devspace resource
	devSpace := &v1alpha1.DevSpace{}
	if err := c.BindJSON(devSpace); err != nil {
		c.Error(err)
	} else {
		result, err := s.KClient.LinuxsurenV1alpha1().DevSpaces("default").Create(c.Request.Context(), devSpace, metav1.CreateOptions{})
		if err != nil {
			c.Error(err)
		} else {
			c.JSON(http.StatusOK, result)
		}
	}

	// query the status of the devspace resource

	// return the space address
}

func (s *Server) ListDevSpace(c *gin.Context) {
	result, err := s.KClient.LinuxsurenV1alpha1().DevSpaces("default").List(c.Request.Context(), metav1.ListOptions{})
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (s *Server) DeleteDevSpace(c *gin.Context) {
	name := c.Params.ByName("devspace")
	err := s.KClient.LinuxsurenV1alpha1().DevSpaces("default").Delete(c.Request.Context(), name, metav1.DeleteOptions{})
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func (s *Server) UpdateDevSpace(c *gin.Context) {
	name := c.Params.ByName("devspace")
	devSpace := &v1alpha1.DevSpace{}
	devSpace.Name = name
	if err := c.BindJSON(devSpace); err != nil {
		c.Error(err)
	} else {
		result, err := s.KClient.LinuxsurenV1alpha1().DevSpaces("default").Update(c.Request.Context(), devSpace, metav1.UpdateOptions{})
		if err != nil {
			c.Error(err)
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func (s *Server) GetDevSpace(c *gin.Context) {
	name := c.Params.ByName("devspace")
	result, err := s.KClient.LinuxsurenV1alpha1().DevSpaces("default").Get(c.Request.Context(), name, metav1.GetOptions{})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (s *Server) Install(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.DefaultQuery("namespace", "default")
	deploy := getDeployment()

	_, err := s.Client.AppsV1().Deployments(namespace).Create(ctx, deploy, metav1.CreateOptions{})

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func getDeployment() *appsv1.Deployment {
	var err error
	data, _ := config.GetFile("manager/manager.yaml")

	deploy := &appsv1.Deployment{}
	if data, err = yaml.ToJSON(data); err == nil {
		yaml.Unmarshal(data, deploy)
	}
	deploy.SetName("kde-controller")
	return deploy
}

type InstanceStatus struct {
	Component string `json:"component"`
	Status    string `json:"status"`
}

func (s *Server) InstanceStatus(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.DefaultQuery("namespace", "default")
	instanceStatus := []InstanceStatus{}

	_, err := s.ExtClient.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, "devspaces.linuxsuren.github.io", metav1.GetOptions{})
	if err == nil {
		instanceStatus = append(instanceStatus, InstanceStatus{
			Component: "crd",
			Status:    "installed",
		})
	} else {
		instanceStatus = append(instanceStatus, InstanceStatus{
			Component: "crd",
			Status:    "not installed",
		})
	}

	controllerDeploy := getDeployment()
	if controllerDeploy, err = s.Client.AppsV1().Deployments(namespace).Get(ctx, controllerDeploy.Name, metav1.GetOptions{}); err == nil {
		instanceStatus = append(instanceStatus, InstanceStatus{
			Component: "controller",
			Status:    controllerDeploy.Status.String(),
		})
	} else {
		instanceStatus = append(instanceStatus, InstanceStatus{
			Component: "controller",
			Status:    "not installed",
		})
	}

	c.JSON(http.StatusOK, instanceStatus)
}

func (s *Server) Namespaces(c *gin.Context) {
	nsList, err := s.Client.CoreV1().Namespaces().List(c.Request.Context(), metav1.ListOptions{})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, nsList)
	}
}
