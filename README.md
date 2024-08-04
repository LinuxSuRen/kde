# kde
Kubernetes based remote IDE

## Setup Project

I setup this project using the following commands:

* kubebuilder init --plugins=go.kubebuilder.io/v4
* kubebuilder create api --group linuxsuren.github.io --version v1alpha1 --kind DevSpace

## Install OpenEBS

```shell
kubectl apply -f https://openebs.github.io/charts/openebs-operator.yaml
```
