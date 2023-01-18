package auth

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func GetUserCredential(logger *utility.Logger, idata interface{}) (external_models.GetUserCredentialResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUserCredentialResponse
	)

	data, ok := idata.(external_models.GetUserCredentialModel)
	if !ok {
		logger.Info("get user credential", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get user credential", data)
	err := external.SendRequest(logger, "service", "get_user_credential", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("get user credential", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("get user credential", outBoundResponse)

	return outBoundResponse, nil
}

func CreateUserCredential(logger *utility.Logger, idata interface{}) (external_models.GetUserCredentialResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUserCredentialResponse
	)

	data, ok := idata.(external_models.CreateUserCredentialModel)
	if !ok {
		logger.Info("create user credential", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("create user credential", data)
	err := external.SendRequest(logger, "service", "create_user_credential", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("create user credential", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("create user credential", outBoundResponse)

	return outBoundResponse, nil
}

func UpdateUserCredential(logger *utility.Logger, idata interface{}) (external_models.GetUserCredentialResponse, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUserCredentialResponse
	)

	data, ok := idata.(external_models.UpdateUserCredentialModel)
	if !ok {
		logger.Info("update user credential", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("update user credential", data)
	err := external.SendRequest(logger, "service", "update_user_credential", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("update user credential", outBoundResponse, err.Error())
		return outBoundResponse, err
	}
	logger.Info("update user credential", outBoundResponse)

	return outBoundResponse, nil
}
