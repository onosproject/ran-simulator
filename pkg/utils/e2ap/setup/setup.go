// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package setup

import (
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-e2t/api/e2ap/v2"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/pdubuilder"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
	"github.com/onosproject/ran-simulator/pkg/utils"

	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-commondatatypes"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("servicemodel", "utils", "setup")

// Setup setup request
type Setup struct {
	ranFunctions                e2aptypes.RanFunctions
	plmnID                      ransimtypes.Uint24
	e2NodeID                    uint64
	componentConfigAdditionList *e2appducontents.E2NodeComponentConfigAdditionList
	transactionID               int32
}

// NewSetupRequest creates a new setup request
func NewSetupRequest(options ...func(*Setup)) *Setup {
	setup := &Setup{}

	for _, option := range options {
		option(setup)
	}

	return setup
}

// WithRanFunctions sets ran functions
func WithRanFunctions(ranFunctions e2aptypes.RanFunctions) func(*Setup) {
	return func(request *Setup) {
		request.ranFunctions = ranFunctions
	}
}

// WithPlmnID sets plmnID
func WithPlmnID(plmnID ransimtypes.Uint24) func(*Setup) {
	return func(request *Setup) {
		request.plmnID = plmnID

	}
}

// WithE2NodeID sets E2 node ID
func WithE2NodeID(e2NodeID uint64) func(*Setup) {
	return func(request *Setup) {
		request.e2NodeID = e2NodeID
	}
}

// WithComponentConfigUpdateList sets E2 node component config update list
func WithComponentConfigUpdateList(componentConfigAdditionList *e2appducontents.E2NodeComponentConfigAdditionList) func(setup *Setup) {
	return func(request *Setup) {
		request.componentConfigAdditionList = componentConfigAdditionList
	}
}

// WithTransactionID sets transaction ID
func WithTransactionID(transID int32) func(setup *Setup) {
	return func(request *Setup) {
		request.transactionID = transID
	}
}

// Build builds e2ap setup request
func (request *Setup) Build() (setupRequest *e2appducontents.E2SetupRequest, err error) {
	//plmnID := types.NewUint24(request.plmnID)
	ge2nID, err := pdubuilder.CreateGlobalE2nodeIDGnb(types.PlmnID(request.plmnID), &asn1.BitString{
		Value: utils.Uint64ToBitString(request.e2NodeID, 28),
		Len:   28,
	})
	if err != nil {
		return nil, err
	}

	cal := &e2appducontents.E2SetupRequestIes{
		Id:          int32(v2.ProtocolIeIDE2nodeComponentConfigAddition),
		Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
		Value: &e2appducontents.E2SetupRequestIe{
			E2SetupRequestIe: &e2appducontents.E2SetupRequestIe_E2NodeComponentConfigAddition{
				E2NodeComponentConfigAddition: request.componentConfigAdditionList,
			},
		},
	}

	setupRequest = &e2appducontents.E2SetupRequest{
		ProtocolIes: make([]*e2appducontents.E2SetupRequestIes, 0),
	}

	setupRequest.SetGlobalE2nodeID(ge2nID).SetRanFunctionsAdded(request.ranFunctions).
		SetTransactionID(request.transactionID)
	setupRequest.ProtocolIes = append(setupRequest.ProtocolIes, cal)

	// TODO enable it when it is available
	/*err = setupRequest.Validate()
	if err != nil {
		log.Warnf("Validation error %s", err.Error())
		return nil, err
	}*/
	log.Debugf("Created E2SetupRequest %v", setupRequest)
	return setupRequest, nil
}
