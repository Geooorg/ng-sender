package common

import "fmt"

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
		return "WarningMessage"
	case ExerciseWarningMessage:
		return "ExerciseInfoMessage"
	case InfoMessage:
		return "InfoMessage"
	case ExerciseInfoMessage:
		return "ExerciseInfoMessage"
	default:
		return fmt.Sprintf("Unknown type: (%d)", m)
	}
}
