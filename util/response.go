package util

import "net/http"

type ResponseData struct {
	Code ErrorCode   `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func SuccessResponse(data interface{}) (int, *ResponseData) {
	return CustomResponse(NothingErrorTemplate, "", data)
}

func UnKnowResponse(msg string) (int, *ResponseData) {
	return CustomResponse(UnKnowErrorTemplate, msg, nil)
}

func CustomResponse(err *ErrorString, msg string, data interface{}) (int, *ResponseData) {
	temp := &ResponseData{
		Code: err.Code,
		Msg:  err.Template,
		Data: data,
	}
	if len(msg) > 0 {
		temp.Msg = msg
	}
	return http.StatusOK, temp
}
