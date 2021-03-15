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
````
2) Deploy sd-ran chart using the following command:

```bash
kubectl create namespace sd-ran
cd sdran-helm-charts
helm install sd-ran sd-ran -n sd-ran
```
if you deploy sd-ran chart successfully, you should be able to see
sdran subsystems deployed successfully as follows:
```bash
 kubectl get pods -n sd-ran
NAME                           READY   STATUS    RESTARTS   AGE
onos-cli-5d8b489f69-595l7      1/1     Running   0          4m31s
onos-config-79b7f4bb65-9qllf   1/1     Running   0          4m31s
onos-consensus-db-1-0          1/1     Running   0          4m31s
onos-e2sub-696574c74d-dnhqf    1/1     Running   0          4m31s
onos-e2t-5646b945cd-zv8d6      1/1     Running   0          4m31s
onos-topo-84d9847bd4-669qn     1/1     Running   0          4m31s
```

3) Deploy ran-simulator helm chart using the following command:
Ran simulator is not enabled in the sd-ran chart by default. You can enable it
when you deploy sd-ran helm chart using the following command 

```bash
helm install sd-ran sd-ran -n sd-ran --set import.ran-simulator.enabled=true
````   
or you can deploy the sd-ran chart first and then
deploy ran-simulator using the following command:
   
```bash
helm install ran-simulator ran-simulator -n sd-ran
```

if you deploy RAN simulator successfully, you should be able to see it
in the list of deployments:

```bash
kubectl get pods -n sd-ran
NAME                             READY   STATUS    RESTARTS   AGE
onos-cli-5d8b489f69-595l7        1/1     Running   0          12m
onos-config-79b7f4bb65-9qllf     1/1     Running   0          12m
onos-consensus-db-1-0            1/1     Running   0          12m
onos-e2sub-696574c74d-dnhqf      1/1     Running   0          12m
onos-e2t-5646b945cd-zv8d6        1/1     Running   0          12m
onos-topo-84d9847bd4-669qn       1/1     Running   0          12m
ran-simulator-6d9c89cdc7-6t6tl   1/1     Running   0          99s
```

After deploying ran-simulator, it loads the models, create E2 nodes and make connections using SCTP
to onos-e2t endpoint which is specified in the model. 

To verify simulated e2 nodes are connected to the E2T endpoint successfully, 
you can use onos-cli to check list of E2 connections using the following command:

```bash
> onos e2t list connections
Global ID            PLNM ID   IP Addr        Port    Conn Type
0000000000340422:0   1279014   10.244.0.247   39772   G_NB
00000000003020f9:0   1279014   10.244.0.247   53406   G_NB
```

or use RAN simulator CLI to check status of E2 nodes:
```bash
> onos ransim get nodes
EnbID            Status   Service Models   E2T Controllers      Cell ECGIs
5153             Running  kpm,rc           e2t-1                21458294227473,21458294227474,21458294227475
5154             Running  kpm,rc           e2t-1                21458294227489,21458294227490,21458294227475
```


[Helm]: https://helm.sh/
[Kubernetes]: https://kubernetes.io/
[kind]: https://kind.sigs.k8s.io
[Directions API]: https://developers.google.com/maps/documentation/directions/start
[sdran-helm-charts]: https://github.com/onosproject/sdran-helm-charts