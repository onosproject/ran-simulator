module github.com/onosproject/ran-simulator

go 1.14

require (
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/magiconair/properties v1.8.1
	github.com/onosproject/config-models/modelplugin/e2node-1.0.0 v0.0.0-20200511074107-7166e3a5247d
	github.com/onosproject/onos-config v0.6.16
	github.com/onosproject/onos-e2t v0.6.12
	github.com/onosproject/onos-lib-go v0.6.25
	github.com/onosproject/onos-ric v0.6.7
	github.com/onosproject/onos-topo v0.6.20
	github.com/openconfig/gnmi v0.0.0-20200617225440-d2b4e6a45802
	github.com/pmcxs/hexgrid v0.0.0-20190126214921-42796ac894ab
	github.com/prometheus/client_golang v1.4.1
	github.com/prometheus/common v0.9.1 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sys v0.0.0-20200803210538-64077c9b5642 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200731012542-8145dea6a485 // indirect
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0 // indirect
	googlemaps.github.io/maps v0.0.0-20200124220646-5b7f2815585f
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
)

replace github.com/pmcxs/hexgrid v0.0.0-20190126214921-42796ac894ab => github.com/SeanCondon/hexgrid v0.0.0-20200424141352-c3819a378a18
