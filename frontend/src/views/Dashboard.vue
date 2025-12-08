<template>
  <div class="dashboard">
    <!-- 欢迎卡片 -->
    <el-card class="welcome-card">
      <div class="welcome-content">
        <div class="welcome-text">
          <h1>欢迎回来! 👋</h1>
          <p>{{ greetingMessage }}</p>
          <div class="datetime">
            <span>{{ currentDate }}</span>
            <span class="time">{{ currentTime }}</span>
          </div>
        </div>
        <div class="welcome-stats">
          <el-progress type="dashboard" :percentage="completionRate" :width="120">
            <template #default="{ percentage }">
              <span class="percentage-value">{{ percentage }}%</span>
              <span class="percentage-label">完成度</span>
            </template>
          </el-progress>
        </div>
      </div>
    </el-card>

    <!-- 统计卡片 - 拟物风格 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="6" v-for="(stat, index) in statCards" :key="index">
        <el-card class="stat-card neumorphism" @click="$router.push(stat.route)">
          <div class="stat-content">
            <div class="stat-icon-neu">
              <component :is="stat.icon" style="font-size: 40px;" :style="{ color: stat.color }" />
            </div>
            <div class="stat-info">
              <h3>{{ stat.value }}</h3>
              <p>{{ stat.label }}</p>
              <span class="stat-trend" :class="stat.trend > 0 ? 'up' : 'down'">
                {{ stat.trend > 0 ? '↑' : '↓' }} {{ Math.abs(stat.trend) }}%
              </span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 今日待办 -->
      <el-col :span="8">
        <el-card class="content-card">
          <template #header>
            <div class="card-header">
              <span>📝 今日待办</span>
              <el-button size="small" text @click="$router.push('/todos')">查看全部</el-button>
            </div>
          </template>
          <div v-if="todayTodos.length === 0" class="empty-state">
            <el-empty description="今天没有待办事项" :image-size="100" />
          </div>
          <div v-else class="todo-list">
            <div v-for="todo in todayTodos" :key="todo.id" class="todo-item">
              <el-checkbox v-model="todo.completed" @change="toggleTodo(todo)" />
              <span :class="{ completed: todo.completed }">{{ todo.title }}</span>
              <el-tag :type="getPriorityType(todo.priority)" size="small">
                {{ todo.priority }}
              </el-tag>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 最近文件 -->
      <el-col :span="8">
        <el-card class="content-card">
          <template #header>
            <div class="card-header">
              <span>📁 最近文件</span>
              <el-button size="small" text @click="$router.push('/files')">查看全部</el-button>
            </div>
          </template>
          <div v-if="recentFiles.length === 0" class="empty-state">
            <el-empty description="暂无文件" :image-size="100" />
          </div>
          <div v-else class="file-list">
            <div v-for="file in recentFiles" :key="file.id" class="file-item" @click="openFile(file)">
              <div class="file-icon" :style="{ background: getFileColor(file.extension) }">
                {{ file.extension || '?' }}
              </div>
              <div class="file-info">
                <span class="file-name">{{ file.file_name }}</span>
                <span class="file-meta">{{ formatSize(file.file_size) }} · {{ formatTime(file.created_at) }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 快速操作 -->
      <el-col :span="8">
        <el-card class="content-card">
          <template #header>
            <span>⚡ 快速操作</span>
          </template>
          <div class="quick-actions">
            <div class="action-btn" @click="$router.push('/files')">
              <el-icon size="24"><el-icon-upload /></el-icon>
              <span>上传文件</span>
            </div>
            <div class="action-btn" @click="$router.push('/code')">
              <el-icon size="24"><el-icon-document-add /></el-icon>
              <span>新建代码</span>
            </div>
            <div class="action-btn" @click="$router.push('/notes')">
              <el-icon size="24"><el-icon-edit /></el-icon>
              <span>写笔记</span>
            </div>
            <div class="action-btn" @click="$router.push('/chat')">
              <el-icon size="24"><el-icon-chat-line-round /></el-icon>
              <span>AI 对话</span>
            </div>
            <div class="action-btn" @click="$router.push('/calculator')">
              <el-icon size="24"><el-icon-promotion /></el-icon>
              <span>计算器</span>
            </div>
            <div class="action-btn" @click="$router.push('/data')">
              <el-icon size="24"><el-icon-data-analysis /></el-icon>
              <span>数据分析</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 活动时间线 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>📊 数据概览</span>
          </template>
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="overview-item">
                <h4>存储使用</h4>
                <el-progress :percentage="storageUsage" :color="storageColor" />
                <p>{{ usedStorage }} / {{ totalStorage }}</p>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="overview-item">
                <h4>本周活跃度</h4>
                <el-progress :percentage="weekActivity" color="#67c23a" />
                <p>{{ weekActivity }}% 比上周提升</p>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="overview-item">
                <h4>AI 使用量</h4>
                <el-progress :percentage="aiUsage" color="#409eff" />
                <p>本月已使用 {{ aiUsage }}%</p>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Folder, Document, DocumentCopy, List } from '@element-plus/icons-vue'
