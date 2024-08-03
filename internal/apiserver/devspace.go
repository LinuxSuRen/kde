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
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	"github.com/linuxsuren/kde/config"
	kdeClient "github.com/linuxsuren/kde/pkg/client/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
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
	if namespace == "" {
		namespace = "default"
	}

	sa := getServiceAccount("service_account.yaml")
	sa.SetNamespace(namespace)
	_, saErr := s.Client.CoreV1().ServiceAccounts(namespace).Create(ctx, sa, metav1.CreateOptions{})

	cm := getConfigMap("config.yaml")
	cm.SetNamespace(namespace)
	_, cmErr := s.Client.CoreV1().ConfigMaps(namespace).Create(ctx, cm, metav1.CreateOptions{})

	deploy := getDeployment("manager.yaml")
	deploy.SetNamespace(namespace)
	_, deployErr := s.Client.AppsV1().Deployments(namespace).Create(ctx, deploy, metav1.CreateOptions{})

	apiserverDeploy := getDeployment("apiserver-deploy.yaml")
	apiserverDeploy.SetNamespace(namespace)
	_, apiserverDeployErr := s.Client.AppsV1().Deployments(namespace).Create(ctx, apiserverDeploy, metav1.CreateOptions{})

	service := getService("apiserver-service.yaml")
	service.SetNamespace(namespace)
	_, serviceErr := s.Client.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})

	ingress := getIngress("ingress.yaml")
	ingress.SetNamespace(namespace)
	_, ingressErr := s.Client.NetworkingV1().Ingresses(namespace).Create(ctx, ingress, metav1.CreateOptions{})

	err := errors.Join(saErr, cmErr, deployErr, apiserverDeployErr, serviceErr, ingressErr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func (s *Server) Uninstall(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.DefaultQuery("namespace", "default")
	if namespace == "" {
		namespace = "default"
	}

	sa := getServiceAccount("service_account.yaml")
	saErr := s.Client.CoreV1().ServiceAccounts(namespace).Delete(ctx, sa.GetName(), metav1.DeleteOptions{})

	cm := getConfigMap("config.yaml")
	cmErr := s.Client.CoreV1().ConfigMaps(namespace).Delete(ctx, cm.GetName(), metav1.DeleteOptions{})

	deploy := getDeployment("manager.yaml")
	deployErr := s.Client.AppsV1().Deployments(namespace).Delete(ctx, deploy.GetName(), metav1.DeleteOptions{})

	apiserverDeploy := getDeployment("apiserver-deploy.yaml")
	apiserverDeployErr := s.Client.AppsV1().Deployments(namespace).Delete(ctx, apiserverDeploy.GetName(), metav1.DeleteOptions{})

	service := getService("apiserver-service.yaml")
	serviceErr := s.Client.CoreV1().Services(namespace).Delete(ctx, service.GetName(), metav1.DeleteOptions{})

	ingress := getIngress("ingress.yaml")
	ingressErr := s.Client.NetworkingV1().Ingresses(namespace).Delete(ctx, ingress.GetName(), metav1.DeleteOptions{})

	err := errors.Join(saErr, cmErr, deployErr, apiserverDeployErr, serviceErr, ingressErr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func getDeployment(name string) *appsv1.Deployment {
	var err error
	data, _ := config.GetFile(filepath.Join("manager", name))

	deploy := &appsv1.Deployment{}
	if data, err = yaml.ToJSON(data); err == nil {
		yaml.Unmarshal(data, deploy)
	}
	return deploy
}

func getService(name string) *corev1.Service {
	var err error
	data, _ := config.GetFile(filepath.Join("manager", name))
	svc := &corev1.Service{}
	if data, err = yaml.ToJSON(data); err == nil {
		yaml.Unmarshal(data, svc)
	}
	return svc
}

func getIngress(name string) *networkingv1.Ingress {
	var err error
	data, _ := config.GetFile(filepath.Join("manager", name))
	ingress := &networkingv1.Ingress{}
	if data, err = yaml.ToJSON(data); err == nil {
		yaml.Unmarshal(data, ingress)
	}
	return ingress
}

func getConfigMap(name string) *corev1.ConfigMap {
	var err error
	data, _ := config.GetFile(filepath.Join("manager", name))
	cm := &corev1.ConfigMap{}
	if data, err = yaml.ToJSON(data); err == nil {
		yaml.Unmarshal(data, cm)
	}
	return cm
}

func getServiceAccount(name string) *corev1.ServiceAccount {
	var err error
	data, _ := config.GetFile(filepath.Join("rbac", name))
	sa := &corev1.ServiceAccount{}
	if data, err = yaml.ToJSON(data); err == nil {
		yaml.Unmarshal(data, sa)
	}
	return sa
}

type InstanceStatus struct {
	Component string `json:"component"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}

func (s *Server) InstanceStatus(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.DefaultQuery("namespace", "default")
	if namespace == "" {
		namespace = "default"
	}
	instanceStatus := []InstanceStatus{}
	instanceStatus = append(instanceStatus,
		s.getCRDStatus(ctx, "devspaces.linuxsuren.github.io"),
		s.getCRDStatus(ctx, "users.linuxsuren.github.io"),
		s.getDeploymentStatus(ctx, namespace, "manager"),
		s.getDeploymentStatus(ctx, namespace, "apiserver-deploy"),
		s.getServiceStatus(ctx, namespace, "apiserver-service"),
		s.getConfigmapStatus(ctx, namespace, "config"),
		s.getIngressStatus(ctx, namespace, "apiserver"),
	)

	c.JSON(http.StatusOK, instanceStatus)
}

func (s *Server) getCRDStatus(ctx context.Context, name string) InstanceStatus {
	_, err := s.ExtClient.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		return InstanceStatus{
			Component: "CRD",
			Name:      name,
			Status:    "installed",
		}
	} else {
		return InstanceStatus{
			Component: "CRD",
			Name:      name,
			Status:    "not installed",
		}
	}
}

func (s *Server) getDeploymentStatus(ctx context.Context, namespace, name string) InstanceStatus {
	deploy := getDeployment(name + ".yaml")
	if _, err := s.Client.AppsV1().Deployments(namespace).Get(ctx, deploy.Name, metav1.GetOptions{}); err == nil {
		return InstanceStatus{
			Component: "Deployment",
			Name:      name,
			Status:    "installed",
		}
	} else {
		return InstanceStatus{
			Component: "Deployment",
			Name:      name,
			Status:    "not installed",
		}
	}
}

func (s *Server) getServiceStatus(ctx context.Context, namespace, name string) InstanceStatus {
	if _, err := s.Client.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{}); err == nil {
		return InstanceStatus{
			Component: "Service",
			Name:      name,
			Status:    "installed",
		}
	} else {
		return InstanceStatus{
			Component: "Service",
			Name:      name,
			Status:    "not installed",
		}
	}
}

func (s *Server) getConfigmapStatus(ctx context.Context, namespace, name string) InstanceStatus {
	if _, err := s.Client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{}); err == nil {
		return InstanceStatus{
			Component: "ConfigMap",
			Name:      name,
			Status:    "installed",
		}
	} else {
		return InstanceStatus{
			Component: "ConfigMap",
			Name:      name,
			Status:    "not installed",
		}
	}
}

// getIngressStatus
func (s *Server) getIngressStatus(ctx context.Context, namespace, name string) InstanceStatus {
	if _, err := s.Client.NetworkingV1().Ingresses(namespace).Get(ctx, name, metav1.GetOptions{}); err == nil {
		return InstanceStatus{
			Component: "Ingress",
			Name:      name,
			Status:    "installed",
		}
	} else {
		return InstanceStatus{
			Component: "Ingress",
			Name:      name,
			Status:    "not installed",
		}
	}
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
