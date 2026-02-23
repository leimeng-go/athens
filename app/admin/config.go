package admin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/leimeng-go/athens/pkg/config"
)

// UpdateRequest 配置更新请求结构
type UpdateRequest struct {
	// 支持更新的配置项
	GoEnv            *string `json:"go_env,omitempty"`
	GoBinary         *string `json:"go_binary,omitempty"`
	GoGetWorkers     *int    `json:"go_get_workers,omitempty"`
	ProtocolWorkers  *int    `json:"protocol_workers,omitempty"`
	LogLevel         *string `json:"log_level,omitempty"`
	CloudRuntime     *string `json:"cloud_runtime,omitempty"`
	EnablePprof      *bool   `json:"enable_pprof,omitempty"`
	PprofPort        *string `json:"pprof_port,omitempty"`
	StorageType      *string `json:"storage_type,omitempty"`
	Port             *string `json:"port,omitempty"`
	BasicAuthUser    *string `json:"basic_auth_user,omitempty"`
	BasicAuthPass    *string `json:"basic_auth_pass,omitempty"`
	ForceSSL         *bool   `json:"force_ssl,omitempty"`
	NetworkMode      *string `json:"network_mode,omitempty"`
	SingleFlightType *string `json:"single_flight_type,omitempty"`
	IndexType        *string `json:"index_type,omitempty"`
	ShutdownTimeout  *int    `json:"shutdown_timeout,omitempty"`
}

// 配置管理API处理函数

