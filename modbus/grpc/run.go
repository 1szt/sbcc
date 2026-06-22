// ============================================================
//
//	gRPC 服务模块 — SBCC 控制中心
//	功能：启动 gRPC 服务器，监听指定端口，处理 RPC 请求
//	依赖：google.golang.org/grpc
//	启动方式：由 modbus/main 模块统一调用
//
// ============================================================
package grpc

import (
	"log" // 日志输出
	"net" // 网络监听（TCP）
	"time"

	"modbus/env" // 环境配置（读取 GRPC_PORT）

	"google.golang.org/grpc" // gRPC 框架
)

// run 启动 gRPC 服务器
// 步骤：
//  1. 初始化 gRPC 配置（GRPC_PORT）
//  2. 从配置读取端口号
//  3. 创建 TCP 监听器（绑定端口）
//  4. [待办] 注册具体的 gRPC 服务
//  5. 协程启动服务器，主协程通过 select 检测启动结果
func Run() {
	// 第一步：初始化 gRPC 配置
	// env.Init() 仅在首次调用时写入默认值，后续读取已有配置
	env.Init([][]string{
		{"GRPC_PORT", "50051", "gRPC 监听端口"},
	})

	// 第二步：从配置读取端口号
	grpcPort := env.Get("GRPC_PORT")

	// 第三步：在指定端口上创建 TCP 监听器
	// net.Listen("tcp", ":50051") 会监听所有网卡的 50051 端口
	// 如果端口被占用或权限不足，会返回错误
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		// log.Fatalf 会打印错误信息并调用 os.Exit(1) 退出程序
		// 常见错误：端口已被占用、无权限绑定低端口（<1024）
		log.Fatalf("❌ [gRPC] TCP 监听失败（端口可能被占用）: %v", err)
	}

	// 第四步：创建 gRPC 服务器实例
	// grpc.NewServer() 可传入可选参数，如拦截器（Interceptor）、
	// 最大消息大小、TLS 凭证等
	grpcServer := grpc.NewServer()

	// Todo: 注册你的 gRPC 服务
	// 示例：pb.RegisterYourServiceServer(grpcServer, &YourServiceImpl{})
	// Register your gRPC services here

	// 创建错误探测通道，缓冲区为 1 防止协程泄漏
	errChan := make(chan error, 1)

	// 第五步：协程启动服务器
	// grpcServer.Serve(lis) 会阻塞当前协程，持续接受和处理 RPC 请求
	// 服务器会一直运行，直到收到终止信号或发生致命错误
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	// 第六步：等待 100ms 确认服务是否成功启动
	select {
	case err := <-errChan:
		// 如果通道里收到错误，说明端口启动失败（如 Address already in use）
		log.Fatalf("❌ [gRPC] 服务启动失败: %v", err)
		return
	case <-time.After(100 * time.Millisecond):
		// 100ms 过去了没报错，说明端口占领成功
		log.Printf("✅ [gRPC] 服务器已启动，监听端口 :%s", grpcPort)
	}
}
