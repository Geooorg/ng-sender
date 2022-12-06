package workflow

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"time"
)

func SendToReceiversWF(ctx workflow.Context, envelopeAsJson []byte, uuid string, hostAndPort string) (time.Time, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("SendToReceiversWF starting for UUID " + uuid + " and stationId on " + hostAndPort)

	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 1.05,
		MaximumInterval:    time.Second * 10,
		MaximumAttempts:    0, // Unlimited
	}
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)
	var a *WarningMessageActivities

	var finishedAt time.Time

	if err := workflow.ExecuteActivity(ctx, a.SendWarningMessageToHost, envelopeAsJson, uuid, hostAndPort).Get(ctx, &finishedAt); err != nil {
		logger.Error("SendWarningMessageToHost Activity failed.", "Error", err.Error())
		return time.Now(), err
	}

	if err := workflow.ExecuteActivity(ctx, a.PublishEvent, envelopeAsJson, uuid).Get(ctx, &finishedAt); err != nil {
		logger.Error("PublishEvent failed.", "Error", err.Error())
		return time.Now(), err
	}

	return time.Now(), nil
}
