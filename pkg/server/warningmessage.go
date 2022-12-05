package server

import (
	"github.com/google/uuid"
	"net/http"
	"ng-receiver/pkg/common"
)

func (s *Server) SendWarningMessageToAllReceivers(w http.ResponseWriter, r *http.Request) {

	println("received warning message")

	newUUID, _ := uuid.NewUUID()
	s.LogMessage(common.WarningMessage, newUUID, "HH-GS-42")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) extractTargetHosts() []string {

	var hosts = []string{}
	for _, r := range s.stationListCache.StationsList.Receivers {
		if len(r.Hosts) > 0 {
			for _, host := range r.Hosts {
				hosts = append(hosts, host.Hostname+host.Port)
			}
		}
	}

	return hosts
}
