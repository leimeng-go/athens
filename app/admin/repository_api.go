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

// 仓库类型
var repositoryTypes = []string{"git", "svn", "mercurial", "proxy"}

// 仓库状态
var repositoryStatuses = []string{"active", "syncing", "error", "inactive"}

// 模拟的仓库数据
var mockRepositories []RepositoryData

// 用于保护仓库数据的并发访问
var repositoriesMutex sync.RWMutex

// 初始化模拟数据
func init() {
	// 生成模拟仓库数据
	mockRepositories = generateMockRepositories()
}

// generateMockRepositories 生成模拟仓库数据
func generateMockRepositories() []RepositoryData {
	repositories := make([]RepositoryData, 0, 10)

	// 添加一些常见的仓库
	repositories = append(repositories, RepositoryData{
		ID:          uuid.New().String(),
		Name:        "GitHub",
		URL:         "https://github.com",
		Type:        "git",
		ModuleCount: 1000 + rand.Intn(9000),
		CreatedAt:   time.Now().Add(-time.Duration(30+rand.Intn(60)) * 24 * time.Hour),
		LastSync:    time.Now().Add(-time.Duration(rand.Intn(24)) * time.Hour),
		Status:      "active",
	})

	repositories = append(repositories, RepositoryData{
		ID:          uuid.New().String(),
		Name:        "GitLab",
		URL:         "https://gitlab.com",
		Type:        "git",
		ModuleCount: 500 + rand.Intn(4500),
		CreatedAt:   time.Now().Add(-time.Duration(20+rand.Intn(40)) * 24 * time.Hour),
		LastSync:    time.Now().Add(-time.Duration(rand.Intn(12)) * time.Hour),
		Status:      "active",
	})

	repositories = append(repositories, RepositoryData{
		ID:          uuid.New().String(),
		Name:        "Bitbucket",
		URL:         "https://bitbucket.org",
		Type:        "git",
		ModuleCount: 300 + rand.Intn(2700),
		CreatedAt:   time.Now().Add(-time.Duration(15+rand.Intn(30)) * 24 * time.Hour),
		LastSync:    time.Now().Add(-time.Duration(rand.Intn(48)) * time.Hour),
		Status:      "active",
	})

	repositories = append(repositories, RepositoryData{
		ID:          uuid.New().String(),
		Name:        "Go Modules Proxy",
		URL:         "https://proxy.golang.org",
		Type:        "proxy",
		ModuleCount: 5000 + rand.Intn(15000),
		CreatedAt:   time.Now().Add(-time.Duration(60+rand.Intn(30)) * 24 * time.Hour),
		LastSync:    time.Now().Add(-time.Duration(rand.Intn(6)) * time.Hour),
		Status:      "active",
	})

	// 添加一些随机仓库
	for i := 0; i < 6; i++ {
		name := "Repository-" + strconv.Itoa(i+1)
		repoType := repositoryTypes[rand.Intn(len(repositoryTypes))]
		status := repositoryStatuses[rand.Intn(len(repositoryStatuses))]
		
		repositories = append(repositories, RepositoryData{
			ID:          uuid.New().String(),
			Name:        name,
			URL:         "https://example.com/" + strings.ToLower(name),
			Type:        repoType,
			ModuleCount: rand.Intn(1000),
			CreatedAt:   time.Now().Add(-time.Duration(rand.Intn(90)) * 24 * time.Hour),
			LastSync:    time.Now().Add(-time.Duration(rand.Intn(72)) * time.Hour),
			Status:      status,
		})
	}

	return repositories
}

