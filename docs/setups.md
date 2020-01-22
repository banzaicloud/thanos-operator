
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