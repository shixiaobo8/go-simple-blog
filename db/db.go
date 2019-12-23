package db

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	Db *gorm.DB
)

func init() {
	var err error

	dbHost := beego.AppConfig.String("dbHost")
	dbPort, _ := beego.AppConfig.Int("dbPort")
	dbName := beego.AppConfig.String("dbName")
	dbUser := beego.AppConfig.String("dbUser")
	dbPassword := beego.AppConfig.String("dbPassword")

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	Db, err = gorm.Open("mysql", connString)

	if err != nil {
		panic("failed to connect database")
	}

	Db.LogMode(true)
}
