package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/leimeng-go/athens/pkg/build"
	"github.com/gorilla/mux"
)

// RegisterAPIHandlers 注册API路由处理函数
func RegisterAPIHandlers(r *mux.Router) {
	// 系统状态API
	r.HandleFunc(Admin+"/api/system/status", systemStatusAPIHandler)
	
	// 仪表盘数据API
	r.HandleFunc(Admin+"/api/dashboard", dashboardAPIHandler)
	
	// 最近活动API
	r.HandleFunc(Admin+"/api/activities/recent", recentActivitiesAPIHandler)
	
	// 系统设置API
	r.HandleFunc(Admin+"/api/settings", systemSettingsAPIHandler)
	
	// 模块下载相关API
	r.HandleFunc(Admin+"/api/download/modules", downloadModulesAPIHandler)
	r.HandleFunc(Admin+"/api/download/modules/{path:.+}/versions", downloadModuleVersionsAPIHandler)
	r.HandleFunc(Admin+"/api/download/modules/{path:.+}", downloadModuleDetailAPIHandler)
	r.HandleFunc(Admin+"/api/download/stats", downloadStatsAPIHandler)
	r.HandleFunc(Admin+"/api/download/popular", downloadPopularAPIHandler)
	r.HandleFunc(Admin+"/api/download/recent", downloadRecentAPIHandler)
	
	// 仓库相关API
	r.HandleFunc(Admin+"/api/repositories", repositoriesAPIHandler)
	r.HandleFunc(Admin+"/api/repositories/batch-delete", repositoryBatchDeleteAPIHandler)
	r.HandleFunc(Admin+"/api/repositories/{id}", repositoryDetailAPIHandler)
	
	// 模块上传相关API
	r.HandleFunc(Admin+"/api/upload/module", uploadModuleAPIHandler)
	r.HandleFunc(Admin+"/api/upload/import-url", uploadImportURLAPIHandler)
	r.HandleFunc(Admin+"/api/upload/tasks", uploadTasksAPIHandler)
	r.HandleFunc(Admin+"/api/upload/tasks/{taskId}/cancel", uploadTaskCancelAPIHandler)
	r.HandleFunc(Admin+"/api/upload/tasks/{taskId}", uploadTaskDetailAPIHandler)
}

