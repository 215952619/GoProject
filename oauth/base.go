package oauth

import (
	"GoProject/util"
	"errors"
	"fmt"
)

type scopeHandler func() string
type codeHandler func(ip string) string
type tokenHandler func(ip string, state string, code string) (interface{}, error)

type Idp struct {
	ClientId             string
	ClientSecret         string
	Platform             string
	AuthorizeUrl         string
	AuthorizeCallbackUrl string
	TokenUrl             string
	RedirectUrl          string
	GetCodeHandler       codeHandler
	GetScopeHandler      scopeHandler
	GetTokenHandler      tokenHandler
}

func (idp *Idp) GetStateKey(ip string) string {
	return fmt.Sprintf("%s__%s__state", ip, idp.Platform)
}

func (idp *Idp) GetTokenKey(ip string) string {
	return fmt.Sprintf("%s__%s__token", ip, idp.Platform)
}

func (idp *Idp) GetState(ip string) string {
	key := idp.GetStateKey(ip)
	data, exists := util.GetCache(key)
	if exists {
		return data.(string)
	} else {
		return ""
	}
}

func (idp *Idp) GetScope() string {
	if idp.GetScopeHandler == nil {
		return defaultGetScopeHandler()
	}
	return idp.GetScopeHandler()
}

//RedirectIdpAuthorizeUrl TODO 自动跳转cors解决之前，先返回url由前端自动跳转
func (idp *Idp) RedirectIdpAuthorizeUrl(ip string) string {
	if idp.GetCodeHandler == nil {
		return defaultGetCodeHandler(idp, ip)
	}
	return idp.GetCodeHandler(ip)
}

func (idp *Idp) GetToken(ip string, state string, code string) (interface{}, error) {
	if idp.GetTokenHandler == nil {
		return nil, errors.New("not implement")
	}
	return idp.GetTokenHandler(ip, state, code)
}

func (idp *Idp) GenerateState(ip string) string {
	key := idp.GetStateKey(ip)
	state := util.RandomString(20)
	util.SetCacheWithDefault(key, state)
	return state
}

func (idp *Idp) RemoveState(ip string) {
	key := idp.GetStateKey(ip)
	util.RemoveCache(key)
}

func defaultGetCodeHandler(idp *Idp, ip string) string {
	state := idp.GenerateState(ip)
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&state=%s", idp.AuthorizeUrl, idp.ClientId, idp.AuthorizeCallbackUrl, idp.GetScope(), state)
}

func defaultGetScopeHandler() string {
	return "user"
}
