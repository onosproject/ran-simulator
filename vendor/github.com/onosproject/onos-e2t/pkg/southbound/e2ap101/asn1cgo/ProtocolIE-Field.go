// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package asn1cgo

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "ProtocolIE-Field.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	"unsafe"
)

func newRicSubscriptionDeleteFailureIe1Cause(rsdfCauseIe *e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes1) (*C.RICsubscriptionDeleteFailure_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsdfCauseIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDCause)
	if err != nil {
		return nil, err
	}

	choiceC := [64]byte{} // The size of the RICsubscriptionDeleteFailure_IEs__value

	rsdfCauseC, err := newCause(rsdfCauseIe.GetValue())
	if err != nil {
		return nil, err
	}

	binary.LittleEndian.PutUint64(choiceC[0:], uint64(rsdfCauseC.present))
	copy(choiceC[8:16], rsdfCauseC.choice[:8])

	ie := C.RICsubscriptionDeleteFailure_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionDeleteFailure_IEs__value{
			present: C.RICsubscriptionDeleteFailure_IEs__value_PR_Cause,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newErrorIndicationIe1Cause(eiCauseIe *e2appducontents.ErrorIndicationIes_ErrorIndicationIes1) (*C.ErrorIndication_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(eiCauseIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDCause)
	if err != nil {
		return nil, err
	}

	//TODO: Size should be double-checked
	choiceC := [64]byte{} // The size of the RICsubscriptionDeleteFailure_IEs__value

	rsdfCauseC, err := newCause(eiCauseIe.GetValue())
	if err != nil {
		return nil, err
	}

	binary.LittleEndian.PutUint64(choiceC[0:], uint64(rsdfCauseC.present))
	copy(choiceC[8:16], rsdfCauseC.choice[:8])

	ie := C.ErrorIndication_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_ErrorIndication_IEs__value{
			present: C.ErrorIndication_IEs__value_PR_Cause,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newE2setupFailureIe1Cause(e2sfCauseIe *e2appducontents.E2SetupFailureIes_E2SetupFailureIes1) (*C.E2setupFailureIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(e2sfCauseIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDCause)
	if err != nil {
		return nil, err
	}

	//TODO: Size should be double-checked
	choiceC := [80]byte{} // The size of the RICsubscriptionDeleteFailure_IEs__value

	e2sfCauseC, err := newCause(e2sfCauseIe.GetValue())
	if err != nil {
		return nil, err
	}

	binary.LittleEndian.PutUint64(choiceC[0:], uint64(e2sfCauseC.present))
	copy(choiceC[8:16], e2sfCauseC.choice[:8])

	ie := C.E2setupFailureIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_E2setupFailureIEs__value{
			present: C.E2setupFailureIEs__value_PR_Cause,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newE2connectionUpdateFailureIes1Cause(e2cuaIe *e2appducontents.E2ConnectionUpdateFailureIes_E2ConnectionUpdateFailureIes1) (*C.E2connectionUpdateFailure_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2nodeConfigurationUpdateFailureIes1Cause(e2cuaIe *e2appducontents.E2NodeConfigurationUpdateFailureIes_E2NodeConfigurationUpdateFailureIes1) (*C.E2nodeConfigurationUpdateFailure_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newRicSubscriptionDeleteFailureIe2CriticalityDiagnostics(rsdfCritDiagsIe *e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes2) (*C.RICsubscriptionDeleteFailure_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsdfCritDiagsIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDCriticalityDiagnostics)
	if err != nil {
		return nil, err
	}

	choiceC := [64]byte{} // The size of the RICsubscriptionDeleteFailure_IEs__value

	rsdfCritDiagsC, err := newCriticalityDiagnostics(rsdfCritDiagsIe.GetValue())
	if err != nil {
		return nil, err
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(rsdfCritDiagsC.procedureCode))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(uintptr(unsafe.Pointer(rsdfCritDiagsC.triggeringMessage))))
	binary.LittleEndian.PutUint64(choiceC[16:], uint64(uintptr(unsafe.Pointer(rsdfCritDiagsC.procedureCriticality))))
	binary.LittleEndian.PutUint64(choiceC[24:], uint64(uintptr(unsafe.Pointer(rsdfCritDiagsC.ricRequestorID))))
	binary.LittleEndian.PutUint64(choiceC[40:], uint64(uintptr(unsafe.Pointer(rsdfCritDiagsC.iEsCriticalityDiagnostics))))

	ie := C.RICsubscriptionDeleteFailure_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionDeleteFailure_IEs__value{
			present: C.RICsubscriptionDeleteFailure_IEs__value_PR_CriticalityDiagnostics,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newErrorIndicationIe2CriticalityDiagnostics(eiCritDiagsIe *e2appducontents.ErrorIndicationIes_ErrorIndicationIes2) (*C.ErrorIndication_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(eiCritDiagsIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDCriticalityDiagnostics)
	if err != nil {
		return nil, err
	}

	//TODO: Size should be double-checked
	choiceC := [64]byte{} // The size of the ErrorIndication_IEs__value

	eiCritDiagsC, err := newCriticalityDiagnostics(eiCritDiagsIe.GetValue())
	if err != nil {
		return nil, err
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(eiCritDiagsC.procedureCode))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(uintptr(unsafe.Pointer(eiCritDiagsC.triggeringMessage))))
	binary.LittleEndian.PutUint64(choiceC[16:], uint64(uintptr(unsafe.Pointer(eiCritDiagsC.procedureCriticality))))
	binary.LittleEndian.PutUint64(choiceC[24:], uint64(uintptr(unsafe.Pointer(eiCritDiagsC.ricRequestorID))))
	binary.LittleEndian.PutUint64(choiceC[40:], uint64(uintptr(unsafe.Pointer(eiCritDiagsC.iEsCriticalityDiagnostics))))

	ie := C.ErrorIndication_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_ErrorIndication_IEs__value{
			present: C.ErrorIndication_IEs__value_PR_CriticalityDiagnostics,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newE2setupIe2CriticalityDiagnostics(e2sfCritDiagsIe *e2appducontents.E2SetupFailureIes_E2SetupFailureIes2) (*C.E2setupFailureIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(e2sfCritDiagsIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDCriticalityDiagnostics)
	if err != nil {
		return nil, err
	}

	//TODO: Size should be double-checked
	choiceC := [80]byte{} // The size of the ErrorIndication_IEs__value

	e2sfCritDiagsC, err := newCriticalityDiagnostics(e2sfCritDiagsIe.GetValue())
	if err != nil {
		return nil, err
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(e2sfCritDiagsC.procedureCode))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(uintptr(unsafe.Pointer(e2sfCritDiagsC.triggeringMessage))))
	binary.LittleEndian.PutUint64(choiceC[16:], uint64(uintptr(unsafe.Pointer(e2sfCritDiagsC.procedureCriticality))))
	binary.LittleEndian.PutUint64(choiceC[24:], uint64(uintptr(unsafe.Pointer(e2sfCritDiagsC.ricRequestorID))))
	binary.LittleEndian.PutUint64(choiceC[40:], uint64(uintptr(unsafe.Pointer(e2sfCritDiagsC.iEsCriticalityDiagnostics))))

	ie := C.E2setupFailureIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_E2setupFailureIEs__value{
			present: C.E2setupFailureIEs__value_PR_CriticalityDiagnostics,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionFailureIe2CriticalityDiagnostics(rsfCritDiagsIe *e2appducontents.RicsubscriptionFailureIes_RicsubscriptionFailureIes2) (*C.RICsubscriptionFailure_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsfCritDiagsIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDCriticalityDiagnostics)
	if err != nil {
		return nil, err
	}

	choiceC := [64]byte{} // The size of the RICsubscriptionFailure_IEs__value

	rsfCritDiagsC, err := newCriticalityDiagnostics(rsfCritDiagsIe.GetValue())
	if err != nil {
		return nil, err
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(rsfCritDiagsC.procedureCode))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(uintptr(unsafe.Pointer(rsfCritDiagsC.triggeringMessage))))
	binary.LittleEndian.PutUint64(choiceC[16:], uint64(uintptr(unsafe.Pointer(rsfCritDiagsC.procedureCriticality))))
	binary.LittleEndian.PutUint64(choiceC[24:], uint64(uintptr(unsafe.Pointer(rsfCritDiagsC.ricRequestorID))))
	binary.LittleEndian.PutUint64(choiceC[40:], uint64(uintptr(unsafe.Pointer(rsfCritDiagsC.iEsCriticalityDiagnostics))))

	ie := C.RICsubscriptionFailure_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionFailure_IEs__value{
			present: C.RICsubscriptionFailure_IEs__value_PR_CriticalityDiagnostics,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newE2connectionUpdateFailureIes2CriticalityDiagnostics(e2cufIe *e2appducontents.E2ConnectionUpdateFailureIes_E2ConnectionUpdateFailureIes2) (*C.E2connectionUpdateFailure_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2nodeConfigurationUpdateFailureIes2CriticalityDiagnostics(e2ncufIe *e2appducontents.E2NodeConfigurationUpdateFailureIes_E2NodeConfigurationUpdateFailureIes2) (*C.E2nodeConfigurationUpdateFailure_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2setupRequestIe3GlobalE2NodeID(esIe *e2appducontents.E2SetupRequestIes_E2SetupRequestIes3) (*C.E2setupRequestIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(esIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDGlobalE2nodeID)
	if err != nil {
		return nil, err
	}

	choiceC := [48]byte{} // The size of the E2setupRequestIEs__value_u

	globalNodeIDC, err := newGlobalE2nodeID(esIe.GetValue())
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Assigning to choice of E2setupRequestIE %v %v %v %v %v\n",
	//	globalNodeIDC, globalNodeIDC.present, &globalNodeIDC.choice,
	//	unsafe.Sizeof(globalNodeIDC.present), unsafe.Sizeof(globalNodeIDC.choice))
	binary.LittleEndian.PutUint32(choiceC[0:], uint32(globalNodeIDC.present))
	for i := 0; i < 8; i++ {
		choiceC[i+8] = globalNodeIDC.choice[i]
	}

	ie := C.E2setupRequestIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_E2setupRequestIEs__value{
			present: C.E2setupRequestIEs__value_PR_GlobalE2node_ID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newE2setupResponseIe4GlobalRicID(esIe *e2appducontents.E2SetupResponseIes_E2SetupResponseIes4) (*C.E2setupResponseIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(esIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDGlobalRicID)
	if err != nil {
		return nil, err
	}

	choiceC := [112]byte{} // The size of the E2setupResponseIEs__value_u

	globalRicIDC, err := newGlobalRicID(esIe.Value)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Assigning to choice of E2setupReponseIE %v \n", globalRicIDC)

	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(globalRicIDC.pLMN_Identity.buf))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(globalRicIDC.pLMN_Identity.size))
	binary.LittleEndian.PutUint64(choiceC[40:], uint64(uintptr(unsafe.Pointer(globalRicIDC.ric_ID.buf))))
	binary.LittleEndian.PutUint64(choiceC[48:], uint64(globalRicIDC.ric_ID.size))
	binary.LittleEndian.PutUint32(choiceC[56:], uint32(globalRicIDC.ric_ID.bits_unused))

	ie := C.E2setupResponseIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_E2setupResponseIEs__value{
			present: C.E2setupResponseIEs__value_PR_GlobalRIC_ID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionRequestIe5RanFunctionID(rsrRfIe *e2appducontents.RicsubscriptionRequestIes_RicsubscriptionRequestIes5) (*C.RICsubscriptionRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrRfIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionID)
	if err != nil {
		return nil, err
	}

	choiceC := [112]byte{} // The size of the E2setupResponseIEs__value_u

	ranFunctionIDC := newRanFunctionID(rsrRfIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionRequestIE %v \n", ranFunctionIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ranFunctionIDC))

	ie := C.RICsubscriptionRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionRequest_IEs__value{
			present: C.RICsubscriptionRequest_IEs__value_PR_RANfunctionID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicControlRequestIe5RanFunctionID(rcrRfIe *e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes5) (*C.RICcontrolRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rcrRfIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ranFunctionIDC := newRanFunctionID(rcrRfIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionRequestIE %v \n", ranFunctionIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ranFunctionIDC))

	ie := C.RICcontrolRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICcontrolRequest_IEs__value{
			present: C.RICcontrolRequest_IEs__value_PR_RANfunctionID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicControlAcknowledgeIe5RanFunctionID(rcaRfIe *e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes5) (*C.RICcontrolAcknowledge_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rcaRfIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ranFunctionIDC := newRanFunctionID(rcaRfIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionRequestIE %v \n", ranFunctionIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ranFunctionIDC))

	ie := C.RICcontrolAcknowledge_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICcontrolAcknowledge_IEs__value{
			present: C.RICcontrolAcknowledge_IEs__value_PR_RANfunctionID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionResponseIe5RanFunctionID(rsrRfIe *e2appducontents.RicsubscriptionResponseIes_RicsubscriptionResponseIes5) (*C.RICsubscriptionResponse_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrRfIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionID)
	if err != nil {
		return nil, err
	}

	choiceC := [48]byte{} // The size of the E2setupResponseIEs__value_u

	ranFunctionIDC := newRanFunctionID(rsrRfIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionResponseIE %v \n", ranFunctionIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ranFunctionIDC))

	ie := C.RICsubscriptionResponse_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionResponse_IEs__value{
			present: C.RICsubscriptionResponse_IEs__value_PR_RANfunctionID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicIndicationIe5RanFunctionID(rsrRfIe *e2appducontents.RicindicationIes_RicindicationIes5) (*C.RICsubscriptionRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrRfIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionID)
	if err != nil {
		return nil, err
	}

	choiceC := [112]byte{} // The size of the E2setupResponseIEs__value_u

	ranFunctionIDC := newRanFunctionID(rsrRfIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionRequestIE %v \n", ranFunctionIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ranFunctionIDC))

	ie := C.RICsubscriptionRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionRequest_IEs__value{
			present: C.RICsubscriptionRequest_IEs__value_PR_RANfunctionID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionDeleteRequestIe5RanFunctionID(rsdrRfIe *e2appducontents.RicsubscriptionDeleteRequestIes_RicsubscriptionDeleteRequestIes5) (*C.RICsubscriptionDeleteRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsdrRfIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ranFunctionIDC := newRanFunctionID(rsdrRfIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionRequestIE %v \n", ranFunctionIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ranFunctionIDC))

	ie := C.RICsubscriptionDeleteRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionDeleteRequest_IEs__value{
			present: C.RICsubscriptionDeleteRequest_IEs__value_PR_RANfunctionID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionDeleteResponseIe5RanFunctionID(rsdrRfIe *e2appducontents.RicsubscriptionDeleteResponseIes_RicsubscriptionDeleteResponseIes5) (*C.RICsubscriptionDeleteResponse_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsdrRfIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ranFunctionIDC := newRanFunctionID(rsdrRfIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionResponseIE %v \n", ranFunctionIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ranFunctionIDC))

	ie := C.RICsubscriptionDeleteResponse_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionDeleteResponse_IEs__value{
			present: C.RICsubscriptionDeleteResponse_IEs__value_PR_RANfunctionID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionDeleteFailureIe5RanFunctionID(rsdfRfIe *e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes5) (*C.RICsubscriptionDeleteFailure_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsdfRfIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionID)
	if err != nil {
		return nil, err
	}

	choiceC := [64]byte{} // The size of the RICsubscriptionDeleteFailure_IEs__value_u

	ranFunctionIDC := newRanFunctionID(rsdfRfIe.Value)

	//fmt.Printf("Assigning to choice of RICsubscriptionDeleteFailureIE %v \n", ranFunctionIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ranFunctionIDC))

	ie := C.RICsubscriptionDeleteFailure_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionDeleteFailure_IEs__value{
			present: C.RICsubscriptionDeleteFailure_IEs__value_PR_RANfunctionID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionFailureIe5RanFunctionID(rsfRfIe *e2appducontents.RicsubscriptionFailureIes_RicsubscriptionFailureIes5) (*C.RICsubscriptionFailure_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsfRfIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionID)
	if err != nil {
		return nil, err
	}

	choiceC := [64]byte{} // The size of the RICsubscriptionDeleteFailure_IEs__value_u

	ranFunctionIDC := newRanFunctionID(rsfRfIe.Value)

	//fmt.Printf("Assigning to choice of RICsubscriptionDeleteFailureIE %v \n", ranFunctionIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ranFunctionIDC))

	ie := C.RICsubscriptionFailure_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionFailure_IEs__value{
			present: C.RICsubscriptionFailure_IEs__value_PR_RANfunctionID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newErrorIndicationIe5RanFunctionID(eiRfIe *e2appducontents.ErrorIndicationIes_ErrorIndicationIes5) (*C.ErrorIndication_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(eiRfIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionID)
	if err != nil {
		return nil, err
	}

	//TODO: Size should be double-checked
	choiceC := [64]byte{} // The size of the ErrorIndication_IEs__value_u

	ranFunctionIDC := newRanFunctionID(eiRfIe.Value)

	//fmt.Printf("Assigning to choice of ErrorIndicationIE %v \n", ranFunctionIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ranFunctionIDC))

	ie := C.ErrorIndication_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_ErrorIndication_IEs__value{
			present: C.ErrorIndication_IEs__value_PR_RANfunctionID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newE2setupResponseIe9RanFunctionsAccepted(esIe *e2appducontents.E2SetupResponseIes_E2SetupResponseIes9) (*C.E2setupResponseIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(esIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionsAccepted)
	if err != nil {
		return nil, err
	}

	choiceC := [112]byte{} // The size of the E2setupResponseIEs__value_u

	ranFunctionsIDListC, err := newRanFunctionsIDList(esIe.Value)
	if err != nil {
		return nil, fmt.Errorf("newRanFunctionsIDList() %s", err.Error())
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(ranFunctionsIDListC.list.array))))
	binary.LittleEndian.PutUint32(choiceC[8:], uint32(ranFunctionsIDListC.list.count))
	binary.LittleEndian.PutUint32(choiceC[12:], uint32(ranFunctionsIDListC.list.size))

	ie := C.E2setupResponseIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_E2setupResponseIEs__value{
			present: C.E2setupResponseIEs__value_PR_RANfunctionsID_List,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newE2setupRequestIe10RanFunctionList(esIe *e2appducontents.E2SetupRequestIes_E2SetupRequestIes10) (*C.E2setupRequestIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(esIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionsAdded)
	if err != nil {
		return nil, err
	}

	listC := [48]byte{} // The size of the E2setupRequestIEs__value_u

	ranFunctionsListC, err := newRanFunctionsList(esIe.GetValue())
	if err != nil {
		return nil, fmt.Errorf("newRanFunctionsList() %s", err.Error())
	}
	binary.LittleEndian.PutUint64(listC[0:], uint64(uintptr(unsafe.Pointer(ranFunctionsListC.list.array))))
	binary.LittleEndian.PutUint32(listC[8:], uint32(ranFunctionsListC.list.count))
	binary.LittleEndian.PutUint32(listC[12:], uint32(ranFunctionsListC.list.size))

	ie := C.E2setupRequestIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_E2setupRequestIEs__value{
			present: C.E2setupRequestIEs__value_PR_RANfunctions_List,
			choice:  listC,
		},
	}

	return &ie, nil
}

func newE2setupResponseIe13RanFunctionsRejected(esIe *e2appducontents.E2SetupResponseIes_E2SetupResponseIes13) (*C.E2setupResponseIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(esIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionsRejected)
	if err != nil {
		return nil, err
	}

	choiceC := [112]byte{} // The size of the E2setupResponseIEs__value_u

	ranFunctionsIDCauseList, err := newRanFunctionsIDcauseList(esIe.Value)
	if err != nil {
		return nil, fmt.Errorf("newRanFunctionsIDcauseList() %s", err.Error())
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(ranFunctionsIDCauseList.list.array))))
	binary.LittleEndian.PutUint32(choiceC[8:], uint32(ranFunctionsIDCauseList.list.count))
	binary.LittleEndian.PutUint32(choiceC[12:], uint32(ranFunctionsIDCauseList.list.size))

	ie := C.E2setupResponseIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_E2setupResponseIEs__value{
			present: C.E2setupResponseIEs__value_PR_RANfunctionsIDcause_List,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicIndicationIe15RicActionID(riIe *e2appducontents.RicindicationIes_RicindicationIes15) (*C.RICindication_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(riIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicactionID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the RICindication_IEs__value_u

	ricActionID := newRicActionID(riIe.Value)

	binary.LittleEndian.PutUint64(choiceC[0:], uint64(*ricActionID))

	ie := C.RICindication_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICindication_IEs__value{
			present: C.RICindication_IEs__value_PR_RICactionID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionResponseIe17RactionAdmittedList(rsrRrIe *e2appducontents.RicsubscriptionResponseIes_RicsubscriptionResponseIes17) (*C.RICsubscriptionResponse_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrRrIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicactionsAdmitted)
	if err != nil {
		return nil, err
	}

	listC := [48]byte{} // The size of the E2setupResponseIEs__value_u

	ricActionAdmittedListC, err := newRicActionAdmittedList(rsrRrIe.Value)
	if err != nil {
		return nil, fmt.Errorf("newRicActionAdmittedList() %s", err.Error())
	}
	binary.LittleEndian.PutUint64(listC[0:], uint64(uintptr(unsafe.Pointer(ricActionAdmittedListC.list.array))))
	binary.LittleEndian.PutUint32(listC[8:], uint32(ricActionAdmittedListC.list.count))
	binary.LittleEndian.PutUint32(listC[12:], uint32(ricActionAdmittedListC.list.size))

	ie := C.RICsubscriptionResponse_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionResponse_IEs__value{
			present: C.RICsubscriptionResponse_IEs__value_PR_RICaction_Admitted_List,
			choice:  listC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionFailureIe18RicActionNotAdmittedList(rsfRanaIe *e2appducontents.RicsubscriptionFailureIes_RicsubscriptionFailureIes18) (*C.RICsubscriptionFailure_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsfRanaIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicactionsNotAdmitted)
	if err != nil {
		return nil, err
	}

	listC := [64]byte{} // The size of the E2setupResponseIEs__value_u

	ricActionNotAdmittedListC, err := newRicActionNotAdmittedList(rsfRanaIe.Value)
	if err != nil {
		return nil, fmt.Errorf("newRicActionAdmittedList() %s", err.Error())
	}
	binary.LittleEndian.PutUint64(listC[0:], uint64(uintptr(unsafe.Pointer(ricActionNotAdmittedListC.list.array))))
	binary.LittleEndian.PutUint32(listC[8:], uint32(ricActionNotAdmittedListC.list.count))
	binary.LittleEndian.PutUint32(listC[12:], uint32(ricActionNotAdmittedListC.list.size))

	ie := C.RICsubscriptionFailure_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionFailure_IEs__value{
			present: C.RICsubscriptionFailure_IEs__value_PR_RICaction_NotAdmitted_List,
			choice:  listC,
		},
	}

	return &ie, nil
}

func newRicIndicationIe20RiccallProcessID(riIe20 *e2appducontents.RicindicationIes_RicindicationIes20) (*C.RICindication_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(riIe20.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRiccallProcessID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ricCallProcessID := newRicCallProcessID(riIe20.Value)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(ricCallProcessID.buf))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricCallProcessID.size))

	ie := C.RICindication_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICindication_IEs__value{
			present: C.RICindication_IEs__value_PR_RICcallProcessID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicControlRequestIe20RiccallProcessID(rcrIe20 *e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes20) (*C.RICcontrolRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rcrIe20.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRiccallProcessID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ricCallProcessIDC := newRicCallProcessID(rcrIe20.Value)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(ricCallProcessIDC.buf))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricCallProcessIDC.size))

	ie := C.RICcontrolRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICcontrolRequest_IEs__value{
			present: C.RICcontrolRequest_IEs__value_PR_RICcallProcessID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicControlAcknowledgeIe20RiccallProcessID(rcrIe20 *e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes20) (*C.RICcontrolAcknowledge_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rcrIe20.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRiccallProcessID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ricCallProcessIDC := newRicCallProcessID(rcrIe20.Value)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(ricCallProcessIDC.buf))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricCallProcessIDC.size))

	ie := C.RICcontrolAcknowledge_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICcontrolAcknowledge_IEs__value{
			present: C.RICcontrolAcknowledge_IEs__value_PR_RICcallProcessID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicControlRequestIe21RiccontrolAckRequest(rcrIe21 *e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes21) (*C.RICcontrolRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rcrIe21.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRiccontrolAckRequest)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ricControlAckRequestC, err := newRicControlAckRequest(rcrIe21.Value)
	if err != nil {
		return nil, fmt.Errorf("newRicControlAckRequest() %s", err.Error())
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(*ricControlAckRequestC))

	ie := C.RICcontrolRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICcontrolRequest_IEs__value{
			present: C.RICcontrolRequest_IEs__value_PR_RICcontrolAckRequest,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicControlRequestIe22RiccontrolHeader(rcrIe22 *e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes22) (*C.RICcontrolRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rcrIe22.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRiccontrolHeader)
	if err != nil {
		return nil, err
	}

	//ToDo - double-check number of bytes required here
	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ricControlHeaderC := newRicControlHeader(rcrIe22.Value)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(ricControlHeaderC.buf))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricControlHeaderC.size))

	ie := C.RICcontrolRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICcontrolRequest_IEs__value{
			present: C.RICcontrolRequest_IEs__value_PR_RICcontrolHeader,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicControlRequestIe23RiccontrolMessage(rcrIe23 *e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes23) (*C.RICcontrolRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rcrIe23.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRiccontrolMessage)
	if err != nil {
		return nil, err
	}

	//ToDo - double-check number of bytes required here
	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ricControlMessageC := newRicControlMessage(rcrIe23.Value)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(ricControlMessageC.buf))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricControlMessageC.size))

	ie := C.RICcontrolRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICcontrolRequest_IEs__value{
			present: C.RICcontrolRequest_IEs__value_PR_RICcontrolMessage,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicControlAcknowledgeIe24RiccontrolStatus(rcrIe24 *e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes24) (*C.RICcontrolAcknowledge_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rcrIe24.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRiccontrolStatus)
	if err != nil {
		return nil, err
	}

	//ToDo - double-check number of bytes required here
	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ricControlStatusC, err := newRicControlStatus(rcrIe24.Value)
	if err != nil {
		return nil, fmt.Errorf("newRicControlAckRequest() %s", err.Error())
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(*ricControlStatusC))

	ie := C.RICcontrolAcknowledge_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICcontrolAcknowledge_IEs__value{
			present: C.RICcontrolAcknowledge_IEs__value_PR_RICcontrolStatus,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicIndicationIe25RicIndicationHeader(rihIe *e2appducontents.RicindicationIes_RicindicationIes25) (*C.RICindication_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rihIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicindicationHeader)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the RICindication_IEs__value_u

	ricIndicationHeader := newRicIndicationHeader(rihIe.Value)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(ricIndicationHeader.buf))))
	binary.LittleEndian.PutUint32(choiceC[8:], uint32(ricIndicationHeader.size))

	ie := C.RICindication_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICindication_IEs__value{
			present: C.RICindication_IEs__value_PR_RICindicationHeader,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicIndicationIe26RicIndicationMessage(rimIe *e2appducontents.RicindicationIes_RicindicationIes26) (*C.RICindication_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rimIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicindicationMessage)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the RICindication_IEs__value_u

	ricIndicationMessage := newRicIndicationMessage(rimIe.Value)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(ricIndicationMessage.buf))))
	binary.LittleEndian.PutUint32(choiceC[8:], uint32(ricIndicationMessage.size))

	ie := C.RICindication_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICindication_IEs__value{
			present: C.RICindication_IEs__value_PR_RICindicationMessage,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicIndicationIe27RicIndicationSn(risnIe *e2appducontents.RicindicationIes_RicindicationIes27) (*C.RICindication_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(risnIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicindicationSn)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the RICindication_IEs__value_u

	ricIndicationSequenceNumber := newRicIndicationSn(risnIe.Value)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(*ricIndicationSequenceNumber))

	ie := C.RICindication_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICindication_IEs__value{
			present: C.RICindication_IEs__value_PR_RICindicationSN,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicIndicationIe28RicIndicationType(ritIe *e2appducontents.RicindicationIes_RicindicationIes28) (*C.RICindication_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(ritIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicindicationType)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the RICindication_IEs__value_u

	ricIndicationTypeC, err := newRicIndicationType(&ritIe.Value)
	if err != nil {
		return nil, fmt.Errorf("newRicIndicationType() %s", err.Error())
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(*ricIndicationTypeC))

	ie := C.RICindication_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICindication_IEs__value{
			present: C.RICindication_IEs__value_PR_RICindicationType,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicIndicationIe29RicRequestID(rsrRrIDIe *e2appducontents.RicindicationIes_RicindicationIes29) (*C.RICindication_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrRrIDIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicrequestID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the RICindication_IEs__value_u

	ricRequestIDC := newRicRequestID(rsrRrIDIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionRequestIE %v \n", ricRequestIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ricRequestIDC.ricRequestorID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricRequestIDC.ricInstanceID))

	ie := C.RICindication_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICindication_IEs__value{
			present: C.RICindication_IEs__value_PR_RICrequestID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicControlRequestIe29RicRequestID(rcrRrIDIe *e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes29) (*C.RICcontrolRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rcrRrIDIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicrequestID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the RICindication_IEs__value_u

	ricRequestIDC := newRicRequestID(rcrRrIDIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionRequestIE %v \n", ricRequestIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ricRequestIDC.ricRequestorID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricRequestIDC.ricInstanceID))

	ie := C.RICcontrolRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICcontrolRequest_IEs__value{
			present: C.RICcontrolRequest_IEs__value_PR_RICrequestID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionRequestIe29RicRequestID(rsrRrIDIe *e2appducontents.RicsubscriptionRequestIes_RicsubscriptionRequestIes29) (*C.RICsubscriptionRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrRrIDIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicrequestID)
	if err != nil {
		return nil, err
	}

	choiceC := [112]byte{} // The size of the E2setupResponseIEs__value_u

	ricRequestIDC := newRicRequestID(rsrRrIDIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionRequestIE %v \n", ricRequestIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ricRequestIDC.ricRequestorID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricRequestIDC.ricInstanceID))

	ie := C.RICsubscriptionRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionRequest_IEs__value{
			present: C.RICsubscriptionRequest_IEs__value_PR_RICrequestID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionResponseIe29RicRequestID(rsrRrIDIe *e2appducontents.RicsubscriptionResponseIes_RicsubscriptionResponseIes29) (*C.RICsubscriptionResponse_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrRrIDIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicrequestID)
	if err != nil {
		return nil, err
	}

	choiceC := [48]byte{}

	ricRequestIDC := newRicRequestID(rsrRrIDIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionResponseIE %v \n", ricRequestIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ricRequestIDC.ricRequestorID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricRequestIDC.ricInstanceID))

	ie := C.RICsubscriptionResponse_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionResponse_IEs__value{
			present: C.RICsubscriptionResponse_IEs__value_PR_RICrequestID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionDeleteRequestIe29RicRequestID(rsrdRrIDIe *e2appducontents.RicsubscriptionDeleteRequestIes_RicsubscriptionDeleteRequestIes29) (*C.RICsubscriptionDeleteRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrdRrIDIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicrequestID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ricRequestIDC := newRicRequestID(rsrdRrIDIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionRequestIE %v \n", ricRequestIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ricRequestIDC.ricRequestorID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricRequestIDC.ricInstanceID))

	ie := C.RICsubscriptionDeleteRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionDeleteRequest_IEs__value{
			present: C.RICsubscriptionDeleteRequest_IEs__value_PR_RICrequestID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionDeleteResponseIe29RicRequestID(rsrRrIDIe *e2appducontents.RicsubscriptionDeleteResponseIes_RicsubscriptionDeleteResponseIes29) (*C.RICsubscriptionDeleteResponse_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrRrIDIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicrequestID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{}

	ricRequestIDC := newRicRequestID(rsrRrIDIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionResponseIE %v \n", ricRequestIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ricRequestIDC.ricRequestorID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricRequestIDC.ricInstanceID))

	ie := C.RICsubscriptionDeleteResponse_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionDeleteResponse_IEs__value{
			present: C.RICsubscriptionDeleteResponse_IEs__value_PR_RICrequestID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionDeleteFailureIe29RicRequestID(rsrRrIDIe *e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes29) (*C.RICsubscriptionDeleteFailure_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrRrIDIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicrequestID)
	if err != nil {
		return nil, err
	}

	choiceC := [64]byte{}

	ricRequestIDC := newRicRequestID(rsrRrIDIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionResponseIE %v \n", ricRequestIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ricRequestIDC.ricRequestorID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricRequestIDC.ricInstanceID))

	ie := C.RICsubscriptionDeleteFailure_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionDeleteFailure_IEs__value{
			present: C.RICsubscriptionDeleteFailure_IEs__value_PR_RICrequestID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newErrorIndicationIe29RicRequestID(eiRrIDIe *e2appducontents.ErrorIndicationIes_ErrorIndicationIes29) (*C.ErrorIndication_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(eiRrIDIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicrequestID)
	if err != nil {
		return nil, err
	}

	//TODO: Size should be double-checked
	choiceC := [64]byte{}

	ricRequestIDC := newRicRequestID(eiRrIDIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionResponseIE %v \n", ricRequestIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ricRequestIDC.ricRequestorID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricRequestIDC.ricInstanceID))

	ie := C.ErrorIndication_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_ErrorIndication_IEs__value{
			present: C.ErrorIndication_IEs__value_PR_RICrequestID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionFailureIe29RicRequestID(rsrRrIDIe *e2appducontents.RicsubscriptionFailureIes_RicsubscriptionFailureIes29) (*C.RICsubscriptionFailure_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrRrIDIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicrequestID)
	if err != nil {
		return nil, err
	}

	choiceC := [64]byte{}

	ricRequestIDC := newRicRequestID(rsrRrIDIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionResponseIE %v \n", ricRequestIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ricRequestIDC.ricRequestorID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricRequestIDC.ricInstanceID))

	ie := C.RICsubscriptionFailure_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionFailure_IEs__value{
			present: C.RICsubscriptionFailure_IEs__value_PR_RICrequestID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicControlAcknowledgeIe29RicRequestID(rsrRrIDIe *e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes29) (*C.RICcontrolAcknowledge_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrRrIDIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicrequestID)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{}

	ricRequestIDC := newRicRequestID(rsrRrIDIe.Value)

	//fmt.Printf("Assigning to choice of RicSubscriptionResponseIE %v \n", ricRequestIDC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ricRequestIDC.ricRequestorID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricRequestIDC.ricInstanceID))

	ie := C.RICcontrolAcknowledge_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICcontrolAcknowledge_IEs__value{
			present: C.RICcontrolAcknowledge_IEs__value_PR_RICrequestID,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newRicSubscriptionRequestIe30RicSubscriptionDetails(rsrDetIe *e2appducontents.RicsubscriptionRequestIes_RicsubscriptionRequestIes30) (*C.RICsubscriptionRequest_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rsrDetIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicsubscriptionDetails)
	if err != nil {
		return nil, err
	}

	choiceC := [112]byte{} // The size of the E2setupResponseIEs__value_u

	rsrDetC, err := newRicSubscriptionDetails(rsrDetIe.GetValue())
	if err != nil {
		return nil, err
	}

	//fmt.Printf("Assigning to choice of RicSubscriptionRequestIE %v \n", rsrDetC)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(rsrDetC.ricEventTriggerDefinition.buf))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(rsrDetC.ricEventTriggerDefinition.size))
	binary.LittleEndian.PutUint64(choiceC[40:], uint64(uintptr(unsafe.Pointer(rsrDetC.ricAction_ToBeSetup_List.list.array))))
	binary.LittleEndian.PutUint32(choiceC[48:], uint32(rsrDetC.ricAction_ToBeSetup_List.list.count))
	binary.LittleEndian.PutUint32(choiceC[52:], uint32(rsrDetC.ricAction_ToBeSetup_List.list.size))

	ie := C.RICsubscriptionRequest_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICsubscriptionRequest_IEs__value{
			present: C.RICsubscriptionRequest_IEs__value_PR_RICsubscriptionDetails,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newE2setupFailureIe31TimeToWait(e2sfIe *e2appducontents.E2SetupFailureIes_E2SetupFailureIes31) (*C.E2setupFailureIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(e2sfIe.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDTimeToWait)
	if err != nil {
		return nil, err
	}

	//TODO: Size should be double-checked
	choiceC := [80]byte{} // The size of the RICsubscriptionDeleteFailure_IEs__value

	e2sfTtwC, err := newTimeToWait(e2sfIe.GetValue())
	if err != nil {
		return nil, err
	}

	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(e2sfTtwC))))
	//copy(choiceC[8:16], e2sfCauseC.choice[:8])

	ie := C.E2setupFailureIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_E2setupFailureIEs__value{
			present: C.E2setupFailureIEs__value_PR_TimeToWait,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newE2connectionUpdateFailureIes31TimeToWait(e2cuaIe *e2appducontents.E2ConnectionUpdateFailureIes_E2ConnectionUpdateFailureIes31) (*C.E2connectionUpdateFailure_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2nodeConfigurationUpdateFailureIes31TimeToWait(e2cuaIe *e2appducontents.E2NodeConfigurationUpdateFailureIes_E2NodeConfigurationUpdateFailureIes31) (*C.E2nodeConfigurationUpdateFailure_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newRicControlAcknowledgeIe32RiccontrolOutcome(rcrIe32 *e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes32) (*C.RICcontrolAcknowledge_IEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rcrIe32.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRiccontrolOutcome)
	if err != nil {
		return nil, err
	}

	//ToDo - double-check number of bytes required here
	choiceC := [40]byte{} // The size of the E2setupResponseIEs__value_u

	ricControlOutcomeC := newRicControlOutcome(rcrIe32.Value)
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(uintptr(unsafe.Pointer(ricControlOutcomeC.buf))))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ricControlOutcomeC.size))

	ie := C.RICcontrolAcknowledge_IEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICcontrolAcknowledge_IEs__value{
			present: C.RICcontrolAcknowledge_IEs__value_PR_RICcontrolOutcome,
			choice:  choiceC,
		},
	}

	return &ie, nil
}

func newE2setupRequestIe33E2nodeComponentConfigUpdateList(e2sfIe *e2appducontents.E2SetupRequestIes_E2SetupRequestIes33) (*C.E2setupRequestIEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2nodeConfigurationUpdateIe33E2nodeComponentConfigUpdateList(e2ncuIe *e2appducontents.E2NodeConfigurationUpdateIes) (*C.E2nodeConfigurationUpdate_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2setupResponseIe35E2nodeComponentConfigUpdateAckList(e2sfIe *e2appducontents.E2SetupResponseIes_E2SetupResponseIes35) (*C.E2setupResponseIEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2nodeConfigurationUpdateIe35E2nodeComponentConfigUpdateAckList(e2sfIe *e2appducontents.E2NodeConfigurationUpdateAcknowledgeIes) (*C.E2nodeConfigurationUpdateAcknowledge_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2connectionUpdateAck39E2connectionUpdateList(e2cuaIe *e2appducontents.E2ConnectionUpdateAckIes_E2ConnectionUpdateAckIes39) (*C.E2connectionUpdateAck_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2connectionUpdateAck40E2connectionSetupFailedList(e2cuaIe *e2appducontents.E2ConnectionUpdateAckIes_E2ConnectionUpdateAckIes40) (*C.E2connectionUpdateAck_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2connectionUpdateIe44E2connectionUpdateList(e2cuIe *e2appducontents.E2ConnectionUpdateIes_E2ConnectionUpdateIes44) (*C.E2connectionUpdate_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2connectionUpdateIe45E2connectionUpdateList(e2cuIe *e2appducontents.E2ConnectionUpdateIes_E2ConnectionUpdateIes45) (*C.E2connectionUpdate_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2connectionUpdateIe46E2connectionUpdateRemoveList(e2cuIe *e2appducontents.E2ConnectionUpdateIes_E2ConnectionUpdateIes46) (*C.E2connectionUpdate_IEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newE2setupFailureIe48Tnlinformation(e2sfIe *e2appducontents.E2SetupFailureIes_E2SetupFailureIes48) (*C.E2setupResponseIEs_t, error) {
	// TODO new for E2AP 1.0.1
	return nil, fmt.Errorf("not yet implemented - new for E2AP 1.0.1")
}

func newRANfunctionItemIEs(rfItemIes *e2appducontents.RanfunctionItemIes) (*C.RANfunction_ItemIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rfItemIes.GetE2ApProtocolIes10().GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionItem)
	if err != nil {
		return nil, err
	}

	choiceC := [88]byte{} // The size of the RANfunction_ItemIEs__value_u
	rfItemC := newRanFunctionItem(rfItemIes.GetE2ApProtocolIes10().GetValue())
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(rfItemC.ranFunctionID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(uintptr(unsafe.Pointer(rfItemC.ranFunctionDefinition.buf))))
	binary.LittleEndian.PutUint64(choiceC[16:], uint64(rfItemC.ranFunctionDefinition.size))
	// Gap of 24 for the asn_struct_ctx_t belonging to OCTET STRING
	binary.LittleEndian.PutUint64(choiceC[48:], uint64(rfItemC.ranFunctionRevision))
	binary.LittleEndian.PutUint64(choiceC[56:], uint64(uintptr(unsafe.Pointer(rfItemC.ranFunctionOID))))

	rfItemIesC := C.RANfunction_ItemIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RANfunction_ItemIEs__value{
			present: C.RANfunction_ItemIEs__value_PR_RANfunction_Item,
			choice:  choiceC,
		},
	}

	return &rfItemIesC, nil
}

func newRANfunctionIDItemIEs(rfIDItemIes *e2appducontents.RanfunctionIdItemIes) (*C.RANfunctionID_ItemIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rfIDItemIes.GetRanFunctionIdItemIes6().GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionIDItem)
	if err != nil {
		return nil, err
	}

	choiceC := [40]byte{} // The size of the RANfunction_ItemIEs__value_u
	rfIDItemC := newRanFunctionIDItem(rfIDItemIes.GetRanFunctionIdItemIes6().GetValue())
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(rfIDItemC.ranFunctionID))
	binary.LittleEndian.PutUint64(choiceC[8:16], uint64(rfIDItemC.ranFunctionRevision))

	rfItemIesC := C.RANfunctionID_ItemIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RANfunctionID_ItemIEs__value{
			present: C.RANfunctionID_ItemIEs__value_PR_RANfunctionID_Item,
			choice:  choiceC,
		},
	}

	return &rfItemIesC, nil
}

func newRANfunctionIDCauseItemIEs(rfIDItemIes *e2appducontents.RanfunctionIdcauseItemIes) (*C.RANfunctionIDcause_ItemIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(rfIDItemIes.GetRanFunctionIdcauseItemIes7().GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRanfunctionIeCauseItem)
	if err != nil {
		return nil, err
	}

	choiceC := [72]byte{} // The size of the RANfunction_ItemIEs__value_u
	rfIDItemC, err := newRanFunctionIDCauseItem(rfIDItemIes.GetRanFunctionIdcauseItemIes7().GetValue())
	if err != nil {
		return nil, fmt.Errorf("newRanFunctionIDCauseItem() error %s", err.Error())
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(rfIDItemC.ranFunctionID))
	binary.LittleEndian.PutUint64(choiceC[8:16], uint64(rfIDItemC.cause.present))
	copy(choiceC[16:24], rfIDItemC.cause.choice[:])

	rfItemIesC := C.RANfunctionIDcause_ItemIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RANfunctionIDcause_ItemIEs__value{
			present: C.RANfunctionIDcause_ItemIEs__value_PR_RANfunctionIDcause_Item,
			choice:  choiceC,
		},
	}

	return &rfItemIesC, nil
}

func newRicActionAdmittedItemIEs(raaItemIes *e2appducontents.RicactionAdmittedItemIes) (*C.RICaction_Admitted_ItemIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(raaItemIes.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicactionAdmittedItem)
	if err != nil {
		return nil, err
	}

	choiceC := [32]byte{} // The size of the RANfunction_ItemIEs__value_u
	rfItemC := newRicActionAdmittedItem(raaItemIes.GetValue())
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(rfItemC.ricActionID))

	rfItemIesC := C.RICaction_Admitted_ItemIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICaction_Admitted_ItemIEs__value{
			present: C.RICaction_Admitted_ItemIEs__value_PR_RICaction_Admitted_Item,
			choice:  choiceC,
		},
	}

	return &rfItemIesC, nil
}

func newRicActionNotAdmittedItemIEs(ranaItemIes *e2appducontents.RicactionNotAdmittedItemIes) (*C.RICaction_NotAdmitted_ItemIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(ranaItemIes.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicactionNotAdmittedItem)
	if err != nil {
		return nil, err
	}

	choiceC := [72]byte{} // The size of the RANfunction_ItemIEs__value_u
	rfItemC, err := newRicActionNotAdmittedItem(ranaItemIes.GetValue())
	if err != nil {
		return nil, fmt.Errorf("newRicActionNotAdmittedItem() %s", err.Error())
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(rfItemC.ricActionID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(rfItemC.cause.present))
	copy(choiceC[16:24], rfItemC.cause.choice[:])

	rfItemIesC := C.RICaction_NotAdmitted_ItemIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICaction_NotAdmitted_ItemIEs__value{
			present: C.RICaction_NotAdmitted_ItemIEs__value_PR_RICaction_NotAdmitted_Item,
			choice:  choiceC,
		},
	}

	return &rfItemIesC, nil
}

func newRicActionToBeSetupItemIEs(ratbsItemIes *e2appducontents.RicactionToBeSetupItemIes) (*C.RICaction_ToBeSetup_ItemIEs_t, error) {
	critC, err := criticalityToC(e2ap_commondatatypes.Criticality(ratbsItemIes.GetCriticality()))
	if err != nil {
		return nil, err
	}
	idC, err := protocolIeIDToC(v1beta2.ProtocolIeIDRicactionToBeSetupItem)
	if err != nil {
		return nil, err
	}

	choiceC := [56]byte{} // The size of the RANfunction_ItemIEs__value_u
	ratbsItemC, err := newRicActionToBeSetupItem(ratbsItemIes.GetValue())
	if err != nil {
		return nil, err
	}
	binary.LittleEndian.PutUint64(choiceC[0:], uint64(ratbsItemC.ricActionID))
	binary.LittleEndian.PutUint64(choiceC[8:], uint64(ratbsItemC.ricActionType))
	binary.LittleEndian.PutUint64(choiceC[16:], uint64(uintptr(unsafe.Pointer(ratbsItemC.ricActionDefinition))))
	binary.LittleEndian.PutUint64(choiceC[24:], uint64(uintptr(unsafe.Pointer(ratbsItemC.ricSubsequentAction))))

	rfItemIesC := C.RICaction_ToBeSetup_ItemIEs_t{
		id:          idC,
		criticality: critC,
		value: C.struct_RICaction_ToBeSetup_ItemIEs__value{
			present: C.RICaction_ToBeSetup_ItemIEs__value_PR_RICaction_ToBeSetup_Item,
			choice:  choiceC,
		},
	}

	return &rfItemIesC, nil
}

func decodeE2setupRequestIE(e2srIeC *C.E2setupRequestIEs_t) (*e2appducontents.E2SetupRequestIes, error) {
	//fmt.Printf("Handling E2SetupReqIE %+v\n", e2srIeC)
	ret := new(e2appducontents.E2SetupRequestIes)

	switch e2srIeC.value.present {
	case C.E2setupRequestIEs__value_PR_GlobalE2node_ID:
		gE2nID, err := decodeGlobalE2NodeID(e2srIeC.value.choice)
		if err != nil {
			return nil, err
		}
		ret.E2ApProtocolIes3 = &e2appducontents.E2SetupRequestIes_E2SetupRequestIes3{
			Id:          int32(v1beta2.ProtocolIeIDGlobalE2nodeID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       gE2nID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.E2setupRequestIEs__value_PR_RANfunctions_List:
		rfl, err := decodeRanFunctionsListBytes(e2srIeC.value.choice)
		if err != nil {
			return nil, err
		}
		ret.E2ApProtocolIes10 = &e2appducontents.E2SetupRequestIes_E2SetupRequestIes10{
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionsAdded),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rfl,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}
	case C.E2setupRequestIEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeE2setupRequestIE(). %v not yet implemneted", e2srIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeE2setupRequestIE(). unexpected choice %v", e2srIeC.value.present)
	}

	return ret, nil
}

func decodeE2setupResponseIE(e2srIeC *C.E2setupResponseIEs_t) (*e2appducontents.E2SetupResponseIes, error) {
	//fmt.Printf("Handling E2SetupReqIE %+v\n", e2srIeC)
	ret := new(e2appducontents.E2SetupResponseIes)

	switch e2srIeC.value.present {
	case C.E2setupResponseIEs__value_PR_GlobalRIC_ID:
		gE2nID, err := decodeGlobalRicIDBytes(e2srIeC.value.choice)
		if err != nil {
			return nil, err
		}
		ret.E2ApProtocolIes4 = &e2appducontents.E2SetupResponseIes_E2SetupResponseIes4{
			Id:          int32(v1beta2.ProtocolIeIDGlobalRicID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       gE2nID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}
	case C.E2setupResponseIEs__value_PR_RANfunctionsID_List:
		rfAccepted, err := decodeRanFunctionsIDListBytes(e2srIeC.value.choice)
		if err != nil {
			return nil, err
		}
		ret.E2ApProtocolIes9 = &e2appducontents.E2SetupResponseIes_E2SetupResponseIes9{
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionsAccepted),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rfAccepted,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}
	case C.E2setupResponseIEs__value_PR_RANfunctionsIDcause_List:
		rfRejected, err := decodeRanFunctionsIDCauseListBytes(e2srIeC.value.choice)
		if err != nil {
			return nil, err
		}
		ret.E2ApProtocolIes13 = &e2appducontents.E2SetupResponseIes_E2SetupResponseIes13{
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionsRejected),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rfRejected,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}
	case C.E2setupResponseIEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeE2setupResponseIE(). %v not yet implemneted", e2srIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeE2setupResponseIE(). unexpected choice %v", e2srIeC.value.present)
	}

	return ret, nil
}

func decodeRicSubscriptionRequestIE(rsrIeC *C.RICsubscriptionRequest_IEs_t) (*e2appducontents.RicsubscriptionRequestIes, error) {
	//	//fmt.Printf("Handling RicSubscriptionResp %+v\n", rsrIeC)
	ret := new(e2appducontents.RicsubscriptionRequestIes)
	//
	switch rsrIeC.value.present {
	case C.RICsubscriptionRequest_IEs__value_PR_RICrequestID:
		ret.E2ApProtocolIes29 = &e2appducontents.RicsubscriptionRequestIes_RicsubscriptionRequestIes29{
			Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
			Value:       decodeRicRequestIDBytes(rsrIeC.value.choice[:16]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}
	case C.RICsubscriptionRequest_IEs__value_PR_RANfunctionID:
		ret.E2ApProtocolIes5 = &e2appducontents.RicsubscriptionRequestIes_RicsubscriptionRequestIes5{
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       decodeRanFunctionIDBytes(rsrIeC.value.choice[0:8]),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}
	case C.RICsubscriptionRequest_IEs__value_PR_RICsubscriptionDetails:
		rsDet, err := decodeRicSubscriptionDetailsBytes(rsrIeC.value.choice[0:64])
		if err != nil {
			return nil, fmt.Errorf("decodeRicSubscriptionDetailsBytes() %s", err.Error())
		}
		ret.E2ApProtocolIes30 = &e2appducontents.RicsubscriptionRequestIes_RicsubscriptionRequestIes30{
			Id:          int32(v1beta2.ProtocolIeIDRicsubscriptionDetails),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rsDet,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}
	case C.RICsubscriptionRequest_IEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeRicSubscriptionRequestIE(). %v not yet implemneted", rsrIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeRicSubscriptionRequestIE(). unexpected choice %v", rsrIeC.value.present)
	}

	return ret, nil
}

func decodeRicSubscriptionResponseIE(rsrIeC *C.RICsubscriptionResponse_IEs_t) (*e2appducontents.RicsubscriptionResponseIes, error) {
	//fmt.Printf("Handling RicSubscriptionResp %+v\n", rsrIeC)
	ret := new(e2appducontents.RicsubscriptionResponseIes)

	switch rsrIeC.value.present {
	case C.RICsubscriptionResponse_IEs__value_PR_RANfunctionID:
		ret.E2ApProtocolIes5 = &e2appducontents.RicsubscriptionResponseIes_RicsubscriptionResponseIes5{
			Value:       decodeRanFunctionIDBytes(rsrIeC.value.choice[:8]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
		}

	case C.RICsubscriptionResponse_IEs__value_PR_RICrequestID:
		ret.E2ApProtocolIes29 = &e2appducontents.RicsubscriptionResponseIes_RicsubscriptionResponseIes29{
			Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
			Value:       decodeRicRequestIDBytes(rsrIeC.value.choice[:16]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICsubscriptionResponse_IEs__value_PR_RICaction_Admitted_List:
		raal, err := decodeRicActionAdmittedListBytes(rsrIeC.value.choice[:48])
		if err != nil {
			return nil, err
		}
		ret.E2ApProtocolIes17 = &e2appducontents.RicsubscriptionResponseIes_RicsubscriptionResponseIes17{
			Id:          int32(v1beta2.ProtocolIeIDRicactionsAdmitted),
			Value:       raal,
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICsubscriptionResponse_IEs__value_PR_RICaction_NotAdmitted_List:
		ranal, err := decodeRicActionNotAdmittedListBytes(rsrIeC.value.choice[:48])
		if err != nil {
			return nil, err
		}
		ret.E2ApProtocolIes18 = &e2appducontents.RicsubscriptionResponseIes_RicsubscriptionResponseIes18{
			Id:          int32(v1beta2.ProtocolIeIDRicactionsNotAdmitted),
			Value:       ranal,
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.RICsubscriptionResponse_IEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeRicSubscriptionResponseIE(). Unexpected value.\n%v", rsrIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeRicSubscriptionResponseIE(). unexpected choice %v", rsrIeC.value.present)
	}

	return ret, nil
}

func decodeRANfunctionItemIes(rfiIesValC *C.struct_RANfunction_ItemIEs__value) (*e2appducontents.RanfunctionItemIes, error) {
	//fmt.Printf("Value %T %v\n", rfiIesValC, rfiIesValC)

	switch present := rfiIesValC.present; present {
	case C.RANfunction_ItemIEs__value_PR_RANfunction_Item:

		rfiIes := e2appducontents.RanfunctionItemIes{
			E2ApProtocolIes10: &e2appducontents.RanfunctionItemIes_RanfunctionItemIes8{
				Id:          int32(v1beta2.ProtocolIeIDRanfunctionItem),
				Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
				Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
			},
		}
		rfi, err := decodeRanFunctionItemBytes(rfiIesValC.choice)
		if err != nil {
			return nil, fmt.Errorf("decodeRANfunctionItemIes() %s", err.Error())
		}
		rfiIes.GetE2ApProtocolIes10().Value = rfi
		return &rfiIes, nil
	default:
		return nil, fmt.Errorf("error decoding RanFunctionItemIE - present %v not supported", present)
	}
}

func decodeRANfunctionIDItemIes(rfIDiIesValC *C.struct_RANfunctionID_ItemIEs__value) (*e2appducontents.RanfunctionIdItemIes, error) {
	//fmt.Printf("Value %T %v\n", rfIDiIesValC, rfIDiIesValC)

	switch present := rfIDiIesValC.present; present {
	case C.RANfunctionID_ItemIEs__value_PR_RANfunctionID_Item:

		rfIDiIes := e2appducontents.RanfunctionIdItemIes{
			RanFunctionIdItemIes6: &e2appducontents.RanfunctionIdItemIes_RanfunctionIdItemIes6{
				Id:          int32(v1beta2.ProtocolIeIDRanfunctionIDItem),
				Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
				Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
			},
		}
		rfi, err := decodeRanFunctionIDItemBytes(rfIDiIesValC.choice)
		if err != nil {
			return nil, fmt.Errorf("decodeRANfunctionIdItemIes() %s", err.Error())
		}
		rfIDiIes.GetRanFunctionIdItemIes6().Value = rfi
		return &rfIDiIes, nil
	default:
		return nil, fmt.Errorf("error decoding RanFunctionIDItemIE - present %v not supported", present)
	}
}

func decodeRANfunctionIDCauseItemIes(rfIDciIesValC *C.struct_RANfunctionIDcause_ItemIEs__value) (*e2appducontents.RanfunctionIdcauseItemIes, error) {
	//fmt.Printf("Value %T %v\n", rfIDciIesValC, rfIDciIesValC)

	switch present := rfIDciIesValC.present; present {
	case C.RANfunctionIDcause_ItemIEs__value_PR_RANfunctionIDcause_Item:

		rfIDiIes := e2appducontents.RanfunctionIdcauseItemIes{
			RanFunctionIdcauseItemIes7: &e2appducontents.RanfunctionIdcauseItemIes_RanfunctionIdcauseItemIes7{
				Id:          int32(v1beta2.ProtocolIeIDRanfunctionIeCauseItem),
				Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
				Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
			},
		}
		rfi, err := decodeRanFunctionIDcauseItemBytes(rfIDciIesValC.choice)
		if err != nil {
			return nil, fmt.Errorf("decodeRANfunctionIdcauseItemIes() %s", err.Error())
		}
		rfIDiIes.GetRanFunctionIdcauseItemIes7().Value = rfi
		return &rfIDiIes, nil
	default:
		return nil, fmt.Errorf("error decoding RanFunctionIDCauseItemIE - present %v not supported", present)
	}
}

func decodeRicActionAdmittedIDItemIes(raaiIesValC *C.struct_RICaction_Admitted_ItemIEs__value) (*e2appducontents.RicactionAdmittedItemIes, error) {
	//fmt.Printf("Value %T %v\n", raaiIesValC, raaiIesValC)

	switch present := raaiIesValC.present; present {
	case C.RICaction_Admitted_ItemIEs__value_PR_RICaction_Admitted_Item:

		raaiIes := e2appducontents.RicactionAdmittedItemIes{
			Id:          int32(v1beta2.ProtocolIeIDRicactionAdmittedItem),
			Value:       decodeRicActionAdmittedItemBytes(raaiIesValC.choice),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}
		return &raaiIes, nil
	default:
		return nil, fmt.Errorf("error decoding RicactionAdmittedItemIes - present %v. not supported", present)
	}
}

func decodeRicActionNotAdmittedIDItemIes(ranaiIesValC *C.struct_RICaction_NotAdmitted_ItemIEs__value) (*e2appducontents.RicactionNotAdmittedItemIes, error) {
	//fmt.Printf("Value %T %v\n", ranaiIesValC, ranaiIesValC)

	switch present := ranaiIesValC.present; present {
	case C.RICaction_NotAdmitted_ItemIEs__value_PR_RICaction_NotAdmitted_Item:
		rana, err := decodeRicActionNotAdmittedItemBytes(ranaiIesValC.choice[:24])
		if err != nil {
			return nil, fmt.Errorf("decodeRicActionNotAdmittedItemBytes() %s", err.Error())
		}
		ranaiIes := e2appducontents.RicactionNotAdmittedItemIes{
			Id:          int32(v1beta2.ProtocolIeIDRicactionNotAdmittedItem),
			Value:       rana,
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}
		return &ranaiIes, nil
	default:
		return nil, fmt.Errorf("error decoding RicactionNotAdmittedItemIes - present %v. not supported", present)
	}
}

func decodeRicIndicationIE(riIeC *C.RICindication_IEs_t) (*e2appducontents.RicindicationIes, error) {
	//fmt.Printf("Handling E2SetupReqIE %+v\n", riIeC)
	ret := new(e2appducontents.RicindicationIes)

	switch riIeC.value.present {
	case C.RICindication_IEs__value_PR_RANfunctionID:
		rfID := decodeRanFunctionIDBytes(riIeC.value.choice[0:8])
		ret.E2ApProtocolIes5 = &e2appducontents.RicindicationIes_RicindicationIes5{
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rfID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICindication_IEs__value_PR_RICactionID:
		raID := decodeRicActionIDBytes(riIeC.value.choice[0:8])
		ret.E2ApProtocolIes15 = &e2appducontents.RicindicationIes_RicindicationIes15{
			Id:          int32(v1beta2.ProtocolIeIDRicactionID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       raID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICindication_IEs__value_PR_RICcallProcessID:
		rcpID := decodeRicCallProcessIDBytes(riIeC.value.choice[0:16])
		ret.E2ApProtocolIes20 = &e2appducontents.RicindicationIes_RicindicationIes20{
			Id:          int32(v1beta2.ProtocolIeIDRiccallProcessID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rcpID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.RICindication_IEs__value_PR_RICindicationHeader:
		rih := decodeRicIndicationHeaderBytes(riIeC.value.choice[0:16])
		ret.E2ApProtocolIes25 = &e2appducontents.RicindicationIes_RicindicationIes25{
			Id:          int32(v1beta2.ProtocolIeIDRicindicationHeader),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rih,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICindication_IEs__value_PR_RICindicationMessage:
		rim := decodeRicIndicationMessageBytes(riIeC.value.choice[0:16])
		ret.E2ApProtocolIes26 = &e2appducontents.RicindicationIes_RicindicationIes26{
			Id:          int32(v1beta2.ProtocolIeIDRicindicationMessage),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rim,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICindication_IEs__value_PR_RICindicationSN:
		risn := decodeRicIndicationSnBytes(riIeC.value.choice[0:8])
		ret.E2ApProtocolIes27 = &e2appducontents.RicindicationIes_RicindicationIes27{
			Id:          int32(v1beta2.ProtocolIeIDRicindicationSn),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       risn,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.RICindication_IEs__value_PR_RICindicationType:
		rit := decodeRicIndicationTypeBytes(riIeC.value.choice[0:8])
		ret.E2ApProtocolIes28 = &e2appducontents.RicindicationIes_RicindicationIes28{
			Id:          int32(v1beta2.ProtocolIeIDRicindicationType),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rit,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICindication_IEs__value_PR_RICrequestID:
		rrID := decodeRicRequestIDBytes(riIeC.value.choice[0:16])
		ret.E2ApProtocolIes29 = &e2appducontents.RicindicationIes_RicindicationIes29{
			Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rrID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICindication_IEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeRicIndicationIE(). %v not yet implemneted", riIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeRicIndicationIE(). unexpected choice %v", riIeC.value.present)
	}

	return ret, nil
}

func decodeRicControlRequestIE(rcRIeC *C.RICcontrolRequest_IEs_t) (*e2appducontents.RiccontrolRequestIes, error) {
	//fmt.Printf("Handling E2SetupReqIE %+v\n", riIeC)
	ret := new(e2appducontents.RiccontrolRequestIes)

	switch rcRIeC.value.present {
	case C.RICcontrolRequest_IEs__value_PR_RANfunctionID:
		rfID := decodeRanFunctionIDBytes(rcRIeC.value.choice[0:8])
		ret.E2ApProtocolIes5 = &e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes5{
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rfID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICcontrolRequest_IEs__value_PR_RICcallProcessID:
		rcpID := decodeRicCallProcessIDBytes(rcRIeC.value.choice[0:16])
		ret.E2ApProtocolIes20 = &e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes20{
			Id:          int32(v1beta2.ProtocolIeIDRiccallProcessID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rcpID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.RICcontrolRequest_IEs__value_PR_RICrequestID:
		rrID := decodeRicRequestIDBytes(rcRIeC.value.choice[0:16])
		ret.E2ApProtocolIes29 = &e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes29{
			Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rrID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICcontrolRequest_IEs__value_PR_RICcontrolHeader:
		rch := decodeRicControlHeaderBytes(rcRIeC.value.choice[0:16])
		ret.E2ApProtocolIes22 = &e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes22{
			Id:          int32(v1beta2.ProtocolIeIDRiccontrolHeader),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rch,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICcontrolRequest_IEs__value_PR_RICcontrolMessage:
		rcm := decodeRicControlMessageBytes(rcRIeC.value.choice[0:16])
		ret.E2ApProtocolIes23 = &e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes23{
			Id:          int32(v1beta2.ProtocolIeIDRiccontrolMessage),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rcm,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICcontrolRequest_IEs__value_PR_RICcontrolAckRequest:
		rcar := decodeRicControlAckRequestBytes(rcRIeC.value.choice[0:16])
		ret.E2ApProtocolIes21 = &e2appducontents.RiccontrolRequestIes_RiccontrolRequestIes21{
			Id:          int32(v1beta2.ProtocolIeIDRiccontrolAckRequest),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rcar,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.RICcontrolRequest_IEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeRicControlRequestIE(). %v not yet implemneted", rcRIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeRicControlRequestIE(). unexpected choice %v", rcRIeC.value.present)
	}

	return ret, nil
}

func decodeRicControlAcknowledgeIE(rcaIeC *C.RICcontrolAcknowledge_IEs_t) (*e2appducontents.RiccontrolAcknowledgeIes, error) {
	//fmt.Printf("Handling E2SetupReqIE %+v\n", riIeC)
	ret := new(e2appducontents.RiccontrolAcknowledgeIes)

	switch rcaIeC.value.present {
	case C.RICcontrolAcknowledge_IEs__value_PR_RANfunctionID:
		rfID := decodeRanFunctionIDBytes(rcaIeC.value.choice[0:8])
		ret.E2ApProtocolIes5 = &e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes5{
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rfID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICcontrolAcknowledge_IEs__value_PR_RICcallProcessID:
		rcpID := decodeRicCallProcessIDBytes(rcaIeC.value.choice[0:16])
		ret.E2ApProtocolIes20 = &e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes20{
			Id:          int32(v1beta2.ProtocolIeIDRiccallProcessID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rcpID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.RICcontrolAcknowledge_IEs__value_PR_RICcontrolStatus:
		rcs := decodeRicControlStatusBytes(rcaIeC.value.choice[0:16])
		ret.E2ApProtocolIes24 = &e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes24{
			Id:          int32(v1beta2.ProtocolIeIDRiccontrolStatus),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rcs,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICcontrolAcknowledge_IEs__value_PR_RICrequestID:
		rrID := decodeRicRequestIDBytes(rcaIeC.value.choice[0:16])
		ret.E2ApProtocolIes29 = &e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes29{
			Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rrID,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICcontrolAcknowledge_IEs__value_PR_RICcontrolOutcome:
		rco := decodeRicControlOutcomeBytes(rcaIeC.value.choice[0:16])
		ret.E2ApProtocolIes32 = &e2appducontents.RiccontrolAcknowledgeIes_RiccontrolAcknowledgeIes32{
			Id:          int32(v1beta2.ProtocolIeIDRiccontrolOutcome),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       rco,
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.RICcontrolAcknowledge_IEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeRicControlAcknowledgeIE(). %v not yet implemneted", rcaIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeRicControlAcknowledgeIE(). unexpected choice %v", rcaIeC.value.present)
	}

	return ret, nil
}

func decodeRicActionToBeSetupItemIes(ratbsIesValC *C.struct_RICaction_ToBeSetup_ItemIEs__value) (*e2appducontents.RicactionToBeSetupItemIes, error) {
	//fmt.Printf("Value %T %v\n", ratbsIesValC, ratbsIesValC)

	switch present := ratbsIesValC.present; present {
	case C.RICaction_ToBeSetup_ItemIEs__value_PR_RICaction_ToBeSetup_Item:
		ratbsIIes := e2appducontents.RicactionToBeSetupItemIes{
			Id:          int32(v1beta2.ProtocolIeIDRicactionToBeSetupItem),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}
		ratbsI, err := decodeRicActionToBeSetupItemBytes(ratbsIesValC.choice)
		if err != nil {
			return nil, fmt.Errorf("decodeRicActionToBeSetupItemBytes() %s", err.Error())
		}
		ratbsIIes.Value = ratbsI
		return &ratbsIIes, nil
	default:
		return nil, fmt.Errorf("error decoding RicactionToBeSetupItemIes - present %v not supported", present)
	}
}

func decodeRicSubscriptionDeleteRequestIE(rsrdIeC *C.RICsubscriptionDeleteRequest_IEs_t) (*e2appducontents.RicsubscriptionDeleteRequestIes, error) {
	//	//fmt.Printf("Handling RicSubscriptionResp %+v\n", rsrdIeC)
	ret := new(e2appducontents.RicsubscriptionDeleteRequestIes)
	//
	switch rsrdIeC.value.present {
	case C.RICsubscriptionDeleteRequest_IEs__value_PR_RICrequestID:
		ret.E2ApProtocolIes29 = &e2appducontents.RicsubscriptionDeleteRequestIes_RicsubscriptionDeleteRequestIes29{
			Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
			Value:       decodeRicRequestIDBytes(rsrdIeC.value.choice[:16]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}
	case C.RICsubscriptionDeleteRequest_IEs__value_PR_RANfunctionID:
		ret.E2ApProtocolIes5 = &e2appducontents.RicsubscriptionDeleteRequestIes_RicsubscriptionDeleteRequestIes5{
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Value:       decodeRanFunctionIDBytes(rsrdIeC.value.choice[0:8]),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}
	case C.RICsubscriptionDeleteRequest_IEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeRicSubscriptionDeleteRequestIE(). %v not yet implemneted", rsrdIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeRicSubscriptionDeleteRequestIE(). unexpected choice %v", rsrdIeC.value.present)
	}

	return ret, nil
}

func decodeRicSubscriptionDeleteResponseIE(rsdrIeC *C.RICsubscriptionDeleteResponse_IEs_t) (*e2appducontents.RicsubscriptionDeleteResponseIes, error) {
	//fmt.Printf("Handling RicSubscriptionResp %+v\n", rsdrIeC)
	ret := new(e2appducontents.RicsubscriptionDeleteResponseIes)

	switch rsdrIeC.value.present {
	case C.RICsubscriptionDeleteResponse_IEs__value_PR_RANfunctionID:
		ret.E2ApProtocolIes5 = &e2appducontents.RicsubscriptionDeleteResponseIes_RicsubscriptionDeleteResponseIes5{
			Value:       decodeRanFunctionIDBytes(rsdrIeC.value.choice[:8]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
		}

	case C.RICsubscriptionDeleteResponse_IEs__value_PR_RICrequestID:
		ret.E2ApProtocolIes29 = &e2appducontents.RicsubscriptionDeleteResponseIes_RicsubscriptionDeleteResponseIes29{
			Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
			Value:       decodeRicRequestIDBytes(rsdrIeC.value.choice[:16]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICsubscriptionDeleteResponse_IEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeRicSubscriptionDeleteResponseIE(). %v not yet implemneted", rsdrIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeRicSubscriptionDeleteResponseIE(). unexpected choice %v", rsdrIeC.value.present)
	}

	return ret, nil
}

func decodeRicSubscriptionDeleteFailureIE(rsdfIeC *C.RICsubscriptionDeleteFailure_IEs_t) (*e2appducontents.RicsubscriptionDeleteFailureIes, error) {
	//fmt.Printf("Handling RicSubscriptionResp %+v\n", rsdfIeC)
	ret := new(e2appducontents.RicsubscriptionDeleteFailureIes)

	switch rsdfIeC.value.present {
	case C.RICsubscriptionDeleteFailure_IEs__value_PR_RANfunctionID:
		ret.E2ApProtocolIes5 = &e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes5{
			Value:       decodeRanFunctionIDBytes(rsdfIeC.value.choice[:8]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
		}

	case C.RICsubscriptionDeleteFailure_IEs__value_PR_RICrequestID:
		ret.E2ApProtocolIes29 = &e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes29{
			Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
			Value:       decodeRicRequestIDBytes(rsdfIeC.value.choice[:16]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICsubscriptionDeleteFailure_IEs__value_PR_Cause:
		cause, err := decodeCauseBytes(rsdfIeC.value.choice[:16])
		if err != nil {
			return nil, fmt.Errorf("decodeCauseBytes() %s", err.Error())
		}
		ret.E2ApProtocolIes1 = &e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes1{
			Id:          int32(v1beta2.ProtocolIeIDCause),
			Value:       cause,
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICsubscriptionDeleteFailure_IEs__value_PR_CriticalityDiagnostics:
		cd, err := decodeCriticalityDiagnosticsBytes(rsdfIeC.value.choice[:48])
		if err != nil {
			return nil, fmt.Errorf("decodeCriticalityDiagnosticsBytes() %s", err.Error())
		}
		ret.E2ApProtocolIes2 = &e2appducontents.RicsubscriptionDeleteFailureIes_RicsubscriptionDeleteFailureIes2{
			Id:          int32(v1beta2.ProtocolIeIDCriticalityDiagnostics),
			Value:       cd,
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.RICsubscriptionDeleteFailure_IEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeRicSubscriptionDeleteFailureIE(). %v not yet implemneted", rsdfIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeRicSubscriptionDeleteFailureIE(). unexpected choice %v", rsdfIeC.value.present)
	}

	return ret, nil
}

func decodeRicSubscriptionFailureIE(rsfIeC *C.RICsubscriptionFailure_IEs_t) (*e2appducontents.RicsubscriptionFailureIes, error) {
	//fmt.Printf("Handling RicSubscriptionResp %+v\n", rsfIeC)
	ret := new(e2appducontents.RicsubscriptionFailureIes)

	switch rsfIeC.value.present {
	case C.RICsubscriptionFailure_IEs__value_PR_RANfunctionID:
		ret.E2ApProtocolIes5 = &e2appducontents.RicsubscriptionFailureIes_RicsubscriptionFailureIes5{
			Value:       decodeRanFunctionIDBytes(rsfIeC.value.choice[:8]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
		}

	case C.RICsubscriptionFailure_IEs__value_PR_RICrequestID:
		ret.E2ApProtocolIes29 = &e2appducontents.RicsubscriptionFailureIes_RicsubscriptionFailureIes29{
			Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
			Value:       decodeRicRequestIDBytes(rsfIeC.value.choice[:16]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICsubscriptionFailure_IEs__value_PR_CriticalityDiagnostics:
		cd, err := decodeCriticalityDiagnosticsBytes(rsfIeC.value.choice[:48])
		if err != nil {
			return nil, fmt.Errorf("decodeCriticalityDiagnosticsBytes() %s", err.Error())
		}
		ret.E2ApProtocolIes2 = &e2appducontents.RicsubscriptionFailureIes_RicsubscriptionFailureIes2{
			Id:          int32(v1beta2.ProtocolIeIDCriticalityDiagnostics),
			Value:       cd,
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.RICsubscriptionFailure_IEs__value_PR_RICaction_NotAdmitted_List:
		ranaL, err := decodeRicActionNotAdmittedListBytes(rsfIeC.value.choice[:48])
		if err != nil {
			return nil, fmt.Errorf("decodeRicActionNotAdmittedListBytes() %s", err.Error())
		}
		ret.E2ApProtocolIes18 = &e2appducontents.RicsubscriptionFailureIes_RicsubscriptionFailureIes18{
			Id:          int32(v1beta2.ProtocolIeIDRicactionsNotAdmitted),
			Value:       ranaL,
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.RICsubscriptionFailure_IEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeRicSubscriptionFailureIE(). %v not yet implemneted", rsfIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeRicSubscriptionFailureIE(). unexpected choice %v", rsfIeC.value.present)
	}

	return ret, nil
}

func decodeErrorIndicationIE(eiIeC *C.ErrorIndication_IEs_t) (*e2appducontents.ErrorIndicationIes, error) {
	//fmt.Printf("Handling ErrorIndication %+v\n", rsfIeC)
	ret := new(e2appducontents.ErrorIndicationIes)

	switch eiIeC.value.present {
	case C.ErrorIndication_IEs__value_PR_RANfunctionID:
		ret.E2ApProtocolIes5 = &e2appducontents.ErrorIndicationIes_ErrorIndicationIes5{
			Value:       decodeRanFunctionIDBytes(eiIeC.value.choice[:8]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
			Id:          int32(v1beta2.ProtocolIeIDRanfunctionID),
		}

	case C.ErrorIndication_IEs__value_PR_RICrequestID:
		ret.E2ApProtocolIes29 = &e2appducontents.ErrorIndicationIes_ErrorIndicationIes29{
			Id:          int32(v1beta2.ProtocolIeIDRicrequestID),
			Value:       decodeRicRequestIDBytes(eiIeC.value.choice[:16]),
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_REJECT),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.ErrorIndication_IEs__value_PR_CriticalityDiagnostics:
		cd, err := decodeCriticalityDiagnosticsBytes(eiIeC.value.choice[:48])
		if err != nil {
			return nil, fmt.Errorf("decodeCriticalityDiagnosticsBytes() %s", err.Error())
		}
		ret.E2ApProtocolIes2 = &e2appducontents.ErrorIndicationIes_ErrorIndicationIes2{
			Id:          int32(v1beta2.ProtocolIeIDCriticalityDiagnostics),
			Value:       cd,
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.ErrorIndication_IEs__value_PR_Cause:
		cause, err := decodeCauseBytes(eiIeC.value.choice[:48])
		if err != nil {
			return nil, fmt.Errorf("decodeCauseBytes() %s", err.Error())
		}
		ret.E2ApProtocolIes1 = &e2appducontents.ErrorIndicationIes_ErrorIndicationIes1{
			Id:          int32(v1beta2.ProtocolIeIDCause),
			Value:       cause,
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.ErrorIndication_IEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeErrorIndicationIE(). %v not yet implemneted", eiIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeErrorIndicationIE(). unexpected choice %v", eiIeC.value.present)
	}

	return ret, nil
}

func decodeE2setupFailureIE(eiIeC *C.E2setupFailureIEs_t) (*e2appducontents.E2SetupFailureIes, error) {
	//fmt.Printf("Handling ErrorIndication %+v\n", rsfIeC)
	ret := new(e2appducontents.E2SetupFailureIes)

	switch eiIeC.value.present {
	case C.E2setupFailureIEs__value_PR_TimeToWait:
		ret.E2ApProtocolIes31 = &e2appducontents.E2SetupFailureIes_E2SetupFailureIes31{
			Value:       decodeTimeToWaitBytes(eiIeC.value.choice[:8]), //TODO: See RICtimeToWait
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
			Id:          int32(v1beta2.ProtocolIeIDTimeToWait),
		}

	case C.E2setupFailureIEs__value_PR_CriticalityDiagnostics:
		cd, err := decodeCriticalityDiagnosticsBytes(eiIeC.value.choice[:48])
		if err != nil {
			return nil, fmt.Errorf("decodeCriticalityDiagnosticsBytes() %s", err.Error())
		}
		ret.E2ApProtocolIes2 = &e2appducontents.E2SetupFailureIes_E2SetupFailureIes2{
			Id:          int32(v1beta2.ProtocolIeIDCriticalityDiagnostics),
			Value:       cd,
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_OPTIONAL),
		}

	case C.E2setupFailureIEs__value_PR_Cause:
		cause, err := decodeCauseBytes(eiIeC.value.choice[:48])
		if err != nil {
			return nil, fmt.Errorf("decodeCauseBytes() %s", err.Error())
		}
		ret.E2ApProtocolIes1 = &e2appducontents.E2SetupFailureIes_E2SetupFailureIes1{
			Id:          int32(v1beta2.ProtocolIeIDCause),
			Value:       cause,
			Criticality: int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE),
			Presence:    int32(e2ap_commondatatypes.Presence_PRESENCE_MANDATORY),
		}

	case C.E2setupFailureIEs__value_PR_NOTHING:
		return nil, fmt.Errorf("decodeErrorIndicationIE(). %v not yet implemneted", eiIeC.value.present)

	default:
		return nil, fmt.Errorf("decodeErrorIndicationIE(). unexpected choice %v", eiIeC.value.present)
	}

	return ret, nil
}
