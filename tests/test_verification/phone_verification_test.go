package test_verification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/mocks/auth_mocks"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/internal/models"
	"github.com/vesicash/verification-ms/pkg/controller/verification"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/verification-ms/tests"
	"github.com/vesicash/verification-ms/utility"
)

func TestRequestPhoneVerification(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
		accountID = uint(utility.GetRandomNumbersInRange(1000000000, 9999999999))
		testUser  = external_models.User{
			ID:           uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
			AccountID:    accountID,
			EmailAddress: fmt.Sprintf("testuser%v@qa.team", muuid.String()),
			PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
			AccountType:  "individual",
			Firstname:    "test",
			Lastname:     "user",
			Username:     fmt.Sprintf("test_username%v", muuid.String()),
		}
	)
	auth_mocks.User = &testUser

	veri := verification.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	r := gin.Default()

	tests := []struct {
		Name         string
		RequestBody  models.RequestPhoneVerificationRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK phone verification with account id",
			RequestBody: models.RequestPhoneVerificationRequest{
				AccountID: int(accountID),
			},
			ExpectedCode: http.StatusOK,
			Message:      "Phone number verification requested",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "OK phone verification with phone number",
			RequestBody: models.RequestPhoneVerificationRequest{
				PhoneNumber: testUser.PhoneNumber,
			},
			ExpectedCode: http.StatusOK,
			Message:      "Phone number verification requested",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "both account_id and phone number",
			RequestBody: models.RequestPhoneVerificationRequest{
				AccountID:   int(accountID),
				PhoneNumber: testUser.PhoneNumber,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name:         "no input",
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	verificationTypeUrl := r.Group(fmt.Sprintf("%v", "v2"))
	{
		verificationTypeUrl.POST("/phone", veri.RequestPhoneVerification)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/phone"}

			req, err := http.NewRequest(http.MethodPost, URI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)

			code := int(data["code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}

func TestVerifyPhone(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
		accountID = uint(utility.GetRandomNumbersInRange(1000000000, 9999999999))
		testUser  = external_models.User{
			ID:           uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
			AccountID:    accountID,
			EmailAddress: fmt.Sprintf("testuser%v@qa.team", muuid.String()),
			PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
			AccountType:  "individual",
			Firstname:    "test",
			Lastname:     "user",
			Username:     fmt.Sprintf("test_username%v", muuid.String()),
		}
		verificationType = "phone"
	)
	auth_mocks.User = &testUser

	veri := verification.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	r := gin.Default()
	verificationTypeUrl := r.Group(fmt.Sprintf("%v", "v2"))
	{
		verificationTypeUrl.POST("/phone", veri.RequestPhoneVerification)
		verificationTypeUrl.POST("/phone/verify", veri.VerifyPhone)

	}

	tests := []struct {
		Name         string
		RequestBody  models.VerifyPhoneRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
		Create       bool
		Code         bool
		Token        bool
	}{
		{
			Name: "OK account_id and code",
			RequestBody: models.VerifyPhoneRequest{
				AccountID: int(accountID),
			},
			ExpectedCode: http.StatusOK,
			Message:      "Phone Number Verified",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Create: true,
			Code:   true,
			Token:  false,
		}, {
			Name: "OK account_id and token",
			RequestBody: models.VerifyPhoneRequest{
				AccountID: int(accountID),
			},
			ExpectedCode: http.StatusOK,
			Message:      "Phone Number Verified",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Create: true,
			Code:   false,
			Token:  true,
		},
		{
			Name:         "no account id",
			RequestBody:  models.VerifyPhoneRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Create: false,
			Code:   false,
			Token:  true,
		},
		{
			Name: "no code and email",
			RequestBody: models.VerifyPhoneRequest{
				AccountID: int(accountID),
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Create: false,
			Code:   false,
			Token:  false,
		},
		{
			Name: "both token and code",
			RequestBody: models.VerifyPhoneRequest{
				AccountID: int(accountID),
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Create: false,
			Code:   true,
			Token:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			verification := models.Verification{AccountID: int(testUser.AccountID), VerificationType: verificationType}
			verificationCode := models.VerificationCode{}

			if test.Create {
				vCode := utility.GetRandomNumbersInRange(111111, 999999)
				vToken, _ := utility.ShaHash(utility.RandomString(20))
				verificationCode = models.VerificationCode{
					AccountID: int(testUser.AccountID),
					Code:      vCode,
					Token:     vToken,
					ExpiresAt: strconv.Itoa(int(time.Now().Add(15 * time.Minute).Unix())),
				}
				err := verificationCode.CreateVerificationCode(db.Verification)
				if err != nil {
					t.Fatal(err)
				}
				verification = models.Verification{
					AccountID:          verificationCode.AccountID,
					VerificationCodeId: int(verificationCode.ID),
					VerificationType:   verificationType,
					IsVerified:         false,
				}
				err = verification.CreateVerification(db.Verification)
				if err != nil {
					t.Fatal(err)
				}
			}

			if test.Code {
				test.RequestBody.Code = verificationCode.Code
			}
			if test.Token {
				test.RequestBody.Token = verificationCode.Token
			}

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/phone/verify"}

			req, err := http.NewRequest(http.MethodPost, URI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)

			code := int(data["code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}
