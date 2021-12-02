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

package receiver

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	monitoringv1alpha1 "github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"github.com/imdario/mergo"
	prometheus "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func NewComponent(scheme *runtime.Scheme) *Component {
	objects := []runtime.Object{
		&appsv1.StatefulSet{},
		&corev1.Service{},
		&prometheus.ServiceMonitor{},
		&netv1.Ingress{},
	}

	purgeTypes := make([]schema.GroupVersionKind, 0, len(objects))
	for _, obj := range objects {
		gvk, err := apiutil.GVKForObject(obj, scheme)
		if err != nil {
			panic(err)
		}
		purgeTypes = append(purgeTypes, gvk)
	}

	return &Component{
		purgeTypes: purgeTypes,
	}
}

type Component struct {
	purgeTypes []schema.GroupVersionKind
}

func (Component) ResourceBuilders(parent reconciler.ResourceOwner, config interface{}) (builders []reconciler.ResourceBuilder) {
	receiver, _ := parent.(*monitoringv1alpha1.Receiver)
	if receiver == nil {
		return
	}

	r := extend(receiver)
	builders = append(builders,
		r.commonService,
		r.hashringConfig,
	)

	for _, group := range receiver.Spec.ReceiverGroups {
		group := group
		if err := mergo.Merge(&group, monitoringv1alpha1.DefaultReceiverGroup); err != nil {
			return
		}
		g := receiverInstance{r, &group}
		builders = append(builders,
			g.statefulset,
			g.service,
			g.serviceMonitor,
			g.ingressGRPC,
			g.ingressHTTP,
		)
	}

	return
}

func (Component) RegisterWatches(b *builder.Builder) {
	b.For(&monitoringv1alpha1.Receiver{}).
		Owns(&corev1.Service{}).
		Owns(&netv1.Ingress{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.ConfigMap{})
}

func (c Component) PurgeTypes() []schema.GroupVersionKind {
	return append([]schema.GroupVersionKind(nil), c.purgeTypes...)
}