// systemStatusAPIHandler 处理系统状态API请求
func systemStatusAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	// 获取系统信息
	systemInfo := getSystemInfo()
	
	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(systemInfo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// 系统启动时间
var startTime = time.Now()

// SystemInfo 表示系统状态信息
type SystemInfo struct {
	Status      string `json:"status"`
	Uptime      string `json:"uptime"`
	Version     string `json:"version"`
	GoVersion   string `json:"goVersion"`
	MemoryUsage string `json:"memoryUsage"`
	CPUUsage    string `json:"cpuUsage"`
}

// getSystemInfo 获取系统状态信息
func getSystemInfo() *SystemInfo {
	// 计算运行时间
	uptime := time.Since(startTime)
	uptimeStr := formatUptime(uptime)
	
	// 获取内存使用情况
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memoryUsage := formatMemory(memStats.Alloc)
	
	// 获取版本信息
	buildInfo := build.Data()
	
	// 创建系统状态对象
	return &SystemInfo{
		Status:      "healthy",
		Uptime:      uptimeStr,
		Version:     buildInfo.Version,
		GoVersion:   runtime.Version(),
		MemoryUsage: memoryUsage,
		CPUUsage:    "N/A", // 简化实现，实际获取CPU使用率需要更复杂的逻辑
	}
}

// formatUptime 格式化运行时间
func formatUptime(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// formatMemory 格式化内存大小
func formatMemory(bytes uint64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	
	var value float64
	var unit string
	
	switch {
	case bytes >= GB:
		value = float64(bytes) / GB
		unit = "GB"
	case bytes >= MB:
		value = float64(bytes) / MB
		unit = "MB"
	case bytes >= KB:
		value = float64(bytes) / KB
		unit = "KB"
	default:
		value = float64(bytes)
		unit = "B"
	}
	
	return fmt.Sprintf("%.2f %s", value, unit)
}

// dashboardAPIHandler 处理仪表盘数据API请求
func dashboardAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	// 获取仪表盘数据
	dashboardData := getDashboardData()
	
	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(dashboardData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// getDashboardData 获取仪表盘数据
// 返回仪表盘数据，包括统计信息、下载趋势、热门模块和最近活动
func getDashboardData() *DashboardData {
	// 这里应该从数据库或其他数据源获取实际数据
	// 目前返回模拟数据用于演示
	return &DashboardData{
		Stats: Stats{
			TotalModules:     1250,
			TotalDownloads:   45678,
			TotalRepositories: 35,
			StorageUsed:      "2.3 GB",
		},
		DownloadTrend: []DownloadTrend{
			{Date: "2023-01-01", Count: 120},
			{Date: "2023-01-02", Count: 145},
			{Date: "2023-01-03", Count: 132},
			{Date: "2023-01-04", Count: 167},
			{Date: "2023-01-05", Count: 189},
			{Date: "2023-01-06", Count: 156},
			{Date: "2023-01-07", Count: 123},
		},
		PopularModules: []PopularModule{
			{Path: "github.com/gorilla/mux", Downloads: 12345},
			{Path: "github.com/gin-gonic/gin", Downloads: 10234},
			{Path: "github.com/spf13/cobra", Downloads: 8765},
			{Path: "github.com/stretchr/testify", Downloads: 7654},
			{Path: "github.com/prometheus/client_golang", Downloads: 6543},
		},
		RecentActivities: []RecentActivity{
			{ID: "1", Type: "download", Message: "模块 github.com/gorilla/mux@v1.8.0 被下载", Timestamp: "2023-01-07T12:34:56Z"},
			{ID: "2", Type: "upload", Message: "模块 github.com/gin-gonic/gin@v1.9.0 被上传", Timestamp: "2023-01-07T10:23:45Z"},
			{ID: "3", Type: "system", Message: "系统更新完成", Timestamp: "2023-01-06T22:12:34Z"},
			{ID: "4", Type: "download", Message: "模块 github.com/spf13/cobra@v1.6.1 被下载", Timestamp: "2023-01-06T18:45:12Z"},
			{ID: "5", Type: "upload", Message: "模块 github.com/stretchr/testify@v1.8.1 被上传", Timestamp: "2023-01-06T15:34:23Z"},
		},
	}
}

// recentActivitiesAPIHandler 处理最近活动API请求
func recentActivitiesAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	// 获取查询参数
	limit := 10 // 默认限制为10条
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	// 获取最近活动数据
	activities := getRecentActivities(limit)
	
	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(activities); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// getRecentActivities 获取最近活动数据
// 根据指定的限制数量返回最近的活动数据
func getRecentActivities(limit int) []RecentActivity {
	// 这里应该从数据库或其他数据源获取实际数据
	// 目前返回模拟数据用于演示
	allActivities := []RecentActivity{
		{ID: "1", Type: "download", Message: "模块 github.com/gorilla/mux@v1.8.0 被下载", Timestamp: "2023-01-07T12:34:56Z"},
		{ID: "2", Type: "upload", Message: "模块 github.com/gin-gonic/gin@v1.9.0 被上传", Timestamp: "2023-01-07T10:23:45Z"},
		{ID: "3", Type: "system", Message: "系统更新完成", Timestamp: "2023-01-06T22:12:34Z"},
		{ID: "4", Type: "download", Message: "模块 github.com/spf13/cobra@v1.6.1 被下载", Timestamp: "2023-01-06T18:45:12Z"},
		{ID: "5", Type: "upload", Message: "模块 github.com/stretchr/testify@v1.8.1 被上传", Timestamp: "2023-01-06T15:34:23Z"},
		{ID: "6", Type: "download", Message: "模块 github.com/sirupsen/logrus@v1.9.0 被下载", Timestamp: "2023-01-06T14:23:12Z"},
		{ID: "7", Type: "upload", Message: "模块 github.com/go-chi/chi@v5.0.8 被上传", Timestamp: "2023-01-06T12:11:45Z"},
		{ID: "8", Type: "system", Message: "系统备份完成", Timestamp: "2023-01-05T23:45:12Z"},
		{ID: "9", Type: "download", Message: "模块 github.com/pkg/errors@v0.9.1 被下载", Timestamp: "2023-01-05T18:34:23Z"},
		{ID: "10", Type: "upload", Message: "模块 github.com/go-redis/redis@v8.11.5 被上传", Timestamp: "2023-01-05T15:12:34Z"},
	}
	
	// 限制返回数量
	if limit > len(allActivities) {
		limit = len(allActivities)
	}
	
	return allActivities[:limit]
}