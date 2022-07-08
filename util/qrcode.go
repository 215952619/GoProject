package util

import (
	"errors"
	"github.com/skip2/go-qrcode"
)

func NewQrCode() (png []byte, err error) {
	const context = "test text"
	png, err = qrcode.Encode(context, qrcode.Medium, 256)
	if err != nil {
		//Logger.WithFields(logrus.Fields{
		//	"err": err,
		//}).Error("create qrcode error")
		return nil, errors.New("create qrcode error")
	}
	return
}
