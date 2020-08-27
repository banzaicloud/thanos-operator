---
title: StoreEndpointSpec
weight: 200
generated_file: true
---

### StoreEndpointSpec
#### StoreEndpointSpec defines the desired state of StoreEndpoint

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| url | string | No | - | Foo is an example field of StoreEndpoint. Edit StoreEndpoint_types.go to remove/update<br> |
| selector | *KubernetesSelector | No | - |  |
| config | secret.Secret | No | - |  |
| thanos | string | Yes | - |  |
| ingress | *Ingress | No | - |  |
### KubernetesSelector
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| namespaces | string | No | - |  |
| labels | map[string]string | No | - |  |
| annotations | map[string]string | No | - |  |
| httpPort | int32 | No | - |  |
| grpcPort | int32 | No | - |  |
### StoreEndpointStatus
#### StoreEndpointStatus defines the observed state of StoreEndpoint

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
### StoreEndpoint
#### StoreEndpoint is the Schema for the storeendpoints API

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | StoreEndpointSpec | No | - |  |
| status | StoreEndpointStatus | No | - |  |
### StoreEndpointList
#### StoreEndpointList contains a list of StoreEndpoint

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []StoreEndpoint | Yes | - |  |
