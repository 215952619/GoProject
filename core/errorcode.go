package core

import "errors"

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
)

var (
	NothingError     = errors.New("请求成功")
	UnKnowError      = errors.New("未知错误")
	InBlackListError = errors.New("您已被列入黑名单")
)
