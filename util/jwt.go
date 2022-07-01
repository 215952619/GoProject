package util

import (
	"GoProject/global"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func GetToken(c *gin.Context) (string, error) {
	const minLength = len("Bearer ")

	token := c.Query(global.TokenQueryKey)
	if len(token) == 0 {
		token = c.GetHeader(global.TokenHeaderKey)
	}

	if len(token) <= minLength {
		return "", errors.New("not match a valid token")
	}
	return token, nil
}

func GenerateJwt(u *global.User) (string, error) {
	expiredTime := time.Now().Add(global.DefaultValidityPeriod)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expiredTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        string(u.ID),
		Issuer:    global.AppIssuer,
	}
	uClaims := global.UserStdClaims{StandardClaims: stdClaims, User: u}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uClaims)
	tokenString, err := token.SignedString(global.AppSecret)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err":  err,
			"user": u,
		}).Error("generate jwt error")
		return "", err
	}
	return tokenString, nil
}

func ParseJwt(c *gin.Context) (*global.User, error) {
	tokenString, err := GetToken(c)
	if err != nil {
		return nil, err
	}
	if tokenString == "" {
		global.Logger.Debug("tokenString not allow nil")
		return nil, TokenEmptyError
	}
	claims := global.UserStdClaims{}
	_, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(global.AppSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims.User, err
}
