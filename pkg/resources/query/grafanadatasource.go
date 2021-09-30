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

package query

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/operator-tools/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type strIfMap = map[string]interface{}

func (q *Query) grafanaDatasource() (runtime.Object, reconciler.DesiredState, error) {

	grafanaDatasource := unstructured.Unstructured{}
	grafanaDatasource.SetAPIVersion("integreatly.org/v1alpha1")
	grafanaDatasource.SetKind("GrafanaDataSource")

	grafanaDatasource.SetName(q.getName())
	grafanaDatasource.SetNamespace(q.Thanos.Namespace)

	state := reconciler.StatePresent
	if q.Thanos.Spec.Query == nil || !q.Thanos.Spec.Query.GrafanaDatasource {
		return &grafanaDatasource, reconciler.StateAbsent, nil
	}

	grafanaDatasource.SetOwnerReferences([]metav1.OwnerReference{
		{
			APIVersion: q.Thanos.APIVersion,
			Kind:       q.Thanos.Kind,
			Name:       q.Thanos.Name,
			UID:        q.Thanos.UID,
			Controller: utils.BoolPointer(true),
		},
	})

	grafanaDatasource.SetLabels(q.GetCommonLabels())

	grafanaDatasource.Object["spec"] = strIfMap{
		"datasources": []strIfMap{
			{
				"access":    "proxy",
				"editable":  true,
				"isDefault": false,
				"name":      q.getName(),
				"type":      "prometheus",
				"url":       q.GetHTTPServiceURL(),
				"version":   1,
			},
		},
		"name": q.getName(),
	}

	return &grafanaDatasource, state, nil
}
