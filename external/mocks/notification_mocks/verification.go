package notification_mocks

import (
	"fmt"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/utility"
)

func VerificationFailedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	_, ok := idata.(external_models.VerificationFailedModel)
	if !ok {
		logger.Error("verification failed notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	logger.Info("verification failed notification", outBoundResponse)

	return nil, nil
}

func VerificationSuccessfulNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	_, ok := idata.(external_models.VerificationSuccessfulModel)
	if !ok {
		logger.Error("verification successful notification", "incorrect data format", idata)
		return nil, fmt.Errorf("request data format error")
	}

	logger.Info("verification successful notification", outBoundResponse)

	return nil, nil
}
