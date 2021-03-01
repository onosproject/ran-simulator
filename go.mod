module github.com/onosproject/ran-simulator

go 1.15

require (
	github.com/Microsoft/go-winio v0.4.15 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/garyburd/redigo v1.1.1-0.20170914051019-70e1b1943d4f // indirect
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/onosproject/helmit v0.6.8
	github.com/onosproject/onos-api/go v0.7.7
	github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm v0.7.7
	github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre v0.7.7
	github.com/onosproject/onos-e2t v0.7.5
	github.com/onosproject/onos-lib-go v0.7.0
	github.com/onosproject/onos-ric-sdk-go v0.7.9
	github.com/onosproject/onos-test v0.6.4
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.3.0
	gotest.tools v2.2.0+incompatible

)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
