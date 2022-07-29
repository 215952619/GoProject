package user

import (
	"GoProject/database"
	"GoProject/global"
	"GoProject/oauth"
	"GoProject/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "not Implementation"})
}

func login(c *gin.Context, _ *database.User) (data interface{}, err error) {
	var form logonRequest
	err = c.ShouldBind(&form)
	if err != nil {
		return nil, util.DefaultError(util.ParamsParseFailed)
	}
	//return nil, nil
	result, err := form.CheckCode()
	if err != nil {
		return nil, util.DefaultError(util.CaptchaParseFailed)
	}
	if !result {
		return nil, util.DefaultError(util.CaptchaInvalid)
	}

	user, ok := form.CheckPwd()
	if ok != nil {
		return nil, util.UnKnowError("校验密码失败")
	}

	jwt, err := user.GenJwt()
	if err != nil {
		//c.JSON(util.UnKnowResponse("生成密钥失败"))
		return nil, util.UnKnowError("生成密钥失败")
	}
	data = gin.H{"token": jwt}
	return
}

func getCode(c *gin.Context, _ *database.User) (data interface{}, err error) {
	platform := c.Param("platform")
	ip := c.ClientIP()
	switch platform {
	case string(oauth.GithubPlatform):
		return oauth.Github.RedirectIdpAuthorizeUrl(ip), nil
	case string(oauth.GiteePlatform):
		return oauth.Gitee.RedirectIdpAuthorizeUrl(ip), nil
	default:
		return nil, util.UnKnowError("未知的身份提供商")
	}
}

func ssoRedirect(c *gin.Context, _ *database.User) (data interface{}, err error) {
	platform := c.Param("platform")
	state := c.Query("state")
	code := c.Query("code")

	switch platform {
	case string(oauth.GithubPlatform):
		token, err := oauth.Github.GetToken(c.ClientIP(), state, code)
		if err != nil {
			global.Logger.WithFields(logrus.Fields{
				"err":      err,
				"platform": platform,
			}).Error("sso get idp info error")
			return nil, util.UnKnowError("获取token失败")
		}

		var data oauth.GithubTokenResponse
		if err = json.Unmarshal(token.([]byte), &data); err != nil {
			return nil, util.UnKnowError(err.Error())
		}
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s?%s", oauth.Github.RedirectUrl, data.String()))
		return nil, nil
	case string(oauth.GiteePlatform):
		token, err := oauth.Gitee.GetToken(c.ClientIP(), state, code)
		if err != nil {
			global.Logger.WithFields(logrus.Fields{
				"err":      err,
				"platform": platform,
			}).Error("sso get idp info error")
			return nil, util.UnKnowError("获取token失败")
		}

		var data oauth.GiteeTokenResponse
		if err = json.Unmarshal(token.([]byte), &data); err != nil {
			return nil, util.UnKnowError(err.Error())
		}
		//if data.Error != "" {
		//	return nil, util.UnKnowError(data.ErrorDescription)
		//}
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s?%s", oauth.Gitee.RedirectUrl, data.String()))
		return nil, nil
	default:
		return nil, util.UnKnowError("未知的身份提供者")
	}
}

func createUser(c *gin.Context, _ *database.User) (data interface{}, err error) {
	var newUser NewUserRequest
	if err := c.ShouldBind(&newUser); err != nil {
		if err := database.DBM.Create(&newUser); err != nil {
			return nil, util.UnKnowError(err.Error())
		}
	}
	return nil, nil
}

func userList(c *gin.Context, _ *database.User) (data interface{}, err error) {
	var users []database.User
	if err := database.DBM.List(&users); err != nil {
		return nil, util.UnKnowError(err.Error())
	} else {
		return users, nil
	}
}

func userDetail(c *gin.Context, _ *database.User) (data interface{}, err error) {
	id := c.Param("id")
	if len(id) > 0 {
		var user database.User
		if err := database.DBM.First(&user, "id=?", id); err != nil {
			return nil, util.UnKnowError("您查找的用户未找到")
		} else {
			return user, nil
		}
	} else {
		return nil, util.UnKnowError("未提供用户标识")
	}
}
