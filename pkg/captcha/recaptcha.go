package captcha

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const ReCaptchaVerifyUrl = "https://recaptcha.net/recaptcha/api/siteverify"

type CaptchaClient interface {
	Verify(string) (bool, error)
}

type ReCaptchaClient struct {
	secret string
}

func (r *ReCaptchaClient) Verify(responseToken string) (bool, error) {
	// secret	Required. The shared key between your site and reCAPTCHA.
	// response	Required. The user response token provided by the reCAPTCHA client-side integration on your site.
	reqData := url.Values{
		"secret":   {r.secret},
		"response": {responseToken},
	}
	resp, err := http.PostForm(ReCaptchaVerifyUrl, reqData)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	type captchaResponse struct {
		Success    bool     `json:"success"`
		ErrorCodes []string `json:"error-codes"`
	}
	captchaResp := &captchaResponse{}
	err = json.Unmarshal(body, captchaResp)
	if err != nil {
		return false, err
	}

	if len(captchaResp.ErrorCodes) > 0 {
		return false, errors.New(strings.Join(captchaResp.ErrorCodes, ","))
	}

	return captchaResp.Success, nil
}

func NewReCaptchaClient(secret string) CaptchaClient {
	return &ReCaptchaClient{secret}
}
