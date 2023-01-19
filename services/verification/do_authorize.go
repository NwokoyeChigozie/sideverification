package verification

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gofrs/uuid"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/utility"
)

func DoAuthorizeService(extReq request.ExternalRequest, logger *utility.Logger, ip, device string, db postgresql.Databases, user external_models.User) (models.DoAuthorizeResponse, string, int, error) {

	var (
		score    = 0
		token, _ = uuid.NewV4()
	)
	ipStackInterface, err := extReq.SendExternalRequest(request.IpstackResolveIp, ip)
	if err != nil {
		return models.DoAuthorizeResponse{}, err.Error(), http.StatusInternalServerError, err
	}

	ipStackResponse, ok := ipStackInterface.(external_models.IPStackResolveIPResponse)
	if !ok {
		return models.DoAuthorizeResponse{}, err.Error(), http.StatusInternalServerError, fmt.Errorf("response data format error")
	}

	location := ipStackResponse.City + ", " + ipStackResponse.CountryName
	userProfile, err := GetUserProfileByAccountID(extReq, logger, int(user.AccountID))
	if err != nil {
		return models.DoAuthorizeResponse{AccountId: int(user.AccountID)}, "Authorized", http.StatusOK, nil
	}

	if userProfile.IpAddress != ip {
		score += 1
	}

	r, _ := regexp.Compile(`(?i)` + userProfile.City)
	if !r.MatchString(ipStackResponse.City) {
		score += 1
	}

	r, _ = regexp.Compile(`(?i)` + userProfile.Country)
	if !r.MatchString(ipStackResponse.CountryName) {
		score += 1
	}

	getAuthorize := external_models.GetAuthorizeModel{
		AccountID:  user.AccountID,
		Authorized: true,
		IpAddress:  ip,
		Browser:    device,
	}
	authorize, err := GetAuthorize(extReq, logger, getAuthorize)
	if err != nil {
		createAuthorize := external_models.CreateAuthorizeModel{
			AccountID:  user.AccountID,
			Authorized: false,
			Token:      token.String(),
			IpAddress:  ip,
			Browser:    device,
			Location:   location,
		}

		aut, err := CreateAuthorize(extReq, logger, createAuthorize)
		if err != nil {
			return models.DoAuthorizeResponse{}, err.Error(), http.StatusInternalServerError, err
		}
		authorize = aut
	} else {
		return models.DoAuthorizeResponse{AccountId: int(user.AccountID), AuthorizedAt: authorize.AuthorizedAt}, "Authorized", http.StatusOK, nil
	}

	if score >= 2 {
		updateAuthorize := external_models.UpdateAuthorizeModel{
			ID:         authorize.ID,
			Authorized: true,
			Token:      token.String(),
		}

		aut, err := UpdateAuthorize(extReq, logger, updateAuthorize)
		if err != nil {
			return models.DoAuthorizeResponse{}, err.Error(), http.StatusInternalServerError, err
		}
		authorize = aut

		extReq.SendExternalRequest(request.SendAuthorizedNotification, external_models.AuthorizeNotificationRequest{
			AccountID: int(user.AccountID),
			Token:     token.String(),
			Ip:        ip,
			Location:  location,
			Device:    device,
		})
	}

	statusInterface, err := extReq.SendExternalRequest(request.SetUserAuthorizationRequiredStatus, external_models.SetUserAuthorizationRequiredStatusModel{
		AccountID: user.AccountID,
		Status:    true,
	})
	if err != nil {
		logger.Info(err.Error())
		return models.DoAuthorizeResponse{}, err.Error(), http.StatusInternalServerError, err
	}

	status, ok := statusInterface.(bool)
	if !ok {
		return models.DoAuthorizeResponse{}, "response data format error", http.StatusInternalServerError, fmt.Errorf("response data format error")

	}

	if !status {
		return models.DoAuthorizeResponse{}, "error updating authorization required", http.StatusInternalServerError, fmt.Errorf("error updating authorization required")
	}

	extReq.SendExternalRequest(request.SendAuthorizationNotification, external_models.AuthorizeNotificationRequest{
		AccountID: int(user.AccountID),
		Token:     token.String(),
		Ip:        ip,
		Location:  location,
		Device:    device,
	})

	return models.DoAuthorizeResponse{
		AccountId: int(user.AccountID),
		Token:     token.String(),
	}, "Authorization Required", http.StatusOK, nil
}

func GetAuthorize(extReq request.ExternalRequest, logger *utility.Logger, data external_models.GetAuthorizeModel) (external_models.Authorize, error) {
	authorizeInterface, err := extReq.SendExternalRequest(request.GetAuthorize, data)
	if err != nil {
		logger.Info(err.Error())
		return external_models.Authorize{}, err
	}

	authorizeResp, ok := authorizeInterface.(external_models.GetAuthorizeResponse)
	if !ok {
		return external_models.Authorize{}, fmt.Errorf("response data format error")
	}

	if authorizeResp.Data.ID == 0 {
		return external_models.Authorize{}, fmt.Errorf("authorize not found")
	}

	return authorizeResp.Data, nil
}

func CreateAuthorize(extReq request.ExternalRequest, logger *utility.Logger, data external_models.CreateAuthorizeModel) (external_models.Authorize, error) {
	authorizeInterface, err := extReq.SendExternalRequest(request.CreateAuthorize, data)
	if err != nil {
		logger.Info(err.Error())
		return external_models.Authorize{}, err
	}

	authorizeResp, ok := authorizeInterface.(external_models.GetAuthorizeResponse)
	if !ok {
		return external_models.Authorize{}, fmt.Errorf("response data format error")
	}

	if authorizeResp.Data.ID == 0 {
		return external_models.Authorize{}, fmt.Errorf("authorize not created")
	}

	return authorizeResp.Data, nil
}

func UpdateAuthorize(extReq request.ExternalRequest, logger *utility.Logger, data external_models.UpdateAuthorizeModel) (external_models.Authorize, error) {
	authorizeInterface, err := extReq.SendExternalRequest(request.UpdateAuthorize, data)
	if err != nil {
		logger.Info(err.Error())
		return external_models.Authorize{}, err
	}

	authorizeResp, ok := authorizeInterface.(external_models.GetAuthorizeResponse)
	if !ok {
		return external_models.Authorize{}, fmt.Errorf("response data format error")
	}

	if authorizeResp.Data.ID == 0 {
		return external_models.Authorize{}, fmt.Errorf("authorize not created")
	}

	return authorizeResp.Data, nil
}
