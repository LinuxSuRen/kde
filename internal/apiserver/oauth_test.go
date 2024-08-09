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
	"testing"

	"github.com/linuxsuren/kde/internal/apiserver"
	"github.com/linuxsuren/kde/pkg/http"
	"github.com/stretchr/testify/assert"
)

func TestRegisterOAuth(t *testing.T) {
	t.Run("provider is empty", func(t *testing.T) {
		assert.NoError(t, apiserver.RegisterOAuth(nil, "", "", ""))
	})

	t.Run("not support oauth provider", func(t *testing.T) {
		assert.Error(t, apiserver.RegisterOAuth(nil, "not-support", "", ""))
	})

	t.Run("clientID or clientSecret is empty", func(t *testing.T) {
		assert.Error(t, apiserver.RegisterOAuth(nil, "github", "", "b"))
		assert.Error(t, apiserver.RegisterOAuth(nil, "github", "a", ""))
	})

	t.Run("normal", func(t *testing.T) {
		ginEngine := http.NewFakeGinEngine()

		assert.NoError(t, apiserver.RegisterOAuth(ginEngine, "github", "a", "b"))
	})
}

func TestOAuthHandler(t *testing.T) {
	assert.Nil(t, apiserver.OAuthHandler(""))
}
