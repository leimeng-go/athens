package examples
package main

import (
	"log"
	"time"

	"github.com/leimeng-go/athens/pkg/config"
)

// 演示配置热更新功能
func main() {
	// 创建配置管理器，监听 config.dev.toml
	configManager, err := config.NewConfigManager("../config.dev.toml")
	if err != nil {
		log.Fatalf("Failed to create config manager: %v", err)
	}

	log.Println("=== Athens 配置热更新演示 ===")
	log.Println("配置文件: config.dev.toml")
	log.Println("尝试修改配置文件，观察自动重载...")
	log.Println("")

	// 注册配置变更回调
	configManager.OnChange(func(newConf *config.Config) {
		log.Printf("✅ 配置已更新:")
		log.Printf("  - LogLevel: %s", newConf.LogLevel)
		log.Printf("  - StorageType: %s", newConf.StorageType)
		log.Printf("  - Port: %s", newConf.Port)
		log.Printf("  - GoGetWorkers: %d", newConf.GoGetWorkers)
		log.Println("")
	})

	// 定期输出当前配置
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("📋 当前配置:")
	printConfig(configManager.Get())

	for {
		select {
		case <-ticker.C:
			log.Println("📋 当前配置:")
			printConfig(configManager.Get())

		case <-configManager.ChangeNotify():
			// 配置变更通知已经在回调中处理
		}
	}
}

func printConfig(conf *config.Config) {
	log.Printf("  - LogLevel: %s", conf.LogLevel)
	log.Printf("  - StorageType: %s", conf.StorageType)
	log.Printf("  - Port: %s", conf.Port)
	log.Printf("  - GoGetWorkers: %d", conf.GoGetWorkers)
	log.Println("")
}
