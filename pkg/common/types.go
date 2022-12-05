package common

import (
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
	UUID    uuid.UUID
	content []byte
}
