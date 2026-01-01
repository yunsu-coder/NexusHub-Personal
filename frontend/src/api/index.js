import axios from 'axios'
import { ElMessage } from 'element-plus'
import { ERROR_MESSAGES } from '../config/constants'
import router from '../router'
import config from '../config'

// 请求缓存对象
const requestCache = new Map()

// 缓存配置
const cacheConfig = {
  enabled: true,
  ttl: 5 * 60 * 1000, // 缓存有效期5分钟
  // 不需要缓存的请求方法
  noCacheMethods: ['post', 'put', 'delete', 'patch'],
  // 需要缓存的请求路径
  cachePaths: ['/files', '/notes', '/todos', '/bookmarks']
}

const api = axios.create({
  baseURL: config.api.baseURL,
  timeout: config.api.timeout,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 生成请求缓存键
const generateCacheKey = (config) => {
  const { url, method, params, data } = config
  return `${method}:${url}:${JSON.stringify(params)}:${JSON.stringify(data)}`
}

// 检查请求是否需要缓存
const shouldCacheRequest = (config) => {
  if (!cacheConfig.enabled) return false
  if (cacheConfig.noCacheMethods.includes(config.method.toLowerCase())) return false
  
  // 检查路径是否需要缓存
  const url = config.url || ''
  return cacheConfig.cachePaths.some(path => url.includes(path))
}

// 请求拦截器 - 公共模式下不附加认证头
api.interceptors.request.use(
  apiConfig => {
    try {
      window.dispatchEvent(new CustomEvent('api:loading', { detail: true }))
    } catch {}
    
    // 检查是否有缓存
    if (shouldCacheRequest(apiConfig)) {
      const cacheKey = generateCacheKey(apiConfig)
      const cached = requestCache.get(cacheKey)
      
      if (cached) {
        // 检查缓存是否过期
        if (Date.now() - cached.timestamp < cacheConfig.ttl) {
          // 直接返回缓存数据
          window.dispatchEvent(new CustomEvent('api:loading', { detail: false }))
          return Promise.resolve(cached.response)
        } else {
          // 缓存过期，移除缓存
          requestCache.delete(cacheKey)
        }
      }
    }
    
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
    let responseData
    if (data && typeof data.code !== 'undefined') {
      if (data.code === 0) {
        // 成功响应,返回data字段
        responseData = data.data !== undefined ? data.data : data
      } else {
        // 业务错误
        const errorMsg = data.message || 'Request failed'
        ElMessage.error(errorMsg)
        return Promise.reject(new Error(errorMsg))
      }
    } else {
      // 兼容旧格式
      responseData = data
    }
    
    // 缓存成功响应
    if (shouldCacheRequest(response.config)) {
      const cacheKey = generateCacheKey(response.config)
      // 存储响应数据和时间戳
      requestCache.set(cacheKey, {
        response: {
          ...response,
          data: responseData
        },
        timestamp: Date.now()
      })
    }
    
    return responseData
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
