# SPDX-License-Identifier: Apache-2.0
# Copyright 2019 Open Networking Foundation
# Copyright 2024 Intel Corporation

export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build

RAN_SIMULATOR_VERSION ?= latest
ONOS_PROTOC_VERSION := v0.6.9

OUTPUT_DIR=./build/_output

GOLANG_CI_VERSION := v1.52.2

all: build docker-build

build: # @HELP build the Go binaries and run all validations (default)
build:
	go build ${BUILD_FLAGS} -o ${OUTPUT_DIR}/ransim ./cmd/ransim
	go build ${BUILD_FLAGS} -o ${OUTPUT_DIR}/honeycomb ./cmd/honeycomb

test: # @HELP run the unit tests and source code validation producing a golang style report
test: build lint license
	go test -race github.com/onosproject/ran-simulator/...

docker-build-ran-simulator: # @HELP build ran-simulator Docker image
	@go mod vendor
	docker build . -f build/ran-simulator/Dockerfile \
		-t onosproject/ran-simulator:${RAN_SIMULATOR_VERSION}

docker-build: # @HELP build all Docker images
docker-build: build docker-build-ran-simulator

docker-push-ran-simulator: # @HELP push onos-pci Docker image
	docker push onosproject/ran-simulator:${RAN_SIMULATOR_VERSION}

docker-push: # @HELP push docker images
docker-push: docker-push-ran-simulator

lint: # @HELP examines Go source code and reports coding problems
	golangci-lint --version | grep $(GOLANG_CI_VERSION) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b `go env GOPATH`/bin $(GOLANG_CI_VERSION)
	golangci-lint run --timeout 15m

license: # @HELP run license checks
	rm -rf venv
	python3 -m venv venv
	. ./venv/bin/activate;\
	python3 -m pip install --upgrade pip;\
	python3 -m pip install reuse;\
	reuse lint

check-version: # @HELP check version is duplicated
	./build/bin/version_check.sh all

clean:: # @HELP remove all the build artifacts
	rm -rf ${OUTPUT_DIR} ./cmd/trafficsim/trafficsim ./cmd/ransim/ransim
	go clean github.com/onosproject/ran-simulator/...

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '

model-files: # @HELP generate various model and model-topo YAML files in sdran-helm-charts/ran-simulator
	go run cmd/honeycomb/honeycomb.go topo --plmnid 314628 --towers 2  --ue-count 10 --controller-yaml ../sdran-helm-charts/ran-simulator/files/topo/model-topo.yaml ../sdran-helm-charts/ran-simulator/files/model/model.yaml
	go run cmd/honeycomb/honeycomb.go topo --plmnid 314628 --towers 12 --ue-count 100 --sectors-per-tower 6 --controller-yaml ../sdran-helm-charts/ran-simulator/files/topo/scale-model-topo.yaml ../sdran-helm-charts/ran-simulator/files/model/scale-model.yaml
	go run cmd/honeycomb/honeycomb.go topo --plmnid 314628 --towers 1 --ue-count 5 --controller-yaml ../sdran-helm-charts/ran-simulator/files/topo/three-cell-model-topo.yaml ../sdran-helm-charts/ran-simulator/files/model/three-cell-model.yaml