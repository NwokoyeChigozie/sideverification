package models

import "time"

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
