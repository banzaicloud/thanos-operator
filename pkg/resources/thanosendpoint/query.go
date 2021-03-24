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
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/typeoverride"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r Reconciler) query() (runtime.Object, reconciler.DesiredState, error) {

	meta := r.endpoint.Spec.MetaOverrides.Merge(r.getDescendantMeta())

	query := &v1alpha1.Thanos{
		ObjectMeta: meta,
		Spec: v1alpha1.ThanosSpec{
			Query: &v1alpha1.Query{
				GRPCIngress: &v1alpha1.Ingress{
					IngressOverrides: &typeoverride.IngressNetworkingV1beta1{
						ObjectMeta: typeoverride.ObjectMeta{
							Annotations: map[string]string{
								"nginx.ingress.kubernetes.io/backend-protocol":       "GRPC",
								"nginx.ingress.kubernetes.io/auth-tls-verify-client": "on",
								"nginx.ingress.kubernetes.io/auth-tls-secret":        meta.Namespace + "/" + r.endpoint.Spec.CABundle,
							},
						},
					},
					Certificate: r.endpoint.Spec.Certificate,
					Host:        r.endpoint.Name,
				},
			},
		},
	}

	if r.endpoint.Spec.QueryOverrides != nil {
		if err := merge.Merge(query.Spec.Query, r.endpoint.Spec.QueryOverrides); err != nil {
			return nil, nil, errors.WrapIf(err, "failed to merge overrides to base query resource")
		}
	}

	return query, reconciler.StatePresent, nil
}
