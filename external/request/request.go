package request

import (
	"fmt"

	"github.com/vesicash/verification-ms/external/microservice/auth"
	"github.com/vesicash/verification-ms/external/microservice/notification"
	"github.com/vesicash/verification-ms/external/mocks"
	"github.com/vesicash/verification-ms/utility"
)

type ExternalRequest struct {
	Logger *utility.Logger
	Test   bool
}

var (
	GetUserReq            string = "get_user"
	GetAccessTokenReq     string = "get_access_token"
	SendVerificationEmail string = "send_verification_email"
)

func (er ExternalRequest) SendExternalRequest(name string, data interface{}) (interface{}, error) {
	if !er.Test {
		switch name {
		case "get_user":
			return auth.GetUser(er.Logger, data)
		case "get_access_token":
			return auth.GetAccessToken(er.Logger)
		case "send_verification_email":
			return notification.SendVerificationEmail(er.Logger, data)
		case "validate_on_auth":
			return auth.ValidateOnAuth(er.Logger, data)
		default:
			return nil, fmt.Errorf("request not found")
		}

	} else {
		mer := mocks.ExternalRequest{Logger: er.Logger, Test: true}
		return mer.SendExternalRequest(name, data)
	}
}
