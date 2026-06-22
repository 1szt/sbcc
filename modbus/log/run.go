// log 带模块前缀的日志工具
//
// 用法：
//
//	// main.go 中配置（可选，不调用则默认输出到控制台）
//	log.SetOutput(file)            // 输出到文件
//	log.SetFlags(log.Ldate|log.Ltime)  // 设置日期格式
//
//	// 各模块创建自己的日志器
//	var log = log.New("CHI")
//	log.Printf("%s端口占领成功", port)
//
// 输出：2026/06/22 09:42:18 [CHI] 9081端口占领成功
package log

import (
	"io"
	"log"
	"os"
)

// 全局默认配置，New() 创建日志器时从这个配置读取
var (
	defaultOutput io.Writer = os.Stdout
	defaultFlags  int       = log.Ldate | log.Ltime
)

// New 创建带模块前缀的日志器
func New(module string) *log.Logger {
	return log.New(defaultOutput, "["+module+"] ", defaultFlags)
}

// SetOutput 设置日志输出位置（默认 os.Stdout）
// 在 main.go 中调用，必须在各模块 New() 之前设置
func SetOutput(w io.Writer) {
	defaultOutput = w
}

// SetFlags 设置日志格式（默认 log.Ldate | log.Ltime）
// 可选值：log.Ldate, log.Ltime, log.Lmicroseconds, log.Lshortfile 等
func SetFlags(flag int) {
	defaultFlags = flag
}
