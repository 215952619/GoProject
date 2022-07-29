package oauth

import (
	"GoProject/global"
	"GoProject/util"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

type GithubTokenResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorUri         string `json:"error_uri"`
	AccessToken      string `json:"access_token"`
	Scope            string `json:"scope"`
	TokenType        string `json:"token_type"`
}

func (response *GithubTokenResponse) String() string {
	if response.Error != "" {
		return fmt.Sprintf("%s=%s", "error_description", response.ErrorDescription)
	}
	return fmt.Sprintf("%s=%s&%s=%s&%s=%s", "AccessToken", response.AccessToken, "Scope", response.Scope, "TokenType", response.TokenType)
}

var Github *Idp

func init() {
	Github = &Idp{
		ClientId:             global.GithubClientId,
		ClientSecret:         global.GithubClientSecret,
		Platform:             GithubPlatform,
		AuthorizeUrl:         "https://github.com/login/oauth/authorize",
		AuthorizeCallbackUrl: "http://10.0.7.112:9507/api/user/sso/github/redirect",
		TokenUrl:             "https://github.com/login/oauth/access_token",
		RedirectUrl:          fmt.Sprintf("%s/github", global.SsoRedirectUrl),
		GetTokenHandler:      getGithubToken(),
	}
}

func getGithubToken() tokenHandler {
	return func(ip string, state string, code string) (interface{}, error) {
		defer func() {
			global.Logger.WithFields(logrus.Fields{
				"ip": ip,
			}).Debug("remove state")
			Github.RemoveState(ip)
		}()

		if Github.GetState(ip) != state {
			global.Logger.WithFields(logrus.Fields{
				"ip":    ip,
				"state": state,
				"real":  Github.GetState(ip),
			}).Errorf("valid state err")
			return nil, errors.New("not invalid state")
		}

		var params = map[string]string{
			"client_id":     Github.ClientId,
			"client_secret": Github.ClientSecret,
			"code":          code,
			"redirect_uri":  Github.AuthorizeCallbackUrl,
		}

		resp, err := util.PostJsonRequest(Github.TokenUrl, params, nil)
		if err != nil {
			global.Logger.WithFields(logrus.Fields{
				"err":    err,
				"url":    Github.TokenUrl,
				"params": params,
			}).Error("request receive error response")
			return nil, err
		}
		return resp, nil
	}
}
