# ca-controller-for-strimzi

When dealing with [strimzi-kafka-operator](https://github.com/strimzi/strimzi-kafka-operator/), it is possible to use
you own CA for cluster and clients.
Strimzi requires the CA key split from the rest of a Kubernetes secret of type TLS.
To avoid doing this manually e.g. when using [cert-manager](https://cert-manager.io/) for managing the CAs, this
controller can be used...

Please check the `example-ca.yaml` for how to use the controller after deploying it and using it with cert-manager, but
it also works with normal Kubernetes secrets of type TLS.

> :warning: **This project is build on my own and is not related to Strimzi**: I asked "a" Strimzi core maintainer about
> contributing this controller upstream but it makes more sense to implement the logic in strimzi itself.
