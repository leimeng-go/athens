import { get, put } from './index'
import { SystemSettings } from '@/types'

// 模拟数据 - 用于开发阶段
const mockSystemSettings: SystemSettings = {
  storagePath: '/var/lib/athens',
  maxUploadSize: 50 * 1024 * 1024, // 50MB
  enablePrivateModules: true,
  enableDownloadLogging: true,
  proxyTimeout: 30, // 30秒
  cacheExpiration: 24 // 24小时
}

/**
 * 获取系统设置
 * @returns 系统设置信息
 */
export function getSystemSettings(): Promise<SystemSettings> {
  // 在开发环境中使用模拟数据
  if (import.meta.env.DEV) {
    return Promise.resolve(mockSystemSettings)
  }
  return get<SystemSettings>('/settings')
}

/**
 * 更新系统设置
 * @param settings 要更新的系统设置
 * @returns 更新后的系统设置
 */
export function updateSystemSettings(settings: SystemSettings): Promise<SystemSettings> {
  // 在开发环境中使用模拟数据
  if (import.meta.env.DEV) {
    return Promise.resolve({
      ...mockSystemSettings,
      ...settings
    })
  }
  return put<SystemSettings>('/settings', settings)
}