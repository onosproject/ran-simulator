# RAN simulator CLI

RAN simulator is equipped with a command line interface which is integrated with
[onos-cli](https://github.com/onosproject/onos-cli) that allows to interact with Ransim to retrieve required information from data stores,
monitor RAN environment changes, create/remove/update RAN entities, metrics, etc.
The list of ransim commands is documented here [ransim-cli](https://github.com/onosproject/onos-cli/blob/master/docs/cli/onos_ransim.md) 

## Usage 

1) Follow the instructions in [Quick Start](quick_start.md) to deploy 
sd-ran subsystems and RAN simulator. 
   
2) Use the following command to access to onos-cli that you can run onos ransim 
commands:

```bash
kubectl exec -it -n sd-ran onos-cli-5d8b489f69-nvfcm -- /bin/bash
```

3) Run onos ransim --help to see the list of ransim commands:
```bash
ONOS RAN simulator commands

Usage:
  onos ransim [command]

Available Commands:
  config      Manage the CLI configuration
  create      Commands for creating simulated entities
  delete      Commands for deleting simulated entities
  get         Commands for retrieving RAN simulator model and other information
  log         logging api commands
  set         Commands for setting RAN simulator model metrics and other information
  start       Start E2 node agent
  stop        Stop E2 node agent

Flags:
      --auth-header string       Auth header in the form 'Bearer <base64>'
  -h, --help                     help for ransim
      --no-tls                   if present, do not use TLS
      --service-address string   the gRPC endpoint (default "ran-simulator:5150")
      --tls-cert-path string     the path to the TLS certificate
      --tls-key-path string      the path to the TLS key

Use "onos ransim [command] --help" for more information about a command.
```

4) For example, the following command lists all of  E2 nodes that are running
```bash
> onos ransim get nodes
EnbID            Status   Service Models   E2T Controllers      Cell ECGIs
5153             Running  kpm,rc           e2t-1                21458294227473,21458294227474,21458294227475
5154             Running  kpm,rc           e2t-1                21458294227489,21458294227490,21458294227475
```

5) As another example, you can use the following to create an E2 node:
```bash
onos ransim create node 5155 --service-models kpm --service-models rc --controllers e2t-1 --cells 21458294
227489 --cells 21458294227490 --cells 21458294227491
```

after running the above command, an e2 node will be created and will be 
connected to e2t-1 which is specified in the model as a E2T endpoint.
```bash
> onos ransim get nodes 
EnbID            Status   Service Models   E2T Controllers      Cell ECGIs
5153             Running  kpm,rc           e2t-1                21458294227473,21458294227474,21458294227475
5154             Running  kpm,rc           e2t-1                21458294227489,21458294227490,21458294227475
5155             Running  kpm,rc           e2t-1                21458294227489,21458294227490,21458294227491
```