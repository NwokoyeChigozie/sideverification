package auth

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func GetBankDetails(logger *utility.Logger, idata interface{}) (external_models.BankDetail, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetBankDetailResponse
	)

	data, ok := idata.(external_models.GetBankDetailModel)
	if !ok {
		logger.Info("get bank detail", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get bank detail", data)
	err := external.SendRequest(logger, "service", "get_bank_details", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("get bank detail", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get bank detail", outBoundResponse)

	return outBoundResponse.Data, nil
}
