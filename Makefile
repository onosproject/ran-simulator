# SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build

RAN_SIMULATOR_VERSION := latest
ONOS_PROTOC_VERSION := v0.6.9

OUTPUT_DIR=./build/_output

build: # @HELP build the Go binaries and run all validations (default)
build:
	go build ${BUILD_FLAGS} -o ${OUTPUT_DIR}/ransim ./cmd/ransim
	go build ${BUILD_FLAGS} -o ${OUTPUT_DIR}/honeycomb ./cmd/honeycomb

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

debug: BUILD_FLAGS += -gcflags=all="-N -l"
debug: build # @HELP build the Go binaries with debug symbols

test: # @HELP run the unit tests and source code validation producing a golang style report
test: build deps linters license
	go test -race github.com/onosproject/ran-simulator/...

jenkins-test:  # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: build deps license linters
	TEST_PACKAGES=github.com/onosproject/ran-simulator/pkg/... ./build/build-tools/build/jenkins/make-unit

integration-tests: # @HELP run helmit integration tests
	@kubectl delete ns test; kubectl create ns test
	helmit test -n test ./cmd/ransim-tests --timeout 30m --no-teardown \
		--secret sd-ran-username=${repo_user} --secret sd-ran-password=${repo_password}

model-files: # @HELP generate various model and model-topo YAML files in sdran-helm-charts/ran-simulator
	go run cmd/honeycomb/honeycomb.go topo --plmnid 314628 --towers 2  --ue-count 10 --controller-yaml ../sdran-helm-charts/ran-simulator/files/topo/model-topo.yaml ../sdran-helm-charts/ran-simulator/files/model/model.yaml
	go run cmd/honeycomb/honeycomb.go topo --plmnid 314628 --towers 12 --ue-count 100 --sectors-per-tower 6 --controller-yaml ../sdran-helm-charts/ran-simulator/files/topo/scale-model-topo.yaml ../sdran-helm-charts/ran-simulator/files/model/scale-model.yaml
	go run cmd/honeycomb/honeycomb.go topo --plmnid 314628 --towers 1 --ue-count 5 --controller-yaml ../sdran-helm-charts/ran-simulator/files/topo/three-cell-model-topo.yaml ../sdran-helm-charts/ran-simulator/files/model/three-cell-model.yaml

ran-simulator-docker: # @HELP build ran-simulator Docker image
	@go mod vendor
	docker build . -f build/ran-simulator/Dockerfile \
		-t onosproject/ran-simulator:${RAN_SIMULATOR_VERSION}

images: # @HELP build all Docker images
images: ran-simulator-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/ran-simulator:${RAN_SIMULATOR_VERSION}

all: build images

publish: # @HELP publish version on github and dockerhub
	./build/build-tools/publish-version ${VERSION} onosproject/ran-simulator

jenkins-publish: build-tools jenkins-tools # @HELP Jenkins calls this to publish artifacts
	./build/bin/push-images
	./build/build-tools/release-merge-commit

clean:: # @HELP remove all the build artifacts
	rm -rf ${OUTPUT_DIR} ./cmd/trafficsim/trafficsim ./cmd/ransim/ransim
	go clean -testcache github.com/onosproject/ran-simulator/...

