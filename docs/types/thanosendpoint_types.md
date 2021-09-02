## ThanosEndpoint

###  (metav1.TypeMeta, required) {#thanosendpoint-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#thanosendpoint-metadata}

Default: -

### spec (ThanosEndpointSpec, optional) {#thanosendpoint-spec}

See [ThanosEndpointSpec](#thanosendpointspec) 

Default: -

### status (ThanosEndpointStatus, optional) {#thanosendpoint-status}

See [ThanosEndpointStatus](#thanosendpointstatus) 

Default: -


## ThanosEndpointList

###  (metav1.TypeMeta, required) {#thanosendpointlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#thanosendpointlist-metadata}

Default: -

### items ([]ThanosEndpoint, required) {#thanosendpointlist-items}

Default: -


## ThanosEndpointSpec

### certificate (string, optional) {#thanosendpointspec-certificate}

The endpoint should use this server certificate (tls.crt, tls.key) in the current namespace 

Default: -

### ingressClassName (string, optional) {#thanosendpointspec-ingressclassname}

Reference the given ingressClass resource explicitly 

Default: -

### caBundle (string, optional) {#thanosendpointspec-cabundle}

Name of the secret that contains the CA certificate in ca.crt to verify client certs in the current namespace 

Default: -

### stores ([]string, optional) {#thanosendpointspec-stores}

List of statically configured store addresses 

Default: -

### replicaLabels ([]string, optional) {#thanosendpointspec-replicalabels}

Custom replica labels if the default doesn't apply 

Default: -

### metaOverrides (typeoverride.ObjectMeta, optional) {#thanosendpointspec-metaoverrides}

[Override metadata](../overrides/override/#objectmeta) for managed resources 

Default: -

### queryOverrides (*Query, optional) {#thanosendpointspec-queryoverrides}

Override any of the [Query parameters](../thanos_types/#query) 

Default: -

### storeEndpointOverrides ([]StoreEndpointSpec, optional) {#thanosendpointspec-storeendpointoverrides}

Override any of the [StoreEndpoint parameters](../storeendpoint_types/) 

Default: -


## ThanosEndpointStatus

### endpointAddress (string, optional) {#thanosendpointstatus-endpointaddress}

Host (or IP) and port of the exposed Thanos endpoint 

Default: -


