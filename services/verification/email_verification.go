package verification

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/utility"
)

func RequestEmailVerificationService(extReq request.ExternalRequest, logger *utility.Logger, accountID int, emailAddress string, db postgresql.Databases) (int, error) {
	var (
		user             = external_models.User{}
		verificationType = "email"
	)
	if accountID == 0 && emailAddress == "" {
		return http.StatusBadRequest, fmt.Errorf("enter either account id or email address")
	}

	if accountID != 0 {
		usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{AccountID: uint(accountID)})
		if err != nil {
			return http.StatusInternalServerError, err
		}

		us, ok := usItf.(external_models.User)
		if !ok {
			return http.StatusInternalServerError, fmt.Errorf("response data format error")
		}

		if us.ID == 0 {
			return http.StatusInternalServerError, fmt.Errorf("user not found")
		}
		user = us
	} else if emailAddress != "" {
		usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{EmailAddress: emailAddress})
		if err != nil {
			return http.StatusInternalServerError, err
		}

		us, ok := usItf.(external_models.User)
		if !ok {
			return http.StatusInternalServerError, fmt.Errorf("response data format error")
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

	verification := models.Verification{AccountID: int(user.AccountID), VerificationType: verificationType}
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

	extReq.SendExternalRequest(request.SendVerificationEmail, external_models.EmailNotificationRequest{
		EmailAddress: user.EmailAddress,
		AccountId:    user.AccountID,
		Code:         uint(verificationCode.Code),
		Token:        verificationCode.Token,
	})

	return http.StatusOK, nil
}

func VerifyEmailService(extReq request.ExternalRequest, logger *utility.Logger, req models.VerifyEmailRequest, db postgresql.Databases) (int, error) {
	var (
		user             = external_models.User{}
		verificationType = "email"
	)
	if req.AccountID == 0 && req.EmailAddress == "" {
		return http.StatusBadRequest, fmt.Errorf("enter either account id or email address")
	}
	if req.Code == 0 && req.Token == "" {
		return http.StatusBadRequest, fmt.Errorf("enter either code or token")
	}

	if req.Code != 0 && req.Token != "" {
		return http.StatusBadRequest, fmt.Errorf("enter either code or token")
	}

	if req.AccountID != 0 {
		usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{AccountID: uint(req.AccountID)})
		if err != nil {
			return http.StatusInternalServerError, err
		}

		us, ok := usItf.(external_models.User)
		if !ok {
			return http.StatusInternalServerError, fmt.Errorf("response data format error")
		}

		if us.ID == 0 {
			return http.StatusInternalServerError, fmt.Errorf("user not found")
		}
		user = us
	} else if req.EmailAddress != "" {
		usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{EmailAddress: req.EmailAddress})
		if err != nil {
			return http.StatusInternalServerError, err
		}

		us, ok := usItf.(external_models.User)
		if !ok {
			return http.StatusInternalServerError, fmt.Errorf("response data format error")
		}

		if us.ID == 0 {
			return http.StatusInternalServerError, fmt.Errorf("user not found")
		}
		user = us
	}

	verification := models.Verification{AccountID: int(user.AccountID), VerificationType: verificationType}
	code, err := verification.GetVerificationByAccountIDAndType(db.Verification)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}
		return code, fmt.Errorf("verification not requested")
	}

	verificationCode := models.VerificationCode{AccountID: int(user.AccountID), Code: req.Code, Token: req.Token}
	if req.Token != "" && req.Code == 0 {
		code, err := verificationCode.GetVerificationCodeByAccountIDAndToken(db.Verification)
		if err != nil {
			if code == http.StatusInternalServerError {
				return code, err
			}
			return code, fmt.Errorf("invalid token")
		}

		parseInt, err := strconv.Atoi(verificationCode.ExpiresAt)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		unixTimeUTC := time.Unix(int64(parseInt), 0)
		if time.Now().After(unixTimeUTC) {
			return http.StatusBadRequest, fmt.Errorf("expired token")
		}

	}

	if req.Code != 0 && req.Token == "" {
		code, err := verificationCode.GetVerificationCodeByAccountIDAndCode(db.Verification)
		if err != nil {
			if code == http.StatusInternalServerError {
				return code, err
			}
			return code, fmt.Errorf("invalid code")
		}

		parseInt, err := strconv.Atoi(verificationCode.ExpiresAt)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		unixTimeUTC := time.Unix(int64(parseInt), 0)
		if time.Now().After(unixTimeUTC) {
			return http.StatusBadRequest, fmt.Errorf("expired code")
		}

	}

	verification.IsVerified = true
	verification.VerifiedAt = time.Now()
	err = verification.UpdateAllFields(db.Verification)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extReq.SendExternalRequest(request.SendWelcomeEmail, external_models.AccountIDRequestModel{
		AccountId: user.AccountID,
	})
	extReq.SendExternalRequest(request.SendEmailVerifiedNotification, external_models.AccountIDRequestModel{
		AccountId: user.AccountID,
	})

	return http.StatusOK, nil
}
