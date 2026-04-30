package main

import (
	"fmt"
	"sbcc/internal/home"
	"sbcc/internal/sub"
	"sbcc/internal/web"
)

func main() {

	// motd
	fmt.Println("========================================")
	fmt.Println("             SBCC 控制中心启动！          ")
	fmt.Println("========================================")

	// Web 引擎启动
	web.Run()

	// 功能模块启动
	home.Run()
	sub.Run()

	// gin.Run()

	// 阻塞主进程
	select {}
}
