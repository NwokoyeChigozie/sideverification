package auth

import (
	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func GetUser(data external_models.GetUserRequestModel) (external_models.User, error) {
	logger := utility.NewLogger()
	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUserModel
	)

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get user", data)
	err := external.SendRequest(logger, "service", "get_user", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("get user", outBoundResponse, err)
		return outBoundResponse.Data.User, err
	}
	logger.Info("get user", outBoundResponse)

	return outBoundResponse.Data.User, nil
}
