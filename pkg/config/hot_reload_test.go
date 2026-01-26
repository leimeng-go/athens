package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConfigManager_HotReload(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test.toml")

	// 写入初始配置
	initialConfig := `
GoEnv = "development"
GoBinary = "go"
GoGetWorkers = 10
ProtocolWorkers = 30
LogLevel = "debug"
CloudRuntime = "none"
Timeout = 300
StorageType = "mongo"
Port = ":3000"
NetworkMode = "strict"

[Storage]
    [Storage.Mongo]
        URL = "mongodb://127.0.0.1:27017"
`
	if err := os.WriteFile(configFile, []byte(initialConfig), 0600); err != nil {
		t.Fatalf("Failed to write initial config: %v", err)
	}

	// 创建配置管理器
	cm, err := NewConfigManager(configFile)
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// 验证初始配置
	conf := cm.Get()
	if conf.LogLevel != "debug" {
		t.Errorf("Expected LogLevel=debug, got %s", conf.LogLevel)
	}
	if conf.GoGetWorkers != 10 {
		t.Errorf("Expected GoGetWorkers=10, got %d", conf.GoGetWorkers)
	}

	// 注册变更回调
	changed := make(chan *Config, 1)
	cm.OnChange(func(newConf *Config) {
		changed <- newConf
	})

	// 修改配置文件
	time.Sleep(100 * time.Millisecond) // 等待文件监听器就绪
	updatedConfig := `
GoEnv = "development"
GoBinary = "go"
GoGetWorkers = 20
ProtocolWorkers = 30
LogLevel = "info"
CloudRuntime = "none"
Timeout = 300
StorageType = "mongo"
Port = ":3000"
NetworkMode = "strict"

[Storage]
    [Storage.Mongo]
        URL = "mongodb://127.0.0.1:27017"
`
	if err := os.WriteFile(configFile, []byte(updatedConfig), 0600); err != nil {
		t.Fatalf("Failed to write updated config: %v", err)
	}

	// 等待配置变更通知
	select {
	case newConf := <-changed:
		if newConf.LogLevel != "info" {
			t.Errorf("Expected LogLevel=info after reload, got %s", newConf.LogLevel)
		}
		if newConf.GoGetWorkers != 20 {
			t.Errorf("Expected GoGetWorkers=20 after reload, got %d", newConf.GoGetWorkers)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for config reload")
	}

	// 验证 Get() 返回最新配置
	latestConf := cm.Get()
	if latestConf.LogLevel != "info" {
		t.Errorf("Expected Get() to return LogLevel=info, got %s", latestConf.LogLevel)
	}
}

func TestConfigManager_NoConfigFile(t *testing.T) {
	// 测试无配置文件时的行为
	cm, err := NewConfigManager("")
	if err != nil {
		t.Fatalf("Failed to create config manager with no config file: %v", err)
	}

	conf := cm.Get()
	if conf == nil {
		t.Fatal("Expected default config, got nil")
	}

	// 验证返回的是默认配置
	if conf.Port != ":3000" {
		t.Errorf("Expected default port :3000, got %s", conf.Port)
	}
}

func TestConfigManager_ThreadSafety(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test.toml")

	initialConfig := `
GoEnv = "development"
GoBinary = "go"
GoGetWorkers = 10
ProtocolWorkers = 30
LogLevel = "debug"
CloudRuntime = "none"
Timeout = 300
StorageType = "mongo"
NetworkMode = "strict"

[Storage]
    [Storage.Mongo]
        URL = "mongodb://127.0.0.1:27017"
`
	if err := os.WriteFile(configFile, []byte(initialConfig), 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	cm, err := NewConfigManager(configFile)
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// 并发读取配置
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				conf := cm.Get()
				if conf == nil {
					t.Error("Got nil config")
				}
			}
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}
}