import api from '../api'

const currentDate = ref('')
const currentTime = ref('')
const stats = ref({
  fileCount: 0,
  noteCount: 0,
  codeCount: 0,
  chatCount: 0,
  todoCount: 0,
  completedTodos: 0
})

const recentFiles = ref([])
const todayTodos = ref([])

const greetingMessage = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了,注意休息'
  if (hour < 12) return '早上好,新的一天开始了'
  if (hour < 14) return '中午好,记得吃午饭哦'
  if (hour < 18) return '下午好,继续加油'
  if (hour < 22) return '晚上好,今天辛苦了'
  return '夜深了,早点休息吧'
})

const completionRate = computed(() => {
  if (stats.value.todoCount === 0) return 100
  return Math.round((stats.value.completedTodos / stats.value.todoCount) * 100)
})

const statCards = computed(() => [
  {
    icon: Folder,
    label: '文件总数',
    value: stats.value.fileCount,
    color: '#667eea',
    route: '/files',
    trend: 12
  },
  {
    icon: Document,
    label: '笔记',
    value: stats.value.noteCount,
    color: '#f5576c',
    route: '/notes',
    trend: 8
  },
  {
    icon: DocumentCopy,
    label: '代码片段',
    value: stats.value.codeCount,
    color: '#00f2fe',
    route: '/code',
    trend: -3
  },
  {
    icon: List,
    label: 'TODO',
    value: stats.value.todoCount,
    color: '#fee140',
    route: '/todos',
    trend: 5
  }
])

const storageUsage = ref(45)
const weekActivity = ref(78)
const aiUsage = ref(32)

const usedStorage = computed(() => '2.3 GB')
const totalStorage = computed(() => '5 GB')

const storageColor = computed(() => {
  if (storageUsage.value > 80) return '#f56c6c'
  if (storageUsage.value > 60) return '#e6a23c'
  return '#67c23a'
})

