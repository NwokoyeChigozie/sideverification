package test_verification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/vesicash/verification-ms/external/external_models"
	"github.com/vesicash/verification-ms/external/mocks/auth_mocks"
	"github.com/vesicash/verification-ms/external/request"
	"github.com/vesicash/verification-ms/pkg/controller/verification"
	"github.com/vesicash/verification-ms/pkg/middleware"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/verification-ms/tests"
	"github.com/vesicash/verification-ms/utility"
)

func TestDoAuthorize(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _    = uuid.NewV4()
		token, _    = uuid.NewV4()
		AutToken, _ = uuid.NewV4()
		accountID   = uint(utility.GetRandomNumbersInRange(1000000000, 9999999999))
		testUser    = external_models.User{
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

	auth_mocks.ValidateAuthorizationRes = &external_models.ValidateAuthorizationDataModel{
		Status:  true,
		Message: "authorized",
		Data:    testUser,
	}

	auth_mocks.Authorize = &external_models.Authorize{
		ID:         uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		AccountID:  int(testUser.AccountID),
		Authorized: false,
		Token:      AutToken.String(),
		IpAddress:  "197.333.22.11",
		Browser:    "chrome",
		Os:         "macOs",
	}

	veri := verification.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	r := gin.Default()

	type requestBody struct {
		Ip     string `json:"ip"`
		Device string `json:"device"`
	}

	tests := []struct {
		Name         string
		RequestBody  requestBody
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK do authorize",
			RequestBody: requestBody{
				Ip:     "197.333.22.11",
				Device: "chrome",
			},
			ExpectedCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "no ip",
			RequestBody: requestBody{
				Device: "chrome",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "no device",
			RequestBody: requestBody{
				Ip: "197.333.22.11",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name:         "no input",
			RequestBody:  requestBody{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
	}

	verificationAuthUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, veri.ExtReq, middleware.AuthType))
	{
		verificationAuthUrl.POST("/authorize", veri.DoAuthorize)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/authorize"}

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
