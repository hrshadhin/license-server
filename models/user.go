package models

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	u "github.com/hrshadhin/license-server/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to represent users
type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(50);not null" json:"name"`
	Email    string `gorm:"type:varchar(100);not null;unique_index" json:"email"`
	Password string `gorm:"not null"`
	Token    string `gorm:"-" json:"token"`
}

//Validate incoming user details...
func (user *User) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(user.Email, "@") {
		return u.Message(false, "Email address is required."), false
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Password is required. Min: 6 characters."), false
	}

	//Email must be unique
	temp := &User{}

	//check for errors and duplicate emails
	err := GetDB().Table("users").Where("email = ?", user.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry."), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (user *User) Create() map[string]interface{} {

	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	GetDB().Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create user, connection error.")
	}

	user.Password = "" //delete password

	response := u.Message(true, "User has been created.")
	response["user"] = user
	return response
}

func Login(email, password string) map[string]interface{} {

	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found.")
		}
		return u.Message(false, "Connection error. Please retry.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	user.Password = ""

	//Create JWT token
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(u.MustGetEnv("token_password")))
	user.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

func FetchAllUsers() map[string]interface{} {

	var users []User
	err := GetDB().Table("users").Select("id, name, email, created_at, updated_at, deleted_at").Order("id").Find(&users).Error
	if err != nil {
		return u.Message(false, "Connection error. Please retry.")
	}

	resp := u.Message(true, "User List.")
	resp["users"] = users
	return resp
}

func GetUser(userId string) map[string]interface{} {
	user := &User{}
	err := GetDB().Table("users").Where("id = ?", userId).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "User not found.")
		}
		return u.Message(false, "Connection error. Please retry.")
	}

	resp := u.Message(true, "User details.")
	user.Password = ""
	resp["user"] = user
	return resp
}

func UpdateUser(userId, name, email string) map[string]interface{} {
	user := &User{}
	err := GetDB().Table("users").Where("id = ?", userId).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "User not found.")
		}
		return u.Message(false, "Connection error. Please retry.")
	}

	temp := &User{}
	count := 0
	//check for errors and duplicate emails
	GetDB().Table("users").Where("id <> ?", userId).Where("email = ?", email).Find(&temp).Count(&count)

	if count > 0 {
		return u.Message(false, "Email address already in use by another user.")
	}
	user.Email = email
	user.Name = name
	GetDB().Model(&user).Updates(user)

	resp := u.Message(true, "User updated.")
	user.Password = ""
	resp["user"] = user
	return resp
}

func ChangePassword(userId, password string) map[string]interface{} {
	user := &User{}
	err := GetDB().Table("users").Where("id = ?", userId).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "User not found.")
		}
		return u.Message(false, "Connection error. Please retry.")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	GetDB().Model(&user).Updates(user)

	return u.Message(true, "User password changed.")
}

func DeleteUser(userId string) map[string]interface{} {
	user := &User{}
	err := GetDB().Table("users").Where("id = ?", userId).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "User not found.")
		}
		return u.Message(false, "Connection error. Please retry.")
	}

	GetDB().Delete(&user)

	return u.Message(true, "User deleted.")
}
