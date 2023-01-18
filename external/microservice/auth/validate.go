package auth

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func ValidateOnAuth(logger *utility.Logger, idata interface{}) (bool, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.ValidateOnDBReqModel
	)

	data, ok := idata.(external_models.ValidateOnDBReq)
	if !ok {
		logger.Info("validate on auth", idata, "request data format error")
		return false, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("validate on auth", data)
	err := external.SendRequest(logger, "service", "validate_on_auth", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("validate on auth", outBoundResponse, err.Error())
		return false, err
	}
	logger.Info("validate on auth", outBoundResponse)

	return outBoundResponse.Data, nil
}

func ValidateAuthorization(logger *utility.Logger, idata interface{}) (external_models.ValidateAuthorizationDataModel, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.ValidateAuthorizationModel
	)

	data, ok := idata.(external_models.ValidateAuthorizationReq)
	if !ok {
		logger.Info("validate authorization", idata, "request data format error")
		return external_models.ValidateAuthorizationDataModel{}, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("validate authorization", data)
	err := external.SendRequest(logger, "service", "validate_authorization", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("validate authorization", outBoundResponse, err.Error())
		return external_models.ValidateAuthorizationDataModel{}, err
	}
	logger.Info("validate authorization", outBoundResponse.Data)

	return outBoundResponse.Data, nil
}
