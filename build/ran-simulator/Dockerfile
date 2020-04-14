ARG GMAP_RAN_BASE_VERSION=latest

FROM onosproject/ran-simulator-base:$GMAP_RAN_BASE_VERSION as base

FROM alpine:3.11
RUN apk add libc6-compat

USER nobody

COPY --from=base /go/src/github.com/onosproject/ran-simulator/build/_output/trafficsim /usr/local/bin/trafficsim

ENTRYPOINT ["trafficsim"]
