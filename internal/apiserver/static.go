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
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	ginhttp "github.com/linuxsuren/kde/pkg/http"

	kdeui "github.com/linuxsuren/kde/ui/kde-ui"
)

func handleStaticFilesRequest(c *gin.Context) {
	staticFilePath := c.Request.URL.Path
	if staticFilePath == "/" {
		staticFilePath = "index.html"
	}
	data, err := kdeui.GetFile(filepath.Join("dist", staticFilePath))
	if err == nil {
		switch {
		case strings.HasSuffix(staticFilePath, ".js"):
			c.Writer.Header().Set("Content-Type", "application/javascript")
		case strings.HasSuffix(staticFilePath, ".css"):
			c.Writer.Header().Set("Content-Type", "text/css")
		case strings.HasSuffix(staticFilePath, ".png"):
			c.Writer.Header().Set("Content-Type", "image/png")
		case strings.HasSuffix(staticFilePath, ".svg"):
			c.Writer.Header().Set("Content-Type", "image/svg+xml")
		case strings.HasSuffix(staticFilePath, ".html"):
			c.Writer.Header().Set("Content-Type", "text/html")
		}
		// set the content length
		c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		c.Writer.Write(data)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
}

func RegisterStaticFilesHandle(r ginhttp.GinEngine) {
	r.GET("/", handleStaticFilesRequest)
	r.GET("/index.html", handleStaticFilesRequest)
	r.GET("/favicon.ico", handleStaticFilesRequest)
	r.GET("/assets/:asset", handleStaticFilesRequest)
}
