/*
Copyright 2021.

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

package controllers

import (
	"context"

	"github.com/flanksource/commons/logger"
	"github.com/flanksource/kommons"
	konfigmanagerv1 "github.com/flanksource/konfig-manager/api/v1"
	"github.com/flanksource/konfig-manager/pkg"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HierarchyConfigReconciler reconciles a Konfig object
type HierarchyConfigReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Kommons *kommons.Client
	Log     logr.Logger
}

//+kubebuilder:rbac:groups=konfigmanager.flanksource.com,resources=hierarchyconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=konfigmanager.flanksource.com,resources=hierarchyconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=konfigmanager.flanksource.com,resources=hierarchyconfigs/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources="configmaps,secrets",verbs="*",namespace="*"

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Konfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *HierarchyConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	hierarchyConfig := &konfigmanagerv1.Konfig{}
	err := r.Get(ctx, req.NamespacedName, hierarchyConfig)
	if err != nil {
		r.Log.Error(err, "error fetching hierarchy config object name: %v namespace", req.Name, req.Namespace)
		return ctrl.Result{}, nil
	}
	config := pkg.Config{Hierarchy: hierarchyConfig.Spec.Hierarchy}
	resources, err := r.getResources(config)
	if err != nil {
		r.Log.Error(err, "error fetching resources")
		return ctrl.Result{}, nil
	}
	err = r.createOutputObject(hierarchyConfig.Spec.Output, config, resources)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HierarchyConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Kommons = kommons.NewClient(mgr.GetConfig(), logger.StandardLogger())
	return ctrl.NewControllerManagedBy(mgr).
		For(&konfigmanagerv1.Konfig{}).
		Complete(r)
}
