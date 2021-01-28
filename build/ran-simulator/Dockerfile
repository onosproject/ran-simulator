ARG GMAP_RAN_BASE_VERSION=latest

FROM onosproject/ran-simulator-base:$GMAP_RAN_BASE_VERSION as base

FROM frolvlad/alpine-glibc:alpine-3.12

USER nobody

COPY --from=base /go/src/github.com/onosproject/onos-e2t/build/_output/ransim /usr/local/bin/ransim
# COPY pkg/config/*.yaml /etc/onos/config/

ENTRYPOINT ["ransim"]
