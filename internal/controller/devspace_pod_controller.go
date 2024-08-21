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

	"github.com/go-logr/logr"
	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DevSpacePodPodReconciler struct {
	ctx context.Context
	client.Client
	log logr.Logger
}

//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// below rbac should in the apiserver
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;delete;create;update;watch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;delete;create;update
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;delete;create;update
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch
// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;delete;create;update
// +kubebuilder:rbac:groups="apiextensions.k8s.io",resources=namespaces,verbs=get;list;delete;create;update
// +kubebuilder:rbac:groups="apiextensions.k8s.io",resources=customresourcedefinitions,verbs=get;list;delete;create;update
// +kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;list;delete;create;update;watch
// +kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=clusterroles,verbs=get;list;delete;create;update
// +kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=clusterrolebindings,verbs=get;list;delete;create;update
// below rbac required when retrieving the resource lock for leader election
// +kubebuilder:rbac:groups="coordination.k8s.io",resources=leases,verbs=get;create;update
// +kubebuilder:rbac:groups="",resources=events,verbs=create
// below rbac required when retrieving the metrics
// +kubebuilder:rbac:groups="metrics.k8s.io",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="metrics.k8s.io",resources=nodes,verbs=get;list;watch

func (r *DevSpacePodPodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	r.ctx = ctx
	r.log = log.FromContext(ctx)
	r.log.Info("start to reconcile gitpod pod")

	pod := &v1.Pod{}
	if err = r.Get(r.ctx, req.NamespacedName, pod); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	// find its repository owner
	repoName := pod.Labels[LabelApp]
	if repoName == "" {
		return
	}

	devspace := &v1alpha1.DevSpace{}
	if err = r.Get(r.ctx, types.NamespacedName{
		Name:      repoName,
		Namespace: pod.Namespace,
	}, devspace); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	// find all related pods
	pods := &v1.PodList{}
	if err = r.List(r.ctx, pods, client.MatchingLabels{LabelApp: repoName}); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	podCount := len(pods.Items)
	devspace.Status.DeployStatus = ""
	devspace.Status.Pods = make([]v1.LocalObjectReference, podCount)
	for i, p := range pods.Items {
		devspace.Status.Pods[i] = v1.LocalObjectReference{Name: p.Name}
		if p.Status.Phase == v1.PodRunning {
			devspace.Status.DeployStatus = string(v1.PodRunning)
		}
	}
	if devspace.Status.DeployStatus == "" && podCount > 0 {
		devspace.Status.DeployStatus = string(pods.Items[0].Status.Phase)
	}
	err = r.Status().Update(r.ctx, devspace)
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *DevSpacePodPodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return setupWithManagerAndLabels(mgr, &v1.Pod{}, map[string]string{LabelAppKind: "devspace"}, r)
}
