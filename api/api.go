package api

import (
	"GoProject/oauth"
	"GoProject/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func logon(c *gin.Context) {
	var form logonRequest
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(util.CustomResponse(util.ParamsParseFailedErrorTemplate, err.Error(), nil))
		return
	}
	result, err := form.CheckCode()
	if err != nil {
		c.JSON(util.CustomResponse(util.CaptchaParseFailedErrorTemplate, "", nil))
		return
	}
	if !result {
		c.JSON(util.CustomResponse(util.CaptchaInvalidErrorTemplate, "", nil))
		return
	}

	user, ok := form.CheckPwd()
	if ok != nil {
		c.JSON(util.UnKnowResponse("校验密码失败"))
		return
	}

	jwt, err := user.GenJwt()
	if err != nil {
		c.JSON(util.UnKnowResponse("生成密钥失败"))
		return
	}

	c.JSON(util.SuccessResponse(gin.H{"token": jwt}))
	return
}

func getCode(c *gin.Context) {
	platform := c.Param("platform")
	ip := c.ClientIP()
	switch platform {
	case "github":
		c.JSON(util.SuccessResponse(oauth.Github.RedirectIdpAuthorizeUrl(ip)))
		return
	case "gitee":
		c.JSON(util.SuccessResponse(oauth.Gitee.RedirectIdpAuthorizeUrl(ip)))
		return
	default:
		c.JSON(util.UnKnowResponse("未知的身份提供商"))
	}
}

func ssoRedirect(c *gin.Context) {
	platform := c.Param("platform")
	state := c.Query("state")
	code := c.Query("code")

	switch platform {
	case "github":
		token, _ := oauth.Github.GetToken(c.ClientIP(), state, code)
		//if err != nil {
		//	global.Logger.WithFields(logrus.Fields{
		//		"err":      err,
		//		"platform": "github",
		//	}).Error("sso get idp info error")
		//	c.JSON(util.UnKnowResponse("获取token失败"))
		//	return
		//}

		//userInfo, err := github.GetUserInfo(token)
		//if err != nil {
		//	c.JSON(util.UnKnowResponse("获取用户信息失败"))
		//	return
		//}
		//var user database.User
		//database.DBM.First(&user, "GithubTag=?", userInfo)
		marshal, err := json.Marshal(token)
		if err != nil {
			panic("marshal token panic")
			return
		}
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s?abc=123&token=%s", oauth.Github.RedirectUrl, marshal))
		return
	case "gitee":
		token, err := oauth.Gitee.GetToken(c.ClientIP(), state, code)
		if err != nil {
			c.JSON(util.UnKnowResponse("获取token失败"))
			return
		}

		fmt.Println(token)
		//userInfo, err := github.GetUserInfo(token)
		//if err != nil {
		//	c.JSON(util.UnKnowResponse("获取用户信息失败"))
		//	return
		//}
		//var user database.User
		//database.DBM.First(&user, "GithubTag=?", userInfo)
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s?abc=123", oauth.Gitee.RedirectUrl))
		return
	default:
		c.JSON(util.UnKnowResponse("未知的身份提供者"))
		return
	}
}
