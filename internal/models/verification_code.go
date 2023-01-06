package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type VerificationCode struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID int       `gorm:"column:account_id; type:int; not null; comment: Account id of the user being verified" json:"account_id"`
	Code      int       `gorm:"column:code; type:int; not null; comment: Verification Code" json:"code"`
	Token     string    `gorm:"column:token; type:varchar(250); not null" json:"token"`
	ExpiresAt string    `gorm:"column:expires_at; type:varchar(250); not null; comment: Time the code lives for.. 15 minutes by default" json:"expires_at"`
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (v *VerificationCode) GetVerificationCodeByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &v, "id = ? ", v.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (v *VerificationCode) CreateVerificationCode(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &v)
	if err != nil {
		return fmt.Errorf("verification code creation failed: %v", err.Error())
	}
	return nil
}

func (v *VerificationCode) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &v)
	return err
}
