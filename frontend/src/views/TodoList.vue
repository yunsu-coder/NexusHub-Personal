<template>
  <div class="todo-page">
    <el-row :gutter="20">
      <!-- 左侧: TODO 列表 -->
      <el-col :span="16">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>📝 我的目标清单</span>
              <el-button type="primary" @click="showAddDialog = true">添加目标</el-button>
            </div>
          </template>

          <!-- 过滤器 -->
          <div class="filter-bar">
            <el-radio-group v-model="filterStatus" @change="filterTodos">
              <el-radio-button label="all">全部 ({{ todos.length }})</el-radio-button>
              <el-radio-button label="pending">待办 ({{ pendingCount }})</el-radio-button>
              <el-radio-button label="in_progress">进行中 ({{ inProgressCount }})</el-radio-button>
              <el-radio-button label="completed">已完成 ({{ completedCount }})</el-radio-button>
            </el-radio-group>

            <el-select v-model="filterPriority" placeholder="优先级" style="width: 120px; margin-left: 10px" @change="filterTodos">
              <el-option label="全部" value="all" />
              <el-option label="高" value="high" />
              <el-option label="中" value="medium" />
              <el-option label="低" value="low" />
            </el-select>
          </div>

          <!-- TODO 列表 -->
          <div class="todo-list">
            <div
              v-for="todo in filteredTodos"
              :key="todo.id"
              class="todo-item"
              :class="{ completed: todo.status === 'completed' }"
            >
              <div class="todo-checkbox">
                <el-checkbox
                  :model-value="todo.status === 'completed'"
                  @change="toggleComplete(todo)"
                  size="large"
                />
              </div>

              <div class="todo-content" @click="editTodo(todo)">
                <div class="todo-title">
                  <span :class="{ 'line-through': todo.status === 'completed' }">{{ todo.title }}</span>
                  <el-tag
                    :type="getPriorityType(todo.priority)"
                    size="small"
                    style="margin-left: 10px"
                  >
                    {{ getPriorityText(todo.priority) }}
                  </el-tag>
                  <el-tag
                    :type="getStatusType(todo.status)"
                    size="small"
                    style="margin-left: 5px"
                  >
                    {{ getStatusText(todo.status) }}
                  </el-tag>
                </div>

                <div class="todo-description" v-if="todo.description">
                  {{ todo.description }}
                </div>

                <div class="todo-meta">
                  <span v-if="todo.category">📁 {{ todo.category }}</span>
                  <span v-if="todo.due_date">📅 {{ formatDate(todo.due_date) }}</span>
                  <span>🕐 {{ formatDate(todo.created_at) }}</span>
                </div>
              </div>

              <div class="todo-actions">
                <el-button-group>
                  <el-button
                    size="small"
                    @click="changeStatus(todo, 'pending')"
                    :disabled="todo.status === 'pending'"
                  >
                    待办
                  </el-button>
                  <el-button
                    size="small"
                    @click="changeStatus(todo, 'in_progress')"
                    :disabled="todo.status === 'in_progress'"
                  >
                    进行中
                  </el-button>
                  <el-button
                    size="small"
                    type="danger"
                    @click="deleteTodo(todo)"
                  >
                    删除
                  </el-button>
                </el-button-group>
              </div>
            </div>

            <el-empty v-if="filteredTodos.length === 0" description="暂无目标" />
          </div>
        </el-card>
      </el-col>

      <!-- 右侧: 统计和进度 -->
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>📊 完成统计</span>
          </template>

          <div class="progress-section">
            <el-progress
              :percentage="completionRate"
              :color="progressColor"
              :stroke-width="20"
            >
              <span style="font-size: 16px; font-weight: bold">{{ completionRate }}%</span>
            </el-progress>

            <el-descriptions :column="1" border style="margin-top: 20px">
              <el-descriptions-item label="总目标">{{ todos.length }}</el-descriptions-item>
              <el-descriptions-item label="待办">{{ pendingCount }}</el-descriptions-item>
              <el-descriptions-item label="进行中">{{ inProgressCount }}</el-descriptions-item>
              <el-descriptions-item label="已完成">{{ completedCount }}</el-descriptions-item>
            </el-descriptions>
          </div>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>
            <span>🏆 今日目标</span>
          </template>

          <div class="today-goals">
            <div v-for="todo in todayTodos" :key="todo.id" class="today-item">
              <el-checkbox
                :model-value="todo.status === 'completed'"
                @change="toggleComplete(todo)"
              />
              <span style="margin-left: 10px">{{ todo.title }}</span>
            </div>

            <el-empty v-if="todayTodos.length === 0" description="今天没有目标" />
          </div>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>
            <span>📅 本周概览</span>
          </template>

          <div class="week-overview">
            <el-statistic title="本周新增" :value="weekStats.added" />
            <el-statistic title="本周完成" :value="weekStats.completed" style="margin-top: 15px" />
            <el-statistic
              title="完成率"
              :value="weekStats.rate"
              suffix="%"
              style="margin-top: 15px"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="currentTodo?.id ? '编辑目标' : '添加目标'"
      width="500px"
    >
      <el-form :model="currentTodo" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="currentTodo.title" placeholder="输入目标标题..." />
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="currentTodo.description"
            type="textarea"
            :rows="3"
            placeholder="详细描述..."
          />
        </el-form-item>

        <el-form-item label="优先级">
          <el-select v-model="currentTodo.priority" style="width: 100%">
            <el-option label="高" value="high" />
            <el-option label="中" value="medium" />
            <el-option label="低" value="low" />
          </el-select>
        </el-form-item>

        <el-form-item label="分类">
          <el-input v-model="currentTodo.category" placeholder="工作/学习/生活..." />
        </el-form-item>

        <el-form-item label="截止日期">
          <el-date-picker
            v-model="currentTodo.due_date"
            type="date"
            placeholder="选择日期"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveTodo">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

