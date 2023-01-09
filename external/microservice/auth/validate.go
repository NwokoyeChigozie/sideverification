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
		return false, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("validate on auth", data)
	err := external.SendRequest(logger, "service", "validate_on_auth", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("validate on auth", outBoundResponse, err)
		return false, err
	}
	logger.Info("validate on auth", outBoundResponse)

	return outBoundResponse.Data, nil
}
