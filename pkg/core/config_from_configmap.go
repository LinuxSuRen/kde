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

package core

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func GetConfigFromConfigMap(ctx context.Context, client typedcorev1.ConfigMapInterface, name string) (config *Config, err error) {
	var cm *corev1.ConfigMap
	if cm, err = client.Get(ctx, name, metav1.GetOptions{}); err == nil {
		config, err = ReadConfigFromConfigMap(cm)
	}
	return
}

// ReadConfigFromConfigMap reads the config data from a ConfigMap.
// Return nil if the ConfigMap is empty.
func ReadConfigFromConfigMap(cm *corev1.ConfigMap) (config *Config, err error) {
	data := []byte(cm.Data[ConfigFileName])
	if len(data) > 0 {
		return ParseConfigAsJSON(data)
	}
	return
}
