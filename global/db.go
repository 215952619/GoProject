package global

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

var (
	DB  *gorm.DB
	DBM *DbManager
)

type DbManager struct {
	db *gorm.DB
}

func (dm *DbManager) Migrate() {
	if err := dm.db.AutoMigrate(
		&User{},
		&Article{},
	); err != nil {
		Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("migrate database panic")
	}
	Logger.Info("migrate database success")
}

func InitDb() {
	db, err := gorm.Open(sqlite.Open(DbPath))
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("init db panic")
	}

	DBM = &DbManager{db}
	DBM.Migrate()
	DB = db
}

type User struct {
	gorm.Model
	Name      string `json:"name"`
	Password  string `json:"-"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Role      UserRoles
	Status    UserStatus
	GithubTag int `json:"github_tag,omitempty"`
}

type UserStdClaims struct {
	StandardClaims jwt.StandardClaims
	User           *User
}

func (c UserStdClaims) Valid() (err error) {
	if c.StandardClaims.ExpiresAt < time.Now().Unix() {
		return errors.New("token is expired")
	}
	if c.StandardClaims.Issuer != AppIssuer {
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
	SuperAdmin UserRoles = iota
	Admin
	Normal

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
		Logger.WithFields(logrus.Fields{
			"userRole": ur,
		}).Error("parse user role error")
		return "", errors.New("parse user role error")
	}
}

type Article struct {
	gorm.Model
}
