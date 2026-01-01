<template>
  <div class="dashboard">
    <!-- 顶部统计总览 -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-label">笔记</div>
        <div class="stat-value">{{ notesCount }}</div>
        <div class="stat-sub">最近 3 篇已加载</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">未完成待办</div>
        <div class="stat-value">{{ openTodosCount }}</div>
        <div class="stat-sub">今日聚焦任务</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">文件</div>
        <div class="stat-value">{{ filesCount }}</div>
        <div class="stat-sub">已上传总数</div>
      </div>
      <div class="stat-card system-card">
        <div class="system-row">
          <span class="stat-label">系统状态</span>
          <span class="status-pill" :class="{ online: isBackendOnline, offline: !isBackendOnline }">
            {{ isBackendOnline ? '在线' : '离线' }}
          </span>
        </div>
        <div class="stat-sub">版本 {{ appVersion }}</div>
      </div>
    </div>

    <!-- 快捷链接区域 -->
    <div class="quick-links-section">
      <QuickLinks />
    </div>

    <!-- 桌面小组件网格 -->
    <div class="desktop-grid">
      <!-- 天气卡片 -->
      <div class="widget-item weather-widget">
        <div class="weather-content" v-if="weatherData">
          <div class="weather-icon">{{ getWeatherIcon(weatherData.weathercode) }}</div>
          <div class="weather-info">
            <div class="temperature">{{ weatherData.temperature }}°C</div>
            <div class="location">{{ locationName }} (实时)</div>
          </div>
        </div>
        <div class="weather-content" v-else>
          <div class="weather-icon">⌛</div>
          <div class="weather-info">
            <div class="temperature">--°C</div>
            <div class="location">获取中...</div>
          </div>
        </div>
      </div>

      <!-- 3D 宠物卡片 -->
      <div class="widget-item pet-widget" @click="petInteract">
        <div class="pet-container">
          <div class="pet-3d-model" :class="{ bounce: isPetBouncing }">
            <div class="pet-figure">{{ petEmoji }}</div>
          </div>
          <div class="pet-message" v-if="petMessage">{{ petMessage }}</div>
        </div>
      </div>

      <!-- 待办事项 (列表样式) -->
      <div class="widget-item todo-widget wide-widget">
        <div class="widget-header">
          <span>📝 待办事项</span>
          <el-tooltip content="查看所有待办事项" placement="top" :offset="10">
            <el-button link size="small" @click="$router.push('/todos')">更多</el-button>
          </el-tooltip>
        </div>
        <div class="todo-list-mini">
           <div v-for="todo in todayTodos" :key="todo.id" class="mini-todo-item">
             <div class="todo-dot" :class="{ completed: todo.completed }"></div>
             <span class="todo-title" :class="{ completed: todo.completed }">{{ todo.title }}</span>
           </div>
           <div v-if="todayTodos.length === 0" class="empty-mini">无待办</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import QuickLinks from '../components/QuickLinks.vue'
import api from '../api'
import config from '../config'

// 组件状态管理
const todayTodos = ref([])
const notesCount = ref(0)
const openTodosCount = ref(0)
const filesCount = ref(0)
const isBackendOnline = ref(true)
const appVersion = ref('3.0.1')
const petEmoji = ref('🐱')
const petMessage = ref('')
const isPetBouncing = ref(false)
const weatherData = ref(null)
const locationName = ref('上海')
const userLocation = ref(null)
const rssFeed = ref([]) // 添加rssFeed变量定义

const fetchWeather = async () => {
  try {
    // 获取用户IP地址和地理位置信息
    const geoRes = await fetch('https://ipapi.co/json/')
    const geoData = await geoRes.json()
    
    if (geoData.latitude && geoData.longitude) {
      const { latitude, longitude, city, country_name } = geoData
      locationName.value = `${city}, ${country_name}`
      
      // 使用获取的坐标获取天气数据
      const weatherRes = await fetch(`https://api.open-meteo.com/v1/forecast?latitude=${latitude}&longitude=${longitude}&current_weather=true`)
      const weatherJson = await weatherRes.json()
      
      if (weatherJson.current_weather) {
        weatherData.value = weatherJson.current_weather
      }
    }
  } catch (e) {
    console.error('Weather fetch failed', e)
    // 失败时使用默认位置（上海）
    locationName.value = '上海'
    const defaultRes = await fetch('https://api.open-meteo.com/v1/forecast?latitude=31.22&longitude=121.46&current_weather=true')
    const defaultData = await defaultRes.json()
    if (defaultData.current_weather) {
      weatherData.value = defaultData.current_weather
    }
  }
}

const fetchRSS = async () => {
  try {
    // 默认使用 少数派 RSS
    const res = await api.post('/rss/feed', { url: 'https://sspai.com/feed' })
    rssFeed.value = res.items
  } catch (e) {
    rssFeed.value = []
  }
}

const getWeatherIcon = (code) => {
  if (code === 0) return '☀️'
  if (code <= 3) return '⛅'
  if (code <= 48) return '🌫️'
  if (code <= 67) return '🌧️'
  if (code <= 77) return '❄️'
  return '🌦️'
}

