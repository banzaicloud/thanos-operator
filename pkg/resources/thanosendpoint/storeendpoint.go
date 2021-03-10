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
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r Reconciler) storeEndpoint() (runtime.Object, reconciler.DesiredState, error) {

	meta := r.endpoint.Spec.MetaOverrides.Merge(r.getMeta())

	storeEndpoint := &v1alpha1.StoreEndpoint{
		ObjectMeta: meta,
		Spec: v1alpha1.StoreEndpointSpec{
			Selector: &v1alpha1.KubernetesSelector{},
			Thanos:   meta.Name,
		},
	}

	if r.endpoint.Spec.StoreEndpointOverrides != nil {
		if err := merge.Merge(storeEndpoint.Spec, r.endpoint.Spec.StoreEndpointOverrides); err != nil {
			return nil, nil, errors.Wrap(err, "failed to merge storeendpoint overrides")
		}
	}

	return storeEndpoint, reconciler.StatePresent, nil
}
