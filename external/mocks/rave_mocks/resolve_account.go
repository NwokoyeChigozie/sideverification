package rave_mocks

import (
	"fmt"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/utility"
)

var (
	AccountName string
)

func RaveResolveBankAccount(logger *utility.Logger, idata interface{}) (external_models.ResolveAccountSuccessResponseData, error) {
	response := external_models.ResolveAccountSuccessResponseData{
		AccountName:   "Account Name",
		AccountNumber: "0983746574",
	}
	_, ok := idata.(external_models.ResolveAccountRequest)
	if !ok {
		logger.Error("rave resolve bank account", idata, "request data format error")
		return response, fmt.Errorf("request data format error")
	}

	if AccountName == "" {
		logger.Error("rave resolve bank account", "account name not provided", AccountName)
		return response, fmt.Errorf("account name not provided")
	}

	logger.Info("rave resolve bank account", AccountName)

	return response, nil
}
