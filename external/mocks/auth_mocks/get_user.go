package auth_mocks

import (
	"fmt"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/utility"
)

var (
	User *external_models.User
)

func GetUser(logger *utility.Logger, idata interface{}) (external_models.User, error) {
	_, ok := idata.(external_models.GetUserRequestModel)
	if !ok {
		logger.Error("get user", idata, "request data format error")
		return external_models.User{}, fmt.Errorf("request data format error")
	}

	if User == nil {
		logger.Error("get user", User, "user not provided")
		return external_models.User{}, fmt.Errorf("user not provided")
	}

	logger.Info("get user", User, "user found")
	return *User, nil
}

func SetUserAuthorizationRequiredStatus(logger *utility.Logger, idata interface{}) (bool, error) {

	data, ok := idata.(external_models.SetUserAuthorizationRequiredStatusModel)
	if !ok {
		logger.Error("set user authorization required status", idata, "request data format error")
		return false, fmt.Errorf("request data format error")
	}

	logger.Info("set user authorization required status", data)

	return true, nil
}
