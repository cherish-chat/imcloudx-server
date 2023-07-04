package pb

import (
	"encoding/json"
	"errors"
	"strings"
)

// This is done this way because of a linter.
const (
	sdpTypeOfferStr    = "offer"
	sdpTypePranswerStr = "pranswer"
	sdpTypeAnswerStr   = "answer"
	sdpTypeRollbackStr = "rollback"
)

func (t SDPType) ToString() string {
	switch t {
	case SDPType_SDPTypeOffer:
		return sdpTypeOfferStr
	case SDPType_SDPTypePranswer:
		return sdpTypePranswerStr
	case SDPType_SDPTypeAnswer:
		return sdpTypeAnswerStr
	case SDPType_SDPTypeRollback:
		return sdpTypeRollbackStr
	default:
		return "Unknown"
	}
}

// MarshalJSON enables JSON marshaling of a SDPType
func (t SDPType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.ToString())
}

// UnmarshalJSON enables JSON unmarshaling of a SDPType
func (t *SDPType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		return errors.New("Unknown")
	case "offer":
		*t = SDPType_SDPTypeOffer
	case "pranswer":
		*t = SDPType_SDPTypePranswer
	case "answer":
		*t = SDPType_SDPTypeAnswer
	case "rollback":
		*t = SDPType_SDPTypeRollback
	}

	return nil
}
