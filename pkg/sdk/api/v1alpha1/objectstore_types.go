// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"net"
	"strings"
	"time"

	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ObjectStoreSpec defines the desired state of ObjectStore
type ObjectStoreSpec struct {
	// Config
	Config    secret.Secret `json:"config"`
	Compactor *Compactor    `json:"compactor,omitempty"`
	BucketWeb *BucketWeb    `json:"bucketWeb,omitempty"`
}

type Compactor struct {
	BaseObject `json:",inline"`
	Enabled    bool `json:"enabled,omitempty"`
	// Listen host:port for HTTP endpoints.
	HTTPAddress string `json:"httpAddress,omitempty"`
	// Time to wait after an interrupt received for HTTP Server.
	HTTPGracePeriod metav1.Duration `json:"httpGracePeriod,omitempty"`
	// Data directory in which to cache blocks and process compactions.
	DataDir string `json:"dataDir,omitempty"`
	// Minimum age of fresh (non-compacted) blocks before they are being processed.
	// Malformed blocks older than the maximum of consistency-delay and 48h0m0s will be removed.
	ConsistencyDelay metav1.Duration `json:"consistencyDelay,omitempty"`
	// How long to retain raw samples in bucket. 0d - disables this retention.
	RetentionResolutionRaw metav1.Duration `json:"retentionResolutionRaw,omitempty"`
	// How long to retain samples of resolution 1 (5 minutes) in bucket. 0d - disables this retention.
	RetentionResolution5m metav1.Duration `json:"retentionResolution5m,omitempty"`
	// How long to retain samples of resolution 2 (1 hour) in bucket. 0d - disables this retention.
	RetentionResolution1h metav1.Duration `json:"retentionResolution1h,omitempty"`
	// Do not exit after all compactions have been processed and wait for new work.
	Wait bool `json:"wait,omitempty"`
	// Disables downsampling. This is not recommended as querying long time ranges without non-downsampleddata
	// is not efficient and useful e.g it is not possible to render all samples for a human eye anyway.
	DownsamplingDisable bool `json:"downsamplingDisable,omitempty"`
	// Number of goroutines to use when syncing block metadata from object storage.
	// +kubebuilder:validation:Minimum=1
	BlockSyncConcurrency int `json:"blockSyncConcurrency,omitempty"`
	// Number of goroutines to use when compacting groups.
	// +kubebuilder:validation:Minimum=1
	CompactConcurrency int `json:"compactConcurrency,omitempty"`
}

func (in *Compactor) SetDefaults() (*Compactor, error) {
	out := in.DeepCopy()

	// BaseObject
	if out.Image.Repository == "" {
		out.Image.Repository = "quay.io/thanos/thanos"
	}
	if out.Image.Tag == "" {
		out.Image.Tag = "v0.9.0"
	}
	switch strings.ToLower(out.Image.PullPolicy) {
	case strings.ToLower(string(corev1.PullAlways)):
		out.Image.PullPolicy = string(corev1.PullAlways)
	case strings.ToLower(string(corev1.PullNever)):
		out.Image.PullPolicy = string(corev1.PullNever)
	default:
		out.Image.PullPolicy = string(corev1.PullIfNotPresent)
	}

	// HTTPAddress
	if out.HTTPAddress == "" {
		out.HTTPAddress = "0.0.0.0:10902"
	}
	host, port, err := net.SplitHostPort(out.HTTPAddress)
	if err != nil {
		return nil, err
	}
	out.HTTPAddress = net.JoinHostPort(host, port)

	if out.HTTPGracePeriod.Duration < time.Nanosecond {
		out.HTTPGracePeriod.Duration = 2 * time.Minute
	}

	if out.DataDir == "" {
		out.DataDir = "./data"
	}

	if out.ConsistencyDelay.Duration < time.Nanosecond {
		out.ConsistencyDelay.Duration = 30 * time.Minute
	}

	if out.RetentionResolutionRaw.Duration < 0 {
		out.RetentionResolutionRaw.Duration = 0
	}

	if out.RetentionResolution5m.Duration < 0 {
		out.RetentionResolution5m.Duration = 0
	}

	if out.RetentionResolution1h.Duration < 0 {
		out.RetentionResolution1h.Duration = 0
	}

	if out.BlockSyncConcurrency <= 0 {
		out.BlockSyncConcurrency = 20
	}

	if out.CompactConcurrency < 1 {
		out.CompactConcurrency = 1
	}

	return out, nil
}

