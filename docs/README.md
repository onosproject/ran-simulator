# Ran Simulator

The Ran Simulator is part of µONOS and is meant to work alongside `onos-ric` and
`onos-gui`, `onos-topo` and `onos-config`.

The simulator mimics a collection of Cell Towers and a set of UE's moving along
routes between different locations.

The simulator has 2 main gRPC interfaces:

1. **trafficsim** - for communicating with the `onos-gui` - this is exposed on port `5150`
1. **e2** & **gnmi** - for communicating to `onos-ric` and `onos-config` - there
is a separate port opened per cell site (tower) usually starting at port `5152-`

The number and position of towers is got from `onos-topo`. By default no towers
are present in `onos-topo` when it is started from new. When a device (with type
**E2Node** and version **1.0.0**) is created in `onos-topo`, the `ran-simulator`
detects it and will:

1. Create a gRPC server exposing E2 and gNMI interfaces for the specified **gnmiport**.
1. call on the **Kubernetes API** to open the **gnmiport**.

> For the Kubernetes service to be manipulated like this, the 'update' RBAC permission
> has to be given to the `ran-simulator` pod. This is done through the Helm Chart
> as the role `ran-simulator-service-role` and the rolebinding `ran-simulator-access-services`

The application is very tunable through startup parameters. It can be deployed
only in a **Kubernetes** cluster. See [deployment](./deployment.md).

### Google Maps API Key
The RAN Simulator can connect to Google's [Directions API] with a Google API Key.
Google charges $5.00 per 1000 requests to the [Directions API], and so we do not put
our API key up in the public domain.  

> This feature will create Routes that follow known road and street layouts in the
> real world and will make a better demo when used with `onos-gui`

Without the API key Routes will be created randomly between 2 locations
(usually in the form of a zig-zag line).

## Startup parameters
Run `ran-simulator` with `-help` parameter to show the start parameters
and their defaults
```bash
docker run -it onosproject/ran-simulator:latest -help
...
Usage of trafficsim:
  -addK8sSvcPorts
    	Add K8S service ports per tower (default true)
  -caPath string
    	path to CA certificate
  -certPath string
    	path to client certificate
  -fade
    	Show map as faded on start (default true)
  -googleAPIKey string
    	your google maps api key
  -keyPath string
    	path to client private key
  -locationsScale float
    	Ratio of random locations diameter to tower grid width (default 1.25)
  -maxUEs uint
    	Max number of UEs for complete simulation (default 300)
  -metricsAllHoEvents
    	Export all HO events in metrics (only historgram if false) (default true)
  -metricsPort uint
    	port for Prometheus metrics (default 9090)
  -minUEs uint
    	Max number of UEs for complete simulation (default 3)
  -showPower
    	Show power as circle on start (default true)
  -showRoutes
    	Show routes on start (default true)
  -stepDelayMs uint
    	delay between steps on route (default 1000)
  -topoEndpoint string
    	Endpoint for the onos-topo service (default "onos-topo:5150")
  -zoom float
    	The starting Zoom level (default 13)
```

> Some of these only have an effect when the onos-gui MapView is active
>
> e.g. `fade` is used to control whether the map is displayed with full opacity or
>faded at startup
>
> `-showRoutes`, `-showPower` and `-zoom` are also only related to the display

See [deployment.md](deployment.md) for how to change these for a Kubernetes deployment.

## Creating the tower/cell configuration files
The YAML files can be created by hand, or by application. There is an
application to create towers with sectors in a honeycomb (hexagonal) layout.

Sample outputs from this tool are in
[ran-simulator/pkg/config](https://github.com/onosproject/ran-simulator/tree/master/pkg/config).

> The sample configurations can be copied over to the `onos-cli` pod and run from there
> with the `kubectl cp` command.

There are 2 types of file:

1. *-topo.yaml files - these can be loaded in to `onos-topo` using the `onos-cli`
command
    1. like `onos topo load yaml <filename>-topo.yaml`
1. *-gnmi.yaml files - these can be loaded in to `onos-config` using the `onos-cli`
command
    1. like `onos config load yaml <filename>-gnmi.yaml`

To run the **honeycomb** generation tool, first get it with:
```bash
go get github.com/onosproject/ran-simulator/cmd/honeycomb
```

and run it (for topo) like:
```bash
go run github.com/onosproject/ran-simulator/cmd/honeycomb topo pkg/config/berlin-honeycomb-331-3-topo.yaml \
     --towers 331 --sectors-per-tower 3 -a 52.52 -g 13.405 -i 0.03
```

### Adding devices individually
Individual cells can be added to `onos-topo` using the `onos topo add device` but
be aware that the Type must be "E2Node", the version "1.0.0" and the 6 attributes must
be added "plmnid", "ecid", "longitude", "latitude", "azimuth" & "arc", or else the
topo device will be ignored.

```
onos topo add device 315010-0001234 -a ran-simulator:4660 -t E2Node -v 1.0.0 \
--insecure -d "New Tower" --attributes ecid=0001234 --attributes plmnid=315010 --attributes azimuth=0 \
--attributes arc=120 --attributes grpcport=4660 --attributes latitude=52.468038 --attributes longitude=13.355697
```

## gNMI access
Each Cell supports configuration through a [gNMI] interface, according to the YANG
model [E2Node](https://github.com/onosproject/config-models/tree/master/modelplugin/e2node-1.0.0/yang).

Usually we connect to `onos-config` and allow it to propagate the config changes
through to each Cell. `onos-config` is aware of the existence of the gNMI interface
on the Cell, because it is listed in `onos-topo`.

It is also possible to connect directly to the gNMI interface on the Cell e.g. at `ran-simulator:5162`.

For example to do a gNMI Get from inside the `onos-cli` pod - run:
```bash
gnmi_cli -get -address ran-simulator:5155 -proto "prefix: <>" -timeout 5s -en PROTO -alsologtostderr -insecure -client_crt /etc/ssl/certs/client1.crt -client_key /etc/ssl/certs/clie
nt1.key -ca_crt /etc/ssl/certs/onfca.crt
```

To do a set
```bash
gnmi_cli -set -address ran-simulator:5155 \
-proto "prefix: <elem: <name: 'e2node'> elem: <name: 'intervals'>> update: < path: <elem: <name: 'RadioMeasReportPerUe'>> val: <uint_val: 21>> update: < path: <elem: <name: 'SchedMeasReportPerUe'>> val: <uint_val: 22>>" \
-timeout 5s -en PROTO -alsologtostderr -insecure \
-client_crt /etc/ssl/certs/client1.crt -client_key /etc/ssl/certs/client1.key -ca_crt /etc/ssl/certs/onfca.crt
```

See [onos-config](https://docs.onosproject.org/onos-config/docs/gnmi/) for more
details on how to install and use the `gnmi_cli` tool.

## Browser access
When deployed with the **onos-gui** application, the simulation can be accessed
from a browser.

The [Map View](https://docs.onosproject.org/onos-gui/docs/ran-gui/#map-view) is linked directly to the `ran-simulator`

[Directions API]: https://developers.google.com/maps/documentation/directions/start
[gNMI]: https://datatracker.ietf.org/meeting/98/materials/slides-98-rtgwg-gnmi-intro-draft-openconfig-rtgwg-gnmi-spec-00