package monnify

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

type RequestObj struct {
	Name         string
	Path         string
	Method       string
	SuccessCode  int
	RequestData  interface{}
	DecodeMethod string
	Logger       *utility.Logger
}

var (
	JsonDecodeMethod    string = "json"
	PhpSerializerMethod string = "phpserializer"
)

func (r *RequestObj) getNewSendRequestObject(data interface{}, headers map[string]string, urlprefix string) *external.SendRequestObject {
	return external.GetNewSendRequestObject(r.Logger, r.Name, r.Path, r.Method, urlprefix, r.DecodeMethod, headers, r.SuccessCode, data)
}

func (r *RequestObj) getMonnifyLoginObject() *RequestObj {
	var (
		config = config.GetConfig()
	)
	return &RequestObj{
		Name:         "monnify_login",
		Path:         fmt.Sprintf("%v/api/v1/auth/login", config.Monnify.MonnifyApi),
		Method:       "POST",
		SuccessCode:  200,
		DecodeMethod: JsonDecodeMethod,
		RequestData:  nil,
		Logger:       r.Logger,
	}
}
