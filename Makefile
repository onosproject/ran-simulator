export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build

RAN_SIMULATOR_VERSION := latest
ONOS_BUILD_VERSION := v0.6.7
ONOS_PROTOC_VERSION := v0.6.7

build: # @HELP build the Go binaries and run all validations (default)
build:
	export GOPRIVATE="github.com/onosproject/*"
	go build -o build/_output/ransim ./cmd/ransim

test: # @HELP run the unit tests and source code validation producing a golang style report
test: build deps linters license_check
	go test -race github.com/onosproject/ran-simulator/...

jenkins-test: build-tools # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: build deps license_check linters
	export GOPRIVATE="github.com/onosproject/*"
	TEST_PACKAGES=github.com/onosproject/ran-simulator/pkg/... ./../build-tools/build/jenkins/make-unit

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	export GOPRIVATE="github.com/onosproject/*"
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

golang-ci: # @HELP install golang-ci if not present
	golangci-lint --version || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b `go env GOPATH`/bin v1.36.0

license_check: build-tools # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}/pkg --boilerplate LicenseRef-ONF-Member-1.0

gofmt: # @HELP run the Go format validation
	bash -c "diff -u <(echo -n) <(gofmt -d pkg/ cmd/ tests/)"

protos: # @HELP compile the protobuf files (using protoc-go Docker)
	docker run -it -v `pwd`:/go/src/github.com/onosproject/ran-simulator \
		-v `pwd`/../build-tools/licensing:/build-tools/licensing \
		-w /go/src/github.com/onosproject/ran-simulator \
		--entrypoint build/bin/compile-protos.sh \
		onosproject/protoc-go:${ONOS_PROTOC_VERSION}

ran-simulator-docker: # @HELP build ran-simulator Docker image
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

bumponosdeps: # @HELP update "onosproject" go dependencies and push patch to git.
	./../build-tools/bump-onos-deps ${VERSION}

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./cmd/trafficsim/trafficsim ./cmd/ransim/ransim
	go clean -testcache github.com/onosproject/ran-simulator/...

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
