
# Example deployment modes for Thanos

## Single

- Install Prometheus Operator
- Install Thanos Operator

Thanos component discovery
- Discover all Thanos component provisioned by operator
- Discover prometheus sidecars

## Remote Thanos for Prometheus

- Remote URLs
  - format: `http(s)://<fqdn>:<port>`
  - tls: ?

## Multi Thanos Observer

- Select endpoints based on labels
  - Namespace?
  - Labels?
  - Pods/Svc?
  
## Query discovery
Automatically discover all query on the cluster
- Only creted by Operator?

## Sidecar SD
As sidecars is only useful for Query it should be under query key
sidecars:
  - namespaces (optional) | []string | default: same namespace
    labels (optional) | default: app=prometheus
    url (optional) | default: ""
    
Based on `namespaces` and `labels` the operator creates services and adds
them to query parameters.

The `url` attribute has priority and requires FQDN