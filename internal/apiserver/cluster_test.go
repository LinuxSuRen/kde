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
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestClusterInfo(t *testing.T) {
	s := &Server{
		Client: fake.NewSimpleClientset(&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-1",
			},
			Status: corev1.NodeStatus{
				NodeInfo: corev1.NodeSystemInfo{
					OSImage:                 "Ubuntu 20.04.4 LTS",
					Architecture:            "amd64",
					ContainerRuntimeVersion: "containerd://1.6.4",
					OperatingSystem:         "linux",
				},
			},
		}),
	}
	engine := gin.New()
	engine.GET("/cluster", s.ClusterInfo)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/cluster", nil)
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, clusterInfo, w.Body.String())
}

//go:embed testdata/clusterInfo.json
var clusterInfo string
