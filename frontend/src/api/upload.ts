import { post, get } from './index'
import { UploadTaskListResponse, UploadTaskDetail } from '@/types'

// 上传模块文件
export function uploadModule(formData: FormData) {
  return post<{ taskId: string }>('/upload/module', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    // 上传进度回调
    onUploadProgress: (progressEvent) => {
      const percentCompleted = Math.round(
        (progressEvent.loaded * 100) / (progressEvent.total || 1)
      )
      console.log('上传进度:', percentCompleted)
    }
  })
}

// 从URL导入模块
export function importModuleFromUrl(data: { url: string; version?: string }) {
  return post<{ taskId: string }>('/upload/import-url', data)
}

// 获取上传任务列表
export function getUploadTasks(params: {
  page: number
  pageSize: number
  status?: 'pending' | 'processing' | 'success' | 'failed'
}) {
  return get<UploadTaskListResponse>('/upload/tasks', params)
}

// 获取上传任务详情
export function getUploadTaskDetail(taskId: string) {
  return get<UploadTaskDetail>(`/upload/tasks/${taskId}`)
}

// 取消上传任务
export function cancelUploadTask(taskId: string) {
  return post(`/upload/tasks/${taskId}/cancel`)
}