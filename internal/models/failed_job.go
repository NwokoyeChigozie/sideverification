package models

import "time"

type FailedJob struct {
	ID         uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	Connection string    `gorm:"column:connection; type:text; not null" json:"connection"`
	Queue      string    `gorm:"column:queue; type:text; not null" json:"queue"`
	Payload    string    `gorm:"column:payload; type:text; not null" json:"payload"`
	Exception  string    `gorm:"column:exception; type:text; not null" json:"exception"`
	FailedAt   time.Time `gorm:"column:failed_at; autoCreateTime" json:"failed_at"`
}
