### ServiceMonitor
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - | Standard object’s metadata. More info:<br>https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata<br>+k8s:openapi-gen=false<br> |
| spec | ServiceMonitorSpec | Yes | - | Specification of desired Service selection for target discrovery by<br>Prometheus.<br> |
### ServiceMonitorList
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - | Standard list metadata<br>More info: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#metadata<br> |
| items | []*ServiceMonitor | Yes | - | List of ServiceMonitors<br> |
### ServiceMonitorSpec
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| jobLabel | string | No | - | The label to use to retrieve the job name from.<br> |
| targetLabels | []string | No | - | TargetLabels transfers labels on the Kubernetes Service onto the target.<br> |
| podTargetLabels | []string | No | - | PodTargetLabels transfers labels on the Kubernetes Pod onto the target.<br> |
| endpoints | []Endpoint | Yes | - | A list of endpoints allowed as part of this ServiceMonitor.<br> |
| selector | metav1.LabelSelector | Yes | - | Selector to select Endpoints objects.<br> |
| namespaceSelector | NamespaceSelector | No | - | Selector to select which namespaces the Endpoints objects are discovered from.<br> |
| sampleLimit | uint64 | No | - | SampleLimit defines per-scrape limit on number of scraped samples that will be accepted.<br> |
### Endpoint
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| port | string | No | - | Name of the service port this endpoint refers to. Mutually exclusive with targetPort.<br> |
| targetPort | *intstr.IntOrString | No | - | Name or number of the target port of the endpoint. Mutually exclusive with port.<br> |
| path | string | No | - | HTTP path to scrape for metrics.<br> |
| scheme | string | No | - | HTTP scheme to use for scraping.<br> |
| params | map[string][]string | No | - | Optional HTTP URL parameters<br> |
| interval | string | No | - | Interval at which metrics should be scraped<br> |
| scrapeTimeout | string | No | - | Timeout after which the scrape is ended<br> |
| tlsConfig | *TLSConfig | No | - | TLS configuration to use when scraping the endpoint<br> |
| bearerTokenFile | string | No | - | File to read bearer token for scraping targets.<br> |
| honorLabels | bool | No | - | HonorLabels chooses the metric's labels on collisions with target labels.<br> |
| basicAuth | *BasicAuth | No | - | BasicAuth allow an endpoint to authenticate over basic authentication<br>More info: https://prometheus.io/docs/operating/configuration/#endpoints<br> |
| metricRelabelings | []*RelabelConfig | No | - | MetricRelabelConfigs to apply to samples before ingestion.<br> |
| relabelings | []*RelabelConfig | No | - | RelabelConfigs to apply to samples before ingestion.<br>More info: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#<relabel_config><br> |
| proxyUrl | *string | No | - | ProxyURL eg http://proxyserver:2195 Directs scrapes to proxy through this endpoint.<br> |
### TLSConfig
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| caFile | string | No | - | The CA cert to use for the targets.<br> |
| certFile | string | No | - | The client cert file for the targets.<br> |
| keyFile | string | No | - | The client key file for the targets.<br> |
| serverName | string | No | - | Used to verify the hostname for the targets.<br> |
| insecureSkipVerify | bool | No | - | Disable target certificate validation.<br> |
### BasicAuth
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| username | v1.SecretKeySelector | No | - | The secret that contains the username for authenticate<br> |
| password | v1.SecretKeySelector | No | - | The secret that contains the password for authenticate<br> |
### RelabelConfig
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| sourceLabels | []string | No | - | The source labels select values from existing labels. Their content is concatenated<br>using the configured separator and matched against the configured regular expression<br>for the replace, keep, and drop actions.<br> |
| separator | string | No | - | Separator placed between concatenated source label values. default is ';'.<br> |
| targetLabel | string | No | - | Label to which the resulting value is written in a replace action.<br>It is mandatory for replace actions. Regex capture groups are available.<br> |
| regex | string | No | - | Regular expression against which the extracted value is matched. defailt is '(.*)'<br> |
| modulus | uint64 | No | - | Modulus to take of the hash of the source label values.<br> |
| replacement | string | No | - | Replacement value against which a regex replace is performed if the<br>regular expression matches. Regex capture groups are available. Default is '$1'<br> |
| action | string | No | - | Action to perform based on regex matching. Default is 'replace'<br> |
### NamespaceSelector
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| any | bool | No | - | Boolean describing whether all namespaces are selected in contrast to a<br>list restricting them.<br> |
| matchNames | []string | No | - | List of namespace names.<br> |
