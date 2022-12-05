package server

import (
	"errors"
	"github.com/gorilla/mux"
	temporal "go.temporal.io/sdk/client"
	"log"
	"net/http"
	c "ng-receiver/pkg/common"
	"os"
	"time"
)

type Server struct {
	TemporalClient   *temporal.Client
	StationsEndpoint string
	Port             string
	LogDirectory     string
	stationListCache StationListCache
}

type StationListCache struct {
	StationsList StationsListDto
	LastUpdated  time.Time
}

func (s *Server) RegisterHandlersAndServe() error {

	s.createMessageLogFiles()
	router := mux.NewRouter()

	router.HandleFunc("/warningmessage", s.SendWarningMessageToAllReceivers).Methods("POST")

	s.updateStationsListIfNeeded()

	println("Server listening on port " + s.Port)
	err := http.ListenAndServe(":"+s.Port, router)
	if err != nil {
		log.Fatal("HTTP server could not be started", err)
	}
	return err
}

func (s *Server) updateStationsListIfNeeded() {
	filteredStations := StationsListDto{}
	const receiverTypeStationFilter = "STATION"

	if s.stationListNeedsUpdate() {

		stations, err := s.fetchStationsList()

		if err == nil && len(stations.Receivers) > 0 {
			for _, r := range stations.Receivers {

				if r.ReceiverType.Category == receiverTypeStationFilter {
					filteredStations.Receivers = append(filteredStations.Receivers, r)
				}
			}

			s.stationListCache.StationsList = filteredStations
			s.stationListCache.LastUpdated = time.Now()
		} else {
			println("WARN: Could not update the stations list or list was empty")
		}
	} else {
		println("INFO: No need to update Station list")
	}

}

func (s *Server) stationListNeedsUpdate() bool {
	isOlderThanOneHour := time.Since(s.stationListCache.LastUpdated) > 1*time.Hour
	return len(s.stationListCache.StationsList.Receivers) == 0 || isOlderThanOneHour
}

func (s *Server) createMessageLogFiles() {

	if _, err := os.Stat(s.LogDirectory); os.IsNotExist(err) {
		log.Fatal("Log Directory does not exist", err)
	}

	logFileTypes := []c.MessageType{c.WarningMessage, c.ExerciseWarningMessage, c.InfoMessage, c.ExerciseInfoMessage}

	for messageType := range logFileTypes {
		logFileName := s.determineLogFileName(c.MessageType(messageType))

		if _, err := os.Stat(logFileName); errors.Is(err, os.ErrNotExist) {
			_, err := os.Create(logFileName)
			if err != nil {
				log.Fatalln("Could not create log file "+logFileName, err)
			}
		}

	}
}

func (s *Server) determineLogFileName(messageType c.MessageType) string {
	typeName := c.AsString(messageType)
	return s.LogDirectory + "/" + typeName + ".log"
}
