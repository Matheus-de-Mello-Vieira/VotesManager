package controller

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type recaptchaBody struct {
	CaptchaToken string `json:"captcha_token"`
}
type recaptchaResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error_codes"`
}

func captchaMiddleware(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		response, err := getCaptchaTokenFromBody(responseWriter, request)
		if err != nil {
			return
		}

		if !getVerifyCaptcha(*response) {
			writeCaptchaError(responseWriter)
			return
		}

		next(responseWriter, request)
	})
}

func getCaptchaTokenFromBody(responseWriter http.ResponseWriter, request *http.Request) (*recaptchaResponse, error) {
	body := recaptchaBody{}
	err := loadBody(responseWriter, request, &body)
	if err != nil {
		return nil, err
	}

	decodedJSON, err := base64.StdEncoding.DecodeString(body.CaptchaToken)
	if err != nil {
		writeCaptchaError(responseWriter)
		return nil, err
	}

	var decodedRecaptchaBody recaptchaResponse
	err = json.Unmarshal(decodedJSON, &decodedRecaptchaBody)
	if err != nil {
		writeCaptchaError(responseWriter)
		return nil, err
	}

	return &decodedRecaptchaBody, nil
}

func writeCaptchaError(responseWriter http.ResponseWriter) {
	http.Error(responseWriter, "Invalid Captcha", http.StatusUnauthorized)
}

func getVerifyCaptcha(response recaptchaResponse) bool {
	return response.Success && len(response.ErrorCodes) == 0
}
