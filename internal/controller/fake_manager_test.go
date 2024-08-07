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
	"context"
	"testing"

	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/meta"
	clientfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func TestFakeManager(t *testing.T) {
	schema, err := v1alpha1.SchemeBuilder.Register().Build()
	assert.Nil(t, err)
	client := clientfake.NewClientBuilder().WithScheme(schema).Build()

	fake := FakeManager{
		Client: client,
		Scheme: schema,
	}
	assert.Nil(t, fake.Add(nil))
	assert.Nil(t, fake.Elected())
	assert.Nil(t, fake.SetFields(nil))
	assert.Nil(t, fake.AddMetricsExtraHandler("", nil))
	assert.Nil(t, fake.AddHealthzCheck("", nil))
	assert.Nil(t, fake.AddReadyzCheck("", nil))
	assert.Nil(t, fake.Start(context.Background()))
	assert.Nil(t, fake.GetConfig())
	assert.Equal(t, schema, fake.GetScheme())
	assert.Equal(t, client, fake.GetClient())
	assert.Nil(t, fake.GetFieldIndexer())
	assert.Nil(t, fake.GetCache())
	assert.Nil(t, fake.GetEventRecorderFor(""))
	assert.Equal(t, meta.FirstHitRESTMapper{}, fake.GetRESTMapper())
	assert.Equal(t, client, fake.GetAPIReader())
	assert.Nil(t, fake.GetWebhookServer())
	assert.Equal(t, logr.New(log.NullLogSink{}), fake.GetLogger())
	assert.NotNil(t, fake.GetControllerOptions())
}
