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

//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list

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

	gitpod := &v1alpha1.DevSpace{}
	if err = r.Get(r.ctx, types.NamespacedName{
		Name:      repoName,
		Namespace: pod.Namespace,
	}, gitpod); err != nil {
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
	gitpod.Status.DeployStatus = ""
	gitpod.Status.Pods = make([]v1.LocalObjectReference, podCount)
	for i, p := range pods.Items {
		gitpod.Status.Pods[i] = v1.LocalObjectReference{Name: p.Name}
		if p.Status.Phase == v1.PodRunning {
			gitpod.Status.DeployStatus = string(v1.PodRunning)
		}
	}
	if gitpod.Status.DeployStatus == "" && podCount > 0 {
		gitpod.Status.DeployStatus = string(pods.Items[0].Status.Phase)
	}
	err = r.Status().Update(r.ctx, gitpod)
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *DevSpacePodPodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return setupWithManagerAndLabels(mgr, &v1.Pod{}, map[string]string{LabelAppKind: "gitpod"}, r)
}
