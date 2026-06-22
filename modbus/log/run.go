// log 带模块前缀的日志工具
//
// 用法：
//
//	// main.go 中初始化（从 .env 读取配置）
//	log.Init()
//
//	// 各模块创建自己的日志器
//	var log = log.New("CHI")
//	log.Printf("%s端口占领成功", port)
//
// .env 配置：
//
//	LOG_FILE=data/app.log    # 输出到文件（不设置则输出到控制台）
//
// 输出：2026/06/22 09:42:18 [CHI] 9081端口占领成功
package log

import (
	"io"
	"log"
	"os"

	"modbus/env"
)

// 全局默认配置，New() 创建日志器时从这个配置读取
var (
	defaultOutput io.Writer = os.Stdout
	defaultFlags  int       = log.Ldate | log.Ltime
)

// Init 从 .env 读取日志配置并生效。
// 必须在各模块 New() 之前调用。
func Init() {
	// 注册日志模块自己的配置项到 .env
	env.Init([][]string{
		{"LOG_FILE", "", "日志文件路径（空=输出到控制台）"},
	})

	// 如果设置了 LOG_FILE，输出到文件
	if path := env.GetConfig("LOG_FILE"); path != "" {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			defaultOutput = f
		}
	}
}

// New 创建带模块前缀的日志器
func New(module string) *log.Logger {
	return log.New(defaultOutput, "["+module+"] ", defaultFlags)
}

// SetOutput 设置日志输出位置（默认 os.Stdout）
func SetOutput(w io.Writer) {
	defaultOutput = w
}

// SetFlags 设置日志格式（默认 log.Ldate | log.Ltime）
// 可选值：log.Ldate, log.Ltime, log.Lmicroseconds, log.Lshortfile 等
func SetFlags(flag int) {
	defaultFlags = flag
}
