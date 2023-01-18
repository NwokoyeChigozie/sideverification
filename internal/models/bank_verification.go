package models

type VerifyBankAccountRequest struct {
	BankCode      string `json:"bank_code" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
}
