package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/elliotchance/phpserialize"
	"github.com/vesicash/verification-ms/utility"
)

var (
	ResponseCode int
	ResponseBody string
)

func SendRequest(logger *utility.Logger, reqType, name string, headers map[string]string, data interface{}, response interface{}, urlPrefix ...string) error {
	var (
		reqObject = RequestObj{}
		err       error
	)
	if reqType == "service" {
		reqObject, err = FindMicroserviceRequest(name, headers, data)
	} else if reqType == "third_party" {
		reqObject, err = FindThirdPartyRequest(name, headers, data)
	} else {
		err = fmt.Errorf("not implemented")
	}

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(data)
	if err != nil {
		logger.Error("encoding error", name, err.Error())
	}

	logger.Info(name, reqObject.Path, data, buf)
	if len(urlPrefix) > 0 {
		reqObject.Path += urlPrefix[0]
	}

	client := &http.Client{}
	req, err := http.NewRequest(reqObject.Method, reqObject.Path, buf)
	if err != nil {
		logger.Error("request creation error", name, err.Error())
		return err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	logger.Error("request", name, reqObject.Path, reqObject.Method, reqObject.Headers)

	res, err := client.Do(req)
	if err != nil {
		logger.Error("client do", name, err.Error())
		return err
	}

	if reqObject.DecodeMethod != PhpSerializerMethod {
		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			logger.Error("json decoding error", name, err.Error())
			return err
		}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("readin body error", name, err.Error())
		return err
	}

	logger.Info("response body", name, reqObject.Path, string(body))

	if reqObject.DecodeMethod == PhpSerializerMethod {
		err := phpserialize.Unmarshal(body, response)
		if err != nil {
			logger.Error("php serializer decoding error", name, err.Error())
			return err
		}
	}

	defer res.Body.Close()
	ResponseCode = res.StatusCode

	if res.StatusCode == reqObject.SuccessCode {
		return nil
	}

	if res.StatusCode < 200 && res.StatusCode > 299 {
		return fmt.Errorf("Error " + strconv.Itoa(res.StatusCode))
	}

	return nil
}
