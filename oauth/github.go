package oauth

import (
	"GoProject/database"
	"GoProject/global"
	"GoProject/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

var Github *Idp

func init() {
	go func() {
		select {
		case <-global.InitChan:
			Github = &Idp{
				ClientId:             global.GithubClientId,
				ClientSecret:         global.GithubClientSecret,
				Platform:             GithubPlatform,
				AuthorizeUrl:         "https://github.com/login/oauth/authorize",
				AuthorizeCallbackUrl: fmt.Sprintf("%sapi/user/sso/github/redirect", global.Local),
				TokenUrl:             "https://github.com/login/oauth/access_token",
				RedirectUrl:          "%s/github",
				GetScopeHandler:      getGithubScope(),
				GetCodeHandler:       getGithubCode(),
				GetTokenHandler:      getGithubToken(),
				GetUserInfoHandler:   getGithubUserInfo(),
			}
		}
		defer func() {
			global.InitChan <- true
		}()
	}()
}

func getGithubToken() tokenHandler {
	return func(ip string, state string, code string) (interface{}, string, error) {
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
			return nil, "", errors.New("not invalid state")
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
			return nil, "", err
		}
		return resp, Github.GetRefer(ip), nil
	}
}

func getGithubCode() codeHandler {
	return func(ip string, refer string) string {
		state := Github.GenerateState(ip, refer)
		return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&state=%s", Github.AuthorizeUrl, Github.ClientId, Github.AuthorizeCallbackUrl, Github.GetScope(), state)

	}
}

func getGithubScope() scopeHandler {
	return func() string {
		return "user"
	}
}

func getGithubUserInfo() userInfoHandler {
	return func(token string) (*database.OauthUser, error) {
		res, err := util.GetRequest("https://api.github.com/user", nil, map[string]string{"Authorization": fmt.Sprintf("token %s", token)})

		if err != nil {
			return nil, err
		}

		var data *database.OauthUser
		if err = json.Unmarshal(res, &data); err != nil {
			return nil, err
		}

		return data, nil
	}
}
