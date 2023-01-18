package verification

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/utility"
)

func VerifyBankAccountService(extReq request.ExternalRequest, Logger *utility.Logger, req models.VerifyBankAccountRequest, db postgresql.Databases, user external_models.User) (bool, int, error) {
	accountNameInterface, err := extReq.SendExternalRequest(request.RaveResolveBankAccount, external_models.ResolveAccountRequest{
		AccountBank:   req.BankCode,
		AccountNumber: req.AccountNumber,
	})

	if err != nil {
		return false, http.StatusBadRequest, fmt.Errorf("invalid bank code or account number")
	}

	accountName, ok := accountNameInterface.(string)
	if !ok {
		return false, http.StatusInternalServerError, fmt.Errorf("error verifying account")
	}

	if accountName == "" {
		return false, http.StatusInternalServerError, fmt.Errorf("error retrieving account name")
	}

	accountName = strings.Join(strings.Fields(accountName), " ")
	req.AccountName = strings.Join(strings.Fields(req.AccountName), " ")
	if !strings.EqualFold(accountName, req.AccountName) {
		return false, http.StatusBadRequest, fmt.Errorf("account Name does not tally with account number")
	}

	return true, http.StatusOK, nil
}
