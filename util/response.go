package util

type ResponseData struct {
	Code ErrorCode   `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (rd *ResponseData) SetDefault() {
	rd.Msg = "success"
}

func (rd *ResponseData) Merge(data *ResponseData) {
	rd.Code = data.Code
	rd.Msg = data.Msg
	rd.Data = data.Data
}

func SuccessResponse(data interface{}) *ResponseData {
	return &ResponseData{
		Code: Nothing,
		Msg:  NothingError.Error(),
		Data: data,
	}
}
