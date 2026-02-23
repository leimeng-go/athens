package admin

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// 模拟的上传任务数据
var mockUploadTasks []UploadTask

// 用于保护上传任务数据的并发访问
var uploadTasksMutex sync.RWMutex

// 初始化模拟数据
func init() {
	// 生成模拟上传任务数据
	mockUploadTasks = generateMockUploadTasks()

	// 启动模拟任务处理器
	go simulateTaskProcessing()
}

// generateMockUploadTasks 生成模拟上传任务数据
func generateMockUploadTasks() []UploadTask {
	tasks := make([]UploadTask, 0, 10)

	// 添加一些已完成的任务
	for i := 0; i < 5; i++ {
		modulePath := popularModulePaths[rand.Intn(len(popularModulePaths))]
		version := moduleVersions[rand.Intn(len(moduleVersions))]
		source := "file"
		if rand.Intn(2) == 1 {
			source = "url"
		}

		createdAt := time.Now().Add(-time.Duration(rand.Intn(72)) * time.Hour)
		completedAt := createdAt.Add(time.Duration(rand.Intn(30)) * time.Minute)

		tasks = append(tasks, UploadTask{
			ID:          uuid.New().String(),
			ModulePath:  modulePath,
			Version:     version,
			Source:      source,
			Status:      "completed",
			Progress:    100,
			CreatedAt:   createdAt,
			CompletedAt: completedAt,
			FileSize:    int64(10000 + rand.Intn(990000)), // 10KB-1MB大小
		})
	}

	// 添加一些失败的任务
	for i := 0; i < 2; i++ {
		modulePath := popularModulePaths[rand.Intn(len(popularModulePaths))]
		version := moduleVersions[rand.Intn(len(moduleVersions))]
		source := "file"
		if rand.Intn(2) == 1 {
			source = "url"
		}

		createdAt := time.Now().Add(-time.Duration(rand.Intn(48)) * time.Hour)
		completedAt := createdAt.Add(time.Duration(rand.Intn(20)) * time.Minute)

		tasks = append(tasks, UploadTask{
			ID:          uuid.New().String(),
			ModulePath:  modulePath,
			Version:     version,
			Source:      source,
			Status:      "failed",
			Progress:    rand.Intn(90), // 0-89
			Error:       "模块验证失败",
			CreatedAt:   createdAt,
			CompletedAt: completedAt,
			FileSize:    int64(10000 + rand.Intn(990000)), // 10KB-1MB大小
		})
	}

	// 添加一个正在处理的任务
	modulePath := popularModulePaths[rand.Intn(len(popularModulePaths))]
	version := moduleVersions[rand.Intn(len(moduleVersions))]
	source := "file"
	if rand.Intn(2) == 1 {
		source = "url"
	}

	tasks = append(tasks, UploadTask{
		ID:         uuid.New().String(),
		ModulePath: modulePath,
		Version:    version,
		Source:     source,
		Status:     "processing",
		Progress:   rand.Intn(90) + 10, // 10-99
		CreatedAt:  time.Now().Add(-time.Duration(rand.Intn(30)) * time.Minute),
		FileSize:   int64(10000 + rand.Intn(990000)), // 10KB-1MB大小
	})

	return tasks
}

// simulateTaskProcessing 模拟任务处理
func simulateTaskProcessing() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 获取上传任务列表（写锁）
		uploadTasksMutex.Lock()

		// 更新处理中的任务
		for i, task := range mockUploadTasks {
			if task.Status == "processing" {
				// 增加进度
				task.Progress += rand.Intn(10) + 1

				// 检查是否完成
				if task.Progress >= 100 {
					task.Progress = 100
					task.Status = "completed"
					task.CompletedAt = time.Now()
				}

				// 更新任务
				mockUploadTasks[i] = task
			}
		}

		uploadTasksMutex.Unlock()
	}
}

