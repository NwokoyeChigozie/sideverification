package rave

import (
	"fmt"
	"strings"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
)

func (r *RequestObj) RaveResolveBankAccount() (external_models.ResolveAccountSuccessResponseData, error) {

	var (
		outBoundResponse external_models.ResolveAccountSuccessResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + config.GetConfig().Rave.SecretKey,
	}

	data, ok := idata.(external_models.ResolveAccountRequest)
	if !ok {
		logger.Error("rave resolve bank account", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("rave resolve bank account", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	if !strings.EqualFold(outBoundResponse.Status, "success") {
		logger.Error("rave resolve bank account", outBoundResponse, err.Error())
		return outBoundResponse.Data, fmt.Errorf("verification failed")
	}
	logger.Info("rave resolve bank account", outBoundResponse)

	return outBoundResponse.Data, nil
}
