package rave

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func RaveResolveBankAccount(logger *utility.Logger, idata interface{}) (string, error) {

	var (
		outBoundResponse external_models.ResolveAccountSuccessResponse
	)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + config.GetConfig().Rave.SecretKey,
	}

	data, ok := idata.(external_models.ResolveAccountRequest)
	if !ok {
		logger.Info("rave resolve bank account", idata, "request data format error")
		return "", fmt.Errorf("request data format error")
	}

	err := external.SendRequest(logger, "third_party", "rave_resolve_bank_account", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("rave resolve bank account", outBoundResponse, err.Error())
		return "", err
	}
	logger.Info("rave resolve bank account", outBoundResponse)

	return outBoundResponse.Data.AccountName, nil
}
