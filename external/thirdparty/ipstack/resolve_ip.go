package ipstack

import (
	"fmt"

	"github.com/vesicash/verification-ms/external"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/utility"
)

func IpstackResolveIp(logger *utility.Logger, idata interface{}) (external_models.IPStackResolveIPResponse, error) {

	var (
		key              = config.GetConfig().IPStack.Key
		outBoundResponse external_models.IPStackResolveIPResponse
	)

	ip, ok := idata.(string)
	if !ok {
		logger.Info("ipstack resolve ip", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	path := "/" + ip + "?access_key=" + key

	logger.Info("ipstack resolve ip", ip)
	err := external.SendRequest(logger, "third_party", "ipstack_resolve_ip", map[string]string{}, nil, &outBoundResponse, path)
	if err != nil {
		logger.Info("ipstack resolve ip", outBoundResponse, err.Error())
		return outBoundResponse, err
	}

	return outBoundResponse, nil
}
