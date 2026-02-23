package admin

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// 模拟数据 - 常用的Go模块路径
var popularModulePaths = []string{
	"github.com/gin-gonic/gin",
	"github.com/gorilla/mux",
	"github.com/spf13/cobra",
	"github.com/stretchr/testify",
	"github.com/prometheus/client_golang",
	"github.com/sirupsen/logrus",
	"github.com/go-redis/redis",
	"github.com/jinzhu/gorm",
	"github.com/golang/protobuf",
	"github.com/pkg/errors",
	"github.com/uber-go/zap",
	"github.com/go-sql-driver/mysql",
	"github.com/lib/pq",
	"github.com/dgrijalva/jwt-go",
	"github.com/go-chi/chi",
}

// 模拟数据 - 模块版本
var moduleVersions = []string{
	"v1.0.0", "v1.1.0", "v1.2.0", "v2.0.0", "v2.1.0",
	"v0.1.0", "v0.2.0", "v0.3.0", "v0.4.0", "v0.5.0",
}

// 模拟数据 - 模块描述
var moduleDescriptions = []string{
	"HTTP web framework",
	"A powerful HTTP router and URL matcher",
	"A Commander for modern Go CLI interactions",
	"A toolkit with common assertions and mocks",
	"Prometheus instrumentation library for Go applications",
	"Structured, pluggable logging for Go",
	"Type-safe Redis client for Golang",
	"The fantastic ORM library for Golang",
	"Go support for Protocol Buffers",
	"Simple error handling primitives",
	"Blazing fast, structured, leveled logging",
	"MySQL driver for Go's database/sql package",
	"Pure Go Postgres driver for database/sql",
	"Golang implementation of JSON Web Tokens",
	"Lightweight, idiomatic and composable router for Go",
}

// 模拟的模块数据缓存
var mockModules []ModuleData

// 初始化模拟数据
func init() {
	// 生成模拟模块数据
	mockModules = generateMockModules()
}

// generateMockModules 生成模拟模块数据
func generateMockModules() []ModuleData {
	modules := make([]ModuleData, 0, 50)

	// 为每个流行模块生成多个版本
	for _, path := range popularModulePaths {
		descIdx := rand.Intn(len(moduleDescriptions))
		description := moduleDescriptions[descIdx]

		// 每个模块生成2-5个版本
		versionCount := 2 + rand.Intn(4)
		for i := 0; i < versionCount; i++ {
			version := moduleVersions[rand.Intn(len(moduleVersions))]
			downloads := 10 + rand.Intn(990) // 10-1000次下载
			size := int64(10000 + rand.Intn(990000)) // 10KB-1MB大小
			
			// 最近访问时间在过去30天内
			daysAgo := rand.Intn(30)
			lastAccess := time.Now().Add(-time.Duration(daysAgo) * 24 * time.Hour)
			
			modules = append(modules, ModuleData{
				Path:        path,
				Version:     version,
				Size:        size,
				Downloads:   downloads,
				LastAccess:  lastAccess,
				Description: description,
			})
		}
	}

	return modules
}

// downloadModulesAPIHandler 处理模块列表API请求
func downloadModulesAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 获取查询参数
	query := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// 解析limit和offset参数
	limit := 20 // 默认限制为20条
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // 默认偏移为0
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// 过滤和分页模块数据
	var filteredModules []ModuleData
	for _, module := range mockModules {
		// 如果有查询参数，则过滤模块
		if query != "" && !strings.Contains(strings.ToLower(module.Path), strings.ToLower(query)) {
			continue
		}
		filteredModules = append(filteredModules, module)
	}

	// 计算总数
	total := len(filteredModules)

	// 应用分页
	start := offset
	end := offset + limit
	if start >= total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedModules := filteredModules
	if start < end {
		pagedModules = filteredModules[start:end]
	} else {
		pagedModules = []ModuleData{}
	}

	// 构建响应
	response := map[string]interface{}{
		"modules": pagedModules,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	}

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// downloadModuleDetailAPIHandler 处理模块详情API请求
func downloadModuleDetailAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 获取路径参数
	vars := mux.Vars(r)
	modulePath := vars["path"]

	// 查找模块
	var foundModules []ModuleData
	for _, module := range mockModules {
		if module.Path == modulePath {
			foundModules = append(foundModules, module)
		}
	}

	if len(foundModules) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "模块未找到"})
		return
	}

	// 返回模块详情
	response := map[string]interface{}{
		"module":  foundModules[0], // 返回第一个版本作为模块详情
		"versions": len(foundModules),
	}

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// downloadModuleVersionsAPIHandler 处理模块版本列表API请求
func downloadModuleVersionsAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 获取路径参数
	vars := mux.Vars(r)
	modulePath := vars["path"]

	// 查找模块的所有版本
	var versions []ModuleData
	for _, module := range mockModules {
		if module.Path == modulePath {
			versions = append(versions, module)
		}
	}

	if len(versions) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "模块未找到"})
		return
	}

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(versions); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// downloadStatsAPIHandler 处理下载统计API请求
func downloadStatsAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 获取下载统计数据
	stats := getDownloadStats()

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// getDownloadStats 获取下载统计数据
func getDownloadStats() *DownloadStats {
	// 计算总下载量
	totalDownloads := 0
	for _, module := range mockModules {
		totalDownloads += module.Downloads
	}

	// 获取热门模块
	popularModules := getPopularModules(10)

	// 获取下载趋势
	downloadTrends := getDownloadTrends()

	// 获取最近下载的模块
	recentModules := getRecentModules(10)

	return &DownloadStats{
		TotalDownloads: totalDownloads,
		TotalModules:   len(mockModules),
		PopularModules: popularModules,
		DownloadTrends: downloadTrends,
		RecentModules:  recentModules,
	}
}

