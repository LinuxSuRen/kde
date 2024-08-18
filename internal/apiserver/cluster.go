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
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Server) ClusterInfo(c *gin.Context) {
	ctx := c.Request.Context()
	nodeList, err := s.Client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		c.Error(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	cluster := Cluster{
		Nodes: make([]ClusterNode, 0, len(nodeList.Items)),
	}
	for _, node := range nodeList.Items {
		cluster.Nodes = append(cluster.Nodes, ClusterNode{
			Name:             node.Name,
			Arch:             node.Status.NodeInfo.Architecture,
			ContaienrRuntime: node.Status.NodeInfo.ContainerRuntimeVersion,
			OS:               node.Status.NodeInfo.OperatingSystem,
			OSImage:          node.Status.NodeInfo.OSImage,
			Allocatable: NodeResource{
				CPU:              node.Status.Allocatable.Cpu().String(),
				Memory:           node.Status.Allocatable.Memory().String(),
				StorageEphemeral: node.Status.Allocatable.StorageEphemeral().String(),
				Pods:             int(node.Status.Allocatable.Pods().Value()),
			},
			Capacity: NodeResource{
				CPU:              node.Status.Capacity.Cpu().String(),
				Memory:           node.Status.Capacity.Memory().String(),
				StorageEphemeral: node.Status.Capacity.StorageEphemeral().String(),
				Pods:             int(node.Status.Capacity.Pods().Value()),
			},
		})
	}
	c.JSON(200, cluster)
}

type Cluster struct {
	Nodes []ClusterNode `json:"nodes"`
}

type ClusterNode struct {
	Name             string       `json:"name"`
	Arch             string       `json:"arch"`
	ContaienrRuntime string       `json:"containerRuntime"`
	OS               string       `json:"os"`
	OSImage          string       `json:"osImage"`
	Allocatable      NodeResource `json:"allocatable"`
	Capacity         NodeResource `json:"capacity"`
}

type NodeResource struct {
	CPU              string `json:"cpu"`
	Memory           string `json:"memory"`
	StorageEphemeral string `json:"storageEphemeral"`
	Pods             int    `json:"pods"`
}
