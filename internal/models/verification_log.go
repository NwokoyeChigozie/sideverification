package models

import "time"

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
