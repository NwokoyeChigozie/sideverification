package verification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/services/verification"
	"github.com/vesicash/verification-ms/utility"
)

func (base *Controller) VerifyBVN(c *gin.Context) {
	var (
		req struct {
			Bvn string `json:"bvn"  validate:"required"`
			Dob string `json:"dob"  validate:"required"`
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

	status, code, err := verification.VerifyBVNService(base.ExtReq, base.Logger, base.Db, req.Bvn, req.Dob, *user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "BVN verification completed.", gin.H{"verified": status})
	c.JSON(http.StatusOK, rd)

}
