package admin

import "time"

// ModuleData 表示模块信息
type ModuleData struct {
	Path        string    `json:"path"`
	Version     string    `json:"version"`
	Size        int64     `json:"size"`
	Downloads   int       `json:"downloads"`
	LastAccess  time.Time `json:"lastAccess"`
	Description string    `json:"description"`
}

// DownloadStats 表示下载统计信息
type DownloadStats struct {
	TotalDownloads int                  `json:"totalDownloads"`
	TotalModules   int                  `json:"totalModules"`
	PopularModules []PopularModule      `json:"popularModules"`
	DownloadTrends []DownloadTrend      `json:"downloadTrends"`
	RecentModules  []RecentModuleAccess `json:"recentModules"`
}

// RecentModuleAccess 表示最近访问的模块
type RecentModuleAccess struct {
	Path       string    `json:"path"`
	Version    string    `json:"version"`
	AccessTime time.Time `json:"accessTime"`
}