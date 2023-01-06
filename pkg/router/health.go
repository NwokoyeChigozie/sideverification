package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/verification-ms/pkg/controller/health"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
)

func Health(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases) *gin.Engine {
	health := health.Controller{Db: db, Validator: validator}

	healthUrl := r.Group(fmt.Sprintf("%v/verification", ApiVersion))
	{
		healthUrl.POST("/health", health.Post)
		healthUrl.GET("/health", health.Get)
	}
	return r
}
