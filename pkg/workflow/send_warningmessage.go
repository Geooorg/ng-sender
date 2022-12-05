package workflow

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"ng-sender/pkg/common"
	"time"
)

func SendToReceiversWF(ctx workflow.Context, warningMessage common.MessageEnvelope, stationId string, hostAndPort string) {
	//logger := workflow.GetLogger(ctx)

	//var activities *WarningMessageActivities

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
}
