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
	"encoding/json"
	"os"
)

type Config struct {
	StorageClassName string `json:"storageClassName"`
	VolumeMode       string `json:"volumeMode"`
	VolumeAccessMode string `json:"volumeAccessMode"`
	IngressMode      string `json:"ingressMode"`
	ImagePullPolicy  string `json:"imagePullPolicy"`
}

func ReadConfigFromJSONFile(file string) (config *Config, err error) {
	var data []byte
	data, err = os.ReadFile(file)
	if err == nil {
		config, err = ParseConfigAsJSON(data)
	}
	return
}

func ParseConfigAsJSON(data []byte) (config *Config, err error) {
	config = &Config{}
	err = json.Unmarshal(data, &config)
	return
}

func (c *Config) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}
