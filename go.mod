module github.com/onosproject/ran-simulator

go 1.14

require (
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.3
	github.com/googleapis/gnostic v0.3.0 // indirect
	github.com/onosproject/onos-lib-go v0.6.2
	github.com/onosproject/onos-ric v0.6.4
	github.com/onosproject/onos-topo v0.6.0
	github.com/openconfig/gnmi v0.0.0-20190823184014-89b2bf29312c
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pmcxs/hexgrid v0.0.0-20190126214921-42796ac894ab
	github.com/prometheus/client_golang v1.4.1
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/cobra v0.0.6
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.6.2
	go.uber.org/multierr v1.4.0 // indirect
	golang.org/x/sys v0.0.0-20200212091648-12a6c2dcc1e4 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/genproto v0.0.0-20200212174721-66ed5ce911ce // indirect
	google.golang.org/grpc v1.27.1
	googlemaps.github.io/maps v0.0.0-20200124220646-5b7f2815585f
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
)

replace github.com/pmcxs/hexgrid v0.0.0-20190126214921-42796ac894ab => github.com/SeanCondon/hexgrid v0.0.0-20200424141352-c3819a378a18
