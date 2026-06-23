// ============================================================
//
//	WebSocket 服务模块 — SBCC 控制中心
//	功能：启动 WebSocket 服务器，管理实时双向通信连接
//	依赖：github.com/gorilla/websocket
//	启动方式：由 modbus/main 模块统一调用
//
// ============================================================
package websocket

import (
	"fmt"      // 格式化输出
	"log"      // 日志输出
	"net/http" // HTTP 服务
	"time"

	"modbus/env" // 环境配置（读取 WEBSOCKET_PORT）
)

// Hub 全局实例（包级变量）
// 所有 WebSocket 连接管理通过此实例进行
var GlobalHub *Hub

// Run 启动 WebSocket 服务
// 步骤：
//  1. 初始化配置（WEBSOCKET_PORT）
//  2. 创建全局 Hub 连接管理器
//  3. 注册 HTTP 路由（WebSocket 握手端点）
//  4. 协程启动 HTTP 服务器
//  5. 确认服务启动成功
func Run() {
	// 第一步：初始化 WebSocket 配置
	env.Init([][]string{
		{"WEBSOCKET_PORT", "9082", "WebSocket 服务端口"},
		{"WEBSOCKET_PATH", "/ws", "WebSocket 连接路径"},
	})

	// 第二步：从配置读取端口和路径
	wsPort := env.Get("WEBSOCKET_PORT")
	wsPath := env.Get("WEBSOCKET_PATH")

	// 第三步：创建全局 Hub 连接管理器
	GlobalHub = NewHub()

	// 第四步：创建 HTTP 路由
	// WebSocket 握手端点
	mux := http.NewServeMux()
	mux.HandleFunc(wsPath, GlobalHub.HandleWebSocket)

	// 健康检查端点
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"status":"ok","connections":%d}`, GlobalHub.Count())))
	})

	// 创建错误探测通道
	errChan := make(chan error, 1)

	// 第五步：协程启动 HTTP 服务器（用于 WebSocket 升级）
	go func() {
		if err := http.ListenAndServe(":"+wsPort, mux); err != nil {
			errChan <- err
		}
	}()

	// 第六步：等待 100ms 确认服务是否成功启动
	select {
	case err := <-errChan:
		log.Fatalf("❌ [WebSocket] 服务启动失败: %v", err)
		return
	case <-time.After(100 * time.Millisecond):
		log.Printf("✅ [WebSocket] 服务器已启动，监听端口 :%s", wsPort)
		log.Printf("🔌 [WebSocket] 连接端点: ws://localhost:%s%s", wsPort, wsPath)
	}
}
