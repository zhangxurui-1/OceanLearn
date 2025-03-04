package common

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/url"
	"oceanlearn/model"
)

var DB *gorm.DB

func InitDB() {

	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	location := viper.GetString("datasource.location")
	// 要加parseTime=true！！！
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s", username, password, host, port, database, charset, url.QueryEscape(location))

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
