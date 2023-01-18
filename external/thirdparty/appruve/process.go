package appruve

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func AppruveVerifyID(logger *utility.Logger, idata interface{}) (int, error) {

	var (
		token            = config.GetConfig().Appruve.AccessToken
		outBoundResponse map[string]interface{}
	)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	fdata, ok := idata.(external_models.AppruveReqModelFirst)
	if !ok {
		logger.Info("appruve_verify_id", idata, "request data format error")
		return http.StatusInternalServerError, fmt.Errorf("request data format error")
	}

	data := external_models.AppruveReqModelMain{
		ID:           fdata.ID,
		FirstName:    fdata.FirstName,
		LastName:     fdata.LastName,
		MiddleName:   fdata.MiddleName,
		Gender:       fdata.Gender,
		Phone_number: fdata.Phone_number,
		DateOfBirth:  fdata.DateOfBirth,
	}
	logger.Info("appruve_verify_id", data)
	endpoint := "/" + strings.ToLower(fdata.CountryCode) + "/" + fdata.Endpoint

	err := external.SendRequest(logger, "third_party", "appruve_verify_id", headers, data, &outBoundResponse, endpoint)
	if err != nil {
		logger.Info("appruve_verify_id", outBoundResponse, err.Error())
		code := http.StatusInternalServerError
		if external.ResponseCode != 0 {
			code = external.ResponseCode
		}
		return code, err
	}
	logger.Info("appruve_verify_id", outBoundResponse)

	return external.ResponseCode, nil
}
