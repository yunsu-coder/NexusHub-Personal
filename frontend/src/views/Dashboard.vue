<template>
  <div class="dashboard">
    <!-- 顶部统计总览 -->
    <div class="stats-row">
      <div class="stat-card glass-card">
        <div class="stat-label">笔记</div>
        <div class="stat-value">{{ notesCount }}</div>
        <div class="stat-sub">最近 3 篇已加载</div>
      </div>
      <div class="stat-card glass-card">
        <div class="stat-label">未完成待办</div>
        <div class="stat-value">{{ openTodosCount }}</div>
        <div class="stat-sub">今日聚焦任务</div>
      </div>
      <div class="stat-card glass-card">
        <div class="stat-label">文件</div>
        <div class="stat-value">{{ filesCount }}</div>
        <div class="stat-sub">已上传总数</div>
      </div>
      <div class="stat-card glass-card system-card">
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
      <div class="widget-item weather-widget glass-card">
        <div class="weather-content" v-if="weatherData">
          <div class="weather-icon">{{ getWeatherIcon(weatherData.weathercode) }}</div>
          <div class="weather-info">
            <div class="temperature">{{ weatherData.temperature }}°C</div>
            <div class="location">上海 (实时)</div>
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

      <!-- 宠物卡片 -->
      <div class="widget-item pet-widget glass-card" @click="petInteract">
        <div class="pet-emoji" :class="{ bounce: isPetBouncing }">{{ petEmoji }}</div>
        <div class="pet-message" v-if="petMessage">{{ petMessage }}</div>
      </div>

      <!-- RSS 新闻卡片 -->
      <div class="widget-item rss-widget glass-card">
        <div class="widget-header">
          <span>📰 科技动态</span>
          <el-button link size="small" @click="fetchRSS">刷新</el-button>
        </div>
        <div class="rss-list" v-if="rssFeed.length > 0">
          <a 
            v-for="(item, index) in rssFeed" 
            :key="index" 
            :href="item.link" 
            target="_blank" 
            class="rss-item"
          >
            <div class="rss-dot"></div>
            <div class="rss-title">{{ item.title }}</div>
          </a>
        </div>
        <div class="rss-loading" v-else>
          加载中...
        </div>
      </div>

      <!-- 待办事项 (列表样式) -->
      <div class="widget-item todo-widget glass-card wide-widget">
        <div class="widget-header">
          <span>📝 待办事项</span>
          <el-button link size="small" @click="$router.push('/todos')">更多</el-button>
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

// ... (保留原有逻辑，简化部分)
const todayTodos = ref([])
const notesCount = ref(0)
const openTodosCount = ref(0)
const filesCount = ref(0)
const isBackendOnline = ref(true)
const appVersion = ref('3.0.1')
const petEmoji = ref('🐱')
const petMessage = ref('')
const isPetBouncing = ref(false)
const rssFeed = ref([])
const weatherData = ref(null)

const fetchWeather = async () => {
  try {
    // 上海坐标: 31.22, 121.46
    const res = await fetch('https://api.open-meteo.com/v1/forecast?latitude=31.22&longitude=121.46&current_weather=true')
    const data = await res.json()
    if (data.current_weather) {
      weatherData.value = data.current_weather
    }
  } catch (e) {
    console.error('Weather fetch failed', e)
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
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.stat-label {
  font-size: 12px;
  opacity: 0.7;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 22px;
  font-weight: 700;
}

.stat-sub {
  margin-top: 4px;
  font-size: 12px;
  opacity: 0.7;
}

.system-card {
  align-items: flex-start;
}

.system-row {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-pill {
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 11px;
  border: 1px solid rgba(255,255,255,0.3);
}

.status-pill.online {
  background: rgba(103, 194, 58, 0.2);
}

.status-pill.offline {
  background: rgba(245, 108, 108, 0.2);
}

.quick-links-section {
  margin-bottom: 30px;
}

.desktop-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  grid-auto-rows: 180px;
}

.widget-item {
  padding: 20px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  transition: transform 0.2s;
  cursor: default;
  position: relative;
  overflow: hidden;
}

.widget-item:hover {
  transform: translateY(-5px);
  background: rgba(255,255,255,0.15);
}

.wide-widget {
  grid-column: span 2;
  align-items: flex-start;
  justify-content: flex-start;
}

/* 天气 */
.weather-content { text-align: center; }
.weather-icon { font-size: 40px; margin-bottom: 10px; }
.temperature { font-size: 24px; font-weight: bold; }
.location { font-size: 12px; opacity: 0.8; }

/* 宠物 */
.pet-widget { cursor: pointer; }
.pet-emoji { font-size: 50px; transition: transform 0.2s; }
.pet-emoji.bounce { transform: scale(1.2); }
.pet-message { 
  position: absolute; top: 10px; right: 10px; 
  background: #fff; color: #333; padding: 4px 8px; 
  border-radius: 8px; font-size: 12px; 
}

/* 名言/RSS */
.rss-widget { cursor: default; }
.rss-list { width: 100%; display: flex; flex-direction: column; gap: 8px; }
.rss-item { display: flex; align-items: center; gap: 8px; text-decoration: none; color: #fff; opacity: 0.9; transition: opacity 0.2s; }
.rss-item:hover { opacity: 1; }
.rss-dot { width: 6px; height: 6px; background: #e6a23c; border-radius: 50%; flex-shrink: 0; }
.rss-title { font-size: 12px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.rss-loading { font-size: 12px; opacity: 0.6; }

/* 待办 */
.widget-header {
  width: 100%;
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  font-weight: bold;
}
.todo-list-mini { width: 100%; }
.mini-todo-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 5px 0;
  border-bottom: 1px solid rgba(255,255,255,0.1);
}
.todo-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #f56c6c;
}
.todo-dot.completed { background: #67c23a; }
.todo-title { font-size: 13px; flex: 1; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.todo-title.completed { text-decoration: line-through; opacity: 0.6; }
.empty-mini { font-size: 12px; color: #aaa; text-align: center; margin-top: 20px; }

@media (max-width: 768px) {
  .desktop-grid { grid-template-columns: 1fr; }
  .wide-widget { grid-column: span 1; }
}
</style>
