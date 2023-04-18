package verification

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/services/verification"
	"github.com/vesicash/verification-ms/utility"
)

// CheckVerification
// FetchUserVerifications
func (base *Controller) CheckVerification(c *gin.Context) {
	var (
		req struct {
			Type string `json:"type"  validate:"required"`
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

	user := models.MyIdentity
	if user == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "error retrieving authenticated user", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	data, code, err := verification.CheckVerificationService(base.ExtReq, base.Logger, req.Type, int(user.AccountID), base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "success", data)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) CheckVerificationApp(c *gin.Context) {
	var (
		req struct {
			AccountID int    `json:"account_id"  validate:"required"  pgvalidate:"exists=auth$users$account_id"`
			Type      string `json:"type"  validate:"required"`
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

	user, err := verification.GetUserWithAccountID(base.ExtReq, req.AccountID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "error retrieving user", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	data, code, err := verification.CheckVerificationService(base.ExtReq, base.Logger, req.Type, int(user.AccountID), base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "success", data)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) FetchUserVerifications(c *gin.Context) {

	user := models.MyIdentity
	if user == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "error retrieving authenticated user", fmt.Errorf("error retrieving authenticated user"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	verification := models.Verification{AccountID: int(user.AccountID)}
	verifications, err := verification.GetAllByAccountID(base.Db.Verification)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "success", verifications)
	c.JSON(http.StatusOK, rd)

}
