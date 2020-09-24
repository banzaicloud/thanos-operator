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
	"github.com/banzaicloud/operator-tools/pkg/secret"
	"github.com/banzaicloud/operator-tools/pkg/types"
	"github.com/banzaicloud/operator-tools/pkg/volume"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DefaultReceiverGroup = &ReceiverGroup{
	Metrics: &Metrics{
		Interval:       "15s",
		Timeout:        "5s",
		Path:           "/metrics",
		ServiceMonitor: false,
	},
	HTTPAddress:        "0.0.0.0:10909",
	RemoteWriteAddress: "0.0.0.0:10908",
	GRPCAddress:        "0.0.0.0:10907",
	TSDBPath:           "/data",
}

type ReceiverSpec struct {
	ReceiverGroups []ReceiverGroup `json:"receiverGroups,omitempty"`
	ClusterDomain  string          `json:"clusterDomain,omitempty"`
}

// Defines a Receiver group
// Tenants are the Hard tenants of the receiver group
// Replicas are the number of instances in this receiver group
type ReceiverGroup struct {
	GroupName             string               `json:"groupName"`
	Namespace             string               `json:"namespace"`
	Tenants               []string             `json:"tenants,omitempty"`
	Config                secret.Secret        `json:"config"`
	Replicas              int32                `json:"replicas,omitempty"`
	MetaOverrides         *types.MetaBase      `json:"metaOverrides,omitempty"`
	WorkloadMetaOverrides *types.MetaBase      `json:"workloadMetaOverrides,omitempty"`
	WorkloadOverrides     *types.PodSpecBase   `json:"workloadOverrides,omitempty"`
	ContainerOverrides    *types.ContainerBase `json:"containerOverrides,omitempty"`
	HTTPIngress           *Ingress             `json:"HTTPIngress,omitempty"`
	// Secret name for HTTP Server certificate (Kubernetes TLS secret type)
	HTTPServerCertificate string `json:"HTTPServerCertificate,omitempty"`
	// Secret name for HTTP Client certificate (Kubernetes TLS secret type)
	HTTPClientCertificate string   `json:"HTTPClientCertificate,omitempty"`
	GRPCIngress           *Ingress `json:"GRPCIngress,omitempty"`
	// Secret name for GRPC Server certificate (Kubernetes TLS secret type)
	GRPCClientCertificate string `json:"GRPCClientCertificate,omitempty"`
	// Secret name for GRPC Client certificate (Kubernetes TLS secret type)
	GRPCServerCertificate string `json:"GRPCServerCertificate,omitempty"`
	// Server name to verify the hostname on the returned gRPC certificates. See https://tools.ietf.org/html/rfc4366#section-3.1
	RemoteWriteClientServerName string   `json:"remoteWriteClientServerName,omitempty" thanos:"--remote-write.client-server-name=%s"`
	Metrics                     *Metrics `json:"metrics,omitempty"`
	// Listen host:port for HTTP endpoints.
	HTTPAddress string `json:"httpAddress,omitempty" thanos:"--http-address=%s"`
	// Time to wait after an interrupt received for HTTP Server.
	HTTPGracePeriod metav1.Duration `json:"httpGracePeriod,omitempty"`
	// Listen ip:port address for gRPC endpoints
	GRPCAddress string `json:"grpcAddress,omitempty" thanos:"--grpc-address=%s"`
	// Time to wait after an interrupt received for GRPC Server.
	GRPCGracePeriod string `json:"grpcGracePeriod,omitempty" thanos:"--grpc-grace-period=%s"`
	// Address to listen on for remote write requests.
	RemoteWriteAddress string `json:"remoteWriteAddress,omitempty" thanos:"--remote-write.address=%s"`
	// External labels to announce. This flag will be removed in the future when handling multiple tsdb instances is added.
	Labels map[string]string `json:"labels,omitempty"`
	// Kubernetes volume abstraction refers to different types of volumes to be mounted to pods: emptyDir, hostPath, pvc.
	DataVolume *volume.KubernetesVolume `json:"dataVolume,omitempty"`
	TSDBPath   string                   `json:"tsdbPath,omitempty" thanos:"--tsdb.path=%s"`
	// How long to retain raw samples on local storage. 0d - disables this retention.
	TSDBRetention string `json:"tsdbRetention,omitempty" thanos:"--tsdb.retention=%s"`
	// Refresh interval to re-read the hashring configuration file. (used as a fallback)
	ReceiveHashringsFileRefreshInterval string `json:"receiveHashringsFileRefreshInterval,omitempty" thanos:"--receive.hashrings-file-refresh-interval=%s"`
	// HTTP header to determine tenant for write requests.
	ReceiveTenantHeader string `json:"receiveTenantHeader,omitempty" thanos:"--receive.tenant-header=%s"`
	// Default tenant ID to use when none is provided via a header.
	ReceiveDefaultTenantID string `json:"receiveDefaultTenantId,omitempty" thanos:"--receive.default-tenant-id=%s"`
	// Label name through which the tenant will be announced.
	ReceiveTenantLabelName string `json:"receiveTenantLabelName,omitempty" thanos:"--receive.tenant-label-name=%s"`
	// HTTP header specifying the replica number of a write request.
	ReceiveReplicaHeader string `json:"receiveReplicaHeader,omitempty" thanos:"--receive.replica-header=%s"`
	// How many times to replicate incoming write requests.
	ReceiveReplicationFactor int `json:"receiveReplicationFactor,omitempty" thanos:"--receive.replication-factor=%d"`
	// Compress the tsdb WAL.
	TSDBWalCompression *bool `json:"tsdbWalCompression,omitempty" thanos:"--tsdb.wal-compression"`
	// Do not create lockfile in TSDB data directory. In any case, the lockfiles will be deleted on next startup.
	TSDBNoLockfile *bool `json:"tsdbNoLockfile,omitempty" thanos:"--tsdb.no-lockfile"`
}

// ObjectStoreStatus defines the observed state of ObjectStore
type ReceiverStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Receiver is the Schema for the receiver cluster
type Receiver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReceiverSpec   `json:"spec,omitempty"`
	Status ReceiverStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ObjectStoreList contains a list of ObjectStore
type ReceiverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Receiver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Receiver{}, &ReceiverList{})
}
