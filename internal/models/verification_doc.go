package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type VerificationDoc struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID int       `gorm:"column:account_id; type:int; not null" json:"account_id"`
	Type      string    `gorm:"column:type; type:varchar(250); not null" json:"type"`
	Value     string    `gorm:"column:value; type:varchar(250)" json:"value"`
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	Meta      string    `gorm:"column:meta; type:varchar(255); comment: Used to save additional information regarding" json:"meta"`
}

func (v *VerificationDoc) GetVerificationDocByAccountIDAndType(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &v, "account_id = ? and type = ?", v.AccountID, v.Type)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (v *VerificationDoc) CreateVerificationDoc(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &v)
	if err != nil {
		return fmt.Errorf("verification doc creation failed: %v", err.Error())
	}
	return nil
}

func (v *VerificationDoc) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &v)
	return err
}
func (v *VerificationDoc) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &v)
	if err != nil {
		return fmt.Errorf("verification doc delete failed: %v", err.Error())
	}
	return nil
}
