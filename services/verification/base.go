package verification

import (
	"fmt"
	"net/http"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/utility"
)

func CreateUserCredential(accountID int, vType, vData, bvn string, extReq request.ExternalRequest) (external_models.UsersCredential, int, error) {
	userCredential := external_models.UsersCredential{}
	usItf, err := extReq.SendExternalRequest(request.CreateUserCredential, external_models.CreateUserCredentialModel{
		AccountID:          uint(accountID),
		IdentificationType: vType,
		IdentificationData: vData,
		Bvn:                bvn,
	})

	if err != nil {
		return userCredential, http.StatusInternalServerError, err
	}

	usC, ok := usItf.(external_models.GetUserCredentialResponse)
	if !ok {
		return userCredential, http.StatusInternalServerError, fmt.Errorf("response data format error")
	}

	if usC.Code != http.StatusOK {
		return userCredential, http.StatusInternalServerError, fmt.Errorf("create user credential failed: %v", usC.Message)
	}

	return usC.Data, http.StatusOK, nil

}

func UpdateUserCredential(id, accountID int, vType, vData, bvn string, extReq request.ExternalRequest) (external_models.UsersCredential, int, error) {
	userCredential := external_models.UsersCredential{}
	usItf, err := extReq.SendExternalRequest(request.UpdateUserCredential, external_models.UpdateUserCredentialModel{
		ID:                 uint(id),
		AccountID:          uint(accountID),
		IdentificationType: vType,
		IdentificationData: vData,
		Bvn:                bvn,
	})

	if err != nil {
		return userCredential, http.StatusInternalServerError, err
	}

	usC, ok := usItf.(external_models.GetUserCredentialResponse)
	if !ok {
		return userCredential, http.StatusInternalServerError, fmt.Errorf("response data format error")
	}

	if usC.Code != http.StatusOK {
		return userCredential, http.StatusInternalServerError, fmt.Errorf("update user credential failed: %v", usC.Message)
	}

	return usC.Data, http.StatusOK, nil

}

func GetUserCredential(id, accountID int, vType string, extReq request.ExternalRequest) (external_models.UsersCredential, int, error) {
	userCredential := external_models.UsersCredential{}

	usCItf, err := extReq.SendExternalRequest(request.GetUserCredential, external_models.GetUserCredentialModel{
		ID:                 uint(id),
		AccountID:          uint(accountID),
		IdentificationType: vType,
	})

	if err != nil {
		return userCredential, http.StatusInternalServerError, err
	}

	usC, ok := usCItf.(external_models.GetUserCredentialResponse)
	if !ok {
		return userCredential, http.StatusInternalServerError, fmt.Errorf("response data format error")
	}

	if usC.Code != http.StatusOK {
		if usC.Code == http.StatusInternalServerError {
			return userCredential, http.StatusInternalServerError, fmt.Errorf(usC.Message)
		}
		return userCredential, http.StatusBadRequest, fmt.Errorf("user credential not found")
	}

	return usC.Data, http.StatusOK, nil

}

