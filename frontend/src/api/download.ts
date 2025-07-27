import { get } from './index'
import { ModuleListResponse, ModuleDetail, DownloadStats } from '@/types'

// 获取模块列表
export function getModules(params: {
  page: number
  pageSize: number
  search?: string
  sort?: string
  order?: 'asc' | 'desc'
}) {
  return get<ModuleListResponse>('/download/modules', params)
}

// 获取模块详情
export function getModuleDetail(path: string) {
  return get<ModuleDetail>(`/download/modules/${encodeURIComponent(path)}`)
}

// 获取模块版本列表
export function getModuleVersions(path: string) {
  return get<string[]>(`/download/modules/${encodeURIComponent(path)}/versions`)
}

// 获取下载统计信息
export function getDownloadStats() {
  return get<DownloadStats>('/download/stats')
}

// 获取热门模块
export function getPopularModules(limit: number = 10) {
  return get<{ path: string; downloads: number }[]>('/download/popular', { limit })
}

// 获取最近下载的模块
export function getRecentDownloads(limit: number = 10) {
  return get<{ path: string; version: string; timestamp: string }[]>('/download/recent', { limit })
}