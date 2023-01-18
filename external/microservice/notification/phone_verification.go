package notification

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/microservice/auth"
	"github.com/vesicash/verification-ms/utility"
)

func SendSendSMSToPhone(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(external_models.SMSToPhoneNotificationRequest)
	if !ok {
		logger.Info("sms to phone", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := auth.GetAccessToken(logger)
	if err != nil {
		logger.Info("send sms to phone", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("send sms to phone", data)
	err = external.SendRequest(logger, "service", "send_sms_to_phone", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("send sms to phone", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("send sms to phone", outBoundResponse)

	return nil, nil
}
