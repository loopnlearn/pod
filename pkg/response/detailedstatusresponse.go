package response

import (
	"encoding/hex"
)

type DetailedStatusResponse struct {
	LastProgSeqNum      uint8
	Reservoir           uint16
	Alerts              uint8
	BolusActive         bool
	TempBasalActive     bool
	BasalActive         bool
	ExtendedBolusActive bool
	PodProgress         PodProgress
	Delivered           uint16
	BolusRemaining      uint16
	MinutesActive       uint16
	FaultEvent          uint8
	FaultEventTime      uint16
}

func (r *DetailedStatusResponse) Marshal() ([]byte, error) {

	// OFF 1  2  3  4  5 6  7  8 9 10 1112 1314 1516 17 18 19 20 21 2223
	// 02 16 02 0J 0K LLLL MM NNNN PP QQQQ RRRR SSSS TT UU VV WW XX YYYY
	// 02 16 02 08 02 0000 00 01b2 00 0000 03ff 01cc 00 00 00 00 00 0000

	response, _ := hex.DecodeString("021602080200000001b200000003ff01cc00000000000000")

	// PodProgress
	if r.FaultEvent == 0 {
		response[3] = byte(r.PodProgress)
	} else {
		response[3] = PodProgressFault
	}

	// Delivery bits
	response[4] = 0 // suspended
	if r.ExtendedBolusActive {
		// extended bolus bit exclusive of bolus active bit
		response[4] |= 0b1000
	} else if r.BolusActive {
		response[4] |= 0b0100
	}
	if r.TempBasalActive {
		// temp basal active bit exclusive of basal active bit
		response[4] |= 0b0010
	} else if r.BasalActive {
		response[4] |= 0b0001
	}

	// Bolus remaining pulses
	response[5] = byte(r.BolusRemaining >> 8)
	response[6] = byte(r.BolusRemaining & 0xff)

	// LastProgSeqNum
	response[7] = r.LastProgSeqNum

	// Total delivered pulses
	response[8] = byte(r.Delivered >> 8)
	response[9] = byte(r.Delivered & 0xff)

	// Fault event
	response[10] = r.FaultEvent

	// Fault Event Time
	response[11] = byte(r.FaultEventTime >> 8)
	response[12] = byte(r.FaultEventTime & 0xff)

	// Reservoir
	if r.Reservoir <= (50 / 0.05) {
		response[13] = byte(r.Reservoir >> 8)
		response[14] = byte(r.Reservoir & 0xff)
	} else {
		response[13] = 0x03
		response[14] = 0xff
	}

	// Minutes since activation
	response[15] = byte(r.MinutesActive >> 8)
	response[16] = byte(r.MinutesActive & 0xff)

	// Set active alert slot bits
	response[17] = r.Alerts

	if r.FaultEvent != 0 {
		// previous PodProgress returned in low nibble of VV byte
		response[19] = byte(r.PodProgress) & 0b1111
		if r.BolusActive {
			response[19] |= 0b00010000
		}
	}

	return response, nil
}
