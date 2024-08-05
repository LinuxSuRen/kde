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
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/linuxsuren/kde/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Querier interface {
	DefaultQuery(key, defaultValue string) string
}

func getNamespaceFromQuery(querier Querier) string {
	namespace := querier.DefaultQuery("namespace", "default")
	if namespace == "" {
		namespace = "default"
	}
	return namespace
}

func (s *Server) Install(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := getNamespaceFromQuery(c)

	crdDevSpace := getCRD("linuxsuren.github.io_devspaces.yaml")
	_, crdDevSpaceErr := s.ExtClient.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, crdDevSpace, metav1.CreateOptions{})

	crdUser := getCRD("linuxsuren.github.io_users.yaml")
	_, crdUserErr := s.ExtClient.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, crdUser, metav1.CreateOptions{})

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

	err := errors.Join(client.IgnoreAlreadyExists(crdDevSpaceErr), client.IgnoreAlreadyExists(crdUserErr),
		client.IgnoreAlreadyExists(saErr), client.IgnoreAlreadyExists(cmErr),
		client.IgnoreAlreadyExists(deployErr), client.IgnoreAlreadyExists(apiserverDeployErr),
		client.IgnoreAlreadyExists(serviceErr), client.IgnoreAlreadyExists(ingressErr))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func (s *Server) Uninstall(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := getNamespaceFromQuery(c)

	crdDevSpace := getCRD("linuxsuren.github.io_devspaces.yaml")
	crdDevSpaceErr := s.ExtClient.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, crdDevSpace.GetName(), metav1.DeleteOptions{})

	crdUser := getCRD("linuxsuren.github.io_users.yaml")
	crdUserErr := s.ExtClient.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, crdUser.GetName(), metav1.DeleteOptions{})

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

	err := errors.Join(client.IgnoreNotFound(crdDevSpaceErr), client.IgnoreNotFound(crdUserErr),
		client.IgnoreNotFound(saErr), client.IgnoreNotFound(cmErr), client.IgnoreNotFound(deployErr),
		client.IgnoreNotFound(apiserverDeployErr), client.IgnoreNotFound(serviceErr),
		client.IgnoreNotFound(ingressErr))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, "")
	}
}

func getCRD(name string) *extv1.CustomResourceDefinition {
	var err error
	data, _ := config.GetFile(filepath.Join("crd/bases", name))

	crd := &extv1.CustomResourceDefinition{}
	if data, err = yaml.ToJSON(data); err == nil {
		yaml.Unmarshal(data, crd)
	}
	return crd
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
	c.JSON(http.StatusOK, s.getInstanceStatus(ctx, namespace))
}

func (s *Server) InstanceStatusWS(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := getNamespaceFromQuery(c)

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	for {
		instanceStatus := s.getInstanceStatus(ctx, namespace)
		conn.WriteJSON(instanceStatus)
		time.Sleep(time.Second)
	}
}

func (s *Server) getInstanceStatus(ctx context.Context, namespace string) []InstanceStatus {
	return []InstanceStatus{
		s.getCRDStatus(ctx, "devspaces.linuxsuren.github.io"),
		s.getCRDStatus(ctx, "users.linuxsuren.github.io"),
		s.getDeploymentStatus(ctx, namespace, "manager"),
		s.getDeploymentStatus(ctx, namespace, "apiserver-deploy"),
		s.getServiceStatus(ctx, namespace, "apiserver-service"),
		s.getConfigmapStatus(ctx, namespace, "config"),
		s.getIngressStatus(ctx, namespace, "apiserver"),
	}
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
