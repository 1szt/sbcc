package main

import (
	"fmt"
	"modbus/grom"
	"modbus/home"
	"modbus/sql"
	"modbus/sqlx"
	"modbus/sub"
	"modbus/web"
)

func main() {

	// motd
	fmt.Println("========================================")
	fmt.Println("             SBCC 控制中心启动！          ")
	fmt.Println("========================================")

	// Web 引擎启动
	web.Run()

	// 数据库模块启动
	sql.Run()
	sqlx.Run()
	grom.Run()

	// 功能模块启动
	home.Run()
	sub.Run()

	// gin.Run()

	// 阻塞主进程
	select {}
}
