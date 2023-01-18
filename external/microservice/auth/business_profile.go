package auth

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func GetBusinessProfile(logger *utility.Logger, idata interface{}) (external_models.BusinessProfile, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetBusinessProfileResponse
	)

	data, ok := idata.(external_models.GetBusinessProfileModel)
	if !ok {
		logger.Info("get business profile", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get business profile", data)
	err := external.SendRequest(logger, "service", "get_business_profile", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("get business profile", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get business profile", outBoundResponse)

	return outBoundResponse.Data, nil
}
