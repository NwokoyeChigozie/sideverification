package migrations

import "github.com/vesicash/verification-ms/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.FailedJob{},
		models.SwitchWay{},
		models.VerificationCode{},
		models.VerificationDoc{},
		models.VerificationLog{},
		models.Verification{},
	}
}
