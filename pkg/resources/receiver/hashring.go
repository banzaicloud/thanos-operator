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
	"encoding/json"
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type HashRingGroup struct {
	HashRing  string   `json:"hashring,omitempty"`
	Endpoints []string `json:"endpoints,omitempty"`
	Tenants   []string `json:"tenants,omitempty"`
}

//func (r *Receiver) GetEndpointsForTenant(tenantName string) []string{
//	var endpoints []string
//	for _, receiverGroup := range r.ReceiverGroups {
//		if contains(receiverGroup.Tenants, tenantName) {
//			endpoints = append(endpoints, receiverGroup.Tenants...)
//		}
//	}
//	return endpoints
//}
//
//// Check if a string contains an element
//func contains(list []string, element string) bool {
//	for _, e := range list {
//		if e == element {
//			return true
//		}
//	}
//	return false
//}

func (r *Receiver) generateEndpointsForGroup(group v1alpha1.ReceiverGroup) []string {
	var endpoints []string
	replicas := int(group.Replicas)
	if replicas == 0 {
		replicas = 1
	}
	for i := 0; i < replicas; i++ {
		name := (&receiverInstance{
			Receiver:      r,
			receiverGroup: &group,
		}).getName(group.Name)
		endpoints = append(endpoints, fmt.Sprintf("%s-%d.%s:10907", name, i, name))
	}
	return endpoints
}

func (r *Receiver) generateHashring() (string, error) {
	hashringConfig := make([]HashRingGroup, len(r.Spec.ReceiverGroups))
	for i, receiverGroup := range r.Spec.ReceiverGroups {
		hashringConfig[i].HashRing = receiverGroup.Name
		hashringConfig[i].Tenants = receiverGroup.Tenants
		hashringConfig[i].Endpoints = r.generateEndpointsForGroup(receiverGroup)
	}
	result, err := json.Marshal(hashringConfig)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (r *Receiver) hashring() (runtime.Object, reconciler.DesiredState, error) {
	configuration, err := r.generateHashring()
	if err != nil {
		return nil, nil, err
	}
	configmap := &v1.ConfigMap{
		ObjectMeta: r.GetObjectMeta("hashring-config"),
		Data: map[string]string{
			"hashring.json": configuration,
		},
	}
	return configmap, reconciler.StatePresent, nil
}
