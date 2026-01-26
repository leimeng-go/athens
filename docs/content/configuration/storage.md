---
title: Configuring Storage
description: Configuring Storage in Athens
weight: 3
---

## Storage

Athens proxy now **only supports MongoDB** as storage backend. All other storage types have been removed to simplify the codebase.

## Table of Contents

- [Why MongoDB Only?](#why-mongodb-only)
- [Configuration](#configuration)
- [Connection String Examples](#connection-string-examples)
- [Running Multiple Athens Instances](#running-multiple-athens-instances)
- [Single Flight Mechanisms](#single-flight-mechanisms)
  - [Memory](#memory-default)
  - [Etcd](#etcd)
  - [Redis](#redis)
  - [Redis Sentinel](#redis-sentinel)

---

## Why MongoDB Only?

MongoDB is chosen as the only supported storage backend because:

- ✅ **Data Comparison**: Easy to compare and sync modules between external and internal networks
- ✅ **Flexible Queries**: Advanced querying capabilities for module management  
- ✅ **Scalability**: Supports replica sets for high availability
- ✅ **GridFS**: Efficient storage for large module files
- ✅ **Enterprise Ready**: Production-proven for mission-critical applications

---

## Configuration

On startup, Athens will create an `athens` database and `modules` collection on your MongoDB server.

```toml
# StorageType sets the type of storage backend
# Only "mongo" is supported
# Env override: ATHENS_STORAGE_TYPE
StorageType = "mongo"

[Storage]
    [Storage.Mongo]
        # Full URL for mongo storage
        # Env override: ATHENS_MONGO_STORAGE_URL
        URL = "mongodb://127.0.0.1:27017"

        # Optional: Path to certificate to use for the mongo connection
        # Env override: ATHENS_MONGO_CERT_PATH
        CertPath = "/path/to/cert/file"

        # Optional: Allows for insecure SSL / http connections
        # Should be used for testing or development only
        # Env override: ATHENS_MONGO_INSECURE
        Insecure = false

        # Optional: Custom database name (default: "athens")
        # Env override: ATHENS_MONGO_DEFAULT_DATABASE
        DefaultDBName = "athens"

        # Optional: Custom collection name (default: "modules")
        # Env override: ATHENS_MONGO_DEFAULT_COLLECTION
        DefaultCollectionName = "modules"
```

---

## Connection String Examples

**Single Instance:**
```
mongodb://admin:password@localhost:27017/?authSource=admin
```

**Replica Set:**
```
mongodb://admin:password@mongo-rs-1:27017,mongo-rs-2:27017,mongo-rs-3:27017/?authSource=admin&replicaSet=rs0
```

---

## Running Multiple Athens Instances

When running multiple Athens instances pointed at the same MongoDB storage, you need to configure a distributed lock mechanism (**single flight**) to prevent concurrent writes to the same module.

---

## Single Flight Mechanisms

### Memory (Default)

In-memory locking, only suitable for single instance deployments.

```toml
SingleFlightType = "memory"
```

---

### Etcd

Etcd provides distributed locking for multiple Athens servers. When multiple Athens instances receive requests for the same module simultaneously, etcd ensures only one instance downloads and stores the module.

```toml
SingleFlightType = "etcd"

[SingleFlight]
    [SingleFlight.Etcd]
        # Comma separated URLs of etcd servers
        # Env override: ATHENS_ETCD_ENDPOINTS
        Endpoints = "localhost:2379,localhost:22379,localhost:32379"
```

---

### Redis

Athens supports two mechanisms: direct connection or connecting via redis sentinels.

#### Direct Connection

Using a direct connection to redis only requires a single `redis-server`:

```toml
SingleFlightType = "redis"

[SingleFlight]
    [SingleFlight.Redis]
        # Endpoint is the redis endpoint
        # Env override: ATHENS_REDIS_ENDPOINT
        Endpoint = "127.0.0.1:6379"

        # Optional: Password for authentication
        # Env override: ATHENS_REDIS_PASSWORD
        Password = ""
```

**Redis URL format is also supported:**

```toml
SingleFlightType = "redis"

[SingleFlight]
    [SingleFlight.Redis]
        # Use rediss:// for TLS
        Endpoint = "redis://user:password@127.0.0.1:6379/0?protocol=3"
```

**Customizing Lock Configuration:**

If you need to customize the distributed lock behavior:

```toml
[SingleFlight.Redis]
    ...
    [SingleFlight.Redis.LockConfig]
        # TTL for the lock in seconds (default: 900 = 15 minutes)
        # Env override: ATHENS_REDIS_LOCK_TTL
        TTL = 900
        # Timeout for acquiring the lock in seconds (default: 15)
        # Env override: ATHENS_REDIS_LOCK_TIMEOUT
        Timeout = 15
        # Max retries while acquiring the lock (default: 10)
        # Env override: ATHENS_REDIS_LOCK_MAX_RETRIES
        MaxRetries = 10
```

---

### Redis Sentinel

**NOTE**: Redis Sentinel requires a working knowledge of Redis and is not recommended for beginners.

Redis Sentinel provides high-availability for Redis through automated monitoring, replication, and failover.

For more details, see the [Redis Sentinel documentation](https://redis.io/topics/sentinel).

**Required Configuration:**

- `Endpoints`: List of redis-sentinel endpoints (typically 3 or more)
- `MasterName`: Named master instance from sentinel configuration

```toml
SingleFlightType = "redis-sentinel"

[SingleFlight]
  [SingleFlight.RedisSentinel]
      # Redis sentinel endpoints
      # Env override: ATHENS_REDIS_SENTINEL_ENDPOINTS
      Endpoints = ["127.0.0.1:26379"]
      
      # Master name from sentinel config
      MasterName = "redis-1"
      
      # Optional: Sentinel authentication password
      SentinelPassword = "sekret"
      
      # Optional: Redis master credentials
      # Env override: ATHENS_REDIS_USERNAME
      RedisUsername = ""
      # Env override: ATHENS_REDIS_PASSWORD
      RedisPassword = ""
```

Lock configuration can be customized for Redis Sentinel in the same way as regular Redis.
