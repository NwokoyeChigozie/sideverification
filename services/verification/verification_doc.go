package verification

import (
	"fmt"
	"net/http"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/utility"
)

func UploadVerificationDocService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, vType, data string, user external_models.User) (int, error) {
	if typesStringerr, status := validateType(vType); !status {
		return http.StatusBadRequest, fmt.Errorf(typesStringerr)
	}

	verificationDocs := models.VerificationDoc{AccountID: int(user.AccountID), Type: vType}
	code, err := verificationDocs.GetVerificationDocByAccountIDAndType(db.Verification)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}

		verificationDocs.Value = data
		err := verificationDocs.CreateVerificationDoc(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		verificationDocs.Value = data
		err := verificationDocs.UpdateAllFields(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	_, code, err = GetUserCredential(0, int(user.AccountID), vType, extReq)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}

		_, code, err = CreateUserCredential(int(user.AccountID), vType, data, "", extReq)
		if err != nil {
			return code, err
		}
	}

	verification := models.Verification{AccountID: int(user.AccountID), VerificationType: vType}
	code, err = verification.GetVerificationByAccountIDAndType(db.Verification)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}

		verification.VerificationDocId = int(verificationDocs.ID)
		verification.IsVerified = false
		err := verification.CreateVerification(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		verification.VerificationDocId = int(verificationDocs.ID)
		verification.IsVerified = false
		err := verification.UpdateAllFields(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil

}

func validateType(vType string) (string, bool) {
	types := []string{
		"cac",
		"passport",
		"drivers_license",
		"national_id",
		"bvn",
		"nin",
		"email",
		"phone",
		"cac",
		"passport",
		"nin",
		"id",
		"nationalid",
		"driverslicense",
		"utilitybill",
		"voter_id",
	}
	typesString := "invalid type, try any of these types ("

	for _, v := range types {
		if v == vType {
			return "", true
		}
	}

	for _, v := range types {
		typesString += v + ", "
	}

	typesString += typesString + ")"

	return typesString, false
}
