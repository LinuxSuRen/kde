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
	"os"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/kde/internal/apiserver"
	kdeClient "github.com/linuxsuren/kde/pkg/client/clientset/versioned"
	kdeui "github.com/linuxsuren/kde/ui/kde-ui"
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
	flags := cmd.Flags()
	flags.StringVar(&opt.address, "address", ":8080", "The address to listen")
	flags.StringVar(&opt.kubeConfig, "kube-config", os.ExpandEnv("$HOME/.kube/config"), "The kube config file")
	flags.StringVar(&opt.providerName, "oauth-provider-name", "", "The OAuth provider name")
	flags.StringVar(&opt.clientID, "oauth-client-id", "", "The OAuth client ID")
	flags.StringVar(&opt.clientSecret, "oauth-client-secret", "", "The OAuth client secret")
	flags.StringVar(&opt.systemNamespace, "system-namespace", "kde-system", "The system namespace")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type option struct {
	address                              string
	kubeConfig                           string
	providerName, clientID, clientSecret string
	systemNamespace                      string
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
		Client:          clientset,
		KClient:         kClient,
		DClient:         dyClient,
		ExtClient:       extClient,
		SystemNamespace: o.systemNamespace,
	}

	r := gin.Default()
	apiserver.RegisterStaticFilesHandle(r.Use(func(ctx *gin.Context) {
		ctx.Set("reader", kdeui.NewembedReader())
	}))
	if err = apiserver.RegisterOAuth(r, o.providerName, o.clientID, o.clientSecret); err != nil {
		return
	}
	apiserver.RegisterHealthEndpoint(r)
	r.POST("/webhook", server.IDEWebhook)

	authorizedAPI := r.Group("/api", apiserver.OAuthHandler(o.providerName))
	authorizedAPI.GET("/devspace", server.ListDevSpace)
	authorizedAPI.POST("/devspace", server.CreateDevSpace)
	authorizedAPI.DELETE("/devspace/:devspace", server.DeleteDevSpace)
	authorizedAPI.PUT("/devspace/:devspace", server.UpdateDevSpace)
	authorizedAPI.PUT("/devspace/:devspace/restart", server.RestartDevSpace)
	authorizedAPI.PUT("/devspace/:devspace/replicas", server.SetDevSpaceReplicas)
	authorizedAPI.GET("/devspace/:devspace", server.GetDevSpace)
	authorizedAPI.GET("/languages", server.GetDevSpaceLanguages)
	authorizedAPI.GET("/serverImages", server.ServerImages)
	authorizedAPI.POST("/install", server.Install)
	authorizedAPI.DELETE("/uninstall", server.Uninstall)
	authorizedAPI.GET("/instanceStatus", server.InstanceStatus)
	authorizedAPI.GET("/ws/instanceStatus", server.InstanceStatusWS)
	authorizedAPI.GET("/namespaces", server.Namespaces)
	authorizedAPI.GET("/images", server.Images)
	authorizedAPI.GET("/config", server.GetConfig)
	authorizedAPI.PUT("/config", server.UpdateConfig)
	authorizedAPI.GET("/cluster/info", server.ClusterInfo)
	r.Run(o.address)
}
