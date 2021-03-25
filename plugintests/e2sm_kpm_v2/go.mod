module github.com/onosproject/ran-simulator/plugintests/e2sm_kpm_v2

go 1.15

require (
	github.com/onosproject/onos-api/go v0.7.21
	github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2 v0.7.17
	github.com/onosproject/ran-simulator v0.7.21
	google.golang.org/protobuf v1.25.0
	gotest.tools v2.2.0+incompatible
)

replace github.com/onosproject/ran-simulator => ../../