const loadStats = async () => {
  try {
    // 加载 TODO 概览
    const todos = await api.get('/todos').catch(() => [])
    openTodosCount.value = todos.filter(t => t.status !== 'completed').length
    todayTodos.value = todos.slice(0, 3).map(t => ({ ...t, completed: t.status === 'completed' }))

    // 加载笔记数量
    const notes = await api.get('/notes').catch(() => [])
    notesCount.value = Array.isArray(notes) ? notes.length : 0

    // 加载文件数量
    const files = await api.get('/files').catch(() => [])
    filesCount.value = Array.isArray(files) ? files.length : 0

    // 健康检查
    await api.get('/health')
    isBackendOnline.value = true
  } catch (e) {
    isBackendOnline.value = false
  }
}

const petInteract = () => {
  isPetBouncing.value = true
  const messages = ['喵~', '蹭蹭', '开心!', 'Zzz...']
  petMessage.value = messages[Math.floor(Math.random() * messages.length)]
  setTimeout(() => { isPetBouncing.value = false; petMessage.value = '' }, 2000)
}

onMounted(() => {
  // 安全读取版本号，避免 config 未准备好导致的异常
  try {
    appVersion.value = config?.app?.version || '3.0.1'
  } catch {
    appVersion.value = '3.0.1'
  }
  loadStats()
  fetchRSS()
  fetchWeather()
})
</script>

<style scoped>
.dashboard {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  color: var(--text-primary);
  background-color: var(--bg-color);
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  padding: 20px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  background: var(--card-bg);
  border-radius: 16px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--primary-color);
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 8px;
  font-weight: 500;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
}

.stat-sub {
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-secondary);
}

.system-card {
  align-items: flex-start;
}

.system-row {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.status-pill {
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 500;
  border: 1px solid var(--border-color);
  transition: all 0.3s ease;
}

.status-pill.online {
  background: rgba(103, 194, 58, 0.1);
  color: var(--success-color);
  border-color: rgba(103, 194, 58, 0.3);
}

.status-pill.offline {
  background: rgba(245, 108, 108, 0.1);
  color: var(--danger-color);
  border-color: rgba(245, 108, 108, 0.3);
}

.quick-links-section {
  margin-bottom: 24px;
}

.desktop-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 20px;
  grid-auto-rows: 180px;
}

.widget-item {
  padding: 20px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: var(--card-bg);
  border-radius: 16px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: default;
  position: relative;
  overflow: hidden;
}

.widget-item:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--primary-color);
}

.wide-widget {
  grid-column: span 2;
  align-items: flex-start;
  justify-content: flex-start;
}

/* 天气 */
.weather-content {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.weather-icon {
  font-size: 56px;
  line-height: 1;
}

.temperature {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
}

.location {
  font-size: 14px;
  color: var(--text-secondary);
}

/* 3D 宠物 */
.pet-widget {
  cursor: pointer;
  perspective: 1000px;
}

.pet-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  position: relative;
  transform-style: preserve-3d;
}

.pet-3d-model {
  position: relative;
  width: 100px;
  height: 100px;
  transform-style: preserve-3d;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  animation: rotate 6s infinite linear;
}

.pet-figure {
  position: absolute;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 72px;
  transform-style: preserve-3d;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.pet-3d-model:hover {
  transform: scale(1.1) rotateY(180deg);
}

.pet-3d-model.bounce {
  animation: bounce3d 0.6s ease-in-out;
}

.pet-message {
  background: var(--primary-color);
  color: white;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  animation: fadeInUp 0.3s ease-in;
  box-shadow: var(--shadow-md);
  transform: translateZ(10px);
}

/* 3D 动画 */
@keyframes rotate {
  from {
    transform: rotateY(0deg);
  }
  to {
    transform: rotateY(360deg);
  }
}

@keyframes bounce3d {
  0%, 100% {
    transform: translateY(0) rotateY(0deg);
  }
  50% {
    transform: translateY(-20px) rotateY(180deg) scale(1.3);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px) translateZ(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0) translateZ(10px);
  }
}

/* 待办 */
.widget-header {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  font-weight: 600;
  font-size: 16px;
  color: var(--text-primary);
}

.todo-list-mini {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.mini-todo-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.2s ease;
}

.mini-todo-item:last-child {
  border-bottom: none;
}

.todo-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--danger-color);
  flex-shrink: 0;
  transition: all 0.3s ease;
}

.todo-dot.completed {
  background: var(--success-color);
}

.todo-title {
  font-size: 14px;
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.todo-title.completed {
  text-decoration: line-through;
  color: var(--text-secondary);
}

.empty-mini {
  font-size: 14px;
  color: var(--text-secondary);
  text-align: center;
  padding: 16px 0;
  background: var(--bg-light);
  border-radius: 8px;
}

/* 动画 */
@keyframes bounce {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.3);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .dashboard {
    padding: 16px;
  }
  
  .stats-row {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .desktop-grid {
    grid-template-columns: 1fr;
    gap: 16px;
    grid-auto-rows: auto;
  }
  
  .wide-widget {
    grid-column: span 1;
  }
  
  .pet-emoji {
    font-size: 56px;
  }
  
  .temperature {
    font-size: 24px;
  }
}
</style>
