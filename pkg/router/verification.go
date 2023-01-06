package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/verification-ms/pkg/controller/verification"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/utility"
)

func Verification(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases, logger *utility.Logger) *gin.Engine {
	verification := verification.Controller{Db: db, Validator: validator, Logger: logger}

	verificationUrl := r.Group(fmt.Sprintf("%v/verification", ApiVersion))
	{
		verificationUrl.POST("/email", verification.RequestEmailVerification)
	}
	return r
}
