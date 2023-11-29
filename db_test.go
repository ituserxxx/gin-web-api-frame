package main

import (
	"flag"
	"fmt"
	"gin-web-api-ws-mqtt-frame/db"
	"gin-web-api-ws-mqtt-frame/router"
	"gin-web-api-ws-mqtt-frame/tools/config"
	"testing"
)

func TestData2(t *testing.T) {
	config.InitConfig("")

	l, err := router.GetCameraDevices()
	if err != nil {
		fmt.Printf("init camera device err :%#v", err.Error())
		return
	}
	for _, info := range l {
		fmt.Printf("url %#v", info)

	}
}
func TestData(t *testing.T) {
	var env string
	// go run main.go --env=test //读取.env.test,默认读取.env
	flag.StringVar(&env, "env", "", "--env=test")
	flag.Parse()
	config.InitConfig(env)
	db.Init(config.Get("MYSQL_DSN"))
	var dinfo *db.Dcc
	db.DB.Model(db.Dcc{}).Where("dcc_number=?", "xxx").Find(&dinfo)
	fmt.Printf("%#v", dinfo)
}
