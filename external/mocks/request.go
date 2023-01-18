package mocks

import (
	"fmt"

	"github.com/vesicash/verification-ms/external/mocks/appruve_mocks"
	"github.com/vesicash/verification-ms/external/mocks/auth_mocks"
	"github.com/vesicash/verification-ms/external/mocks/monnify_mocks"
	"github.com/vesicash/verification-ms/external/mocks/notification_mocks"
	"github.com/vesicash/verification-ms/external/mocks/rave_mocks"
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
	case "get_user_credential":
		return auth_mocks.GetUserCredential(er.Logger, data)
	case "create_user_credential":
		return auth_mocks.CreateUserCredential(er.Logger, data)
	case "update_user_credential":
		return auth_mocks.UpdateUserCredential(er.Logger, data)
	case "get_user_profile":
		return auth_mocks.GetUserProfile(er.Logger, data)
	case "get_business_profile":
		return auth_mocks.GetBusinessProfile(er.Logger, data)
	case "get_country":
		return auth_mocks.GetCountry(er.Logger, data)
	case "get_bank_details":
		return auth_mocks.GetBankDetails(er.Logger, data)
	case "get_access_token":
		return auth_mocks.GetAccessToken(er.Logger)
	case "validate_on_auth":
		return auth_mocks.ValidateOnAuth(er.Logger, data)
	case "validate_authorization":
		return auth_mocks.ValidateAuthorization(er.Logger, data)
	case "send_verification_email":
		return notification_mocks.SendVerificationEmail(er.Logger, data)
	case "send_welcome_email":
		return notification_mocks.SendWelcomeEmail(er.Logger, data)
	case "send_email_verified_notification":
		return notification_mocks.SendEmailVerifiedNotification(er.Logger, data)
	case "send_sms_to_phone":
		return notification_mocks.SendSendSMSToPhone(er.Logger, data)
	case "monnify_login":
		return monnify_mocks.MonnifyLogin(er.Logger, data)
	case "monnify_match_bvn_details":
		return monnify_mocks.MonnifyMatchBvnDetails(er.Logger, data)
	case "appruve_verify_id":
		return appruve_mocks.AppruveVerifyID(er.Logger, data)
	case "verification_failed_notification":
		return notification_mocks.VerificationFailedNotification(er.Logger, data)
	case "verification_successful_notification":
		return notification_mocks.VerificationSuccessfulNotification(er.Logger, data)
	case "rave_resolve_bank_account":
		return rave_mocks.RaveResolveBankAccount(er.Logger, data)
	default:
		return nil, fmt.Errorf("request not found")
	}
}
