package notification

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/microservice/auth"
	"github.com/vesicash/verification-ms/utility"
)

func SendAuthorizedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.AuthorizeNotificationRequest)
	if !ok {
		logger.Info("authorized notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := auth.GetAccessToken(logger)
	if err != nil {
		logger.Info("authorized notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("authorized notification", data)
	err = external.SendRequest(logger, "service", "send_authorized_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("authorized notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("authorized notification", outBoundResponse)

	return nil, nil
}

func SendAuthorizationNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.AuthorizeNotificationRequest)
	if !ok {
		logger.Info("authorization notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := auth.GetAccessToken(logger)
	if err != nil {
		logger.Info("authorization notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("verification failed notification", data)
	err = external.SendRequest(logger, "service", "send_authorization_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("authorization notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("authorization notification", outBoundResponse)

	return nil, nil
}
