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
	"errors"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	_ "embed"

	"github.com/go-logr/logr"
	"github.com/johnaoss/htpasswd/apr1"
	"github.com/linuxsuren/kde/api/linuxsuren.github.io/v1alpha1"
	"github.com/linuxsuren/kde/pkg/core"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
)

// DevSpaceReconciler reconciles a DevSpace object
type DevSpaceReconciler struct {
	client.Client
	Scheme          *runtime.Scheme
	Recorder        record.EventRecorder
	SystemNamespace string
	// inner fields
	ctx context.Context
	log logr.Logger
}

// +kubebuilder:rbac:groups=linuxsuren.github.io,resources=devspaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=linuxsuren.github.io,resources=devspaces/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=linuxsuren.github.io,resources=devspaces/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;delete;create;update;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;delete;create;update;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DevSpace object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *DevSpaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	r.ctx = ctx
	r.log = log.FromContext(ctx)
	r.log.Info("start to reconcile gitpod")

	devSpace := &v1alpha1.DevSpace{}
	if err = r.Get(ctx, req.NamespacedName, devSpace); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	config := &core.Config{}
	configCM := &corev1.ConfigMap{}
	if cfgErr := r.Get(ctx, types.NamespacedName{
		Name:      "config",
		Namespace: r.SystemNamespace,
	}, configCM); cfgErr == nil {
		if config, cfgErr = core.ReadConfigFromConfigMap(configCM); err != nil {
			r.log.Error(cfgErr, "failed to parse config")
		}
	} else {
		r.log.Error(cfgErr, "failed to get config")
	}

	setDefaultValueForGitPod(devSpace, config.Host)
	devSpace = r.updateStatus(devSpace)

	_ = r.Status().Update(ctx, devSpace.DeepCopy())

	if err = r.Get(ctx, req.NamespacedName, devSpace); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}
	setDefaultValueForGitPod(devSpace, config.Host)
	devSpace.Annotations[v1alpha1.AnnoKeyServiceNamespace] = r.SystemNamespace
	devSpace.Annotations[v1alpha1.AnnoKeyServiceName] = "apiserver"
	configmap, configmapErr := turnTemplateToUnstructured(gitpodConfigMap, devSpace)
	secret, secretErr := turnTemplateToUnstructured(gitpodSecret, devSpace)
	pvc, pvcErr := turnTemplateToUnstructured(gitpodPvc, devSpace)
	deploy, deployErr := turnTemplateToUnstructured(gitpodDeployment, devSpace)
	service, serviceErr := turnTemplateToUnstructured(gitpodService, devSpace)
	ingress, ingressErr := turnTemplateToUnstructured(gitpodIngress, devSpace)
	exposeIngress, exposeIngressErr := turnTemplateToUnstructured(gitpodExposeIngress, devSpace)

	// check the object templates render result
	if err = errors.Join(configmapErr, secretErr, pvcErr, deployErr, serviceErr, ingressErr, exposeIngressErr); err != nil {
		r.Recorder.Event(devSpace, v1.EventTypeWarning, "Render", err.Error())
		return
	}

	result.RequeueAfter = time.Minute

	auth := devSpace.Spec.Auth.BasicAuth
	if auth != nil {
		passwd := auth.Password
		if passwd != "" && auth.Username != "" {
			var hash string
			hash, err = apr1.Hash(passwd, "")
			if err != nil {
				return
			}

			base64Str := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", auth.Username, hash)))
			devSpace.Annotations[v1alpha1.AnnoKeyBasicAuth] = base64Str
			auth.Password = "" // keep the password safe
			if err = r.Update(ctx, devSpace); err != nil {
				return
			}
		}

		err = createOrUpdate(ctx, r.Client, secret)
	} else {
		err = r.Delete(ctx, secret)
		err = client.IgnoreNotFound(err)
	}

	err = errors.Join(err, createOrUpdateObjs(ctx, r.Client, configmap, pvc, deploy, service, ingress, exposeIngress))
	return
}

func setDefaultValueForGitPod(devspace *v1alpha1.DevSpace, ingress string) {
	if devspace.Spec.Image == "" {
		devspace.Spec.Image = "ghcr.io/linuxsuren/openvscode-server-full:v0.0.8"
	}
	if devspace.Spec.Host == "" {
		devspace.Spec.Host = ingress
	}

	if devspace.Annotations == nil {
		devspace.Annotations = map[string]string{}
	}
	if policy := devspace.Annotations[v1alpha1.AnnoKeyImagePullPolicy]; policy != PolicyAlways {
		devspace.Annotations[v1alpha1.AnnoKeyImagePullPolicy] = PolicyIfNotPresent
	}
}

