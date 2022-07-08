package user

import "GoProject/database"

type NewUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Role     database.UserRoles
}
