// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package labelinfo

import (
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
	"github.com/onosproject/onos-lib-go/pkg/errors"
)

// LabelInfo labels info
type LabelInfo struct {
	plmnID             ransimtypes.Uint24
	sst                []byte
	sd                 []byte
	fiveQI             int32
	qfi                int32
	qci                int32
	qciMax             int32
	qciMin             int32
	arpMax             int32
	arpMin             int32
	bitrateRange       int32
	layerMuMimo        int32
	distX              int32
	distY              int32
	distZ              int32
	startEndIndication e2smkpmv2.StartEndInd
}

// NewLabelInfo creates a new label info
func NewLabelInfo(options ...func(*LabelInfo) error) (*LabelInfo, error) {
	labelInfo := &LabelInfo{}
	for _, option := range options {
		err := option(labelInfo)
		if err != nil {
			return nil, err
		}
	}

	return labelInfo, nil
}

// WithPlmnID sets plmn ID
func WithPlmnID(plmnID ransimtypes.Uint24) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		info.plmnID = plmnID
		return nil
	}
}

// WithSST sets SST
func WithSST(sst []byte) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		if len(sst) != 1 {
			return errors.NewInvalid("sst must be 1 byte")
		}
		info.sst = sst
		return nil
	}
}

// WithSD sets SD
func WithSD(sd []byte) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		if len(sd) != 3 {
			return errors.NewInvalid("sst must be 3 bytes")
		}
		info.sd = sd
		return nil
	}
}

// WithFiveQI sets five QI
func WithFiveQI(fiveQI int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		if fiveQI < 0 || fiveQI > 255 {
			return errors.NewInvalid("five QI values must be in rang [0, 255]")
		}
		info.fiveQI = fiveQI
		return nil
	}
}

// WithQFI sets QFI
func WithQFI(qfi int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		if qfi < 0 || qfi > 63 {
			return errors.NewInvalid("qfi values must be in rang [0, 63]")

		}
		info.qfi = qfi
		return nil
	}
}

// WithQCI sets qci
func WithQCI(qci int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		if qci < 0 || qci > 255 {
			return errors.NewInvalid("QCI values must be in rang [0, 255]")
		}
		info.qci = qci
		return nil
	}

}

// WithQCIMax sets maximum qci
func WithQCIMax(qciMax int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		if qciMax < 0 || qciMax > 255 {
			return errors.NewInvalid("QCI Max values must be in rang [0, 255]")
		}
		info.qciMax = qciMax
		return nil
	}

}

// WithQCIMin sets minimum qci
func WithQCIMin(qciMin int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		if qciMin < 0 || qciMin > 255 {
			return errors.NewInvalid("QCI Max values must be in rang [0, 255]")
		}
		info.qciMin = qciMin
		return nil
	}

}

// WithArpMax sets arp max
func WithArpMax(arpMax int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		if arpMax < 0 || arpMax > 15 {
			return errors.NewInvalid("Arp Max values must be in rang [0, 15]")
		}

		info.arpMax = arpMax
		return nil
	}
}

// WithArpMin sets arp min
func WithArpMin(arpMin int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		if arpMin < 0 || arpMin > 15 {
			return errors.NewInvalid("Arp Max values must be in rang [0, 15]")
		}
		info.arpMin = arpMin
		return nil
	}
}

// WithBitRateRange sets bit rate range
func WithBitRateRange(bitrateRange int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		info.bitrateRange = bitrateRange
		return nil
	}
}

// WithLayerMuMimo sets layer muMimo
func WithLayerMuMimo(layerMuMimo int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		info.layerMuMimo = layerMuMimo
		return nil
	}
}

// WithDistX sets distX
func WithDistX(distX int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		info.distX = distX
		return nil
	}
}

// WithDistY sets distY
func WithDistY(distY int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		info.distY = distY
		return nil
	}
}

// WithDistZ sets distZ
func WithDistZ(distZ int32) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		info.distZ = distZ
		return nil
	}
}

// WithStartEndIndication sets start
func WithStartEndIndication(startEndIndication e2smkpmv2.StartEndInd) func(info *LabelInfo) error {
	return func(info *LabelInfo) error {
		info.startEndIndication = startEndIndication
		return nil
	}
}

// Build builds label information item
func (l *LabelInfo) Build() (*e2smkpmv2.LabelInfoItem, error) {
	sum := e2smkpmv2.SUM_SUM_TRUE
	prelabelOverride := e2smkpmv2.PreLabelOverride_PRE_LABEL_OVERRIDE_TRUE

	labelInfoItem := e2smkpmv2.LabelInfoItem{
		MeasLabel: &e2smkpmv2.MeasurementLabel{
			PlmnId: &e2smkpmv2.PlmnIdentity{
				Value: l.plmnID.ToBytes(),
			},
			SliceId: &e2smkpmv2.Snssai{
				SD:  l.sd,
				SSt: l.sst,
			},
			FiveQi: &e2smkpmv2.FiveQi{
				Value: l.fiveQI,
			},
			QFi: &e2smkpmv2.Qfi{
				Value: l.qfi,
			},
			QCi: &e2smkpmv2.Qci{
				Value: l.qci,
			},
			QCimax: &e2smkpmv2.Qci{
				Value: l.qciMax,
			},
			QCimin: &e2smkpmv2.Qci{
				Value: l.qciMin,
			},
			ARpmax: &e2smkpmv2.Arp{
				Value: l.arpMax,
			},
			ARpmin: &e2smkpmv2.Arp{
				Value: l.arpMin,
			},
			BitrateRange:     &l.bitrateRange,
			LayerMuMimo:      &l.layerMuMimo,
			SUm:              &sum,
			DistBinX:         &l.distX,
			DistBinY:         &l.distY,
			DistBinZ:         &l.distZ,
			PreLabelOverride: &prelabelOverride,
			StartEndInd:      &l.startEndIndication,
		},
	}

	// FIXME: Add back when ready
	//if err := labelInfoItem.Validate(); err != nil {
	//	return nil, errors.New(errors.Invalid, err.Error())
	//}

	return &labelInfoItem, nil

}
