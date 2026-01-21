# Google Cloud Storage配置

<cite>
**本文档引用的文件**
- [storage.go](file://pkg/config/storage.go)
- [config.go](file://pkg/config/config.go)
- [singleflight.go](file://pkg/config/singleflight.go)
- [storage.go](file://cmd/proxy/actions/storage.go)
- [storage.md](file://docs/content/configuration/storage.md)
</cite>

## 更新摘要
**所做更改**
- 移除了所有关于Google Cloud Storage配置的内容
- 更新了架构说明以反映GCS功能已被完全移除
- 删除了GCS相关的配置示例和部署指导
- 更新了支持的存储类型列表，仅保留MongoDB选项

## 目录
1. [简介](#简介)
2. [当前支持的存储类型](#当前支持的存储类型)
3. [MongoDB存储配置](#mongodb存储配置)
4. [GCS功能状态说明](#gcs功能状态说明)
5. [迁移指南](#迁移指南)
6. [结论](#结论)

## 简介
本文档原本旨在提供Google Cloud Storage配置的详细说明，但根据最新的代码变更，Google Cloud Storage配置功能已被完全移除，不再支持GCS存储选项。本文档现更新为反映这一变更，并提供相应的迁移指导。

## 当前支持的存储类型
根据最新的代码实现，Athens当前仅支持以下存储类型：

- **内存存储（Memory）**：用于开发目的的临时存储
- **磁盘存储（Disk）**：基于文件系统的持久化存储
- **MongoDB存储（Mongo）**：生产环境推荐的数据库存储方案

**章节来源**
- [storage.go](file://pkg/config/storage.go#L3-L6)
- [config.go](file://pkg/config/config.go#L399-L404)

## MongoDB存储配置
MongoDB是当前唯一受支持的生产存储解决方案，配置方法如下：

### 基础配置
```toml
# 设置存储类型为MongoDB
StorageType = "mongo"

[Storage]
    [Storage.Mongo]
        # MongoDB连接URL
        URL = "mongodb://127.0.0.1:27017"
        
        # 可选：证书路径
        CertPath = "/path/to/cert/file"
        
        # 可选：允许不安全SSL连接（仅测试环境）
        Insecure = false
        
        # 可选：自定义数据库名称
        DefaultDBName = "athens"
        
        # 可选：自定义集合名称
        DefaultCollectionName = "modules"
```

### 环境变量配置
- `ATHENS_STORAGE_TYPE=mongo`
- `ATHENS_MONGO_STORAGE_URL`：MongoDB连接URL
- `ATHENS_MONGO_CERT_PATH`：证书路径
- `ATHENS_MONGO_INSECURE`：是否允许不安全连接
- `ATHENS_MONGO_DEFAULT_DATABASE`：自定义数据库名称
- `ATHENS_MONGO_DEFAULT_COLLECTION`：自定义集合名称

**章节来源**
- [storage.md](file://docs/content/configuration/storage.md#L71-L107)

## GCS功能状态说明
### 功能状态
Google Cloud Storage配置功能已被完全移除，不再支持GCS存储选项。

### 代码变更分析
1. **配置结构变更**：`pkg/config/storage.go`中的`Storage`结构体仅包含`Mongo *MongoConfig`
2. **存储选择逻辑**：`cmd/proxy/actions/storage.go`中的`GetStorage`函数仅支持"mongo"类型
3. **验证逻辑**：`pkg/config/config.go`中的`validateStorage`函数强制要求存储类型为"mongo"
4. **文档更新**：`docs/content/configuration/storage.md`中移除了Google Cloud Storage相关章节

### 影响范围
- 移除了所有GCS相关的配置参数（ProjectID、Bucket、ServiceAccountJSON、Endpoint）
- 移除了GCS认证方式说明（服务账号、应用默认凭据、工作负载身份）
- 移除了GCS存储的ACL配置说明
- 移除了GCS部署配置示例
- 移除了GCS成本优化和性能特征说明

**章节来源**
- [storage.go](file://pkg/config/storage.go#L3-L6)
- [storage.go](file://cmd/proxy/actions/storage.go#L14-L19)
- [config.go](file://pkg/config/config.go#L299-L304)

## 迁移指南
如果您的系统之前使用Google Cloud Storage，建议按以下步骤迁移到MongoDB：

### 1. 数据迁移
```bash
# 导出现有GCS数据
# 使用GCS客户端工具导出所有模块文件

# 导入到MongoDB
# 使用MongoDB导入工具将数据转换为MongoDB格式
```

### 2. 配置迁移
将原有的GCS配置转换为MongoDB配置：

**从GCS配置**：
```toml
StorageType = "gcp"
[Storage]
    [Storage.GCP]
        ProjectID = "your-project-id"
        Bucket = "your-bucket-name"
```

**转换为MongoDB配置**：
```toml
StorageType = "mongo"
[Storage]
    [Storage.Mongo]
        URL = "mongodb://localhost:27017"
```

### 3. 环境变量迁移
- 移除GCS相关环境变量：`GOOGLE_CLOUD_PROJECT`、`ATHENS_STORAGE_GCP_BUCKET`、`ATHENS_STORAGE_GCP_JSON_KEY`
- 添加MongoDB相关环境变量：`ATHENS_MONGO_STORAGE_URL`

### 4. 部署更新
```yaml
# Kubernetes部署示例
apiVersion: apps/v1
kind: Deployment
metadata:
  name: athens-proxy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: athens-proxy
  template:
    metadata:
      labels:
        app: athens-proxy
    spec:
      containers:
      - name: athens-proxy
        image: ghcr.io/leimeng-go/athens:latest
        env:
        - name: ATHENS_STORAGE_TYPE
          value: "mongo"
        - name: ATHENS_MONGO_STORAGE_URL
          value: "mongodb://mongo-service:27017"
```

## 结论
Google Cloud Storage配置功能已被完全移除，这是Athens项目的一个重要变更。当前系统仅支持MongoDB存储方案，用户需要：

1. **立即行动**：如果使用GCS，需要尽快迁移到MongoDB
2. **数据迁移**：制定详细的迁移计划，确保数据完整性和业务连续性
3. **配置更新**：更新所有环境配置和部署文件
4. **测试验证**：在生产环境部署前进行全面的功能和性能测试

虽然GCS功能已移除，但MongoDB提供了更好的生产就绪特性和企业级功能，包括更好的性能、可扩展性和管理工具支持。建议新用户直接采用MongoDB存储方案，以获得最佳的使用体验。