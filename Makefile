export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

GMAP_RAN_VERSION := latest
GMAP_RAN_DEBUG_VERSION := debug
ONOS_BUILD_VERSION := stable

build: # @HELP build the Go binaries and run all validations (default)
build:
	CGO_ENABLED=1 go build -o build/_output/trafficsim ./cmd/trafficsim

build-gui:
	cd web/sd-ran-gui && ng build --prod

test: # @HELP run the unit tests and source code validation
test: build deps linters license_check
	#CGO_ENABLED=1 go test -race github.com/onosproject/ran-simulator/pkg/...
	#CGO_ENABLED=1 go test -race github.com/OpenNetworkingFoundation/trafficsim/cmd/...
	#CGO_ENABLED=1 go test -race github.com/OpenNetworkingFoundation/trafficsim/api/...

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	./build/bin/coveralls-coverage

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run --timeout 30m

license_check: # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}

gofmt: # @HELP run the Go format validation
	bash -c "diff -u <(echo -n) <(gofmt -d pkg/ cmd/ tests/)"

protos: # @HELP compile the protobuf files (using protoc-go Docker)
	docker run -it -v `pwd`:/go/src/github.com/OpenNetworkingFoundation/gmap-ran \
		-w /go/src/github.com/OpenNetworkingFoundation/gmap-ran \
		--entrypoint build/bin/compile-protos.sh \
		onosproject/protoc-go:stable

trafficsim-base-docker: # @HELP build trafficsim base Docker image
	@go mod vendor
	docker build . -f build/base/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		--build-arg ONOS_MAKE_TARGET=build \
		-t onosproject/trafficsim-base:${GMAP_RAN_VERSION}
	@rm -rf vendor

trafficsim-docker: trafficsim-base-docker # @HELP build trafficsim Docker image
	docker build . -f build/trafficsim/Dockerfile \
		--build-arg GMAP_RAN_BASE_VERSION=${GMAP_RAN_VERSION} \
		-t onosproject/trafficsim:${GMAP_RAN_VERSION}

images: # @HELP build all Docker images
images: trafficsim-docker build-gui

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/trafficsim:${GMAP_RAN_VERSION}

all: build images

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/trafficsim/trafficsim ./cmd/onos/onos
	go clean -testcache github.com/OpenNetworkingFoundation/trafficsim/...

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
