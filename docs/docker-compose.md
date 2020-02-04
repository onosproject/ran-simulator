# Deploy with docker-compose

`docker-compose` is a simple tool that allow a collection of Docker images to be
launched together. It can be used as a simple alternative for deploying to a Kubernetes
cluster, or running the docker images individually.

> It is well suited to running in a demo scenario on a laptop

> It is only possible to run on docker-compose because none of the components below
> have a dependency on Atomix at this point in time (Feb '20). Atomix is not
> supported on docker-compose. When onos-ran changes to use Atomix, this page
> will be deleted. 

In this project the [ran-simulator/build/dockercompose.yaml](ran-simulator/build/dockercompose.yaml)
file launches 4 applications:

1. onos-ran (the main ORAN application - exposes the C1 interface northbound, accesses the E2 interface on the southbound )
1. ran-simulator (a service for simulating UEs and base stations - exposes the E2 interface northbound, and GUI services)
1. envoy-proxy (a grpc-web proxy that converts between ran-simulator gRPC GUI services and the GUI)
1. sd-ran-gui (a Google Maps based GUI that displays Towers, UE's, routes and links from ran-simulator)

## Google Maps API Key
The RAN Simulator connects to Google's [Directions API] and so needs a Google API Key.
Google charges $5.00 per 1000 requests to the [Directions API], and so we do not put
our API key up in the public domain.

You **must** enter your own key in to `build/docker-compose.yaml` before you
start `docker-compose`, or else the `ran-simulator` service will fail to start.

## Running
From the Ran Simulator directory run
```bash
docker-compose -f build/docker-compose.yaml up
```

The first time this is run, it may pull the application images down from the internet.
On subsequent runs it will use images cached on your system.

To see the services running (in a separate terminal window) use:
```bash
> docker ps
CONTAINER ID        IMAGE                              COMMAND                  CREATED             STATUS              PORTS                                                         NAMES
4838449725ec        onosproject/ran-simulator:latest   "trafficsim -googleA…"   25 seconds ago      Up 19 seconds       0.0.0.0:15150->5150/tcp                                       build_ran-simulator_1
c96f93e5fb71        onosproject/onos-ran:latest        "onos-ran -certPath=…"   11 minutes ago      Up 20 seconds       0.0.0.0:25150->5150/tcp                                       build_onos-ran_1
06a026a4f556        onosproject/sd-ran-gui:latest      "nginx -g 'daemon of…"   20 hours ago        Up 21 seconds       0.0.0.0:4200->80/tcp                                          build_sd-ran-gui_1
a298d1214eb1        envoyproxy/envoy-alpine:v1.11.1    "/docker-entrypoint.…"   3 days ago          Up 21 seconds       10000/tcp, 0.0.0.0:18080->8080/tcp, 0.0.0.0:19901->9901/tcp   build_envoy-proxy_1
```

Ensure 4 all are running. If `ran-simulator` is not running, you might not have
entered a valid Google API Key. Check the logs shown in the startup terminal.

## Browser access
With the 4 services running, the GUI will be available on [http://localhost:4200]

[Directions API]: https://developers.google.com/maps/documentation/directions/start
