package admin

import (
	"encoding/json"
	"net/http"
	"sync"
)

// 用于保护设置的并发访问
var settingsMutex sync.RWMutex

// 默认系统设置
var defaultSettings = SystemSettings{
	StoragePath:           "/var/lib/athens",
	MaxUploadSize:         50 * 1024 * 1024, // 50MB
	EnablePrivateModules:  true,
	EnableDownloadLogging: true,
	ProxyTimeout:          30, // 30秒
	CacheExpiration:       24, // 24小时
}

// 当前系统设置，初始化为默认设置
var currentSettings = defaultSettings

/**
 * systemSettingsAPIHandler 处理系统设置API请求
 * 支持GET和PUT方法，分别用于获取和更新系统设置
 */
func systemSettingsAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		// 获取系统设置
		handleGetSystemSettings(w, r)
	case http.MethodPut:
		// 更新系统设置
		handleUpdateSystemSettings(w, r)
	default:
		// 不支持的HTTP方法
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

/**
 * handleGetSystemSettings 处理获取系统设置的请求
 */
func handleGetSystemSettings(w http.ResponseWriter, r *http.Request) {
	// 获取当前设置（读锁）
	settingsMutex.RLock()
	settings := currentSettings
	settingsMutex.RUnlock()

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(settings); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

/**
 * handleUpdateSystemSettings 处理更新系统设置的请求
 */
func handleUpdateSystemSettings(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var newSettings SystemSettings
	if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
		return
	}

	// 验证设置（可以根据需要添加更多验证）
	if newSettings.MaxUploadSize <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "最大上传大小必须大于0"})
		return
	}

	if newSettings.ProxyTimeout <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "代理超时必须大于0"})
		return
	}

	if newSettings.CacheExpiration <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "缓存过期时间必须大于0"})
		return
	}

	// 更新设置（写锁）
	settingsMutex.Lock()
	currentSettings = newSettings
	settingsMutex.Unlock()

	// 返回更新后的设置
	if err := json.NewEncoder(w).Encode(newSettings); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

/**
 * GetSystemSettings 获取当前系统设置
 * 供内部使用的辅助函数
 */
func GetSystemSettings() SystemSettings {
	settingsMutex.RLock()
	defer settingsMutex.RUnlock()
	return currentSettings
}

/**
 * UpdateSystemSettings 更新系统设置
 * 供内部使用的辅助函数
 */
func UpdateSystemSettings(settings SystemSettings) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings = settings
}