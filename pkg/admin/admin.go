package admin

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
)

const Admin="/admin"

// 嵌入模板文件
//go:embed templates/layout.html
var layoutTemplate string

//go:embed templates/styles.html
var stylesTemplate string

//go:embed templates/dashboard-styles.html
var dashboardStylesTemplate string

//go:embed templates/dashboard_refactored.html
var dashboardTemplate string

//go:embed templates/repository-styles.html
var repositoryStylesTemplate string

//go:embed templates/repositories_refactored.html
var repositoriesTemplate string

//go:embed templates/upload-styles.html
var uploadStylesTemplate string

//go:embed templates/upload_refactored.html
var uploadTemplate string

//go:embed templates/download-styles.html
var downloadStylesTemplate string

//go:embed templates/download_refactored.html
var downloadTemplate string

//go:embed templates/admin_refactored.html
var adminTemplate string

//go:embed templates/header.html
var headerTemplate string

//go:embed templates/nav.html
var navTemplate string

//go:embed templates/sidebar.html
var sidebarTemplate string

//go:embed templates/footer.html
var footerTemplate string

// AdminPageData 定义模板数据结构
type AdminPageData struct {
	Title          string
	SystemName     string
	Username       string
	PageTitle      string
	WelcomeMessage string
}

// DashboardData 仪表盘数据结构
type DashboardData struct {
	Title           string
	SystemName      string
	Username        string
	Stats           DashboardStats
	RecentActivities []Activity
}

// DashboardStats 统计数据
type DashboardStats struct {
	TotalModules        int
	TotalDownloads      int
	ActiveUsers         int
	StorageUsed         string
	NewModulesToday     int
	ActiveRepositories  int
	RepositoryUptime    int
	TodayDownloads      int
	DownloadGrowth      int
	StorageUsagePercent int
}

// Activity 活动记录
type Activity struct {
	Type        string
	Description string
	Time        string
	User        string
	Icon        string
	Title       string
}

// RepositoryData 仓库管理数据结构
type RepositoryData struct {
	Title        string
	SystemName   string
	Username     string
	PageTitle    string
	Repositories []Repository
}

// Repository 仓库信息
type Repository struct {
	Name        string
	Type        string
	URL         string
	Status      string
	ModuleCount int
	LastUpdated string
}

// UploadData 上传页面数据结构
type UploadData struct {
	Title         string
	SystemName    string
	Username      string
	PageTitle     string
	Repositories  []Repository
	RecentUploads []Upload
}

// Upload 上传记录
type Upload struct {
	Module     string
	Version    string
	Repository string
	UploadTime string
	Status     string
}

// DownloadData 下载页面数据结构
type DownloadData struct {
	Title           string
	SystemName      string
	Username        string
	PageTitle       string
	Query           string
	Modules         []Module
	ResultCount     int
	CurrentPage     int
	TotalPages      int
	PageNumbers     []int
	RecentDownloads []Download
}

// Module 模块信息
type Module struct {
	Name          string
	Version       string
	Description   string
	Status        string
	UpdatedAt     string
	DownloadCount int
	Stars         int
	Tags          []string
}

// Download 下载记录
type Download struct {
	Module       string
	Version      string
	DownloadTime string
	FileSize     string
}

func AdminHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置内容类型为HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		log.Println("AdminHandler: Starting template parsing")
		
		// 解析嵌入的模板
		tmpl, err := template.New("layout").Parse(layoutTemplate)
		if err != nil {
			log.Printf("AdminHandler: Error parsing layout template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing layout template: " + err.Error()))
			return
		}
		log.Println("AdminHandler: Layout template parsed successfully")
		
		// 添加其他模板
		log.Println("AdminHandler: Parsing styles template")
		tmpl, err = tmpl.New("base-styles").Parse(stylesTemplate)
		if err != nil {
			log.Printf("AdminHandler: Error parsing styles template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing styles template: " + err.Error()))
			return
		}
		
		log.Println("AdminHandler: Parsing header template")
		tmpl, err = tmpl.New("header").Parse(headerTemplate)
		if err != nil {
			log.Printf("AdminHandler: Error parsing header template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing header template: " + err.Error()))
			return
		}
		
		log.Println("AdminHandler: Parsing nav template")
		tmpl, err = tmpl.New("nav").Parse(navTemplate)
		if err != nil {
			log.Printf("AdminHandler: Error parsing nav template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing nav template: " + err.Error()))
			return
		}
		
		log.Println("AdminHandler: Parsing sidebar template")
		tmpl, err = tmpl.New("sidebar").Parse(sidebarTemplate)
		if err != nil {
			log.Printf("AdminHandler: Error parsing sidebar template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing sidebar template: " + err.Error()))
			return
		}
		
		log.Println("AdminHandler: Parsing footer template")
		tmpl, err = tmpl.New("footer").Parse(footerTemplate)
		if err != nil {
			log.Printf("AdminHandler: Error parsing footer template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing footer template: " + err.Error()))
			return
		}
		
		// 添加页面模板
		log.Println("AdminHandler: Parsing admin template")
		tmpl, err = tmpl.New("admin_refactored").Parse(adminTemplate)
		if err != nil {
			log.Printf("AdminHandler: Error parsing admin template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing admin template: " + err.Error()))
			return
		}
		
		// 准备模板数据
		log.Println("AdminHandler: Preparing template data")
		data := AdminPageData{
			Title:          "Athens仓库管理系统",
			SystemName:     "Athens Go模块仓库管理系统",
			Username:       "管理员",
			PageTitle:      "Go模块仓库管理",
			WelcomeMessage: "欢迎使用Athens仓库管理系统，您可以在此管理所有Go模块仓库。本系统仅支持hosted模式的Go模块仓库。",
		}
		
		// 渲染模板
		log.Println("AdminHandler: Executing template")
		err = tmpl.ExecuteTemplate(w, "admin_refactored", data)
		if err != nil {
			log.Printf("AdminHandler: Error executing admin template: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error executing admin template: " + err.Error()))
			return
		}
		log.Println("AdminHandler: Template executed successfully")
    })
}