const todos = ref([])
const filteredTodos = ref([])
const filterStatus = ref('all')
const filterPriority = ref('all')
const showAddDialog = ref(false)
const currentTodo = ref({
  title: '',
  description: '',
  priority: 'medium',
  status: 'pending',
  category: '',
  due_date: null
})

// 统计
const pendingCount = computed(() => todos.value.filter(t => t.status === 'pending').length)
const inProgressCount = computed(() => todos.value.filter(t => t.status === 'in_progress').length)
const completedCount = computed(() => todos.value.filter(t => t.status === 'completed').length)
const completionRate = computed(() => {
  if (todos.value.length === 0) return 0
  return Math.round((completedCount.value / todos.value.length) * 100)
})

const progressColor = computed(() => {
  if (completionRate.value < 30) return '#f56c6c'
  if (completionRate.value < 70) return '#e6a23c'
  return '#67c23a'
})

const todayTodos = computed(() => {
  const today = new Date().toISOString().split('T')[0]
  return todos.value.filter(t => {
    return t.due_date && t.due_date.startsWith(today) && t.status !== 'completed'
  })
})

const weekStats = computed(() => {
  const now = new Date()
  const weekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)

  const added = todos.value.filter(t => new Date(t.created_at) > weekAgo).length
  const completed = todos.value.filter(t => {
    return t.status === 'completed' && t.updated_at && new Date(t.updated_at) > weekAgo
  }).length

  const rate = added > 0 ? Math.round((completed / added) * 100) : 0

  return { added, completed, rate }
})

const loadTodos = async () => {
  try {
    const response = await api.get('/todos')
    todos.value = response
    filterTodos()
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const filterTodos = () => {
  let filtered = todos.value

  if (filterStatus.value !== 'all') {
    filtered = filtered.filter(t => t.status === filterStatus.value)
  }

  if (filterPriority.value !== 'all') {
    filtered = filtered.filter(t => t.priority === filterPriority.value)
  }

  filteredTodos.value = filtered
}

const saveTodo = async () => {
  try {
    if (!currentTodo.value.title) {
      ElMessage.warning('请输入标题')
      return
    }

    if (currentTodo.value.id) {
      await api.put(`/todos/${currentTodo.value.id}`, currentTodo.value)
      ElMessage.success('更新成功')
    } else {
      await api.post('/todos', currentTodo.value)
      ElMessage.success('添加成功')
    }

    showAddDialog.value = false
    currentTodo.value = {
      title: '',
      description: '',
      priority: 'medium',
      status: 'pending',
      category: '',
      due_date: null
    }
    loadTodos()
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

const editTodo = (todo) => {
  currentTodo.value = { ...todo }
  showAddDialog.value = true
}

const toggleComplete = async (todo) => {
  const newStatus = todo.status === 'completed' ? 'pending' : 'completed'
  await changeStatus(todo, newStatus)
}

const changeStatus = async (todo, status) => {
  try {
    await api.put(`/todos/${todo.id}`, { ...todo, status })
    ElMessage.success('状态已更新')
    loadTodos()
  } catch (error) {
    ElMessage.error('更新失败')
  }
}

const deleteTodo = async (todo) => {
  try {
    await ElMessageBox.confirm('确定删除这个目标吗?', '提示', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await api.delete(`/todos/${todo.id}`)
    ElMessage.success('删除成功')
    loadTodos()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const getPriorityType = (priority) => {
  const types = { high: 'danger', medium: 'warning', low: 'info' }
  return types[priority] || 'info'
}

const getPriorityText = (priority) => {
  const texts = { high: '高', medium: '中', low: '低' }
  return texts[priority] || priority
}

const getStatusType = (status) => {
  const types = { pending: 'info', in_progress: 'warning', completed: 'success' }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = { pending: '待办', in_progress: '进行中', completed: '已完成' }
  return texts[status] || status
}

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN')
}

onMounted(loadTodos)
</script>

<style scoped>
.todo-page {
  animation: fadeIn 0.5s;
}

.filter-bar {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}

.todo-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.todo-item {
  display: flex;
  align-items: flex-start;
  gap: 15px;
  padding: 15px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  transition: all 0.3s;
}

.todo-item:hover {
  border-color: #409EFF;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
}

.todo-item.completed {
  opacity: 0.6;
}

.todo-checkbox {
  flex-shrink: 0;
  padding-top: 2px;
}

.todo-content {
  flex: 1;
  cursor: pointer;
}

.todo-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 8px;
  display: flex;
  align-items: center;
}

.todo-title .line-through {
  text-decoration: line-through;
}

.todo-description {
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 8px;
  line-height: 1.5;
}

.todo-meta {
  display: flex;
  gap: 15px;
  font-size: 12px;
  color: var(--text-secondary);
}

.todo-actions {
  flex-shrink: 0;
}

.progress-section {
  padding: 10px 0;
}

.today-goals {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.today-item {
  display: flex;
  align-items: center;
  padding: 10px;
  background: var(--hover-bg);
  border-radius: 6px;
}

.week-overview {
  padding: 10px 0;
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
