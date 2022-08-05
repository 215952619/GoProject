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

var Gitee *Idp

func init() {
	go func() {
		select {
		case <-global.InitChan:
			Gitee = &Idp{
				ClientId:             global.GiteeClientId,
				ClientSecret:         global.GiteeClientSecret,
				Platform:             GiteePlatform,
				AuthorizeUrl:         "https://gitee.com/oauth/authorize",
				AuthorizeCallbackUrl: fmt.Sprintf("%sapi/user/sso/gitee/redirect", global.Local),
				TokenUrl:             "https://gitee.com/oauth/token",
				RedirectUrl:          "%s/gitee",
				GetScopeHandler:      getGiteeScope(),
				GetCodeHandler:       getGiteeCode(),
				GetTokenHandler:      getGiteeToken(),
				GetUserInfoHandler:   getGiteeUserInfo(),
			}
		}
		defer func() {
			global.InitChan <- true
		}()
	}()
}

func getGiteeScope() scopeHandler {
	return func() string {
		return "user_info emails"
	}
}

func getGiteeCode() codeHandler {
	return func(ip string, refer string) string {
		state := Gitee.GenerateState(ip, refer)
		return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&state=%s&response_type=code", Gitee.AuthorizeUrl, Gitee.ClientId, Gitee.AuthorizeCallbackUrl, Gitee.GetScope(), state)
	}
}

func getGiteeToken() tokenHandler {
	return func(ip string, state string, code string) (interface{}, string, error) {
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
			return nil, "", errors.New("not invalid state")
		}

		var params = map[string]string{
			"grant_type":    "authorization_code",
			"code":          code,
			"client_id":     Gitee.ClientId,
			"redirect_uri":  Gitee.AuthorizeCallbackUrl,
			"client_secret": Gitee.ClientSecret,
		}

		resp, err := util.PostJsonRequest(Gitee.TokenUrl, params, nil)
		if err != nil {
			global.Logger.WithFields(logrus.Fields{
				"err":  err,
				"url":  Gitee.TokenUrl,
				"resp": resp,
			}).Error("request receive error response")
			return nil, "", err
		}

		return resp, Gitee.GetRefer(ip), nil
	}
}

func getGiteeUserInfo() userInfoHandler {
	return func(token string) (*database.OauthUser, error) {
		res, err := util.GetRequest("https://gitee.com/api/v5/user", map[string]string{"access_token": token}, nil)

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
