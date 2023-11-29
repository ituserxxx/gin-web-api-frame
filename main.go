package main

import (
	"flag"
	"gin-web-api-ws-mqtt-frame/db"
	"gin-web-api-ws-mqtt-frame/router"
	"gin-web-api-ws-mqtt-frame/service"
	"gin-web-api-ws-mqtt-frame/tools/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var Loc, _ = time.LoadLocation("Asia/Shanghai")

func main() {
	//设置时区
	time.Local = Loc
	var env string
	// go run main.go --env=test //读取.env.test,默认读取.env
	flag.StringVar(&env, "env", "", "--env=test")
	flag.Parse()
	config.InitConfig(env)
	db.Init(config.Get("MYSQL_DSN"))
	r := gin.Default()
	r.Use(cors.Default())
	go service.Manager.WSStart()
	go service.MqttServerRun()
	router.InitRouter(r)
	r.Run(":8077")
}
