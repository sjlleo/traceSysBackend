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
	Username string `gorm:"type:varchar(50);uniqueIndex" json:"username"`
	Role     uint   `gorm:"type:int" json:"role"`
	Email    string `gorm:"type:varchar(80)" json:"email"`
	Phone    string `gorm:"type:varchar(20)" json:"phone"`
	Password string `gorm:"type:varchar(50)" json:"password"`

	Lastseen time.Time
}

func (t *Users) TableName() string { return "users" }

func ValidUser(username string, password string) (ID uint, roleID uint, err error) {
	db := database.GetDB()
	u := Users{}

	err = db.Where(&Users{Username: username, Password: util.MD5(password)}).First(&u).Error
	return u.ID, u.Role, err
}

func CreateUser(u Users) (err error) {
	db := database.GetDB()
	u.Password = util.MD5(u.Password)

	err = db.Create(&u).Error

	if err != nil && strings.Contains(err.Error(), "Error 1062") {
		return errors.New("用户名已存在")
	}

	return nil
}

type UsersRes struct {
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"type:varchar(50);uniqueIndex" json:"username"`
	Role     int    `gorm:"type:int" json:"role"`
	Email    string `gorm:"type:varchar(80)" json:"email"`
	Phone    string `gorm:"type:varchar(20)" json:"phone"`
}

func ListUsers(p *PaginationQ) error {
	var users []UsersRes
	db := database.GetDB()
	tx := db.Model(&Users{})
	// 查找条件
	if p.Parm != "" {
		tx = tx.Where("username like ?", "%"+p.Parm+"%")
	}
	total, err := crudAll(p, tx, &users)
	if err != nil {
		return err
	} else {
		p.Total = uint(total)
		p.Data = users
		return nil
	}
}

func UpdateUser(user Users) error {
	if user.Password != "" {
		user.Password = util.MD5(user.Password)
	}
	db := database.GetDB()
	tx := db.Model(&user)
	// 查找条件
	err := tx.Updates(&user).Error
	return err
}

func DelUser(user Users) error {
	db := database.GetDB()
	err := db.Delete(&user).Error
	return err
}

type Password struct {
	UserID    uint
	BeforePwd string `json:"beforePassword"`
	AfterPwd  string `json:"newPassword"`
}

func UpdatePassword(p *Password) error {
	// 查询密码是否匹配
	db := database.GetDB()
	u := Users{}
	u.ID = p.UserID
	err := db.Model(&u).Take(&u).Error
	if err != nil {
		return errors.New("用户未找到")
	}

	if util.MD5(p.BeforePwd) != u.Password {
		return errors.New("原密码输入不正确")
	}
	err = db.Model(&u).Update("password", util.MD5(p.AfterPwd)).Error
	return err
}
