module github.com/onosproject/ran-simulator

go 1.14

require (
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.3
	github.com/onosproject/config-models/modelplugin/e2node-1.0.0 v0.0.0-20200511074107-7166e3a5247d
	github.com/onosproject/onos-config v0.6.5
	github.com/onosproject/onos-lib-go v0.6.5
	github.com/onosproject/onos-ric v0.6.7
	github.com/onosproject/onos-topo v0.6.9
	github.com/openconfig/gnmi v0.0.0-20190823184014-89b2bf29312c
	github.com/pmcxs/hexgrid v0.0.0-20190126214921-42796ac894ab
	github.com/prometheus/client_golang v1.4.1
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	google.golang.org/grpc v1.29.1
	googlemaps.github.io/maps v0.0.0-20200124220646-5b7f2815585f
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
)

replace github.com/pmcxs/hexgrid v0.0.0-20190126214921-42796ac894ab => github.com/SeanCondon/hexgrid v0.0.0-20200424141352-c3819a378a18
