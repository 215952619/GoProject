package oauth

import (
	"GoProject/global"
	"GoProject/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type GiteeTokenResponse struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	Scope            string `json:"scope,omitempty"`
	CreatedAt        int64  `json:"created_at,omitempty"`
}

func (response *GiteeTokenResponse) String() string {
	if response.Error != "" {
		return fmt.Sprintf("%s=%s", "error_description", response.ErrorDescription)
	}
	return fmt.Sprintf("%s=%s&%s=%s&%s=%s&%s=%s&%s=%d", "access_token", response.AccessToken, "token_type", response.TokenType, "refresh_token", response.RefreshToken, "scope", response.Scope, "created_at", response.CreatedAt)
}

var Gitee *Idp

func init() {
	Gitee = &Idp{
		ClientId:             global.GiteeClientId,
		ClientSecret:         global.GiteeClientSecret,
		Platform:             GiteePlatform,
		AuthorizeUrl:         "https://gitee.com/oauth/authorize",
		AuthorizeCallbackUrl: "http://10.0.7.112:9507/api/user/sso/gitee/redirect",
		TokenUrl:             "https://gitee.com/oauth/token",
		RedirectUrl:          fmt.Sprintf("%s/gitee", global.SsoRedirectUrl),
		GetScopeHandler:      getGiteeScope(),
		GetCodeHandler:       getGiteeCode(),
		GetTokenHandler:      getGiteeToken(),
	}
}

func getGiteeScope() scopeHandler {
	return func() string {
		return "user_info emails"
	}
}

func getGiteeCode() codeHandler {
	return func(ip string) string {
		state := Gitee.GenerateState(ip)
		return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&state=%s&response_type=code", Gitee.AuthorizeUrl, Gitee.ClientId, Gitee.AuthorizeCallbackUrl, Gitee.GetScope(), state)
	}
}

func getGiteeToken() tokenHandler {
	return func(ip string, state string, code string) (interface{}, error) {
		defer func() {
			global.Logger.WithFields(logrus.Fields{
				"ip": ip,
			}).Debug("remove state")
			Gitee.RemoveState(ip)
		}()

		if Gitee.GetState(ip) != state {
			global.Logger.WithFields(logrus.Fields{
				"ip":    ip,
				"state": state,
				"real":  Gitee.GetState(ip),
			}).Errorf("valid state err")
			return nil, errors.New("not invalid state")
		}

		var params = map[string]string{
			"grant_type":    "authorization_code",
			"code":          code,
			"client_id":     Gitee.ClientId,
			"redirect_uri":  Gitee.AuthorizeCallbackUrl,
			"client_secret": Gitee.ClientSecret,
		}
		paramsBytes, err := json.Marshal(params)
		if err != nil {
			global.Logger.WithFields(logrus.Fields{
				"params": params,
				"err":    err,
			}).Errorf("marshal params error")
			return nil, err
		}

		resp, err := util.PostRequest(Gitee.TokenUrl, paramsBytes, func(header *http.Header) {
			header.Set("Content-type", "application/json")
		})
		if err != nil {
			global.Logger.WithFields(logrus.Fields{
				"err":  err,
				"url":  Gitee.TokenUrl,
				"resp": resp,
			}).Error("request receive error response")
			return nil, err
		}

		return resp, nil
	}
}
