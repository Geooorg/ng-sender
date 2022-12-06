package workflow

import (
	"bytes"
	"context"
	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/activity"
	"net/http"
	"time"
)

type WarningMessageActivities struct {
	NatsClient *nats.Conn
}

const warningMessagePath = "/warningMessage"
const topic = "warningMessageSent_Central"

func (a *WarningMessageActivities) SendWarningMessageToHost(ctx context.Context, envelopeAsJson []byte, uuid string, hostAndPort string) (time.Time, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("executing activity 'sendWarningMessageToHost' for WarningMessage with UUID " + uuid + " to Host " + hostAndPort)

	url := "http://" + hostAndPort + warningMessagePath

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(envelopeAsJson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Warn("Could not send warning message", err)
		return time.Now(), err
	}
	defer resp.Body.Close()

	return time.Now(), nil
}

func (a *WarningMessageActivities) PublishEvent(ctx context.Context, envelopeAsJson []byte, uuid string) error {
	logger := activity.GetLogger(ctx)

	err := a.NatsClient.Publish(topic, envelopeAsJson)
	if err != nil {
		logger.Warn("WARN: Could not publish message "+uuid+" to queue "+topic, err)
	}

	return err
}
