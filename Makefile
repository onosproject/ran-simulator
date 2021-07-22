export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build

RAN_SIMULATOR_VERSION := latest
ONOS_PROTOC_VERSION := v0.6.7

OUTPUT_DIR=./build/_output

build: # @HELP build the Go binaries and run all validations (default)
build:
	go build ${BUILD_FLAGS} -o ${OUTPUT_DIR}/ransim ./cmd/ransim
	go build ${BUILD_FLAGS} -o ${OUTPUT_DIR}/honeycomb ./cmd/honeycomb

debug: BUILD_FLAGS += -gcflags=all="-N -l"
debug: build # @HELP build the Go binaries with debug symbols

test: # @HELP run the unit tests and source code validation producing a golang style report
test: build deps linters license_check
	go test -race github.com/onosproject/ran-simulator/...

jenkins-test:  # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: build-tools build deps license_check linters
	TEST_PACKAGES=github.com/onosproject/ran-simulator/pkg/... ./../build-tools/build/jenkins/make-unit

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	go test -covermode=count -coverprofile=onos.coverprofile github.com/onosproject/ran-simulator/pkg/...
	cd .. && go get github.com/mattn/goveralls && cd ran-simulator
	grep -v .pb.go onos.coverprofile >onos-nogrpc.coverprofile
	goveralls -coverprofile=onos-nogrpc.coverprofile -service travis-pro -repotoken xHYC7gvqJdxaScSObicSox1E6sraczouC

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: golang-ci # @HELP examines Go source code and reports coding problems
	golangci-lint run --timeout 5m

build-tools: # @HELP install the ONOS build tools if needed
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi

jenkins-tools: # @HELP installs tooling needed for Jenkins
	cd .. && go get -u github.com/jstemmer/go-junit-report && go get github.com/t-yuki/gocover-cobertura

golang-ci: # @HELP install golang-ci if not present
	golangci-lint --version || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b `go env GOPATH`/bin v1.36.0

license_check: build-tools # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}/pkg --boilerplate LicenseRef-ONF-Member-1.0

gofmt: # @HELP run the Go format validation
	bash -c "diff -u <(echo -n) <(gofmt -d pkg/ cmd/ tests/)"

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
	./../build-tools/publish-version ${VERSION} onosproject/ran-simulator

jenkins-publish: build-tools jenkins-tools # @HELP Jenkins calls this to publish artifacts
	./build/bin/push-images
	../build-tools/release-merge-commit

bumponosdeps: # @HELP update "onosproject" go dependencies and push patch to git.
	./../build-tools/bump-onos-deps ${VERSION}

clean: # @HELP remove all the build artifacts
	rm -rf ${OUTPUT_DIR} ./cmd/trafficsim/trafficsim ./cmd/ransim/ransim
	go clean -testcache github.com/onosproject/ran-simulator/...

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
