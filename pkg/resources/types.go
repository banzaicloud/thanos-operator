// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
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
	StoreEndpoint  = "monitoring.banzaicloud.io/storeendpoint"

	HealthCheckPath = "/-/healthy"
	ReadyCheckPath  = "/-/ready"
)

// ComponentReconciler reconciler interface
type ComponentReconciler func() (*reconcile.Result, error)

// Resource redeclaration of function with return type kubernetes Object
type Resource func() (runtime.Object, reconciler.DesiredState, error)

type Labels map[string]string

func (l Labels) Merge(labelGroups ...Labels) Labels {
	for _, labels := range labelGroups {
		for k, v := range labels {
			l[k] = v
		}
	}
	return l
}
