package auth

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func GetUser(logger *utility.Logger, idata interface{}) (external_models.User, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUserModel
	)

	data, ok := idata.(external_models.GetUserRequestModel)
	if !ok {
		logger.Info("get user", idata, "request data format error")
		return outBoundResponse.Data.User, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get user", data)
	err := external.SendRequest(logger, "service", "get_user", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("get user", outBoundResponse, err.Error())
		return outBoundResponse.Data.User, err
	}
	logger.Info("get user", outBoundResponse)

	return outBoundResponse.Data.User, nil
}

func SetUserAuthorizationRequiredStatus(logger *utility.Logger, idata interface{}) (bool, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.SetUserAuthorizationRequiredStatusResponse
	)

	data, ok := idata.(external_models.SetUserAuthorizationRequiredStatusModel)
	if !ok {
		logger.Info("set user authorization required status", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("set user authorization required status", data)
	err := external.SendRequest(logger, "service", "set_user_authorization_required_status", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("set user authorization required status", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("set user authorization required status", outBoundResponse)

	return outBoundResponse.Data, nil
}
