package pages

import (
	"github.com/go-acme/lego/v3/platform/config/env"
	"github.com/haisum/recaptcha"
)

var (
	siteKey = "6LfWhgQTAAAAABh_abnRpp4OgxB8CG8iUJpPE_OR"
	re      = recaptcha.R{
		Secret: recaptchaSecret(),
	}
)

func recaptchaSecret() string {
	return env.GetOrDefaultString("RECAPTCHA_SECRET", "")
}
