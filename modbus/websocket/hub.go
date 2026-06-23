// ============================================================
//
//	WebSocket Hub 连接管理器 — SBCC 控制中心
//	功能：管理所有 WebSocket 客户端连接，支持广播、注册、注销
//	依赖：gorilla/websocket
//	启动方式：由 websocket.Run() 自动创建
//
// ============================================================
package websocket

import (
	"log"      // 日志输出
	"net/http" // HTTP 处理
	"sync"     // 并发安全（读写锁）

	"github.com/gorilla/websocket" // WebSocket 库
)

// Hub 连接管理器，负责管理所有 WebSocket 客户端连接
type Hub struct {
	// 读写锁，保证并发安全
	mu sync.RWMutex

	// 所有活跃连接的集合（map[客户端ID]连接对象）
	clients map[string]*Client

	// HTTP 升级器（将 HTTP 连接升级为 WebSocket）
	upgrader websocket.Upgrader
}

// NewHub 创建并初始化一个新的 Hub 实例
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		upgrader: websocket.Upgrader{
			// 允许所有来源的连接（开发环境）
			// 生产环境应限制 CheckOrigin
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			// 读缓冲区大小
			ReadBufferSize: 1024,
			// 写缓冲区大小
			WriteBufferSize: 1024,
		},
	}
}

// Register 注册一个新的客户端连接
// client: 已握手成功的 WebSocket 客户端
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client.ID] = client
	log.Printf("🔗 [WebSocket] 客户端已连接: %s（当前在线: %d）", client.ID, len(h.clients))
}

// Unregister 注销一个客户端连接
// clientID: 要移除的客户端 ID
func (h *Hub) Unregister(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client, ok := h.clients[clientID]; ok {
		client.Conn.Close()
		delete(h.clients, clientID)
		log.Printf("🔌 [WebSocket] 客户端已断开: %s（当前在线: %d）", clientID, len(h.clients))
	}
}

// Broadcast 向所有连接的客户端广播消息
// message: 要发送的消息内容（字节数组）
func (h *Hub) Broadcast(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for id, client := range h.clients {
		select {
		case client.Send <- message:
		default:
			// 客户端发送缓冲区已满，跳过
			log.Printf("⚠️ [WebSocket] 客户端 %s 发送缓冲区已满，跳过消息", id)
		}
	}
}

// Count 返回当前在线客户端数量
func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// HandleWebSocket 处理 WebSocket 握手请求
// 将 HTTP 连接升级为 WebSocket，并启动读写协程
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 协议
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ [WebSocket] 升级失败: %v", err)
		return
	}

	// 创建客户端实例
	client := NewClient(conn, h)

	// 注册客户端
	h.Register(client)

	// 启动客户端的读写协程
	go client.WritePump()
	go client.ReadPump()
}
