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
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func createOrUpdate(ctx context.Context, cli client.Client, obj *unstructured.Unstructured) (err error) {
	copiedObj := obj.DeepCopy()
	objKey := client.ObjectKeyFromObject(copiedObj)
	if err = cli.Get(ctx, objKey, copiedObj); err != nil {
		if obj.GetName() == "" {
			// skip those resources only have the generatedName
			err = nil
		}

		if err = client.IgnoreNotFound(err); err == nil {
			if err = cli.Create(ctx, obj); err != nil {
				err = fmt.Errorf("failed to create %v, error: %v", objKey, err)
			}
		} else {
			err = fmt.Errorf("failed to get %v, error: %v", objKey, err)
		}
	} else {
		toUpdate := obj.DeepCopy()
		toUpdate.SetResourceVersion(copiedObj.GetResourceVersion())

		if err = cli.Update(ctx, toUpdate); err != nil {
			err = fmt.Errorf("failed to update %v, error: %v", objKey, err)
		}
	}
	return
}

func createOrUpdateObjs(ctx context.Context, cli client.Client, objs ...*unstructured.Unstructured) (err error) {
	for _, obj := range objs {
		if obj == nil {
			continue
		}
		err = errors.Join(err, createOrUpdate(ctx, cli, obj))
	}
	return
}

func setupWithManagerAndLabels(mgr ctrl.Manager, object client.Object, labels map[string]string, r reconcile.Reconciler) error {
	labelPredicate, err := predicate.LabelSelectorPredicate(metav1.LabelSelector{
		MatchLabels: labels,
	})
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(object).WithEventFilter(labelPredicate).
		Complete(r)
}

// StringToIntSlice converts a string to a slice of int.
// It will remove the duplicated numbers.
// for example: "1,2,3" -> [1,2,3]
func StringToIntSlice(array string) []int {
	if array == "" {
		return nil
	}
	unique := make(map[int]bool)
	for _, s := range strings.Split(array, ",") {
		if s == "" {
			continue
		}
		if i, err := strconv.Atoi(s); err == nil {
			if _, ok := unique[i]; !ok {
				unique[i] = true
			}
		}
	}

	if len(unique) > 0 {
		result := make([]int, 0, len(unique))
		for k := range unique {
			result = append(result, k)
		}
		sort.Ints(result)
		return result
	}

	return nil
}
