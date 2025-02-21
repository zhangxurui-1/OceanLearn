package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"oceanlearn/model"
)

var DB *gorm.DB

func InitDB() {

	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "123456"
	charset := "utf8"
	// 要加parseTime=true！！！
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true", username, password, host, port, database, charset)

	var err error
	DB, err = gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err = DB.AutoMigrate(&model.User{}); err != nil {
		log.Printf("failed to auto migrate database: %v", err)
		return
	}
}

func GetDB() *gorm.DB {
	return DB
}
