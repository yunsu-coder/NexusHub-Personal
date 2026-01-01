export default {
  api: {
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
    timeout: 30000 // 30 seconds
  },
  storage: {
    token: 'nexushub_token',
    user: 'nexushub_user',
    theme: 'nexushub_theme'
  },
  // 存储提供商配置
  storageProviders: {
    // 默认存储方式
    default: 'local',
    // 支持的存储提供商列表
    providers: {
      // 本地存储
      local: {
        name: '本地存储',
        type: 'local'
      },
      // 阿里云OSS
      oss: {
        name: '阿里云OSS',
        type: 'oss',
        endpoint: import.meta.env.VITE_OSS_ENDPOINT || '',
        accessKeyId: import.meta.env.VITE_OSS_ACCESS_KEY_ID || '',
        accessKeySecret: import.meta.env.VITE_OSS_ACCESS_KEY_SECRET || '',
        bucket: import.meta.env.VITE_OSS_BUCKET || '',
        region: import.meta.env.VITE_OSS_REGION || ''
      },
      // 腾讯云COS
      cos: {
        name: '腾讯云COS',
        type: 'cos',
        secretId: import.meta.env.VITE_COS_SECRET_ID || '',
        secretKey: import.meta.env.VITE_COS_SECRET_KEY || '',
        region: import.meta.env.VITE_COS_REGION || '',
        bucket: import.meta.env.VITE_COS_BUCKET || ''
      },
      // AWS S3
      s3: {
        name: 'AWS S3',
        type: 's3',
        accessKeyId: import.meta.env.VITE_S3_ACCESS_KEY_ID || '',
        secretAccessKey: import.meta.env.VITE_S3_SECRET_ACCESS_KEY || '',
        region: import.meta.env.VITE_S3_REGION || '',
        bucket: import.meta.env.VITE_S3_BUCKET || ''
      },
      // 七牛云Kodo
      qiniu: {
        name: '七牛云Kodo',
        type: 'qiniu',
        accessKey: import.meta.env.VITE_QINIU_ACCESS_KEY || '',
        secretKey: import.meta.env.VITE_QINIU_SECRET_KEY || '',
        bucket: import.meta.env.VITE_QINIU_BUCKET || '',
        region: import.meta.env.VITE_QINIU_REGION || '',
        domain: import.meta.env.VITE_QINIU_DOMAIN || ''
      }
    }
  },
  upload: {
    maxSize: 100 * 1024 * 1024, // 100MB
    timeout: 60000, // 60秒超时
    retryCount: 3, // 重试次数
    concurrentUploads: 3, // 并发上传数量
    chunkSize: 5 * 1024 * 1024, // 分片大小（5MB）
    allowedExtensions: true, // 是否允许所有文件扩展名
   思源笔记: {
      // 思源笔记常用配置选项
      enable: true,
      pathFormat: '{YYYY}/{MM}/{DD}', // 路径格式
      fileNameFormat: '{name}-{timestamp}', // 文件名格式
      autoRename: true, // 自动重命名冲突文件
      preserveOriginalName: true, // 保留原始文件名
      excludeExtensions: ['.tmp', '.swp', '.bak'] // 排除的文件扩展名
    },
    acceptedTypes: {
      image: ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.svg'],
      document: ['.pdf', '.doc', '.docx', '.ppt', '.pptx', '.xls', '.xlsx', '.txt', '.md', '.markdown'],
      code: ['.py', '.java', '.js', '.json', '.html', '.css', '.ts', '.cpp', '.go'],
      audio: ['.mp3', '.wav', '.ogg', '.flac'],
      video: ['.mp4', '.webm', '.ogg', '.avi', '.mov'],
      archive: ['.zip', '.rar', '.7z', '.tar', '.gz'],
      note: ['.sy', '.syp', '.syl'] // 思源笔记文件类型
    }
  }
}
