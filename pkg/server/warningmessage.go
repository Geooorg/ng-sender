package server

import (
	"github.com/google/uuid"
	"net/http"
	"ng-receiver/pkg/common"
)

func (s *Server) CreateWarningMessage(w http.ResponseWriter, r *http.Request) {

	println("received warning message")

	newUUID, _ := uuid.NewUUID()
	s.LogMessage(common.WarningMessage, newUUID, "HH-GS-42")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
