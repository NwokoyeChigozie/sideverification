package auth_mocks

import (
	"fmt"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/utility"
)

func ValidateOnAuth(logger *utility.Logger, idata interface{}) (bool, error) {

	_, ok := idata.(external_models.ValidateOnDBReq)
	if !ok {
		logger.Info("validate on auth", idata, "request data format error")
		return false, fmt.Errorf("request data format error")
	}

	logger.Info("validate on auth", true)

	return true, nil
}
