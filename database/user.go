package database

import (
	"GoProject/global"
	"GoProject/util"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name      string `json:"name"`
	Password  string `json:"-"`
	Phone     string `json:"phone" gorm:"unique"`
	Email     string `json:"email" gorm:"unique"`
	Role      UserRoles
	Status    UserStatus
	GithubTag int `json:"github_tag,omitempty"`
}

func (u *User) CheckPwd(pwd string) bool {
	pwdSha256 := util.GetSha256(pwd)
	return u.Password == pwdSha256
}

func (u *User) GenJwt() (string, error) {
	//return util.GenerateJwt(u)
	expiredTime := time.Now().Add(global.DefaultValidityPeriod)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expiredTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        string(u.ID),
		Issuer:    global.AppIssuer,
	}
	uClaims := UserStdClaims{StandardClaims: stdClaims, User: u}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uClaims)
	tokenString, err := token.SignedString(global.AppJwtSecret)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err":  err,
			"user": u,
		}).Error("generate jwt error")
		return "", err
	}
	return tokenString, nil
}

type UserStdClaims struct {
	StandardClaims jwt.StandardClaims
	User           *User
}

func (c UserStdClaims) Valid() (err error) {
	if c.StandardClaims.ExpiresAt < time.Now().Unix() {
		return errors.New("token is expired")
	}
	if c.StandardClaims.Issuer != global.AppIssuer {
		return errors.New("token's issuer is wrong")
	}
	if c.User.ID < 1 {
		return errors.New("invalid user in jwt")
	}
	return nil
}

type UserRoles int
type UserStatus int

const (
	Normal UserRoles = iota
	Admin
	SuperAdmin

	Pending UserStatus = iota
	Active
	Closed
)

var (
	AllRoles  = [3]UserRoles{SuperAdmin, Admin, Normal}
	AllStatus = [3]UserStatus{Pending, Active, Closed}
)

func (ur UserRoles) String() (string, error) {
	switch ur {
	case 0:
		return "superAdmin", nil
	case 1:
		return "admin", nil
	case 2:
		return "normal", nil
	default:
		global.Logger.WithFields(logrus.Fields{
			"userRole": ur,
		}).Error("parse user role error")
		return "", errors.New("parse user role error")
	}
}
