# gmap-ran
## trafficsim
The main application is the **Traffic Simulator** which generates UE locations randomly
around a set of cell towers. This application is written in Go, and exposes a gRPC
interface on port 5150.

## sd-ran-gui
There is a Google Map based GUI then that displays these UEs (cars) on a map as they move.
It connects to the Traffic Simulator over the gRPC port (proxied through an Envoy
Web Proxy with does a grpc-web translation on the messages to convert them to
HTTP 1.1). The services available to it are

* List Towers
* List Ues
* List Routes

This is an [Angular] appliation and written in [TypeScript] in `web/sd-ran-gui`.
See the [README.md](web/sd-ran-gui/README.md) in that folder for more info.

## Building the application
```bash
make all
```

Prerequisites:
```bash
npm i grpc-web
npm i google-protobuf
```

## Running the application
To run the App use docker compose
```bash
docker-compose -f build/docker-compose.yaml up
```

and open your browser at localhost:4200

## Next stage
The Traffic Simulator will relay information about the UE's distance (and signal power)
back to the Ran simulator.

[Angular]: https://angular.io/
[TypeScript]: https://www.typescriptlang.org/
