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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	networkv1 "k8s.io/api/networking/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
)

func TestGitPodTemplateRender(t *testing.T) {
	gitpod := createDefaultGitPod()

	t.Run("deploy", func(t *testing.T) {
		deploy := turnTemplateToUnstructured(gitpodDeployment, gitpod)

		data, err := deploy.MarshalJSON()
		assert.NoError(t, err, err)
		assert.Contains(t, string(data), "initContainers", string(data))

		deployment := &appsv1.Deployment{}
		err = json.Unmarshal(data, deployment)
		assert.NoError(t, err)
		assert.NotEmpty(t, deployment.Spec.Template.Spec.Containers[0].Image)
	})

	t.Run("pvc", func(t *testing.T) {
		pvc := turnTemplateToUnstructured(gitpodPvc, gitpod)

		data, err := pvc.MarshalJSON()
		assert.NoError(t, err, err)
		assert.Contains(t, string(data), `"storageClassName":"storageClassName"`, string(data))
	})

	t.Run("service", func(t *testing.T) {
		service := turnTemplateToUnstructured(gitpodService, gitpod)
		data, err := service.MarshalJSON()
		assert.NoError(t, err, err)
		assert.Contains(t, string(data), `"name":"8080"`, string(data))
		assert.Contains(t, string(data), `"name":"9090"`, string(data))
	})

	t.Run("ingress", func(t *testing.T) {
		ingress := turnTemplateToUnstructured(gitpodExposeIngress, gitpod)
		data, err := ingress.MarshalJSON()
		assert.NoError(t, err, err)
		assert.Contains(t, string(data), `"host":"8080.demo.gitpod.linuxsuren.github.io"`, string(data))
	})
}

func TestUpdateStatus(t *testing.T) {
	reconciler := &DevSpaceReconciler{
		Ingress: "gitpod.linuxsuren.github.io",
	}
	gitpod := createDefaultGitPod()
	gitpod = reconciler.updateStatus(gitpod)
	assert.Equal(t, "demo.gitpod.linuxsuren.github.io", gitpod.Status.Link)
	assert.Equal(t, []v1alpha1.ExposeLink{
		{
			Port: 8080,
			Link: "8080.demo.gitpod.linuxsuren.github.io",
		},
		{
			Port: 9090,
			Link: "9090.demo.gitpod.linuxsuren.github.io",
		},
	}, gitpod.Status.ExposeLinks)
}

func createDefaultGitPod() *v1alpha1.DevSpace {
	return &v1alpha1.DevSpace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo",
			Namespace: "default",
			Annotations: map[string]string{
				"storageClassName":           "storageClassName",
				"volumeMode":                 "volumeMode",
				v1alpha1.AnnoKeyExposePorts:  strings.Join([]string{"8080", "9090", "abc", "8080", "-1", "", "65536", "3000"}, ","),
				v1alpha1.AnnoKeyMaintainMode: PolicyAlways,
			},
		},
		Spec: v1alpha1.DevSpaceSpec{
			Image:  "linuxsuren.github.io/image:tag",
			Memory: "1Gi",
			CPU:    "1",
			Host:   "gitpod.linuxsuren.github.io",
			Windows: []v1alpha1.Window{{
				From: "01:00",
				To:   "15:00",
			}},
		},
		Status: v1alpha1.DevSpaceStatus{
			Link: "demo.gitpod.linuxsuren.github.io",
			ExposeLinks: []v1alpha1.ExposeLink{
				{
					Port: 8080,
					Link: "8080.demo.gitpod.linuxsuren.github.io",
				},
				{
					Port: 9090,
					Link: "9090.demo.gitpod.linuxsuren.github.io",
				},
			},
		},
	}
}

