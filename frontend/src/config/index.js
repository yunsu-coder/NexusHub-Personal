// 应用配置文件
const config = {
  // API 配置
  api: {
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
    timeout: 30000,
  },

  // 应用配置
  app: {
    name: 'NexusHub',
    version: '3.0.1',
  },

  // 存储键名
  storage: {
    token: 'token',
    user: 'user',
    theme: 'theme',
  },

  // 文件上传配置
  upload: {
    maxSize: 100 * 1024 * 1024, // 100MB
    acceptedTypes: {
      image: ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.svg', '.bmp', '.ico', '.heic', '.tiff'],
      document: ['.pdf', '.doc', '.docx', '.txt', '.md', '.rtf', '.xls', '.xlsx', '.ppt', '.pptx'],
      code: ['.js', '.ts', '.vue', '.jsx', '.tsx', '.go', '.py', '.java', '.c', '.cpp', '.cs', '.php', '.rb', '.rs', '.sh', '.json', '.xml', '.yaml', '.yml'],
      audio: ['.mp3', '.wav', '.flac', '.ogg', '.aac', '.m4a'],
      video: ['.mp4', '.webm', '.avi', '.mov', '.wmv', '.flv', '.mkv', '.m4v'],
      archive: ['.zip', '.rar', '.7z', '.tar', '.gz', '.bz2'],
      all: '*',
    },
  },

  // 分页配置
  pagination: {
    defaultPageSize: 20,
    pageSizes: [10, 20, 50, 100],
  },
}

export default config
