# Ran Simulator

The Ran Simulator is part of µONOS and is meant to work alongside `onos-ric` and
`sd-ran-gui`.

The simulator mimics a collection of Cell Towers and a set of UE's moving along
routes between different locations.

The application is very tunable through startup parameters, and can be run:
 
* as a standalone application in docker
* deployed in a **Kubernetes** cluster

## Google Maps API Key
The RAN Simulator can connect to Google's [Directions API] with a Google API Key.
Google charges $5.00 per 1000 requests to the [Directions API], and so we do not put
our API key up in the public domain. Without the API key directions will be created
randomly between 2 locations (usually in the form of a zig-zag line).

The SD RAN GUI also (separately) accesses Google [Maps API] and incurs an additional
cost of $7.00 per 1000 requests.

## Startup parameters
Supplying `ran-simulator` with a bogus parameter gets it to show the start parameters
and their defaults
```bash
docker run -it onosproject/ran-simulator:latest -test
...
Usage of trafficsim:
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
  -mapCenterLat float
    	Map center latitude (default 52.52)
  -mapCenterLng float
    	Map center longitude (default 13.405)
  -maxUEsPerTower
        Max num of UEs per tower
  -numLocations int
    	Number of locations (default 10)
  -numRoutes int
    	Number of routes (default 3)
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
    	The starting Zoom level (default 12)
```

See [deployment.md](deployment.md) for how to change these for a Kubernetes deployment.

## Browser access
When deployed with the **sd-ran-gui** application, the simulation can be accessed
from a browser.

[Directions API]: https://developers.google.com/maps/documentation/directions/start
