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

package resources

import (
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	NameLabel      = "app.kubernetes.io/name"
	InstanceLabel  = "app.kubernetes.io/instance"
	VersionLabel   = "app.kubernetes.io/version"
	ComponentLabel = "app.kubernetes.io/component"
	ManagedByLabel = "app.kubernetes.io/managed-by"
	ManagedByValue = "thanos-operator"
	StoreEndpoint  = "monitoring.banzaicloud.io/storeendpoint"

	HealthCheckPath = "/-/healthy"
	ReadyCheckPath  = "/-/ready"
)

// ComponentReconciler reconciler interface
type ComponentReconciler func() (*reconcile.Result, error)

// Resource redeclaration of function with return type kubernetes Object
type Resource func() (runtime.Object, reconciler.DesiredState, error)

type Labels = map[string]string

func Dispatch(rr reconciler.ResourceReconciler, resourceList []Resource) (*reconcile.Result, error) {
	for _, res := range resourceList {
		o, state, err := res()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create desired object")
		}
		if o == nil {
			return nil, errors.Errorf("Reconcile error! Resource %#v returns with nil object", res)
		}
		result, err := rr.ReconcileResource(o, state)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to reconcile resource")
		}
		if result != nil {
			return result, nil
		}
	}
	return nil, nil
}
