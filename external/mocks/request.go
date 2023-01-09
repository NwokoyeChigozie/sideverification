package mocks

import (
	"fmt"

	"github.com/vesicash/verification-ms/external/mocks/auth_mocks"
	"github.com/vesicash/verification-ms/external/mocks/notification_mocks"
	"github.com/vesicash/verification-ms/utility"
)

type ExternalRequest struct {
	Logger *utility.Logger
	Test   bool
}

func (er ExternalRequest) SendExternalRequest(name string, data interface{}) (interface{}, error) {
	switch name {
	case "get_user":
		return auth_mocks.GetUser(er.Logger, data)
	case "get_access_token":
		return auth_mocks.GetAccessToken(er.Logger)
	case "validate_on_auth":
		return auth_mocks.ValidateOnAuth(er.Logger, data)
	case "send_verification_email":
		return notification_mocks.SendVerificationEmail(er.Logger, data)
	default:
		return nil, fmt.Errorf("request not found")
	}
}
