export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

RAN_SIMULATOR_VERSION := latest
ONOS_BUILD_VERSION := v0.6.3
ONOS_PROTOC_VERSION := v0.6.3

build: # @HELP build the Go binaries and run all validations (default)
build:
	export GOPRIVATE="github.com/onosproject/*"
	CGO_ENABLED=1 go build -o build/_output/trafficsim ./cmd/trafficsim

test: # @HELP run the unit tests and source code validation
test: build deps linters license_check
	CGO_ENABLED=1 go test -race github.com/onosproject/ran-simulator/pkg/...
	CGO_ENABLED=1 go test -race github.com/onosproject/ran-simulator/cmd/...
	CGO_ENABLED=1 go test -race github.com/onosproject/ran-simulator/api/...

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	./../build-tools/build/coveralls/coveralls-coverage ran-simulator xHYC7gvqJdxaScSObicSox1E6sraczouC

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run --timeout 30m

license_check: # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR} --boilerplate LicenseRef-ONF-Member-1.0

gofmt: # @HELP run the Go format validation
	bash -c "diff -u <(echo -n) <(gofmt -d pkg/ cmd/ tests/)"

protos: # @HELP compile the protobuf files (using protoc-go Docker)
	docker run -it -v `pwd`:/go/src/github.com/onosproject/ran-simulator \
		-v `pwd`/../build-tools/licensing:/build-tools/licensing \
		-w /go/src/github.com/onosproject/ran-simulator \
		--entrypoint build/bin/compile-protos.sh \
		onosproject/protoc-go:${ONOS_PROTOC_VERSION}

ran-simulator-base-docker: # @HELP build ran-simulator base Docker image
	@go mod vendor
	docker build . -f build/base/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		--build-arg ONOS_MAKE_TARGET=build \
		-t onosproject/ran-simulator-base:${RAN_SIMULATOR_VERSION}
	@rm -rf vendor

ran-simulator-docker: ran-simulator-base-docker # @HELP build ran-simulator Docker image
	docker build . -f build/ran-simulator/Dockerfile \
		--build-arg GMAP_RAN_BASE_VERSION=${RAN_SIMULATOR_VERSION} \
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

bumponosdeps: # @HELP update "onosproject" go dependencies and push patch to git.
	./../build-tools/bump-onos-deps ${VERSION}

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/trafficsim/trafficsim
	go clean -testcache github.com/onosproject/ran-simulator/...

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
