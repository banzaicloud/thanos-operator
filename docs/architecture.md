<p align="center"><img src="./img/logo/thanos_operator_vertical.svg" width="260"></p>
<p align="center">

# Thanos Operator
 
## Concepts
 
## Custom Resources

### ObjectStore

This resource responsible for operations done per object store.
- Compactor
  - Must run 1 instance per bucket
  - Responsible for downsampling and compacting
- Bucket
  - Reporting tool for buckets
  - Simple Web UI
- `config`
  - Credential and configuration for object storage

*example*
```
apiVersion: monitoring.banzaicloud.io/v1alpha1
kind: ObjectStore
metadata:
  name: objectstore-sample
spec:
  config:
    mountFrom:
      secretKeyRef:
        name: thanos
        key: object-store.yaml
  bucketWeb: {}
  compactor: {}
    dataVolume:
      pvc:
        spec:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 1Gi
```

### Thanos

- Query
  - HTTP frontend
- StoreGateway
  - 
- Rule
  - 
  
### StoreEndpoint
- Define local or remote StoreAPI providers
- Sidecar, Rule, Query, ...
- Local (Selector) or Remote (URL) configuration
- Each `Thanos` and `StoreEndpoint` make a unique Stack with separate instances

## Deployment modes