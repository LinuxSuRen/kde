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
package core_test

import (
	"context"
	"testing"

	"github.com/linuxsuren/kde/pkg/core"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestReadConfigFromConfigMap(t *testing.T) {
	t.Run("empty configmap", func(t *testing.T) {
		config, err := core.ReadConfigFromConfigMap(&corev1.ConfigMap{})
		assert.NoError(t, err)
		assert.Nil(t, config)
	})

	t.Run("invalid json format", func(t *testing.T) {
		config, err := core.ReadConfigFromConfigMap(&corev1.ConfigMap{
			Data: map[string]string{
				core.ConfigFileName: "invalid json",
			},
		})
		assert.Error(t, err)
		assert.NotNil(t, config)
	})
}

func TestGetConfigFromConfigMap(t *testing.T) {
	t.Run("not found configmap", func(t *testing.T) {
		client := fake.NewSimpleClientset()
		config, err := core.GetConfigFromConfigMap(context.Background(), client.CoreV1().ConfigMaps("default"), "test")
		assert.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("normal", func(t *testing.T) {
		client := fake.NewSimpleClientset(
			&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Data: map[string]string{
					core.ConfigFileName: "{}",
				},
			},
		)
		config, err := core.GetConfigFromConfigMap(context.Background(), client.CoreV1().ConfigMaps("default"), "test")
		assert.NoError(t, err)
		assert.NotNil(t, config)
	})
}
