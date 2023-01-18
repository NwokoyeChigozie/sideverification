package notification

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/microservice/auth"
	"github.com/vesicash/verification-ms/utility"
)

func VerificationFailedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.VerificationFailedModel)
	if !ok {
		logger.Info("verification failed notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := auth.GetAccessToken(logger)
	if err != nil {
		logger.Info("verification failed notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("verification failed notification", data)
	err = external.SendRequest(logger, "service", "verification_failed_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("verification failed notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("verification failed notification", outBoundResponse)

	return nil, nil
}

func VerificationSuccessfulNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.VerificationSuccessfulModel)
	if !ok {
		logger.Info("verification successful notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := auth.GetAccessToken(logger)
	if err != nil {
		logger.Info("verification successful notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("verification successful notification", data)
	err = external.SendRequest(logger, "service", "verification_successful_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("verification successful notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("verification successful notification", outBoundResponse)

	return nil, nil
}
