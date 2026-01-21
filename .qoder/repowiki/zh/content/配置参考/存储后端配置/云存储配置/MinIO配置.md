# MinIO配置

<cite>
**本文档引用的文件**
- [pkg/config/storage.go](file://pkg/config/storage.go)
- [pkg/config/config.go](file://pkg/config/config.go)
- [pkg/config/config_test.go](file://pkg/config/config_test.go)
- [config.dev.toml](file://config.dev.toml)
- [docs/content/configuration/storage.md](file://docs/content/configuration/storage.md)
</cite>

## 更新摘要
**所做更改**
- 移除了所有关于MinIO存储配置的内容，因为该功能已被完全移除
- 更新了配置验证逻辑，现在只支持MongoDB存储类型
- 移除了MinIO相关的配置示例和文档
- 更新了架构图以反映当前的存储配置结构

## 目录
1. [简介](#简介)
2. [当前存储配置状态](#当前存储配置状态)
3. [支持的存储类型](#支持的存储类型)
4. [配置验证机制](#配置验证机制)
5. [迁移指南](#迁移指南)
6. [结论](#结论)

## 简介
本文档原本旨在说明MinIO存储配置，但根据最新的代码变更，MinIO存储配置功能已被完全移除。当前版本的Athens只支持MongoDB作为存储后端。本文档将更新为反映这一变更，并提供相应的迁移指导。

## 当前存储配置状态
**重要变更**：MinIO存储配置功能已被完全移除，不再支持MinIO存储选项。

### 已移除的功能
- MinioConfig结构体定义
- MinIO存储后端实现
- MinIO相关的配置验证逻辑
- MinIO配置示例和文档

### 当前支持的存储类型
根据最新的配置验证逻辑，当前版本只支持以下存储类型：

```mermaid
graph TB
Config["Config<br/>pkg/config/config.go"] --> StorageType["StorageType<br/>必须为'mongo'"]
StorageType --> Storage["Storage<br/>pkg/config/storage.go"]
Storage --> MongoConfig["MongoConfig<br/>唯一支持的存储配置"]
```

**图表来源**
- [pkg/config/config.go](file://pkg/config/config.go#L299-L304)
- [pkg/config/storage.go](file://pkg/config/storage.go#L3-L6)

**章节来源**
- [pkg/config/config.go](file://pkg/config/config.go#L299-L304)
- [pkg/config/storage.go](file://pkg/config/storage.go#L3-L6)

## 支持的存储类型
### MongoDB存储（唯一支持的存储类型）

当前版本的Athens只支持MongoDB作为存储后端，配置结构非常简洁：

```mermaid
classDiagram
class Storage {
+Mongo *MongoConfig
}
class MongoConfig {
+string URL
+string CertPath
+bool Insecure
+string DefaultDBName
+string DefaultCollectionName
}
```

**图表来源**
- [pkg/config/storage.go](file://pkg/config/storage.go#L3-L6)

**章节来源**
- [pkg/config/storage.go](file://pkg/config/storage.go#L3-L6)
- [config.dev.toml](file://config.dev.toml#L454-L471)

## 配置验证机制
### 存储类型验证

当前的配置验证逻辑严格限制存储类型为MongoDB：

```mermaid
flowchart TD
Start["开始配置验证"] --> CheckType{"StorageType是否为'mongo'?"}
CheckType --> |是| ValidateMongo["验证MongoConfig"]
CheckType --> |否| Error["返回错误: 仅支持'mongo'存储类型"]
ValidateMongo --> Success["验证成功"]
Error --> End["结束"]
Success --> End
```

**图表来源**
- [pkg/config/config.go](file://pkg/config/config.go#L299-L304)

**章节来源**
- [pkg/config/config.go](file://pkg/config/config.go#L299-L304)

### 配置示例

完整的MongoDB存储配置示例：

```toml
[Storage]
    [Storage.Mongo]
        URL = "mongodb://admin:password123@mongo-rs-1:27017,mongo-rs-2:27017,mongo-rs-3:27017/?authSource=admin&replicaSet=rs0"
        DefaultDBName = "athens"
        CertPath = ""
        Insecure = false
```

**章节来源**
- [config.dev.toml](file://config.dev.toml#L454-L471)

## 迁移指南
### 从MinIO迁移到MongoDB

由于MinIO功能已被移除，用户需要进行以下迁移步骤：

1. **停止使用MinIO配置**
   - 移除所有MinIO相关的配置项
   - 删除MinIO存储桶和数据

2. **设置MongoDB存储**
   - 配置MongoDB连接字符串
   - 设置数据库和集合名称
   - 配置SSL证书路径（如需要）

3. **数据迁移**
   - 将MinIO中的模块数据迁移到MongoDB
   - 验证数据完整性和一致性

4. **配置更新**
   - 更新StorageType为"mongo"
   - 移除MinIO相关的环境变量
   - 添加MongoDB配置参数

**章节来源**
- [config.dev.toml](file://config.dev.toml#L122-L125)
- [config.dev.toml](file://config.dev.toml#L454-L471)

## 结论
MinIO存储配置功能已被完全移除，当前版本的Athens只支持MongoDB存储后端。用户需要按照迁移指南将现有的MinIO配置迁移到MongoDB。新的配置结构更加简洁，只包含必要的MongoDB配置参数，简化了部署和维护工作。

对于需要对象存储功能的用户，建议考虑以下替代方案：
- 使用MongoDB的GridFS功能
- 部署外部对象存储服务并通过自定义存储后端集成
- 使用其他支持的存储后端（如Azure Blob Storage、Google Cloud Storage等）