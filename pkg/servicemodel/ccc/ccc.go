// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package ccc

import (
	"context"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/pdubuilder"
	e2smcccsm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/servicemodel"
	e2smccc "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/v1/e2sm-ccc-ies"
	e2smcommon "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_ccc/v1/e2sm-common-ies"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	e2apIndicationUtils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/indication"
	ranfuncdescription "github.com/onosproject/ran-simulator/pkg/utils/e2sm/ccc/ranfunctiondefinition"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/ccc/controloutcome"
	cccIndicationHeader "github.com/onosproject/ran-simulator/pkg/utils/e2sm/ccc/indication"
	cccMessageFormat1 "github.com/onosproject/ran-simulator/pkg/utils/e2sm/ccc/indication/messageformat1"

	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
	e2apControlutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/control"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"
	"google.golang.org/protobuf/proto"
)

var _ servicemodel.Client = &Client{}

var log = logging.GetLogger()

const (
	modelVersion                             = "v1"
	ranConfStructureName                     = "O-RRMPolicyRatio"
	attName                                  = "PolicyRatio_1"
	eventTriggerStyleType                    = 1
	eventTriggerStyleName                    = "E2 Node Configuration Change"
	eventTriggerFormatType                   = 1
	reportServiceStyleType                   = 1
	reportServiceStyleName                   = "Node-Level Configuration"
	reportServiceActionDefinitionFormatType  = 1
	reportServiceIndicationHeaderFormatType  = 1
	reportServiceIndicationMessageFormatType = 1
	controlServiceStyleType                  = 1
	controlServiceStyleName                  = "Node Configuration and Control"
	controlServiceHeaderFormatType           = 1
	controlServiceMessageFormatType          = 1
	controlServiceControlOutcomeFormatType   = 1
	ranFunctionDescription                   = "Cell Configuration and Control"
	ranFunctionShortName                     = "ORAN-E2SM-CCC"
	ranFunctionE2SmOid                       = "1.3.6.1.4.1.53148.1.1.2.4"
	ranFunctionInstance                      = 1
)

// Client ccc service model client
type Client struct {
	ServiceModel *registry.ServiceModel
}

