# RAN Simulator
## ran-simulator
The main application is the **Ran Simulator** which generates UE locations randomly
around a set of cell towers. This application is written in Go, and exposes a gRPC
interface on port 5150 with
* E2 interface
* GUI interface

## sd-ran-gui
There is a Google Map based GUI then that displays these UEs (cars) on a map as they move.
It connects to the Traffic Simulator over the gRPC port (proxied through an [Envoy Proxy]
with does a [grpc-web] translation on the messages to convert them to
HTTP 1.1). The services available to it are

This is an [Angular] appliation and written in [TypeScript] in `web/sd-ran-gui`.
See the [README.md](web/sd-ran-gui/README.md) in that folder for more info.

## Running the applications
See [demo.md](docs/demo.md) for details of running the application.
It can be run simply with [docker-compose](docs/docker-compose.md)

Alternatively it can be run on Kubernetes - see [ran-simulator.md](docs/deployment.md) and
[sd-ran-gui-deployment.md](docs/sd-ran-gui-deployment.md)

## Building the application
> You do not need to build the applications to run them, as they have already been
>built and loaded in to Docker hub 

Follow the prerequisites for both:

* [Go Projects] at [https://docs.onosproject.org/onos-docs/docs/content/developers/prerequisites/](https://docs.onosproject.org/onos-docs/docs/content/developers/prerequisites/)
* [GUI Projects] at [https://docs.onosproject.org/onos-gui/docs/prerequisites/](https://docs.onosproject.org/onos-gui/docs/prerequisites/)

Add your Google maps API Key in to:
 
* `gmap-ran/web/sd-ran-gui/src/index.html`
AND
* `gmap-ran/build/docker-compose.yaml`
 
> Do not persist your API Key in GitHub, as it will probably be "borrowed" by others,
>and will incur charges from Google.
 
 and build with:
```bash
make images
```


[Angular]: https://angular.io/
[TypeScript]: https://www.typescriptlang.org/
[docker-compose]: https://docs.docker.com/compose/
[grpc-web]: https://github.com/grpc/grpc-web
[Envoy Proxy]: https://www.envoyproxy.io/
