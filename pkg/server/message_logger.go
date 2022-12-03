package server

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	c "ng-receiver/pkg/common"
	"os"
	"strings"
	"time"
)

func (s *Server) LogMessage(messageType c.MessageType, uuid uuid.UUID, stationId string) {

	timestamp := fmt.Sprintf("%s", time.Now().Format(time.RFC3339))
	fields := []string{timestamp, fmt.Sprint(c.AsString(messageType)), uuid.String(), stationId}
	msg := strings.Join(fields, "\t") + "\n"
	println(msg)

	logFile := s.determineLogFileName(messageType)

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		println("WARN: Could not append to logfile " + logFile)
	}

	defer f.Close()

	if _, err = io.WriteString(f, msg); err != nil {
		println("WARN: Could not write warning message event to logfile " + logFile)
	}

	f.Sync()
}
