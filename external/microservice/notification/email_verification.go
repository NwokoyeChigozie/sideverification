package notification

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/microservice/auth"
	"github.com/vesicash/verification-ms/utility"
)

func SendVerificationEmail(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.EmailNotificationRequest)
	if !ok {
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := auth.GetAccessToken(logger)
	if err != nil {
		logger.Info("verification email", outBoundResponse, err)
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("verification email", data)
	err = external.SendRequest(logger, "service", "send_verification_email", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("verification email", outBoundResponse, err)
		return nil, err
	}
	logger.Info("verification email", outBoundResponse)

	return nil, nil
}

func SendWelcomeEmail(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.AccountIDRequestModel)
	if !ok {
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := auth.GetAccessToken(logger)
	if err != nil {
		logger.Info("welcome email", outBoundResponse, err)
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("welcome email", data)
	err = external.SendRequest(logger, "service", "send_welcome_email", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("welcome email", outBoundResponse, err)
		return nil, err
	}
	logger.Info("welcome email", outBoundResponse)

	return nil, nil
}

func SendEmailVerifiedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.AccountIDRequestModel)
	if !ok {
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := auth.GetAccessToken(logger)
	if err != nil {
		logger.Info("email verified notification", outBoundResponse, err)
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("email verified notification", data)
	err = external.SendRequest(logger, "service", "send_email_verified_notification", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("email verified notification", outBoundResponse, err)
		return nil, err
	}
	logger.Info("email verified notification", outBoundResponse)

	return nil, nil
}
