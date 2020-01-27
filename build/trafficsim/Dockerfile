ARG GMAP_RAN_BASE_VERSION=latest

FROM onosproject/trafficsim-base:$GMAP_RAN_BASE_VERSION as base

FROM alpine:3.9
RUN apk add libc6-compat

USER nobody

COPY --from=base /go/src/github.com/OpenNetworkingFoundation/gmap-ran/build/_output/trafficsim /usr/local/bin/trafficsim

ENTRYPOINT ["trafficsim"]
