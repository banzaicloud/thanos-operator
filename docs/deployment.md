<p align="center"><img src="./img/logo/thanos_operator_vertical.svg" width="260"></p>
<p align="center">

# Requirements

- Thanos operator requires Kubernetes v1.14.x or later.
- For the [Helm basde installation](#deploy-thanos-operator-with-helm) you need Helm v3.0.2 or later.


# Deploy Thanos operator with Helm

<p align="center"><img src="./img/logo/helm.svg" width="150"></p>
<p align="center">

Complete the following steps to deploy the Logging operator using Helm. Alternatively, you can also [install the operator using Kubernetes manifests](./Readme.md).
> Note: For the [Helm base installation](#deploy-thanos-operator-with-helm) you need Helm v3 or later.


1. Create `monitor` namespace
    ```bash
    kubectl create namespace monitor
    ```
1. Add the operator chart repository.
    ```bash
    helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com
    helm repo update
    ```
1. Install the Thanos Operator
    ```bash
    helm install thanos-operator --namespace monitor banzaicloud-stable/thanos-operator --set manageCrds=false
    ```

---

# Check the Thanos operator deployment

To verify that the installation was successful, complete the following steps.

1. Check the status of the pods. You should see a new thanos-operator pod.
    ```bash
    $ kubectl -n monitor get pods
    NAME                                        READY   STATUS    RESTARTS   AGE
    thanos-operator-7df8485bf6-gf5gk   1/1     Running   0          13s
    ```
1. Check the CRDs. You should see the following three new CRDs.
    ```bash
    $  kubectl get crd
    NAME                                    CREATED AT
    objectstores.monitoring.banzaicloud.io      2020-02-07T21:48:20Z
    storeendpoints.monitoring.banzaicloud.io    2020-02-07T21:48:20Z
    thanos.monitoring.banzaicloud.io            2020-02-07T21:48:20Z
    ```
