package models

type RequestPhoneVerificationRequest struct {
	AccountID   int    `json:"account_id" pgvalidate:"exists=auth$users$account_id"`
	PhoneNumber string `json:"phone_number"`
}

type VerifyPhoneRequest struct {
	AccountID int    `json:"account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	Code      int    `json:"code"`
	Token     string `json:"token"`
}
