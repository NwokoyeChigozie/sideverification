package verification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/services/verification"
	"github.com/vesicash/verification-ms/utility"
)

func (base *Controller) RequestEmailVerification(c *gin.Context) {
	var (
		req struct {
			AccountID    int    `json:"account_id" pgvalidate:"exists=auth$users$account_id"`
			EmailAddress string `json:"email_address" pgvalidate:"exists=auth$users$email_address"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	vr := postgresql.ValidateRequestM{Logger: base.Logger, Test: base.ExtReq.Test}
	err = vr.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	code, err := verification.RequestEmailVerificationService(base.ExtReq, base.Logger, req.AccountID, req.EmailAddress, base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "E-mail verification request completed.", nil)
	c.JSON(http.StatusOK, rd)

}
