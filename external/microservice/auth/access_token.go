package auth

import (
	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func GetAccessToken(logger *utility.Logger) (external_models.AccessToken, error) {
	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetAccessTokenModel
	)

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	err := external.SendRequest(logger, "service", "get_access_token", headers, nil, &outBoundResponse)
	if err != nil {
		logger.Info("get access_token", outBoundResponse, err)
		return outBoundResponse.Data, err
	}
	logger.Info("get access_token", outBoundResponse)

	return outBoundResponse.Data, nil
}
