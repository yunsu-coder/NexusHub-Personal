import api from '../api'
import config from '../config'
import { ElMessage } from 'element-plus'

// 获取用户配置的存储设置
const getUserStorageSettings = () => {
  try {
    const saved = localStorage.getItem('storage_settings')
    if (saved) {
      return JSON.parse(saved)
    }
  } catch (error) {
    console.error('Failed to get user storage settings:', error)
  }
  return null
}

/**
 * 存储服务基类
 */
class StorageService {
  constructor(type) {
    this.type = type
    // 合并默认配置和用户配置
    this.userSettings = getUserStorageSettings()
    this.config = this.mergeConfig()
    this.uploadConfig = config.upload
  }
  
  /**
   * 合并默认配置和用户配置
   * @returns {Object} 合并后的配置
   */
  mergeConfig() {
    const defaultConfig = config.storageProviders.providers[this.type]
    const userConfig = this.userSettings?.providers?.[this.type] || {}
    return { ...defaultConfig, ...userConfig }
  }

  /**
   * 上传文件
   * @param {File} file - 要上传的文件
   * @param {Object} options - 上传选项
   * @returns {Promise<Object>} - 上传结果
   */
  async upload(file, options = {}) {
    throw new Error('Not implemented')
  }

  /**
   * 获取文件URL
   * @param {String} fileId - 文件ID
   * @returns {String} - 文件URL
   */
  getFileUrl(fileId) {
    throw new Error('Not implemented')
  }
  
  /**
   * 列出存储中的文件
   * @param {Object} options - 列出选项
   * @returns {Promise<Array>} - 文件列表
   */
  async listFiles(options = {}) {
    throw new Error('Not implemented')
  }
  
  /**
   * 下载文件
   * @param {String} fileId - 文件ID
   * @param {Object} options - 下载选项
   * @returns {Promise<Blob>} - 文件Blob对象
   */
  async downloadFile(fileId, options = {}) {
    throw new Error('Not implemented')
  }
  
  /**
   * 删除文件
   * @param {String} fileId - 文件ID
   * @returns {Promise<Object>} - 删除结果
   */
  async deleteFile(fileId) {
    throw new Error('Not implemented')
  }
  
  /**
   * 同步文件
   * @param {Object} options - 同步选项
   * @returns {Promise<Object>} - 同步结果
   */
  async sync(options = {}) {
    throw new Error('Not implemented')
  }
  
  /**
   * 比较本地和远程文件差异
   * @param {Array} localFiles - 本地文件列表
   * @param {Array} remoteFiles - 远程文件列表
   * @returns {Object} - 文件差异
   */
  compareFiles(localFiles, remoteFiles) {
    // 比较本地和远程文件，返回差异对象
    const added = remoteFiles.filter(remote => !localFiles.find(local => local.name === remote.name))
    const removed = localFiles.filter(local => !remoteFiles.find(remote => remote.name === local.name))
    const modified = localFiles.filter(local => {
      const remote = remoteFiles.find(r => r.name === local.name)
      return remote && remote.updatedAt > local.updatedAt
    })
    
    return { added, removed, modified }
  }
}

/**
 * 本地存储服务
 */
class LocalStorageService extends StorageService {
  constructor() {
    super('local')
  }

  async upload(file, options = {}) {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('storage_type', this.type)
    
    if (options.path) {
      formData.append('path', options.path)
    }

    const headers = {
      'Content-Type': undefined
    }
    
    return api.post('/files/upload', formData, { headers })
  }

  getFileUrl(fileId) {
    return `${config.api.baseURL}/files/download/${fileId}`
  }
}

/**
 * 阿里云OSS存储服务
 */
class OSSStorageService extends StorageService {
  constructor() {
    super('oss')
    this.config = this.mergeConfig()
    this.uploadConfig = config.upload
  }

  async upload(file, options = {}) {
    try {
      const formData = new FormData()
      formData.append('file', file)
      formData.append('storage_type', this.type)
      
      // 添加OSS配置
      formData.append('endpoint', this.config.endpoint)
      formData.append('access_key_id', this.config.accessKeyId)
      formData.append('access_key_secret', this.config.accessKeySecret)
      formData.append('bucket', this.config.bucket)
      formData.append('region', this.config.region)
      
      if (options.path) {
        formData.append('path', options.path)
      }

      const headers = {
        'Content-Type': undefined
      }
      
      // 使用通用上传API路径
      return await api.post('/files/upload', formData, {
        headers,
        timeout: this.uploadConfig.timeout || 60000
      })
    } catch (err) {
      console.error('OSS Upload Error:', err)
      // 处理特定的API错误
      if (err.response && err.response.status === 404) {
        ElMessage.error('OSS上传功能未实现，请检查后端服务配置')
      }
      throw err
    }
  }

