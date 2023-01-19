package auth

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func GetAuthorize(logger *utility.Logger, idata interface{}) (external_models.GetAuthorizeResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetAuthorizeResponse
	)

	data, ok := idata.(external_models.GetAuthorizeModel)
	if !ok {
		logger.Info("get authorize", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get user credential", data)
	err := external.SendRequest(logger, "service", "get_authorize", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("get authorize", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("get authorize", outBoundResponse)

	return outBoundResponse, nil
}

func CreateAuthorize(logger *utility.Logger, idata interface{}) (external_models.GetAuthorizeResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetAuthorizeResponse
	)

	data, ok := idata.(external_models.CreateAuthorizeModel)
	if !ok {
		logger.Info("create authorize", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("create authorize", data)
	err := external.SendRequest(logger, "service", "create_authorize", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("create authorize", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("create authorize", outBoundResponse)

	return outBoundResponse, nil
}

func UpdateAuthorize(logger *utility.Logger, idata interface{}) (external_models.GetAuthorizeResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetAuthorizeResponse
	)

	data, ok := idata.(external_models.UpdateAuthorizeModel)
	if !ok {
		logger.Info("update authorize", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("update authorize", data)
	err := external.SendRequest(logger, "service", "update_authorize", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("update authorize", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("update authorize", outBoundResponse)

	return outBoundResponse, nil
}
