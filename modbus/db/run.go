package db

import (
	"database/sql"
	"fmt"
	"modbus/env"
	"time"

	_ "github.com/go-sql-driver/mysql" // 换成你实际使用的驱动
)

var DB *sql.DB

func Run() {
	// 一次性初始化数据库配置
	env.Init([][]string{
		// { "KEY", "DEFAULT", "COMMENT..." }
		{"DB_TYPE", "sqlite", "数据库驱动类型", "支持: sqlite (文件模式), pgsql/mysql (服务器模式)"},
		{"DB_NAME", "trade_data.db", "数据库名", "若使用 sqlite，此处为文件名；若使用 pgsql/mysql，此处为数据库实例名"},

		// 以下为网络数据库专用配置
		{"DB_HOST", "127.0.0.1", "数据库地址", "仅针对 pgsql/mysql，sqlite 模式下会忽略此项"},
		{"DB_PORT", "5432", "数据库端口", "pgsql 默认 5432, mysql 默认 3306, sqlite 忽略"},
		{"DB_USER", "postgres", "数据库用户名"},
		{"DB_PASSWORD", "your_password", "数据库密码"},
	})

	dsn := "123456@tcp(127.0.0.1:3306)/my_db?parseTime=True"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("⚠️ [DB] 数据库模块 配置错误！")
	}

	// 间隔2s
	time.Sleep(2 * time.Second)

	// 开始无限循环重连
	for {

		err = DB.Ping()
		if err == nil {
			fmt.Printf("✅ [DB] 数据库 连接成功！\n")
			break // 连上了，跳出循环
		}

		fmt.Printf("🔄 [DB] 数据库连接失败: %v。 9秒后尝试重连...\n", err)
		time.Sleep(9 * time.Second)
		// fmt.Printf("❌ [DB] 数据库 重连中...\n")
	}

	// 设置连接池（连上后再设置）
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(20)
}
