package models

import (
	"github.com/jinzhu/gorm"
	"go-blog/db"
)

func init() {
	db.Db.AutoMigrate(&Roles{})
}

type Roles struct {
	gorm.Model
	Name      	  string `gorm:"type:varchar(20);unique_index;comment:\"角色名\""`
	Info          string `gorm:"type:text;comment:\"角色介绍\""`
}

func (r *Roles) TableName() string {
	return "roles"
}

