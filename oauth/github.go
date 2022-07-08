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

type GithubTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

var Github *Idp

func init() {
	Github = &Idp{
		ClientId:             global.GithubClientId,
		ClientSecret:         global.GithubClientSecret,
		Platform:             "github",
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
			"redirect_uri":  Github.RedirectUrl,
		}
		paramsBytes, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}

		resp, err := util.PostRequest(Github.TokenUrl, paramsBytes, func(header *http.Header) {
			header.Set("Content-type", "application/json")
		})
		if err != nil {
			global.Logger.WithFields(logrus.Fields{
				"err":    err,
				"url":    Github.TokenUrl,
				"params": params,
			}).Error("request receive error response")
			return nil, err
		}
		var data GithubTokenResponse
		err = json.Unmarshal(resp, &data)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
}
