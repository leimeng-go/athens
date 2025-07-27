package admin


// DashboardData 表示仪表盘数据
type DashboardData struct {
	Stats           Stats            `json:"stats"`
	DownloadTrend   []DownloadTrend  `json:"downloadTrend"`
	PopularModules  []PopularModule  `json:"popularModules"`
	RecentActivities []RecentActivity `json:"recentActivities"`
}

// Stats 表示统计数据
type Stats struct {
	TotalModules     int    `json:"totalModules"`
	TotalDownloads   int    `json:"totalDownloads"`
	TotalRepositories int    `json:"totalRepositories"`
	StorageUsed      string `json:"storageUsed"`
}

// DownloadTrend 表示下载趋势数据
type DownloadTrend struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// PopularModule 表示热门模块数据
type PopularModule struct {
	Path      string `json:"path"`
	Downloads int    `json:"downloads"`
}

// RecentActivity 表示最近活动数据
type RecentActivity struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"` // 'download' | 'upload' | 'system'
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}
