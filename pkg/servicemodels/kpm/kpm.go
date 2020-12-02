// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm

import "github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"

// ServiceModel kpm service model struct
type ServiceModel struct {
}

// ProcessSubscriptionDelete ...
func (sm ServiceModel) ProcessSubscriptionDelete(request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	return nil, nil, nil
}

// ProcessSubscription ...
func (sm ServiceModel) ProcessSubscription(request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	return nil, nil, nil
}

// ProcessControl ...
func (sm ServiceModel) ProcessControl(request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	return nil, nil, nil
}
