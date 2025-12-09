// 常量定义

// HTTP 状态码
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  INTERNAL_SERVER_ERROR: 500,
}

// 任务状态
export const TASK_STATUS = {
  PENDING: 'pending',
  IN_PROGRESS: 'in_progress',
  COMPLETED: 'completed',
}

// 任务优先级
export const TASK_PRIORITY = {
  LOW: 'low',
  MEDIUM: 'medium',
  HIGH: 'high',
}

// 文件类型图标映射
export const FILE_TYPE_ICONS = {
  folder: 'Folder',
  image: 'Picture',
  document: 'Document',
  code: 'DocumentCopy',
  archive: 'Box',
  video: 'VideoCamera',
  audio: 'Headset',
  unknown: 'Document',
}

// 文件扩展名分类
export const FILE_EXTENSIONS = {
  image: ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.svg', '.webp', '.ico'],
  document: ['.pdf', '.doc', '.docx', '.xls', '.xlsx', '.ppt', '.pptx', '.txt', '.md'],
  code: ['.js', '.ts', '.vue', '.jsx', '.tsx', '.go', '.py', '.java', '.c', '.cpp', '.html', '.css', '.json', '.xml'],
  archive: ['.zip', '.rar', '.7z', '.tar', '.gz'],
  video: ['.mp4', '.avi', '.mov', '.wmv', '.flv', '.mkv'],
  audio: ['.mp3', '.wav', '.flac', '.aac', '.ogg'],
}

// 主题配置
export const THEMES = {
  LIGHT: 'light',
  DARK: 'dark',
  AUTO: 'auto',
}

// 错误消息
export const ERROR_MESSAGES = {
  NETWORK_ERROR: '网络连接失败，请检查网络设置',
  TIMEOUT_ERROR: '请求超时，请稍后重试',
  UNAUTHORIZED: '登录已过期，请重新登录',
  FORBIDDEN: '没有权限执行此操作',
  NOT_FOUND: '请求的资源不存在',
  SERVER_ERROR: '服务器错误，请稍后重试',
  UNKNOWN_ERROR: '发生未知错误',
}

// 成功消息
export const SUCCESS_MESSAGES = {
  SAVE_SUCCESS: '保存成功',
  DELETE_SUCCESS: '删除成功',
  UPDATE_SUCCESS: '更新成功',
  CREATE_SUCCESS: '创建成功',
  UPLOAD_SUCCESS: '上传成功',
  COPY_SUCCESS: '复制成功',
}

// 路由路径
export const ROUTES = {
  HOME: '/',
  AUTH: '/auth',
  DASHBOARD: '/dashboard',
  NOTES: '/notes',
  TASKS: '/tasks',
  CODE: '/code',
  BOOKMARKS: '/bookmarks',

  FILE_MANAGER: '/files',
  DATA_ANALYSIS: '/data',
  SETTINGS: '/settings',
  WELCOME: '/',
}

// 数据分析图表类型
export const CHART_TYPES = {
  LINE: 'line',
  BAR: 'bar',
  PIE: 'pie',
  SCATTER: 'scatter',
  AREA: 'area',
}
