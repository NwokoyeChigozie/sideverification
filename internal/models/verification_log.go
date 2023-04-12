package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type VerificationLog struct {
	ID             uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	VerificationId int       `gorm:"column:verification_id; type:int; comment: Verification Code Id from the verification codes table" json:"verification_id"`
	Payload        string    `gorm:"column:payload; type:text; not null" json:"payload"`
	CreatedAt      time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	Strategy       string    `gorm:"column:strategy; type:varchar(250); comment: The service/strategy used in executing the process" json:"strategy"`
	Type           string    `gorm:"column:type; type:varchar(250); comment: The verification type" json:"type"`
	AccountId      string    `gorm:"column:account_id; type:varchar(250)" json:"account_id"`
	Attempts       int       `gorm:"column:attempts; type:int; default: 0; comment: The number of times this verification has been attempted" json:"attempts"`
	Status         string    `gorm:"column:status; type:varchar(250)" json:"status"`
}

func (v *VerificationLog) GetVerificationLogByAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &v, "account_id = ? ", v.AccountId)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (v *VerificationLog) GetVerificationLogByAccountIDAndType(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &v, "account_id = ? and LOWER(type)=?", v.AccountId, strings.ToLower(v.Type))
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (v *VerificationLog) CreateVerificationLog(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &v)
	if err != nil {
		return fmt.Errorf("verification log creation failed: %v", err.Error())
	}
	return nil
}

func (v *VerificationLog) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &v)
	return err
}

func (v *VerificationLog) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &v)
	if err != nil {
		return fmt.Errorf("verification log delete failed: %v", err.Error())
	}
	return nil
}

func (v *VerificationLog) GetAllByStatus(db *gorm.DB) ([]VerificationLog, error) {
	details := []VerificationLog{}
	err := postgresql.SelectAllFromDb(db, "asc", &details, "status = ? ", v.Status)
	if err != nil {
		return details, err
	}
	return details, nil
}
