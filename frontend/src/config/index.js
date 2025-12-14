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
    maxSize: 50 * 1024 * 1024, // 50MB
    acceptedTypes: {
      image: ['.jpg', '.jpeg', '.png', '.gif', '.webp'],
      document: ['.pdf', '.doc', '.docx', '.txt', '.md'],
      code: ['.js', '.ts', '.vue', '.jsx', '.tsx', '.go', '.py', '.java'],
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
