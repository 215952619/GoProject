package util

import (
	"GoProject/global"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/wenlng/go-captcha/captcha"
)

func NewCaptcha() (string, string, string, error) {
	capt := captcha.GetCaptcha()
	dots, base64, thumbBase64, _, err := capt.Generate()
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("generate captcha error")
		return "", "", "", errors.New("generate captcha error")
	} else {
		str, _ := json.Marshal(dots)
		key, _ := AesEncrypt(str, global.DefaultAesSecret)
		return base64, thumbBase64, string(key), nil
	}
}

func CheckCaptcha() bool {
	return true
}
