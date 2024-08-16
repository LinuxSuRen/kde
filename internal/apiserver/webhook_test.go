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

package apiserver_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/linuxsuren/kde/internal/apiserver"
	"github.com/linuxsuren/kde/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
)

func TestIDEWebhook(t *testing.T) {
	t.Run("ns or name is empty", func(t *testing.T) {
		engine := gin.New()
		server := apiserver.Server{}
		engine.POST("/webhook", server.IDEWebhook)
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/webhook", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("not found devspace", func(t *testing.T) {
		engine := gin.New()
		server := apiserver.Server{
			KClient: fake.NewSimpleClientset(),
		}
		engine.POST("/webhook", server.IDEWebhook)
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/webhook?namespace=ns&devspace=fake", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("token is missing", func(t *testing.T) {
		engine := gin.New()
		server := apiserver.Server{
			KClient: fake.NewSimpleClientset(createDefaultDevSpace().DeepCopy()),
		}
		engine.POST("/webhook", server.IDEWebhook)
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/webhook?namespace=default&devspace=fake", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
	})

	t.Run("request payload is empty", func(t *testing.T) {
		engine := gin.New()
		server := apiserver.Server{
			KClient: fake.NewSimpleClientset(createDefaultDevSpace().DeepCopy()),
		}
		engine.POST("/webhook", server.IDEWebhook)
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/webhook?namespace=default&devspace=fake&token=token", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode, w.Body)
	})

	t.Run("normal", func(t *testing.T) {
		withoutAnnos := createDefaultDevSpace().DeepCopy()
		withoutAnnos.Annotations = nil

		engine := gin.New()
		server := apiserver.Server{
			KClient: fake.NewSimpleClientset(withoutAnnos.DeepCopy()),
		}
		engine.POST("/webhook", server.IDEWebhook)
		w := httptest.NewRecorder()

		payload := bytes.NewBufferString(`{"ports":[8080]}`)
		req, _ := http.NewRequest(http.MethodPost, "/webhook?namespace=default&devspace=fake", payload)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode, w.Body)
	})
}

func createDefaultDevSpace() *v1alpha1.DevSpace {
	return &v1alpha1.DevSpace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake",
			Namespace: "default",
			Annotations: map[string]string{
				v1alpha1.AnnoKeyWebhookToken: "token",
			},
		},
		Status: v1alpha1.DevSpaceStatus{
			ExposeLinks: []v1alpha1.ExposeLink{{
				Link: "example.com",
				Port: 8080,
			}},
		},
	}
}
