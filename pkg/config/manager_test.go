package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

// testLogger 实现 Logger 接口用于测试
type testLoggerImpl struct {
	*logrus.Logger
}

func (l *testLoggerImpl) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *testLoggerImpl) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *testLoggerImpl) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *testLoggerImpl) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

func init() {
	// 设置测试用的 logger 构造逻辑（这里直接使用 logrus）
	// 具体测试中通过 NewConfigManager(configFile, &testLoggerImpl{...}) 传入
}

// createValidTestConfig 创建一个通过验证的测试配置
func createValidTestConfig() *Config {
	config := defaultConfig()
	config.StorageType = "mongo" // 修改为支持的存储类型
	if config.Storage == nil {
		config.Storage = &Storage{}
	}
	if config.Storage.Mongo == nil {
		config.Storage.Mongo = &MongoConfig{
			URL: "mongodb://localhost:27017",
		}
	}
	return config
}

func TestConfigManager_Update(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.toml")

	// 创建初始配置
	initialConfig := createValidTestConfig()
	initialConfig.Port = ":3000"
	initialConfig.LogLevel = "info"

	// 保存初始配置
	if err := saveConfigToFile(configFile, initialConfig); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// 创建ConfigManager
	cm, err := NewConfigManager(configFile, &testLoggerImpl{Logger: logrus.New()})
	if err != nil {
		t.Fatalf("Failed to create ConfigManager: %v", err)
	}

	// 验证初始配置
	config := cm.Get()
	if config.Port != ":3000" {
		t.Errorf("Expected port :3000, got %s", config.Port)
	}

	// 测试Update方法
	err = cm.Update(func(c *Config) error {
		c.Port = ":8080"
		c.LogLevel = "debug"
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// 验证更新后的配置
	config = cm.Get()
	if config.Port != ":8080" {
		t.Errorf("Expected port :8080, got %s", config.Port)
	}
	if config.LogLevel != "debug" {
		t.Errorf("Expected log level debug, got %s", config.LogLevel)
	}

	// 验证配置已保存到文件
	reloadedConfig, err := ParseConfigFile(configFile)
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}
	if reloadedConfig.Port != ":8080" {
		t.Errorf("Expected reloaded port :8080, got %s", reloadedConfig.Port)
	}
}

func TestConfigManager_Set(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.toml")

	// 创建初始配置
	initialConfig := createValidTestConfig()
	if err := saveConfigToFile(configFile, initialConfig); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// 创建ConfigManager
	cm, err := NewConfigManager(configFile, &testLoggerImpl{Logger: logrus.New()})
	if err != nil {
		t.Fatalf("Failed to create ConfigManager: %v", err)
	}

	// 创建新配置
	newConfig := createValidTestConfig()
	newConfig.Port = ":9000"
	newConfig.LogLevel = "error"

	// 测试Set方法
	err = cm.Set(newConfig)
	if err != nil {
		t.Fatalf("Failed to set config: %v", err)
	}

	// 验证更新后的配置
	config := cm.Get()
	if config.Port != ":9000" {
		t.Errorf("Expected port :9000, got %s", config.Port)
	}
	if config.LogLevel != "error" {
		t.Errorf("Expected log level error, got %s", config.LogLevel)
	}
}

func TestConfigManager_Reload(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.toml")

	// 创建初始配置
	initialConfig := createValidTestConfig()
	initialConfig.Port = ":3000"
	if err := saveConfigToFile(configFile, initialConfig); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// 创建ConfigManager
	cm, err := NewConfigManager(configFile, &testLoggerImpl{Logger: logrus.New()})
	if err != nil {
		t.Fatalf("Failed to create ConfigManager: %v", err)
	}

	// 外部修改配置文件
	externalConfig := createValidTestConfig()
	externalConfig.Port = ":7000"
	if err := saveConfigToFile(configFile, externalConfig); err != nil {
		t.Fatalf("Failed to save external config: %v", err)
	}

	// 手动重新加载
	err = cm.Reload()
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	// 验证重新加载后的配置
	config := cm.Get()
	if config.Port != ":7000" {
		t.Errorf("Expected port :7000, got %s", config.Port)
	}
}

func TestConfigManager_OnChange(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.toml")

	// 创建初始配置
	initialConfig := createValidTestConfig()
	if err := saveConfigToFile(configFile, initialConfig); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// 创建ConfigManager
	cm, err := NewConfigManager(configFile, &testLoggerImpl{Logger: logrus.New()})
	if err != nil {
		t.Fatalf("Failed to create ConfigManager: %v", err)
	}

	// 注册回调
	callbackCalled := false
	var callbackConfig *Config
	cm.OnChange(func(c *Config) {
		callbackCalled = true
		callbackConfig = c
	})

	// 更新配置
	err = cm.Update(func(c *Config) error {
		c.Port = ":5000"
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// 验证回调被调用
	if !callbackCalled {
		t.Error("Expected callback to be called")
	}
	if callbackConfig == nil || callbackConfig.Port != ":5000" {
		t.Error("Expected callback to receive updated config")
	}
}

func TestConfigManager_ChangeNotify(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.toml")

	// 创建初始配置
	initialConfig := createValidTestConfig()
	if err := saveConfigToFile(configFile, initialConfig); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// 创建ConfigManager
	cm, err := NewConfigManager(configFile, &testLoggerImpl{Logger: logrus.New()})
	if err != nil {
		t.Fatalf("Failed to create ConfigManager: %v", err)
	}

	// 监听变更通知
	notifyChan := cm.ChangeNotify()

	// 在goroutine中更新配置
	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = cm.Update(func(c *Config) error {
			c.Port = ":6000"
			return nil
		})
	}()

	// 等待通知
	select {
	case <-notifyChan:
		// 收到通知，验证配置已更新
		config := cm.Get()
		if config.Port != ":6000" {
			t.Errorf("Expected port :6000, got %s", config.Port)
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for change notification")
	}
}

func TestConfigManager_UpdateWithError(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.toml")

	// 创建初始配置
	initialConfig := createValidTestConfig()
	initialConfig.Port = ":3000"
	if err := saveConfigToFile(configFile, initialConfig); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// 创建ConfigManager
	cm, err := NewConfigManager(configFile, &testLoggerImpl{Logger: logrus.New()})
	if err != nil {
		t.Fatalf("Failed to create ConfigManager: %v", err)
	}

	// 测试Update方法返回错误
	expectedErr := os.ErrInvalid
	err = cm.Update(func(c *Config) error {
		return expectedErr
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// 验证配置未被修改
	config := cm.Get()
	if config.Port != ":3000" {
		t.Errorf("Expected port :3000 (unchanged), got %s", config.Port)
	}
}

func TestConfigManager_NoConfigFile(t *testing.T) {
	// 创建不带配置文件的ConfigManager
	cm := &ConfigManager{
		config:       createValidTestConfig(),
		configFile:   "",
		changeNotify: make(chan struct{}, 1),
	}

	// 测试Update方法应该返回错误
	err := cm.Update(func(c *Config) error {
		c.Port = ":8080"
		return nil
	})
	if err == nil {
		t.Error("Expected error when updating without config file")
	}

	// 测试Set方法应该返回错误
	err = cm.Set(createValidTestConfig())
	if err == nil {
		t.Error("Expected error when setting without config file")
	}
}

// saveConfigToFile 辅助函数：保存配置到文件
func saveConfigToFile(filename string, config *Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(config)
}
