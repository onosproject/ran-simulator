<!--
SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>

SPDX-License-Identifier: Apache-2.0
-->

# RAN simulator CLI

RAN simulator is equipped with a command line interface which is integrated with
[onos-cli](https://github.com/onosproject/onos-cli) that allows to interact with Ransim to retrieve required information from data stores,
monitor RAN environment changes, create/remove/update RAN entities, metrics, etc.
The list of ransim commands is documented here [ransim-cli](https://github.com/onosproject/onos-cli/blob/master/docs/cli/onos_ransim.md) 

## Usage 

1) Follow the instructions in [Quick Start](quick_start.md) to deploy 
sd-ran subsystems and RAN simulator. 
   
2) Use the following command to access to `onos-cli` that you can run the [RAN simulator commands][ransim-cli]:

```bash
kubectl exec -it -n sd-ran onos-cli-5d8b489f69-nvfcm -- /bin/bash
```

_TODO: Update CLI usage to include routes and UEs!_

3) Run onos ransim --help to see the list of ransim commands:
```bash
ONOS RAN simulator commands

Usage:
  onos ransim [command]

Available Commands:
  clear       Clear the simulated nodes, cells and metrics
  config      Manage the CLI configuration
  create      Commands for creating simulated entities
  delete      Commands for deleting simulated entities
  get         Commands for retrieving RAN simulator model and other information
  load        Load model and/or metric data
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
GnbID            Status   Service Models   E2T Controllers      Cell NCGIs
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
GnbID            Status   Service Models   E2T Controllers      Cell NCGIs
5153             Running  kpm,rc           e2t-1                21458294227473,21458294227474,21458294227475
5154             Running  kpm,rc           e2t-1                21458294227489,21458294227490,21458294227475
5155             Running  kpm,rc           e2t-1                21458294227489,21458294227490,21458294227491
```

6) The following command displays the cell information, including the number of UEs per cell that are in RRC State Idle and Connected.

```bash
$ onos ransim get cells
NCGI                 #UEs Max UEs    TxDB       Lat       Lng Azimuth     Arc   A3Offset     TTT  A3Hyst CellOffset FreqOffset      PCI    Color Idle Conn Neighbors
13842601454c002         0   99999   11.00    52.486    13.412     120     120          0       0       0          0          0      148    green   49,   17, 13842601454c001,13842601454c003,1
38426014550002
13842601454c003         0   99999   11.00    52.486    13.412     240     120          0       0       0          0          0      480    green   92,  102, 13842601454c001,13842601454c002,1
```


[ransim-cli]: https://github.com/onosproject/onos-cli/blob/master/docs/cli/onos_ransim.md
