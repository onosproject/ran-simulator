ARG ONOS_BUILD_VERSION=undefined

FROM onosproject/golang-build:v0.6.7 as build

ENV GO111MODULE=on
ARG ONOS_MAKE_TARGET=build

COPY Makefile go.mod go.sum /build/
COPY cmd/ /build/cmd/
COPY pkg/ /build/pkg/
COPY vendor/ /build/vendor/

WORKDIR /build

RUN GOFLAGS=-mod=vendor make ${ONOS_MAKE_TARGET}


FROM frolvlad/alpine-glibc:alpine-3.12

USER nobody

COPY --from=build /build/build/_output/ransim /usr/local/bin/ransim
# COPY pkg/config/*.yaml /etc/onos/config/

ENTRYPOINT ["ransim"]
