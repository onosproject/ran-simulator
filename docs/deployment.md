# Deploying ran-simulator

This guide deploys `ran-simulator` through it's [Helm] chart assumes you have a
[Kubernetes] cluster running deployed in a namespace.

`ran-simulator` Helm chart is based on Helm 3.0 version, with no need for the Tiller pod to be present. 
If you don't have a cluster running and want to try on your local machine please follow first 
the [Kubernetes] setup steps outlined in [deploy with Helm](https://docs.onosproject.org/developers/deploy_with_helm/).
The following steps assume you have the setup outlined in that page, including the `micro-onos` namespace configured. 

## Google Maps API Key
The RAN Simulator can connect to Google's [Directions API] and with a Google API Key.
Google charges $5.00 per 1000 requests to the [Directions API], and so we do not put
our API key up in the public domain.

> If a Google API key is not given, the simulator will revert to a built in random route generator

## Tuning parameters
[README.md](README.md) shows the list of startup parameters - these can be overridden using the
```--set param=value``` syntax.

## Installing the Chart
To install the chart in the `micro-onos` namespace run from the root directory of
the `onos-helm-charts` repo the command:
```bash
helm install -n micro-onos ran-simulator ran-simulator
```
The output should be:
```bash
NAME: ran-simulator
LAST DEPLOYED: Tue Feb  4 08:02:57 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

To start with a Google API key run the command like:
```bash
helm install -n micro-onos ran-simulator ran-simulator --set googleAPIKey=YOUR_API_KEY_HERE
```

`helm install` assigns a unique name to the chart and displays all the k8s resources that were
created by it. To list the charts that are installed and view their statuses, run `helm ls`:

```bash
> helm ls
NAME         	NAMESPACE 	REVISION	UPDATED                                	STATUS  	CHART              	APP VERSION
onos-cli     	micro-onos	1       	2020-02-04 08:01:54.860813386 +0000 UTC	deployed	onos-cli-0.0.1     	1          
onos-ric     	micro-onos	1       	2020-02-04 08:02:17.663782372 +0000 UTC	deployed	onos-ric-0.0.1     	1          
ran-simulator	micro-onos	1       	2020-02-04 09:32:21.533299519 +0000 UTC	deployed	ran-simulator-0.0.1	1          
sd-ran-gui   	micro-onos	1       	2020-02-04 09:32:49.018099586 +0000 UTC	deployed	sd-ran-gui-0.0.1   	1  
```

> Here the service is shown running alongside `onos-cli`, `onos-ric` and the `sd-ran-gui`
> as these are usually deployed together to give a demo scenario. See the individual
> deployment instructions for these services.

To check whether the service is running correctly use `kubectl`
```bash
> kubectl get pods -n micro-onos
NAME                             READY   STATUS             RESTARTS   AGE
onos-cli-68bbf4f674-ssjt4        1/1     Running            0          18m
onos-ric-5fb8c6bdd7-xmcmq        1/1     Running            0          18m
ran-simulator-6f577597d8-5lcv8   1/1     Running            0          82s
sd-ran-gui-76ff54d85-fh72j       2/2     Running            0          54s
```

See Troubleshooting below if the `Status` is not `Running`

### Installing the chart in a different namespace.

Issue the `helm install` command substituting `micro-onos` with your namespace.
```bash
helm install -n <your_name_space> ran-simulator ran-simulator
```

### Troubleshoot
If your chart does not install or the pod is not running for some reason and/or you modified values Helm offers two flags to help you
debug your chart:  

* `--dry-run` check the chart without actually installing the pod. 
* `--debug` prints out more information about your chart

```bash
helm install -n micro-onos ran-simulator --debug --dry-run ran-simulator/
```

#### CrashLoopBackOff
If the ran-simulator has a status of `CrashLoopBackOff` it is usually an indication
that the Google API Key is not valid. Un-deploy the ran-simulator, update the key
in `values.yaml` and redeploy it again.

Looking at the logs of the pod shows where this is the case
```bash
> kubectl -n micro-onos logs ran-simulator-6f65c8b57-rq4gf
E0204 08:24:24.305652       1 trafficsim.go:69] Cant' avoid double Error logging no such flag -alsologtostderr
I0204 08:24:24.305777       1 trafficsim.go:111] Starting trafficsim
I0204 08:24:24.305804       1 manager.go:40] Creating Manager
I0204 08:24:24.305818       1 manager.go:51] Starting Manager with {3 3 0.02 0.02} {10} {3 YOUR_API_KEY_HERE 1s}
I0204 08:24:24.305953       1 dispatcher.go:41] User Equipment Event listener initialized
I0204 08:24:24.306070       1 dispatcher.go:79] Route Event listener initialized
F0204 08:24:24.553960       1 manager.go:62] Error calculating routes maps: REQUEST_DENIED - The provided API key is invalid.
``` 

## Uninstalling the chart.

To remove the `ran-simulator` pod issue
```bash
 helm delete -n micro-onos ran-simulator
```

## Pod Information

To view the pods that are deployed, run `kubectl -n micro-onos get pods`.

[Helm]: https://helm.sh/
[Kubernetes]: https://kubernetes.io/
[kind]: https://kind.sigs.k8s.io
[Directions API]: https://developers.google.com/maps/documentation/directions/start
