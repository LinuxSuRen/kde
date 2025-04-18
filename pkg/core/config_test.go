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
	"testing"

	_ "embed"

	"github.com/linuxsuren/kde/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestParseConfigAsJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    core.Config
		wantErr bool
	}{
		{
			name:    "valid config",
			data:    sampleConfigData,
			want:    sampleConfig,
			wantErr: false,
		},
		{
			name:    "invalid config",
			data:    []byte("invalid"),
			want:    core.Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := core.ParseConfigAsJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfigAsJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, *got, tt.want)
		})
	}
}

func TestReadConfigFromJSONFile(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    core.Config
		wantErr bool
	}{
		{
			name: "valid config",
			file: "testdata/config.json",
			want: sampleConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := core.ReadConfigFromJSONFile(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfigFromJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, *got, tt.want)

			_, err = got.ToJSON()
			if err != nil {
				t.Errorf("Failed to convert to JSON: %v", err)
			}
		})
	}
}

func TestCleanInvalidLanguages(t *testing.T) {
	languages := core.CleanInvalidLanguages([]core.Language{{
		Name: "go",
	}, {
		Name:  "go",
		Image: "golang:latest",
	}, {
		Image: "golang:latest",
	}})

	assert.Len(t, languages, 1)
}

var sampleConfig = core.Config{
	IngressMode:      "nginx",
	StorageClassName: "local-path",
	VolumeAccessMode: "ReadWriteOnce",
	VolumeMode:       "Filesystem",
}

//go:embed testdata/config.json
var sampleConfigData []byte
