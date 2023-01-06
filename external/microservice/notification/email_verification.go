package notification

import (
	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/microservice/auth"
	"github.com/vesicash/verification-ms/utility"
)

func SendVerificationEmail(data external_models.EmailNotificationRequest) error {
	logger := utility.NewLogger()
	var (
		outBoundResponse map[string]interface{}
	)
	accessToken, err := auth.GetAccessToken(logger)
	if err != nil {
		logger.Info("verification email", outBoundResponse, err)
		return err
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
		return err
	}
	logger.Info("verification email", outBoundResponse)

	return nil
}
