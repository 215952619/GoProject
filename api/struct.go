package api

import "GoProject/util"

type logonRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"pwd"`
	Dots       int    `json:"dots"`
	ValidCode  string `json:"valid_code"`
}

func (lr *logonRequest) CheckCode() (bool, bool) {
	return util.CheckCaptcha(lr.Dots, lr.ValidCode)
}
