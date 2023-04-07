package verification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
)

type VerifcationJobModel struct {
	User         external_models.User
	ExtReq       request.ExternalRequest
	Request      models.VerifyIdCardRequest
	Strategies   []string
	Verification models.Verification
	Headers      map[string]string
	Country      external_models.Country
	UserProfile  external_models.UserProfile
	Db           postgresql.Databases
}

func (v *VerifcationJobModel) VerificationJob() (int, error) {
	if v == nil {
		v.ExtReq.Logger.Info((http.StatusInternalServerError), "verification job model is empty")
		return http.StatusInternalServerError, fmt.Errorf("verification job model is empty")
	}
	var (
		appruveReq = AppruveReq{
			ID:           v.Request.ID,
			FirstName:    v.User.Firstname,
			LastName:     v.User.Lastname,
			MiddleName:   v.User.Middlename,
			Phone_number: v.User.PhoneNumber,
			CountryCode:  v.Country.CountryCode,
			DateOfBirth:  v.UserProfile.Dob,
		}
	)
	// fmt.Println(fmt.Sprintf("%v, new: %v", appruveReq, *v))

	statusCode, _ := appruveReq.Process(*v)
	fmt.Println("Status code", statusCode)
	statusCode = 500
	if statusCode != http.StatusOK {
		if statusCode >= 500 && statusCode <= 599 {
			code, err := saveVerificationLogs(*v, appruveReq)
			if err != nil {
				return code, err
			}
		} else {
			v.ExtReq.SendExternalRequest(request.VerificationFailedNotification, external_models.VerificationFailedModel{
				AccountID: v.User.AccountID,
				Type:      v.Request.Type,
			})

			err := v.Verification.Delete(v.Db.Verification)
			if err != nil {
				v.ExtReq.Logger.Info((http.StatusInternalServerError), err.Error())
				return http.StatusInternalServerError, err
			}

			verificationDoc := models.VerificationDoc{AccountID: int(v.User.AccountID), Type: v.Request.Type}
			verificationDoc.GetVerificationDocByAccountIDAndType(v.Db.Verification)
			if err == nil {
				err := verificationDoc.Delete(v.Db.Verification)
				if err != nil {
					v.ExtReq.Logger.Info(strconv.Itoa(http.StatusInternalServerError), err.Error())
					return http.StatusInternalServerError, err
				}
			}

		}

	} else {
		v.ExtReq.SendExternalRequest(request.VerificationSuccessfulNotification, external_models.VerificationSuccessfulModel{
			AccountID: v.User.AccountID,
			Type:      v.Request.Type,
		})

		v.Verification.IsVerified = true
		err := v.Verification.UpdateAllFields(v.Db.Verification)
		if err != nil {
			v.ExtReq.Logger.Info(strconv.Itoa(http.StatusInternalServerError), err.Error())
			return http.StatusInternalServerError, err
		}
		verificationLog := models.VerificationLog{
			AccountId: strconv.Itoa(int(v.User.AccountID)),
			Status:    "failed",
			Type:      v.Request.Type,
		}
		_, err = verificationLog.GetVerificationLogByAccountIDAndType(v.Db.Verification)
		if err == nil {
			verificationLog.Status = "success"
			err := verificationLog.UpdateAllFields(v.Db.Verification)
			if err != nil {
				v.ExtReq.Logger.Info(strconv.Itoa(http.StatusInternalServerError), err.Error())
				return http.StatusInternalServerError, err
			}
		}

	}

	return http.StatusOK, nil
}

func saveVerificationLogs(v VerifcationJobModel, appruveReq AppruveReq) (int, error) {
	verificationLog := models.VerificationLog{
		AccountId: strconv.Itoa(int(v.User.AccountID)),
		Status:    "failed",
		Type:      v.Request.Type,
	}
	code, err := verificationLog.GetVerificationLogByAccountIDAndType(v.Db.Verification)
	if err != nil {
		if code == http.StatusInternalServerError {
			v.ExtReq.Logger.Info(strconv.Itoa(code), err.Error())
			return code, err
		}

		jsonByte, _ := json.Marshal(appruveReq)
		verificationLog.Strategy = "appruve"
		verificationLog.Type = v.Request.Type
		verificationLog.AccountId = strconv.Itoa(int(v.User.AccountID))
		verificationLog.Status = "failed"
		verificationLog.Payload = string(jsonByte)
		err := verificationLog.CreateVerificationLog(v.Db.Verification)
		if err != nil {
			v.ExtReq.Logger.Info(strconv.Itoa(http.StatusInternalServerError), err.Error())
			return http.StatusInternalServerError, err
		}
	}

	if verificationLog.Attempts >= 10 {
		v.ExtReq.SendExternalRequest(request.VerificationFailedNotification, external_models.VerificationFailedModel{
			AccountID: v.User.AccountID,
			Type:      v.Request.Type,
		})
		err := v.Verification.Delete(v.Db.Verification)
		if err != nil {
			v.ExtReq.Logger.Info(strconv.Itoa(http.StatusInternalServerError), err.Error())
			return http.StatusInternalServerError, err
		}
		verificationLog.Delete(v.Db.Verification)
	} else {
		verificationLog.Attempts += 1
		verificationLog.UpdateAllFields(v.Db.Verification)
		if err != nil {
			v.ExtReq.Logger.Info(strconv.Itoa(http.StatusInternalServerError), err.Error())
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil

}
