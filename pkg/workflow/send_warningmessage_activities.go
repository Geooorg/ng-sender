package workflow

import (
	"context"
	"go.temporal.io/sdk/activity"

	"ng-sender/pkg/common"
)

type WarningMessageActivities struct {
}

func (a *WarningMessageActivities) sendWarningMessageToHost(ctx context.Context, warningMessage common.MessageEnvelope, stationId string, hostAndPortString string) {
	logger := activity.GetLogger(ctx)
	logger.Info("executing activity 'sendWarningMessageToHost' for WarningMessage, Station and Host", warningMessage.UUID.String(), stationId, hostAndPortString)

}
