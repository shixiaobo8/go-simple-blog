package models

import (
	"github.com/jinzhu/gorm"
	"go-blog/db"
)

func init() {
	db.Db.AutoMigrate(&AdminRoles{})
	db.Db.Model(&AdminRoles{}).AddUniqueIndex("idx_admin_role", "admin_id", "role_id")
}

type AdminRoles struct {
	gorm.Model
	AdminId	    int
	RoleId		int
}

func (a *AdminRoles) TableName() string {
	return "admin_roles"
}