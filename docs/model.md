# Simulation Models

RAN simulator defines two levels of simulation models:

* **Generic model:** that defines E2 nodes, cells, service models, E2T end points in a Yaml file. RAN simulator reads this file to create E2 nodes and initializing data stores. A sample [model](https://github.com/onosproject/sdran-helm-charts/blob/master/ran-simulator/files/model/model.yaml) has been created using the [Honeycomb Topology Generator](topology_generator.md) and added to the [ran-simulator helm chart][RAN simulator helm chart]. 


* **Use Case Specific Models**: The simulation information that is not common between use cases can be added as new service models will be introduced. These models can be added to the [ran-simulator helm chart][RAN simulator helm chart] and can be loaded by RAN simulator. 


[RAN simulator helm chart]: https://github.com/onosproject/sdran-helm-charts/tree/master/ran-simulator