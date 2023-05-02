package main

import (
	"WeToken/api"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	dsn string
)

func main() {
	viper.AddConfigPath("./config/")
	viper.SetConfigName("key")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("配置文件读取出错")
		return
	}

	user := viper.GetString("mysql.user")
	pass := viper.GetString("mysql.pass")
	host := viper.GetString("mysql.host")
	dbname := viper.GetString("mysql.dbname")
	dsn = user + ":" + pass + "@tcp(" + host + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	api.Kdb, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	sqlDB, _ := api.Kdb.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(30)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(300)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(60 * time.Second)

	port := viper.GetString("port") //这是运行端口不是数据库端口
	r := gin.Default()
	r.GET("/", api.MpVerify)           //路径自己设置，我用的独立二级域名所以直接用根目录
	r.POST("/", api.ChatToken)         //同上
	r.POST("/verify", api.TokenVerify) //验证token是否有效

	r.Run(":" + port)
}
