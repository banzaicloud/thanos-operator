### ThanosPeer
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | ThanosPeerSpec | No | - |  |
| status | ThanosPeerStatus | No | - |  |
### ThanosPeerList
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []Thanos | Yes | - |  |
### ThanosPeerSpec
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| endpointAddress | string | Yes | - | Host (or IP) and port of the remote Thanos endpoint<br> |
| peerEndpointAlias | string | No | - | Optional alias for the remote endpoint in case we have to access it through a different name.<br>This is typically needed if the remote endpoint has a certificate created for a predefined hostname.<br>The controller should create an externalName service for this backed buy the actual peer endpoint host<br>or a k8s service with a manually crafted k8s endpoint if EndpointAddress doesn't have a host but only an IP.<br> |
| caBundle | string | No | - | CA certificate to verify the server cert<br> |
| replicaLabels | []string | No | - | Custom replica labels if the default doesn't apply<br> |
### ThanosPeerStatus
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
