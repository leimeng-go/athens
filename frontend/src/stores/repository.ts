import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getRepositories, getRepositoryStats, deleteRepository } from '@/api/repository'
import { RepositoryData, RepositoryListResponse } from '@/types'
import { ElMessage } from 'element-plus'

export const useRepositoryStore = defineStore('repository', () => {
  // 状态
  const repositories = ref<RepositoryData[]>([])
  const total = ref(0)
  const loading = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(10)
  const searchQuery = ref('')
  const sortBy = ref('lastUpdated')
  const sortOrder = ref<'asc' | 'desc'>('desc')
  const stats = ref({
    totalRepositories: 0,
    totalModules: 0,
    totalSize: '0 MB'
  })

  // 计算属性
  const isEmpty = computed(() => repositories.value.length === 0 && !loading.value)

  // 获取仓库列表
  async function fetchRepositories() {
    loading.value = true
    try {
      const response = await getRepositories({
        page: currentPage.value,
        pageSize: pageSize.value,
        search: searchQuery.value || undefined,
        sort: sortBy.value,
        order: sortOrder.value
      })
      repositories.value = response.items
      total.value = response.total
    } catch (error) {
      console.error('Failed to fetch repositories:', error)
      ElMessage.error('获取仓库列表失败')
    } finally {
      loading.value = false
    }
  }

  // 获取仓库统计信息
  async function fetchRepositoryStats() {
    try {
      stats.value = await getRepositoryStats()
    } catch (error) {
      console.error('Failed to fetch repository stats:', error)
    }
  }

  // 删除仓库
  async function removeRepository(id: string) {
    try {
      await deleteRepository(id)
      ElMessage.success('仓库删除成功')
      // 重新获取仓库列表
      await fetchRepositories()
      // 更新统计信息
      await fetchRepositoryStats()
    } catch (error) {
      console.error('Failed to delete repository:', error)
      ElMessage.error('删除仓库失败')
    }
  }

  // 更新分页和排序
  function updatePagination(page: number, size: number) {
    currentPage.value = page
    pageSize.value = size
    fetchRepositories()
  }

  // 更新搜索和排序
  function updateSearch(search: string) {
    searchQuery.value = search
    currentPage.value = 1 // 重置到第一页
    fetchRepositories()
  }

  function updateSort(sort: string, order: 'asc' | 'desc') {
    sortBy.value = sort
    sortOrder.value = order
    fetchRepositories()
  }

  return {
    repositories,
    total,
    loading,
    currentPage,
    pageSize,
    searchQuery,
    sortBy,
    sortOrder,
    stats,
    isEmpty,
    fetchRepositories,
    fetchRepositoryStats,
    removeRepository,
    updatePagination,
    updateSearch,
    updateSort
  }
})