package auth

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func GetCountry(logger *utility.Logger, idata interface{}) (external_models.Country, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetCountryResponse
	)

	data, ok := idata.(external_models.GetCountryModel)
	if !ok {
		logger.Info("get country", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get country", data)
	err := external.SendRequest(logger, "service", "get_country", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("get country", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get country", outBoundResponse)

	return outBoundResponse.Data, nil
}
