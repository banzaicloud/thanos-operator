### ThanosEndpoint
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | ThanosEndpointSpec | No | - |  |
| status | ThanosEndpointStatus | No | - |  |
### ThanosEndpointList
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []ThanosEndpoint | Yes | - |  |
### ThanosEndpointSpec
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| certificate | string | No | - | The endpoint should use this server certificate<br> |
| caBundle | string | No | - | CA certificate to verify client certs<br> |
| stores | []string | No | - | List of statically configured store addresses<br> |
| replicaLabels | []string | No | - | Custom replica labels if the default doesn't apply<br> |
### ThanosEndpointStatus
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| endpointAddress | string | No | - | Host (or IP) and port of the exposed Thanos endpoint<br> |