// getConfigAPIHandler 获取当前配置
func getConfigAPIHandler(w http.ResponseWriter, r *http.Request) {
	if configManager == nil {
		http.Error(w, "配置管理器未初始化", http.StatusInternalServerError)
		return
	}

	conf := configManager.Get()

	response := map[string]interface{}{
		"config":    conf,
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// updateConfigAPIHandler 更新配置
func updateConfigAPIHandler(w http.ResponseWriter, r *http.Request) {
	if configManager == nil {
		http.Error(w, "配置管理器未初始化", http.StatusInternalServerError)
		return
	}

	var req UpdateRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("读取请求体失败: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, fmt.Sprintf("解析 JSON 失败: %v", err), http.StatusBadRequest)
		return
	}

	// 验证请求
	if err := validateUpdateRequestImpl(&req); err != nil {
		http.Error(w, fmt.Sprintf("配置验证失败: %v", err), http.StatusBadRequest)
		return
	}

	// 更新配置文件
	if err := applyConfigUpdate(&req); err != nil {
		http.Error(w, fmt.Sprintf("更新配置失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回更新后的配置
	updatedConfig := configManager.Get()
	response := map[string]interface{}{
		"message":   "配置更新成功",
		"config":    updatedConfig,
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// resetConfigAPIHandler 重置配置到默认值
func resetConfigAPIHandler(w http.ResponseWriter, r *http.Request) {
	if configManager == nil {
		http.Error(w, "配置管理器未初始化", http.StatusInternalServerError)
		return
	}

	currentConfig := configManager.Get()
	defaultConfig := &config.Config{
		GoEnv:            "development",
		GoBinary:         "go",
		GoGetWorkers:     10,
		ProtocolWorkers:  30,
		LogLevel:         "debug",
		CloudRuntime:     "none",
		EnablePprof:      false,
		PprofPort:        ":3001",
		StorageType:      "memory",
		Port:             ":3000",
		NetworkMode:      "strict",
		SingleFlightType: "memory",
		IndexType:        "none",
		ShutdownTimeout:  60,
	}

	// 保留必要的配置项（如存储配置等）
	if currentConfig.Storage != nil {
		defaultConfig.Storage = currentConfig.Storage
	}
	if currentConfig.Index != nil {
		defaultConfig.Index = currentConfig.Index
	}
	if currentConfig.SingleFlight != nil {
		defaultConfig.SingleFlight = currentConfig.SingleFlight
	}

	// 写入默认配置到文件
	if err := writeConfigToFile(defaultConfig); err != nil {
		http.Error(w, fmt.Sprintf("重置配置失败: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":   "配置已重置为默认值",
		"config":    defaultConfig,
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// validateConfigAPIHandler 验证配置而不应用
func validateConfigAPIHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("读取请求体失败: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, fmt.Sprintf("解析 JSON 失败: %v", err), http.StatusBadRequest)
		return
	}

	if err := validateUpdateRequestImpl(&req); err != nil {
		http.Error(w, fmt.Sprintf("配置验证失败: %v", err), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"message":   "配置验证通过",
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// validateUpdateRequestImpl 验证更新请求的合法性
func validateUpdateRequestImpl(req *UpdateRequest) error {
	// 验证日志级别
	if req.LogLevel != nil {
		validLevels := map[string]bool{
			"debug": true, "info": true, "warn": true, "error": true,
		}
		if !validLevels[strings.ToLower(*req.LogLevel)] {
			return fmt.Errorf("无效的日志级别: %s", *req.LogLevel)
		}
	}

	// 验证网络模式
	if req.NetworkMode != nil {
		validModes := map[string]bool{
			"strict": true, "offline": true, "fallback": true,
		}
		if !validModes[strings.ToLower(*req.NetworkMode)] {
			return fmt.Errorf("无效的网络模式: %s", *req.NetworkMode)
		}
	}

	// 验证数字范围
	if req.GoGetWorkers != nil && *req.GoGetWorkers <= 0 {
		return fmt.Errorf("GoGetWorkers 必须大于 0")
	}
	if req.ProtocolWorkers != nil && *req.ProtocolWorkers <= 0 {
		return fmt.Errorf("ProtocolWorkers 必须大于 0")
	}
	if req.ShutdownTimeout != nil && *req.ShutdownTimeout < 0 {
		return fmt.Errorf("ShutdownTimeout 不能为负数")
	}

	// 验证端口格式
	if req.Port != nil && !strings.HasPrefix(*req.Port, ":") {
		return fmt.Errorf("端口必须以冒号开头，例如 :3000")
	}
	if req.PprofPort != nil && !strings.HasPrefix(*req.PprofPort, ":") {
		return fmt.Errorf("PprofPort 必须以冒号开头，例如 :3001")
	}

	return nil
}

// applyConfigUpdate 应用配置更新到文件
func applyConfigUpdate(req *UpdateRequest) error {
	if configManager.ConfigFile() == "" {
		return fmt.Errorf("未指定配置文件，无法更新")
	}

	// 读取当前配置
	_, err := os.ReadFile(configManager.ConfigFile())
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析现有配置
	currentConfig, err := config.ParseConfigFile(configManager.ConfigFile())
	if err != nil {
		return fmt.Errorf("解析当前配置失败: %w", err)
	}

	// 应用更新
	applyUpdatesToConfigImpl(currentConfig, req)

	// 生成新的 TOML 内容
	newContent, err := generateTomlContent(currentConfig)
	if err != nil {
		return fmt.Errorf("生成配置内容失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(configManager.ConfigFile(), []byte(newContent), 0600); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// applyUpdatesToConfigImpl 将更新应用到配置结构体
func applyUpdatesToConfigImpl(conf *config.Config, req *UpdateRequest) {
	if req.GoEnv != nil {
		conf.GoEnv = *req.GoEnv
	}
	if req.GoBinary != nil {
		conf.GoBinary = *req.GoBinary
	}
	if req.GoGetWorkers != nil {
		conf.GoGetWorkers = *req.GoGetWorkers
	}
	if req.ProtocolWorkers != nil {
		conf.ProtocolWorkers = *req.ProtocolWorkers
	}
	if req.LogLevel != nil {
		conf.LogLevel = *req.LogLevel
	}
	if req.CloudRuntime != nil {
		conf.CloudRuntime = *req.CloudRuntime
	}
	if req.EnablePprof != nil {
		conf.EnablePprof = *req.EnablePprof
	}
	if req.PprofPort != nil {
		conf.PprofPort = *req.PprofPort
	}
	if req.StorageType != nil {
		conf.StorageType = *req.StorageType
	}
	if req.Port != nil {
		conf.Port = *req.Port
	}
	if req.BasicAuthUser != nil {
		conf.BasicAuthUser = *req.BasicAuthUser
	}
	if req.BasicAuthPass != nil {
		conf.BasicAuthPass = *req.BasicAuthPass
	}
	if req.ForceSSL != nil {
		conf.ForceSSL = *req.ForceSSL
	}
	if req.NetworkMode != nil {
		conf.NetworkMode = *req.NetworkMode
	}
	if req.SingleFlightType != nil {
		conf.SingleFlightType = *req.SingleFlightType
	}
	if req.IndexType != nil {
		conf.IndexType = *req.IndexType
	}
	if req.ShutdownTimeout != nil {
		conf.ShutdownTimeout = *req.ShutdownTimeout
	}
}

// generateTomlContent 生成 TOML 格式的配置内容
func generateTomlContent(conf *config.Config) (string, error) {
	var builder strings.Builder

	// 基础配置
	builder.WriteString(fmt.Sprintf("GoEnv = \"%s\"\n", conf.GoEnv))
	builder.WriteString(fmt.Sprintf("GoBinary = \"%s\"\n", conf.GoBinary))
	builder.WriteString(fmt.Sprintf("GoGetWorkers = %d\n", conf.GoGetWorkers))
	builder.WriteString(fmt.Sprintf("ProtocolWorkers = %d\n", conf.ProtocolWorkers))
	builder.WriteString(fmt.Sprintf("LogLevel = \"%s\"\n", conf.LogLevel))
	builder.WriteString(fmt.Sprintf("CloudRuntime = \"%s\"\n", conf.CloudRuntime))
	builder.WriteString(fmt.Sprintf("EnablePprof = %t\n", conf.EnablePprof))
	builder.WriteString(fmt.Sprintf("PprofPort = \"%s\"\n", conf.PprofPort))
	builder.WriteString(fmt.Sprintf("StorageType = \"%s\"\n", conf.StorageType))
	builder.WriteString(fmt.Sprintf("Port = \"%s\"\n", conf.Port))
	builder.WriteString(fmt.Sprintf("BasicAuthUser = \"%s\"\n", conf.BasicAuthUser))
	builder.WriteString(fmt.Sprintf("BasicAuthPass = \"%s\"\n", conf.BasicAuthPass))
	builder.WriteString(fmt.Sprintf("ForceSSL = %t\n", conf.ForceSSL))
	builder.WriteString(fmt.Sprintf("NetworkMode = \"%s\"\n", conf.NetworkMode))
	builder.WriteString(fmt.Sprintf("SingleFlightType = \"%s\"\n", conf.SingleFlightType))
	builder.WriteString(fmt.Sprintf("IndexType = \"%s\"\n", conf.IndexType))
	builder.WriteString(fmt.Sprintf("ShutdownTimeout = %d\n", conf.ShutdownTimeout))

	// TODO: 添加 Storage 和 Index 配置的序列化
	// 这里简化处理，实际项目中需要完整实现

	return builder.String(), nil
}

// writeConfigToFile 将配置写入文件
func writeConfigToFile(conf *config.Config) error {
	if configManager.ConfigFile() == "" {
		return fmt.Errorf("未指定配置文件路径")
	}

	content, err := generateTomlContent(conf)
	if err != nil {
		return fmt.Errorf("生成配置内容失败: %w", err)
	}

	return os.WriteFile(configManager.ConfigFile(), []byte(content), 0600)
}
