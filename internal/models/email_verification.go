package models

type VerifyEmailRequest struct {
	AccountID    int    `json:"account_id" pgvalidate:"exists=auth$users$account_id"`
	EmailAddress string `json:"email_address" pgvalidate:"exists=auth$users$email_address"`
	Code         int    `json:"code"`
	Token        string `json:"token"`
}
