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
			"err": err.Error(),
		}).Error("generate captcha error")
		return "", "", "", errors.New("generate captcha error")
	} else {
		str, _ := json.Marshal(dots)
		key, _ := AesEncrypt(str, []byte(global.AesSecret)[:global.AesLength])
		return base64, thumbBase64, string(key), nil
	}
}

func CheckCaptcha(dots interface{}, secret string) (bool, error) {
	type dot [][]int
	dotMap, ok := dots.(dot)
	if !ok {
		return false, errors.New("校验数据格式错误")
	}
	realDots, err := AesDecrypt([]byte(secret), []byte(global.AesSecret)[:global.AesLength])
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"dots":   dotMap,
			"secret": secret,
			"err":    err.Error(),
		}).Error("Decrypt secret error")
		return false, err
	}
	var realDotMap []captcha.CharDot
	err = json.Unmarshal(realDots, &realDotMap)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"realDots": string(realDots),
			"err":      err.Error(),
		}).Error("dots unmarshal error")
		return false, err
	}

	if len(dotMap) != len(realDotMap) {
		return false, nil
	}

	for index, dot := range realDotMap {
		target := dotMap[index]
		if !captcha.CheckPointDistWithPadding(int64(target[0]), int64(target[1]), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height), 5) {
			return false, nil
		}
	}

	return true, nil
}
