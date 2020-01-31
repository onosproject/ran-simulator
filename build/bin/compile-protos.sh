#!/bin/sh

proto_imports=".:${GOPATH}/src/github.com/gogo/protobuf/protobuf:${GOPATH}/src/github.com/gogo/protobuf:${GOPATH}/src/github.com/google/protobuf/src:${GOPATH}/src:${GOPATH}/src/github.com/onosproject/ran-simulator/"

rm -f api/e2/e2-interface.proto api/e2/e2-interface.pb.go
wget https://raw.githubusercontent.com/onosproject/onos-ran/master/api/sb/e2-interface.proto -P api/e2
protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/e2,plugins=grpc:. api/e2/*.proto
protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/trafficsim,plugins=grpc:. api/trafficsim/*.proto
protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/types,plugins=grpc:. api/types/*.proto

# Warning this required protoc v3.9.0 or greater
protoc -I=$proto_imports --js_out=import_style=commonjs:. ${GOPATH}/src/github.com/onosproject/ran-simulator/api/types/types.proto
protoc -I=$proto_imports --js_out=import_style=commonjs:. ${GOPATH}/src/github.com/onosproject/ran-simulator/api/trafficsim/trafficsim.proto

# Currently a bug in the below command outputs to "Github.com" (uppercase G)
# The below uses grpcwebtext as Google implementation does not fully support server side streaming yet (Aug'19)
# See https://grpc.io/blog/state-of-grpc-web/
protoc -I=$proto_imports --grpc-web_out=import_style=typescript,mode=grpcwebtext:. ${GOPATH}/src/github.com/onosproject/ran-simulator/api/types/types.proto
protoc -I=$proto_imports --grpc-web_out=import_style=typescript,mode=grpcwebtext:. ${GOPATH}/src/github.com/onosproject/ran-simulator/api/trafficsim/trafficsim.proto

cp -r github.com/onosproject/ran-simulator/* web/sd-ran-gui/src/app/onos-sdran/proto/github.com/onosproject/ran-simulator/
rm -rf github.com
cp -r Github.com/onosproject/ran-simulator/* web/sd-ran-gui/src/app/onos-sdran/proto/github.com/onosproject/ran-simulator/
rm -rf Github.com

# Add the license text to generated files
for f in $(find web/sd-ran-gui/src/app/onos-*/proto/github.com/ -type f -name "*.d.ts"); do
  cat /build-tools/licensing/boilerplate.generatego.txt | sed -e '$a\\' | cat - $f > tempf && mv tempf $f
done

# Remove unused import for gogoproto
for f in $(find web/sd-ran-gui/src/app/onos-* -type f -name "*ts"); do
  sed -i "s/import \* as gogoproto_gogo_pb from '..\/..\/..\/..\/..\/..\/gogoproto\/gogo_pb';//g" $f
done
