package external

import (
	"fmt"
)

func FindThirdPartyRequest(name string, headers map[string]string, data interface{}) (RequestObj, error) {
	var (
	// config = config.GetConfig()
	)
	switch name {
	default:
		return RequestObj{}, fmt.Errorf("request not found")
	}
}
