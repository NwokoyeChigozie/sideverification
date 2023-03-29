package main

import (
	"fmt"
	"log"

	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/internal/models/migrations"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"

	"github.com/vesicash/verification-ms/utility"

	"github.com/go-playground/validator/v10"
	"github.com/vesicash/verification-ms/pkg/router"
)

func main() {
	logger := utility.NewLogger() //Warning !!!!! Do not recreate this action anywhere on the app

	configuration := config.Setup(logger, "./app")

	postgresql.ConnectToDatabases(logger, configuration.Databases)
	validatorRef := validator.New()
	db := postgresql.Connection()

	if configuration.Databases.Migrate {
		migrations.RunAllMigrations(db)
	}

	r := router.Setup(logger, validatorRef, db, &configuration.App)
	rM := router.SetupMetrics(&configuration.App)

	go func(logger *utility.Logger, metricsPort string) {
		utility.LogAndPrint(logger, fmt.Sprintf("Metric Server is starting at 127.0.0.1:%s", metricsPort))
		err := rM.Run(":" + metricsPort)
		if err != nil {
			log.Fatal(err)
		}
	}(logger, configuration.Server.MetricsPort)

	utility.LogAndPrint(logger, fmt.Sprintf("Server is starting at 127.0.0.1:%s", configuration.Server.Port))
	log.Fatal(r.Run(":" + configuration.Server.Port))
}
