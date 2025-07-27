import { get, post, del } from './index'
import { RepositoryData, RepositoryListResponse } from '@/types'

// 获取仓库列表
export function getRepositories(params: {
  page: number
  pageSize: number
  search?: string
  sort?: string
  order?: 'asc' | 'desc'
}) {
  return get<RepositoryListResponse>('/repositories', params)
}

// 获取仓库详情
export function getRepositoryDetail(id: string) {
  return get<RepositoryData>(`/repositories/${id}`)
}

// 删除仓库
export function deleteRepository(id: string) {
  return del(`/repositories/${id}`)
}

// 批量删除仓库
export function batchDeleteRepositories(ids: string[]) {
  return post('/repositories/batch-delete', { ids })
}

// 获取仓库统计信息
export function getRepositoryStats() {
  return get('/repositories/stats')
}