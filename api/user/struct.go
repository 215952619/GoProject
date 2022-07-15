package user

import (
	"GoProject/database"
	"GoProject/util"
	"errors"
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

func (lr *logonRequest) CheckPwd() (user database.User, err error) {
	if err = database.DBM.First(&user, "name=? or phone=? or email=?", lr.Identifier, lr.Identifier, lr.Identifier); err != nil {
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