// repositoriesAPIHandler 处理仓库列表API请求
func repositoriesAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		// 获取仓库列表
		handleGetRepositories(w, r)
	case http.MethodPost:
		// 创建新仓库
		handleCreateRepository(w, r)
	default:
		// 不支持的HTTP方法
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleGetRepositories 处理获取仓库列表的请求
func handleGetRepositories(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数
	query := r.URL.Query().Get("q")
	typeFilter := r.URL.Query().Get("type")
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

	// 获取仓库列表（读锁）
	repositoriesMutex.RLock()
	repositories := mockRepositories
	repositoriesMutex.RUnlock()

	// 过滤仓库数据
	var filteredRepositories []RepositoryData
	for _, repo := range repositories {
		// 如果有查询参数，则过滤仓库名称和URL
		if query != "" && !strings.Contains(strings.ToLower(repo.Name), strings.ToLower(query)) && !strings.Contains(strings.ToLower(repo.URL), strings.ToLower(query)) {
			continue
		}

		// 如果有类型过滤，则过滤仓库类型
		if typeFilter != "" && repo.Type != typeFilter {
			continue
		}

		// 如果有状态过滤，则过滤仓库状态
		if statusFilter != "" && repo.Status != statusFilter {
			continue
		}

		filteredRepositories = append(filteredRepositories, repo)
	}

	// 计算总数
	total := len(filteredRepositories)

	// 应用分页
	start := offset
	end := offset + limit
	if start >= total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedRepositories := filteredRepositories
	if start < end {
		pagedRepositories = filteredRepositories[start:end]
	} else {
		pagedRepositories = []RepositoryData{}
	}

	// 构建响应
	response := map[string]interface{}{
		"repositories": pagedRepositories,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	}

	// 返回JSON响应
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// handleCreateRepository 处理创建新仓库的请求
func handleCreateRepository(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var newRepo RepositoryData
	if err := json.NewDecoder(r.Body).Decode(&newRepo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
		return
	}

	// 验证必填字段
	if newRepo.Name == "" || newRepo.URL == "" || newRepo.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "名称、URL和类型为必填字段"})
		return
	}

	// 设置ID和时间戳
	newRepo.ID = uuid.New().String()
	newRepo.CreatedAt = time.Now()
	newRepo.LastSync = time.Now()
	newRepo.Status = "active"

	// 添加到仓库列表（写锁）
	repositoriesMutex.Lock()
	mockRepositories = append(mockRepositories, newRepo)
	repositoriesMutex.Unlock()

	// 返回创建的仓库
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newRepo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// repositoryDetailAPIHandler 处理仓库详情API请求
func repositoryDetailAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 获取路径参数
	vars := mux.Vars(r)
	repoID := vars["id"]

	switch r.Method {
	case http.MethodGet:
		// 获取仓库详情
		handleGetRepositoryDetail(w, r, repoID)
	case http.MethodPut:
		// 更新仓库
		handleUpdateRepository(w, r, repoID)
	case http.MethodDelete:
		// 删除仓库
		handleDeleteRepository(w, r, repoID)
	default:
		// 不支持的HTTP方法
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleGetRepositoryDetail 处理获取仓库详情的请求
func handleGetRepositoryDetail(w http.ResponseWriter, r *http.Request, repoID string) {
	// 获取仓库列表（读锁）
	repositoriesMutex.RLock()
	repositories := mockRepositories
	repositoriesMutex.RUnlock()

	// 查找仓库
	var foundRepo *RepositoryData
	for i, repo := range repositories {
		if repo.ID == repoID {
			foundRepo = &repositories[i]
			break
		}
	}

	if foundRepo == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "仓库未找到"})
		return
	}

	// 返回仓库详情
	if err := json.NewEncoder(w).Encode(foundRepo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// handleUpdateRepository 处理更新仓库的请求
func handleUpdateRepository(w http.ResponseWriter, r *http.Request, repoID string) {
	// 解析请求体
	var updatedRepo RepositoryData
	if err := json.NewDecoder(r.Body).Decode(&updatedRepo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
		return
	}

	// 验证必填字段
	if updatedRepo.Name == "" || updatedRepo.URL == "" || updatedRepo.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "名称、URL和类型为必填字段"})
		return
	}

	// 更新仓库（写锁）
	repositoriesMutex.Lock()
	defer repositoriesMutex.Unlock()

	// 查找仓库
	var foundIndex = -1
	for i, repo := range mockRepositories {
		if repo.ID == repoID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "仓库未找到"})
		return
	}

	// 保留原始ID和创建时间
	updatedRepo.ID = repoID
	updatedRepo.CreatedAt = mockRepositories[foundIndex].CreatedAt

	// 更新最后同步时间
	updatedRepo.LastSync = time.Now()

	// 更新仓库
	mockRepositories[foundIndex] = updatedRepo

	// 返回更新后的仓库
	if err := json.NewEncoder(w).Encode(updatedRepo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// handleDeleteRepository 处理删除仓库的请求
func handleDeleteRepository(w http.ResponseWriter, r *http.Request, repoID string) {
	// 删除仓库（写锁）
	repositoriesMutex.Lock()
	defer repositoriesMutex.Unlock()

	// 查找仓库
	var foundIndex = -1
	for i, repo := range mockRepositories {
		if repo.ID == repoID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "仓库未找到"})
		return
	}

	// 删除仓库
	mockRepositories = append(mockRepositories[:foundIndex], mockRepositories[foundIndex+1:]...)

	// 返回成功响应
	w.WriteHeader(http.StatusNoContent)
}

// repositoryBatchDeleteAPIHandler 处理批量删除仓库API请求
func repositoryBatchDeleteAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 只支持DELETE方法
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var request struct {
		IDs []string `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "无效的请求数据"})
		return
	}

	// 验证ID列表
	if len(request.IDs) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID列表不能为空"})
		return
	}

	// 批量删除仓库（写锁）
	repositoriesMutex.Lock()
	defer repositoriesMutex.Unlock()

	// 创建ID集合，用于快速查找
	idSet := make(map[string]bool)
	for _, id := range request.IDs {
		idSet[id] = true
	}

	// 过滤出未删除的仓库
	var remainingRepositories []RepositoryData
	for _, repo := range mockRepositories {
		if !idSet[repo.ID] {
			remainingRepositories = append(remainingRepositories, repo)
		}
	}

	// 更新仓库列表
	mockRepositories = remainingRepositories

	// 返回成功响应
	w.WriteHeader(http.StatusNoContent)
}