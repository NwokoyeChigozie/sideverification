package verification

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/microservice/auth"
	"github.com/vesicash/verification-ms/external/microservice/notification"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/utility"
)

func RequestEmailVerificationService(logger *utility.Logger, accountID int, emailAddress string, db postgresql.Databases) (int, error) {
	var (
		user             = external_models.User{}
		verificationType = "email"
	)
	if accountID == 0 && emailAddress == "" {
		return http.StatusBadRequest, fmt.Errorf("enter either account id or email address")
	}

	if accountID != 0 {
		us, err := auth.GetUser(logger, external_models.GetUserRequestModel{AccountID: uint(accountID)})
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if us.ID == 0 {
			return http.StatusInternalServerError, fmt.Errorf("user not found")
		}
		user = us
	} else if emailAddress != "" {
		us, err := auth.GetUser(logger, external_models.GetUserRequestModel{EmailAddress: emailAddress})
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if us.ID == 0 {
			return http.StatusInternalServerError, fmt.Errorf("user not found")
		}
		user = us
	}

	vCode := utility.GetRandomNumbersInRange(111111, 999999)
	vToken, err := utility.ShaHash(utility.RandomString(20))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	verification := models.Verification{}
	verificationCode := models.VerificationCode{}
	code, err := verification.GetVerificationByAccountIDAndType(db.Verification)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}

		verificationCode = models.VerificationCode{
			AccountID: int(user.AccountID),
			Code:      vCode,
			Token:     vToken,
			ExpiresAt: strconv.Itoa(int(time.Now().Add(15 * time.Minute).Unix())),
		}
		err = verificationCode.CreateVerificationCode(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		verification = models.Verification{
			AccountID:          verificationCode.AccountID,
			VerificationCodeId: int(verificationCode.ID),
			VerificationType:   verificationType,
			IsVerified:         false,
		}
		err = verification.CreateVerification(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}

	} else {
		verification.IsVerified = false
		err := verification.UpdateAllFields(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		verificationCode = models.VerificationCode{ID: verification.ID}
		code, err := verificationCode.GetVerificationCodeByID(db.Verification)
		if err != nil {
			return code, err
		}

		verificationCode.AccountID = int(user.AccountID)
		verificationCode.Code = vCode
		verificationCode.Token = vToken
		verificationCode.ExpiresAt = strconv.Itoa(int(time.Now().Add(15 * time.Minute).Unix()))
		err = verificationCode.UpdateAllFields(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}

	}

	notification.SendVerificationEmail(logger, external_models.EmailNotificationRequest{
		EmailAddress: user.EmailAddress,
		AccountId:    user.AccountID,
		Code:         uint(verificationCode.Code),
		Token:        verificationCode.Token,
	})

	return http.StatusOK, nil
}
