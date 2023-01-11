package verification

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/verification-ms/utility"
)

func VerifyBVNService(extReq request.ExternalRequest, Logger *utility.Logger, db postgresql.Databases, bvn, dob string, user external_models.User) (bool, int, error) {
	var (
		verificationType = "bvn"
		verification     = models.Verification{AccountID: int(user.AccountID), VerificationType: verificationType}
	)

	code, err := verification.GetVerificationByAccountIDAndType(db.Verification)
	if err != nil {
		if code == http.StatusInternalServerError {
			return false, code, err
		}

		verification.IsVerified = false
		err := verification.CreateVerification(db.Verification)
		if err != nil {
			return false, http.StatusInternalServerError, err
		}
	}

	if err := validateDobFormat(dob); err != nil {
		return false, http.StatusBadRequest, err
	}

	dob, err = utility.FormatDate(dob, "2006-01-02", "02-Jan-2006")
	if err != nil {
		return false, http.StatusInternalServerError, err
	}

	bvnMatchReq := external_models.MonnifyMatchBvnDetailsReq{
		Bvn:         bvn,
		DateOfBirth: dob,
		Name:        getValidatedName(user.Lastname, user.Middlename, user.Firstname),
		MobileNo:    getValidatedPhone(user.PhoneNumber),
	}
	status, err := extReq.SendExternalRequest(request.MonnifyMatchBvnDetails, bvnMatchReq)
	aStatus := status.(bool)
	if !aStatus {
		return false, http.StatusBadRequest, err
	}

	verification.IsVerified = true
	err = verification.UpdateAllFields(db.Verification)
	if !aStatus {
		return false, http.StatusBadRequest, err
	}

	return verification.IsVerified, http.StatusOK, nil

}

func getValidatedPhone(phoneNumber string) string {
	if phoneNumber == "" {
		phoneNumber = fmt.Sprintf("0%v", utility.GetRandomNumbersInRange(7000000000, 9099999999))
	}

	phoneNumber, status := utility.PhoneValid(phoneNumber)
	if !status {
		phoneNumber = fmt.Sprintf("0%v", utility.GetRandomNumbersInRange(7000000000, 9099999999))
	}

	phoneNumber = strings.Replace(phoneNumber, "+", "", 1)
	return phoneNumber
}

func getValidatedName(last, middle, first string) string {
	name := ""
	if last != "" {
		name = last
	}

	if middle != "" {
		if name == "" {
			name = middle
		} else {
			name += " " + middle
		}
	}

	if first != "" {
		if name == "" {
			name = first
		} else {
			name += " " + first
		}
	}

	if name == "" {
		name = "ab"
	}
	return name
}
func validateDobFormat(dob string) error {
	var (
		errMsg   = fmt.Errorf("incorrect date of birth format try this format 1997-03-26")
		dobSlice = strings.Split(dob, "-")
	)

	if len(dobSlice) != 3 {
		return errMsg
	}

	if len(dobSlice[0]) != 4 {
		return errMsg
	}

	if len(dobSlice[1]) != 2 || len(dobSlice[2]) != 2 {
		return errMsg
	}

	_, err := strconv.Atoi(dobSlice[0])
	if err != nil {
		return fmt.Errorf("first value is year, try this format 1997-03-26")
	}

	if len(dobSlice[0]) != 4 {
		return fmt.Errorf("first value is year, try this format 1997-03-26")
	}

	month, err := strconv.Atoi(dobSlice[1])
	if err != nil {
		return fmt.Errorf("first value is month, try this format 1997-03-26")
	}

	if month > 12 || month < 1 {
		return fmt.Errorf("first value is month, try this format 1997-03-26")
	}

	day, err := strconv.Atoi(dobSlice[2])
	if err != nil {
		return fmt.Errorf("first value is day, try this format 1997-03-26")
	}

	if day > 31 || day < 1 {
		return fmt.Errorf("first value is day, try this format 1997-03-26")
	}
	return nil
}
