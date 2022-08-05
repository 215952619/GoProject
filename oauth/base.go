package oauth

import (
	"GoProject/database"
	"GoProject/util"
	"errors"
	"fmt"
)

type Platform string
type StateMap map[string]string
type scopeHandler func() string
type codeHandler func(ip string, refer string) string
type tokenHandler func(ip string, state string, code string) (interface{}, string, error)
type userInfoHandler func(token string) (*database.OauthUser, error)
type TokenResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	AccessToken      string `json:"access_token"`
	Scope            string `json:"scope"`
	TokenType        string `json:"token_type"`
}

func (response *TokenResponse) String() string {
	if response.Error != "" {
		return fmt.Sprintf("%s=%s", "error_description", response.ErrorDescription)
	}
	return fmt.Sprintf("%s=%s&%s=%s&%s=%s", "AccessToken", response.AccessToken, "Scope", response.Scope, "TokenType", response.TokenType)
}

var (
	GithubPlatform Platform = "github"
	GiteePlatform  Platform = "gitee"
)

type Idp struct {
	ClientId             string
	ClientSecret         string
	Platform             Platform
	AuthorizeUrl         string
	AuthorizeCallbackUrl string
	TokenUrl             string
	RedirectUrl          string
	GetCodeHandler       codeHandler
	GetScopeHandler      scopeHandler
	GetTokenHandler      tokenHandler
	GetUserInfoHandler   userInfoHandler
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
		return data.(StateMap)["state"]
	} else {
		return ""
	}
}

func (idp *Idp) GetRefer(ip string) string {
	key := idp.GetStateKey(ip)
	data, exists := util.GetCache(key)
	if exists {
		return data.(StateMap)["refer"]
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
func (idp *Idp) RedirectIdpAuthorizeUrl(ip string, refer string) (string, error) {
	if idp.GetCodeHandler == nil {
		return "", errors.New("not implement")
	}
	return idp.GetCodeHandler(ip, refer), nil
}

func (idp *Idp) GetToken(ip string, state string, code string) (interface{}, string, error) {
	if idp.GetTokenHandler == nil {
		return nil, "", errors.New("not implement")
	}
	return idp.GetTokenHandler(ip, state, code)
}

func (idp *Idp) GetUserInfo(token string) (*database.OauthUser, error) {
	if idp.GetUserInfoHandler == nil {
		return nil, errors.New("not implement")
	}
	if user, err := idp.GetUserInfoHandler(token); err != nil {
		return nil, err
	} else {
		if err = database.DBM.Upsert(&user, user, database.OauthUser{ID: user.ID, Platform: string(idp.Platform)}); err != nil {
			return nil, err
		}
		return user, nil
	}
}

func (idp *Idp) GenerateState(ip string, refer string) string {
	key := idp.GetStateKey(ip)
	state := util.RandomString(20)
	util.SetCacheWithDefault(key, StateMap{"state": state, "refer": refer})
	return state
}

func (idp *Idp) RemoveState(ip string) {
	key := idp.GetStateKey(ip)
	util.RemoveCache(key)
}

func defaultGetScopeHandler() string {
	return "user"
}
