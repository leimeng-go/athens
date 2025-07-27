package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/leimeng-go/athens/pkg/build"
)

// SystemStatus 表示系统状态信息
type SystemStatus struct {
	Status      string `json:"status"`
	Uptime      string `json:"uptime"`
	Version     string `json:"version"`
	GoVersion   string `json:"goVersion"`
	MemoryUsage string `json:"memoryUsage"`
	CPUUsage    string `json:"cpuUsage"`
}

// 系统启动时间
var startTime = time.Now()

// systemStatusHandler 处理系统状态API请求
func systemStatusHandler(w http.ResponseWriter, r *http.Request) {
	// 获取系统状态信息
	status := getSystemStatus()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(status); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// getSystemStatus 获取系统状态信息
func getSystemStatus() *SystemStatus {
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
	return &SystemStatus{
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