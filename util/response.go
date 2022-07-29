package util

import (
	"net/http"
)

type ResponseData struct {
	Code ErrorCode   `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func DefaultResponse(code ErrorCode) (int, *ResponseData) {
	return CustomResponse(DefaultError(code), "", nil)
}

func SuccessResponse(data interface{}) (int, *ResponseData) {
	return CustomResponse(NewError(Nothing, ""), "", data)
}

func UnKnowResponse(msg string) (int, *ResponseData) {
	return CustomResponse(UnKnowError(""), msg, nil)
}

func CustomResponse(err *ErrorString, msg string, data interface{}) (int, *ResponseData) {
	temp := &ResponseData{
		Code: err.Code,
		Msg:  err.Message,
		Data: data,
	}
	if len(msg) > 0 {
		temp.Msg = msg
	}
	return http.StatusOK, temp
}
