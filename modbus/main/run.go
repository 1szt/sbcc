package main

import (
	"fmt"
	"modbus/home"
	"modbus/sub"
	"modbus/web"
	"modbus/db"
)

func main() {

	// motd
	fmt.Println("========================================")
	fmt.Println("             SBCC 控制中心启动！          ")
	fmt.Println("========================================")

	// Web 引擎启动
	web.Run()

	// 数据库模块启动
	db.Run()

	// 功能模块启动
	home.Run()
	sub.Run()

	// gin.Run()

	// 阻塞主进程
	select {}
}
