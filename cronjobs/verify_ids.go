package cronjobs

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/services/verification"
)

func VerifyIDs(extReq request.ExternalRequest, db postgresql.Databases) {
	verificationLog := models.VerificationLog{Status: "failed"}
	verificationLogs, err := verificationLog.GetAllByStatus(db.Verification)
	if err != nil {
		extReq.Logger.Error("verify by id cronjob error: ", err.Error())
	}

	extReq.Logger.Info(fmt.Sprintf("Running Autoverify Scheduler for %v records", len(verificationLogs)))

	for _, log := range verificationLogs {
		var (
			accountID, _    = strconv.Atoi(log.AccountId)
			VerificationJob = verification.VerifcationJobModel{
				ExtReq: extReq,
				Db:     db,
			}
			appruveReq verification.AppruveReq
		)

		verification := models.Verification{AccountID: accountID, VerificationType: log.Type}
		_, err := verification.GetVerificationByAccountIDAndType(db.Verification)
		if err == nil {
			err = json.Unmarshal([]byte(log.Payload), &appruveReq)
			if err == nil {
				VerificationJob.User = external_models.User{
					AccountID:   uint(accountID),
					Firstname:   appruveReq.FirstName,
					Lastname:    appruveReq.LastName,
					Middlename:  appruveReq.MiddleName,
					PhoneNumber: appruveReq.Phone_number,
				}

				VerificationJob.UserProfile = external_models.UserProfile{
					Dob: appruveReq.DateOfBirth,
				}

				VerificationJob.Country = external_models.Country{
					CountryCode: appruveReq.CountryCode,
				}
				VerificationJob.Request = models.VerifyIdCardRequest{
					ID:   appruveReq.ID,
					Type: log.Type,
				}
				VerificationJob.Verification = verification

				VerificationJob.VerificationJob()
			}
		}

	}

}
