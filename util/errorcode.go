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

var (
	NothingErrorTemplate            = New(Nothing, "请求成功")
	UnKnowErrorTemplate             = New(UnKnow, "未知错误")
	InBlackListErrorTemplate        = New(InBlackList, "您已被列入黑名单，请联系站长解除")
	TokenEmptyErrorTemplate         = New(TokenEmpty, "找不到用户凭据")
	InvalidSignErrorTemplate        = New(InvalidSign, "用户凭据解析失败")
	NotFoundErrorTemplate           = New(NotFound, "未找到指定资源")
	NotLogonErrorTemplate           = New(NotLogon, "未登录")
	PermissionDeniedErrorTemplate   = New(PermissionDenied, "没有权限访问该资源")
	ParamsParseFailedErrorTemplate  = New(ParamsParseFailed, "参数错误")
	CaptchaParseFailedErrorTemplate = New(CaptchaParseFailed, "验证码解析失败")
	CaptchaInvalidErrorTemplate     = New(CaptchaInvalid, "验证码校对失败，请重试")
)

type ErrorString struct {
	Code     ErrorCode `json:"code"`
	Template string    `json:"template,omitempty"`
}

func (e *ErrorString) Error() string {
	return e.Template
}

func (e *ErrorString) ToError() error {
	return e
}

func New(code ErrorCode, text string) *ErrorString {
	return &ErrorString{Code: code, Template: text}
}
