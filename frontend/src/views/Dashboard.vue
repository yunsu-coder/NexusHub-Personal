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
          <el-skeleton v-if="isLoading" rows="4" animated />
          <div v-else-if="todayTodos.length === 0" class="empty-state">
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
          <el-skeleton v-if="isLoading" rows="4" animated />
          <div v-else-if="recentFiles.length === 0" class="empty-state">
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
const isLoading = ref(true)

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
    const [files, notes, todos] = await Promise.all([
      api.get('/files'),
      api.get('/notes'),
      api.get('/todos').catch(() => [])
    ])

    stats.value.fileCount = files.length
    stats.value.noteCount = notes.length
    stats.value.codeCount = 0 // 暂时设置为0，等待后端实现code接口
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
  finally {
    isLoading.value = false
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
  padding: 24px;
  animation: fadeIn 0.5s ease;
}

/* 欢迎卡片 */
.welcome-card {
  background: var(--gradient-primary);
  border: 1px solid var(--border-color);
  color: white;
  backdrop-filter: blur(15px);
  border-radius: 16px;
  box-shadow: var(--shadow-lg);
  margin-bottom: 24px;
  animation: slideInUp 0.5s ease;
  position: relative;
  overflow: hidden;
}

.welcome-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -25%;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.2) 0%, transparent 70%);
  border-radius: 50%;
  animation: pulse 3s infinite;
}

.welcome-card :deep(.el-card__body) {
  padding: 36px;
  position: relative;
  z-index: 1;
}

.welcome-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 30px;
}

.welcome-text h1 {
  margin: 0 0 16px 0;
  font-size: 36px;
  font-weight: 700;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  animation: slideInLeft 0.5s ease 0.1s both;
}

.welcome-text p {
  margin: 0 0 20px 0;
  font-size: 20px;
  opacity: 0.95;
  font-weight: 500;
  animation: slideInLeft 0.5s ease 0.2s both;
}

.datetime {
  display: flex;
  gap: 20px;
  font-size: 16px;
  opacity: 0.9;
  animation: slideInLeft 0.5s ease 0.3s both;
}

.time {
  font-weight: 700;
  font-size: 18px;
}

.welcome-stats {
  animation: slideInRight 0.5s ease 0.4s both;
}

.percentage-value {
  display: block;
  font-size: 28px;
  font-weight: 700;
}

.percentage-label {
  display: block;
  font-size: 14px;
  margin-top: 5px;
  opacity: 0.9;
}

/* 统计卡片 */
.stat-card {
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: none;
  overflow: hidden;
  margin-bottom: 24px;
  position: relative;
}

/* 拟物风格卡片 */
.neumorphism {
  background: var(--gradient-card);
  box-shadow: var(--shadow-md);
  border-radius: 16px;
  backdrop-filter: blur(15px);
  border: 1px solid var(--border-color);
  animation: slideInUp 0.5s ease calc(var(--index) * 0.1s + 0.2s) both;
}

.neumorphism:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: var(--shadow-xl);
  border-color: rgba(102, 126, 234, 0.3);
}

