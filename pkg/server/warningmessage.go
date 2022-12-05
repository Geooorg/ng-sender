package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	temporal "go.temporal.io/sdk/client"
	"log"
	"net/http"
	"ng-sender/pkg/common"
	"ng-sender/pkg/workflow"
)

func (s *Server) SendWarningMessageToAllReceivers(w http.ResponseWriter, r *http.Request) {

	println("received warning message")

	// TODO extract MessageEnvelope, write HTTP test file
	newUUID, _ := uuid.NewUUID()
	s.LogMessage(common.WarningMessage, newUUID, "HH-GS-42")

	s.sendToReceivers(context.Background(), nil)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

const queueName = "warningMessages"
const workflowIDPrefix = "warningMessage-"

func (s *Server) sendToReceivers(ctx context.Context, warningMessage common.MessageEnvelope) {

	s.updateStationsListIfNeeded()

	// start workflow for each host to send message to
	for _, receiver := range s.stationListCache.StationsList.Receivers {
		for _, host := range receiver.Hosts {

			hostAndPort := host.Hostname + host.Port
			log.Println(fmt.Printf("INFO: starting workflow for warningMessage %s, station %s on host %s",
				warningMessage.UUID.String(), receiver.ID, hostAndPort,
			))

			workflowID := workflowIDPrefix + warningMessage.UUID.String() + "_" + receiver.ID + "_" + hostAndPort

			workflowExecution, err := (*s.TemporalClient).ExecuteWorkflow(ctx, temporal.StartWorkflowOptions{
				ID:        workflowID,
				TaskQueue: queueName,
			}, workflow.SendToReceiversWF, warningMessage, hostAndPort)

			if err != nil {
				log.Println("WARN: Unable to execute workflow with ID", workflowID, err)
			} else {
				log.Println("INFO: Workflow execution started with ID / RunID: ", "ID", workflowExecution.GetID(), workflowExecution.GetRunID())
			}
		}

	}

}
