import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getModules, getDownloadStats, getPopularModules, getRecentDownloads } from '@/api/download'
import { ModuleData, DownloadStats } from '@/types'
import { ElMessage } from 'element-plus'

export const useDownloadStore = defineStore('download', () => {
  // 状态
  const modules = ref<ModuleData[]>([])
  const total = ref(0)
  const loading = ref(false)
  const statsLoading = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(10)
  const searchQuery = ref('')
  const sortBy = ref('downloads')
  const sortOrder = ref<'asc' | 'desc'>('desc')
  const stats = ref<DownloadStats>({
    totalDownloads: 0,
    todayDownloads: 0,
    weeklyDownloads: 0,
    monthlyDownloads: 0,
    downloadTrend: []
  })
  const popularModules = ref<{ path: string; downloads: number }[]>([])
  const recentDownloads = ref<{ path: string; version: string; timestamp: string }[]>([])

  // 计算属性
  const isEmpty = computed(() => modules.value.length === 0 && !loading.value)

  // 获取模块列表
  async function fetchModules() {
    loading.value = true
    try {
      const response = await getModules({
        page: currentPage.value,
        pageSize: pageSize.value,
        search: searchQuery.value || undefined,
        sort: sortBy.value,
        order: sortOrder.value
      })
      modules.value = response.items
      total.value = response.total
    } catch (error) {
      console.error('Failed to fetch modules:', error)
      ElMessage.error('获取模块列表失败')
    } finally {
      loading.value = false
    }
  }

  // 获取下载统计信息
  async function fetchDownloadStats() {
    statsLoading.value = true
    try {
      stats.value = await getDownloadStats()
    } catch (error) {
      console.error('Failed to fetch download stats:', error)
      ElMessage.error('获取下载统计信息失败')
    } finally {
      statsLoading.value = false
    }
  }

  // 获取热门模块
  async function fetchPopularModules(limit: number = 10) {
    try {
      popularModules.value = await getPopularModules(limit)
    } catch (error) {
      console.error('Failed to fetch popular modules:', error)
    }
  }

  // 获取最近下载的模块
  async function fetchRecentDownloads(limit: number = 10) {
    try {
      recentDownloads.value = await getRecentDownloads(limit)
    } catch (error) {
      console.error('Failed to fetch recent downloads:', error)
    }
  }

  // 更新分页和排序
  function updatePagination(page: number, size: number) {
    currentPage.value = page
    pageSize.value = size
    fetchModules()
  }

  // 更新搜索和排序
  function updateSearch(search: string) {
    searchQuery.value = search
    currentPage.value = 1 // 重置到第一页
    fetchModules()
  }

  function updateSort(sort: string, order: 'asc' | 'desc') {
    sortBy.value = sort
    sortOrder.value = order
    fetchModules()
  }

  return {
    modules,
    total,
    loading,
    statsLoading,
    currentPage,
    pageSize,
    searchQuery,
    sortBy,
    sortOrder,
    stats,
    popularModules,
    recentDownloads,
    isEmpty,
    fetchModules,
    fetchDownloadStats,
    fetchPopularModules,
    fetchRecentDownloads,
    updatePagination,
    updateSearch,
    updateSort
  }
})