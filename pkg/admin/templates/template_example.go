package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// TemplateData 定义传递给模板的数据结构
type TemplateData struct {
	Title       string
	CurrentPage string
	User        User
	Stats       DashboardStats
	Activities  []Activity
	// 仓库相关数据
	Repositories []Repository
	Pagination   *Pagination
	// 下载相关数据
	Query          string
	Modules        []Module
	PopularModules []Module
}

// User 用户信息
type User struct {
	Name   string
	Avatar string
}

// DashboardStats 仪表盘统计信息
type DashboardStats struct {
	TotalModules       int
	TotalDownloads     int
	NewModulesToday    int
	ActiveRepositories int
}

// Activity 活动信息
type Activity struct {
	Icon        string
	Title       string
	Description string
	Time        string
}

// Repository 仓库信息
type Repository struct {
	ID          string
	Name        string
	Description string
	Type        string // "public" or "private"
	Version     string
	UpdatedAt   string
}

// Module 模块信息
type Module struct {
	Path        string
	Version     string
	Description string
	Downloads   int
	UpdatedAt   string
	License     string
	IsPrivate   bool
}

// Pagination 分页信息
type Pagination struct {
	CurrentPage int
	TotalPages  int
	HasNext     bool
	HasPrev     bool
	NextPage    int
	PrevPage    int
	Pages       []int
}

// TemplateManager 模板管理器
type TemplateManager struct {
	TemplateDir string
	Templates   map[string]*template.Template
}

// ParseTemplates 解析所有模板文件
func (tm *TemplateManager) ParseTemplates() error {
	// 解析基础布局和组件模板
	baseTemplates := []string{
		"layout.html",
		"header.html",
		"nav.html",
		"sidebar.html",
		"footer.html",
		"styles.html",
	}

	// 解析样式模板
	styleTemplates := []string{
		"dashboard-styles.html",
		"repository-styles.html",
		"upload-styles.html",
		"download-styles.html",
	}

	// 解析页面模板
	pageTemplates := []string{
		"dashboard_nested.html",
		"dashboard_refactored.html",
		"repositories_refactored.html",
		"upload_refactored.html",
		"download_refactored.html",
		"simple_page.html",
	}

	for _, pageTemplate := range pageTemplates {
		// 创建模板集合，包含基础模板、样式模板和当前页面模板
		templateFiles := make([]string, 0)
		
		// 添加基础模板
		for _, base := range baseTemplates {
			templateFiles = append(templateFiles, filepath.Join(tm.TemplateDir, base))
		}
		
		// 添加样式模板
		for _, style := range styleTemplates {
			templateFiles = append(templateFiles, filepath.Join(tm.TemplateDir, style))
		}
		
		// 添加页面模板
		templateFiles = append(templateFiles, filepath.Join(tm.TemplateDir, pageTemplate))

		// 解析模板
		tmpl, err := template.ParseFiles(templateFiles...)
		if err != nil {
			return err
		}

		// 存储模板
		tm.Templates[pageTemplate] = tmpl
	}

	return nil
}

