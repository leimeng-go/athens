import axios, { AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  
  timeout: 15000
})

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    // 可以在这里添加认证信息等
    return config
  },
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    
    // 根据后端API的响应格式进行处理
    if (response.status !== 200) {
      ElMessage({
        message: res.message || 'Error',
        type: 'error',
        duration: 5 * 1000
      })
      return Promise.reject(new Error(res.message || 'Error'))
    } else {
      return res
    }
  },
  (error) => {
    console.error('Response error:', error)
    ElMessage({
      message: error.message || '请求失败',
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

// 封装GET请求
export function get<T>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> {
  return service.get(url, { params, ...config })
}

// 封装POST请求
export function post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  return service.post(url, data, config)
}

// 封装PUT请求
export function put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  return service.put(url, data, config)
}

// 封装DELETE请求
export function del<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return service.delete(url, config)
}

export default service