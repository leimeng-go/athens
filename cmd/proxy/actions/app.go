package actions

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gomods/athens/pkg/config"
	"github.com/gomods/athens/pkg/log"
	mw "github.com/gomods/athens/pkg/middleware"
	"github.com/gomods/athens/pkg/module"
	"github.com/gomods/athens/pkg/observ"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
	"go.opencensus.io/plugin/ochttp"
)

// Service is the name of the service that we want to tag our processes with.
const Service = "proxy"

// App is where all routes and middleware for the proxy
// should be defined. This is the nerve center of your
// application.
func App(conf *config.Config) (http.Handler, error) {
	//环境配置，开发和生产
	// ENV is used to help switch settings based on where the
	// application is being run. Default is "development".
	ENV := conf.GoEnv

	if conf.GithubToken != "" {
		if conf.NETRCPath != "" {
			fmt.Println("Cannot provide both GithubToken and NETRCPath. Only provide one.")
			os.Exit(1)
		}

		netrcFromToken(conf.GithubToken)
	}
    // 挂载 .netrc git相关配置
	// mount .netrc to home dir
	// to have access to private repos.
	initializeAuthFile(conf.NETRCPath)
    
	//用于指定Mercurial（Hg）的配置文件路径。Mercurial是一种分布式版本控制系统，用于管理代码仓库
	// mount .hgrc to home dir
	// to have access to private repos.
	initializeAuthFile(conf.HGRCPath)

	//配置日志级别
	logLvl, err := logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		return nil, err
	}
	//构建日志实例,CloudRuntime用于配置日志格式                    
	lggr := log.New(conf.CloudRuntime, logLvl)
    //构建路由组
	r := mux.NewRouter()
	r.Use(
		//请求ID中间件
		mw.WithRequestID,
		//日志中间件，拼接一些字段
		mw.LogEntryMiddleware(lggr),
		mw.RequestLogger,
		secure.New(secure.Options{
			SSLRedirect:     conf.ForceSSL,
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}).Handler,
		mw.ContentType,
	    
	)
    //某些负载均衡器的单个检查路由？
	var subRouter *mux.Router
	if prefix := conf.PathPrefix; prefix != "" {
		// certain Ingress Controllers (such as GCP Load Balancer)
		// can not send custom headers and therefore if the proxy
		// is running behind a prefix as well as some authentication
		// mechanism, we should allow the plain / to return 200.
		r.HandleFunc("/", healthHandler).Methods(http.MethodGet)
		subRouter = r.PathPrefix(prefix).Subrouter()
	}
    //相关trance跟踪，例如jaeger
	// RegisterExporter will register an exporter where we will export our traces to.
	// The error from the RegisterExporter would be nil if the tracer was specified by
	// the user and the trace exporter was created successfully.
	// RegisterExporter returns the function that all traces are flushed to the exporter
	// and the exporter needs to be stopped. The function should be called when the exporter
	// is no longer needed.
	flushTraces, err := observ.RegisterExporter(
		conf.TraceExporter,
		conf.TraceExporterURL,
		Service,
		ENV,
	)
	if err != nil {
		lggr.Infof("%s", err)
	} else {
		defer flushTraces()
	}
     //状态监控，例如prometheus
	// RegisterStatsExporter will register an exporter where we will collect our stats.
	// The error from the RegisterStatsExporter would be nil if the proper stats exporter
	// was specified by the user.
	flushStats, err := observ.RegisterStatsExporter(r, conf.StatsExporter, Service)
	if err != nil {
		lggr.Infof("%s", err)
	} else {
		defer flushStats()
	}
    //使用BasicAuth账号密码认证
	user, pass, ok := conf.BasicAuth()
	if ok {
		r.Use(basicAuth(user, pass))
	}
    //使用过滤
	if !conf.FilterOff() {
		mf, err := module.NewFilter(conf.FilterFile)
		if err != nil {
			lggr.Fatal(err)
		}
		r.Use(mw.NewFilterMiddleware(mf, conf.GlobalEndpoint))
	}

	client := &http.Client{
		Transport: &ochttp.Transport{
			Base: http.DefaultTransport,
		},
	}
    //Hook是重新调用地址吗？
	// Having the hook set means we want to use it
	if vHook := conf.ValidatorHook; vHook != "" {
		r.Use(mw.NewValidationMiddleware(client, vHook))
	}
    //下载的依赖存放方式
	store, err := GetStorage(conf.StorageType, conf.Storage, conf.TimeoutDuration(), client)
	if err != nil {
		err = fmt.Errorf("error getting storage configuration: %w", err)
		return nil, err
	}
    //subRouter包含原来的
	proxyRouter := r
	if subRouter != nil {
		proxyRouter = subRouter
	}
	//添加相关handler
	if err := addProxyRoutes(
		proxyRouter,
		store,
		lggr,
		conf,
	); err != nil {
		err = fmt.Errorf("error adding proxy routes: %w", err)
		return nil, err
	}

	h := &ochttp.Handler{
		Handler: r,
	}

	return h, nil
}
