package verification

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/utility"
)

func VerifyIDCardService(extReq request.ExternalRequest, logger *utility.Logger, req models.VerifyIdCardRequest, db postgresql.Databases, user external_models.User) (int, error) {
	if typesStringerr, status := validateType(req.Type); !status {
		return http.StatusBadRequest, fmt.Errorf(typesStringerr)
	}

	if req.Type == "nationalid" {
		req.Type = "national_id"
	}
	if req.Type == "driverslicense" {
		req.Type = "drivers_license"
	}

	verification := models.Verification{AccountID: int(user.AccountID), VerificationType: req.Type}
	code, err := verification.GetVerificationByAccountIDAndType(db.Verification)
	foundVerification := false
	if err == nil {
		foundVerification = true
		if verification.IsVerified {
			return http.StatusOK, nil
		}
	} else {
		if code == http.StatusInternalServerError {
			return code, err
		}
		foundVerification = false
	}

	if req.Type == "bvn" {
		_, err := GetBankDetailsByAccountID(extReq, logger, int(user.AccountID))
		if err != nil {
			return http.StatusBadRequest, err
		}

	}

	if req.Type == "cac" {
		_, err := GetBusinessProfileByAccountID(extReq, logger, int(user.AccountID))
		if err != nil {
			return http.StatusBadRequest, err
		}
	}

	userProfile, err := GetUserProfileByAccountID(extReq, logger, int(user.AccountID))
	if err != nil {
		return http.StatusBadRequest, err
	}

	country, err := GetCountryByCountryAndCurrency(extReq, logger, userProfile.Country, userProfile.Currency)
	if err != nil {
		return http.StatusBadRequest, err
	}
	regex := GetCountryVerificationTypeRegex(country.CountryCode, req.Type)
	fmt.Println("regex", regex)
	if regex != "" {
		if !validateID(req.ID, req.Type, regex) {
			return http.StatusBadRequest, fmt.Errorf("Invalid ID value supplied, make sure it's the correct value format for %v in your country %v, check your card to verify", req.Type, country.Name)
		}
	}

	strategy, status := getCountryVerificationStrategy(country.CountryCode, req.Type)
	if !status {
		return http.StatusBadRequest, fmt.Errorf("The verification type %v is not available for your country %v", req.Type, country.Name)
	}

	verificationDocs := models.VerificationDoc{AccountID: int(user.AccountID), Type: req.Type}
	code, err = verificationDocs.GetVerificationDocByAccountIDAndType(db.Verification)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}
		verificationDocs.Value = req.ID
		if req.Type == "cac" {
			verificationDocs.Meta = req.Meta
		}
		err := verificationDocs.CreateVerificationDoc(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		verificationDocs.Type = req.Type
		verificationDocs.Value = req.ID
		if req.Type == "cac" {
			verificationDocs.Meta = req.Meta
		}
		err := verification.UpdateAllFields(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	if !foundVerification {
		verification.AccountID = int(user.AccountID)
		verification.VerificationType = req.Type
		verification.VerificationDocId = int(verificationDocs.ID)
		verification.IsVerified = false
		err = verification.CreateVerification(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		verification.AccountID = int(user.AccountID)
		verification.VerificationType = req.Type
		verification.VerificationDocId = int(verificationDocs.ID)
		verification.IsVerified = false
		err := verification.UpdateAllFields(db.Verification)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// Todo: goroutine for verification
	vJobModel := VerifcationJobModel{
		User:         user,
		ExtReq:       extReq,
		Request:      req,
		Strategies:   strategy,
		Verification: verification,
		Headers:      map[string]string{},
		Country:      country,
		UserProfile:  userProfile,
		Db:           db,
	}

	go vJobModel.VerificationJob()

	return http.StatusOK, nil
}

func getCountryVerificationStrategy(countryCode, vType string) ([]string, bool) {
	countryCode, vType = strings.ToUpper(countryCode), strings.ToUpper(vType)
	var (
		smileAndApprove = []string{"smile_identity", "appruve"}
		smile           = []string{"smile_identity"}
		// approve         = []string{"appruve"}
	)
	switch countryCode {
	case "NG":
		switch vType {
		case "DRIVERS_LICENSE":
			return smileAndApprove, true
		case "PASSPORT":
			return smileAndApprove, true
		case "NATIONAL_ID":
			return smileAndApprove, true
		case "NIN":
			return smileAndApprove, true
		case "VOTER_ID":
			return smileAndApprove, true
		default:
			return []string{}, false
		}
	case "GH":
		switch vType {
		case "DRIVERS_LICENSE":
			return smileAndApprove, true
		case "PASSPORT":
			return smileAndApprove, true
		case "VOTER_ID":
			return smileAndApprove, true
		default:
			return []string{}, false
		}
	case "KE":
		switch vType {
		case "PASSPORT":
			return smileAndApprove, true
		case "NATIONAL_ID":
			return smileAndApprove, true
		default:
			return []string{}, false
		}
	case "ZA":
		switch vType {
		case "NATIONAL_ID":
			return smile, true
		default:
			return []string{}, false
		}
	default:
		return []string{}, false
	}
}

func validateID(id, vType, regex string) bool {
	if vType == "cac" {
		id = "RC" + id
	}

	r, _ := regexp.Compile(regex)
	return r.MatchString(id)
}

func GetCountryVerificationTypeRegex(countryCode, vType string) string {
	fmt.Println("vv")
	countryCode, vType = strings.ToUpper(countryCode), strings.ToUpper(vType)
	switch countryCode {
	case "NG":
		switch vType {
		case "BVN":
			return `^\d{11}$`
		case "NIN":
			return `^\d{11}$`
		case "DRIVERS_LICENSE":
			return `^\w{3}([ -])?\w{6,12}$`
		case "PASSPORT":
			return `^(?i)[A-Z]{1}([ ])?[0-9]{8}$`
		case "VOTER_ID":
			return `^(?i)[A-Z_0-9 ]{9,23}$`
		case "NATIONAL_ID":
			fmt.Println("vvhere")
			return `^\d{10,11}$`
		case "BANK_ACCOUNT":
			return `^\d{10,}$`
		case "CAC":
			return `^(RC)?[0-9]{5,8}$`
		case "TIN":
			return `^[0-9]{8,}-[0-9]{4,}$`
		default:
			return ""
		}
	case "GH":
		switch vType {
		case "DRIVERS_LICENSE":
			return `^\w{6,10}$`
		case "SSNIT":
			return `^(?i)[A-Z]{1}\w{12,14}$`
		case "VOTER_ID":
			return `^\d{10,12}$`
		case "PASSPORT":
			return `^(?i)G[A-Z_0-9]{7,9}$`
		case "NATIONAL_ID":
			return `^(?i)GHA-[A-Z_0-9]{9}-[A-Z_0-9]{1}$`
		default:
			return ""
		}
	case "KE":
		switch vType {
		case "NATIONAL_ID":
			return `^\d{1,9}$`
		case "PASSPORT":
			return `^[A-Z_0-9]{7,9}$`
		case "ALIEN_CARD":
			return `^\d{6,9}$`
		default:
			return ""
		}
	case "ZA":
		switch vType {
		case "NATIONAL_ID":
			return `^\d{13}$`
		default:
			return ""
		}
	default:
		return ""
	}
}
