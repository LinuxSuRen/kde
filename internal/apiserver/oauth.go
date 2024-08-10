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

	"github.com/gin-gonic/gin"
	ginhttp "github.com/linuxsuren/kde/pkg/http"
	"github.com/linuxsuren/oauth-hub"
	"golang.org/x/oauth2"
)

const ContextKeyUser = "user"

func RegisterOAuth(r ginhttp.GinEngine, providerName, clientID, clientSecret string) (err error) {
	if providerName == "" {
		return
	}

	var provider oauth.OAuthProvider
	if provider = oauth.GetOAuthProvider(providerName); provider == nil {
		err = fmt.Errorf("not support: %s", providerName)
		return
	}

	if clientID == "" || clientSecret == "" {
		err = fmt.Errorf("clientID or clientSecret is empty")
		return
	}

	authHandler := oauth.NewAuth(provider, oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}, true)
	r.GET("/oauth2/login", func(c *gin.Context) {
		authHandler.RequestCode(c.Writer, c.Request, nil)
	})
	r.GET("/oauth2/getLocalCode", func(c *gin.Context) {
		authHandler.RequestLocalCode(c.Writer, c.Request, nil)
	})
	r.GET("/oauth2/getUserInfoFromLocalCode", func(c *gin.Context) {
		authHandler.RequestLocalToken(c.Writer, c.Request, nil)
	})
	r.GET("/oauth2/callback", func(c *gin.Context) {
		authHandler.Callback(c.Writer, c.Request, nil)
	})
	return
}

func OAuthHandler(provider string) func(*gin.Context) {
	// auth is disabled
	if provider == "" {
		return nil
	}

	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		user := oauth.GetUser(token)
		if user == nil {
			c.Header("WWW-Authenticate", "Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(ContextKeyUser, user)
	}
}