const updateDateTime = () => {
  const now = new Date()
  currentDate.value = now.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
  currentTime.value = now.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const loadStats = async () => {
  try {
    const [files, notes, code, todos] = await Promise.all([
      api.get('/files'),
      api.get('/notes'),
      api.get('/code'),
      api.get('/todos').catch(() => [])
    ])

    stats.value.fileCount = files.length
    stats.value.noteCount = notes.length
    stats.value.codeCount = code.length
    stats.value.todoCount = todos.length
    stats.value.completedTodos = todos.filter(t => t.status === 'completed').length

    recentFiles.value = files.slice(0, 5)

    const today = new Date().toISOString().split('T')[0]
    todayTodos.value = todos.filter(t => {
      return t.due_date && t.due_date.startsWith(today)
    }).slice(0, 5).map(t => ({
      ...t,
      completed: t.status === 'completed'
    }))
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

const toggleTodo = async (todo) => {
  try {
    const newStatus = todo.completed ? 'completed' : 'pending'
    await api.put(`/todos/${todo.id}`, { ...todo, status: newStatus })
    ElMessage.success('状态已更新')
    loadStats()
  } catch (error) {
    ElMessage.error('更新失败')
  }
}

const getPriorityType = (priority) => {
  const types = { high: 'danger', medium: 'warning', low: 'info' }
  return types[priority] || 'info'
}

const getFileColor = (ext) => {
  const colors = {
    pdf: '#f56c6c',
    doc: '#409eff',
    docx: '#409eff',
    xls: '#67c23a',
    xlsx: '#67c23a',
    zip: '#e6a23c',
    rar: '#e6a23c',
    jpg: '#f56c6c',
    png: '#f56c6c',
    mp4: '#909399'
  }
  return colors[ext?.toLowerCase()] || '#909399'
}

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const formatTime = (time) => {
  const now = new Date()
  const date = new Date(time)
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  return date.toLocaleDateString('zh-CN')
}

const openFile = (file) => {
  ElMessage.info('文件预览功能开发中...')
}

let timer = null

onMounted(() => {
  updateDateTime()
  timer = setInterval(updateDateTime, 1000)
  loadStats()
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.dashboard {
  animation: fadeIn 0.5s ease;
}

.welcome-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  color: white;
}

.welcome-card :deep(.el-card__body) {
  padding: 30px;
}

.welcome-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.welcome-text h1 {
  margin: 0 0 10px 0;
  font-size: 32px;
  font-weight: bold;
}

.welcome-text p {
  margin: 0 0 15px 0;
  font-size: 18px;
  opacity: 0.9;
}

.datetime {
  display: flex;
  gap: 15px;
  font-size: 14px;
  opacity: 0.8;
}

.time {
  font-weight: bold;
}

.percentage-value {
  display: block;
  font-size: 24px;
  font-weight: bold;
}

.percentage-label {
  display: block;
  font-size: 12px;
  margin-top: 5px;
  opacity: 0.8;
}

.stat-card {
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
  overflow: hidden;
}

/* 拟物风格卡片 */
.neumorphism {
  background: var(--card-bg);
  box-shadow: 8px 8px 16px rgba(0, 0, 0, 0.1),
             -8px -8px 16px rgba(255, 255, 255, 0.05);
  border-radius: 16px;
}

.neumorphism:hover {
  transform: translateY(-3px);
  box-shadow: 12px 12px 24px rgba(0, 0, 0, 0.15),
             -12px -12px 24px rgba(255, 255, 255, 0.05);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

/* 拟物图标容器 */
.stat-icon-neu {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--card-bg);
  box-shadow: inset 4px 4px 8px rgba(0, 0, 0, 0.1),
              inset -4px -4px 8px rgba(255, 255, 255, 0.05);
  flex-shrink: 0;
}

.stat-info {
  flex: 1;
}

.stat-info h3 {
  margin: 0 0 5px 0;
  font-size: 32px;
  font-weight: bold;
  color: var(--text-primary);
}

.stat-info p {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: var(--text-secondary);
  font-weight: 500;
}

.stat-trend {
  font-size: 12px;
  font-weight: bold;
}

.stat-trend.up {
  color: #67c23a;
}

.stat-trend.down {
  color: #f56c6c;
}

.content-card {
  height: 400px;
  display: flex;
  flex-direction: column;
}

.content-card :deep(.el-card__body) {
  flex: 1;
  overflow-y: auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.todo-list, .file-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.todo-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px;
  background: var(--hover-bg);
  border-radius: 8px;
  transition: all 0.3s;
}

.todo-item:hover {
  background: var(--border-color);
  transform: translateX(5px);
}

.todo-item span {
  flex: 1;
  color: var(--text-primary);
}

.todo-item span.completed {
  text-decoration: line-through;
  opacity: 0.6;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--hover-bg);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.file-item:hover {
  background: var(--border-color);
  transform: translateX(5px);
}

.file-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
  font-size: 12px;
  text-transform: uppercase;
  flex-shrink: 0;
}

.file-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.file-name {
  color: var(--text-primary);
  font-weight: 500;
}

.file-meta {
  font-size: 12px;
  color: var(--text-secondary);
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 20px;
  background: linear-gradient(135deg, var(--hover-bg) 0%, var(--card-bg) 100%);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.action-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  border-color: #409eff;
}

.action-btn span {
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 500;
}

.overview-item {
  text-align: center;
  padding: 20px;
}

.overview-item h4 {
  margin: 0 0 15px 0;
  color: var(--text-primary);
  font-size: 16px;
}

.overview-item p {
  margin: 10px 0 0 0;
  font-size: 14px;
  color: var(--text-secondary);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
