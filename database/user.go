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
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Password string `json:"-"`
	Phone    string `json:"phone" gorm:"unique"`
	Role     UserRoles
	Status   UserStatus
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

type UserCollect struct {
	gorm.Model
	UserID    uint    `json:"user_id"`
	User      User    `json:"user"`
	ArticleID uint    `json:"article_id"`
	Article   Article `json:"article"`
}

type UserHistory struct {
	gorm.Model
	UserID    uint    `json:"user_id"`
	User      User    `json:"user"`
	ArticleID uint    `json:"article_id"`
	Article   Article `json:"article"`
}

type OauthUser struct {
	PlatformID        uint      `json:"platform_id" gorm:"autoIncrement:true;uniqueIndex"`
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Platform          string    `json:"platform" gorm:"primaryKey;autoIncrement:false"`
	Login             string    `json:"login"`
	Name              string    `json:"name"`
	AvatarUrl         string    `json:"avatar_url"`
	Url               string    `json:"url"`
	HtmlUrl           string    `json:"html_url"`
	FollowersUrl      string    `json:"followers_url"`
	FollowingUrl      string    `json:"following_url"`
	GistsUrl          string    `json:"gists_url"`
	StarredUrl        string    `json:"starred_url"`
	SubscriptionsUrl  string    `json:"subscriptions_url"`
	OrganizationsUrl  string    `json:"organizations_url"`
	ReposUrl          string    `json:"repos_url"`
	EventsUrl         string    `json:"events_url"`
	ReceivedEventsUrl string    `json:"received_events_url"`
	PublicRepos       uint      `json:"public_repos"`
	PublicGists       uint      `json:"public_gists"`
	Followers         uint      `json:"followers"`
	Watched           uint      `json:"watched"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Email             string    `json:"email"`
}

type UserBind struct {
	gorm.Model
	PlatformId string    `json:"platform_id" `
	OauthUser  OauthUser `json:"oauth_user" gorm:"foreignKey:platform_id;references:platform_id"`
	User       User
	UserId     int
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
