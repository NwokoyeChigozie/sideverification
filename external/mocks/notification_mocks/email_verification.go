package notification_mocks

import (
	"fmt"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/microservice/auth"
	"github.com/vesicash/verification-ms/utility"
)

func SendVerificationEmail(logger *utility.Logger, idata interface{}) (interface{}, error) {

	_, ok := idata.(external_models.EmailNotificationRequest)
	if !ok {
		logger.Info("get user", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	_, err := auth.GetAccessToken(logger)
	if err != nil {
		logger.Info("verification email", nil, err)
		return nil, err
	}

	return nil, nil
}
