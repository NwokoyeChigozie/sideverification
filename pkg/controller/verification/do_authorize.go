package verification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/services/verification"
	"github.com/vesicash/verification-ms/utility"
)

func (base *Controller) DoAuthorize(c *gin.Context) {
	var (
		req struct {
			Ip     string `json:"ip"  validate:"required"`
			Device string `json:"device"  validate:"required"`
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

	data, msg, code, err := verification.DoAuthorizeService(base.ExtReq, base.Logger, req.Ip, req.Device, base.Db, *user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", msg, err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, msg, data)
	c.JSON(http.StatusOK, rd)

}
