package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // 换成你实际使用的驱动
)

var DB *sql.DB

func Run() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/my_db?parseTime=True"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("⚠️ [DB] 数据库模块 配置错误！")
	}

	// 开始无限循环重连
	for {
		err = DB.Ping()
		if err == nil {
			fmt.Println("✅ [DB] 数据库模块 连接成功！")
			break // 连上了，跳出循环
		}

		fmt.Printf("🔄 [DB] 数据库模块 连接失败: %v。 5秒后尝试重连...\n", err)
		time.Sleep(5 * time.Second)
		// fmt.Println("❌ [DB] 数据库模块 重连中...")
	}

	// 设置连接池（连上后再设置）
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(20)
}