// downloadPopularAPIHandler 处理热门模块API请求
func downloadPopularAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 获取查询参数
	limitStr := r.URL.Query().Get("limit")

	// 解析limit参数
	limit := 10 // 默认限制为10条
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// 获取热门模块
	popularModules := getPopularModules(limit)

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(popularModules); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// getPopularModules 获取热门模块
func getPopularModules(limit int) []PopularModule {
	// 创建模块下载量映射
	moduleDownloads := make(map[string]int)
	for _, module := range mockModules {
		moduleDownloads[module.Path] += module.Downloads
	}

	// 转换为PopularModule切片
	popularModules := make([]PopularModule, 0, len(moduleDownloads))
	for path, downloads := range moduleDownloads {
		popularModules = append(popularModules, PopularModule{
			Path:      path,
			Downloads: downloads,
		})
	}

	// 按下载量排序
	sortPopularModulesByDownloads(popularModules)

	// 限制数量
	if len(popularModules) > limit {
		popularModules = popularModules[:limit]
	}

	return popularModules
}

// sortPopularModulesByDownloads 按下载量排序热门模块
func sortPopularModulesByDownloads(modules []PopularModule) {
	// 简单的冒泡排序
	for i := 0; i < len(modules)-1; i++ {
		for j := 0; j < len(modules)-i-1; j++ {
			if modules[j].Downloads < modules[j+1].Downloads {
				modules[j], modules[j+1] = modules[j+1], modules[j]
			}
		}
	}
}

// downloadRecentAPIHandler 处理最近下载模块API请求
func downloadRecentAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 获取查询参数
	limitStr := r.URL.Query().Get("limit")

	// 解析limit参数
	limit := 10 // 默认限制为10条
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// 获取最近下载的模块
	recentModules := getRecentModules(limit)

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(recentModules); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// getRecentModules 获取最近下载的模块
func getRecentModules(limit int) []RecentModuleAccess {
	// 创建RecentModuleAccess切片
	recentModules := make([]RecentModuleAccess, 0, len(mockModules))
	for _, module := range mockModules {
		recentModules = append(recentModules, RecentModuleAccess{
			Path:       module.Path,
			Version:    module.Version,
			AccessTime: module.LastAccess,
		})
	}

	// 按访问时间排序
	sortRecentModulesByAccessTime(recentModules)

	// 限制数量
	if len(recentModules) > limit {
		recentModules = recentModules[:limit]
	}

	return recentModules
}

// sortRecentModulesByAccessTime 按访问时间排序最近下载的模块
func sortRecentModulesByAccessTime(modules []RecentModuleAccess) {
	// 简单的冒泡排序
	for i := 0; i < len(modules)-1; i++ {
		for j := 0; j < len(modules)-i-1; j++ {
			if modules[j].AccessTime.Before(modules[j+1].AccessTime) {
				modules[j], modules[j+1] = modules[j+1], modules[j]
			}
		}
	}
}

// getDownloadTrends 获取下载趋势数据
func getDownloadTrends() []DownloadTrend {
	// 生成过去30天的下载趋势数据
	trends := make([]DownloadTrend, 30)

	for i := 0; i < 30; i++ {
		date := time.Now().AddDate(0, 0, -29+i)
		count := 100 + rand.Intn(900) // 100-1000次下载

		trends[i] = DownloadTrend{
			Date:  date.Format("2006-01-02"),
			Count: count,
		}
	}

	return trends
}