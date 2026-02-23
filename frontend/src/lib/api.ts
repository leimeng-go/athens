/**
 * Athens Admin API 客户端
 * 所有 API 请求都会通过 Vite 代理转发到后端服务 (http://127.0.0.1:3000)
 */

const API_BASE = '/admin/api'

/**
 * 通用 API 请求封装
 */
async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  })

  if (!response.ok) {
    throw new Error(`API request failed: ${response.status} ${response.statusText}`)
  }

  return response.json()
}

/**
 * 系统状态相关 API
 */
export const systemAPI = {
  /**
   * 获取系统状态
   */
  getStatus: () => request<SystemStatus>(`${API_BASE}/system/status`),
}

/**
 * 仪表盘相关 API
 */
export const dashboardAPI = {
  /**
   * 获取仪表盘数据
   */
  getData: () => request<DashboardData>(`${API_BASE}/dashboard`),

  /**
   * 获取最近活动
   */
  getRecentActivities: (limit = 10) =>
    request<RecentActivity[]>(`${API_BASE}/activities/recent?limit=${limit}`),
}

/**
 * 配置管理相关 API
 */
export const configAPI = {
  /**
   * 获取配置
   */
  get: () => request<Config>(`${API_BASE}/config`),

  /**
   * 更新配置
   */
  update: (config: Partial<Config>) =>
    request<void>(`${API_BASE}/config`, {
      method: 'PUT',
      body: JSON.stringify(config),
    }),

  /**
   * 重置配置
   */
  reset: () =>
    request<void>(`${API_BASE}/config/reset`, {
      method: 'POST',
    }),

  /**
   * 验证配置
   */
  validate: (config: Partial<Config>) =>
    request<ValidationResult>(`${API_BASE}/config/validate`, {
      method: 'POST',
      body: JSON.stringify(config),
    }),
}

/**
 * 模块下载相关 API
 */
export const downloadAPI = {
  /**
   * 获取模块列表
   */
  getModules: () => request<Module[]>(`${API_BASE}/download/modules`),

  /**
   * 获取模块版本列表
   */
  getModuleVersions: (modulePath: string) =>
    request<ModuleVersion[]>(`${API_BASE}/download/modules/${modulePath}/versions`),

  /**
   * 获取模块详情
   */
  getModuleDetail: (modulePath: string) =>
    request<ModuleDetail>(`${API_BASE}/download/modules/${modulePath}`),

  /**
   * 获取下载统计
   */
  getStats: () => request<DownloadStats>(`${API_BASE}/download/stats`),

  /**
   * 获取热门模块
   */
  getPopular: () => request<PopularModule[]>(`${API_BASE}/download/popular`),

  /**
   * 获取最近下载
   */
  getRecent: () => request<RecentDownload[]>(`${API_BASE}/download/recent`),
}

/**
 * 仓库管理相关 API
 */
export const repositoryAPI = {
  /**
   * 获取仓库列表
   */
  list: () => request<Repository[]>(`${API_BASE}/repositories`),

  /**
   * 获取仓库详情
   */
  get: (id: string) => request<Repository>(`${API_BASE}/repositories/${id}`),

  /**
   * 批量删除仓库
   */
  batchDelete: (ids: string[]) =>
    request<void>(`${API_BASE}/repositories/batch-delete`, {
      method: 'POST',
      body: JSON.stringify({ ids }),
    }),
}

/**
 * 模块上传相关 API
 */
export const uploadAPI = {
  /**
   * 上传模块
   */
  uploadModule: (formData: FormData) =>
    fetch(`${API_BASE}/upload/module`, {
      method: 'POST',
      body: formData,
    }).then(res => res.json()),

  /**
   * 从 URL 导入模块
   */
  importFromURL: (url: string) =>
    request<UploadTask>(`${API_BASE}/upload/import-url`, {
      method: 'POST',
      body: JSON.stringify({ url }),
    }),

  /**
   * 获取上传任务列表
   */
  getTasks: () => request<UploadTask[]>(`${API_BASE}/upload/tasks`),

  /**
   * 获取上传任务详情
   */
  getTaskDetail: (taskId: string) =>
    request<UploadTask>(`${API_BASE}/upload/tasks/${taskId}`),

  /**
   * 取消上传任务
   */
  cancelTask: (taskId: string) =>
    request<void>(`${API_BASE}/upload/tasks/${taskId}/cancel`, {
      method: 'POST',
    }),
}

/**
 * 健康检查相关 API
 */
export const healthAPI = {
  /**
   * 健康检查
   */
  check: () => fetch('/healthz').then(res => res.ok),

  /**
   * 就绪检查
   */
  ready: () => fetch('/readyz').then(res => res.ok),

  /**
   * 获取版本信息
   */
  version: () => fetch('/version').then(res => res.text()),
}

// ============ 类型定义 ============

export interface SystemStatus {
  status: string
  uptime: string
  version: string
  goVersion: string
  memoryUsage: string
  cpuUsage: string
}

export interface DashboardData {
  stats: Stats
  downloadTrend: DownloadTrend[]
  popularModules: PopularModule[]
  recentActivities: RecentActivity[]
}

export interface Stats {
  totalModules: number
  totalDownloads: number
  totalRepositories: number
  storageUsed: string
}

export interface DownloadTrend {
  date: string
  count: number
}

export interface PopularModule {
  path: string
  downloads: number
}

export interface RecentActivity {
  id: string
  type: string
  message: string
  timestamp: string
}

export interface Config {
  [key: string]: any
}

export interface ValidationResult {
  valid: boolean
  errors?: string[]
}

export interface Module {
  path: string
  latestVersion: string
  description?: string
}

export interface ModuleVersion {
  version: string
  time: string
}

export interface ModuleDetail {
  path: string
  versions: ModuleVersion[]
  description?: string
  readme?: string
}

export interface DownloadStats {
  total: number
  today: number
  thisWeek: number
  thisMonth: number
}

export interface RecentDownload {
  module: string
  version: string
  timestamp: string
}

export interface Repository {
  id: string
  name: string
  url: string
  type: string
  status: string
}

export interface UploadTask {
  id: string
  status: string
  progress: number
  error?: string
  module?: string
  version?: string
}
