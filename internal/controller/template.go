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

package controller

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func turnTemplateToUnstructured(tplText string, app interface{}) (result *unstructured.Unstructured) {
	if tpl, err := template.New("turnTemplateToUnstructured").Parse(tplText); err == nil {
		buf := new(bytes.Buffer)
		if err = tpl.Execute(buf, app); err == nil {
			if strings.TrimSpace(buf.String()) == "" {
				return
			}

			var data []byte
			result = &unstructured.Unstructured{}
			if data, err = yaml.ToJSON(buf.Bytes()); err == nil {
				_ = result.UnmarshalJSON(data)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return
}
