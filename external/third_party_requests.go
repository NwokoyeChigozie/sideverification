package external

import (
	"fmt"

	"github.com/vesicash/verification-ms/internal/config"
)

func FindThirdPartyRequest(name string, headers map[string]string, data interface{}) (RequestObj, error) {
	var (
		config = config.GetConfig()
	)
	switch name {
	case "monnify_login":
		return RequestObj{
			Path:         fmt.Sprintf("%v/api/v1/auth/login", config.Monnify.MonnifyApi),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "monnify_match_bvn_details":
		return RequestObj{
			Path:         fmt.Sprintf("%v/api/v1/vas/bvn-details-match", config.Monnify.MonnifyApi),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "appruve_verify_id":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v1/verifications", config.Appruve.BaseUrl),
			Method:       "POST",
			Headers:      headers,
			SuccessCode:  200,
			RequestData:  data,
			DecodeMethod: JsonDecodeMethod,
		}, nil
	case "rave_resolve_bank_account":
		return RequestObj{
			Path:         fmt.Sprintf("%v/v3/accounts/resolve", config.Rave.BaseUrl),
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
