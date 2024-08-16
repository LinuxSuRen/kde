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
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Server) IDEWebhook(c *gin.Context) {
	ns := c.Query("namespace")
	name := c.Query("devspace")
	token := c.Query("token")

	var devspace *v1alpha1.DevSpace
	var err error
	var httpStatus int
	defer func() {
		if httpStatus == 0 {
			httpStatus = http.StatusBadRequest
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusOK, "")
		}
	}()

	if ns == "" || name == "" {
		err = fmt.Errorf("the query parameter 'namespace' or 'devspace' is missing")
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	devspace, err = s.KClient.LinuxsurenV1alpha1().DevSpaces(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return
	}

	if devspace.Annotations == nil {
		// avoid nil pointer
		devspace.Annotations = make(map[string]string)
	}
	if token != devspace.Annotations[v1alpha1.AnnoKeyWebhookToken] {
		err = fmt.Errorf("token is incorrect")
		httpStatus = http.StatusForbidden
		return
	}

	payload := &webhookPayload{}
	if err = c.BindJSON(payload); err != nil {
		err = fmt.Errorf("failed to read payload: %w", err)
		return
	}

	if len(payload.Ports) > 0 {
		// convert ports into string
		ports := make([]string, len(payload.Ports))
		for i, port := range payload.Ports {
			ports[i] = fmt.Sprintf("%d", port)
		}
		devspace.Annotations[v1alpha1.AnnoKeyExposePorts] = strings.Join(ports, ",")
		if _, err = s.KClient.LinuxsurenV1alpha1().DevSpaces(ns).Update(ctx, devspace, metav1.UpdateOptions{}); err == nil {
			// get the status of it
			if devspace, err = s.KClient.LinuxsurenV1alpha1().DevSpaces(ns).Get(ctx, name, metav1.GetOptions{}); err == nil {
				c.JSON(http.StatusOK, devspace.Status.ExposeLinks)
			}
		}
	}
}

type webhookPayload struct {
	Ports []int `json:"ports"`
}
