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

func TestRequestEmailVerifiaction(t *testing.T) {
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

	type requestBody struct {
		AccountID    int    `json:"account_id"`
		EmailAddress string `json:"email_address"`
	}

	tests := []struct {
		Name         string
		RequestBody  requestBody
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK email verification with account id",
			RequestBody: requestBody{
				AccountID: int(accountID),
			},
			ExpectedCode: http.StatusOK,
			Message:      "E-mail verification request completed.",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "OK email verification with email address",
			RequestBody: requestBody{
				EmailAddress: testUser.EmailAddress,
			},
			ExpectedCode: http.StatusOK,
			Message:      "E-mail verification request completed.",
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
		verificationTypeUrl.POST("/email", veri.RequestEmailVerification)

	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/email"}

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

func TestVerifyEmail(t *testing.T) {
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
		verificationType = "email"
	)
	auth_mocks.User = &testUser

	veri := verification.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	r := gin.Default()
	verificationTypeUrl := r.Group(fmt.Sprintf("%v", "v2"))
	{
		verificationTypeUrl.POST("/email", veri.RequestEmailVerification)
		verificationTypeUrl.POST("/email/verify", veri.VerifyEmail)

	}

	reVReq := struct {
		AccountID    int    `json:"account_id"`
		EmailAddress string `json:"email_address"`
	}{
		AccountID:    int(testUser.AccountID),
		EmailAddress: testUser.EmailAddress,
	}

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(reVReq)
	req, err := http.NewRequest(http.MethodPost, "/v2/email", &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	rrdata := tst.ParseResponse(rr)
	fmt.Println(rr.Code, rrdata)
	if rr.Code != 200 {
		t.Fatal(rr)
	}

	vCode := utility.GetRandomNumbersInRange(111111, 999999)
	vToken, err := utility.ShaHash(utility.RandomString(20))
	verification := models.Verification{AccountID: int(testUser.AccountID), VerificationType: verificationType}
	verificationCode := models.VerificationCode{}
	code, err := verification.GetVerificationByAccountIDAndType(db.Verification)
	if err != nil {
		if code == http.StatusInternalServerError {
			t.Fatal(err)
		}
		verification.IsVerified = false
		err := verification.CreateVerification(db.Verification)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		verificationCode = models.VerificationCode{ID: verification.ID}
		code, err := verificationCode.GetVerificationCodeByID(db.Verification)
		if err != nil {
			if code == http.StatusInternalServerError {
				t.Fatal(err)
			}
			verificationCode.AccountID = int(testUser.AccountID)
			verificationCode.Code = vCode
			verificationCode.Token = vToken
			verificationCode.ExpiresAt = strconv.Itoa(int(time.Now().Add(15 * time.Minute).Unix()))
			err := verificationCode.CreateVerificationCode(db.Verification)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	tests := []struct {
		Name         string
		RequestBody  models.VerifyEmailRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK account_id and code",
			RequestBody: models.VerifyEmailRequest{
				AccountID: int(accountID),
				Code:      verificationCode.Code,
			},
			ExpectedCode: http.StatusOK,
			Message:      "E-mail verification successful.",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, {
			Name: "OK account_id and token",
			RequestBody: models.VerifyEmailRequest{
				AccountID: int(accountID),
				Token:     verificationCode.Token,
			},
			ExpectedCode: http.StatusOK,
			Message:      "E-mail verification successful.",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, {
			Name: "OK email address and code",
			RequestBody: models.VerifyEmailRequest{
				EmailAddress: testUser.EmailAddress,
				Code:         verificationCode.Code,
			},
			ExpectedCode: http.StatusOK,
			Message:      "E-mail verification successful.",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, {
			Name: "OK emmail address and code",
			RequestBody: models.VerifyEmailRequest{
				EmailAddress: testUser.EmailAddress,
				Token:        verificationCode.Token,
			},
			ExpectedCode: http.StatusOK,
			Message:      "E-mail verification successful.",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "no email and account id",
			RequestBody: models.VerifyEmailRequest{
				Token: verificationCode.Token,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "no code and email",
			RequestBody: models.VerifyEmailRequest{
				EmailAddress: testUser.EmailAddress,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Name: "both token and code",
			RequestBody: models.VerifyEmailRequest{
				EmailAddress: testUser.EmailAddress,
				Code:         verificationCode.Code,
				Token:        verificationCode.Token,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/email/verify"}

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