// RenderTemplate 渲染指定的模板
func (tm *TemplateManager) RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) error {
	tmpl, exists := tm.Templates[templateName]
	if !exists {
		return errors.New("template not found: " + templateName)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(w, data)
}

// DashboardHandler 仪表盘页面处理器
func DashboardHandler(tm *TemplateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := TemplateData{
			Title:       "仪表盘 - Athens管理系统",
			CurrentPage: "dashboard",
			User: User{
				Name:   "管理员",
				Avatar: "/static/images/avatar.png",
			},
			Stats: DashboardStats{
				TotalModules:      156,
				TotalDownloads:    2847,
				NewModulesToday:   12,
				ActiveRepositories: 8,
			},
			Activities: []Activity{
				{
					Icon:        "📦",
					Title:       "新模块上传",
					Description: "github.com/example/module v1.2.0",
					Time:        "2分钟前",
				},
				{
					Icon:        "⬇️",
					Title:       "模块下载",
					Description: "github.com/gin-gonic/gin v1.9.1",
					Time:        "5分钟前",
				},
				{
					Icon:        "🔄",
					Title:       "代理同步",
					Description: "从 proxy.golang.org 同步了 23 个模块",
					Time:        "1小时前",
				},
			},
		}

		err := tm.RenderTemplate(w, "dashboard_refactored.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// RepositoriesHandler 仓库管理页面处理器
func RepositoriesHandler(tm *TemplateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := TemplateData{
			Title:       "仓库管理 - Athens管理系统",
			CurrentPage: "repositories",
			User: User{
				Name:   "管理员",
				Avatar: "/static/images/avatar.png",
			},
			Repositories: []Repository{
				{
					ID:          "1",
					Name:        "github.com/gin-gonic/gin",
					Description: "Gin is a HTTP web framework written in Go",
					Type:        "public",
					Version:     "v1.9.1",
					UpdatedAt:   "2小时前",
				},
				{
					ID:          "2",
					Name:        "github.com/gorilla/mux",
					Description: "A powerful HTTP router and URL matcher",
					Type:        "public",
					Version:     "v1.8.0",
					UpdatedAt:   "1天前",
				},
			},
			Pagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  3,
				HasNext:     true,
				HasPrev:     false,
				NextPage:    2,
				Pages:       []int{1, 2, 3},
			},
		}

		err := tm.RenderTemplate(w, "repositories_refactored.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// UploadHandler 上传页面处理器
func UploadHandler(tm *TemplateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := TemplateData{
			Title:       "上传模块 - Athens管理系统",
			CurrentPage: "upload",
			User: User{
				Name:   "管理员",
				Avatar: "/static/images/avatar.png",
			},
		}

		err := tm.RenderTemplate(w, "upload_refactored.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// DownloadHandler 下载页面处理器
func DownloadHandler(tm *TemplateManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		
		data := TemplateData{
			Title:       "下载模块 - Athens管理系统",
			CurrentPage: "download",
			Query:       query,
			User: User{
				Name:   "管理员",
				Avatar: "/static/images/avatar.png",
			},
			Stats: DashboardStats{
				TotalModules:      156,
				TotalDownloads:    2847,
				NewModulesToday:   12,
				ActiveRepositories: 8,
			},
		}

		if query != "" {
			// 模拟搜索结果
			data.Modules = []Module{
				{
					Path:        "github.com/gin-gonic/gin",
					Version:     "v1.9.1",
					Description: "Gin is a HTTP web framework written in Go",
					Downloads:   15420,
					UpdatedAt:   "2023-12-01",
					License:     "MIT",
					IsPrivate:   false,
				},
			}
		} else {
			// 显示热门模块
			data.PopularModules = []Module{
				{
					Path:        "github.com/gin-gonic/gin",
					Version:     "v1.9.1",
					Description: "Gin is a HTTP web framework written in Go",
					Downloads:   15420,
					UpdatedAt:   "2023-12-01",
					License:     "MIT",
					IsPrivate:   false,
				},
				{
					Path:        "github.com/gorilla/mux",
					Version:     "v1.8.0",
					Description: "A powerful HTTP router and URL matcher",
					Downloads:   8930,
					UpdatedAt:   "2023-11-15",
					License:     "BSD-3-Clause",
					IsPrivate:   false,
				},
			}
		}

		err := tm.RenderTemplate(w, "download_refactored.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	// 创建模板管理器
	tm := &TemplateManager{
		TemplateDir: "./templates",
		Templates:   make(map[string]*template.Template),
	}

	// 解析模板
	err := tm.ParseTemplates()
	if err != nil {
		log.Fatal("解析模板失败:", err)
	}

	// 设置路由
	http.HandleFunc("/admin/dashboard", DashboardHandler(tm))
	http.HandleFunc("/admin/repositories", RepositoriesHandler(tm))
	http.HandleFunc("/admin/upload", UploadHandler(tm))
	http.HandleFunc("/admin/download", DownloadHandler(tm))

	// 启动服务器
	log.Println("服务器启动在 :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}