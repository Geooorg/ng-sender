package workflow

import (
	"context"
	"go.temporal.io/sdk/activity"
	"time"

	"ng-sender/pkg/common"
)

type WarningMessageActivities struct {
}

func (a *WarningMessageActivities) SendWarningMessageToHost(ctx context.Context, warningMessage common.MessageEnvelope, stationId string, hostAndPortString string) (time.Time, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("executing activity 'sendWarningMessageToHost' for WarningMessage, Station and Host", warningMessage.UUID.String(), stationId, hostAndPortString)

	return time.Now(), nil
}
