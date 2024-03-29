package migrations

import (
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func RunAllMigrations(db postgresql.Databases) {

	// verification migration
	MigrateModels(db.Verification, AuthMigrationModels())

}

func MigrateModels(db *gorm.DB, models []interface{}) {
	_ = db.AutoMigrate(models...)
}
