import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { uploadModule, importModuleFromUrl, getUploadTasks, cancelUploadTask } from '@/api/upload'
import { UploadTask } from '@/types'
import { ElMessage } from 'element-plus'

export const useUploadStore = defineStore('upload', () => {
  // 状态
  const uploadTasks = ref<UploadTask[]>([])
  const total = ref(0)
  const loading = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(10)
  const statusFilter = ref<'pending' | 'processing' | 'success' | 'failed' | ''>('')
  const uploadProgress = ref(0)
  const isUploading = ref(false)
  const currentTaskId = ref('')

  // 计算属性
  const isEmpty = computed(() => uploadTasks.value.length === 0 && !loading.value)

  // 获取上传任务列表
  async function fetchUploadTasks() {
    loading.value = true
    try {
      const response = await getUploadTasks({
        page: currentPage.value,
        pageSize: pageSize.value,
        status: statusFilter.value || undefined
      })
      uploadTasks.value = response.items
      total.value = response.total
    } catch (error) {
      console.error('Failed to fetch upload tasks:', error)
      ElMessage.error('获取上传任务列表失败')
    } finally {
      loading.value = false
    }
  }

  // 上传模块文件
  async function uploadModuleFile(file: File) {
    const formData = new FormData()
    formData.append('file', file)
    
    isUploading.value = true
    uploadProgress.value = 0
    
    try {
      const response = await uploadModule(formData)
      currentTaskId.value = response.taskId
      ElMessage.success('模块上传成功')
      await fetchUploadTasks() // 刷新任务列表
      return response.taskId
    } catch (error) {
      console.error('Failed to upload module:', error)
      ElMessage.error('模块上传失败')
      return null
    } finally {
      isUploading.value = false
      uploadProgress.value = 0
    }
  }

  // 从URL导入模块
  async function importFromUrl(url: string, version?: string) {
    try {
      const response = await importModuleFromUrl({ url, version })
      ElMessage.success('URL导入任务已创建')
      await fetchUploadTasks() // 刷新任务列表
      return response.taskId
    } catch (error) {
      console.error('Failed to import from URL:', error)
      ElMessage.error('从URL导入模块失败')
      return null
    }
  }

  // 取消上传任务
  async function cancelTask(taskId: string) {
    try {
      await cancelUploadTask(taskId)
      ElMessage.success('上传任务已取消')
      await fetchUploadTasks() // 刷新任务列表
    } catch (error) {
      console.error('Failed to cancel upload task:', error)
      ElMessage.error('取消上传任务失败')
    }
  }

  // 更新分页和过滤器
  function updatePagination(page: number, size: number) {
    currentPage.value = page
    pageSize.value = size
    fetchUploadTasks()
  }

  function updateStatusFilter(status: 'pending' | 'processing' | 'success' | 'failed' | '') {
    statusFilter.value = status
    currentPage.value = 1 // 重置到第一页
    fetchUploadTasks()
  }

  // 更新上传进度
  function updateUploadProgress(progress: number) {
    uploadProgress.value = progress
  }

  return {
    uploadTasks,
    total,
    loading,
    currentPage,
    pageSize,
    statusFilter,
    uploadProgress,
    isUploading,
    currentTaskId,
    isEmpty,
    fetchUploadTasks,
    uploadModuleFile,
    importFromUrl,
    cancelTask,
    updatePagination,
    updateStatusFilter,
    updateUploadProgress
  }
})