  getFileUrl(fileId) {
    return `${config.api.baseURL}/files/download/${fileId}`
  }
}

/**
 * 腾讯云COS存储服务
 */
class COSStorageService extends StorageService {
  constructor() {
    super('cos')
    this.config = this.mergeConfig()
    this.uploadConfig = config.upload
  }

  async upload(file, options = {}) {
    try {
      const formData = new FormData()
      formData.append('file', file)
      formData.append('storage_type', this.type)
      
      // 添加COS配置
      formData.append('secret_id', this.config.secretId)
      formData.append('secret_key', this.config.secretKey)
      formData.append('bucket', this.config.bucket)
      formData.append('region', this.config.region)
      
      if (options.path) {
        formData.append('path', options.path)
      }

      const headers = {
        'Content-Type': undefined
      }
      
      // 使用通用上传API路径
      return await api.post('/files/upload', formData, {
        headers,
        timeout: this.uploadConfig.timeout || 60000
      })
    } catch (err) {
      console.error('COS Upload Error:', err)
      // 处理特定的API错误
      if (err.response && err.response.status === 404) {
        ElMessage.error('COS上传功能未实现，请检查后端服务配置')
      }
      throw err
    }
  }

  getFileUrl(fileId) {
    return `${config.api.baseURL}/files/download/${fileId}`
  }
}

/**
 * AWS S3存储服务
 */
class S3StorageService extends StorageService {
  constructor() {
    super('s3')
    // 使用合并后的配置（默认配置 + 用户配置）
    this.config = this.mergeConfig()
    this.uploadConfig = config.upload
  }

  async upload(file, options = {}) {
    try {
      // 直接使用通用上传API，不调用特定的S3凭证API
      const formData = new FormData()
      formData.append('file', file)
      formData.append('storage_type', this.type)
      
      // 添加S3配置
      formData.append('endpoint', this.config.endpoint)
      formData.append('access_key_id', this.config.accessKeyId)
      formData.append('secret_access_key', this.config.secretAccessKey)
      formData.append('bucket', this.config.bucket)
      formData.append('region', this.config.region)
      
      if (options.path) {
        formData.append('path', options.path)
      }

      const headers = {
        'Content-Type': undefined
      }
      
      // 使用通用上传API路径
      return await api.post('/files/upload', formData, {
        headers,
        timeout: this.uploadConfig.timeout || 60000
      })
    } catch (err) {
      console.error('S3 Upload Error:', err)
      // 处理特定的API错误
      if (err.response && err.response.status === 404) {
        // API不存在，可能是后端没有实现S3上传功能
        ElMessage.error('S3上传功能未实现，请检查后端服务配置')
      }
      throw err
    }
  }

  getFileUrl(fileId) {
    return `${config.api.baseURL}/files/download/${fileId}`
  }
}

/**
 * 七牛云Kodo存储服务
 */
class QiniuStorageService extends StorageService {
  constructor() {
    super('qiniu')
    this.config = this.mergeConfig()
    this.uploadConfig = config.upload
    // 同步状态
    this.syncStatus = {
      isSyncing: false,
      progress: 0,
      status: 'idle' // idle, syncing, success, error
    }
  }

  async upload(file, options = {}) {
    try {
      const formData = new FormData()
      formData.append('file', file)
      formData.append('storage_type', this.type)
      
      // 添加七牛云配置
      formData.append('access_key', this.config.accessKey)
      formData.append('secret_key', this.config.secretKey)
      formData.append('bucket', this.config.bucket)
      formData.append('region', this.config.region)
      formData.append('domain', this.config.domain)
      
      if (options.path) {
        formData.append('path', options.path)
      }

      const headers = {
        'Content-Type': undefined
      }
      
      // 使用通用上传API路径
      return await api.post('/files/upload', formData, {
        headers,
        timeout: this.uploadConfig.timeout || 60000
      })
    } catch (err) {
      console.error('Qiniu Upload Error:', err)
      // 处理特定的API错误
      if (err.response && err.response.status === 404) {
        ElMessage.error('七牛云上传功能未实现，请检查后端服务配置')
      }
      throw err
    }
  }

  getFileUrl(fileId) {
    return `${config.api.baseURL}/files/download/${fileId}`
  }
  
