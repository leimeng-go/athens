# Athens 配置热更新功能

## 功能概述

Athens 现在支持配置文件热更新，无需重启服务即可应用配置变更。基于 Viper 实现自动监听配置文件变化。

## 实现原理

1. **ConfigManager**: 配置管理器，封装了配置加载和热更新逻辑
2. **文件监听**: 使用 Viper 的 WatchConfig 功能监听配置文件变化
3. **线程安全**: 使用 RWMutex 保证并发读写安全
4. **回调机制**: 支持注册配置变更回调函数

## 核心 API

### 创建配置管理器

```go
configManager, err := config.NewConfigManager(configFile)
if err != nil {
    log.Fatalf("Failed to create config manager: %v", err)
}
```

### 获取当前配置

```go
// 线程安全地获取配置
conf := configManager.Get()
```

### 注册配置变更回调

```go
configManager.OnChange(func(newConf *config.Config) {
    // 配置变更时执行的逻辑
    log.Printf("Config updated, new LogLevel: %s", newConf.LogLevel)
    
    // 示例：动态更新日志级别
    if newLogLvl, err := logrus.ParseLevel(newConf.LogLevel); err == nil {
        logger.SetLevel(newLogLvl)
    }
})
```

### 监听配置变更通知

```go
// 可选：监听配置变更通知通道
for {
    select {
    case <-configManager.ChangeNotify():
        log.Println("Config changed!")
    }
}
```

## main.go 中的使用

```go
func main() {
    // 使用 ConfigManager 实现配置热更新
    configManager, err := config.NewConfigManager(*configFile)
    if err != nil {
        stdlog.Fatalf("Could not create config manager: %v", err)
    }

    // 获取初始配置
    conf := configManager.Get()

    // 初始化 logger
    logLvl, _ := logrus.ParseLevel(conf.LogLevel)
    logger := athenslog.New(conf.CloudRuntime, logLvl, conf.LogFormat)

    // 注册配置变更回调：更新日志级别
    configManager.OnChange(func(newConf *config.Config) {
        if newLogLvl, err := logrus.ParseLevel(newConf.LogLevel); err == nil {
            logger.Logger.SetLevel(newLogLvl)
            logger.Infof("Log level updated to: %s", newConf.LogLevel)
        }
    })

    // ... 其他初始化逻辑
}
```

## 支持热更新的配置项

目前已实现的热更新功能：

- ✅ **LogLevel**: 日志级别（debug, info, warn, error）
- 🔄 **其他配置项**: 可以通过 `Get()` 获取最新配置，但需要业务代码配合实现热更新逻辑

## 测试示例

运行热更新演示程序：

```bash
cd examples
go run hot-reload-demo.go
```

然后修改 `config.dev.toml` 中的配置项（如 LogLevel、GoGetWorkers），观察自动重载效果。

## 注意事项

1. **配置验证**: 配置文件变更后会自动验证，验证失败会保留旧配置并输出错误日志
2. **性能影响**: 文件监听对性能影响微乎其微，监听使用操作系统的 inotify/fsevents 机制
3. **无配置文件**: 如果未指定配置文件（使用默认配置），热更新功能不会启用
4. **线程安全**: `ConfigManager.Get()` 是线程安全的，可以在任何 goroutine 中安全调用

## 扩展热更新功能

如果需要让其他组件支持热更新，按以下步骤操作：

```go
// 1. 在初始化时注册回调
configManager.OnChange(func(newConf *config.Config) {
    // 2. 根据新配置更新组件状态
    myComponent.UpdateConfig(newConf)
})
```

示例：热更新存储配置

```go
configManager.OnChange(func(newConf *config.Config) {
    // 注意：某些配置（如存储类型）可能需要重新初始化连接
    // 这种情况下可能需要重启服务
    if newConf.Storage.Mongo.URL != oldURL {
        log.Warn("MongoDB URL changed, please restart the service")
    }
})
```

## 依赖

- `github.com/spf13/viper`: 配置管理和文件监听
- `github.com/fsnotify/fsnotify`: 文件系统事件通知（Viper 依赖）
