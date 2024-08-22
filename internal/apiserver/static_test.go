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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/kde/internal/apiserver"
	"github.com/linuxsuren/kde/pkg/core"
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
			"/assets/index.js", "/assets/index.css", "/assets/demo.png",
			"/assets/demo.svg",
		}
		for _, api := range apis {
			engine := gin.New()
			engine.Use(func(ctx *gin.Context) {
				ctx.Set(core.FileReaderContext, &fakeReader{})
			})
			apiserver.RegisterStaticFilesHandle(engine)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, api, nil)
			engine.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Result().StatusCode, api)
		}
	})

	t.Run("no reader in context", func(t *testing.T) {
		engine := gin.New()
		apiserver.RegisterStaticFilesHandle(engine)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("invalid reader in context", func(t *testing.T) {
		engine := gin.New()
		engine.Use(func(ctx *gin.Context) {
			ctx.Set(core.FileReaderContext, "")
		})
		apiserver.RegisterStaticFilesHandle(engine)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("failed to read file", func(t *testing.T) {
		apis := []string{
			"/", "/index.html",
		}
		for _, api := range apis {
			engine := gin.New()
			engine.Use(func(ctx *gin.Context) {
				ctx.Set(core.FileReaderContext, &fakeReader{err: errors.New("failed to read")})
			})
			apiserver.RegisterStaticFilesHandle(engine)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, api, nil)
			engine.ServeHTTP(w, req)
			assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode, api)
		}
	})

	t.Run("rediect", func(t *testing.T) {
		apis := []string{
			"/dashboard", "/system", "/dev",
		}
		for _, api := range apis {
			engine := gin.New()
			engine.Use(func(ctx *gin.Context) {
				ctx.Set(core.FileReaderContext, &fakeReader{})
			})
			apiserver.RegisterStaticFilesHandle(engine)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, api, nil)
			engine.ServeHTTP(w, req)
			assert.Equal(t, http.StatusMovedPermanently, w.Result().StatusCode, api)
		}
	})
}

type fakeReader struct {
	err error
}

func (r *fakeReader) GetFile(name string) (data []byte, err error) {
	data = []byte("")
	err = r.err
	return
}