// E2ConnectionUpdate implements connection update procedure
func (sm *Client) E2ConnectionUpdate(ctx context.Context, request *e2appducontents.E2ConnectionUpdate) (response *e2appducontents.E2ConnectionUpdateAcknowledge, failure *e2appducontents.E2ConnectionUpdateFailure, err error) {
	return nil, nil, errors.NewNotSupported("connection update is not supported")
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model,
	subStore *subscriptions.Subscriptions, nodeStore nodes.Store) (registry.ServiceModel, error) {
	cccSm := registry.ServiceModel{
		RanFunctionID: registry.Ccc,
		ModelName:     ranFunctionShortName,
		Revision:      1,
		OID:           ranFunctionE2SmOid,
		Version:       modelVersion,
		Node:          node,
		Model:         model,
		Subscriptions: subStore,
		Nodes:         nodeStore,
	}
	cccClient := &Client{
		ServiceModel: &cccSm,
	}

	cccSm.Client = cccClient

	// Event trigger
	eventTriggerStyles := make([]*e2smccc.EventTriggerStyle, 0)

	eventTriggerStyle := &e2smccc.EventTriggerStyle{
		EventTriggerStyleType: &e2smcommon.RicStyleType{
			Value: eventTriggerStyleType,
		},
		EventTriggerStyleName: &e2smcommon.RicStyleName{
			Value: eventTriggerStyleName,
		},
		EventTriggerFormatType: &e2smcommon.RicFormatType{
			Value: eventTriggerFormatType,
		},
	}

	eventTriggerStyles = append(eventTriggerStyles, eventTriggerStyle)

	eventTrigger := &e2smccc.EventTrigger{
		ListOfSupportedEventTriggerStyles: &e2smccc.ListOfSupportedEventTriggerStyles{
			Value: eventTriggerStyles,
		},
	}

	// Report service
	eventTriggerStyleTypeList := make([]*e2smccc.EventTriggerStyleType, 0)
	eventTriggerStyleTypeItem := &e2smccc.EventTriggerStyleType{
		EventTriggerStyleType: &e2smcommon.RicStyleType{
			Value: eventTriggerStyleType,
		},
	}

	eventTriggerStyleTypeList = append(eventTriggerStyleTypeList, eventTriggerStyleTypeItem)

	reportStyles := make([]*e2smccc.ReportStyle, 0)

	reportStyle := &e2smccc.ReportStyle{
		ReportServiceStyleType: &e2smcommon.RicStyleType{
			Value: reportServiceStyleType,
		},
		ReportServiceStyleName: &e2smcommon.RicStyleName{
			Value: reportServiceStyleName,
		},
		ListOfSupportedEventTriggerStylesForReportStyle: &e2smccc.ListOfSupportedEventTriggerStylesForReportStyle{
			Value: eventTriggerStyleTypeList,
		},
		ReportServiceActionDefinitionFormatType: &e2smcommon.RicFormatType{
			Value: reportServiceActionDefinitionFormatType,
		},
		ReportServiceIndicationHeaderFormatType: &e2smcommon.RicFormatType{
			Value: reportServiceIndicationHeaderFormatType,
		},
		ReportServiceIndicationMessageFormatType: &e2smcommon.RicFormatType{
			Value: reportServiceIndicationMessageFormatType,
		},
	}

	reportStyles = append(reportStyles, reportStyle)

	reportService := &e2smccc.ReportService{
		ListOfSupportedReportStyles: &e2smccc.ListOfSupportedReportStyles{
			Value: reportStyles,
		},
	}

	// Control service
	controlStyles := make([]*e2smccc.ControlStyle, 0)

	controlStyle := &e2smccc.ControlStyle{
		ControlServiceStyleType: &e2smcommon.RicStyleType{
			Value: controlServiceStyleType,
		},
		ControlServiceStyleName: &e2smcommon.RicStyleName{
			Value: controlServiceStyleName,
		},
		ControlServiceHeaderFormatType: &e2smcommon.RicFormatType{
			Value: controlServiceHeaderFormatType,
		},
		ControlServiceMessageFormatType: &e2smcommon.RicFormatType{
			Value: controlServiceMessageFormatType,
		},
		ControlServiceControlOutcomeFormatType: &e2smcommon.RicFormatType{
			Value: controlServiceControlOutcomeFormatType,
		},
	}

	controlStyles = append(controlStyles, controlStyle)

	controlService := &e2smccc.ControlService{
		ListOfSupportedControlStyles: &e2smccc.ListOfSupportedControlStyles{
			Value: controlStyles,
		},
	}

	log.Debugf("Currently disabling control service messages due to unknown issue: %v", controlService)

	// RIC services
	ricServices := &e2smccc.Ricservices{}
	ricServices.SetEventTrigger(eventTrigger)
	ricServices.SetReportService(reportService)
	// ricServices.SetControlService(controlService)

	attNamePtr, err := pdubuilder.CreateAttributeName([]byte(attName))
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	attribute, err := pdubuilder.CreateAttribute(attNamePtr, ricServices)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	attributes := make([]*e2smccc.Attribute, 0)
	attributes = append(attributes, attribute)

	listOfSupportedAttributes, err := pdubuilder.CreateListOfSupportedAttributes(attributes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	ranConfigurationStructureName, err := pdubuilder.CreateRanConfigurationStructureName([]byte(ranConfStructureName))
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	ranConfigurationStructures := make([]*e2smccc.RanconfigurationStructure, 0)
	ranConfigurationStructure := &e2smccc.RanconfigurationStructure{
		RanConfigurationStructureName: ranConfigurationStructureName,
		ListOfSupportedAttributes:     listOfSupportedAttributes,
	}
	ranConfigurationStructures = append(ranConfigurationStructures, ranConfigurationStructure)
	listOfSupportedRanConfiguratonStructures, err := pdubuilder.CreateListOfSupportedRanconfigurationStructures(ranConfigurationStructures)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	ranFuncDescPdu, err := ranfuncdescription.NewRANFunctionDefinition(
		ranfuncdescription.WithRANFunctionShortName(ranFunctionShortName),
		ranfuncdescription.WithRANFunctionE2SmOID(ranFunctionE2SmOid),
		ranfuncdescription.WithRANFunctionDefinition(ranFunctionDescription),
		ranfuncdescription.WithRANFunctionInstance(ranFunctionInstance),
		ranfuncdescription.WithListOfSupportedRanconfigurationStructures(listOfSupportedRanConfiguratonStructures)).
		Build()

	log.Debug("PDU for CCC service model Ran Function Definition: %+v", ranFuncDescPdu)

	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	protoBytes, err := proto.Marshal(ranFuncDescPdu)
	log.Debug("Proto bytes of CCC service model Ran Function Definition:", protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	var cccServiceModel e2smcccsm.CCCServiceModel
	ranFuncDescBytes, err := cccServiceModel.RanFuncDescriptionProtoToASN1(protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	cccSm.Description = ranFuncDescBytes
	return cccSm, nil
}

func (sm *Client) createIndicationHeaderBytes(reason e2smccc.IndicationReason) ([]byte, error) {
	// Creates an indication header
	timestamp := make([]byte, 8)
	binary.BigEndian.PutUint64(timestamp, uint64(time.Now().Unix()))

	header := cccIndicationHeader.NewIndicationHeader(
		cccIndicationHeader.WithIndicationReason(reason),
		cccIndicationHeader.WithTimeStamp(timestamp))

	indicationHeaderAsn1Bytes, err := header.ToAsn1Bytes()
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return indicationHeaderAsn1Bytes, nil
}

func (sm *Client) createIndicationMsgFormat1(ctx context.Context, changeType e2smccc.ChangeType) ([]byte, error) {
	log.Debug("Create Indication message format 1 based on action defs for %v", changeType)

	// TODO: GA: There should be a loop here if there are multiple slices for a single policy ratio
	plmnID := ransimtypes.NewUint24(uint32(sm.ServiceModel.Model.PlmnID))
	sstDec := byte(1)
	sdDec := ransimtypes.NewUint24(uint32(012345))
	var maxRatio int32 = 80
	var minRatio int32 = 50
	var dedicatedRatio int32 = 10

	sst, err := pdubuilder.CreateSst([]byte{sstDec})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	sd, err := pdubuilder.CreateSd(sdDec.ToBytes())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	plmnId, err := pdubuilder.CreatePlmnidentity(plmnID.ToBytes())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	policyMember := &e2smccc.RrmPolicyMember{
		PlmnId: plmnId,
		Snssai: &e2smcommon.SNSsai{
			SSt: sst,
			SD:  sd,
		},
	}
	log.Warnf("policyMember: %v", policyMember)

	policyMemberVec := make([]*e2smccc.RrmPolicyMember, 0)
	policyMemberVec = append(policyMemberVec, policyMember)
	policyMemeberElement := &e2smccc.RrmPolicyMemberList{
		Value: policyMemberVec,
	}

	oRrmpolicyRatio := &e2smccc.ORRmpolicyRatio{}
	oRrmpolicyRatio.ResourceType = e2smccc.ResourceType_RESOURCE_TYPE_DRB
	oRrmpolicyRatio.SchedulerType = e2smccc.SchedulerType_SCHEDULER_TYPE_ROUND_ROBIN
	oRrmpolicyRatio.RRmpolicyMemberList = policyMemeberElement
	oRrmpolicyRatio.RRmpolicyMaxRatio = maxRatio
	oRrmpolicyRatio.RRmpolicyMinRatio = minRatio
	oRrmpolicyRatio.RRmpolicyDedicatedRatio = dedicatedRatio
	log.Warnf("oRrmpolicyRatio: %v", oRrmpolicyRatio)

	ranConfigurationStructure, err := pdubuilder.CreateRanConfigurationStructureORrmpolicyRatio(oRrmpolicyRatio)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Warnf("ranConfigurationStructure: %v", ranConfigurationStructure)

	valuesOfAttributes, err := pdubuilder.CreateValuesOfAttributes(ranConfigurationStructure)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	ranConfigurationStructureName, err := pdubuilder.CreateRanConfigurationStructureName([]byte(attName))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Warnf("ranConfigurationStructureName: %v", ranConfigurationStructureName)

	listOfConfStructuresReported := make([]*e2smccc.ConfigurationStructure, 0)
	confStruct := &e2smccc.ConfigurationStructure{
		ChangeType:                    changeType,
		RanConfigurationStructureName: ranConfigurationStructureName,
		ValuesOfAttributes:            valuesOfAttributes,
	}
	log.Warnf("confStruct: %v", confStruct)

	listOfConfStructuresReported = append(listOfConfStructuresReported, confStruct)

	listOfConfReported, err := pdubuilder.CreateListOfConfigurationsReported(listOfConfStructuresReported)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	log.Warnf("listOfConfReported: %v", listOfConfReported)

	// Creating an indication message format 1
	indicationMessage := cccMessageFormat1.NewIndicationMessage(
		cccMessageFormat1.WithConfigurationsReported(listOfConfReported))

	indicationMessageBytes, err := indicationMessage.ToAsn1Bytes()
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return indicationMessageBytes, nil
}

// sendRicIndication creates ric indication  for each cell in the node
func (sm *Client) sendRicIndication(ctx context.Context, subscription *subutils.Subscription, actionDefinitions []*e2smccc.E2SmCCcRIcactionDefinition, reason e2smccc.IndicationReason, changeType e2smccc.ChangeType) error {
	// TODO: GA: For now supporting RIC indication message format 1
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}

	// Create CCC RIC Indication Header
	indicationHeaderBytes, err := sm.createIndicationHeaderBytes(reason)
	if err != nil {
		log.Warn(err)
		return err
	}

	// Create CCC RIC Indication message
	for _, actionDefinition := range actionDefinitions {
		format1 := actionDefinition.GetActionDefinitionFormat().GetE2SmCccActionDefinitionFormat1()
		if format1 != nil {
			ranconfigurationStructureForAdfs := format1.GetListOfNodeLevelRanconfigurationStructuresForAdf().Value
			for _, ranconfigurationStructureForAdf := range ranconfigurationStructureForAdfs {
				ranConfigurationStructureName := string(ranconfigurationStructureForAdf.GetRanConfigurationStructureName().Value)
				log.Debug("RanConfigurationStructureName is: %v", ranConfigurationStructureName)

				if ranConfigurationStructureName == ranConfStructureName {
					log.Debug("Sending indication message for RanConfigurationStructureName: %v", ranConfigurationStructureName)
					indicationMessageBytes, err := sm.createIndicationMsgFormat1(ctx, changeType)
					if err != nil {
						return err
					}

					indication := e2apIndicationUtils.NewIndication(
						e2apIndicationUtils.WithRicInstanceID(subscription.GetRicInstanceID()),
						e2apIndicationUtils.WithRanFuncID(subscription.GetRanFuncID()),
						e2apIndicationUtils.WithRequestID(subscription.GetReqID()),
						e2apIndicationUtils.WithIndicationHeader(indicationHeaderBytes),
						e2apIndicationUtils.WithIndicationMessage(indicationMessageBytes))

					ricIndication, err := indication.Build()
					if err != nil {
						log.Error("creating indication message failed for %v and %v due to %v", reason, changeType, err)
						return err
					}

					err = sub.E2Channel.RICIndication(ctx, ricIndication)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// RICControl implements control handler for ccc service model
func (sm *Client) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	log.Infof("Control Request is received for service model %v and e2 node ID: %d", sm.ServiceModel.ModelName, sm.ServiceModel.Node.GnbID)
	reqID, err := e2apControlutils.GetRequesterID(request)
	if err != nil {
		return nil, nil, err
	}
	ranFuncID, err := e2apControlutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	ricInstanceID, err := e2apControlutils.GetRicInstanceID(request)
	if err != nil {
		return nil, nil, err
	}

	controlMessage, err := sm.getControlMessage(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	controlHeader, err := sm.getControlHeader(request)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	log.Warnf("CCC control header: %v", controlHeader)
	log.Warnf("CCC control message: %v", controlMessage)

	// Received time stamp
	timestamp := make([]byte, 8)
	binary.BigEndian.PutUint64(timestamp, uint64(time.Now().Unix()))

	acceptedPolicies := make([]*e2smccc.ConfigurationStructureAccepted, 0)
	failedPolicies := make([]*e2smccc.ConfigurationStructureFailed, 0)

	// TODO: GA: Here is where the connection between RAN-SIMULATOR and JSON file (OAI) happens
	ricStyleType := controlHeader.GetControlHeaderFormat().GetE2SmCccControlHeaderFormat1().GetRicStyleType().GetValue()
	if ricStyleType == controlServiceStyleType && controlServiceStyleType == 1 {
		controlMessageFormat1 := controlMessage.GetControlMessageFormat().GetE2SmCccControlMessageFormat1()
		log.Debugf("controlMessageFormat1: %v", controlMessageFormat1)

		// Message array to store policyRatio in controlMessage
		var policyMsg Message

		for _, configurationStructureWrite := range controlMessageFormat1.GetListOfConfigurationStructures().GetValue() {
			ranConfigurationStructureName := configurationStructureWrite.GetRanConfigurationStructureName().String()
			policyRatio := configurationStructureWrite.GetNewValuesOfAttributes().GetRanConfigurationStructure().GetORrmpolicyRatio()

			pMsg := new(RrmPolicyRatio)
			if policyRatio != nil {
				resourceType := policyRatio.GetResourceType().String()
				schedulerType := policyRatio.GetSchedulerType().String()
				for _, policyMember := range policyRatio.GetRRmpolicyMemberList().GetValue() {
					plmnId := policyMember.PlmnId.String()
					sst := policyMember.GetSnssai().GetSSt().GetValue()
					sd := policyMember.GetSnssai().GetSD().GetValue()

					sd_hex_string := fmt.Sprintf("%X", sd)
					log.Warnf("SD is %s\n", sd_hex_string)
					sd_deciNum, err := strconv.ParseInt(sd_hex_string, 16, 64)
					if err != nil {
						log.Warnf("SD %s Conversion from string to Int failed: %s\n", sd_hex_string, err)
					}
					pMsg.SST = uint8(sst[0])
					pMsg.SD = uint32(sd_deciNum)
					pMsg.SDFlag = 1
					log.Warnf("New Policy ratio member: plmn: %v, sst: %v, sd: %v", plmnId, sst, sd)
				}
				policyDedicatedRatio := policyRatio.GetRRmpolicyDedicatedRatio()
				policyMinRatio := policyRatio.GetRRmpolicyMinRatio()
				policyMaxRatio := policyRatio.GetRRmpolicyMaxRatio()
				log.Warnf("Policy ratio '%v' has these details: resourceType: %v, schedulerType: %v, policyDedicatedRatio: %v, policyMinRatio: %v, policyMaxRatio: %v", ranConfigurationStructureName, resourceType, schedulerType, policyDedicatedRatio, policyMinRatio, policyMaxRatio)

				pMsg.MinRatio = uint8(policyMinRatio)
				pMsg.MaxRatio = uint8(policyMaxRatio)
				pMsg.DedicatedRatio = uint8(policyDedicatedRatio)
				policyMsg.RrmPolicyRatio = append(policyMsg.RrmPolicyRatio, *pMsg)
			}
			oldPolicyRatio := configurationStructureWrite.GetOldValuesOfAttributes().GetRanConfigurationStructure().GetORrmpolicyRatio()
			if oldPolicyRatio != nil {
				resourceType := oldPolicyRatio.GetResourceType().String()
				schedulerType := oldPolicyRatio.GetSchedulerType().String()
				for _, oldPolicyMember := range oldPolicyRatio.GetRRmpolicyMemberList().GetValue() {
					plmnId := oldPolicyMember.PlmnId.String()
					sst := oldPolicyMember.GetSnssai().GetSSt().String()
					sd := oldPolicyMember.GetSnssai().GetSD().String()
					log.Warnf("Old Policy ratio member: plmn: %v, sst: %v, sd: %v", plmnId, sst, sd)
				}
				policyDedicatedRatio := oldPolicyRatio.GetRRmpolicyDedicatedRatio()
				policyMinRatio := oldPolicyRatio.GetRRmpolicyMinRatio()
				policyMaxRatio := oldPolicyRatio.GetRRmpolicyMaxRatio()
				log.Warnf("Policy ratio '%v' has these details: resourceType: %v, schedulerType: %v, policyDedicatedRatio: %v, policyMinRatio: %v, policyMaxRatio: %v", ranConfigurationStructureName, resourceType, schedulerType, policyDedicatedRatio, policyMinRatio, policyMaxRatio)
			}

			// Accepted Policy
			acceptedPolicy := &e2smccc.ConfigurationStructureAccepted{
				RanConfigurationStructureName: &e2smccc.RanConfigurationStructureName{
					Value: []byte(ranConfigurationStructureName),
				},
				CurrentValuesOfAttributes: &e2smccc.ValuesOfAttributes{
					RanConfigurationStructure: &e2smccc.RanConfigurationStructure{
						RanConfigurationStructure: &e2smccc.RanConfigurationStructure_ORrmpolicyRatio{
							ORrmpolicyRatio: policyRatio,
						},
					},
				},
				OldValuesOfAttributes: &e2smccc.ValuesOfAttributes{
					RanConfigurationStructure: &e2smccc.RanConfigurationStructure{
						RanConfigurationStructure: &e2smccc.RanConfigurationStructure_ORrmpolicyRatio{
							ORrmpolicyRatio: oldPolicyRatio,
						},
					},
				},
			}

			acceptedPolicies = append(acceptedPolicies, acceptedPolicy)

			// // Failed Policy
			// failedPolicy := &e2smccc.ConfigurationStructureFailed{
			// 	RanConfigurationStructureName: &e2smccc.RanConfigurationStructureName{
			// 		Value: []byte(ranConfigurationStructureName),
			// 	},
			// 	RequestedValuesOfAttributes: &e2smccc.ValuesOfAttributes{
			// 		RanConfigurationStructure: &e2smccc.RanConfigurationStructure{
			// 			RanConfigurationStructure: &e2smccc.RanConfigurationStructure_ORrmpolicyRatio{
			// 				ORrmpolicyRatio: policyRatio,
			// 			},
			// 		},
			// 	},
			// 	OldValuesOfAttributes: &e2smccc.ValuesOfAttributes{
			// 		RanConfigurationStructure: &e2smccc.RanConfigurationStructure{
			// 			RanConfigurationStructure: &e2smccc.RanConfigurationStructure_ORrmpolicyRatio{
			// 				ORrmpolicyRatio: oldPolicyRatio,
			// 			},
			// 		},
			// 	},
			// 	Cause: e2smccc.OutcomeCause_OUTCOME_CAUSE_UNSPECIFIED,
			// }

			// failedPolicies = append(failedPolicies, failedPolicy)
		}
		// write policy CCC Msg into .json for OAI
		Msg2Json(policyMsg)
	} else {
		// TODO: GA: Currently NOT supporting format 2
		controlMessageFormat2 := controlMessage.GetControlMessageFormat().GetE2SmCccControlMessageFormat2()
		log.Warnf("Currently not supporting Control Message Format 2: %v", controlMessageFormat2)
	}

	// acceptedList, err := pdubuilder.CreateRanConfigurationStructuresAcceptedList(acceptedPolicies)
	// if err != nil {
	// 	log.Warn(err)
	// }
	// failedList, err := pdubuilder.CreateRanConfigurationStructuresFailedList(failedPolicies)
	// if err != nil {
	// 	log.Warn(err)
	// }

	acceptedList := &e2smccc.RanConfigurationStructuresAcceptedList{
		Value: acceptedPolicies,
	}
	failedList := &e2smccc.RanConfigurationStructuresFailedList{
		Value: failedPolicies,
	}

	// End GA
	if err != nil {
		outcomeAsn1Bytes, err := controloutcome.NewControlOutcome(
			controloutcome.WithReceivedTimestamp(timestamp),
			controloutcome.WithRanConfigurationStructuresAcceptedList(acceptedList),
			controloutcome.WithRanConfigurationStructuresFailedList(failedList)).
			ToAsn1Bytes()
		if err != nil {
			return nil, nil, err
		}
		failure, err = e2apControlutils.NewControl(
			e2apControlutils.WithRanFuncID(*ranFuncID),
			e2apControlutils.WithRequestID(*reqID),
			e2apControlutils.WithRicInstanceID(*ricInstanceID),
			e2apControlutils.WithRicControlOutcome(outcomeAsn1Bytes)).BuildControlFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, failure, nil
	}

	outcomeAsn1Bytes, err := controloutcome.NewControlOutcome(
		controloutcome.WithReceivedTimestamp(timestamp),
		controloutcome.WithRanConfigurationStructuresAcceptedList(acceptedList),
		controloutcome.WithRanConfigurationStructuresFailedList(failedList)).
		ToAsn1Bytes()
	if err != nil {
		log.Warn(err)
		return nil, nil, err
	}
	log.Warnf("outcomeAsn1Bytes: %v", outcomeAsn1Bytes)

	response, err = e2apControlutils.NewControl(
		e2apControlutils.WithRanFuncID(*ranFuncID),
		e2apControlutils.WithRequestID(*reqID),
		e2apControlutils.WithRicInstanceID(*ricInstanceID),
		e2apControlutils.WithRicControlOutcome(outcomeAsn1Bytes)).BuildControlAcknowledge()
	if err != nil {
		log.Warn(err)
		return nil, nil, err
	}
	log.Warnf("response: %v", response)
	return response, nil, nil
}

// RICSubscription implements subscription handler for ccc service model
func (sm *Client) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	log.Infof("RIC Subscription request received for e2 node %d and service model %s:", sm.ServiceModel.Node.GnbID, sm.ServiceModel.ModelName)
	ricActionsAccepted := make([]*e2aptypes.RicActionID, 0)
	ricActionsNotAdmitted := make(map[e2aptypes.RicActionID]*e2apies.Cause)
	actionList := subutils.GetRicActionToBeSetupList(request)
	reqID, err := subutils.GetRequesterID(request)
	if err != nil {
		return nil, nil, err
	}
	ranFuncID, err := subutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	ricInstanceID, err := subutils.GetRicInstanceID(request)
	if err != nil {
		return nil, nil, err
	}

	for _, action := range actionList {
		actionID := e2aptypes.RicActionID(action.GetValue().GetRicactionToBeSetupItem().GetRicActionId().GetValue())
		actionType := action.GetValue().GetRicactionToBeSetupItem().GetRicActionType()
		// ccc service model supports report action and should be added to the
		// list of accepted actions
		if actionType == e2apies.RicactionType_RICACTION_TYPE_REPORT {
			ricActionsAccepted = append(ricActionsAccepted, &actionID)
		}
		// ccc service model does not support INSERT and POLICY actions and
		// should be added into the list of not admitted actions
		if actionType == e2apies.RicactionType_RICACTION_TYPE_INSERT ||
			actionType == e2apies.RicactionType_RICACTION_TYPE_POLICY {
			cause := &e2apies.Cause{
				Cause: &e2apies.Cause_RicRequest{
					RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_ACTION_NOT_SUPPORTED,
				},
			}
			ricActionsNotAdmitted[actionID] = cause
		}
	}

	// At least one required action must be accepted otherwise sends a subscription failure response
	if len(ricActionsAccepted) == 0 {
		log.Warn("no action is accepted")
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_ACTION_NOT_SUPPORTED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, subscriptionFailure, nil
	}

	// reportInterval, err := sm.getReportInterval(request)
	// if err != nil {
	// 	log.Warn(err)
	// 	cause := &e2apies.Cause{
	// 		Cause: &e2apies.Cause_RicRequest{
	// 			RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
	// 		},
	// 	}
	// 	subscription := subutils.NewSubscription(
	// 		subutils.WithRequestID(*reqID),
	// 		subutils.WithRanFuncID(*ranFuncID),
	// 		subutils.WithRicInstanceID(*ricInstanceID),
	// 		subutils.WithCause(cause))
	// 	subscriptionFailure, err := subscription.BuildSubscriptionFailure()
	// 	if err != nil {
	// 		return nil, subscriptionFailure, nil
	// 	}
	// 	return nil, subscriptionFailure, nil
	// }

	actionDefinitions, err := sm.getActionDefinition(actionList, ricActionsAccepted)
	if err != nil {
		log.Warn(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			log.Warn(err)
			return nil, subscriptionFailure, nil
		}
		return nil, subscriptionFailure, nil
	}

	subscription := subutils.NewSubscription(
		subutils.WithRequestID(*reqID),
		subutils.WithRanFuncID(*ranFuncID),
		subutils.WithRicInstanceID(*ricInstanceID),
		subutils.WithActionsAccepted(ricActionsAccepted),
		subutils.WithActionsNotAdmitted(ricActionsNotAdmitted))

	response, err = subscription.BuildSubscriptionResponse()
	if err != nil {
		log.Warn(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, subscriptionFailure, nil
		}
		return nil, subscriptionFailure, nil
	}

	// reason := e2smccc.IndicationReason(e2smccc.IndicationReason_INDICATION_REASON_UPON_SUBSCRIPTION)
	log.Debug("Prepare RIC Indication: Upon Subscription")
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := sm.sendRicIndication(ctx, subscription, actionDefinitions, e2smccc.IndicationReason_INDICATION_REASON_UPON_SUBSCRIPTION, e2smccc.ChangeType_CHANGE_TYPE_ADDITION)
		if err != nil {
			return
		}
	}()

	return response, nil, nil
}

// RICSubscriptionDelete implements subscription delete handler for ccc service model
func (sm *Client) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	log.Infof("RIC subscription delete request is received for e2 node %d and  service model %s:", sm.ServiceModel.Node.GnbID, sm.ServiceModel.ModelName)
	reqID, err := subdeleteutils.GetRequesterID(request)
	if err != nil {
		return nil, nil, err
	}
	ranFuncID, err := subdeleteutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	ricInstanceID, err := subdeleteutils.GetRicInstanceID(request)
	if err != nil {
		return nil, nil, err
	}
	subID := subscriptions.NewID(*ricInstanceID, *reqID, *ranFuncID)
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return nil, nil, err
	}
	subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
		subdeleteutils.WithRequestID(*reqID),
		subdeleteutils.WithRanFuncID(*ranFuncID),
		subdeleteutils.WithRicInstanceID(*ricInstanceID))
	subDeleteResponse, err := subscriptionDelete.BuildSubscriptionDeleteResponse()
	if err != nil {
		return nil, nil, err
	}
	// Stops the goroutine sending the indication messages
	if sub.Ticker != nil {
		sub.Ticker.Stop()
	}
	return subDeleteResponse, nil, nil
}
