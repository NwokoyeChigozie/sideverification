package models

type VerifyIdCardRequest struct {
	Type string `json:"type" validate:"required"`
	ID   string `json:"id" validate:"required"`
	Meta string `json:"meta"`
}
