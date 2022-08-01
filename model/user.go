package model

import (
	"errors"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"gorm.io/gorm"
)

type UserInfo struct {
	gorm.Model
	Folder   []Folder `gorm:"foreignKey:UserID"`
	UserName string   `gorm:"unique"`
	PassWord string   //hash值
}

//
func (d *DB) AddUser(username, password string) (id uint, err error) {
	user := new(UserInfo)
	user.UserName = username
	if len(password) >= 6 {
		user.PassWord = gmd5.MustEncrypt(password)
	} else {
		return 0, errors.New("password len too short")
	}
	err = db.Create(user).Error
	if err != nil {
		return 0, err
	}
	d.InitFolder(user.ID)
	return user.ID, nil
}

func (d *DB) GetUser(username, password string) (userid uint, err error) {
	user := new(UserInfo)
	db.Where("user_name = ?", username).First(&user)
	if user.PassWord == gmd5.MustEncrypt(password) {
		return user.ID, nil
	} else {
		return 0, errors.New("err password")
	}
}

//func (d *DB)DelUser(username, password string) (ok bool, data string) {
//	db.Where("username = ?")
//}
