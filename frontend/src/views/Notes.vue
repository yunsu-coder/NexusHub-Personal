<template>
  <div class="notes-container">
    <!-- Left Sidebar: Note List -->
    <div class="notes-sidebar glass-panel">
      <div class="sidebar-header">
        <div class="header-top">
          <span class="app-title">Nexus Notes</span>
          <el-button circle type="primary" :icon="Plus" @click="createNote" />
        </div>
        <div class="search-box">
          <el-input
            v-model="search"
            placeholder="Search notes..."
            :prefix-icon="Search"
            clearable
          />
        </div>
      </div>

      <div class="notes-list" v-loading="loading">
        <el-scrollbar>
          <div
            v-for="note in filteredNotes"
            :key="note.id"
            class="note-item"
            :class="{ active: current?.id === note.id }"
            @click="selectNote(note)"
          >
            <div class="note-item-header">
              <span class="note-title">{{ note.title || 'Untitled Note' }}</span>
              <el-icon v-if="note.is_pinned" class="pin-icon"><paperclip /></el-icon>
            </div>
            <div class="note-preview">
              {{ getPreview(note.content) }}
            </div>
            <div class="note-meta">
              <span class="note-date">{{ formatDate(note.updated_at) }}</span>
              <el-tag v-if="note.tags" size="small" effect="plain" class="note-tag">
                {{ note.tags.split(',')[0] }}
              </el-tag>
            </div>
          </div>
          <el-empty v-if="filteredNotes.length === 0" description="No notes found" />
        </el-scrollbar>
      </div>
    </div>

    <!-- Right Main: Editor -->
    <div class="notes-main glass-panel">
      <div v-if="current" class="editor-wrapper">
        <div class="editor-header">
          <div class="title-input-wrapper">
            <input
              v-model="current.title"
              class="title-input"
              placeholder="Enter title..."
              @change="save(false)"
            />
          </div>
          <div class="editor-actions">
            <el-tooltip content="Toggle Pin">
              <el-button
                circle
                :type="current.is_pinned ? 'warning' : 'default'"
                :icon="Paperclip"
                @click="togglePin"
              />
            </el-tooltip>
            <el-tooltip content="Save (Ctrl+S)">
              <el-button circle type="success" :icon="Check" @click="save(true)" />
            </el-tooltip>
            <el-dropdown trigger="click" @command="handleCommand">
              <el-button circle :icon="MoreFilled" />
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="export">Export PDF</el-dropdown-item>
                  <el-dropdown-item command="delete" divided style="color: #f56c6c">Delete</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
        
        <div class="tags-input-wrapper">
           <el-input 
             v-model="current.tags" 
             placeholder="Tags (comma separated)..." 
             size="small" 
             class="tags-input"
           >
             <template #prefix><el-icon><Collection /></el-icon></template>
           </el-input>
        </div>

        <div class="editor-body">
          <v-md-editor
            v-model="current.content"
            height="100%"
            placeholder="Start writing..."
            :disabled-menus="[]"
            @save="save(true)"
            @change="handleContentChange"
          ></v-md-editor>
        </div>
      </div>
      <div v-else class="empty-state">
        <el-icon :size="60"><EditPen /></el-icon>
        <h2>Select or create a note</h2>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, Search, Paperclip, Check, MoreFilled, 
  Delete, EditPen, Collection 
} from '@element-plus/icons-vue'
import api from '../api'
import html2pdf from 'html2pdf.js'

const notes = ref([])
const current = ref(null)
const search = ref('')
const loading = ref(false)
let autosaveTimer = null

const filteredNotes = computed(() => {
  // 使用拷贝避免 sort 直接修改原始数组
  let list = [...notes.value]
  if (search.value) {
    const q = search.value.toLowerCase()
    list = list.filter(n => 
      (n.title && n.title.toLowerCase().includes(q)) || 
      (n.content && n.content.toLowerCase().includes(q))
    )
  }
  // Sort: Pinned first, then by updated_at desc
  return list.sort((a, b) => {
    if (a.is_pinned !== b.is_pinned) return b.is_pinned ? 1 : -1
    return new Date(b.updated_at) - new Date(a.updated_at)
  })
})

const loadNotes = async () => {
  loading.value = true
  try {
    const res = await api.get('/notes')
    notes.value = res || []
  } catch (error) {
    ElMessage.error('Failed to load notes')
  } finally {
    loading.value = false
  }
}

const createNote = () => {
  const newNote = {
    title: '',
    content: '',
    tags: '',
    is_pinned: false
  }
  // Temporarily add to list or just set current (will save on first edit/save)
  current.value = newNote
}

const selectNote = (note) => {
  current.value = { ...note } // Clone to avoid direct mutation of list item
}

const getPreview = (content) => {
  if (!content) return 'No content'
  // Strip markdown symbols for preview (simple regex)
  return content.replace(/[#*`]/g, '').substring(0, 60) + (content.length > 60 ? '...' : '')
}

const formatDate = (dateStr) => {
  if (!dateStr) return 'Just now'
  const date = new Date(dateStr)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})
}

const togglePin = () => {
  if (!current.value) return
  current.value.is_pinned = !current.value.is_pinned
  save(false)
}

const handleContentChange = () => {
  if (autosaveTimer) clearTimeout(autosaveTimer)
  autosaveTimer = setTimeout(() => {
    save(false) // Autosave silent
    saveDraft()
  }, 2000)
}

