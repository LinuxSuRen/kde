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
	"strings"
)

type Config struct {
	StorageClassName string     `json:"storageClassName"`
	VolumeMode       string     `json:"volumeMode"`
	VolumeAccessMode string     `json:"volumeAccessMode"`
	IngressMode      string     `json:"ingressMode"`
	ImagePullPolicy  string     `json:"imagePullPolicy"`
	Host             string     `json:"host"`
	Languages        []Language `json:"languages"`
}

type Language struct {
	Name  string `json:"name"`
	Image string `json:"image"`
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

const ConfigFileName = "config.json"

func CleanInvalidLanguages(languages []Language) []Language {
	for i := 0; i < len(languages); {
		languages[i].Image = strings.TrimSpace(languages[i].Image)
		languages[i].Name = strings.TrimSpace(languages[i].Name)
		if languages[i].Image == "" || languages[i].Name == "" {
			if i == len(languages)-1 {
				languages = languages[:i]
			} else {
				languages = append(languages[:i], languages[i+1:]...)
			}
		} else {
			i++
		}
	}
	return languages
}

type FileReader interface {
	GetFile(name string) (data []byte, err error)
}

const FileReaderContext = "reader"
