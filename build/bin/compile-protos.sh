#!/bin/sh

proto_imports=".:${GOPATH}/src/github.com/gogo/protobuf/protobuf:${GOPATH}/src/github.com/gogo/protobuf:${GOPATH}/src/github.com/google/protobuf/src:${GOPATH}/src:${GOPATH}/src/github.com/onosproject/ran-simulator/"

rm -f api/e2/e2-interface.proto api/e2/e2-interface.pb.go
wget https://raw.githubusercontent.com/onosproject/onos-ran/master/api/sb/e2-interface.proto -P api/e2
protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/e2,plugins=grpc:. api/e2/*.proto
protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/trafficsim,plugins=grpc:. api/trafficsim/*.proto
protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/types,plugins=grpc:. api/types/*.proto
