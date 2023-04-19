package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/pkg/controller/verification"
	"github.com/vesicash/verification-ms/pkg/middleware"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/utility"
)

func Verification(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	verification := verification.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	verificationUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		verificationUrl.POST("/email", verification.RequestEmailVerification)
		verificationUrl.POST("/email/verify", verification.VerifyEmail)

		verificationUrl.POST("/phone", verification.RequestPhoneVerification)
		verificationUrl.POST("/phone/verify", verification.VerifyPhone)

	}

	verificationAuthUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, extReq, middleware.AuthType))
	{
		verificationAuthUrl.POST("/bvn/verify", verification.VerifyBVN)

		verificationAuthUrl.POST("/id/upload", verification.UploadVerificationDoc)
		verificationAuthUrl.POST("/id/verify", verification.VerifyIDCard)
		verificationAuthUrl.POST("/bank_account/verify", verification.VerifyBankAccount)

		verificationAuthUrl.POST("/check", verification.CheckVerification)
		verificationAuthUrl.GET("/fetch", verification.FetchUserVerifications)

		verificationAuthUrl.POST("/authorize", verification.DoAuthorize)
	}

	verificationAppUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, extReq, middleware.AppType))
	{
		verificationAppUrl.POST("/check_verification", verification.CheckVerificationApp)
	}

	verificationjobsUrl := r.Group(fmt.Sprintf("%v/jobs", ApiVersion))
	{
		verificationjobsUrl.POST("/start", verification.StartCronJob)
		verificationjobsUrl.POST("/start-bulk", verification.StartCronJobsBulk)
		verificationjobsUrl.POST("/stop", verification.StopCronJob)
		verificationjobsUrl.PATCH("/update_interval", verification.UpdateCronJobInterval)
	}
	return r
}
