<template>
  <div class="todo-page">
    <div class="todo-header glass-panel">
      <div class="header-left">
        <h2>任务看板</h2>
        <div class="view-toggles">
           <el-radio-group v-model="viewMode" size="small">
             <el-radio-button label="kanban"><el-icon><DataBoard /></el-icon> 看板</el-radio-button>
             <el-radio-button label="list"><el-icon><List /></el-icon> 列表</el-radio-button>
           </el-radio-group>
        </div>
      </div>
      <div class="header-right">
        <el-select v-model="statusFilter" size="small" class="filter-select" placeholder="状态">
          <el-option label="全部" value="all" />
          <el-option label="待办" value="pending" />
          <el-option label="进行中" value="in_progress" />
          <el-option label="已完成" value="completed" />
        </el-select>
        <el-select v-model="sortKey" size="small" class="filter-select" placeholder="排序">
          <el-option label="按创建时间" value="created_at" />
          <el-option label="按截止时间" value="due_date" />
          <el-option label="按优先级" value="priority" />
        </el-select>
        <el-input 
           v-model="searchKeyword" 
           placeholder="搜索任务..." 
           prefix-icon="Search"
           class="search-input"
           clearable
        />
        <el-button type="primary" :icon="Plus" @click="openAddDialog()">新建任务</el-button>
      </div>
    </div>

    <!-- Kanban View -->
    <div v-if="viewMode === 'kanban'" class="kanban-board" v-loading="loading">
      <div class="kanban-column pending-col">
        <div class="column-header">
          <span class="dot pending"></span>
          <h3>待办事项</h3>
          <span class="count">{{ pendingTodos.length }}</span>
        </div>
        <div class="column-content">
          <draggable 
            v-model="pendingTodos" 
            group="todos" 
            item-key="id"
            class="drag-area"
            @change="onChange($event, 'pending')"
          >
            <template #item="{ element }">
              <div class="todo-card glass-panel" @click="editTodo(element)">
                <div class="card-badges">
                  <el-tag size="small" :type="getPriorityType(element.priority)">{{ getPriorityText(element.priority) }}</el-tag>
                  <span class="time-ago">{{ formatTimeAgo(element.created_at) }}</span>
                </div>
                <h4 class="card-title">{{ element.title }}</h4>
                <div class="card-footer">
                   <el-avatar :size="20" class="assignee-avatar">Me</el-avatar>
                   <el-button link type="danger" :icon="Delete" @click.stop="deleteTodo(element)" />
                </div>
              </div>
            </template>
          </draggable>
        </div>
      </div>

      <div class="kanban-column progress-col">
        <div class="column-header">
          <span class="dot progress"></span>
          <h3>进行中</h3>
          <span class="count">{{ inProgressTodos.length }}</span>
        </div>
        <div class="column-content">
          <draggable 
            v-model="inProgressTodos" 
            group="todos" 
            item-key="id"
            class="drag-area"
            @change="onChange($event, 'in_progress')"
          >
            <template #item="{ element }">
              <div class="todo-card glass-panel" @click="editTodo(element)">
                <div class="card-badges">
                  <el-tag size="small" :type="getPriorityType(element.priority)">{{ getPriorityText(element.priority) }}</el-tag>
                </div>
                <h4 class="card-title">{{ element.title }}</h4>
                <div class="card-footer">
                   <el-avatar :size="20" class="assignee-avatar" style="background: #e6a23c">Me</el-avatar>
                   <el-button link type="danger" :icon="Delete" @click.stop="deleteTodo(element)" />
                </div>
              </div>
            </template>
          </draggable>
        </div>
      </div>

      <div class="kanban-column completed-col">
        <div class="column-header">
          <span class="dot completed"></span>
          <h3>已完成</h3>
          <span class="count">{{ completedTodos.length }}</span>
        </div>
        <div class="column-content">
           <draggable 
            v-model="completedTodos" 
            group="todos" 
            item-key="id"
            class="drag-area"
            @change="onChange($event, 'completed')"
          >
            <template #item="{ element }">
              <div class="todo-card glass-panel completed" @click="editTodo(element)">
                <div class="card-badges">
                  <el-icon color="#67C23A"><Check /></el-icon>
                  <span class="done-time">{{ formatTimeAgo(element.updated_at) }}</span>
                </div>
                <h4 class="card-title">{{ element.title }}</h4>
                 <div class="card-footer">
                   <el-button link type="danger" :icon="Delete" @click.stop="deleteTodo(element)" />
                </div>
              </div>
            </template>
          </draggable>
        </div>
      </div>
    </div>

    <!-- List View (Legacy Table) -->
    <div v-else class="list-view glass-panel" v-loading="loading">
      <el-table :data="filteredTodos" style="width: 100%; background: transparent;">
        <el-table-column prop="status" label="状态" width="100">
           <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
           </template>
        </el-table-column>
        <el-table-column prop="title" label="任务标题" />
        <el-table-column prop="priority" label="优先级" width="100">
          <template #default="{ row }">
             <el-tag :type="getPriorityType(row.priority)" effect="plain">{{ getPriorityText(row.priority) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
           <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="right">
          <template #default="{ row }">
            <el-button size="small" :icon="Edit" circle @click="editTodo(row)" />
            <el-button size="small" type="danger" :icon="Delete" circle @click="deleteTodo(row)" />
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- Edit Dialog -->
    <el-dialog 
      v-model="showDialog" 
      :title="currentTodo.id ? '编辑任务' : '新建任务'" 
      width="500px"
    >
      <el-form :model="currentTodo" label-position="top">
        <el-form-item label="标题">
          <el-input v-model="currentTodo.title" placeholder="要做什么?" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="currentTodo.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
             <el-form-item label="优先级">
               <el-select v-model="currentTodo.priority" style="width: 100%">
                 <el-option label="高" value="high" />
                 <el-option label="中" value="medium" />
                 <el-option label="低" value="low" />
               </el-select>
             </el-form-item>
          </el-col>
          <el-col :span="12">
             <el-form-item label="截止时间">
                <el-date-picker v-model="currentTodo.due_date" type="datetime" placeholder="选择日期" style="width: 100%" />
             </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveTodo">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import draggable from 'vuedraggable'
import { Plus, Search, Delete, Edit, DataBoard, List, Check } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

const viewMode = ref('kanban')
const searchKeyword = ref('')
const statusFilter = ref('all')
const sortKey = ref('created_at')
const todos = ref([])
const showDialog = ref(false)
const currentTodo = ref({ title: '', priority: 'medium', status: 'pending' })
const loading = ref(false)

// Kanban Columns (Computed with setter logic handled by draggable change event)
const pendingTodos = ref([])
const inProgressTodos = ref([])
const completedTodos = ref([])

const loadTodos = async () => {
  loading.value = true
  try {
    const data = await api.get('/todos')
    todos.value = data || []
    distributeTodos()
  } catch (e) {
    ElMessage.error('加载任务失败')
  } finally {
    loading.value = false
  }
}

const distributeTodos = () => {
  pendingTodos.value = todos.value.filter(t => t.status === 'pending')
  inProgressTodos.value = todos.value.filter(t => t.status === 'in_progress')
  completedTodos.value = todos.value.filter(t => t.status === 'completed')
}

// Handle Drag & Drop Changes
const onChange = (evt, newStatus) => {
  if (evt.added) {
    const todo = evt.added.element
    updateStatus(todo, newStatus)
  }
}

const updateStatus = async (todo, status) => {
  try {
    await api.put(`/todos/${todo.id}`, { ...todo, status })
    // Update local master list
    const idx = todos.value.findIndex(t => t.id === todo.id)
    if (idx !== -1) todos.value[idx].status = status
  } catch (e) {
    ElMessage.error('状态更新失败')
    loadTodos() // Revert on error
  }
}

const filteredTodos = computed(() => {
  let list = [...todos.value]

  // 状态过滤
  if (statusFilter.value !== 'all') {
    list = list.filter(t => t.status === statusFilter.value)
  }

  // 关键词过滤
  if (searchKeyword.value) {
    const q = searchKeyword.value.toLowerCase()
    list = list.filter(t =>
      (t.title && t.title.toLowerCase().includes(q)) ||
      (t.description && t.description.toLowerCase().includes(q))
    )
  }

  // 排序
  const priorityOrder = { high: 3, medium: 2, low: 1 }
  list.sort((a, b) => {
    if (sortKey.value === 'priority') {
      return (priorityOrder[b.priority] || 0) - (priorityOrder[a.priority] || 0)
    }
    const key = sortKey.value
    const av = a[key]
    const bv = b[key]
    if (!av && !bv) return 0
    if (!av) return 1
    if (!bv) return -1
    const ad = new Date(av).getTime()
    const bd = new Date(bv).getTime()
    return bd - ad
  })

  return list
})

const openAddDialog = () => {
  currentTodo.value = { title: '', priority: 'medium', status: 'pending' }
  showDialog.value = true
}

const editTodo = (todo) => {
  currentTodo.value = { ...todo }
  showDialog.value = true
}

const saveTodo = async () => {
  if (!currentTodo.value.title) return ElMessage.warning('标题不能为空')
  
  try {
    if (currentTodo.value.id) {
      await api.put(`/todos/${currentTodo.value.id}`, currentTodo.value)
    } else {
      await api.post('/todos', currentTodo.value)
    }
    ElMessage.success('保存成功')
    showDialog.value = false
    loadTodos()
  } catch (e) {
    ElMessage.error('保存失败')
  }
}

const deleteTodo = async (todo) => {
  try {
    await ElMessageBox.confirm('确定删除此任务?', '警告', { type: 'warning' })
    await api.delete(`/todos/${todo.id}`)
    loadTodos()
    ElMessage.success('已删除')
  } catch (e) {}
}

// Helpers
const getPriorityType = (p) => ({ high: 'danger', medium: 'warning', low: 'info' }[p] || 'info')
const getPriorityText = (p) => ({ high: 'High', medium: 'Medium', low: 'Low' }[p] || p)
const getStatusType = (s) => ({ pending: 'info', in_progress: 'primary', completed: 'success' }[s])
const getStatusText = (s) => ({ pending: '待办', in_progress: '进行中', completed: '已完成' }[s])
const formatDate = (d) => new Date(d).toLocaleString()
const formatTimeAgo = (d) => {
    if (!d) return ''
    const date = new Date(d)
    const now = new Date()
    const diff = Math.floor((now - date) / 1000 / 60) // minutes
    if (diff < 60) return `${diff}m ago`
    if (diff < 1440) return `${Math.floor(diff/60)}h ago`
    return `${Math.floor(diff/1440)}d ago`
}

onMounted(loadTodos)
</script>

<style scoped>
.todo-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 20px;
  gap: 20px;
}

.glass-panel {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
}

.todo-header {
  padding: 15px 25px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-left { display: flex; align-items: center; gap: 20px; }
.header-left h2 { margin: 0; color: #fff; font-size: 20px; }
.header-right { display: flex; gap: 15px; }
.filter-select { width: 130px; }
.search-input { width: 200px; }

/* Kanban Board */
.kanban-board {
  flex: 1;
  display: flex;
  gap: 20px;
  overflow: hidden; /* Allow columns to scroll internally */
}

.kanban-column {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: rgba(0,0,0,0.2);
  border-radius: 16px;
  padding: 15px;
  min-width: 280px;
}

.column-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
  color: #ddd;
}
.column-header h3 { margin: 0; font-size: 16px; font-weight: 600; flex: 1; }
.dot { width: 8px; height: 8px; border-radius: 50%; }
.dot.pending { background: #909399; box-shadow: 0 0 5px #909399; }
.dot.progress { background: #E6A23C; box-shadow: 0 0 5px #E6A23C; }
.dot.completed { background: #67C23A; box-shadow: 0 0 5px #67C23A; }
.count { background: rgba(255,255,255,0.1); padding: 2px 8px; border-radius: 10px; font-size: 12px; }

.column-content {
  flex: 1;
  overflow-y: auto;
}
.drag-area {
  min-height: 100%;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.todo-card {
  padding: 15px;
  cursor: grab;
  transition: transform 0.2s, box-shadow 0.2s;
  border-left: 3px solid transparent;
}
.todo-card:active { cursor: grabbing; }
.todo-card:hover {
  transform: translateY(-2px);
  background: rgba(255,255,255,0.1);
}

.card-badges {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 11px;
  color: #aaa;
}

.card-title {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #eee;
  line-height: 1.4;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
}

.completed .card-title { text-decoration: line-through; opacity: 0.6; }
.completed { border-left-color: #67C23A; }

/* List View */
.list-view {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}
.list-view :deep(.el-table) {
  --el-table-text-color: #ddd;
  --el-table-header-text-color: #fff;
  --el-table-row-hover-bg-color: rgba(255,255,255,0.1);
}
</style>
