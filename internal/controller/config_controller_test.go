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

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	"github.com/linuxsuren/kde/pkg/core"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
)

func TestConfigController(t *testing.T) {
	schema, err := v1alpha1.SchemeBuilder.Register().Build()
	assert.Nil(t, err)
	err = v1alpha1.SchemeBuilder.AddToScheme(schema)
	assert.Nil(t, err)
	err = v1.SchemeBuilder.AddToScheme(schema)
	assert.Nil(t, err)
	err = appsv1.SchemeBuilder.AddToScheme(schema)
	assert.NoError(t, err, err)
	err = networkingv1.SchemeBuilder.AddToScheme(schema)
	assert.NoError(t, err, err)
	defaultRequest := ctrl.Request{NamespacedName: types.NamespacedName{Name: "demo", Namespace: "default"}}

	type fields struct {
		Client   client.Client
		recorder record.EventRecorder
	}
	tests := []struct {
		name   string
		fields fields
		req    ctrl.Request
		verify func(*testing.T, ctrl.Result, client.Client, error)
	}{{
		name: "not found",
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).Build(),
		},
		req: defaultRequest,
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.NoError(t, err, err)
		},
	}, {
		name: "no required item",
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "demo",
					Namespace: "default",
				},
			}).Build(),
		},
		req: defaultRequest,
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.NoError(t, err, err)
		},
	}, {
		name: "invalid JSON item",
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "demo",
					Namespace: "default",
				},
				Data: map[string]string{
					core.ConfigFileName: "",
				},
			}).Build(),
		},
		req: defaultRequest,
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.Error(t, err, err)
		},
	}, {
		name: "host value is empty",
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "demo",
					Namespace: "default",
				},
				Data: map[string]string{
					core.ConfigFileName: `{"host":""}`,
				},
			}).Build(),
		},
		req: defaultRequest,
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.NoError(t, err, err)
		},
	}, {
		name: "ingress is not found",
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "demo",
					Namespace: "default",
				},
				Data: map[string]string{
					core.ConfigFileName: `{"host":"example.com"}`,
				},
			}).Build(),
		},
		req: defaultRequest,
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.NoError(t, err, err)
		},
	}, {
		name: "normal",
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "demo",
					Namespace: "default",
				},
				Data: map[string]string{
					core.ConfigFileName: `{"host":"good.com"}`,
				},
			}, &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kde-apiserver",
					Namespace: "default",
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{
						{Host: "example.com"},
					},
				},
			}).Build(),
		},
		req: defaultRequest,
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.NoError(t, err, err)

			ingress := &networkingv1.Ingress{}
			Client.Get(context.Background(), types.NamespacedName{Name: "kde-apiserver", Namespace: "default"}, ingress)
			assert.Equal(t, "good.com", ingress.Spec.Rules[0].Host)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := &FakeManager{
				Client: tt.fields.Client,
				Scheme: schema,
			}
			r := NewConfigReconciler(mgr)
			r.SetupWithManager(mgr)
			got, err := r.Reconcile(context.Background(), tt.req)
			tt.verify(t, got, tt.fields.Client, err)
		})
	}
}
