<!--
SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>

SPDX-License-Identifier: Apache-2.0
-->

# Simulation Models

RAN simulator defines two levels of simulation models:

* **Generic model:** that defines E2 nodes, cells, service models, E2T end points in a Yaml file. RAN simulator reads this file to create E2 nodes and initializing data stores. A sample [model](https://github.com/onosproject/sdran-helm-charts/blob/master/ran-simulator/files/model/model.yaml) has been created using the [Honeycomb Topology Generator](topology_generator.md) and added to the [ran-simulator helm chart][RAN simulator helm chart]. 


* **Use Case Specific Models**: The simulation information that is not common between use cases can be added as new service models will be introduced. These models can be added to the [ran-simulator helm chart][RAN simulator helm chart] and can be loaded by RAN simulator. 

## MHO specific model
One example of a use case specific model is the **two-cell-two-node-model.yaml** model, is used by **onos-mho** xApplication to emulate UEs moving between pre-determined end-points. As opposed to the generic model that supports UEs moving on random routes (with randomly generated end-points), the two-cell-two-node-model is ideally suited to test handover scenarios deterministically. In order to support this model, RANSim supports the ability to specify the following directives in the model's yaml description:

* routeEndPoints: Start and end end-point coordinates for routes 
* directRoute: Direct route between end-points (as opposed to the default randomly zig-zagging route)
* initialRrcState: Specify the initial RRC state of UEs (as opposed to the default randomly assigned initial state)
* rrcStateChangesDisabled: Disable RRC state changes


[RAN simulator helm chart]: https://github.com/onosproject/sdran-helm-charts/tree/master/ran-simulator
