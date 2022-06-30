package api

import "GoProject/util"

type logonRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"pwd"`
	Code       int    `json:"code"`
	ValidCode  string `json:"valid_code"`
}

func (lr *logonRequest) CheckCode() (bool, error) {
	return util.CheckCaptcha(), nil
}
