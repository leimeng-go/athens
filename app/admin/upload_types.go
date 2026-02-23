package admin

import "time"

// UploadTask 表示模块上传任务
type UploadTask struct {
	ID          string    `json:"id"`
	ModulePath  string    `json:"modulePath"`
	Version     string    `json:"version"`
	Source      string    `json:"source"` // "file" 或 "url"
	Status      string    `json:"status"` // "pending", "processing", "completed", "failed"
	Progress    int       `json:"progress"` // 0-100
	Error       string    `json:"error,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
	FileSize    int64     `json:"fileSize,omitempty"`
}