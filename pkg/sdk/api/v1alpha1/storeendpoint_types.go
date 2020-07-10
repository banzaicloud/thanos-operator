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

package v1alpha1

import (
	"fmt"

	"github.com/banzaicloud/operator-tools/pkg/secret"
	"github.com/banzaicloud/operator-tools/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	QueryName   = "query"
	StoreName   = "store"
	RuleName    = "rule"
	SidecarName = "sidecar"
)

// StoreEndpointSpec defines the desired state of StoreEndpoint
type StoreEndpointSpec struct {
	// Foo is an example field of StoreEndpoint. Edit StoreEndpoint_types.go to remove/update
	MetaOverrides *types.MetaBase     `json:"metaOverrides,omitempty"`
	URL           string              `json:"url,omitempty"`
	Selector      *KubernetesSelector `json:"selector,omitempty"`
	Config        secret.Secret       `json:"config,omitempty"`
	Thanos        string              `json:"thanos"`
	Ingress       *Ingress            `json:"ingress,omitempty"`
}

type KubernetesSelector struct {
	Namespace   string            `json:"namespaces,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	HTTPPort    int32             `json:"httpPort,omitempty"`
	GRPCPort    int32             `json:"grpcPort,omitempty"`
}

// StoreEndpointStatus defines the observed state of StoreEndpoint
type StoreEndpointStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// StoreEndpoint is the Schema for the storeendpoints API
type StoreEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StoreEndpointSpec   `json:"spec,omitempty"`
	Status StoreEndpointStatus `json:"status,omitempty"`
}

func (s *StoreEndpoint) GetServiceURL() string {
	if s.Spec.URL != "" {
		return s.Spec.URL
	}
	if s.Spec.Selector != nil {
		return fmt.Sprintf("dnssrvnoa+_grpc._tcp.%s.%s.svc.cluster.local", fmt.Sprintf("%s-%s", s.Name, SidecarName), s.Namespace)
	}
	return ""
}

// +kubebuilder:object:root=true

// StoreEndpointList contains a list of StoreEndpoint
type StoreEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StoreEndpoint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StoreEndpoint{}, &StoreEndpointList{})
}
