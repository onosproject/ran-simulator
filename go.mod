module github.com/onosproject/ran-simulator

go 1.15

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/gogo/protobuf v1.3.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onosproject/helmit v0.6.8
	github.com/onosproject/onos-api/go v0.7.0
	github.com/onosproject/onos-e2t v0.7.0
	github.com/onosproject/onos-lib-go v0.7.0
	github.com/onosproject/onos-ric-sdk-go v0.7.0
	github.com/onosproject/onos-test v0.6.4
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.33.2
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onosproject/helmit v0.6.8
	github.com/onosproject/onos-api/go v0.7.0
	github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm v0.7.0
	github.com/onosproject/onos-e2t v0.6.12
	github.com/onosproject/onos-lib-go v0.7.0
	github.com/onosproject/onos-ric-sdk-go v0.7.0
	github.com/onosproject/onos-test v0.6.4
	github.com/prometheus/client_golang v1.4.1 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
	googlemaps.github.io/maps v0.0.0-20200124220646-5b7f2815585f // indirect
	gopkg.in/yaml.v2 v2.3.0
	gotest.tools v2.2.0+incompatible

)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3

replace github.com/pmcxs/hexgrid v0.0.0-20190126214921-42796ac894ab => github.com/SeanCondon/hexgrid v0.0.0-20200424141352-c3819a378a18