.stat-card :deep(.el-card__body) {
  padding: 24px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

/* 拟物图标容器 */
.stat-icon-neu {
  width: 85px;
  height: 85px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: var(--card-bg);
  box-shadow: var(--shadow-sm), inset var(--shadow-inset);
  flex-shrink: 0;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.stat-icon-neu::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: var(--gradient-primary);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.stat-card:hover .stat-icon-neu {
  box-shadow: var(--shadow-md), inset var(--shadow-inset-lg);
}

.stat-card:hover .stat-icon-neu::before {
  opacity: 0.1;
}

.stat-info {
  flex: 1;
}

.stat-info h3 {
  margin: 0 0 8px 0;
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -1px;
}

.stat-info p {
  margin: 0 0 10px 0;
  font-size: 15px;
  color: var(--text-secondary);
  font-weight: 500;
}

.stat-trend {
  font-size: 13px;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.stat-trend.up {
  color: #43e97b;
  background: rgba(67, 233, 123, 0.1);
}

.stat-trend.down {
  color: #f5576c;
  background: rgba(245, 87, 108, 0.1);
}

/* 内容卡片 */
.content-card {
  height: 400px;
  display: flex;
  flex-direction: column;
  background: var(--card-bg);
  backdrop-filter: blur(15px);
  border-radius: 16px;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
  animation: slideInUp 0.5s ease calc(var(--index) * 0.1s + 0.4s) both;
}

.content-card :deep(.el-card__body) {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.content-card :deep(.el-card__header) {
  padding: 20px 20px 16px;
  border-bottom: 1px solid var(--border-color);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: var(--text-primary);
  font-weight: 600;
  font-size: 16px;
}

.card-header :deep(.el-button) {
  color: var(--primary-color);
  font-weight: 500;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.empty-state :deep(.el-empty__description) {
  color: var(--text-secondary);
}

/* 待办事项 */
.todo-list, .file-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.todo-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  background: var(--sidebar-hover);
  border-radius: 10px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid transparent;
}

.todo-item:hover {
  background: rgba(102, 126, 234, 0.1);
  transform: translateX(8px);
  border-color: var(--border-color);
  box-shadow: var(--shadow-sm);
}

.todo-item span {
  flex: 1;
  color: var(--text-primary);
  font-weight: 500;
}

.todo-item span.completed {
  text-decoration: line-through;
  opacity: 0.6;
}

.todo-item :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
}

.todo-item :deep(.el-tag) {
  margin-left: auto;
}

/* 文件列表 */
.file-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px;
  background: var(--sidebar-hover);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid transparent;
}

.file-item:hover {
  background: rgba(102, 126, 234, 0.1);
  transform: translateX(8px);
  border-color: var(--border-color);
  box-shadow: var(--shadow-sm);
}

.file-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 14px;
  text-transform: uppercase;
  flex-shrink: 0;
  box-shadow: var(--shadow-sm);
  background: var(--gradient-primary);
  transition: transform 0.3s ease;
}

.file-item:hover .file-icon {
  transform: scale(1.1);
}

.file-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.file-name {
  color: var(--text-primary);
  font-weight: 600;
  font-size: 15px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-meta {
  font-size: 13px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 10px;
}

/* 快速操作 */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 24px 16px;
  background: linear-gradient(135deg, var(--sidebar-hover) 0%, var(--card-bg) 100%);
  border: 1px solid var(--border-color);
  border-radius: 14px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.action-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--gradient-primary);
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 0;
}

.action-btn:hover {
  transform: translateY(-8px) scale(1.05);
  box-shadow: var(--shadow-lg);
  border-color: var(--primary-color);
}

.action-btn:hover::before {
  opacity: 0.1;
}

.action-btn :deep(.el-icon) {
  font-size: 36px;
  color: var(--primary-color);
  transition: all 0.3s ease;
  position: relative;
  z-index: 1;
}

.action-btn:hover :deep(.el-icon) {
  color: var(--primary-color);
  transform: scale(1.15);
}

.action-btn span {
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 600;
  transition: all 0.3s ease;
  position: relative;
  z-index: 1;
}

.action-btn:hover span {
  color: var(--primary-color);
}

/* 数据概览 */
.overview-item {
  text-align: center;
  padding: 24px;
}

.overview-item h4 {
  margin: 0 0 20px 0;
  color: var(--text-primary);
  font-size: 16px;
  font-weight: 600;
}

.overview-item p {
  margin: 12px 0 0 0;
  font-size: 14px;
  color: var(--text-secondary);
  font-weight: 500;
}

.overview-item :deep(.el-progress) {
  margin-bottom: 8px;
}

/* 动画 */
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

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(40px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideInLeft {
  from {
    opacity: 0;
    transform: translateX(-30px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes slideInRight {
  from {
    opacity: 0;
    transform: translateX(30px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 0.5;
  }
  50% {
    transform: scale(1.2);
    opacity: 0.8;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .dashboard {
    padding: 16px;
  }

  .welcome-content {
    flex-direction: column;
    text-align: center;
  }

  .welcome-text h1 {
    font-size: 28px;
  }

  .welcome-text p {
    font-size: 16px;
  }

  .welcome-card :deep(.el-card__body) {
    padding: 28px 20px;
  }

  .quick-actions {
    grid-template-columns: 1fr;
  }

  .stat-content {
    flex-direction: column;
    text-align: center;
    gap: 16px;
  }
}

/* 平滑滚动 */
.content-card :deep(.el-card__body)::-webkit-scrollbar {
  width: 6px;
}

.content-card :deep(.el-card__body)::-webkit-scrollbar-track {
  background: var(--bg-color);
}

.content-card :deep(.el-card__body)::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

.content-card :deep(.el-card__body)::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}
</style>
