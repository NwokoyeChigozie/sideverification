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

func CreateUserCredential(accountID int, vType, vData, bvn string, extReq request.ExternalRequest) (external_models.UsersCredential, int, error) {
	userCredential := external_models.UsersCredential{}
	usItf, err := extReq.SendExternalRequest(request.CreateUserCredential, external_models.CreateUserCredentialModel{
		AccountID:          uint(accountID),
		IdentificationType: vType,
		IdentificationData: vData,
		Bvn:                bvn,
	})

	if err != nil {
		return userCredential, http.StatusInternalServerError, err
	}

	usC, ok := usItf.(external_models.GetUserCredentialResponse)
	if !ok {
		return userCredential, http.StatusInternalServerError, fmt.Errorf("response data format error")
	}

	if usC.Code != http.StatusOK {
		return userCredential, http.StatusInternalServerError, fmt.Errorf("create user credential failed: %v", usC.Message)
	}

	return usC.Data, http.StatusOK, nil

}

func UpdateUserCredential(id, accountID int, vType, vData, bvn string, extReq request.ExternalRequest) (external_models.UsersCredential, int, error) {
	userCredential := external_models.UsersCredential{}
	usItf, err := extReq.SendExternalRequest(request.UpdateUserCredential, external_models.UpdateUserCredentialModel{
		ID:                 uint(id),
		AccountID:          uint(accountID),
		IdentificationType: vType,
		IdentificationData: vData,
		Bvn:                bvn,
	})

	if err != nil {
		return userCredential, http.StatusInternalServerError, err
	}

	usC, ok := usItf.(external_models.GetUserCredentialResponse)
	if !ok {
		return userCredential, http.StatusInternalServerError, fmt.Errorf("response data format error")
	}

	if usC.Code != http.StatusOK {
		return userCredential, http.StatusInternalServerError, fmt.Errorf("update user credential failed: %v", usC.Message)
	}

	return usC.Data, http.StatusOK, nil

}

func GetUserCredential(id, accountID int, vType string, extReq request.ExternalRequest) (external_models.UsersCredential, int, error) {
	userCredential := external_models.UsersCredential{}

	usCItf, err := extReq.SendExternalRequest(request.GetUserCredential, external_models.GetUserCredentialModel{
		ID:                 uint(id),
		AccountID:          uint(accountID),
		IdentificationType: vType,
	})

	if err != nil {
		return userCredential, http.StatusInternalServerError, err
	}

	usC, ok := usCItf.(external_models.GetUserCredentialResponse)
	if !ok {
		return userCredential, http.StatusInternalServerError, fmt.Errorf("response data format error")
	}

	if usC.Code != http.StatusOK {
		if usC.Code == http.StatusInternalServerError {
			return userCredential, http.StatusInternalServerError, fmt.Errorf(usC.Message)
		}
		return userCredential, http.StatusBadRequest, fmt.Errorf("user credential not found")
	}

	return usC.Data, http.StatusOK, nil

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