type BucketWeb struct {
	BaseObject `json:",inline"`
	Enabled    bool `json:"enabled,omitempty"`
	// Listen host:port for HTTP endpoints.
	HTTPAddress string `json:"httpAddress,omitempty"`
	// Time to wait after an interrupt received for HTTP Server.
	HTTPGracePeriod metav1.Duration `json:"httpGracePeriod,omitempty"`
	// Static prefix for all HTML links and redirect URLs in the bucket web UI interface. Actual endpoints are still served on / or the web.route-prefix. This allows thanos bucket web UI to be served behind a reverse proxy that strips a URL sub-path.
	WebExternalPrefix string `json:"web_external_prefix,omitempty"`
	// Name of HTTP request header used for dynamic prefixing of UI links and redirects. This option is ignored if web.external-prefix argument is set. Security risk: enable this option only if a reverse proxy in front of thanos is resetting the header. The --web.prefix-header=X-Forwarded-Prefix option can be useful, for example, if Thanos UI is served via Traefik reverse proxy with PathPrefixStrip option enabled, which sends the stripped prefix value in X-Forwarded-Prefix header. This allows thanos UI to be served on a sub-path.
	WebPrefixHeader string `json:"web_prefix_header,omitempty"`
	// Refresh interval to download metadata from remote storage.
	Refresh metav1.Duration `json:"refresh,omitempty"`
	// Timeout to download metadata from remote.
	Timeout metav1.Duration `json:"timeout,omitempty"`
	// Prometheus label to use as timeline title.
	Label string `json:"label,omitempty"`
}

func (in *BucketWeb) SetDefaults() (*BucketWeb, error) {
	out := in.DeepCopy()

	// BaseObject
	if out.Image.Repository == "" {
		out.Image.Repository = "quay.io/thanos/thanos"
	}
	if out.Image.Tag == "" {
		out.Image.Tag = "v0.9.0"
	}
	switch strings.ToLower(out.Image.PullPolicy) {
	case strings.ToLower(string(corev1.PullAlways)):
		out.Image.PullPolicy = string(corev1.PullAlways)
	case strings.ToLower(string(corev1.PullNever)):
		out.Image.PullPolicy = string(corev1.PullNever)
	default:
		out.Image.PullPolicy = string(corev1.PullIfNotPresent)
	}

	// HTTPAddress
	if out.HTTPAddress == "" {
		out.HTTPAddress = "0.0.0.0:10903"
	}
	host, port, err := net.SplitHostPort(out.HTTPAddress)
	if err != nil {
		return nil, err
	}
	out.HTTPAddress = net.JoinHostPort(host, port)

	if out.HTTPGracePeriod.Duration < time.Nanosecond {
		out.HTTPGracePeriod.Duration = 2 * time.Minute
	}

	if out.Refresh.Duration < time.Nanosecond {
		out.Refresh.Duration = 30 * time.Minute
	}

	if out.Timeout.Duration < time.Nanosecond {
		out.Timeout.Duration = 5 * time.Minute
	}

	return out, nil
}

// ObjectStoreStatus defines the observed state of ObjectStore
type ObjectStoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// ObjectStore is the Schema for the objectstores API
type ObjectStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ObjectStoreSpec   `json:"spec,omitempty"`
	Status ObjectStoreStatus `json:"status,omitempty"`
}

func (in *ObjectStore) SetDefaults() (*ObjectStore, error) {
	var err error
	out := in.DeepCopy()
	if out.Spec.Compactor == nil {
		out.Spec.Compactor = &Compactor{}
	}
	out.Spec.Compactor, err = out.Spec.Compactor.SetDefaults()
	if err != nil {
		return nil, err
	}
	if out.Spec.BucketWeb == nil {
		out.Spec.BucketWeb = &BucketWeb{}
	}
	out.Spec.BucketWeb, err = out.Spec.BucketWeb.SetDefaults()
	if err != nil {
		return nil, err
	}

	return out, nil
}

// +kubebuilder:object:root=true

// ObjectStoreList contains a list of ObjectStore
type ObjectStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ObjectStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ObjectStore{}, &ObjectStoreList{})
}
