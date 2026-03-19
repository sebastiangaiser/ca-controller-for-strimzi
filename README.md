# ca-controller-for-strimzi

> :warning: **This project is not in cooperation with Strimzi**: I asked one of the Strimzi core maintainer about
> contributing this controller upstream but it makes more sense to implement the logic in Strimzi itself.

When dealing with [strimzi-kafka-operator](https://github.com/strimzi/strimzi-kafka-operator/), it is possible to use
your own CA for cluster and clients.
Strimzi requires the CA key split from the rest of a Kubernetes secret of type TLS.
To avoid doing this manually e.g. when using [cert-manager](https://cert-manager.io/) for managing the CAs, this
controller can be used...

Please check the `examples/example-ca.yaml` how to use the controller after deploying it and using it with cert-manager but
it also works with normal Kubernetes secrets of type TLS.

## Configuration

### Required annotations

The controller watches secrets of type `kubernetes.io/tls` that carry the following annotations. Set them via cert-manager's `secretTemplate` or directly on a hand-crafted secret.

| Annotation | Description |
|---|---|
| `sebastian.gaiser.bayern/tls-strimzi-ca: "reconcile"` | Opt the secret into reconciliation |
| `sebastian.gaiser.bayern/target-cluster-name` | Name of the Strimzi Kafka cluster |
| `sebastian.gaiser.bayern/target-secret-name` | Name of the target certificate secret (receives `ca.crt` and `tls.crt`) |
| `sebastian.gaiser.bayern/target-secret-key-name` | Name of the target private key secret (receives `ca.key`) |

### Private key rotation policy

By default the controller keeps the private key secret in sync with the source secret on every reconciliation. If you configure cert-manager to never rotate the private key (`rotationPolicy: Never`), you should tell the controller the same so it does not overwrite the key secret on certificate renewals.

Add the annotation `sebastian.gaiser.bayern/rotation-policy: "Never"` to the cert-manager `secretTemplate`. Use a YAML anchor to reference the value from `privateKey.rotationPolicy` directly so both fields are always in sync:

```yaml
spec:
  privateKey:
    rotationPolicy: &rotationPolicy Never
  secretTemplate:
    annotations:
      sebastian.gaiser.bayern/rotation-policy: *rotationPolicy
```

With this annotation set the controller will:

- **Create** the private key secret on the first reconciliation (Strimzi requires it to exist).
- **Skip updating** the private key secret on subsequent reconciliations, even when the certificate is renewed.
- **Skip creating historical snapshots** of the private key secret.

The certificate secret (`ca.crt`, `tls.crt`) is always kept up to date regardless of this setting.

## Install the controller via Helm

```shell
helm repo add sebastiangaiser-ca-controller-for-strimzi https://sebastiangaiser.github.io/ca-controller-for-strimzi
helm repo update
helm upgrade --install -n kafka ca-controller-for-strimzi sebastiangaiser-ca-controller-for-strimzi/ca-controller-for-strimzi
```

## See the controller in action (using local kind cluster)

```shell
# create cluster
kind create cluster
# install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.20.0/cert-manager.yaml
# bootstrap strimzi in kafka namespace
kubectl create namespace kafka
helm install -n kafka strimzi-cluster-operator oci://quay.io/strimzi-helm/strimzi-kafka-operator
# create certificates
kubectl apply -f examples/example-ca.yaml
# install controller
helm repo add sebastiangaiser-ca-controller-for-strimzi https://sebastiangaiser.github.io/ca-controller-for-strimzi
helm repo update
helm upgrade --install -n kafka ca-controller-for-strimzi sebastiangaiser-ca-controller-for-strimzi/ca-controller-for-strimzi
# create Kafka cluster
kubectl apply -f examples/kafka.yaml -n kafka
```
