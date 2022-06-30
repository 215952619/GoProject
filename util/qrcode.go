package util

import (
	"GoProject/global"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
)

func NewQrCode() (png []byte, err error) {
	const context = "test text"
	png, err = qrcode.Encode(context, qrcode.Medium, 256)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("create qrcode error")
	}
	return
}
