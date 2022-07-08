package middleware

import (
	"GoProject/database"
	"GoProject/global"
	"GoProject/util"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Resolve() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("reslove")
		var user *database.User
		user, err := ParseJwt(c)
		if err != nil {
			if err == util.TokenEmptyErrorTemplate.ToError() {
				global.Logger.Debug("not logon")
				c.Next()
			} else {
				c.JSON(util.CustomResponse(util.InvalidSignErrorTemplate, err.Error(), nil))
				c.Abort()
			}
		} else {
			if database.DBM.First(&user, "id=?", user.ID) != nil {
				c.JSON(util.CustomResponse(util.NotFoundErrorTemplate, err.Error(), nil))
				c.Abort()
			}
			if user.Status == database.Closed {
				c.JSON(util.CustomResponse(util.NotFoundErrorTemplate, "", nil))
				c.Abort()
			}

			c.Set(global.AuthedKey, user)
			c.Next()
		}
	}
}

func GetToken(c *gin.Context) (string, error) {
	const minLength = len("Bearer ")

	token := c.Query(global.TokenQueryKey)
	if len(token) == 0 {
		token = c.GetHeader(global.TokenHeaderKey)
	}

	if len(token) <= minLength {
		return "", util.TokenEmptyErrorTemplate
	}
	return token, nil
}

func ParseJwt(c *gin.Context) (*database.User, error) {
	tokenString, err := GetToken(c)
	if err != nil {
		return nil, err
	}
	//if tokenString == "" {
	//	global.Logger.Debug("tokenString not allow nil")
	//	return nil, util.TokenEmptyErrorTemplate
	//}
	claims := database.UserStdClaims{}
	_, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(global.AppJwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims.User, err
}
