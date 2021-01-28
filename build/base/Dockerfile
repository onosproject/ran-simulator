ARG ONOS_BUILD_VERSION=undefined

FROM onosproject/golang-build:$ONOS_BUILD_VERSION

ENV GO111MODULE=on
ARG ONOS_MAKE_TARGET=build

COPY Makefile go.mod go.sum /go/src/github.com/onosproject/onos-e2t/
COPY api/ /go/src/github.com/onosproject/onos-e2t/api/
COPY cmd/ /go/src/github.com/onosproject/onos-e2t/cmd/
COPY pkg/ /go/src/github.com/onosproject/onos-e2t/pkg/
COPY vendor/ /go/src/github.com/onosproject/onos-e2t/vendor/

RUN cd /go/src/github.com/onosproject/onos-e2t && GOFLAGS=-mod=vendor make ${ONOS_MAKE_TARGET}
