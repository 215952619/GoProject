package util

type ErrorCode int

const (
	Nothing ErrorCode = iota + 1000
	UnKnow
	InBlackList
	InvalidSign
	TimeExpired
	FrequentRequest
	PermissionDenied
	ParamsParseFailed
	NotFound
	NotLogon
	ExpiredPassword
	PasswordWrong
	Upgrading
	NotImplemented
	TokenEmpty
	CaptchaParseFailed
	CaptchaInvalid
)

var errorMessages = map[ErrorCode]string{
	Nothing:            "请求成功",
	UnKnow:             "未知错误",
	InBlackList:        "您已被列入黑名单，请联系站长解除",
	InvalidSign:        "用户凭据解析失败",
	TimeExpired:        "用户凭据已过期",
	FrequentRequest:    "请求太过频繁，请稍后重试",
	PermissionDenied:   "没有权限访问该资源",
	ParamsParseFailed:  "参数错误",
	NotFound:           "未找到指定资源",
	NotLogon:           "未登录",
	ExpiredPassword:    "密码已过期",
	PasswordWrong:      "密码错误",
	Upgrading:          "系统升级中",
	NotImplemented:     "指定功能暂未开放",
	TokenEmpty:         "找不到用户凭据",
	CaptchaParseFailed: "验证码解析失败",
	CaptchaInvalid:     "验证码校对失败，请重试",
}

type ErrorString struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message,omitempty"`
}

func (e *ErrorString) WithMsg(str string) *ErrorString {
	e.Message = str
	return e
}

func (e *ErrorString) Error() string {
	return e.Message
}

func DefaultError(code ErrorCode) *ErrorString {
	return NewError(code, "")
}

func NewError(code ErrorCode, text string) *ErrorString {
	if len(text) <= 0 {
		msg, ok := errorMessages[code]
		if !ok {
			text = "未知错误"
		}
		text = msg
	}
	return &ErrorString{Code: code, Message: text}
}

func UnKnowError(text string) *ErrorString {
	return NewError(UnKnow, text)
}

func ErrorToErrorString(err error) *ErrorString {
	if err == nil {
		return nil
	}
	errorString, ok := err.(*ErrorString)
	if ok {
		return errorString
	}
	return UnKnowError(err.Error())
}
