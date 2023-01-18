package request

import (
	"fmt"

	"github.com/vesicash/verification-ms/external/microservice/auth"
	"github.com/vesicash/verification-ms/external/microservice/notification"
	"github.com/vesicash/verification-ms/external/mocks"
	"github.com/vesicash/verification-ms/external/thirdparty/appruve"
	"github.com/vesicash/verification-ms/external/thirdparty/monnify"
	"github.com/vesicash/verification-ms/utility"
)

type ExternalRequest struct {
	Logger *utility.Logger
	Test   bool
}

var (

	// microservice
	GetUserReq           string = "get_user"
	GetUserCredential    string = "get_user_credential"
	CreateUserCredential string = "create_user_credential"
	UpdateUserCredential string = "update_user_credential"

	GetUserProfile     string = "get_user_profile"
	GetBusinessProfile string = "get_business_profile"
	GetCountry         string = "get_country"
	GetBankDetails     string = "get_bank_details"

	GetAccessTokenReq             string = "get_access_token"
	ValidateOnAuth                string = "validate_on_auth"
	ValidateAuthorization         string = "validate_authorization"
	SendVerificationEmail         string = "send_verification_email"
	SendWelcomeEmail              string = "send_welcome_email"
	SendEmailVerifiedNotification string = "send_email_verified_notification"
	SendSmsToPhone                string = "send_sms_to_phone"

	// third party
	MonnifyLogin           string = "monnify_login"
	MonnifyMatchBvnDetails string = "monnify_match_bvn_details"

	AppruveVerifyId string = "appruve_verify_id"

	VerificationFailedNotification     string = "verification_failed_notification"
	VerificationSuccessfulNotification string = "verification_successful_notification"
)

func (er ExternalRequest) SendExternalRequest(name string, data interface{}) (interface{}, error) {
	if !er.Test {
		switch name {
		case "get_user":
			return auth.GetUser(er.Logger, data)
		case "get_user_credential":
			return auth.GetUserCredential(er.Logger, data)
		case "create_user_credential":
			return auth.CreateUserCredential(er.Logger, data)
		case "update_user_credential":
			return auth.UpdateUserCredential(er.Logger, data)
		case "get_user_profile":
			return auth.GetUserProfile(er.Logger, data)
		case "get_business_profile":
			return auth.GetBusinessProfile(er.Logger, data)
		case "get_country":
			return auth.GetCountry(er.Logger, data)
		case "get_bank_details":
			return auth.GetBankDetails(er.Logger, data)
		case "get_access_token":
			return auth.GetAccessToken(er.Logger)
		case "validate_on_auth":
			return auth.ValidateOnAuth(er.Logger, data)
		case "validate_authorization":
			return auth.ValidateAuthorization(er.Logger, data)
		case "send_verification_email":
			return notification.SendVerificationEmail(er.Logger, data)
		case "send_welcome_email":
			return notification.SendWelcomeEmail(er.Logger, data)
		case "send_email_verified_notification":
			return notification.SendEmailVerifiedNotification(er.Logger, data)
		case "send_sms_to_phone":
			return notification.SendSendSMSToPhone(er.Logger, data)
		case "monnify_login":
			return monnify.MonnifyLogin(er.Logger, data)
		case "monnify_match_bvn_details":
			return monnify.MonnifyMatchBvnDetails(er.Logger, data)
		case "appruve_verify_id":
			return appruve.AppruveVerifyID(er.Logger, data)
		case "verification_failed_notification":
			return notification.VerificationFailedNotification(er.Logger, data)
		case "verification_successful_notification":
			return notification.VerificationSuccessfulNotification(er.Logger, data)
		default:
			return nil, fmt.Errorf("request not found")
		}

	} else {
		mer := mocks.ExternalRequest{Logger: er.Logger, Test: true}
		return mer.SendExternalRequest(name, data)
	}
}
