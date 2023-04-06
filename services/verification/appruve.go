package verification

import (
	"net/http"
	"strings"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/request"
)

type AppruveReq struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MiddleName   string `json:"middle_name"`
	Gender       string `json:"gender"`
	Phone_number string `json:"phone_number"`
	DateOfBirth  string `json:"date_of_birth"`
	CountryCode  string `json:"country_code"`
}

func (a *AppruveReq) Process(verificationObj VerifcationJobModel) (int, string) {
	var (
		endpoint = ""
		vType    = strings.ToUpper(verificationObj.Request.Type)
	)

	switch vType {
	case "DRIVERS_LICENSE":
		endpoint = "driver_license"
	case "VOTER_ID":
		endpoint = "voter"
	case "NATIONAL_ID":
		endpoint = "national_id"
	case "PASSPORT":
		endpoint = "passport"
	case "NIN":
		endpoint = "national_id"
	default:
		return http.StatusNotImplemented, "not implemented"
	}

	statusInterface, err := verificationObj.ExtReq.SendExternalRequest(request.AppruveVerifyId, external_models.AppruveReqModelFirst{
		ID:           a.ID,
		FirstName:    a.FirstName,
		LastName:     a.LastName,
		MiddleName:   a.MiddleName,
		Gender:       a.Gender,
		Phone_number: a.Phone_number,
		DateOfBirth:  a.DateOfBirth,
		CountryCode:  a.CountryCode,
		Endpoint:     endpoint,
	})

	if err != nil {
		verificationObj.ExtReq.Logger.Error("appruve error, ", err.Error())
		return statusInterface.(int), err.Error()
	}

	return statusInterface.(int), ""
}
