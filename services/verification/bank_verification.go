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

func VerifyBankAccountService(extReq request.ExternalRequest, Logger *utility.Logger, req models.VerifyBankAccountRequest, db postgresql.Databases, user external_models.User) (external_models.ResolveAccountSuccessResponseData, int, error) {
	accountNameInterface, err := extReq.SendExternalRequest(request.RaveResolveBankAccount, external_models.ResolveAccountRequest{
		AccountBank:   req.BankCode,
		AccountNumber: req.AccountNumber,
	})

	if err != nil {
		return external_models.ResolveAccountSuccessResponseData{}, http.StatusBadRequest, fmt.Errorf("invalid bank code or account number")
	}

	accountDetails, ok := accountNameInterface.(external_models.ResolveAccountSuccessResponseData)
	if !ok {
		return accountDetails, http.StatusInternalServerError, fmt.Errorf("error verifying account")
	}

	if accountDetails.AccountName == "" {
		return accountDetails, http.StatusBadRequest, fmt.Errorf("could not retrieve account name")
	}

	if req.AccountName == "" {
		return accountDetails, http.StatusOK, nil
	}

	if !MatchNames(accountDetails.AccountName, req.AccountName) {
		return accountDetails, http.StatusBadRequest, fmt.Errorf("account Name does not tally with account number")
	}

	return accountDetails, http.StatusOK, nil
}

func MatchNames(name1, name2 string) bool {
	name1Slice := strings.Fields(name1)
	name2Slice := strings.Fields(name2)
	matchCount := 0

	if len(name1Slice) != len(name2Slice) {
		return false
	}

	for _, n1 := range name1Slice {
		for _, n2 := range name2Slice {
			if strings.EqualFold(n1, n2) {
				matchCount += 1
			}
		}
	}

	return matchCount == len(name1Slice)
}
