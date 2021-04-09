// Copyright 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	prometheus "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type ControllerWithSource struct {
	Controller controller.Controller
	Source     runtime.Object
}

type ServiceMonitorWatchReconciler struct {
	Log         logr.Logger
	Controllers map[string]ControllerWithSource
	Client      client.Client
	added       bool
}

// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch

func (r *ServiceMonitorWatchReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	crd := &v1.CustomResourceDefinition{}
	err := r.Client.Get(ctx, req.NamespacedName, crd)
	if err != nil {
		// Object not found, return.  Created objects are automatically garbage collected.
		// For additional cleanup logic use finalizers.
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	var combinedErr error
	if crd.Spec.Names.Kind == "ServiceMonitor" && !r.added {
		r.Log.Info("adding watch for ServiceMonitor to all the interested controllers")
		for cName, c := range r.Controllers {
			err := c.Controller.Watch(&source.Kind{Type: &prometheus.ServiceMonitor{}}, &handler.EnqueueRequestForOwner{
				OwnerType:    c.Source,
				IsController: true,
			})
			if err != nil {
				r.Log.Error(err, "unable to add ServiceMonitor watch to controller", "controller", cName)
				combinedErr = errors.Combine(combinedErr, err)
			}
		}
		if combinedErr == nil {
			r.added = true
		}
	}

	return ctrl.Result{}, combinedErr
}

func (r *ServiceMonitorWatchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&v1.CustomResourceDefinition{}).Complete(r)
}
