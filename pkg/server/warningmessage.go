package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	temporal "go.temporal.io/sdk/client"
	"io"
	"log"
	"net/http"
	"ng-sender/pkg/common"
	"ng-sender/pkg/workflow"
)

type WarningMessageSender interface {
	sendWarningMessageToAllReceivers(envelope common.MessageEnvelope) error
}

func (s *Server) OnWarningMessageReceivedNATS(m *nats.Msg) {
	println("OnWarningMessageReceivedNATS")

	envelope, err := common.WarningMessageToEnvelope(m.Data)

	if err != nil {
		log.Println("WARN: Could not convert warning message to JSON received by NATS")
	}
	s.sendWarningMessageToAllReceivers(envelope)
}

func (s *Server) OnWarningMessageReceivedHTTP(w http.ResponseWriter, r *http.Request) {

	bytes, err := io.ReadAll(r.Body)
	envelope, err := common.WarningMessageToEnvelope(bytes)
	if err != nil {
		println("WARN: message decoding failed", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.sendWarningMessageToAllReceivers(envelope)
	if err != nil {
		println("WARN: Could not call 'sendWarningMessageToAllReceivers'", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) sendWarningMessageToAllReceivers(envelope common.MessageEnvelope) error {

	println("Sending warning message to all receivers")

	err := validateRequest(envelope)
	if err != nil {
		println("WARN: message validation failed", err)
		return err
	}

	json, err := common.WarningMessageToJson(envelope)
	if err != nil {
		println("WARN: message conversion to JSON failed", err)
		return err
	}

	s.PersistMessage(common.WarningMessage, json, envelope.UUID, envelope.WarnrechnerStationId)

	s.startSendToReceiversWorkflow(context.Background(), envelope)

	return err
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

func (s *Server) startSendToReceiversWorkflow(ctx context.Context, envelope common.MessageEnvelope) error {

	s.updateStationsListIfNeeded()

	envelopeAsJson, err := common.WarningMessageToJson(envelope)
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
