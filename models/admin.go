package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"go-blog/db"
	"time"
)

func init() {
	db.Db.AutoMigrate(&Admin{})
}

type Admin struct {
	gorm.Model
	Username      string `gorm:"type:varchar(20);unique_index"`
	Password      string `gorm:"size:32"`
	Nickname      string `gorm:"type:varchar(8);unique_index"`
	Avatar        string `gorm:"type:varchar(255);not null;default:\"\""`
	LastLoginTime *time.Time
	RoleId		  int `gorm:"not null;default:0"`
}

func (a *Admin) TableName() string {
	return "admins"
}

func PasswordEnc(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func (a *Admin) Add() error {
	var err error
	if a.ID != 0 {
		err = fmt.Errorf("ID must be zero")
		return err
	}
	db.Db.Create(a)
	return err
}

func (a *Admin) UpdateLastLoginTime(time *time.Time) {
	a.LastLoginTime = time
	db.Db.Model(a).UpdateColumn("last_login_time", *time)
}

func (a *Admin) UpdateAdminInfo(nickname string, avatar string) error {
	if a.ID == 0 {
		return errors.New("user not exists")
	}
	db.Db.Model(a).UpdateColumns(Admin{
		Model: gorm.Model{
			UpdatedAt: time.Now(),
		},
		Nickname:      nickname,
		Avatar:        avatar,
	})
	return nil
}

func (a *Admin) List(fields []string) []Admin {
	var admins []Admin
	if fields != nil {
		db.Db.Model(Admin{}).Select(fields).Find(&admins)
	} else {
		db.Db.Model(Admin{}).Find(&admins)
	}

	return admins
}

func (a *Admin) Info(fields []string) error {
	if a.ID <= 0 {
		return errors.New("invalid id")
	}
	db.Db.Model(a).Select(fields).First(a)
	return nil
}