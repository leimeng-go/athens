package actions

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/leimeng-go/athens/pkg/config"
	"github.com/leimeng-go/athens/pkg/download"
	"github.com/leimeng-go/athens/pkg/download/addons"
	"github.com/leimeng-go/athens/pkg/download/mode"
	"github.com/leimeng-go/athens/pkg/index"
	"github.com/leimeng-go/athens/pkg/index/mem"
	"github.com/leimeng-go/athens/pkg/index/mysql"
	"github.com/leimeng-go/athens/pkg/index/nop"
	"github.com/leimeng-go/athens/pkg/index/postgres"
	"github.com/leimeng-go/athens/pkg/log"
	"github.com/leimeng-go/athens/pkg/module"
	"github.com/leimeng-go/athens/pkg/stash"
	"github.com/leimeng-go/athens/pkg/storage"
	"github.com/spf13/afero"
)

func addProxyRoutes(
	r *mux.Router,
	s storage.Backend,
	l *log.Logger,
	c *config.Config,
) error {
	r.HandleFunc("/", proxyHomeHandler(c))
	r.HandleFunc("/healthz", healthHandler)
	r.HandleFunc("/readyz", getReadinessHandler(s))
	r.HandleFunc("/version", versionHandler)
	r.HandleFunc("/catalog", catalogHandler(s))
	r.HandleFunc("/robots.txt", robotsHandler(c))
	//Go模块的索引
	indexer, err := getIndex(c)
	if err != nil {
		return err
	}
	r.HandleFunc("/index", indexHandler(indexer))

	for _, sumdb := range c.SumDBs {
		sumdbURL, err := url.Parse(sumdb)
		if err != nil {
			return err
		}
		if sumdbURL.Scheme != "https" {
			return fmt.Errorf("sumdb: %v must have an https scheme", sumdb)
		}
		supportPath := path.Join("/sumdb", sumdbURL.Host, "/supported")
		r.HandleFunc(supportPath, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		sumHandler := sumdbProxy(sumdbURL, c.NoSumPatterns)
		pathPrefix := "/sumdb/" + sumdbURL.Host
		r.PathPrefix(pathPrefix + "/").Handler(
			http.StripPrefix(strings.TrimSuffix(c.PathPrefix, "/")+pathPrefix, sumHandler),
		)
	}

	// Download Protocol:
	// 下载协议：
	// the download.Protocol and the stash.Stasher interfaces are composable
	// download.Protocol 和 stash.Stasher 接口是可以组合的（类似中间件）
	// in a middleware fashion. Therefore you can separate concerns
	// 以中间件的方式组合，因此可以分离不同的功能关注点
	// by the functionality: a download.Protocol that just takes care
	// 例如：download.Protocol 只负责 go get 相关的功能
	// of "go getting" things, and another Protocol that just takes care
	// 另一个 Protocol 只负责请求池（限流）等功能
	// of "pooling" requests etc.
	//
	// In our case, we'd like to compose both interfaces in a particular
	// 在本项目中，我们希望以特定顺序组合这两个接口
	// order to ensure logical ordering of execution.
	// 以保证执行流程的逻辑性
	//
	// Here's the order of an incoming request to the download.Protocol:
	// 以下是请求到达 download.Protocol 时的处理顺序：
	//
	// 1. The downloadpool gets hit first, and manages concurrent requests
	// 1. 首先由 downloadpool（下载池）处理，负责并发请求的管理
	// 2. The downloadpool passes the request to its parent Protocol: stasher
	// 2. downloadpool 将请求传递给其父 Protocol：stasher
	// 3. The stasher Protocol checks storage first, and if storage is empty
	// 3. stasher Protocol 首先检查本地存储，如果没有数据
	// it makes a Stash request to the stash.Stasher interface.
	// 则向 stash.Stasher 接口发起 Stash 请求
	//
	// Once the stasher picks up an order, here's how the requests go in order:
	// stasher 接收到请求后，内部的处理顺序如下：
	// 1. The singleflight picks up the first request and latches duplicate ones.
	// 1. singleflight 处理第一个请求，并锁定重复的请求（防止重复下载）
	// 2. The singleflight passes the stash to its parent: stashpool.
	// 2. singleflight 将 stash 请求传递给其父级：stashpool
	// 3. The stashpool manages limiting concurrent requests and passes them to stash.
	// 3. stashpool 管理并发数限制，并将请求传递给 stash
	// 4. The plain stash.New just takes a request from upstream and saves it into storage.
	// 4. 最终由 stash.New 执行实际的下载并保存到存储
	//申明一个文件系统相关接口
	fs := afero.NewOsFs()
	//go sumdb校验
	if !c.GoBinaryEnvVars.HasKey("GONOSUMDB") {
		c.GoBinaryEnvVars.Add("GONOSUMDB", strings.Join(c.NoSumPatterns, ","))
	}
	//校验环境变量
	if err := c.GoBinaryEnvVars.Validate(); err != nil {
		return err
	}
	// go get请求组件
	mf, err := module.NewGoGetFetcher(c.GoBinary, c.GoGetDir, c.GoBinaryEnvVars, fs)
	if err != nil {
		return err
	}

	lister := module.NewVCSLister(c.GoBinary, c.GoBinaryEnvVars, fs, c.TimeoutDuration())
	checker := storage.WithChecker(s)
	withSingleFlight, err := getSingleFlight(l, c, s, checker)
	if err != nil {
		return err
	}
	st := stash.New(mf, s, indexer, stash.WithPool(c.GoGetWorkers), withSingleFlight)

	df, err := mode.NewFile(c.DownloadMode, c.DownloadURL)
	if err != nil {
		return err
	}

	dpOpts := &download.Opts{
		Storage:      s,
		Stasher:      st,
		Lister:       lister,
		DownloadFile: df,
		NetworkMode:  c.NetworkMode,
	}

	dp := download.New(dpOpts, addons.WithPool(c.ProtocolWorkers))

	handlerOpts := &download.HandlerOpts{Protocol: dp, Logger: l, DownloadFile: df}
	//注册下载依赖相关路由
	download.RegisterHandlers(r, handlerOpts)

	return nil
}

// athensLoggerForRedis implements pkg/stash.RedisLogger.
type athensLoggerForRedis struct {
	logger *log.Logger
}

func (l *athensLoggerForRedis) Printf(ctx context.Context, format string, v ...any) {
	l.logger.WithContext(ctx).Printf(format, v...)
}

func getSingleFlight(l *log.Logger, c *config.Config, s storage.Backend, checker storage.Checker) (stash.Wrapper, error) {
	switch c.SingleFlightType {
	case "", "memory":
		return stash.WithSingleflight, nil
	case "etcd":
		if c.SingleFlight == nil || c.SingleFlight.Etcd == nil {
			return nil, errors.New("etcd config must be present")
		}
		endpoints := strings.Split(c.SingleFlight.Etcd.Endpoints, ",")
		return stash.WithEtcd(endpoints, checker)
	case "redis":
		if c.SingleFlight == nil || c.SingleFlight.Redis == nil {
			return nil, errors.New("redis config must be present")
		}
		//redis lock
		return stash.WithRedisLock(
			//日志
			&athensLoggerForRedis{logger: l},
			//redis 地址
			c.SingleFlight.Redis.Endpoint,
			//redis密码
			c.SingleFlight.Redis.Password,
			//检查module version
			checker,
			//redis lock相关配置
			c.SingleFlight.Redis.LockConfig)
	case "redis-sentinel":
		if c.SingleFlight == nil || c.SingleFlight.RedisSentinel == nil {
			return nil, errors.New("redis config must be present")
		}
		return stash.WithRedisSentinelLock(
			&athensLoggerForRedis{logger: l},
			c.SingleFlight.RedisSentinel.Endpoints,
			c.SingleFlight.RedisSentinel.MasterName,
			c.SingleFlight.RedisSentinel.SentinelPassword,
			c.SingleFlight.RedisSentinel.RedisUsername,
			c.SingleFlight.RedisSentinel.RedisPassword,
			checker,
			c.SingleFlight.RedisSentinel.LockConfig,
		)
	case "gcp":
		if c.StorageType != "gcp" {
			return nil, fmt.Errorf("gcp SingleFlight only works with a gcp storage type and not: %v", c.StorageType)
		}
		return stash.WithGCSLock(c.SingleFlight.GCP.StaleThreshold, s)
	case "azureblob":
		if c.StorageType != "azureblob" {
			return nil, fmt.Errorf("azureblob SingleFlight only works with a azureblob storage type and not: %v", c.StorageType)
		}
		return stash.WithAzureBlobLock(c.Storage.AzureBlob, c.TimeoutDuration(), checker)
	default:
		return nil, fmt.Errorf("unrecognized single flight type: %v", c.SingleFlightType)
	}
}

func getIndex(c *config.Config) (index.Indexer, error) {
	switch c.IndexType {
	case "", "none":
		return nop.New(), nil
	case "memory":
		return mem.New(), nil
	case "mysql":
		return mysql.New(c.Index.MySQL)
	case "postgres":
		return postgres.New(c.Index.Postgres)
	}
	return nil, fmt.Errorf("unknown index type: %q", c.IndexType)
}
