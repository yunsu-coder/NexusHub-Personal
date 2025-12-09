<template>
  <div class="notes-page">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="notes-list-card">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>笔记列表</span>
              <el-button type="primary" size="small" @click="createNote">新建</el-button>
            </div>
          </template>
          <el-input v-model="search" placeholder="搜索..." clearable style="margin-bottom: 15px" />
          <div class="notes-list">
            <div v-for="note in filteredNotes" :key="note.id" class="note-item" :class="{active: current?.id === note.id}" @click="selectNote(note)">
              <h4>{{ note.title || '无标题' }}</h4>
              <p>{{ note.content?.substring(0,50) }}...</p>
              <small>{{ new Date(note.updated_at).toLocaleString('zh-CN') }}</small>
            </div>
            <el-empty v-if="notes.length === 0" />
          </div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card v-if="current">
          <template #header>
            <div style="display: flex; gap: 10px; align-items: center">
              <el-input v-model="current.title" placeholder="标题..." style="flex: 1" />
              <el-switch v-model="isEditing" active-text="预览" inactive-text="编辑" :active-value="true" :inactive-value="false" />
              <el-button type="primary" @click="save">保存</el-button>
              <el-button type="success" @click="exportPDF" :loading="exporting">导出PDF</el-button>
              <el-button type="danger" @click="del">删除</el-button>
            </div>
          </template>
          
          <!-- Markdown 编辑器 -->
          <div v-if="!isEditing" class="editor-container">
            <textarea 
              v-model="current.content" 
              class="markdown-editor"
              placeholder="使用 Markdown 编写您的笔记...

# 标题
## 二级标题
**加粗** _斜体_

- 列表项1
- 列表项2

```javascript
// 代码块
console.log('Hello World');
```
" 
            ></textarea>
          </div>
          
          <!-- Markdown 预览 -->
          <div v-else class="preview-container">
            <div class="markdown-preview" v-html="renderedContent"></div>
          </div>
          
          <el-input v-model="current.tags" placeholder="标签 (逗号分隔)" style="margin-top: 10px" />
        </el-card>
        <el-empty v-else description="选择笔记" />
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'
import { marked } from 'marked'
import html2pdf from 'html2pdf.js'

const notes = ref([])
const current = ref(null)
const search = ref('')
const isEditing = ref(false)
const renderedContent = ref('')
const exporting = ref(false)

const filteredNotes = computed(() => notes.value.filter(n => !search.value || n.title?.includes(search.value) || n.content?.includes(search.value)))

// 渲染 Markdown
const renderMarkdown = () => {
  if (current.value?.content) {
    renderedContent.value = marked(current.value.content)
  } else {
    renderedContent.value = ''
  }
}

// 监听内容变化自动渲染
watch(() => current.value?.content, () => {
  renderMarkdown()
}, { deep: true })

const load = async () => {
  try {
    notes.value = await api.get('/notes')
  } catch {
    ElMessage.error('加载失败')
  }
}

const createNote = () => {
  current.value = { title: '', content: '', tags: '' }
  renderMarkdown()
}

const selectNote = (n) => {
  current.value = {...n}
  renderMarkdown()
}

const save = async () => {
  if (!current.value.title?.trim() && !current.value.content?.trim()) {
    ElMessage.warning('请填写标题或内容')
    return
  }
  
  try {
    if (current.value.id) {
      await api.put(`/notes/${current.value.id}`, current.value)
    } else {
      current.value = await api.post('/notes', current.value)
    }
    ElMessage.success('保存成功')
    await load()
  } catch {
    ElMessage.error('保存失败')
  }
}

