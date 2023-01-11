package external

import (
	"fmt"

	"github.com/vesicash/verification-ms/internal/config"
)

type RequestObj struct {
	Name         string
	Path         string
	Method       string
	Headers      map[string]string
	SuccessCode  int
	RequestData  interface{}
	DecodeMethod decodemethod
}

type (
	decodemethod string
)

var (
	JsonDecodeMethod    decodemethod = "json"
	PhpSerializerMethod decodemethod = "phpserializer"
)

func FindMicroserviceRequest(name string, headers map[string]string, data interface{}) (RequestObj, error) {
	var (
		config = config.GetConfig()
	)
	switch name {
	case "get_user":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/auth/get_user", config.Microservices.Auth),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "get_access_token":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/auth/get_access_token", config.Microservices.Auth),
			Method:       "GET",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "validate_on_auth":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/auth/validate_on_db", config.Microservices.Auth),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "validate_authorization":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v2/auth/validate_authorization", config.Microservices.Auth),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "send_verification_email":
		return RequestObj{
			Path:         fmt.Sprintf("%v/email/send/email_verification2", config.Microservices.Notification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "send_welcome_email":
		return RequestObj{
			Path:         fmt.Sprintf("%v/email/send/welcome", config.Microservices.Notification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "send_email_verified_notification":
		return RequestObj{
			Path:         fmt.Sprintf("%v/email/send/email_verified", config.Microservices.Notification),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	default:
		return RequestObj{}, fmt.Errorf("request not found")
	}
}
