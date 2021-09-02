## StoreEndpointSpec

StoreEndpointSpec defines the desired state of StoreEndpoint

### metaOverrides (*typeoverride.ObjectMeta, optional) {#storeendpointspec-metaoverrides}

See [ObjectMeta override](../overrides/override/#objectmeta) 

Default: -

### serviceOverrides (*typeoverride.Service, optional) {#storeendpointspec-serviceoverrides}

See [Service override](../overrides/override/#service) 

Default: -

### url (string, optional) {#storeendpointspec-url}

Default: -

### selector (*KubernetesSelector, optional) {#storeendpointspec-selector}

See [KubernetesSelector](#kubernetesselector) 

Default: -

### config (secret.Secret, optional) {#storeendpointspec-config}

Default: -

### thanos (string, required) {#storeendpointspec-thanos}

Default: -

### ingress (*Ingress, optional) {#storeendpointspec-ingress}

Default: -


## KubernetesSelector

### namespaces (string, optional) {#kubernetesselector-namespaces}

Default: -

### labels (map[string]string, optional) {#kubernetesselector-labels}

Default: -

### annotations (map[string]string, optional) {#kubernetesselector-annotations}

Default: -

### httpPort (int32, optional) {#kubernetesselector-httpport}

Default: -

### grpcPort (int32, optional) {#kubernetesselector-grpcport}

Default: -


## StoreEndpointStatus

StoreEndpointStatus defines the observed state of StoreEndpoint


## StoreEndpoint

StoreEndpoint is the Schema for the storeendpoints API

###  (metav1.TypeMeta, required) {#storeendpoint-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#storeendpoint-metadata}

Default: -

### spec (StoreEndpointSpec, optional) {#storeendpoint-spec}

See [StoreEndpointSpec](#storeendpointspec) 

Default: -

### status (StoreEndpointStatus, optional) {#storeendpoint-status}

See [StoreEndpointStatus](#storeendpointstatus) 

Default: -


## StoreEndpointList

StoreEndpointList contains a list of StoreEndpoint

###  (metav1.TypeMeta, required) {#storeendpointlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#storeendpointlist-metadata}

Default: -

### items ([]StoreEndpoint, required) {#storeendpointlist-items}

Default: -


