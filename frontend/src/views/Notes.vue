<template>
  <div class="notes-container">
    <!-- Left Sidebar: Note List -->
    <div class="notes-sidebar glass-panel">
      <div class="sidebar-header">
        <div class="header-top">
          <span class="app-title">Nexus Notes</span>
          <el-tooltip content="新建笔记" placement="bottom">
            <el-button circle type="primary" :icon="Plus" @click="createNote" />
          </el-tooltip>
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
            <el-tooltip content="保存笔记" placement="top">
              <el-button circle type="success" :icon="Check" @click="save(true)" />
            </el-tooltip>
            <el-tooltip content="导出笔记" placement="top">
              <el-button circle type="warning" :icon="Collection" @click="handleCommand('export')" />
            </el-tooltip>
            <el-tooltip content="删除笔记" placement="top">
              <el-button circle type="danger" :icon="Delete" @click="handleCommand('delete')" />
            </el-tooltip>
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

const exportPDF = async () => {
  if (!current.value) {
    ElMessage.warning('请先选择一个笔记')
    return
  }
  
  try {
    ElMessage.info('正在准备PDF内容...')
    
    // Get the markdown content
    const markdownContent = current.value.content || ''
    const title = current.value.title || 'Nexus Notes'
    const tags = current.value.tags || ''
    
    // Convert markdown to HTML with better styling
    const htmlContent = markdownContent
      // Headers
      .replace(/^### (.*$)/gim, '<h3 style="font-size: 18px; margin-top: 16px; color: #333;">$1</h3>')
      .replace(/^## (.*$)/gim, '<h2 style="font-size: 22px; margin-top: 20px; color: #333;">$1</h2>')
      .replace(/^# (.*$)/gim, '<h1 style="font-size: 26px; margin-top: 24px; color: #333;">$1</h1>')
      // Text formatting
      .replace(/\*\*(.*)\*\*/gim, '<strong style="color: #333;">$1</strong>')
      .replace(/\*(.*)\*/gim, '<em style="color: #333;">$1</em>')
      // Code blocks
      .replace(/```([\s\S]*?)```/gim, '<pre style="background-color: #f5f5f5; padding: 15px; border-radius: 6px; border-left: 4px solid #409EFF; margin: 16px 0;"><code style="font-family: Monaco, Consolas, monospace; color: #333;">$1</code></pre>')
      .replace(/`(.*?)`/gim, '<code style="background-color: #f0f0f0; padding: 2px 6px; border-radius: 3px; color: #d73a49; font-family: Monaco, Consolas, monospace;">$1</code>')
      // Line breaks
      .replace(/\n/gim, '<br>')
    
    // Create a temporary element for PDF generation
    const tempDiv = document.createElement('div')
    tempDiv.style.cssText = `
      position: absolute;
      left: -9999px;
      top: -9999px;
      width: 800px;
      padding: 40px;
      background: white;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Fira Sans", "Droid Sans", "Helvetica Neue", Arial, sans-serif;
      line-height: 1.6;
      color: #333;
    `
    
    tempDiv.innerHTML = `
      <div style="text-align: center; margin-bottom: 30px; border-bottom: 2px solid #409EFF; padding-bottom: 20px;">
        <h1 style="font-size: 28px; color: #333; margin: 0;">${title}</h1>
        ${tags ? `<p style="color: #666; margin: 8px 0 0 0; font-size: 14px;">标签: ${tags}</p>` : ''}
        <p style="color: #999; margin: 8px 0 0 0; font-size: 12px;">导出时间: ${new Date().toLocaleString('zh-CN')}</p>
      </div>
      <div style="margin-top: 20px;">
        ${htmlContent}
      </div>
    `
    
    document.body.appendChild(tempDiv)
    
    // Configure PDF options
    const opt = {
      margin: [10, 10, 10, 10],
      filename: `${title.replace(/[^a-zA-Z0-9一-鿿]/g, '_')}.pdf`,
      image: { type: 'jpeg', quality: 0.98 },
      html2canvas: { 
        scale: 2,
        useCORS: true,
        letterRendering: true,
        allowTaint: true
      },
      jsPDF: { 
        unit: 'mm', 
        format: 'a4', 
        orientation: 'portrait',
        compress: true
      },
      pagebreak: { 
        mode: ['avoid-all', 'css', 'legacy'],
        before: '.page-break-before',
        after: '.page-break-after'
      }
    }
    
    ElMessage.info('正在生成PDF文件...')
    
    // Generate PDF
    await html2pdf().from(tempDiv).set(opt).save()
    
    // Clean up
    document.body.removeChild(tempDiv)
    
    ElMessage.success('PDF导出成功！')
  } catch (error) {
    console.error('PDF导出失败:', error)
    ElMessage.error(`PDF导出失败: ${error.message || '请重试'}`)
  }
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
  background: var(--card-bg);
  backdrop-filter: blur(6px);
  border-radius: 16px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
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
  color: var(--text-primary);
}

.notes-list {
  flex: 1;
  overflow: hidden;
}

.note-item {
  padding: 15px 20px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  transition: all 0.2s;
  color: var(--text-primary);
}

.note-item:hover {
  background: var(--sidebar-hover);
}

.note-item.active {
  background: rgba(102, 126, 234, 0.1);
  border-right: 3px solid var(--primary-color);
}

.note-item-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
}

.note-title {
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.pin-icon {
  color: var(--warning-color);
  font-size: 14px;
}

.note-preview {
  font-size: 13px;
  color: var(--text-secondary);
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
  color: var(--text-secondary);
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
  color: var(--text-primary);
  outline: none;
}

.title-input::placeholder {
  color: var(--text-placeholder);
}

.editor-actions {
  display: flex;
  gap: 10px;
}

.tags-input-wrapper {
  padding: 10px 30px;
  border-bottom: 1px solid var(--border-color);
}

.tags-input :deep(.el-input__wrapper) {
  box-shadow: none;
  background: transparent;
  padding: 0;
  color: var(--text-primary);
}

.tags-input :deep(.el-input__input) {
  color: var(--text-primary);
}

.tags-input :deep(.el-input__placeholder) {
  color: var(--text-placeholder);
}

.editor-body {
  flex: 1;
  overflow: hidden;
}

/* v-md-editor overrides - 更高优先级的样式覆盖 */
:deep(.v-md-editor) {
  box-shadow: none !important;
  background: transparent !important;
  color: var(--text-primary) !important;
}

:deep(.v-md-editor *),
:deep(.v-md-editor__editor-wrapper *), 
:deep(.v-md-editor__preview-wrapper *) {
  color: var(--text-primary) !important;
}

:deep(.v-md-editor__editor-wrapper), 
:deep(.v-md-editor__preview-wrapper) {
  background: transparent !important;
  color: var(--text-primary) !important;
}

:deep(.vuepress-markdown-body) {
  background: transparent !important;
  padding: 20px 30px !important;
  color: var(--text-primary) !important;
}

:deep(.vuepress-markdown-body h1),
:deep(.vuepress-markdown-body h2),
:deep(.vuepress-markdown-body h3),
:deep(.vuepress-markdown-body h4),
:deep(.vuepress-markdown-body h5),
:deep(.vuepress-markdown-body h6) {
  color: var(--text-primary) !important;
}

:deep(.vuepress-markdown-body p),
:deep(.vuepress-markdown-body li),
:deep(.vuepress-markdown-body blockquote),
:deep(.vuepress-markdown-body table),
:deep(.vuepress-markdown-body tbody),
:deep(.vuepress-markdown-body td),
:deep(.vuepress-markdown-body th) {
  color: var(--text-primary) !important;
}

:deep(.vuepress-markdown-body code) {
  color: var(--text-primary) !important;
  background-color: rgba(102, 126, 234, 0.1) !important;
}

:deep(.vuepress-markdown-body pre) {
  background-color: var(--card-bg) !important;
  border-color: var(--border-color) !important;
}

:deep(.vuepress-markdown-body pre code) {
  color: var(--text-primary) !important;
  background-color: transparent !important;
}

:deep(.v-md-editor__toolbar) {
  background: var(--card-bg) !important;
  border-color: var(--border-color) !important;
}

:deep(.v-md-editor__toolbar-btn) {
  color: var(--text-primary) !important;
}

:deep(.v-md-editor__toolbar-btn:hover) {
  background: var(--sidebar-hover) !important;
  color: var(--text-primary) !important;
}

:deep(.v-md-editor__textarea) {
  background: transparent !important;
  color: var(--text-primary) !important;
}

:deep(.v-md-editor__textarea::placeholder) {
  color: var(--text-placeholder) !important;
}

:deep(.v-md-editor__counter) {
  color: var(--text-secondary) !important;
}

:deep(.v-md-editor__scroll) {
  color: var(--text-primary) !important;
}

/* 直接覆盖vuepress主题的样式 */
:deep(.content__default),
:deep(.v-md-editor__content) {
  color: var(--text-primary) !important;
}

/* 覆盖prismjs语法高亮的样式 */
:deep(.token) {
  color: var(--text-primary) !important;
}

:deep(.token.comment),
:deep(.token.prolog),
:deep(.token.doctype),
:deep(.token.cdata) {
  color: var(--text-secondary) !important;
}

:deep(.token.punctuation) {
  color: var(--text-primary) !important;
}

:deep(.token.property),
:deep(.token.tag),
:deep(.token.boolean),
:deep(.token.number),
:deep(.token.constant),
:deep(.token.symbol),
:deep(.token.deleted) {
  color: var(--text-primary) !important;
}

:deep(.token.selector),
:deep(.token.attr-name),
:deep(.token.string),
:deep(.token.char),
:deep(.token.builtin),
:deep(.token.inserted) {
  color: var(--text-primary) !important;
}

:deep(.token.operator),
:deep(.token.entity),
:deep(.token.url),
:deep(.language-css .token.string),
:deep(.style .token.string) {
  color: var(--text-primary) !important;
}

:deep(.token.atrule),
:deep(.token.attr-value),
:deep(.token.keyword) {
  color: var(--text-primary) !important;
}

:deep(.token.function),
:deep(.token.class-name) {
  color: var(--text-primary) !important;
}

:deep(.token.regex),
:deep(.token.important),
:deep(.token.variable) {
  color: var(--text-primary) !important;
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


</style>
