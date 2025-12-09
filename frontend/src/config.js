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
  upload: {
    maxSize: 100 * 1024 * 1024, // 100MB
    acceptedTypes: {
      image: ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml'],
      video: ['video/mp4', 'video/webm', 'video/ogg'],
      audio: ['audio/mpeg', 'audio/wav', 'audio/ogg', 'audio/flac'],
      document: ['application/pdf', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'],
      code: ['text/plain', 'text/x-python', 'text/x-java', 'text/javascript', 'application/json']
    }
  }
}
