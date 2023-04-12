package rave

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
)

func (r *RequestObj) RaveResolveBankAccount() (string, error) {

	var (
		outBoundResponse external_models.ResolveAccountSuccessResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + config.GetConfig().Rave.SecretKey,
	}

	data, ok := idata.(external_models.ResolveAccountRequest)
	if !ok {
		logger.Error("rave resolve bank account", idata, "request data format error")
		return "", fmt.Errorf("request data format error")
	}

	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("rave resolve bank account", outBoundResponse, err.Error())
		if external.ResponseCode != nil {
			if *external.ResponseCode == 400 {
				return "", nil
			}
		}
		return "", err
	}
	logger.Info("rave resolve bank account", outBoundResponse)

	return outBoundResponse.Data.AccountName, nil
}
