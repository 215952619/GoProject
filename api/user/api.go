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
	"net/url"
	"strconv"
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

func userBindCheck(c *gin.Context, _ *database.User) (data interface{}, err error) {
	var requestData CheckRequest
	err = c.ShouldBind(&requestData)
	if err != nil {
		return nil, err
	}

	var oauthUser *database.OauthUser
	if err = database.DBM.First(&oauthUser, "platform=? and id=?", requestData.Platform, requestData.ID); err != nil {
		return nil, err
	}
	var userBind *database.UserBind
	exist, err := database.DBM.Exist(&userBind, "platform_id=?", oauthUser.PlatformID)
	if err != nil {
		return nil, err
	}

	if exist {
		jwt, err := userBind.User.GenJwt()
		if err != nil {
			//c.JSON(util.UnKnowResponse("生成密钥失败"))
			return nil, util.UnKnowError("生成密钥失败")
		}
		return gin.H{"exist": exist, "token": jwt}, err
	} else {
		return gin.H{"exist": exist}, err
	}
}

func userBind(c *gin.Context, _ *database.User) (data interface{}, err error) {
	var response *BindRequest
	err = c.ShouldBind(&response)
	if err != nil {
		return nil, err
	}

	user, err := response.CheckPwd()
	if err != nil {
		return nil, err
	}

	jwt, err := user.GenJwt()
	if err != nil {
		return nil, util.UnKnowError("生成密钥失败")
	}
	data = gin.H{"token": jwt}
	return data, nil
}

func getCode(c *gin.Context, _ *database.User) (data interface{}, err error) {
	platform := c.Param("platform")
	ip := c.ClientIP()
	refer := c.GetHeader("Referer")

	var idp *oauth.Idp
	switch platform {
	case string(oauth.GithubPlatform):
		idp = oauth.Github
	case string(oauth.GiteePlatform):
		idp = oauth.Gitee
	default:
		return nil, util.UnKnowError("未知的身份提供商")
	}
	return idp.RedirectIdpAuthorizeUrl(ip, refer)
}

func ssoRedirect(c *gin.Context, _ *database.User) (data interface{}, err error) {
	platform := c.Param("platform")
	state := c.Query("state")
	code := c.Query("code")

	var idp *oauth.Idp
	switch platform {
	case string(oauth.GithubPlatform):
		idp = oauth.Github
	case string(oauth.GiteePlatform):
		idp = oauth.Gitee
	default:
		return nil, util.UnKnowError("未知的身份提供者")
	}

	token, refer, err := idp.GetToken(c.ClientIP(), state, code)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err":      err,
			"platform": platform,
		}).Error("sso get idp info error")
		return nil, util.UnKnowError("获取token失败")
	}

	var response oauth.TokenResponse
	if err = json.Unmarshal(token.([]byte), &response); err != nil {
		return nil, util.UnKnowError(err.Error())
	}

	userinfo, err := idp.GetUserInfo(response.AccessToken)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err":      err,
			"platform": platform,
		}).Error("获取用户信息失败")
		return nil, util.UnKnowError("获取用户信息失败")
	}

	urlObj, err := url.Parse(fmt.Sprintf(idp.RedirectUrl, refer))
	params := url.Values{}
	params.Set("id", strconv.Itoa(int(userinfo.ID)))
	params.Set("name", userinfo.Name)
	params.Set("avatar", userinfo.AvatarUrl)
	params.Set("email", userinfo.Email)
	params.Set("platform", userinfo.Platform)
	urlObj.RawQuery = params.Encode()

	c.Redirect(http.StatusMovedPermanently, urlObj.String())
	return nil, nil
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
