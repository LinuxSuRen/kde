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
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
	"github.com/linuxsuren/kde/config"
	"github.com/linuxsuren/kde/pkg/core"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type configReconciler struct {
	client.Client
	log logr.Logger
}

func NewConfigReconciler(mgr ctrl.Manager) *configReconciler {
	return &configReconciler{
		Client: mgr.GetClient(),
		log:    ctrl.Log.WithName("config-controller"),
	}
}

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list
func (r *configReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	cm := &v1.ConfigMap{}
	if err = r.Get(ctx, req.NamespacedName, cm); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	configStr, ok := cm.Data[core.ConfigFileName]
	if !ok {
		return
	}

	var configObj *core.Config
	if configObj, err = core.ParseConfigAsJSON([]byte(configStr)); err != nil {
		return
	}

	if strings.TrimSpace(configObj.Host) == "" {
		return
	}

	data, _ := config.GetFile(filepath.Join("manager", "ingress.yaml"))
	ingress := &networkingv1.Ingress{}
	if data, err = yaml.ToJSON(data); err == nil {
		yaml.Unmarshal(data, ingress)
	} else {
		return
	}

	if err = r.Get(ctx, client.ObjectKeyFromObject(ingress), ingress); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	if len(ingress.Spec.Rules) > 0 && ingress.Spec.Rules[0].Host != configObj.Host {
		ingress.Spec.Rules[0].Host = configObj.Host
		err = r.Update(ctx, ingress)
	}
	return
}

func (r *configReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return setupWithManagerAndLabels(mgr, &v1.ConfigMap{}, map[string]string{LabelAppKind: "devspace"}, r)
}
