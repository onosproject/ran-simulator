# Simulation Models

RAN simulator defines two levels of simulation models:

* **Generic model:** that defines E2 nodes, cells, service models,
E2T end points in a Yaml file. RAN simulator reads this file to create E2 nodes and initializing
data stores. A sample model is created using [Honeycomb Topology Generator](topology_generator.md) and added to the
  [ran-simulator helm chart][RAN simulator helm chart]. 


* **Use Case Specific Models**: The simulation information that are not
common between use cases can be added as new service models will be 
introduced. These models can be added to the [ran-simulator helm chart][RAN simulator helm chart]
and can be loaded by RAN simulator. 

  
## PCI Use Case Model
The metrics that are defined in the PCI model for each cell are listed as follows:

-  Cell Size: that can have one of the following values:
   * CellSize_CELL_SIZE_ENTERPRISE 
   * CellSize_CELL_SIZE_FEMTO
   * CellSize_CELL_SIZE_MACRO
   * CellSize_CELL_SIZE_OUTDOOR_SMALL
  
-  E-UTRA Absolute Radio Frequency Channel Number (EARFCN:)
-  Physical Cell ID (PCI).
-  PCI Pool: determines a list of PCI ranges that can be used for PCI value.



[RAN simulator helm chart]: https://github.com/onosproject/sdran-helm-charts/tree/master/ran-simulator