package user

import (
	"GoProject/database"
	"GoProject/util"
	"errors"
	"gorm.io/gorm"
)

type logonRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"pwd"`
	Dots       int    `json:"dots"`
	ValidCode  string `json:"valid_code"`
}

func (lr *logonRequest) CheckCode() (bool, error) {
	return util.CheckCaptcha(lr.Dots, lr.ValidCode)
}

func (lr *logonRequest) CheckPwd() (user *database.User, err error) {
	if err = database.DBM.First(&user, "name=? or phone=?", lr.Identifier, lr.Identifier); err != nil {
		return user, err
	} else {
		if user.CheckPwd(lr.Password) {
			return user, nil
		} else {
			return user, errors.New("密码错误")
		}
	}
}

type NewUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     database.UserRoles
}

type CheckRequest struct {
	Platform string `json:"platform"`
	ID       uint   `json:"id"`
}

type BindRequest struct {
	PlatformUser database.OauthUser `json:"oauth"`
	AutoRegister bool               `json:"auto_register"`
	Identifier   string             `json:"identifier"`
	Password     string             `json:"pwd"`
}

func (br *BindRequest) CheckPwd() (user *database.User, err error) {
	if err = database.DBM.First(&user, "name=? or phone=?", br.Identifier, br.Identifier); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if br.AutoRegister {
				var oauthUser database.OauthUser
				if err = database.DBM.Db.Where(br.PlatformUser).First(&oauthUser).Error; err != nil {
					return nil, err
				}

				if err = database.DBM.Upsert(&user, database.User{Name: oauthUser.Name, Avatar: oauthUser.AvatarUrl}, br.PlatformUser); err != nil {
					return nil, err
				}
			} else {
				return nil, util.DefaultError(util.NotFound)
			}
		}
		return nil, err
	} else {
		if user.CheckPwd(br.Password) {
			return user, nil
		} else {
			return nil, util.DefaultError(util.PasswordWrong)
		}
	}
}
