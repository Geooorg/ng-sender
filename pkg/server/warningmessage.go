package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	temporal "go.temporal.io/sdk/client"
	"log"
	"net/http"
	"ng-sender/pkg/common"
	"ng-sender/pkg/workflow"
)

func (s *Server) SendWarningMessageToAllReceivers(w http.ResponseWriter, r *http.Request) {

	println("sending warning message to all receivers")

	var envelope common.MessageEnvelope
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&envelope)
	if err != nil {
		println("WARN: message decoding failed", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	err = validateRequest(envelope)
	if err != nil {
		println("WARN: message validation failed", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	s.LogMessage(common.WarningMessage, envelope.UUID, envelope.WarnrechnerStationId)

	s.sendToReceivers(context.Background(), envelope)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateRequest(e common.MessageEnvelope) error {
	expectedMessageType := common.AsString(common.WarningMessage)

	if e.MessageType != common.AsString(common.WarningMessage) {
		return errors.New("MessageType was not as expected ' " + expectedMessageType + "'.")
	}
	if e.UUID.String() == "" {
		return errors.New("UUID must not be empty")
	}
	if e.WarnrechnerStationId == "" {
		return errors.New("WarnrechnerStationId must not be empty")
	}

	return nil
}

const queueName = "warningMessages"
const workflowIDPrefix = "warningMessage-"

func (s *Server) sendToReceivers(ctx context.Context, envelope common.MessageEnvelope) error {

	s.updateStationsListIfNeeded()

	envelopeAsJson, err := json.Marshal(envelope)
	if err != nil {
		log.Println("WARN: Could not marshall envelope to JSON!", err)
		return err
	}

	// start workflow for each host to send message to
	for _, receiver := range s.stationListCache.StationsList.Receivers {
		for _, host := range receiver.Hosts {

			hostAndPort := host.Hostname + ":" + host.Port
			log.Println(fmt.Printf("INFO: starting workflow for envelope %s, station %s on host %s",
				envelope.UUID.String(), receiver.ID, hostAndPort,
			))

			workflowID := workflowIDPrefix + envelope.UUID.String() + "_" + receiver.ID + "_" + hostAndPort

			workflowExecution, err := (*s.TemporalClient).ExecuteWorkflow(ctx, temporal.StartWorkflowOptions{
				ID:        workflowID,
				TaskQueue: queueName,
			}, workflow.SendToReceiversWF, envelopeAsJson, envelope.UUID.String(), hostAndPort)

			if err != nil {
				log.Println("WARN: Unable to execute workflow with ID", workflowID, err)
			} else {
				log.Println("INFO: Workflow execution started with ID / RunID: ", "ID", workflowExecution.GetID(), workflowExecution.GetRunID())
			}
		}

	}

	return nil
}
