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
		var user *database.User
		user, err := ParseJwt(c)
		if err != nil {
			errString := util.ErrorToErrorString(err)
			if errString.Code == util.TokenEmpty {
				global.Logger.Debug("not logon")
				c.Next()
			} else {
				c.JSON(util.CustomResponse(util.DefaultError(util.InvalidSign), "", nil))
				c.Abort()
			}
		} else {
			if err = database.DBM.First(&user, "id=?", user.ID); err != nil {
				c.JSON(util.CustomResponse(util.DefaultError(util.NotFound), err.Error(), nil))
				c.Abort()
			}
			if user.Status == database.Closed {
				c.JSON(util.CustomResponse(util.DefaultError(util.NotFound), "", nil))
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
		return "", util.DefaultError(util.TokenEmpty)
	}
	return token, nil
}

func ParseJwt(c *gin.Context) (*database.User, error) {
	tokenString, err := GetToken(c)
	if err != nil {
		return nil, err
	}
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
	return claims.User, nil
}