const saveDraft = () => {
  if (!current.value) return
  try {
    const draft = {
      title: current.value.title,
      content: current.value.content,
      tags: current.value.tags,
      is_pinned: current.value.is_pinned,
      id: current.value.id || null
    }
    localStorage.setItem('notes_draft', JSON.stringify(draft))
  } catch {}
}

const loadDraft = () => {
  try {
    const raw = localStorage.getItem('notes_draft')
    if (!raw) return
    if (current.value) return
    const draft = JSON.parse(raw)
    if (!draft || (!draft.title && !draft.content)) return
    current.value = { ...draft }
    ElMessage.info('已从本地草稿恢复未保存的笔记')
  } catch {}
}

const save = async (notify = true) => {
  if (!current.value) return
  
  // If empty, don't save to DB unless it already has ID
  if (!current.value.title && !current.value.content && !current.value.id) return

  try {
    let res
    if (current.value.id) {
      res = await api.put(`/notes/${current.value.id}`, current.value)
    } else {
      res = await api.post('/notes', current.value)
    }
    
    // Update local list
    if (current.value.id) {
      const index = notes.value.findIndex(n => n.id === current.value.id)
      if (index !== -1) notes.value[index] = { ...res }
    } else {
      notes.value.unshift(res)
      current.value = { ...res } // Update current with ID
    }
    
    if (notify) ElMessage.success('Saved')
  } catch (error) {
    console.error(error)
    if (notify) ElMessage.error('Save failed')
  }
}

const handleCommand = (cmd) => {
  if (cmd === 'delete') deleteNote()
  if (cmd === 'export') exportPDF()
}

const deleteNote = async () => {
  if (!current.value || !current.value.id) {
    current.value = null
    return
  }
  
  try {
    await ElMessageBox.confirm('Are you sure you want to delete this note?', 'Warning', {
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      type: 'warning'
    })
    
    await api.delete(`/notes/${current.value.id}`)
    notes.value = notes.value.filter(n => n.id !== current.value.id)
    current.value = null
    ElMessage.success('Deleted')
  } catch (e) {
    // Cancelled or error
  }
}

const exportPDF = () => {
  // Use v-md-editor's preview class to capture content
  // But v-md-editor is in shadow dom or tricky?
  // Actually, we can use a trick: render a hidden div with the content and print that
  // Or grab .v-md-editor-preview element
  const element = document.querySelector('.v-md-editor__preview')
  if (!element) return
  
  const opt = {
    margin: 10,
    filename: (current.value.title || 'note') + '.pdf',
    image: { type: 'jpeg', quality: 0.98 },
    html2canvas: { scale: 2 },
    jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' }
  }
  
  html2pdf().from(element).set(opt).save()
}

onMounted(() => {
  loadNotes()
  loadDraft()
})
</script>

<style scoped>
.notes-container {
  display: flex;
  height: calc(100vh - 140px); /* Adjust based on App.vue layout */
  gap: 20px;
  padding-bottom: 20px;
}

.glass-panel {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* Sidebar */
.notes-sidebar {
  width: 320px;
  flex-shrink: 0;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.app-title {
  font-size: 1.2rem;
  font-weight: 700;
  color: #333;
}

.notes-list {
  flex: 1;
  overflow: hidden;
}

.note-item {
  padding: 15px 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.03);
  cursor: pointer;
  transition: all 0.2s;
}

.note-item:hover {
  background: rgba(255, 255, 255, 0.5);
}

.note-item.active {
  background: rgba(64, 158, 255, 0.1);
  border-right: 3px solid #409EFF;
}

.note-item-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
}

.note-title {
  font-weight: 600;
  color: #2c3e50;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.pin-icon {
  color: #E6A23C;
  font-size: 14px;
}

.note-preview {
  font-size: 13px;
  color: #606266;
  margin-bottom: 8px;
  height: 36px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.note-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.note-date {
  font-size: 12px;
  color: #909399;
}

/* Main Editor */
.notes-main {
  flex: 1;
  position: relative;
}

.editor-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.editor-header {
  padding: 15px 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.title-input-wrapper {
  flex: 1;
  margin-right: 20px;
}

.title-input {
  width: 100%;
  border: none;
  background: transparent;
  font-size: 24px;
  font-weight: 700;
  color: #333;
  outline: none;
}

.title-input::placeholder {
  color: #ccc;
}

.editor-actions {
  display: flex;
  gap: 10px;
}

.tags-input-wrapper {
  padding: 10px 30px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.03);
}

.tags-input :deep(.el-input__wrapper) {
  box-shadow: none;
  background: transparent;
  padding: 0;
}

.editor-body {
  flex: 1;
  overflow: hidden;
}

/* v-md-editor overrides */
:deep(.v-md-editor) {
  box-shadow: none;
  background: transparent;
}

:deep(.v-md-editor__editor-wrapper), :deep(.v-md-editor__preview-wrapper) {
  background: transparent;
}

:deep(.vuepress-markdown-body) {
  background: transparent;
  padding: 20px 30px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
  gap: 20px;
}

/* Dark Mode Adaptation */
@media (prefers-color-scheme: dark) {
  .glass-panel {
    background: rgba(30, 30, 30, 0.6);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .app-title, .title-input, .note-title {
    color: #eee;
  }
  
  .note-preview {
    color: #aaa;
  }
  
  .sidebar-header, .note-item, .editor-header, .tags-input-wrapper {
    border-color: rgba(255, 255, 255, 0.05);
  }
  
  .note-item:hover {
    background: rgba(255, 255, 255, 0.1);
  }
  
  :deep(.vuepress-markdown-body) {
    color: #eee;
  }
}
</style>
