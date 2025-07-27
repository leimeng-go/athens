import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getDashboardData, getSystemStatus, getRecentActivities } from '@/api/dashboard'
import { DashboardData, SystemStatus } from '@/types'
import { ElMessage } from 'element-plus'

export const useDashboardStore = defineStore('dashboard', () => {
  // 状态
  const dashboardData = ref<DashboardData>({
    stats: {
      totalModules: 0,
      totalDownloads: 0,
      totalRepositories: 0,
      storageUsed: '0 MB'
    },
    downloadTrend: [],
    popularModules: [],
    recentActivities: []
  })
  const loading = ref(false)
  const systemStatus = ref<SystemStatus>({
    status: '',
    uptime: '',
    version: '',
    goVersion: '',
    memoryUsage: '',
    cpuUsage: ''
  })
  const systemStatusLoading = ref(false)

  // 获取仪表盘数据
  async function fetchDashboardData() {
    loading.value = true
    try {
      dashboardData.value = await getDashboardData()
    } catch (error) {
      console.error('Failed to fetch dashboard data:', error)
      ElMessage.error('获取仪表盘数据失败')
    } finally {
      loading.value = false
    }
  }

  // 获取系统状态
  async function fetchSystemStatus() {
    systemStatusLoading.value = true
    try {
      systemStatus.value = await getSystemStatus()
    } catch (error) {
      console.error('Failed to fetch system status:', error)
      ElMessage.error('获取系统状态失败')
    } finally {
      systemStatusLoading.value = false
    }
  }

  // 获取最近活动
  async function fetchRecentActivities(limit: number = 10) {
    try {
      const activities = await getRecentActivities(limit)
      dashboardData.value.recentActivities = activities
    } catch (error) {
      console.error('Failed to fetch recent activities:', error)
    }
  }

  // 刷新所有数据
  async function refreshAll() {
    await Promise.all([
      fetchDashboardData(),
      fetchSystemStatus()
    ])
  }

  return {
    dashboardData,
    loading,
    systemStatus,
    systemStatusLoading,
    fetchDashboardData,
    fetchSystemStatus,
    fetchRecentActivities,
    refreshAll
  }
})