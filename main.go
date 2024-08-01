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

package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/kde/internal/apiserver"
	kdeClient "github.com/linuxsuren/kde/pkg/client/clientset/versioned"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", os.ExpandEnv("$HOME/.kube/config"))
		if err != nil {
			panic(err.Error())
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var dyClient dynamic.Interface
	if dyClient, err = dynamic.NewForConfig(config); err != nil {
		panic(err.Error())
	}

	var kClient *kdeClient.Clientset
	if kClient, err = kdeClient.NewForConfig(config); err != nil {
		panic(err.Error())
	}

	server := &apiserver.Server{
		Client:  clientset,
		KClient: kClient,
		DClient: dyClient,
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/devspace", server.ListDevSpace)
	r.POST("/devspace", server.CreateDevSpace)
	r.DELETE("/devspace/:devspace", server.DeleteDevSpace)
	r.PUT("/devspace/:devspace", server.UpdateDevSpace)
	r.GET("/devspace/:devspace", server.GetDevSpace)
	r.GET("/install", server.Install)
	r.Run()
}
