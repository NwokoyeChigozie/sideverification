package main

import (
	"log"

	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/internal/models/migrations"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"

	"github.com/vesicash/verification-ms/utility"

	"github.com/go-playground/validator/v10"
	"github.com/vesicash/verification-ms/pkg/router"
)

func init() {
	config := config.Setup("./app")
	postgresql.ConnectToDatabases(config.Databases)

}

func main() {
	//Load config
	getConfig := config.GetConfig()
	validatorRef := validator.New()
	db := postgresql.Connection()

	if getConfig.Databases.Migrate {
		migrations.RunAllMigrations(db)
	}

	r := router.Setup(validatorRef, db)

	utility.LogAndPrint("Server is starting at 127.0.0.1:%s", getConfig.Server.Port)
	log.Fatal(r.Run(":" + getConfig.Server.Port))
}
