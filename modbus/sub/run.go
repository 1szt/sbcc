package sub

// 依赖模块 ：web
// 订阅api模块
// 订阅api模块负责处理用户访问订阅api路径的请求
// 挂载路径："/api/sub"

import (
	"fmt"
	"modbus/web"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
)

func Run() {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get("/", sub)
	})

	web.Mux.Mount("/api/sub", r)
	fmt.Println("🎉 [Sub] 订阅api模块 加载完成！")

}

func sub(w http.ResponseWriter, r *http.Request) {

	// // 获取?token参数
	// token := r.URL.Query().Get("token")
	// if token == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, "错误：缺少token参数。")
	// 	return
	// }

	ConfigPath := "data/clash.yaml"

	absPath, _ := filepath.Abs(ConfigPath)

	info, err := os.Stat(absPath)
	if err != nil || info.Size() == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "错误：配置文件不存在或为空。路径：%s", absPath)
		return
	}

	// 配置文件显示名称
	filename := "SB控制中心"
	// 更新间隔（单位: 小时）
	update := "1"
	// 流量信息(单位: 字节)

	// 使用int64类型表示流量信息，避免溢出
	total := int64(1024 * 1024 * 1024 * 222)
	upload := int64(1024 * 1024 * 222)
	download := int64(1024 * 1024 * 500)
	// 到期时间(时间戳)
	expire := time.Date(2077, 12, 10, 0, 0, 0, 0, time.Local).Unix()
	weburl := "http://localhost:9081"

	// Header 文档
	//https://www.clashverge.dev/guide/url_schemes.html

	// 配置文件名
	w.Header().Set("content-disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.PathEscape(filename)))

	// 更新间隔
	w.Header().Set("profile-update-interval", update)

	// 流量信息(单位: 字节)、到期信息(时间戳)
	w.Header().Set("subscription-userinfo", fmt.Sprintf("upload=%d; download=%d; total=%d; expire=%d", upload, download, total, expire))

	// 首页URL
	w.Header().Set("profile-web-page-url", weburl)

	// 基础 Header 设置
	w.Header().Set("content-type", "application/x-yaml; charset=utf-8")

	// 发送文件
	http.ServeFile(w, r, absPath)

}
