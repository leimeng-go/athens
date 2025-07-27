import { get } from './index'
import { DashboardData, SystemStatus } from '@/types'

// 模拟数据 - 用于开发阶段
const mockDashboardData: DashboardData = {
  stats: {
    totalModules: 1250,
    totalDownloads: 45678,
    totalRepositories: 35,
    storageUsed: '2.3 GB'
  },
  downloadTrend: [
    { date: '2023-01-01', count: 120 },
    { date: '2023-01-02', count: 145 },
    { date: '2023-01-03', count: 132 },
    { date: '2023-01-04', count: 167 },
    { date: '2023-01-05', count: 189 },
    { date: '2023-01-06', count: 156 },
    { date: '2023-01-07', count: 123 }
  ],
  popularModules: [
    { path: 'github.com/gorilla/mux', downloads: 12345 },
    { path: 'github.com/gin-gonic/gin', downloads: 10234 },
    { path: 'github.com/spf13/cobra', downloads: 8765 },
    { path: 'github.com/stretchr/testify', downloads: 7654 },
    { path: 'github.com/prometheus/client_golang', downloads: 6543 }
  ],
  recentActivities: [
    { id: '1', type: 'download', message: '模块 github.com/gorilla/mux@v1.8.0 被下载', timestamp: '2023-01-07T12:34:56Z' },
    { id: '2', type: 'upload', message: '模块 github.com/gin-gonic/gin@v1.9.0 被上传', timestamp: '2023-01-07T10:23:45Z' },
    { id: '3', type: 'system', message: '系统更新完成', timestamp: '2023-01-06T22:12:34Z' },
    { id: '4', type: 'download', message: '模块 github.com/spf13/cobra@v1.6.1 被下载', timestamp: '2023-01-06T18:45:12Z' },
    { id: '5', type: 'upload', message: '模块 github.com/stretchr/testify@v1.8.1 被上传', timestamp: '2023-01-06T15:34:23Z' }
  ]
};

const mockSystemStatus: SystemStatus = {
  status: 'healthy',
  uptime: '3d 12h 45m',
  version: 'v0.11.0',
  goVersion: 'go1.19.5',
  memoryUsage: '256MB / 1GB',
  cpuUsage: '12%'
};

// 获取仪表盘数据
export function getDashboardData(): Promise<DashboardData> {
  // 在开发环境中使用模拟数据
  if (import.meta.env.DEV) {
    return Promise.resolve(mockDashboardData);
  }
  return get<DashboardData>('/dashboard');
}

// 获取系统状态
export function getSystemStatus(): Promise<SystemStatus> {
  // 在开发环境中使用模拟数据
  if (import.meta.env.DEV) {
    return Promise.resolve(mockSystemStatus);
  }
  return get<SystemStatus>('/system/status');
}

// 获取最近活动
export function getRecentActivities(limit: number = 10): Promise<DashboardData['recentActivities']> {
  // 在开发环境中使用模拟数据
  if (import.meta.env.DEV) {
    return Promise.resolve(mockDashboardData.recentActivities.slice(0, limit));
  }
  return get<DashboardData['recentActivities']>('/activities/recent', { limit });
}