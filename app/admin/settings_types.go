package admin

// SystemSettings 表示系统设置信息
type SystemSettings struct {
	StoragePath          string `json:"storagePath"`
	MaxUploadSize        int    `json:"maxUploadSize"`
	EnablePrivateModules bool   `json:"enablePrivateModules"`
	EnableDownloadLogging bool  `json:"enableDownloadLogging"`
	ProxyTimeout         int    `json:"proxyTimeout"`
	CacheExpiration      int    `json:"cacheExpiration"`
}