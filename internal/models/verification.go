package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Verification struct {
	ID                 uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID          int       `gorm:"column:account_id; type:int; not null; comment: Account id of the user being verified" json:"account_id"`
	VerificationCodeId int       `gorm:"column:verification_code_id; type:int; not null; comment: Verification Code Id from the verification codes table" json:"verification_code_id"`
	VerificationDocId  int       `gorm:"column:verification_doc_id; type:int; not null; comment: Verification Code Id from the verification codes table" json:"verification_doc_id"`
	VerificationType   string    `gorm:"column:verification_type; type:varchar(250); not null; comment: Verification Type (email|phone|bvn)" json:"verification_type"`
	IsVerified         bool      `gorm:"column:is_verified; type:bool; default:false;not null; comment: If user is verified (true|false)" json:"is_verified"`
	VerifiedAt         time.Time `gorm:"column:verified_at" json:"verified_at"`
	DeletedAt          time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt          time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	Tries              int       `gorm:"column:tries; type:int; default: 0" json:"tries"`
}

type CheckVerificationServiceResponseModel struct {
	Verified        bool            `json:"verified"`
	VerificationDoc VerificationDoc `json:"verification_doc"`
}

func (v *Verification) CreateVerification(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &v)
	if err != nil {
		return fmt.Errorf("verification creation failed: %v", err.Error())
	}
	return nil
}

func (v *Verification) GetVerificationByAccountIDAndType(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &v, "account_id = ? and LOWER(verification_type) = ?", v.AccountID, strings.ToLower(v.VerificationType))
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (v *Verification) GetVerificationByAccountIDAndTypeAndIsverified(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &v, "account_id = ? and LOWER(verification_type) = ? and is_verified = ?", v.AccountID, strings.ToLower(v.VerificationType), v.IsVerified)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (v *Verification) GetVerificationOnTypeID(db *gorm.DB) (int, error) {
	q := `account_id = ? and  is_verified = ? and 
	(
		LOWER(verification_type) = 'passport' 
		or LOWER(verification_type) = 'nin' 
		or LOWER(verification_type) = 'national_id' 
		or LOWER(verification_type) = 'nationalid' 
		or LOWER(verification_type) = 'drivers_license' 
		or LOWER(verification_type) = 'driverslicense' 
	) 
	`
	err, nilErr := postgresql.SelectOneFromDb(db, &v, q, v.AccountID, v.IsVerified)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (v *Verification) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &v)
	return err
}

func (v *Verification) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &v)
	if err != nil {
		return fmt.Errorf("verification delete failed: %v", err.Error())
	}
	return nil
}

func (v *Verification) GetAllByAccountID(db *gorm.DB) ([]Verification, error) {
	details := []Verification{}
	err := postgresql.SelectAllFromDb(db, "asc", &details, "account_id = ? ", v.AccountID)
	if err != nil {
		return details, err
	}
	return details, nil
}
