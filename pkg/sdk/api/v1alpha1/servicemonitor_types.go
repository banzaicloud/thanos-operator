package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const (
	ServiceMonitorsKind   = "ServiceMonitor"
	ServiceMonitorName    = "servicemonitors"
	ServiceMonitorKindKey = "servicemonitor"
)

var (
	// GroupVersion is group version used to register these objects
	ServiceMonitorGroupVersion = schema.GroupVersion{Group: "monitoring.coreos.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	ServiceMonitorSchemeBuilder = &scheme.Builder{GroupVersion: ServiceMonitorGroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	ServiceMonitorAddToScheme = ServiceMonitorSchemeBuilder.AddToScheme
)

func init() {
	ServiceMonitorSchemeBuilder.Register(&ServiceMonitor{}, &ServiceMonitorList{})
}

// +kubebuilder:object:root=true
type ServiceMonitor struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	// +k8s:openapi-gen=false
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of desired Service selection for target discrovery by
	// Prometheus.
	Spec ServiceMonitorSpec `json:"spec"`
}

// +kubebuilder:object:root=true
type ServiceMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of ServiceMonitors
	Items []*ServiceMonitor `json:"items"`
}

type ServiceMonitorSpec struct {
	// The label to use to retrieve the job name from.
	JobLabel string `json:"jobLabel,omitempty"`
	// TargetLabels transfers labels on the Kubernetes Service onto the target.
	TargetLabels []string `json:"targetLabels,omitempty"`
	// PodTargetLabels transfers labels on the Kubernetes Pod onto the target.
	PodTargetLabels []string `json:"podTargetLabels,omitempty"`
	// A list of endpoints allowed as part of this ServiceMonitor.
	Endpoints []Endpoint `json:"endpoints"`
	// Selector to select Endpoints objects.
	Selector metav1.LabelSelector `json:"selector"`
	// Selector to select which namespaces the Endpoints objects are discovered from.
	NamespaceSelector NamespaceSelector `json:"namespaceSelector,omitempty"`
	// SampleLimit defines per-scrape limit on number of scraped samples that will be accepted.
	SampleLimit uint64 `json:"sampleLimit,omitempty"`
}

type Endpoint struct {
	// Name of the service port this endpoint refers to. Mutually exclusive with targetPort.
	Port string `json:"port,omitempty"`
	// Name or number of the target port of the endpoint. Mutually exclusive with port.
	TargetPort *intstr.IntOrString `json:"targetPort,omitempty"`
	// HTTP path to scrape for metrics.
	Path string `json:"path,omitempty"`
	// HTTP scheme to use for scraping.
	Scheme string `json:"scheme,omitempty"`
	// Optional HTTP URL parameters
	Params map[string][]string `json:"params,omitempty"`
	// Interval at which metrics should be scraped
	Interval string `json:"interval,omitempty"`
	// Timeout after which the scrape is ended
	ScrapeTimeout string `json:"scrapeTimeout,omitempty"`
	// Certificate configuration to use when scraping the endpoint
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty"`
	// File to read bearer token for scraping targets.
	BearerTokenFile string `json:"bearerTokenFile,omitempty"`
	// HonorLabels chooses the metric's labels on collisions with target labels.
	HonorLabels bool `json:"honorLabels,omitempty"`
	// BasicAuth allow an endpoint to authenticate over basic authentication
	// More info: https://prometheus.io/docs/operating/configuration/#endpoints
	BasicAuth *BasicAuth `json:"basicAuth,omitempty"`
	// MetricRelabelConfigs to apply to samples before ingestion.
	MetricRelabelConfigs []*RelabelConfig `json:"metricRelabelings,omitempty"`
	// RelabelConfigs to apply to samples before ingestion.
	// More info: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#<relabel_config>
	RelabelConfigs []*RelabelConfig `json:"relabelings,omitempty"`
	// ProxyURL eg http://proxyserver:2195 Directs scrapes to proxy through this endpoint.
	ProxyURL *string `json:"proxyUrl,omitempty"`
}

type TLSConfig struct {
	// The CA cert to use for the targets.
	CAFile string `json:"caFile,omitempty"`
	// The client cert file for the targets.
	CertFile string `json:"certFile,omitempty"`
	// The client key file for the targets.
	KeyFile string `json:"keyFile,omitempty"`
	// Used to verify the hostname for the targets.
	ServerName string `json:"serverName,omitempty"`
	// Disable target certificate validation.
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
}

type BasicAuth struct {
	// The secret that contains the username for authenticate
	Username v1.SecretKeySelector `json:"username,omitempty"`
	// The secret that contains the password for authenticate
	Password v1.SecretKeySelector `json:"password,omitempty"`
}

type RelabelConfig struct {
	//The source labels select values from existing labels. Their content is concatenated
	//using the configured separator and matched against the configured regular expression
	//for the replace, keep, and drop actions.
	SourceLabels []string `json:"sourceLabels,omitempty"`
	//Separator placed between concatenated source label values. default is ';'.
	Separator string `json:"separator,omitempty"`
	//Label to which the resulting value is written in a replace action.
	//It is mandatory for replace actions. Regex capture groups are available.
	TargetLabel string `json:"targetLabel,omitempty"`
	//Regular expression against which the extracted value is matched. defailt is '(.*)'
	Regex string `json:"regex,omitempty"`
	// Modulus to take of the hash of the source label values.
	Modulus uint64 `json:"modulus,omitempty"`
	//Replacement value against which a regex replace is performed if the
	//regular expression matches. Regex capture groups are available. Default is '$1'
	Replacement string `json:"replacement,omitempty"`
	// Action to perform based on regex matching. Default is 'replace'
	Action string `json:"action,omitempty"`
}

type NamespaceSelector struct {
	// Boolean describing whether all namespaces are selected in contrast to a
	// list restricting them.
	Any bool `json:"any,omitempty"`
	// List of namespace names.
	MatchNames []string `json:"matchNames,omitempty"`
}