const del = async () => {
  try {
    await ElMessageBox.confirm('确定删除此笔记?')
    await api.delete(`/notes/${current.value.id}`)
    ElMessage.success('删除成功')
    current.value = null
    await load()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 导出为 PDF
const exportPDF = async () => {
  if (!current.value?.content) {
    ElMessage.warning('没有内容可导出')
    return
  }
  
  exporting.value = true
  
  try {
    // 创建临时元素
    const tempElement = document.createElement('div')
    tempElement.innerHTML = renderedContent.value
    tempElement.style.padding = '20px'
    tempElement.style.maxWidth = '800px'
    tempElement.style.margin = '0 auto'
    
    // 添加标题
    const titleElement = document.createElement('h1')
    titleElement.textContent = current.value.title || '无标题'
    titleElement.style.textAlign = 'center'
    titleElement.style.marginBottom = '20px'
    tempElement.insertBefore(titleElement, tempElement.firstChild)
    
    // 配置选项
    const options = {
      margin: 10,
      filename: `${current.value.title || '笔记'}.pdf`,
      image: { type: 'jpeg', quality: 0.98 },
      html2canvas: { scale: 2 },
      jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' }
    }
    
    // 生成 PDF
    await html2pdf().from(tempElement).set(options).save()
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

onMounted(() => {
  load()
})
</script>

<style scoped>
.notes-page { min-height: 70vh; }
.notes-list-card :deep(.el-card__body) { max-height: 70vh; overflow-y: auto; }
.notes-list { display: flex; flex-direction: column; gap: 10px; }
.note-item { padding: 12px; border: 1px solid var(--border-color); border-radius: 6px; cursor: pointer; transition: all 0.3s; }
.note-item:hover, .note-item.active { background: var(--hover-bg); border-color: #409EFF; }
.note-item h4 { margin: 0 0 8px 0; color: var(--text-primary); }
.note-item p { margin: 0 0 8px 0; font-size: 13px; color: var(--text-secondary); }
.note-item small { font-size: 12px; color: var(--text-secondary); }

/* 编辑器和预览样式 */
.editor-container, .preview-container {
  min-height: 400px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  overflow: hidden;
  transition: all 0.3s;
}

.markdown-editor {
  width: 100%;
  height: 400px;
  padding: 16px;
  border: none;
  resize: none;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
  background-color: #ffffff;
  color: #333333;
  outline: none;
  border: 1px solid #e5e7eb;
}

.markdown-preview {
  padding: 20px;
  min-height: 400px;
  background-color: #fff;
  line-height: 1.8;
  overflow-y: auto;
}

/* Markdown 样式 */
.markdown-preview :deep(h1),
.markdown-preview :deep(h2),
.markdown-preview :deep(h3),
.markdown-preview :deep(h4),
.markdown-preview :deep(h5),
.markdown-preview :deep(h6) {
  margin: 20px 0 10px 0;
  font-weight: 600;
  color: var(--text-primary);
}

.markdown-preview :deep(h1) { font-size: 2em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
.markdown-preview :deep(h2) { font-size: 1.5em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
.markdown-preview :deep(h3) { font-size: 1.25em; }

.markdown-preview :deep(p) { margin: 10px 0; }
.markdown-preview :deep(strong) { font-weight: 600; }
.markdown-preview :deep(em) { font-style: italic; }

.markdown-preview :deep(a) {
  color: #409EFF;
  text-decoration: none;
  transition: color 0.3s;
}

.markdown-preview :deep(a:hover) {
  color: #66b1ff;
  text-decoration: underline;
}

.markdown-preview :deep(ul),
.markdown-preview :deep(ol) {
  margin: 10px 0;
  padding-left: 25px;
}

.markdown-preview :deep(li) {
  margin: 5px 0;
}

.markdown-preview :deep(code) {
  padding: 2px 4px;
  background-color: #f6f8fa;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 85%;
  color: #e83e8c;
}

.markdown-preview :deep(pre) {
  background-color: #f6f8fa;
  border-radius: 6px;
  padding: 16px;
  overflow-x: auto;
  margin: 16px 0;
}

.markdown-preview :deep(pre code) {
  background-color: transparent;
  padding: 0;
  color: inherit;
  font-size: 100%;
}

.markdown-preview :deep(blockquote) {
  border-left: 4px solid #dfe2e5;
  padding: 0 15px;
  color: #6a737d;
  margin: 16px 0;
}

.markdown-preview :deep(hr) {
  height: 0.25em;
  padding: 0;
  margin: 24px 0;
  background-color: #e1e4e8;
  border: 0;
}

.markdown-preview :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 16px 0;
}

.markdown-preview :deep(th),
.markdown-preview :deep(td) {
  border: 1px solid #dfe2e5;
  padding: 6px 13px;
  text-align: left;
}

.markdown-preview :deep(tr) {
  background-color: #fff;
  border-top: 1px solid #c6cbd1;
}

.markdown-preview :deep(tr:nth-child(2n)) {
  background-color: #f6f8fa;
}
</style>
