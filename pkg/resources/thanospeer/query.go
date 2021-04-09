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

package thanospeer

import (
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/resources"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r Reconciler) query(peerCert, peerCA string) resources.Resource {
	return func() (runtime.Object, reconciler.DesiredState, error) {
		meta := r.peer.Spec.MetaOverrides.Merge(r.getDescendantMeta())

		serverName := r.peer.Name

		if r.peer.Spec.PeerEndpointAlias != "" {
			serverName = r.peer.Spec.PeerEndpointAlias
		}

		cert := r.peer.Spec.Certificate
		if cert == "" {
			cert = peerCert
		}
		ca := r.peer.Spec.CABundle
		if ca == "" {
			ca = peerCA
		}

		query := &v1alpha1.Thanos{
			ObjectMeta: meta,
			Spec: v1alpha1.ThanosSpec{
				Query: &v1alpha1.Query{
					GRPCClientCertificate: cert,
					GRPCClientCA:          ca,
					GRPCClientServerName:  serverName,
					Stores: []string{
						r.peer.Spec.EndpointAddress,
					},
				},
			},
		}

		if r.peer.Spec.QueryOverrides != nil {
			if err := merge.Merge(query.Spec.Query, r.peer.Spec.QueryOverrides); err != nil {
				return nil, nil, errors.WrapIf(err, "failed to merge overrides to base query resource")
			}
		}

		return query, reconciler.StatePresent, nil
	}
}
