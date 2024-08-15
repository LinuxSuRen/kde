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
	fakehttp "github.com/linuxsuren/kde/pkg/http"
	"github.com/stretchr/testify/assert"
)

func TestRegisterStaticFileHandler(t *testing.T) {
	t.Run("only register", func(t *testing.T) {
		ginEngine := fakehttp.NewFakeGinEngine()
		apiserver.RegisterStaticFilesHandle(ginEngine)
	})

	t.Run("real request", func(t *testing.T) {
		apis := []string{
			"/", "/index.html", "/favicon.ico", "/assets/a.img",
		}
		for _, api := range apis {
			engine := gin.New()
			apiserver.RegisterStaticFilesHandle(engine)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, api, nil)
			engine.ServeHTTP(w, req)
			assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode, api)
		}
	})
}
