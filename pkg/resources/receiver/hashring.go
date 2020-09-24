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
	for i := 0; i < int(group.Replicas); i++ {
		name := (&receiverInstance{
			Receiver:      r,
			receiverGroup: &group,
		}).getName("")
		endpoints = append(endpoints, fmt.Sprintf("%s-%d.%s.%s.svc.%s", name, i, name, group.Namespace, r.GetClusterDomain()))
	}
	return endpoints
}

func (r *Receiver) GenerateHashring() (string, error) {
	hashringConfig := make([]HashRingGroup, len(r.Spec.ReceiverGroups))
	for i, receiverGroup := range r.Spec.ReceiverGroups {
		hashringConfig[i].HashRing = receiverGroup.GroupName
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
	configuration, err := r.GenerateHashring()
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
