package admin

import(
	"net/http"
	"os"
)

const Admin="/admin"

func AdminHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置内容类型为HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		// 读取HTML模板文件
		templatePath := "static/nexus-admin-template.html"
		html, err := os.ReadFile(templatePath)
		if err != nil {
			// 如果模板文件不存在，返回错误信息
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error loading admin template: " + err.Error()))
			return
		}
		
		// 返回HTML内容
		w.Write(html)
    })
}
