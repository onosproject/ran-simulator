# Quick Start 

## Prerequisite

This guide deploys `ran-simulator` through it's [Helm] chart assumes you have a
[Kubernetes] cluster running deployed in a namespace.

`ran-simulator` Helm chart is based on Helm 3.0 version, with no need for the Tiller pod to be present.
If you don't have a cluster running and want to try on your local machine please follow first
the [Kubernetes] setup steps outlined in [deploy with Helm](https://docs.onosproject.org/developers/deploy_with_helm/).
The following steps assume you have the setup outlined in that page, including the `micro-onos` namespace configured.

## Usage with SD-RAN Subsystems

The steps for using Ransim with SD-RAN subsystems are as follows:

1) Clone [sdran-helm-charts][sdran-helm-charts] using the following command:
```bash
git clone git@github.com:onosproject/sdran-helm-charts.git
```

2) Deploy sd-ran chart using the following command:

```bash
kubectl create namespace sd-ran
cd sdran-helm-charts
helm install sd-ran sd-ran -n sd-ran
```
If you deploy the sd-ran chart successfully, you should be able to see sdran subsystems deployed successfully as follows after a period of time:
```bash
> kubectl get pods -n sd-ran
NAME                           READY   STATUS    RESTARTS   AGE
onos-cli-6f45d4b475-cw78c      1/1     Running   0          21s
onos-config-6f648c5b57-27vlk   4/4     Running   0          20s
onos-consensus-db-1-0          1/1     Running   0          21s
onos-consensus-store-1-0       1/1     Running   0          18s
onos-e2t-5698597f6c-rbswj      3/3     Running   0          20s
onos-topo-66c7757f6d-t84r9     3/3     Running   0          21s
onos-uenib-6b8bd5cddf-68nm2    3/3     Running   0          21s
```

3) Deploy ran-simulator helm chart using the following command:

```bash
helm install ran-simulator ran-simulator -n sd-ran
```
RAN simulator is not enabled in the sd-ran chart by default. You can enable it when you deploy sd-ran helm chart using the following command 
```bash
helm install sd-ran sd-ran -n sd-ran --set import.ran-simulator.enabled=true
```

If you deploy RAN simulator successfully, you should be able to see it in the list of deployments:

```bash
> kubectl get pods -n sd-ran
NAME                             READY   STATUS    RESTARTS   AGE
onos-cli-6f45d4b475-cw78c        1/1     Running   0          3m21s
onos-config-6f648c5b57-27vlk     4/4     Running   0          3m20s
onos-consensus-db-1-0            1/1     Running   0          3m21s
onos-consensus-store-1-0         1/1     Running   0          3m18s
onos-e2t-5698597f6c-rbswj        3/3     Running   0          3m20s
onos-topo-66c7757f6d-t84r9       3/3     Running   0          3m21s
onos-uenib-6b8bd5cddf-68nm2      3/3     Running   0          3m21s
ran-simulator-67bb8894cd-8jgbd   1/1     Running   0          3m21s
```

After deploying ran-simulator, it loads the models, create E2 nodes and make connections using SCTP to onos-e2t endpoint which is specified in the model. 


To verify simulated e2 nodes are connected to the E2T endpoint successfully,  you can use onos-cli to check list of E2 connections using the following command:

```bash
> onos e2t get connections
Global ID            PLNM ID   IP Addr        Port    Conn Type
0000000000340422:0   1279014   10.244.0.247   39772   G_NB
00000000003020f9:0   1279014   10.244.0.247   53406   G_NB
```
or use RAN simulator CLI to check status of E2 nodes:
```bash
> onos ransim get nodes
GnbID            Status   Service Models   E2T Controllers      Cell NCGIs
5154             Running  kpm,rcpre2,kpm2,mho e2t-1                138426014550001,138426014550002,138426014550003
5153             Running  kpm,rcpre2,kpm2,mho e2t-1                13842601454c001,13842601454c002,13842601454c003
```

If you have not installed the onos cli as described in the [cli docs](https://docs.onosproject.org/onos-cli/docs/setup/), you can instead run 
```
clipod=$(kubectl -n sd-ran get pods | grep onos-cli | cut -d\  -f1)
kubectl -n sd-ran exec --stdin $clipod -- /usr/local/bin/onos e2t get connections
```
or 
```
clipod=$(kubectl -n sd-ran get pods | grep onos-cli | cut -d\  -f1)
kubectl -n sd-ran exec --stdin $clipod -- /usr/local/bin/onos ransim get nodes
```


[Helm]: https://helm.sh/
[Kubernetes]: https://kubernetes.io/
[kind]: https://kind.sigs.k8s.io
[Directions API]: https://developers.google.com/maps/documentation/directions/start
[sdran-helm-charts]: https://github.com/onosproject/sdran-helm-charts