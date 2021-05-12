### ThanosPeer
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | ThanosPeerSpec | No | - | See [ThanosPeerSpec](#thanospeerspec)<br> |
| status | ThanosPeerStatus | No | - | See [ThanosPeerStatus](#thanospeerstatus)<br> |
### ThanosPeerList
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []ThanosPeer | Yes | - |  |
### ThanosPeerSpec
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| endpointAddress | string | Yes | - | Host (or IP) and port of the remote Thanos endpoint<br> |
| peerEndpointAlias | string | No | - | Optional alias for the remote endpoint in case we have to access it through a different name.<br>This is typically needed if the remote endpoint has a certificate created for a predefined hostname.<br>The controller should create an externalName service for this backed buy the actual peer endpoint host<br>or a k8s service with a manually crafted k8s endpoint if EndpointAddress doesn't have a host but only an IP.<br> |
| certificate | string | No | - | The peer query should use this client certificate (tls.crt, tls.key) in the current namespace<br> |
| caBundle | string | No | - | Name of the secret that contains the CA certificate in ca.crt to verify client certs in the current namespace<br> |
| replicaLabels | []string | No | - | Custom replica labels if the default doesn't apply<br> |
| metaOverrides | typeoverride.ObjectMeta | No | - | [Override metadata](../overrides/override/#objectmeta) for managed resources<br> |
| queryOverrides | *Query | No | - | Override any of the [Query parameters](../thanos_types/#query)<br> |
### ThanosPeerStatus
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| queryHTTPServiceURL | string | No | - | The peer query is available over HTTP on this internal service URL<br> |
