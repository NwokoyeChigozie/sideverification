package monnify

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func MonnifyLogin(logger *utility.Logger, idata interface{}) (string, error) {

	var (
		base64Key        = config.GetConfig().Monnify.MonnifyBase64Key
		outBoundResponse external_models.MonnifyLoginResponse
	)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic " + base64Key,
	}

	err := external.SendRequest(logger, "third_party", "monnify_login", headers, nil, &outBoundResponse)
	if err != nil {
		logger.Info("monnify login", outBoundResponse, err.Error())
		return "", err
	}
	logger.Info("monnify login", outBoundResponse)

	return outBoundResponse.ResponseBody.AccessToken, nil
}

func MonnifyMatchBvnDetails(logger *utility.Logger, idata interface{}) (bool, error) {

	var (
		outBoundResponse external_models.MonnifyMatchBvnDetailsResponse
	)

	token, err := MonnifyLogin(logger, nil)
	if err != nil {
		logger.Info("monnify match bvn details", outBoundResponse, err.Error())
		return false, err
	}

	data, ok := idata.(external_models.MonnifyMatchBvnDetailsReq)
	if !ok {
		logger.Info("monnify match bvn details", idata, "request data format error")
		return false, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	logger.Info("monnify match bvn details", data)
	err = external.SendRequest(logger, "third_party", "monnify_match_bvn_details", headers, data, &outBoundResponse)
	if err != nil {
		logger.Info("monnify match bvn details", outBoundResponse, err.Error())
		return false, err
	}
	logger.Info("monnify match bvn details", outBoundResponse)
	if !outBoundResponse.RequestSuccessful {
		logger.Info("monnify match bvn details", "request not successful: "+outBoundResponse.ResponseMessage)
		return false, fmt.Errorf("request not successful: " + outBoundResponse.ResponseMessage)
	}

	if outBoundResponse.ResponseBody.DateOfBirth != "FULL_MATCH" {
		logger.Info("monnify match bvn details", "bvn does not match date of birth", outBoundResponse.ResponseBody.DateOfBirth)
		return false, fmt.Errorf("bvn does not match date of birth: %v", outBoundResponse.ResponseBody.DateOfBirth)
	}

	return true, nil
}
