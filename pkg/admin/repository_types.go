package admin

import "time"

// RepositoryData 表示仓库信息
type RepositoryData struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Type        string    `json:"type"`
	ModuleCount int       `json:"moduleCount"`
	CreatedAt   time.Time `json:"createdAt"`
	LastSync    time.Time `json:"lastSync"`
	Status      string    `json:"status"`
}