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

func RequestPhoneVerificationService(extReq request.ExternalRequest, logger *utility.Logger, req models.RequestPhoneVerificationRequest, db postgresql.Databases) (int, error) {
	var (
		user             = external_models.User{}
		verificationType = "phone"
	)

	if req.AccountID == 0 && req.PhoneNumber == "" {
		return http.StatusBadRequest, fmt.Errorf("enter either account id or phone number")
	}

	if req.AccountID != 0 && req.PhoneNumber != "" {
		return http.StatusBadRequest, fmt.Errorf("enter either account id or phone number")
	}

	if req.PhoneNumber != "" {
		us, err := GetUserWithPhone(extReq, req.PhoneNumber)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("user not found")
		}
		user = us
	} else if req.AccountID != 0 {
		us, err := GetUserWithAccountID(extReq, req.AccountID)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("user not found")
		}
		user = us
	}

	if (req.AccountID != 0 && req.PhoneNumber != "") && (user.AccountID != uint(req.AccountID)) {
		return http.StatusBadRequest, fmt.Errorf("phone number is in use by another customer")
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

		verificationCode = models.VerificationCode{ID: uint(verification.VerificationCodeId)}
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

	verification.VerificationCodeId = int(verificationCode.ID)
	err = verification.UpdateAllFields(db.Verification)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	var phone string
	if req.PhoneNumber != "" {
		ph, status := utility.PhoneValid(req.PhoneNumber, extReq.Test)
		if !status {
			return http.StatusBadRequest, fmt.Errorf("invalid phone number")
		}
		phone = ph
	}

	if phone == "" {
		phone = user.PhoneNumber
	}

	if phone == "" {
		return http.StatusBadRequest, fmt.Errorf("user has no recorded phone number")
	}

	extReq.SendExternalRequest(request.SendSmsToPhone, external_models.SMSToPhoneNotificationRequest{
		AccountId:   user.AccountID,
		Message:     "Hi. Your Vesicash Phone Number Verification Code is: " + strconv.Itoa(vCode),
		PhoneNumber: phone,
	})

	return http.StatusOK, nil
}

func VerifyPhoneService(extReq request.ExternalRequest, logger *utility.Logger, req models.VerifyPhoneRequest, db postgresql.Databases) (int, error) {
	var (
		user             = external_models.User{}
		verificationType = "phone"
	)

	if req.Code == 0 && req.Token == "" {
		return http.StatusBadRequest, fmt.Errorf("enter either code or token")
	}

	if req.Code != 0 && req.Token != "" {
		return http.StatusBadRequest, fmt.Errorf("enter either code or token")
	}

	user, err := GetUserWithAccountID(extReq, req.AccountID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("user not found")
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

	verification := models.Verification{AccountID: int(user.AccountID), VerificationType: verificationType, VerificationCodeId: int(verificationCode.ID)}
	code, err := verification.GetVerificationByAccountIDAndTypeAndCodeID(db.Verification)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}
		return code, fmt.Errorf("invalid code/token")
	}

	verification.IsVerified = true
	verification.VerifiedAt = time.Now()
	err = verification.UpdateAllFields(db.Verification)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	verificationCode.Delete(db.Verification)

	return http.StatusOK, nil
}
