/*
Copyright 2024 kde authrors.

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
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGitPodPodController(t *testing.T) {
	schema, err := v1alpha1.SchemeBuilder.Register().Build()
	assert.Nil(t, err)
	err = v1.SchemeBuilder.AddToScheme(schema)
	assert.Nil(t, err)
	err = appsv1.SchemeBuilder.AddToScheme(schema)
	assert.NoError(t, err, err)
	defaultRequest := ctrl.Request{NamespacedName: types.NamespacedName{Name: "demo", Namespace: "default"}}

	defaultPod := createDefaultPod()
	withoutRepoName := defaultPod.DeepCopy()
	withoutRepoName.Labels = nil

	pendingPhasePod := defaultPod.DeepCopy()
	pendingPhasePod.Status.Phase = v1.PodPending

	type fields struct {
		Client client.Client
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
		name: "without repo name",
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(withoutRepoName.DeepCopy()).Build(),
		},
		req: defaultRequest,
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.NoError(t, err, err)
		},
	}, {
		name: "without gitpod found",
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(defaultPod.DeepCopy()).Build(),
		},
		req: defaultRequest,
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.NoError(t, err, err)
		},
	}, {
		name: "normal",
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(defaultPod.DeepCopy(), createDefaultGitPod().DeepCopy()).Build(),
		},
		req: defaultRequest,
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.NoError(t, err, err)

			gitpod := &v1alpha1.DevSpace{}
			err = Client.Get(context.TODO(), defaultRequest.NamespacedName, gitpod)
			assert.NoError(t, err, err)
			assert.Equal(t, string(v1.PodRunning), gitpod.Status.DeployStatus)
		},
	}, {
		name: "pending status",
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(pendingPhasePod.DeepCopy(), createDefaultGitPod().DeepCopy()).Build(),
		},
		req: defaultRequest,
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.NoError(t, err, err)

			gitpod := &v1alpha1.DevSpace{}
			err = Client.Get(context.TODO(), defaultRequest.NamespacedName, gitpod)
			assert.NoError(t, err, err)
			assert.Equal(t, string(v1.PodPending), gitpod.Status.DeployStatus)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DevSpacePodPodReconciler{
				Client: tt.fields.Client,
			}
			mgr := &FakeManager{
				Client: tt.fields.Client,
				Scheme: schema,
			}
			r.SetupWithManager(mgr)
			got, err := r.Reconcile(context.Background(), tt.req)
			tt.verify(t, got, tt.fields.Client, err)
		})
	}
}

func createDefaultPod() *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo",
			Namespace: "default",
			Labels: map[string]string{
				LabelApp: "demo",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: "test-container",
				},
			},
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
		},
	}
}
