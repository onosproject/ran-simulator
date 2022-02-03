<!--
SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>

SPDX-License-Identifier: Apache-2.0
-->

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


To verify simulated e2 nodes are connected to the E2T endpoint successfully,  you can use onos-cli to check list of *CONTROL* (E2T to E2 node relation) and *CONTAINS* (E2 node to E2 cell relation) relations, and E2 node and cell entities using the following commands:

```bash
> onos topo get relations 
Relation ID                                 Kind ID    Source ID                      Target ID            Labels   Aspects
uuid:87d8f851-4394-4b6d-8775-95bf9d7d5c88   controls   e2:onos-e2t-56dcfbb685-knc8r   e2:1/5153            <None>   <None>
uuid:03c782b8-d993-62d3-5ada-8cde9bcc8d64   contains   e2:1/5154                      e2:1/5154/14550001   <None>   <None>
uuid:74c84ff1-74c2-388b-107e-8f62180b8aed   contains   e2:1/5153                      e2:1/5153/1454c002   <None>   <None>
uuid:e8d1924d-8a87-3840-ada0-0cacbef26cc5   contains   e2:1/5153                      e2:1/5153/1454c001   <None>   <None>
uuid:abee243a-ff82-4b0a-b037-782341b489ca   controls   e2:onos-e2t-56dcfbb685-knc8r   e2:1/5154            <None>   <None>
uuid:273c7b45-e7f3-ff52-43bd-891e86ff219d   contains   e2:1/5154                      e2:1/5154/14550002   <None>   <None>
uuid:826ab183-a742-79c2-aa83-a288ed68fa34   contains   e2:1/5154                      e2:1/5154/14550003   <None>   <None>
uuid:efe476d6-a6e4-7483-4c55-97c2ca884e73   contains   e2:1/5153                      e2:1/5153/1454c003   <None>   <None>

```bash
> onos topo get entities
Entity ID                      Kind ID   Labels   Aspects
e2:1/5154/14550002             e2cell    <None>   onos.topo.E2Cell
e2:1/5154/14550003             e2cell    <None>   onos.topo.E2Cell
e2:1/5154/14550001             e2cell    <None>   onos.topo.E2Cell
e2:1/5153/1454c003             e2cell    <None>   onos.topo.E2Cell
e2:onos-e2t-56dcfbb685-knc8r   e2t       <None>   onos.topo.Lease,onos.topo.E2TInfo
e2:1/5153                      e2node    <None>   onos.topo.MastershipState,onos.topo.E2Node
e2:1/5153/1454c001             e2cell    <None>   onos.topo.E2Cell
e2:1/5154                      e2node    <None>   onos.topo.E2Node,onos.topo.MastershipState
e2:1/5153/1454c002             e2cell    <None>   onos.topo.E2Cell
```

or use RAN simulator CLI to check status of E2 nodes:
```bash
> onos ransim get nodes
GnbID            Status   Service Models   E2T Controllers      Cell NCGIs
5153             Running  mho,rcpre2,kpm2  e2t-1                13842601454c001,13842601454c002,13842601454c003
5154             Running  mho,rcpre2,kpm2  e2t-1                138426014550001,138426014550002,138426014550003
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