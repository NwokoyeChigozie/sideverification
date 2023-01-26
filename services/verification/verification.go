package verification

import (
	"fmt"
	"net/http"

	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/utility"
)

func CheckVerificationService(extReq request.ExternalRequest, logger *utility.Logger, vType string, accountID int, db postgresql.Databases) (models.CheckVerificationServiceResponseModel, int, error) {
	var (
		verification    = models.Verification{AccountID: accountID, IsVerified: true, VerificationType: vType}
		verificationDoc = models.VerificationDoc{AccountID: accountID, Type: vType}
		verified        = false
	)

	if typesStringerr, status := validateType(vType); !status {
		return models.CheckVerificationServiceResponseModel{}, http.StatusBadRequest, fmt.Errorf(typesStringerr)
	}

	if vType == "id" {
		code, err := verification.GetVerificationOnTypeID(db.Verification)
		if err != nil {
			if code == http.StatusInternalServerError {
				return models.CheckVerificationServiceResponseModel{}, code, err
			}
		} else {
			verified = true
		}

		code, err = verificationDoc.GetVerificationDocOnTypeID(db.Verification)
		if err != nil {
			if code == http.StatusInternalServerError {
				return models.CheckVerificationServiceResponseModel{}, code, err
			}
		}

	} else {
		code, err := verification.GetVerificationByAccountIDAndTypeAndIsverified(db.Verification)
		if err != nil {
			if code == http.StatusInternalServerError {
				return models.CheckVerificationServiceResponseModel{}, code, err
			}
		} else {
			verified = true
		}

		code, err = verificationDoc.GetVerificationDocByAccountIDAndType(db.Verification)
		if err != nil {
			if code == http.StatusInternalServerError {
				return models.CheckVerificationServiceResponseModel{}, code, err
			}
		}
	}

	return models.CheckVerificationServiceResponseModel{
		Verified:        verified,
		VerificationDoc: verificationDoc,
	}, http.StatusOK, nil
}
