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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/kde/internal/apiserver"
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
}