// DashboardHandler 仪表盘处理函数
func DashboardHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		// 解析嵌入的模板
		tmpl, err := template.New("layout").Parse(layoutTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing layout template: " + err.Error()))
			return
		}
		
		// 添加样式模板
		tmpl, err = tmpl.New("base-styles").Parse(stylesTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing styles template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("dashboard-styles").Parse(dashboardStylesTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing dashboard styles template: " + err.Error()))
			return
		}
		
		// 添加组件模板
		tmpl, err = tmpl.New("header").Parse(headerTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing header template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("nav").Parse(navTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing nav template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("sidebar").Parse(sidebarTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing sidebar template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("footer").Parse(footerTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing footer template: " + err.Error()))
			return
		}
		
		// 添加页面模板
		tmpl, err = tmpl.New("dashboard_refactored").Parse(dashboardTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing dashboard template: " + err.Error()))
			return
		}
		
		// 模拟数据
		data := DashboardData{
			Title:      "Athens仓库管理系统 - 仪表盘",
			SystemName: "Athens Go模块仓库管理系统",
			Username:   "管理员",
			Stats: DashboardStats{
				TotalModules:        1250,
				TotalDownloads:      8934,
				ActiveUsers:         45,
				StorageUsed:         "2.3",
				NewModulesToday:     23,
				ActiveRepositories:  8,
				RepositoryUptime:    99,
				TodayDownloads:      456,
				DownloadGrowth:      12,
				StorageUsagePercent: 65,
			},
			RecentActivities: []Activity{
				{Type: "upload", Title: "上传了新模块 github.com/gin-gonic/gin v1.9.1", Description: "上传了新模块 github.com/gin-gonic/gin v1.9.1", Time: "2分钟前", User: "开发者A", Icon: "📦"},
				{Type: "download", Title: "下载了模块 github.com/gorilla/mux v1.8.0", Description: "下载了模块 github.com/gorilla/mux v1.8.0", Time: "5分钟前", User: "开发者B", Icon: "⬇️"},
				{Type: "system", Title: "系统自动清理了过期缓存", Description: "系统自动清理了过期缓存", Time: "10分钟前", User: "系统", Icon: "🔧"},
			},
		}
		
		err = tmpl.ExecuteTemplate(w, "dashboard_refactored", data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error executing dashboard template: " + err.Error()))
			return
		}
	})
}

// RepositoriesHandler 仓库管理处理函数
func RepositoriesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		// 解析嵌入的模板
		tmpl, err := template.New("layout").Parse(layoutTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing layout template: " + err.Error()))
			return
		}
		
		// 添加样式和组件模板
		tmpl, err = tmpl.New("base-styles").Parse(stylesTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing styles template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("repository-styles").Parse(repositoryStylesTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing repository styles template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("header").Parse(headerTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing header template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("nav").Parse(navTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing nav template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("sidebar").Parse(sidebarTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing sidebar template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("footer").Parse(footerTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing footer template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("repositories_refactored").Parse(repositoriesTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing repositories template: " + err.Error()))
			return
		}
		
		// 模拟数据
		data := RepositoryData{
			Title:      "Athens仓库管理系统 - 仓库管理",
			SystemName: "Athens Go模块仓库管理系统",
			Username:   "管理员",
			PageTitle:  "Go模块仓库管理",
			Repositories: []Repository{
				{Name: "public-modules", Type: "hosted", URL: "/var/athens/storage/public", Status: "online", ModuleCount: 856, LastUpdated: "2023-12-01 14:30:00"},
				{Name: "private-modules", Type: "hosted", URL: "/var/athens/storage/private", Status: "online", ModuleCount: 234, LastUpdated: "2023-12-01 13:45:00"},
				{Name: "proxy-cache", Type: "proxy", URL: "/var/athens/cache/proxy", Status: "online", ModuleCount: 1567, LastUpdated: "2023-12-01 15:20:00"},
			},
		}
		
		err = tmpl.ExecuteTemplate(w, "repositories_refactored", data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error executing repositories template: " + err.Error()))
			return
		}
	})
}

// UploadHandler 上传页面处理函数
func UploadHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		// 解析嵌入的模板
		tmpl, err := template.New("layout").Parse(layoutTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing layout template: " + err.Error()))
			return
		}
		
		// 添加样式和组件模板
		tmpl, err = tmpl.New("base-styles").Parse(stylesTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing styles template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("upload-styles").Parse(uploadStylesTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing upload styles template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("header").Parse(headerTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing header template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("nav").Parse(navTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing nav template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("sidebar").Parse(sidebarTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing sidebar template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("footer").Parse(footerTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing footer template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("upload_refactored").Parse(uploadTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing upload template: " + err.Error()))
			return
		}
		
		// 模拟数据
		data := UploadData{
			Title:      "Athens仓库管理系统 - 上传模块",
			SystemName: "Athens Go模块仓库管理系统",
			Username:   "管理员",
			PageTitle:  "上传Go模块",
			Repositories: []Repository{
				{Name: "public-modules", Type: "hosted"},
				{Name: "private-modules", Type: "hosted"},
			},
			RecentUploads: []Upload{
				{Module: "github.com/gin-gonic/gin", Version: "v1.9.1", Repository: "public-modules", UploadTime: "2023-12-01 14:30:00", Status: "success"},
				{Module: "github.com/gorilla/mux", Version: "v1.8.0", Repository: "public-modules", UploadTime: "2023-12-01 13:45:00", Status: "success"},
				{Module: "github.com/company/private-lib", Version: "v2.1.0", Repository: "private-modules", UploadTime: "2023-12-01 12:20:00", Status: "processing"},
			},
		}
		
		err = tmpl.ExecuteTemplate(w, "upload_refactored", data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error executing upload template: " + err.Error()))
			return
		}
	})
}

// DownloadHandler 下载页面处理函数
func DownloadHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		// 创建模板并添加自定义函数
		tmpl := template.New("layout").Funcs(template.FuncMap{
			"add": func(a, b int) int { return a + b },
			"sub": func(a, b int) int { return a - b },
		})
		
		// 解析嵌入的模板
		tmpl, err := tmpl.Parse(layoutTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing layout template: " + err.Error()))
			return
		}
		
		// 添加样式和组件模板
		tmpl, err = tmpl.New("base-styles").Parse(stylesTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing styles template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("download-styles").Parse(downloadStylesTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing download styles template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("header").Parse(headerTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing header template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("nav").Parse(navTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing nav template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("sidebar").Parse(sidebarTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing sidebar template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("footer").Parse(footerTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing footer template: " + err.Error()))
			return
		}
		
		tmpl, err = tmpl.New("download_refactored").Parse(downloadTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing download template: " + err.Error()))
			return
		}
		
		// 获取查询参数
		query := r.URL.Query().Get("query")
		
		// 模拟数据
		data := DownloadData{
			Title:       "Athens仓库管理系统 - 下载模块",
			SystemName:  "Athens Go模块仓库管理系统",
			Username:    "管理员",
			PageTitle:   "搜索和下载Go模块",
			Query:       query,
			ResultCount: 3,
			CurrentPage: 1,
			TotalPages:  1,
			PageNumbers: []int{1},
			Modules: []Module{
				{Name: "github.com/gin-gonic/gin", Version: "v1.9.1", Description: "Gin是一个用Go编写的HTTP Web框架", Status: "available", UpdatedAt: "2023-11-15", DownloadCount: 1234, Stars: 75000, Tags: []string{"web", "framework", "http"}},
				{Name: "github.com/gorilla/mux", Version: "v1.8.0", Description: "强大的HTTP路由器和URL匹配器", Status: "cached", UpdatedAt: "2023-10-20", DownloadCount: 856, Stars: 20000, Tags: []string{"router", "http", "mux"}},
				{Name: "github.com/labstack/echo", Version: "v4.11.3", Description: "高性能、极简的Go Web框架", Status: "available", UpdatedAt: "2023-11-28", DownloadCount: 567, Stars: 28000, Tags: []string{"web", "framework", "fast"}},
			},
			RecentDownloads: []Download{
				{Module: "github.com/gin-gonic/gin", Version: "v1.9.1", DownloadTime: "2023-12-01 14:30:00", FileSize: "2.3 MB"},
				{Module: "github.com/gorilla/mux", Version: "v1.8.0", DownloadTime: "2023-12-01 13:45:00", FileSize: "1.8 MB"},
			},
		}
		
		err = tmpl.ExecuteTemplate(w, "download_refactored", data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error executing download template: " + err.Error()))
			return
		}
	})
}