func TestGitPodController(t *testing.T) {
	schema, err := v1alpha1.SchemeBuilder.Register().Build()
	assert.Nil(t, err)
	err = v1alpha1.SchemeBuilder.AddToScheme(schema)
	assert.Nil(t, err)
	err = v1.SchemeBuilder.AddToScheme(schema)
	assert.Nil(t, err)
	err = appsv1.SchemeBuilder.AddToScheme(schema)
	assert.NoError(t, err, err)
	defaultRequest := ctrl.Request{NamespacedName: types.NamespacedName{Name: "demo", Namespace: "default"}}

	withoutDefaultValue := createDefaultGitPod().DeepCopy()
	withoutDefaultValue.Spec.Image = ""
	withoutDefaultValue.Spec.Host = ""
	withoutDefaultValue.Annotations = map[string]string{
		v1alpha1.AnnoKeyMaintainMode: PolicyAlways,
	}
	withoutDefaultValue.Spec.Auth = v1alpha1.DevSpaceAuth{
		BasicAuth: &v1alpha1.BasicAuth{
			Username: "test",
			Password: "test",
		},
	}

	noMaintainMode := createDefaultGitPod().DeepCopy()
	noMaintainMode.Annotations = nil

	zeroReplicas := createDefaultGitPod().DeepCopy()
	zeroReplicas.Spec.Replicas = 0
	zeroReplicas.Spec.Windows = nil

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
		name: "normal",
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(createDefaultGitPod()).Build(),
		},
		req: defaultRequest,
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.Equal(t, time.Minute, r.RequeueAfter)

			assert.NoError(t, err, err)
			ctx := context.TODO()

			deploy := &appsv1.Deployment{}
			err = Client.Get(ctx, types.NamespacedName{Name: "demo", Namespace: "default"}, deploy)
			assert.NoError(t, err, err)
			assert.Equal(t, v1.PullPolicy("IfNotPresent"), deploy.Spec.Template.Spec.Containers[0].ImagePullPolicy)

			service := &v1.Service{}
			err = Client.Get(ctx, types.NamespacedName{Name: "demo", Namespace: "default"}, service)
			assert.NoError(t, err, err)
			assert.ElementsMatch(t, []v1.ServicePort{{
				Name: "http", Port: 3000, TargetPort: intstr.IntOrString{Type: intstr.String, StrVal: "http"}, Protocol: v1.ProtocolTCP,
			}, {
				Name: "8080", Port: 8080, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080}, Protocol: v1.ProtocolTCP,
			}, {
				Name: "9090", Port: 9090, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9090}, Protocol: v1.ProtocolTCP,
			}}, service.Spec.Ports)

			configmap := &v1.ConfigMap{}
			err = Client.Get(ctx, types.NamespacedName{Name: "demo", Namespace: "default"}, configmap)
			assert.NoError(t, err, err)

			pvc := &v1.PersistentVolumeClaim{}
			err = Client.Get(ctx, types.NamespacedName{Name: "demo", Namespace: "default"}, pvc)
			assert.NoError(t, err, err)

			secret := &v1.Secret{}
			err = Client.Get(ctx, types.NamespacedName{Name: "demo", Namespace: "default"}, secret)
			assert.Error(t, err, err)

			ingress := &networkv1.Ingress{}
			err = Client.Get(ctx, types.NamespacedName{Name: "demo", Namespace: "default"}, ingress)
			assert.Error(t, err, err)
		},
	}, {
		name: "without default value",
		req:  defaultRequest,
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(withoutDefaultValue.DeepCopy()).Build(),
		},
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.NoError(t, err)

			gitpod := &v1alpha1.DevSpace{}
			assert.NoError(t, Client.Get(context.TODO(), defaultRequest.NamespacedName, gitpod))
			assert.Empty(t, gitpod.Spec.Auth.BasicAuth.Password)
			assert.NotEmpty(t, gitpod.Annotations[v1alpha1.AnnoKeyBasicAuth])

			var auth []byte
			auth, err = base64.StdEncoding.DecodeString(gitpod.Annotations[v1alpha1.AnnoKeyBasicAuth])
			assert.NoError(t, err)
			assert.Contains(t, string(auth), gitpod.Spec.Auth.BasicAuth.Username)

			deploy := &appsv1.Deployment{}
			err = Client.Get(context.Background(), defaultRequest.NamespacedName, deploy)
			assert.NoError(t, err)
			assert.NotEmpty(t, deploy.Spec.Template.Spec.Containers[0].Image)
			assert.Contains(t, deploy.Spec.Template.Spec.Containers[0].Image, "openvscode-server-full")
		},
	}, {
		name: "the replicas number is zero",
		req:  defaultRequest,
		fields: fields{
			Client: fake.NewClientBuilder().WithScheme(schema).WithObjects(zeroReplicas.DeepCopy()).Build(),
		},
		verify: func(t *testing.T, r ctrl.Result, Client client.Client, err error) {
			assert.NoError(t, err)

			gitpod := &v1alpha1.DevSpace{}
			assert.NoError(t, Client.Get(context.TODO(), defaultRequest.NamespacedName, gitpod))
			assert.Equal(t, 0, len(gitpod.Status.Pods))
			assert.Empty(t, gitpod.Status.DeployStatus)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DevSpaceReconciler{
				Client:   tt.fields.Client,
				Recorder: tt.fields.recorder,
				Ingress:  "gitpod.linuxsuren.github.io",
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

func TestIsInAliveWindows(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		windows      []v1alpha1.Window
		expectHasWin bool
		expectOk     bool
		expectErr    bool
	}{{
		name:         "no windows",
		windows:      nil,
		expectHasWin: false,
		expectOk:     true,
	}, {
		name: "invalid time format",
		windows: []v1alpha1.Window{{
			From: "aa:bb",
		}, {
			From: "00:12:00",
			To:   "cc:dd",
		}},
		expectHasWin: true,
		expectErr:    true,
		expectOk:     true,
	}, {
		name: "before the window",
		windows: []v1alpha1.Window{{
			From: now.Add(time.Minute).Format(time.TimeOnly),
			To:   now.Add(time.Minute * 2).Format(time.TimeOnly),
		}},
		expectHasWin: true,
		expectErr:    false,
		expectOk:     false,
	}, {
		name: "after the window",
		windows: []v1alpha1.Window{{
			From: now.Add(-2 * time.Minute).Format(time.TimeOnly),
			To:   now.Add(-1 * time.Minute).Format(time.TimeOnly),
		}},
		expectHasWin: true,
		expectErr:    false,
		expectOk:     false,
	}, {
		name: "from bigger than to",
		windows: []v1alpha1.Window{{
			From: now.Add(2 * time.Minute).Format(time.TimeOnly),
			To:   now.Add(-1 * time.Minute).Format(time.TimeOnly),
		}},
		expectHasWin: true,
		expectErr:    true,
		expectOk:     false,
	}, {
		name: "in the window",
		windows: []v1alpha1.Window{{
			From: now.Add(-1 * time.Minute).Format(time.TimeOnly),
			To:   now.Add(1 * time.Minute).Format(time.TimeOnly),
		}},
		expectHasWin: true,
		expectErr:    false,
		expectOk:     true,
	}}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasWin, _, err := isInAliveWindows(now, tt.windows)
			assert.Equal(t, tt.expectHasWin, hasWin, fmt.Sprintf("case %d", i))
			assert.Equal(t, tt.expectErr, err != nil, fmt.Sprintf("case %d, error: %v", i, err))
		})
	}
}
