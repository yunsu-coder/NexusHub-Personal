import axios from 'axios'
import { ElMessage } from 'element-plus'
import { ERROR_MESSAGES } from '../config/constants'
import router from '../router'
import config from '../config'

const api = axios.create({
  baseURL: config.api.baseURL,
  timeout: config.api.timeout,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器 - 公共模式下不附加认证头
api.interceptors.request.use(
  apiConfig => {
    try {
      window.dispatchEvent(new CustomEvent('api:loading', { detail: true }))
    } catch {}
    return apiConfig
  },
  error => {
    console.error('Request Error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器 - 统一错误处理
api.interceptors.response.use(
  response => {
    try {
      window.dispatchEvent(new CustomEvent('api:loading', { detail: false }))
    } catch {}
    // 检查响应格式
    const data = response.data

    // 如果响应有code字段,检查是否成功
    if (data && typeof data.code !== 'undefined') {
      if (data.code === 0) {
        // 成功响应,返回data字段
        return data.data !== undefined ? data.data : data
      } else {
        // 业务错误
        const errorMsg = data.message || 'Request failed'
        ElMessage.error(errorMsg)
        return Promise.reject(new Error(errorMsg))
      }
    }

    // 兼容旧格式
    return data
  },
  error => {
    console.error('Response Error:', error)
    try {
      window.dispatchEvent(new CustomEvent('api:loading', { detail: false }))
    } catch {}

    // 处理网络错误
    if (!error.response) {
      ElMessage.error('Network error. Please check your connection.')
      return Promise.reject(error)
    }

    const { status, data } = error.response
    let message = data?.message || data?.error || 'Request failed'

    // 根据状态码处理
    switch (status) {
      case 400:
        ElMessage.error(`${ERROR_MESSAGES.BAD_REQUEST || 'Bad Request'}: ${message}`)
        break
      case 401:
        ElMessage.error(ERROR_MESSAGES.UNAUTHORIZED || 'Public mode: no login required.')
        break
      case 403:
        ElMessage.error(ERROR_MESSAGES.FORBIDDEN || 'Access denied')
        break
      case 404:
        ElMessage.error(ERROR_MESSAGES.NOT_FOUND || 'Resource not found')
        break
      case 409:
        ElMessage.error(`Conflict: ${message}`)
        break
      case 413:
        ElMessage.error('File too large')
        break
      case 500:
        ElMessage.error(ERROR_MESSAGES.SERVER_ERROR || 'Server error. Please try again later.')
        break
      case 503:
        ElMessage.error('Service unavailable. Please try again later.')
        break
      default:
        ElMessage.error(`${ERROR_MESSAGES.UNKNOWN_ERROR || 'Error'}: ${message}`)
    }

    return Promise.reject(error)
  }
)

export default api
