// 仪表盘数据类型
export interface DashboardData {
  stats: {
    totalModules: number
    totalDownloads: number
    totalRepositories: number
    storageUsed: string
  }
  downloadTrend: {
    date: string
    count: number
  }[]
  popularModules: {
    path: string
    downloads: number
  }[]
  recentActivities: {
    id: string
    type: 'download' | 'upload' | 'system'
    message: string
    timestamp: string
    details?: Record<string, any>
  }[]
}

// 仓库数据类型
export interface RepositoryData {
  id: string
  name: string
  path: string
  moduleCount: number
  size: string
  lastUpdated: string
}

export interface RepositoryListResponse {
  items: RepositoryData[]
  total: number
  page: number
  pageSize: number
}

// 上传任务类型
export interface UploadTask {
  id: string
  modulePath: string
  version: string
  status: 'pending' | 'processing' | 'success' | 'failed'
  createdAt: string
  updatedAt: string
  error?: string
}

export interface UploadTaskListResponse {
  items: UploadTask[]
  total: number
  page: number
  pageSize: number
}

export interface UploadTaskDetail extends UploadTask {
  logs: string[]
  fileSize: string
  uploadedBy: string
}

// 模块数据类型
export interface ModuleData {
  path: string
  latestVersion: string
  versions: number
  downloads: number
  lastDownloaded: string
  size: string
}

export interface ModuleListResponse {
  items: ModuleData[]
  total: number
  page: number
  pageSize: number
}

export interface ModuleDetail {
  path: string
  versions: {
    version: string
    size: string
    createdAt: string
    downloads: number
  }[]
  totalDownloads: number
  firstAdded: string
  lastDownloaded: string
}

// 下载统计类型
export interface DownloadStats {
  totalDownloads: number
  todayDownloads: number
  weeklyDownloads: number
  monthlyDownloads: number
  downloadTrend: {
    date: string
    count: number
  }[]
}

// 系统状态类型
export interface SystemStatus {
  status: string
  uptime: string
  version: string
  goVersion: string
  memoryUsage: string
  cpuUsage: string
}

// 系统设置类型
export interface SystemSettings {
  storagePath: string
  maxUploadSize: number
  enablePrivateModules: boolean
  enableDownloadLogging: boolean
  proxyTimeout: number
  cacheExpiration: number
}