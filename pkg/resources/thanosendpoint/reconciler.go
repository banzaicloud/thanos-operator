// Copyright 2021 Banzai Cloud
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

package thanosendpoint

import (
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	log logr.Logger
	resourceReconciler reconciler.ResourceReconciler
	endpoint *v1alpha1.ThanosEndpoint
}

func NewReconciler(logger logr.Logger, reconciler reconciler.ResourceReconciler, endpoint *v1alpha1.ThanosEndpoint) *Reconciler {
	return &Reconciler{
		log: logger,
		resourceReconciler: reconciler,
		endpoint: endpoint,
	}
}

func (r Reconciler) Reconcile() (*reconcile.Result, error) {
	var resourceList []resources.Resource
	resourceList = append(resourceList, r.query)
	resourceList = append(resourceList, r.storeEndpoint)
	return resources.Dispatch(r.resourceReconciler, resourceList)
}


func (r Reconciler) getMeta(suffix ...string) metav1.ObjectMeta {
	nameSuffix := ""
	if len(suffix) > 0 {
		nameSuffix = suffix[0]
	}
	meta := r.getNameMeta(r.getName(nameSuffix), "")
	meta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: r.endpoint.APIVersion,
			Kind:       r.endpoint.Kind,
			Name:       r.endpoint.Name,
			UID:        r.endpoint.UID,
			Controller: utils.BoolPointer(true),
		},
	}
	meta.Labels = r.getLabels()
	return meta
}

func (r Reconciler) getNameMeta(name string, namespaceOverride string) metav1.ObjectMeta {
	namespace := r.endpoint.Namespace
	if namespaceOverride != "" {
		namespace = namespaceOverride
	}
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
}

func (r Reconciler) getName(suffix ...string) string {
	name := r.qualifiedName(v1alpha1.EndpointName)
	if len(suffix) > 0 && suffix[0] != "" {
		name = name + "-" + suffix[0]
	}
	return name
}

func (r Reconciler) qualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", r.endpoint.Name, name)
}

func (r Reconciler) getLabels() resources.Labels {
	return resources.Labels{
		resources.NameLabel:     v1alpha1.EndpointName,
		resources.InstanceLabel:  r.endpoint.Name,
		resources.ManagedByLabel: resources.ManagedByValue,
	}
}
