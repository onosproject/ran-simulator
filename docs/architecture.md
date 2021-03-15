# RAN simulator Architecture

The RAN simulator application has the following key components as the following figure shows: 

* **E2 nodes**: Upper half of the RAN simulator is responsible to simulate e2 nodes where each E2 node implements an E2 agent using E2AP, and implement service models.

* **RAN Environment**: Lower half of the RAN simulator is  responsible to simulate RAN environment to support required RAN functions
  for implementing E2 service models (e.g. simulating User Equipments (UEs), mobility, etc).

* **Data Stores**: lower half and upper half are connected using data stores that stores information
about E2-nodes, E2-agents, UEs, RAN metrics, E2 subscriptions, etc. 

* **RAN simulator APIs**: RAN simulator provides a variety of gRPC APIs that can be used for controlling E2 nodes and RAN environment. 
You can find more details about RAN simulator APIs here: [RAN simulator APIs](api.md)
  
* **RAN simulator CLI**: RAN simulator is equipped with a command line interface which is integrated with 
 [onos-cli](https://github.com/onosproject/onos-cli) that allows to interact with RAN simulator to retrieve required information from data stores, 
  monitor RAN environment changes, create/remove/update RAN entities, metrics, etc.
The list of ransim commands is documented here [RAN simulator CLI](https://github.com/onosproject/onos-cli/blob/master/docs/cli/onos_ransim.md)  
  

![RAN simulator Architecture](images/ransim_architecture.jpg)