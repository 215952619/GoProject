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
)

var (
	NothingError     = Error(Nothing, "请求成功")
	UnKnowError      = Error(UnKnow, "未知错误")
	InBlackListError = Error(InBlackList, "您已被列入黑名单")
	TokenEmptyError  = Error(TokenEmpty, "找不到用户凭据")
)

func Error(code ErrorCode, text string) error {
	return &errorString{code: code, s: text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	code ErrorCode
	s    string
}

func (e *errorString) Code() ErrorCode {
	return e.code
}

func (e *errorString) Error() string {
	return e.s
}