func isInAliveWindows(specifyTime time.Time, windows []v1alpha1.Window) (hasWin, ok bool, err error) {
	nowTime, _ := time.ParseInLocation(time.TimeOnly, specifyTime.Format(time.TimeOnly), time.Local)
	ok = true
	hasWin = len(windows) > 0
	if hasWin {
		for _, win := range windows {
			from, fErr := time.ParseInLocation(time.TimeOnly, win.From, time.Local)
			to, tErr := time.ParseInLocation(time.TimeOnly, win.To, time.Local)
			if fErr != nil {
				err = errors.Join(err, fmt.Errorf("failed to parse window from time: %q, error: %v", win.From, fErr))
				continue
			}
			if tErr != nil {
				err = errors.Join(err, fmt.Errorf("failed to parse window to time: %q, error: %v", win.To, tErr))
				continue
			}
			if from.After(to) {
				err = errors.Join(err, fmt.Errorf("from (%q) should smaller than to (%q)", win.From, win.To))
				continue
			}

			if nowTime.After(from) && nowTime.Before(to) {
				return
			}
		}
	}
	ok = false
	return
}

func (r *DevSpaceReconciler) updateStatus(devSpace *v1alpha1.DevSpace) *v1alpha1.DevSpace {
	devSpace.Status.Link = fmt.Sprintf("%s.%s", devSpace.Name, devSpace.Spec.Host)
	devSpace.Status.ExposeLinks = nil
	devSpace.Status.Phase = v1alpha1.DevSpacePhaseReady
	ports := devSpace.Annotations[v1alpha1.AnnoKeyExposePorts]
	if ports != "" {
		portSlice := StringToIntSlice(ports)
		for _, port := range portSlice {
			// TODO this is the port of the gitpod, but we should ignore the port via configuration instead of hard code
			if port == 3000 || port == 2376 {
				continue
			}

			if port <= 0 || port > 65535 {
				r.log.Info("invalid port", "port", port, "key", client.ObjectKeyFromObject(devSpace))
				continue
			}

			devSpace.Status.ExposeLinks = append(devSpace.Status.ExposeLinks, v1alpha1.ExposeLink{
				Link: fmt.Sprintf("%d.%s", port, devSpace.Status.Link),
				Port: port,
			})
		}
	}

	// check the alive windows
	var (
		hasWin, ok bool
		wErr       error
	)
	if hasWin, ok, wErr = isInAliveWindows(time.Now(), devSpace.Spec.Windows); hasWin {
		if wErr != nil {
			r.log.Error(wErr, "got error when parsing windows time")
		}

		if !ok {
			devSpace.Status.Phase = v1alpha1.DevSpacePhaseOff
		}
	}
	r.log.Info("check alive window", "enable", hasWin, "in", ok, "now", time.Now().Format(time.TimeOnly))

	// check if all the pods are deleted
	if devSpace.Status.Phase != v1alpha1.DevSpacePhaseOff && devSpace.Spec.Replicas != nil && *(devSpace.Spec.Replicas) <= 0 {
		podList := &corev1.PodList{}
		listErr := r.List(r.ctx, podList, client.MatchingLabels{LabelApp: devSpace.Name})
		if listErr == nil {
			if len(podList.Items) == 0 {
				devSpace.Status.Pods = nil
				devSpace.Status.DeployStatus = ""
			}
		} else {
			r.log.Error(listErr, "failed to list pods", LabelApp, devSpace.Name)
		}
	}
	return devSpace
}

const (
	PolicyAlways       = "Always"
	PolicyIfNotPresent = "IfNotPresent"
)

//go:embed data/deployment.yaml
var gitpodDeployment string

//go:embed data/service.yaml
var gitpodService string

//go:embed data/pvc.yaml
var gitpodPvc string

//go:embed data/ingress.yaml
var gitpodIngress string

//go:embed data/ingress-expose.yaml
var gitpodExposeIngress string

//go:embed data/configmap.yaml
var gitpodConfigMap string

//go:embed data/secret.yaml
var gitpodSecret string

// SetupWithManager sets up the controller with the Manager.
func (r *DevSpaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.DevSpace{}).
		Complete(r)
}