func GetAuthorize(extReq request.ExternalRequest, logger *utility.Logger, data external_models.GetAuthorizeModel) (external_models.Authorize, error) {
	authorizeInterface, err := extReq.SendExternalRequest(request.GetAuthorize, data)
	if err != nil {
		logger.Error(err.Error())
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
		logger.Error(err.Error())
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
		logger.Error(err.Error())
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

func GetCountryByCountryAndCurrency(extReq request.ExternalRequest, logger *utility.Logger, countryCode, currencyCode string) (external_models.Country, error) {

	countryInterface, err := extReq.SendExternalRequest(request.GetCountry, external_models.GetCountryModel{
		CountryCode:  countryCode,
		CurrencyCode: currencyCode,
	})

	if err != nil {
		logger.Error(err.Error())
		return external_models.Country{}, fmt.Errorf("your country could not be resolved, please update your profile")
	}
	country, ok := countryInterface.(external_models.Country)
	if !ok {
		return external_models.Country{}, fmt.Errorf("response data format error")
	}
	if country.ID == 0 {
		return external_models.Country{}, fmt.Errorf("your country could not be resolved, please update your profile")
	}

	return country, nil
}
func GetUserProfileByAccountID(extReq request.ExternalRequest, logger *utility.Logger, accountID int) (external_models.UserProfile, error) {
	userProfileInterface, err := extReq.SendExternalRequest(request.GetUserProfile, external_models.GetUserProfileModel{
		AccountID: uint(accountID),
	})
	if err != nil {
		logger.Error(err.Error())
		return external_models.UserProfile{}, err
	}

	userProfile, ok := userProfileInterface.(external_models.UserProfile)
	if !ok {
		return external_models.UserProfile{}, fmt.Errorf("response data format error")
	}

	if userProfile.ID == 0 {
		return external_models.UserProfile{}, fmt.Errorf("user profile not found")
	}

	return userProfile, nil

}
func GetBusinessProfileByAccountID(extReq request.ExternalRequest, logger *utility.Logger, accountID int) (external_models.BusinessProfile, error) {
	businessProfileInterface, err := extReq.SendExternalRequest(request.GetBusinessProfile, external_models.GetBusinessProfileModel{
		AccountID: uint(accountID),
	})
	if err != nil {
		logger.Error(err.Error())
		return external_models.BusinessProfile{}, fmt.Errorf("business lacks a profile")
	}

	businessProfile, ok := businessProfileInterface.(external_models.BusinessProfile)
	if !ok {
		return external_models.BusinessProfile{}, fmt.Errorf("response data format error")
	}

	if businessProfile.ID == 0 {
		return external_models.BusinessProfile{}, fmt.Errorf("business lacks a profile")
	}
	return businessProfile, nil
}

func GetBankDetailsByAccountID(extReq request.ExternalRequest, logger *utility.Logger, accountID int) (external_models.BankDetail, error) {
	bankDetailsInterface, err := extReq.SendExternalRequest(request.GetBankDetails, external_models.GetBankDetailModel{
		AccountID: uint(accountID),
	})
	if err != nil {
		logger.Error(err.Error())
		return external_models.BankDetail{}, fmt.Errorf("you are yet to add your bank account details")
	}

	bankDetail, ok := bankDetailsInterface.(external_models.BankDetail)
	if !ok {
		return external_models.BankDetail{}, fmt.Errorf("response data format error")
	}

	if bankDetail.ID == 0 {
		return external_models.BankDetail{}, fmt.Errorf("you are yet to add your bank account details")
	}
	return bankDetail, nil
}
func GetUserWithAccountID(extReq request.ExternalRequest, accountID int) (external_models.User, error) {
	usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{AccountID: uint(accountID)})
	if err != nil {
		return external_models.User{}, err
	}

	us, ok := usItf.(external_models.User)
	if !ok {
		return external_models.User{}, fmt.Errorf("response data format error")
	}

	if us.ID == 0 {
		return external_models.User{}, fmt.Errorf("user not found")
	}
	return us, nil
}
func GetUserWithEmail(extReq request.ExternalRequest, email string) (external_models.User, error) {
	usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{EmailAddress: email})
	if err != nil {
		return external_models.User{}, err
	}

	us, ok := usItf.(external_models.User)
	if !ok {
		return external_models.User{}, fmt.Errorf("response data format error")
	}

	if us.ID == 0 {
		return external_models.User{}, fmt.Errorf("user not found")
	}
	return us, nil
}
func GetUserWithPhone(extReq request.ExternalRequest, phone string) (external_models.User, error) {
	usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{PhoneNumber: phone})
	if err != nil {
		return external_models.User{}, err
	}

	us, ok := usItf.(external_models.User)
	if !ok {
		return external_models.User{}, fmt.Errorf("response data format error")
	}

	if us.ID == 0 {
		return external_models.User{}, fmt.Errorf("user not found")
	}
	return us, nil
}
