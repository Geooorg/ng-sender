package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type MessageType int64

const (
	WarningMessage MessageType = iota
	ExerciseWarningMessage
	InfoMessage
	ExerciseInfoMessage
)

func AsString(m MessageType) string {
	switch m {
	case WarningMessage:
		return "warningMessage"
	case ExerciseWarningMessage:
		return "exerciseInfoMessage"
	case InfoMessage:
		return "infoMessage"
	case ExerciseInfoMessage:
		return "exerciseInfoMessage"
	default:
		return fmt.Sprintf("Unknown type: (%d)", m)
	}
}

type MessageEnvelope struct {
	UUID                 uuid.UUID `json:"uuid"`
	MessageType          string    `json:"messageType"`
	ObjectNode           string    `json:"objectNode"` // note this is the actual JSON payload
	SentAtTimestamp      int64     `json:"sentAtTimestamp"`
	WarnrechnerId        string    `json:"warnrechnerId"`
	WarnrechnerStationId string    `json:"warnrechnerStationId"`
	WarnrechnerHostname  string    `json:"warnrechnerHostname"`
}

func WarningMessageToEnvelope(b []byte) (MessageEnvelope, error) {
	var envelope MessageEnvelope
	dec := json.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(&envelope)
	if err != nil {
		println("WARN: message decoding failed", err)

	}
	return envelope, err
}

func WarningMessageToJson(envelope MessageEnvelope) ([]byte, error) {
	return json.Marshal(envelope)
}
