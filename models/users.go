package models

import (
	"errors"
	"strings"
	"time"

	"github.com/sjlleo/traceSysBackend/database"
	"github.com/sjlleo/traceSysBackend/util"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);uniqueIndex"`
	Role     int    `gorm:"type:int"`
	Email    string `gorm:"type:varchar(80)"`
	Phone    string `gorm:"type:varchar(20)"`
	Password string `gorm:"type:varchar(50)"`

	Lastseen time.Time
}

func (t *Users) TableName() string { return "users" }

func ValidUser(username string, password string) (roleID int, err error) {
	db := database.GetDB()

	u := Users{}

	err = db.Where(&Users{Username: username, Password: util.MD5(password)}).First(&u).Error
	return u.Role, err
}

func CreateUser(username string, password string, roleID int) (err error) {
	db := database.GetDB()
	user := Users{
		Username: username,
		Password: util.MD5(password),
		Role:     roleID,
	}

	err = db.Create(&user).Error

	if err != nil && strings.Contains(err.Error(), "Error 1062") {
		return errors.New("用户名已存在")
	}

	return nil
}
