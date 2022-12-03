package server

import (
	"errors"
	"github.com/gorilla/mux"
	temporal "go.temporal.io/sdk/client"
	"log"
	"net/http"
	c "ng-receiver/pkg/common"
	"os"
)

type Server struct {
	TemporalClient *temporal.Client
	Port           string
	LogDirectory   string
}

func (s *Server) RegisterHandlersAndServe() error {

	s.createMessageLogFiles()
	router := mux.NewRouter()

	router.HandleFunc("/warningmessage", s.CreateWarningMessage).Methods("POST")

	println("Server listening on port " + s.Port)
	err := http.ListenAndServe(":"+s.Port, router)
	if err != nil {
		log.Fatal("HTTP server could not be started", err)
	}
	return err
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
