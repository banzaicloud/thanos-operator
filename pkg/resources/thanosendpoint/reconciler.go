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
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/go-logr/logr"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	log                logr.Logger
	resourceReconciler reconciler.ResourceReconciler
	client             client.Client
	endpoint           *v1alpha1.ThanosEndpoint
}

func NewReconciler(logger logr.Logger, client client.Client, reconciler reconciler.ResourceReconciler, endpoint *v1alpha1.ThanosEndpoint) *Reconciler {
	return &Reconciler{
		log:                logger,
		resourceReconciler: reconciler,
		endpoint:           endpoint,
		client:             client,
	}
}

func (r Reconciler) Reconcile() (*reconcile.Result, error) {
	var resourceList []resources.Resource

	endpointDetails := []interface{}{
		"name", r.endpoint.Name,
		"namespace", r.endpoint.Namespace,
	}

	resourceList = append(resourceList, r.query)
	resourceList = append(resourceList, r.storeEndpoint)
	result, err := resources.Dispatch(r.resourceReconciler, resourceList)
	if err != nil {
		return result, err
	}

	ctx := context.Background()
	ingList := &netv1.IngressList{}
	err = r.client.List(ctx, ingList, client.MatchingLabels{
		"app.kubernetes.io/name":       v1alpha1.QueryName,
		"app.kubernetes.io/managed-by": r.getDescendantResourceName(),
	})
	if err != nil {
		return result, errors.WrapWithDetails(err, "failed to list ingresses for ThanosEndpoint", endpointDetails...)
	}

	originalEndpointAddress := r.endpoint.Status.EndpointAddress

	switch len(ingList.Items) {
	case 0:
		return result, nil
	case 1:
		ing := ingList.Items[0]
		switch len(ing.Status.LoadBalancer.Ingress) {
		case 0:
			return result, nil
		case 1:
			addr := originalEndpointAddress
			if ing.Status.LoadBalancer.Ingress[0].Hostname != "" {
				addr = fmt.Sprintf("%s:443", ing.Status.LoadBalancer.Ingress[0].Hostname)
			}
			if ing.Status.LoadBalancer.Ingress[0].IP != "" {
				addr = fmt.Sprintf("%s:443", ing.Status.LoadBalancer.Ingress[0].IP)
			}
			if originalEndpointAddress != addr {
				r.log.Info("updating status", "endpointAddress", addr)
				r.endpoint.Status.EndpointAddress = addr
				if err := r.client.Status().Update(ctx, r.endpoint); err != nil {
					return result, errors.WrapIfWithDetails(err, "failed status update for ThanosEndpoint", endpointDetails...)
				}
			}
		default:
			return result, errors.Errorf("multiple items detected for Ingress %s/%s, should be only one", ing.Namespace, ing.Name)
		}
		return result, nil
	default:
		return result, errors.NewWithDetails("multiple ingress resources found for ThanosEndpoint, should be only one", endpointDetails...)
	}
}

func (r Reconciler) getDescendantMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      r.getDescendantResourceName(),
		Namespace: r.endpoint.Namespace,
	}
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

func (r Reconciler) getDescendantResourceName() string {
	name := r.qualifiedName(v1alpha1.EndpointName)
	return name
}

func (r Reconciler) qualifiedName(name string) string {
	return fmt.Sprintf("%s-%s", r.endpoint.Name, name)
}

func (r Reconciler) getLabels() resources.Labels {
	return resources.Labels{
		resources.NameLabel:      v1alpha1.EndpointName,
		resources.InstanceLabel:  r.endpoint.Name,
		resources.ManagedByLabel: resources.ManagedByValue,
	}
}
