package request

import (
	"fmt"

	"github.com/vesicash/verification-ms/external/microservice/auth"
	"github.com/vesicash/verification-ms/external/microservice/notification"
	"github.com/vesicash/verification-ms/external/mocks"
	rave "github.com/vesicash/verification-ms/external/thirdparty/Rave"
	"github.com/vesicash/verification-ms/external/thirdparty/appruve"
	"github.com/vesicash/verification-ms/external/thirdparty/ipstack"
	"github.com/vesicash/verification-ms/external/thirdparty/monnify"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

type ExternalRequest struct {
	Logger *utility.Logger
	Test   bool
}

var (
	JsonDecodeMethod    string = "json"
	PhpSerializerMethod string = "phpserializer"

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

	RaveResolveBankAccount string = "rave_resolve_bank_account"

	IpstackResolveIp                   string = "ipstack_resolve_ip"
	GetAuthorize                       string = "get_authorize"
	CreateAuthorize                    string = "create_authorize"
	UpdateAuthorize                    string = "update_authorize"
	SendAuthorizedNotification         string = "send_authorized_notification"
	SendAuthorizationNotification      string = "send_authorization_notification"
	SetUserAuthorizationRequiredStatus string = "set_user_authorization_required_status"
)

func (er ExternalRequest) SendExternalRequest(name string, data interface{}) (interface{}, error) {
	var (
		config = config.GetConfig()
	)
	if !er.Test {
		switch name {
		case "get_user":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/get_user", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.GetUser()
		case "get_user_credential":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/get_user_credentials", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.GetUserCredential()
		case "create_user_credential":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/create_user_credentials", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.CreateUserCredential()
		case "update_user_credential":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/update_user_credentials", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.UpdateUserCredential()
		case "get_user_profile":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/get_user_profile", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.GetUserProfile()
		case "get_business_profile":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/get_business_profile", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.GetBusinessProfile()
		case "get_country":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/get_country", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.GetCountry()
		case "get_bank_details":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/get_bank_detail", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.GetBankDetails()
		case "get_access_token":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/get_access_token", config.Microservices.Auth),
				Method:       "GET",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.GetAccessToken()
		case "validate_on_auth":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/validate_on_db", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.ValidateOnAuth()
		case "validate_authorization":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/validate_authorization", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.ValidateAuthorization()
		case "send_verification_email":
			obj := notification.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/send/send_email_verification_mail", config.Microservices.Notification),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.SendVerificationEmail()
		case "send_welcome_email":
			obj := notification.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/send/send_welcome_mail", config.Microservices.Notification),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.SendWelcomeEmail()
		case "send_email_verified_notification":
			obj := notification.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/send/send_email_verified_mail", config.Microservices.Notification),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.SendEmailVerifiedNotification()
		case "send_sms_to_phone":
			obj := notification.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/send/send_sms_to_phone", config.Microservices.Notification),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.SendSendSMSToPhone()
		case "monnify_login":
			obj := monnify.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/api/v1/auth/login", config.Monnify.MonnifyApi),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.MonnifyLogin()
		case "monnify_match_bvn_details":
			obj := monnify.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/api/v1/vas/bvn-details-match", config.Monnify.MonnifyApi),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.MonnifyMatchBvnDetails()
		case "appruve_verify_id":
			obj := appruve.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v1/verifications", config.Appruve.BaseUrl),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.AppruveVerifyID()
		case "verification_failed_notification":
			obj := notification.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/send/send_verification_failed", config.Microservices.Notification),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.VerificationFailedNotification()
		case "verification_successful_notification":
			obj := notification.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/send/send_verification_successful", config.Microservices.Notification),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.VerificationSuccessfulNotification()
		case "rave_resolve_bank_account":
			obj := rave.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v3/accounts/resolve", config.Rave.BaseUrl),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.RaveResolveBankAccount()
		case "ipstack_resolve_ip":
			obj := ipstack.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v", config.IPStack.BaseUrl),
				Method:       "GET",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.IpstackResolveIp()
		case "get_authorize":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/get_authorize", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.GetAuthorize()
		case "create_authorize":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/create_authorize", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.CreateAuthorize()
		case "update_authorize":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/update_authorize", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.UpdateAuthorize()
		case "send_authorized_notification":
			obj := notification.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/send/send_authorized", config.Microservices.Notification),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.SendAuthorizedNotification()
		case "send_authorization_notification":
			obj := notification.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/send/send_authorization", config.Microservices.Notification),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.SendAuthorizationNotification()
		case "set_user_authorization_required_status":
			obj := auth.RequestObj{
				Name:         name,
				Path:         fmt.Sprintf("%v/v2/set_authorization_required", config.Microservices.Auth),
				Method:       "POST",
				SuccessCode:  200,
				DecodeMethod: JsonDecodeMethod,
				RequestData:  data,
				Logger:       er.Logger,
			}
			return obj.SetUserAuthorizationRequiredStatus()
		default:
			return nil, fmt.Errorf("request not found")
		}

	} else {
		mer := mocks.ExternalRequest{Logger: er.Logger, Test: true}
		return mer.SendExternalRequest(name, data)
	}
}
