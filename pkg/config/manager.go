package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Logger 是配置模块使用的日志接口，用于解耦具体实现
type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// defaultLogger 是一个空实现，用于在未显式传入 logger 时兜底
type defaultLogger struct{}

func (l *defaultLogger) Info(args ...interface{})                  {}
func (l *defaultLogger) Infof(format string, args ...interface{})  {}
func (l *defaultLogger) Error(args ...interface{})                 {}
func (l *defaultLogger) Errorf(format string, args ...interface{}) {}

// ConfigManager 管理配置热更新
type ConfigManager struct {
	mu           sync.RWMutex
	config       *Config
	configFile   string
	v            *viper.Viper
	onChange     []func(*Config)
	changeNotify chan struct{}
	logger       Logger
}

// NewConfigManager 创建配置管理器并启用热更新
func NewConfigManager(configFile string, logger Logger) (*ConfigManager, error) {
	if logger == nil {
		logger = &defaultLogger{}
	}

	// 初始加载配置
	initialConfig, err := Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load initial config: %w", err)
	}

	cm := &ConfigManager{
		config:       initialConfig,
		configFile:   configFile,
		changeNotify: make(chan struct{}, 1),
		logger:       logger,
	}

	// 如果没有指定配置文件，不启用热更新
	if configFile == "" {
		cm.logger.Info("No config file specified, hot reload disabled")
		return cm, nil
	}

	// 初始化 Viper
	cm.v = viper.New()
	cm.v.SetConfigFile(configFile)
	cm.v.SetConfigType("toml")

	// 监听配置文件变化
	cm.v.WatchConfig()
	cm.v.OnConfigChange(func(e fsnotify.Event) {
		cm.logger.Infof("Config file changed: %s, reloading...", e.Name)
		if err := cm.reloadConfig(); err != nil {
			cm.logger.Errorf("Failed to reload config: %v", err)
			return
		}
		cm.logger.Info("Config reloaded successfully")

		// 通知配置变更
		select {
		case cm.changeNotify <- struct{}{}:
		default:
		}

		// 调用所有变更回调
		cm.mu.RLock()
		config := cm.config
		callbacks := cm.onChange
		cm.mu.RUnlock()

		for _, callback := range callbacks {
			callback(config)
		}
	})

	cm.logger.Infof("Config hot reload enabled for: %s", configFile)
	return cm, nil
}

// reloadConfig 重新加载配置文件
func (cm *ConfigManager) reloadConfig() error {
	newConfig, err := ParseConfigFile(cm.configFile)
	if err != nil {
		return fmt.Errorf("parse config file failed: %w", err)
	}

	cm.mu.Lock()
	cm.config = newConfig
	cm.mu.Unlock()

	return nil
}

// Get 获取当前配置（线程安全）
func (cm *ConfigManager) Get() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// OnChange 注册配置变更回调函数
func (cm *ConfigManager) OnChange(callback func(*Config)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.onChange = append(cm.onChange, callback)
}

// ChangeNotify 返回配置变更通知通道
func (cm *ConfigManager) ChangeNotify() <-chan struct{} {
	return cm.changeNotify
}

// ConfigFile 返回配置文件路径
func (cm *ConfigManager) ConfigFile() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.configFile
}

// Update 更新配置并保存到文件（线程安全）
// updateFunc 接收当前配置的副本，修改后返回新配置
func (cm *ConfigManager) Update(updateFunc func(*Config) error) error {
	if cm.configFile == "" {
		return fmt.Errorf("no config file specified, cannot update")
	}

	// 获取当前配置的副本
	cm.mu.RLock()
	currentConfig := cm.config
	cm.mu.RUnlock()

	// 深拷贝配置，避免直接修改当前配置
	newConfig := &Config{}
	*newConfig = *currentConfig

	// 执行更新函数
	if err := updateFunc(newConfig); err != nil {
		return fmt.Errorf("update function failed: %w", err)
	}

	// 保存到文件
	if err := cm.saveToFile(newConfig); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// 更新内存中的配置
	cm.mu.Lock()
	cm.config = newConfig
	cm.mu.Unlock()

	// 通知配置变更
	select {
	case cm.changeNotify <- struct{}{}:
	default:
	}

	// 调用所有变更回调
	cm.mu.RLock()
	callbacks := cm.onChange
	cm.mu.RUnlock()

	for _, callback := range callbacks {
		callback(newConfig)
	}

	return nil
}

// Set 直接设置新配置并保存到文件（线程安全）
func (cm *ConfigManager) Set(newConfig *Config) error {
	if cm.configFile == "" {
		return fmt.Errorf("no config file specified, cannot set")
	}

	// 保存到文件
	if err := cm.saveToFile(newConfig); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// 更新内存中的配置
	cm.mu.Lock()
	cm.config = newConfig
	cm.mu.Unlock()

	// 通知配置变更
	select {
	case cm.changeNotify <- struct{}{}:
	default:
	}

	// 调用所有变更回调
	cm.mu.RLock()
	callbacks := cm.onChange
	cm.mu.RUnlock()

	for _, callback := range callbacks {
		callback(newConfig)
	}

	return nil
}

// saveToFile 将配置保存到文件
func (cm *ConfigManager) saveToFile(config *Config) error {
	// 创建或打开配置文件
	file, err := os.Create(cm.configFile)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	// 使用 TOML 编码器保存完整的配置结构
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// Reload 手动重新加载配置文件
func (cm *ConfigManager) Reload() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	newConfig, err := ParseConfigFile(cm.configFile)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	cm.config = newConfig

	// 通知配置变更
	select {
	case cm.changeNotify <- struct{}{}:
	default:
	}

	// 调用所有变更回调
	callbacks := cm.onChange
	for _, callback := range callbacks {
		callback(newConfig)
	}

	return nil
}
