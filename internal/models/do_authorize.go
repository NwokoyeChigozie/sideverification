package models

type DoAuthorizeResponse struct {
	AccountId    int    `json:"account_id"`
	Token        string `json:"token"`
	AuthorizedAt string `json:"authorized_at"`
}
