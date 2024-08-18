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
	"fmt"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
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

	nodeMetrics, err := s.MetricClient.NodeMetricses().List(ctx, metav1.ListOptions{})
	cluster.Message = fmt.Sprintf("%v", err)
	if nodeMetrics == nil {
		nodeMetrics = &metricv1beta1.NodeMetricsList{}
	}

	for _, node := range nodeList.Items {
		clusterNode := ClusterNode{
			Name:             node.Name,
			Arch:             node.Status.NodeInfo.Architecture,
			ContaienrRuntime: node.Status.NodeInfo.ContainerRuntimeVersion,
			OS:               node.Status.NodeInfo.OperatingSystem,
			OSImage:          node.Status.NodeInfo.OSImage,
			Allocatable:      resourceListToNodeResource(node.Status.Allocatable),
			Capacity:         resourceListToNodeResource(node.Status.Capacity),
			Images:           node.Status.Images,
		}

		for _, nodeMetric := range nodeMetrics.Items {
			if nodeMetric.Name == node.Name {
				clusterNode.Usage = resourceListToNodeResource(nodeMetric.Usage)
				break
			}
		}

		cluster.Nodes = append(cluster.Nodes, clusterNode)
	}
	c.JSON(200, cluster)
}

func resourceListToNodeResource(res corev1.ResourceList) NodeResource {
	return NodeResource{
		CPU:              res.Cpu().String(),
		Memory:           res.Memory().String(),
		StorageEphemeral: res.StorageEphemeral().String(),
		Pods:             res.Pods().Value(),
	}
}

type Cluster struct {
	Nodes   []ClusterNode `json:"nodes"`
	Message string        `json:"message"`
}

type ClusterNode struct {
	Name             string                  `json:"name"`
	Arch             string                  `json:"arch"`
	ContaienrRuntime string                  `json:"containerRuntime"`
	OS               string                  `json:"os"`
	OSImage          string                  `json:"osImage"`
	Allocatable      NodeResource            `json:"allocatable"`
	Capacity         NodeResource            `json:"capacity"`
	Usage            NodeResource            `json:"usage"`
	Images           []corev1.ContainerImage `json:"images,omitempty"`
}

type NodeResource struct {
	CPU              string `json:"cpu"`
	Memory           string `json:"memory"`
	StorageEphemeral string `json:"storageEphemeral"`
	Pods             int64  `json:"pods"`
}
