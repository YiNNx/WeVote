package captcha

import (
	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/pkg/captcha"
)

var captchaClient captcha.CaptchaClient

func VerifyCaptchaIfCaptchaOpen(recaptchaToken *string) error {
	if !config.C.Captcha.Open {
		return nil
	}
	if recaptchaToken == nil {
		return errors.ErrCaptchaRequired
	}
	ok, err := captchaClient.Verify(*recaptchaToken)
	if err != nil || !ok {
		return errors.ErrCaptchaInvalid
	}
	return nil
}

func InitCaptchaClient() {
	captchaClient = captcha.NewClient(
		config.C.Captcha.RecaptchaSecret,
	)
}