// uploadModuleAPIHandler 处理模块上传API请求
func uploadModuleAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 只支持POST方法
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var request struct {
		ModulePath string `json:"modulePath"`
		Version    string `json:"version"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
		return
	}

	// 验证必填字段
	if request.ModulePath == "" || request.Version == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "模块路径和版本为必填字段"})
		return
	}

	// 创建上传任务
	task := UploadTask{
		ID:         uuid.New().String(),
		ModulePath: request.ModulePath,
		Version:    request.Version,
		Source:     "file",
		Status:     "pending",
		Progress:   0,
		CreatedAt:  time.Now(),
		FileSize:   int64(10000 + rand.Intn(990000)), // 10KB-1MB大小
	}

	// 添加到任务列表（写锁）
	uploadTasksMutex.Lock()
	mockUploadTasks = append(mockUploadTasks, task)
	uploadTasksMutex.Unlock()

	// 返回创建的任务
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 模拟异步处理任务
	go func() {
		// 等待一段时间
		time.Sleep(2 * time.Second)

		// 更新任务状态（写锁）
		uploadTasksMutex.Lock()
		defer uploadTasksMutex.Unlock()

		// 查找任务
		for i, t := range mockUploadTasks {
			if t.ID == task.ID {
				// 更新状态
				mockUploadTasks[i].Status = "processing"
				mockUploadTasks[i].Progress = 10
				break
			}
		}
	}()
}

// uploadImportURLAPIHandler 处理从URL导入模块API请求
func uploadImportURLAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 只支持POST方法
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var request struct {
		ModulePath string `json:"modulePath"`
		Version    string `json:"version"`
		URL        string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
		return
	}

	// 验证必填字段
	if request.ModulePath == "" || request.Version == "" || request.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "模块路径、版本和URL为必填字段"})
		return
	}

	// 创建上传任务
	task := UploadTask{
		ID:         uuid.New().String(),
		ModulePath: request.ModulePath,
		Version:    request.Version,
		Source:     "url",
		Status:     "pending",
		Progress:   0,
		CreatedAt:  time.Now(),
	}

	// 添加到任务列表（写锁）
	uploadTasksMutex.Lock()
	mockUploadTasks = append(mockUploadTasks, task)
	uploadTasksMutex.Unlock()

	// 返回创建的任务
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 模拟异步处理任务
	go func() {
		// 等待一段时间
		time.Sleep(2 * time.Second)

		// 更新任务状态（写锁）
		uploadTasksMutex.Lock()
		defer uploadTasksMutex.Unlock()

		// 查找任务
		for i, t := range mockUploadTasks {
			if t.ID == task.ID {
				// 更新状态
				mockUploadTasks[i].Status = "processing"
				mockUploadTasks[i].Progress = 10
				break
			}
		}
	}()
}

// uploadTasksAPIHandler 处理上传任务列表API请求
func uploadTasksAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 只支持GET方法
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	query := r.URL.Query().Get("q")
	statusFilter := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// 解析limit和offset参数
	limit := 20 // 默认限制为20条
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // 默认偏移为0
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// 获取上传任务列表（读锁）
	uploadTasksMutex.RLock()
	tasks := mockUploadTasks
	uploadTasksMutex.RUnlock()

	// 过滤任务数据
	var filteredTasks []UploadTask
	for _, task := range tasks {
		// 如果有查询参数，则过滤模块路径和版本
		if query != "" && !strings.Contains(strings.ToLower(task.ModulePath), strings.ToLower(query)) && !strings.Contains(strings.ToLower(task.Version), strings.ToLower(query)) {
			continue
		}

		// 如果有状态过滤，则过滤任务状态
		if statusFilter != "" && task.Status != statusFilter {
			continue
		}

		filteredTasks = append(filteredTasks, task)
	}

	// 按创建时间排序（最新的在前）
	sortTasksByCreatedAt(filteredTasks)

	// 计算总数
	total := len(filteredTasks)

	// 应用分页
	start := offset
	end := offset + limit
	if start >= total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedTasks := filteredTasks
	if start < end {
		pagedTasks = filteredTasks[start:end]
	} else {
		pagedTasks = []UploadTask{}
	}

	// 构建响应
	response := map[string]interface{}{
		"tasks":  pagedTasks,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// sortTasksByCreatedAt 按创建时间排序上传任务
func sortTasksByCreatedAt(tasks []UploadTask) {
	// 简单的冒泡排序
	for i := 0; i < len(tasks)-1; i++ {
		for j := 0; j < len(tasks)-i-1; j++ {
			if tasks[j].CreatedAt.Before(tasks[j+1].CreatedAt) {
				tasks[j], tasks[j+1] = tasks[j+1], tasks[j]
			}
		}
	}
}

// uploadTaskDetailAPIHandler 处理上传任务详情API请求
func uploadTaskDetailAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 获取路径参数
	vars := mux.Vars(r)
	taskID := vars["taskId"]

	switch r.Method {
	case http.MethodGet:
		// 获取任务详情
		handleGetTaskDetail(w, r, taskID)
	default:
		// 不支持的HTTP方法
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleGetTaskDetail 处理获取任务详情的请求
func handleGetTaskDetail(w http.ResponseWriter, r *http.Request, taskID string) {
	// 获取上传任务列表（读锁）
	uploadTasksMutex.RLock()
	tasks := mockUploadTasks
	uploadTasksMutex.RUnlock()

	// 查找任务
	var foundTask *UploadTask
	for i, task := range tasks {
		if task.ID == taskID {
			foundTask = &tasks[i]
			break
		}
	}

	if foundTask == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "任务未找到"})
		return
	}

	// 返回任务详情
	if err := json.NewEncoder(w).Encode(foundTask); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// uploadTaskCancelAPIHandler 处理取消上传任务API请求
func uploadTaskCancelAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 获取路径参数
	vars := mux.Vars(r)
	taskID := vars["taskId"]

	// 只支持POST方法
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 获取上传任务列表（写锁）
	uploadTasksMutex.Lock()
	defer uploadTasksMutex.Unlock()

	// 查找任务
	var foundIndex = -1
	for i, task := range mockUploadTasks {
		if task.ID == taskID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "任务未找到"})
		return
	}

	// 检查任务状态
	if mockUploadTasks[foundIndex].Status == "completed" || mockUploadTasks[foundIndex].Status == "failed" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "无法取消已完成或已失败的任务"})
		return
	}

	// 更新任务状态
	mockUploadTasks[foundIndex].Status = "failed"
	mockUploadTasks[foundIndex].Error = "任务已取消"
	mockUploadTasks[foundIndex].CompletedAt = time.Now()

	// 返回更新后的任务
	if err := json.NewEncoder(w).Encode(mockUploadTasks[foundIndex]); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}