# Deploy with docker-compose

[docker-compose] is a simple tool that allow a collection of Docker images to be
launched together. It can be used as a simple alternative for deploying to a Kubernetes
cluster, or running the docker images individually.

> It is well suited to running in a demo scenario on a laptop

> It is only possible to run on docker-compose because none of the components below
> have a dependency on Atomix at this point in time (Feb '20). Atomix is not
> supported on docker-compose. When onos-ran changes to use Atomix, this page
> will be deleted.

In this project the [../build/docker-compose.yaml](../build/docker-compose.yaml)
file launches 6 applications:

1. onos-ran (the main ORAN application - exposes the C1 interface northbound, accesses the E2 interface on the southbound)
1. onos-ran-ho (the Handover application - attaches to the C1 interface of onos-ran)
1. onos-ran-mlb (the Load balancing application - attaches to the C1 interface of onos-ran)
1. ran-simulator (a service for simulating UEs and base stations - exposes the E2 interface northbound, and GUI services)
1. envoy-proxy (a grpc-web proxy that converts between ran-simulator gRPC GUI services and the GUI)
1. sd-ran-gui (a Google Maps based GUI that displays Towers, UE's, routes and links from ran-simulator)

## Google Maps API Key
The RAN Simulator can be connected to Google's [Directions API] if a valid Google
API Key is given. Google charges $5.00 per 1000 requests to the [Directions API],
and so we do not put our API key up in the public domain.

In addition the `sd-ran-gui` uses the Google [Maps API] to retrieve tiled maps,
at an additional cost of $7.00 per 1000 requests.
 
The key must be specified when running docker-compose as shown below.

> If the environment variable `GOOGLE_API_KEY` is not given, then the `ran-simulator`
> application will run, and will use randomly generated routes. The `sd-ran-gui` will
> exit and will not be accessible.

## Running
From the **Ran Simulator** directory run
```bash
GOOGLE_API_KEY=<YOUR_API_KEY_HERE> docker-compose -f build/docker-compose.yaml up
```

The first time this is run, it may pull the application images down from [Docker Hub].
On subsequent runs it will use images cached on your system.

To see the services running (in a separate terminal window) use:
```bash
> docker ps
CONTAINER ID        IMAGE                              COMMAND                  CREATED              STATUS              PORTS                                                         NAMES
12855e2344c2        onosproject/onos-ran:latest        "onos-ran-mlb -onosr…"   About a minute ago   Up 59 seconds                                                                     build_onos-ran-mlb_1
3ac40daa66b1        onosproject/onos-ran:latest        "onos-ran-ho -onosra…"   About a minute ago   Up 59 seconds                                                                     build_onos-ran-ho_1
36f985e634b2        onosproject/onos-ran:latest        "onos-ran -simulator…"   11 minutes ago       Up About a minute   0.0.0.0:25150->5150/tcp                                       build_onos-ran_1
8df77b786468        onosproject/sd-ran-gui:latest      "nginx -g 'daemon of…"   16 minutes ago       Up 58 seconds       0.0.0.0:4200->80/tcp                                          build_sd-ran-gui_1
7b514e8a21d0        onosproject/ran-simulator:latest   "trafficsim -googleA…"   16 minutes ago       Up 58 seconds       0.0.0.0:15150->5150/tcp                                       build_ran-simulator_1
f47337694054        envoyproxy/envoy-alpine:v1.11.1    "/docker-entrypoint.…"   5 days ago           Up 56 seconds       10000/tcp, 0.0.0.0:18080->8080/tcp, 0.0.0.0:19901->9901/tcp   build_envoy-proxy_1
```

Ensure 6 all are running. Check the logs shown in the startup terminal.

## onos-cli access
The `onos-cli` application has been extended to expose the C1 interface of `onos-ran`
as a set of commands.

To access the cli when running docker compose, run the `onos-cli` on the host machine.
First get the latest:
```bash
go get github.com/onosproject/onos-cli/cmd/onos
```
Then run it, using the forwarded port from `onos-ran` (from `docker-compose.yaml`)
like:
```bash
> go run github.com/onosproject/onos-cli/cmd/onos ran get stations --service-address localhost:25150 --no-tls
go: finding github.com/onosproject/onos-cli latest
ECID      MAX
0000003   5
0000002   5
0000009   5
0000008   5
0000005   5
0000006   5
0000007   5
0000001   5
0000004   5
```

### Tips for running

* See the `docker-compose` documentation for [more tips](https://docs.docker.com/compose/#step-8-experiment-with-some-other-commands)
* To run a subset of the containers, give only the names of the ones you want to start at the end of the command
    * e.g. ```docker-compose -f build/docker-compose.yaml up ran-simulator envoy-proxy```
* To run the containers in the background, use the `-d` flag after `up`
* To connect to the terminal shell of any container, use `docker exec -it <containername> /bin/sh`
* When you make updates to any one of the applications, run `make images` to push
the changes in to your local docker and `docker-compose` will pick it up the next
time it is started

## Browser access
With the `ran-simulator` and `envoy-proxy` services running, the GUI will be available on [http://localhost:4200]

[Directions API]: https://developers.google.com/maps/documentation/directions/start
[Maps API]: https://developers.google.com/maps/documentation/javascript/tutorial
[docker-compose]: https://docs.docker.com/compose/
[Docker Hub]: https://hub.docker.com/orgs/onosproject/repositories
