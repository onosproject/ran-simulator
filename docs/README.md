# Ran Simulator

The Ran Simulator is part of µONOS and is meant to work alongside `onos-ric` and
`onos-gui`.

The simulator mimics a collection of Cell Towers and a set of UE's moving along
routes between different locations.

The simulator has 2 main gRPC interfaces

1. **trafficsim** - for communicating with the `onos-gui` - this is exposed on port `5150`
1. **e2** - for communicating to `onos-ric` - there is a separate port opened per
cell site (tower) starting at port `5152-`

The number and position of towers are defined in a YAML file which is loaded
at startup - use the `-towerConfigName` param.

This can be changed as a value in the Helm chart at deploy time
e.g. `--set towerConfigName=berlin-honeycomb-169-6.yaml`

> After startup the number of towers cannot be changed.


As towers are created at startup, updates are made to:

1. **onos-topo** - a `device` is added to `onos-topo` per cell site (tower), with
an address including the cell-site port number.
1. **Kubernetes API** - the `ran-simulator` service is expanded after startup with
the port numbers added per cell site.

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
  -towerConfigName string
    	the name of a tower configuration (default "berlin-honeycomb-169-6.yaml")
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

## Creating the tower configuration files
The YAML files can be created by hand, or by application. There is an
application to create towers in a honeycomb (hexagonal) layout.

Sample configurations are supplied with the build and stored in `/etc/onos/config`

> These are copied from https://github.com/onosproject/ran-simulator/tree/master/pkg/config at build time.
>
> If a new layout is required that's not in the build, it can be mounted through
> the Helm chart with a "ConfigMap" that mounts at `/etc/onos/config`

To run the tool, first get it with:
```bash
go get github.com/onosproject/ran-simulator/cmd/honeycomb/honeycomb
```

and run it like:
```bash
go run github.com/onosproject/ran-simulator/cmd/honeycomb/honeycomb pkg/config/berlin-honeycomb-331-3.yaml \
     --towers 331 --sectors-per-tower 3 -a 52.52 -g 13.405 -i 0.03
```

## Browser access
When deployed with the **onos-gui** application, the simulation can be accessed
from a browser.

The [Map View](https://docs.onosproject.org/onos-gui/docs/ran-gui/#map-view) is linked directly to the `ran-simulator`

[Directions API]: https://developers.google.com/maps/documentation/directions/start