  /**
   * 列出七牛云存储中的文件
   * @param {Object} options - 列出选项
   * @returns {Promise<Array>} - 文件列表
   */
  async listFiles(options = {}) {
    try {
      // 注意：当前后端未实现七牛云文件列表API，这里返回模拟数据
      // 实际项目中应该调用后端API获取七牛云文件列表
      console.log('Qiniu listFiles called with options:', options)
      
      // 由于后端未实现此API，返回空列表
      // 这将导致同步时认为没有远程文件，会上传所有本地文件
      return []
    } catch (err) {
      console.error('Qiniu List Files Error:', err)
      // 失败时返回空列表
      return []
    }
  }
  
  /**
   * 七牛云文件同步
   * @param {Object} options - 同步选项
   * @returns {Promise<Object>} - 同步结果
   */
  async sync(options = {}) {
    try {
      this.syncStatus.isSyncing = true
      this.syncStatus.progress = 0
      this.syncStatus.status = 'syncing'
      
      // 通知UI同步开始
      this.notifySyncStatus()
      
      // 1. 获取本地文件列表
      const localFiles = await api.get('/files', {
        params: {
          path: options.path || ''
        }
      })
      
      this.syncStatus.progress = 30
      this.notifySyncStatus()
      
      // 2. 注意：当前后端未实现七牛云文件列表API
      // 由于没有远程文件列表，我们将把所有本地文件视为新增文件
      // 在实际项目中，应该比较本地和远程文件差异
      console.log('Qiniu sync: remote file listing not implemented, uploading all local files')
      
      this.syncStatus.progress = 50
      this.notifySyncStatus()
      
      // 3. 执行同步操作
      const syncResult = {
        added: [],
        updated: [],
        deleted: [],
        skipped: 0
      }
      
      // 上传所有本地文件到七牛云
      for (const file of localFiles) {
        try {
          // 下载本地文件
          const response = await fetch(`${config.api.baseURL}/files/download/${file.id}`)
          const blob = await response.blob()
          const fileObj = new File([blob], file.file_name, { type: blob.type })
          
          // 上传到七牛云
          await this.upload(fileObj, {
            path: options.path || ''
          })
          
          syncResult.added.push(file)
        } catch (err) {
          console.error(`Failed to sync file ${file.file_name}:`, err)
          syncResult.skipped++
        }
      }
      
      this.syncStatus.progress = 80
      this.notifySyncStatus()
      
      // 4. 同步完成
      this.syncStatus.isSyncing = false
      this.syncStatus.progress = 100
      this.syncStatus.status = 'success'
      this.notifySyncStatus()
      
      ElMessage.success(`七牛云同步完成：上传 ${syncResult.added.length} 个文件，跳过 ${syncResult.skipped} 个文件`)
      return syncResult
    } catch (err) {
      console.error('Qiniu Sync Error:', err)
      this.syncStatus.isSyncing = false
      this.syncStatus.status = 'error'
      this.notifySyncStatus()
      
      // 处理不同类型的错误
      if (err.response && err.response.status === 404) {
        ElMessage.error('七牛云同步功能未实现，请检查后端服务配置')
      } else if (err.response && err.response.status === 403) {
        ElMessage.error('七牛云同步失败：权限不足，请检查存储配置')
      } else if (err.response && err.response.status === 500) {
        ElMessage.error('七牛云同步失败：服务器内部错误')
      } else {
        ElMessage.error(`七牛云同步失败：${err.message || '未知错误'}`)
      }
      throw err
    }
  }
  
  /**
   * 通知UI同步状态变化
   */
  notifySyncStatus() {
    try {
      window.dispatchEvent(new CustomEvent('storage:syncStatus', {
        detail: {
          storageType: this.type,
          status: this.syncStatus
        }
      }))
    } catch (err) {
      console.error('Failed to notify sync status:', err)
    }
  }
  
  /**
   * 获取当前同步状态
   * @returns {Object} 同步状态
   */
  getSyncStatus() {
    return this.syncStatus
  }
}

/**
 * 存储服务工厂
 */
class StorageServiceFactory {
  /**
   * 获取存储服务实例
   * @param {String} type - 存储类型
   * @returns {StorageService} - 存储服务实例
   */
  static getStorageService(type = null) {
    // 如果没有指定类型，使用用户配置的默认存储方式
    if (!type) {
      const userSettings = getUserStorageSettings()
      type = userSettings?.defaultStorage || config.storageProviders.default
    }
    
    switch (type) {
      case 'local':
        return new LocalStorageService()
      case 'oss':
        return new OSSStorageService()
      case 'cos':
        return new COSStorageService()
      case 's3':
        return new S3StorageService()
      case 'qiniu':
        return new QiniuStorageService()
      default:
        throw new Error(`Unsupported storage type: ${type}`)
    }
  }
}

export default StorageServiceFactory