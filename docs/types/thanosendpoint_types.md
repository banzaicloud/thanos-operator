### ThanosEndpoint
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | ThanosEndpointSpec | No | - | See [ThanosEndpointSpec](#thanosendpointspec)<br> |
| status | ThanosEndpointStatus | No | - | See [ThanosEndpointStatus](#thanosendpointstatus)<br> |
### ThanosEndpointList
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []ThanosEndpoint | Yes | - |  |
### ThanosEndpointSpec
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| certificate | string | No | - | The endpoint should use this server certificate (tls.crt, tls.key) in the current namespace<br> |
| ingressClassName | string | No | - | Reference the given ingressClass resource explicitly<br> |
| caBundle | string | No | - | Name of the secret that contains the CA certificate in ca.crt to verify client certs in the current namespace<br> |
| stores | []string | No | - | List of statically configured store addresses<br> |
| replicaLabels | []string | No | - | Custom replica labels if the default doesn't apply<br> |
| metaOverrides | typeoverride.ObjectMeta | No | - | [Override metadata](../overrides/override/#objectmeta) for managed resources<br> |
| queryOverrides | *Query | No | - | Override any of the [Query parameters](../thanos_types/#query)<br> |
| storeEndpointOverrides | []StoreEndpointSpec | No | - | Override any of the [StoreEndpoint parameters](../storeendpoint_types/)<br> |
### ThanosEndpointStatus
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| endpointAddress | string | No | - | Host (or IP) and port of the exposed Thanos endpoint<br> |
