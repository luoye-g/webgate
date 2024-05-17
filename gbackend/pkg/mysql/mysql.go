package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(host, port, user, pass string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/blog?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return db
}
