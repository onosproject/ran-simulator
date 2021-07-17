module github.com/onosproject/ran-simulator

go 1.15

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96 // indirect
	github.com/garyburd/redigo v1.1.1-0.20170914051019-70e1b1943d4f // indirect
	github.com/godbus/dbus v0.0.0-20190422162347-ade71ed3457e // indirect
	github.com/golangplus/bytes v0.0.0-20160111154220-45c989fe5450 // indirect
	github.com/golangplus/fmt v0.0.0-20150411045040-2a5d6d7d2995 // indirect
	github.com/google/uuid v1.1.2
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/onosproject/helmit v0.6.13
	github.com/onosproject/onos-api/go v0.7.80
	github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm v0.7.48
	github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2 v0.7.48
	github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho v0.7.48
	github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre v0.7.48
	github.com/onosproject/onos-e2t v0.7.10
	github.com/onosproject/onos-lib-go v0.7.13
	github.com/onosproject/onos-ric-sdk-go v0.7.20
	github.com/onosproject/onos-test v0.6.4
	github.com/onosproject/rrm-son-lib v0.0.2
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/runtime-tools v0.0.0-20181011054405-1d69bd0f9c39 // indirect
	github.com/pmcxs/hexgrid v0.0.0-20190126214921-42796ac894ab
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/gocapability v0.0.0-20170704070218-db04d3cc01c8 // indirect
	github.com/xlab/handysort v0.0.0-20150421192137-fb3537ed64a1 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	gonum.org/v1/netlib v0.0.0-20190331212654-76723241ea4e // indirect
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.26.0
	googlemaps.github.io/maps v1.3.2
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kubernetes v1.13.0 // indirect
	sigs.k8s.io/kustomize v2.0.3+incompatible // indirect
	sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06 // indirect
	vbom.ml/util v0.0.0-20160121211510-db5cfe13f5cc // indirect
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3

replace github.com/pmcxs/hexgrid v0.0.0-20190126214921-42796ac894ab => github.com/SeanCondon/hexgrid v0.0.0-20200424141352-c3819a378a18

replace github.com/onosproject/onos-e2t => ../onos-e2t
