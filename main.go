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

package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/kde/internal/apiserver"
	kdeClient "github.com/linuxsuren/kde/pkg/client/clientset/versioned"
	"github.com/linuxsuren/kde/pkg/core"
	"github.com/spf13/cobra"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	opt := &option{}
	cmd := cobra.Command{
		Use:   "kde",
		Short: "A Kubernetes DevSpace manager",
		Run:   opt.runE,
	}
	cmd.Flags().StringVar(&opt.address, "address", ":8080", "The address to listen")
	cmd.Flags().StringVar(&opt.config, "config", "", "The config file")
	cmd.Flags().StringVar(&opt.kubeConfig, "kube-config", os.ExpandEnv("$HOME/.kube/config"), "The kube config file")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type option struct {
	address    string
	config     string
	kubeConfig string
}

func (o *option) runE(cmd *cobra.Command, args []string) {
	// creates the in-cluster config
	config, err := clientcmd.BuildConfigFromFlags("", os.ExpandEnv(o.kubeConfig))
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	var dyClient dynamic.Interface
	if dyClient, err = dynamic.NewForConfig(config); err != nil {
		return
	}

	var kClient *kdeClient.Clientset
	if kClient, err = kdeClient.NewForConfig(config); err != nil {
		return
	}

	var extClient *apiextensionsclientset.Clientset
	if extClient, err = apiextensionsclientset.NewForConfig(config); err != nil {
		return
	}

	server := &apiserver.Server{
		Client:    clientset,
		KClient:   kClient,
		DClient:   dyClient,
		ExtClient: extClient,
	}

	if o.config != "" {
		if server.Config, err = core.ReadConfigFromJSONFile(o.config); err != nil {
			return
		}
	} else {
		server.Config = &core.Config{}
	}

	r := gin.Default()
	apiserver.RegisterStaticFilesHandle(r)
	r.GET("/api/devspace", server.ListDevSpace)
	r.POST("/api/devspace", server.CreateDevSpace)
	r.DELETE("/api/devspace/:devspace", server.DeleteDevSpace)
	r.PUT("/api/devspace/:devspace", server.UpdateDevSpace)
	r.GET("/api/devspace/:devspace", server.GetDevSpace)
	r.GET("/api/languages", server.GetDevSpaceLanguages)
	r.POST("/api/install", server.Install)
	r.DELETE("/api/uninstall", server.Uninstall)
	r.GET("/api/instanceStatus", server.InstanceStatus)
	r.GET("/api/ws/instanceStatus", server.InstanceStatusWS)
	r.GET("/api/namespaces", server.Namespaces)
	r.GET("/api/config", server.GetConfig)
    r.PUT("/api/config", server.UpdateConfig)
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.Run(o.address)
}
