# Ran Simulator

The Ran Simulator is part of µONOS and is meant to work alongside `onos-ric` and
`onos-gui`.

The simulator mimics a collection of Cell Towers and a set of UE's moving along
routes between different locations.

The simulator has 2 main gRPC interfaces

1. **trafficsim** - for communicating with the `onos-gui` - this is exposed on port `5150`
1. **e2** - for communicating to `onos-ric` - there is a separate port opened per
cell site (tower) starting at port `5152-`

> The number of towers can be chosen at startup with the `-towerCols` and the `-towerRows`
> arguments shown below. After startup the number of towers cannot be changed.

As towers are created at startup, updates are made to:

1. **onos-topo** - a `device` is added to `onos-topo` per cell site (tower), with
an address including the cell-site port number.
1. **Kubernetes API** - the `ran-simulator` service is expanded after startup with
the port numbers added per cell site.

> For the Kubernetes service to be manipulated like this, the 'update' RBAC permission
> has to be given to the `ran-simulator` pod. This is done through the Helm Chart
> as the role `ran-simulator-service-role` and the rolebinding `ran-simulator-access-services`

The application is very tunable through startup parameters, and can be deployed
only in a **Kubernetes** cluster. See [deployment](./deployment.md).

## Google Maps API Key
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
Usage of /tmp/go-build089760012/b001/exe/trafficsim:
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
        Ratio of random locations diameter to tower grid width (default 1)
  -mapCenterLat float
        Map center latitude (default 52.52)
  -mapCenterLng float
        Map center longitude (default 13.405)
  -maxUEs int
        Max number of UEs for complete simulation (default 300)
  -maxUEsPerTower int
        Max num of UEs per tower (default 5)
  -metricsPort int
        port for Prometheus metrics (default 9090)
  -minUEs int
        Max number of UEs for complete simulation (default 3)
  -showPower
        Show power as circle on start (default true)
  -showRoutes
        Show routes on start (default true)
  -stepDelayMs int
        delay between steps on route (default 1000)
  -towerCols int
        Number of columns of towers (default 3)
  -towerRows int
        Number of rows of towers (default 3)
  -towerSpacingHoriz float
        Tower spacing horiz in degrees longitude (default 0.02)
  -towerSpacingVert float
        Tower spacing vert in degrees latitude (default 0.02)
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

## Browser access
When deployed with the **onos-gui** application, the simulation can be accessed
from a browser.

[Directions API]: https://developers.google.com/maps/documentation/directions/start
