<template>
  <div class="code-editor-page">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between">
              <span>代码片段</span>
              <el-button type="primary" size="small" @click="createSnippet">新建</el-button>
            </div>
          </template>
          <div class="snippets-list">
            <div v-for="snippet in snippets" :key="snippet.id" class="snippet-item" :class="{active: current?.id === snippet.id}" @click="selectSnippet(snippet)">
              <h4>{{ snippet.title }}</h4>
              <el-tag size="small">{{ snippet.language }}</el-tag>
              <small>{{ new Date(snippet.updated_at).toLocaleString('zh-CN') }}</small>
            </div>
            <el-empty v-if="snippets.length === 0" />
          </div>
        </el-card>
      </el-col>
      <el-col :span="18">
        <el-card v-if="current">
          <template #header>
            <div style="display: flex; gap: 10px; align-items: center">
              <el-input v-model="current.title" placeholder="标题..." style="flex: 1" />
              <el-select v-model="current.language" placeholder="语言" style="width: 150px" @change="updateHighlight">
                <el-option label="JavaScript" value="javascript" />
                <el-option label="Python" value="python" />
                <el-option label="Go" value="go" />
                <el-option label="C++" value="cpp" />
                <el-option label="Java" value="java" />
                <el-option label="HTML" value="html" />
                <el-option label="CSS" value="css" />
                <el-option label="Markdown" value="markdown" />
              </el-select>
              <el-switch v-model="isEditing" active-text="编辑" inactive-text="预览" @change="updateHighlight" :active-value="false" :inactive-value="true" />
              <el-button type="primary" @click="save">保存</el-button>
              <el-button type="danger" @click="del">删除</el-button>
            </div>
          </template>

          <!-- 编辑模式 -->
          <div v-if="!isEditing" class="code-editor-wrapper">
            <textarea
              v-model="current.code"
              @input="updateHighlight"
              @keydown="handleKeyDown"
              placeholder="输入代码..."
              class="code-textarea"
              spellcheck="false"
            ></textarea>
          </div>

          <!-- 预览模式 -->
          <div v-else class="code-preview">
            <pre><code :class="`language-${current.language}`" v-html="highlightedCode"></code></pre>
          </div>

          <div style="margin-top: 10px">
            <el-input v-model="current.tags" placeholder="标签 (逗号分隔)" />
          </div>
        </el-card>
        <el-empty v-else description="选择或创建代码片段" />
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'
import api from '../api'

const snippets = ref([])
const current = ref(null)
const isEditing = ref(false)  // false = 编辑, true = 预览
const highlightedCode = ref('')

const load = async () => { try { snippets.value = await api.get('/code') } catch { ElMessage.error('加载失败') } }

const createSnippet = () => {
  current.value = { title: '', language: 'javascript', code: '', tags: '' }
  isEditing.value = false  // 新建时默认编辑模式
}

const selectSnippet = (s) => {
  current.value = {...s}
  isEditing.value = false  // 选择时默认编辑模式
  updateHighlight()
}

const updateHighlight = () => {
  if (!current.value?.code || !isEditing.value) {
    highlightedCode.value = ''
    return
  }

  try {
    const language = current.value.language === 'cpp' ? 'cpp' : current.value.language
    highlightedCode.value = hljs.highlight(current.value.code, { language }).value
  } catch (error) {
    highlightedCode.value = hljs.highlightAuto(current.value.code).value
  }
}

const handleKeyDown = (event) => {
  // 处理Tab键插入4个空格
  if (event.key === 'Tab') {
    event.preventDefault()
    const textarea = event.target
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const spaces = '    ' // 4个空格

    // 插入4个空格
    current.value.code =
      current.value.code.substring(0, start) +
      spaces +
      current.value.code.substring(end)

    // 设置光标位置
    setTimeout(() => {
      textarea.selectionStart = textarea.selectionEnd = start + spaces.length
    }, 0)
  }
}

const save = async () => {
  try {
    if (current.value.id) {
      await api.put(`/code/${current.value.id}`, current.value)
    } else {
      current.value = await api.post('/code', current.value)
    }
    ElMessage.success('保存成功')
    load()
  } catch {
    ElMessage.error('保存失败')
  }
}

const del = async () => {
  try {
    await ElMessageBox.confirm('确定删除?')
    await api.delete(`/code/${current.value.id}`)
    ElMessage.success('删除成功')
    current.value = null
    load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

watch(() => current.value?.code, () => {
  if (isEditing.value) {
    updateHighlight()
  }
})

onMounted(load)
</script>

<style scoped>
.code-editor-page { min-height: 70vh; }
.snippets-list { max-height: 70vh; overflow-y: auto; display: flex; flex-direction: column; gap: 10px; }
.snippet-item { padding: 12px; border: 1px solid var(--border-color); border-radius: 6px; cursor: pointer; transition: all 0.3s; }
.snippet-item:hover, .snippet-item.active { background: var(--hover-bg); border-color: #409EFF; }
.snippet-item h4 { margin: 0 0 8px 0; color: var(--text-primary); }
.snippet-item small { display: block; margin-top: 8px; font-size: 12px; color: var(--text-secondary); }

.code-editor-wrapper {
  position: relative;
  min-height: 600px;
}

.code-textarea {
  width: 100%;
  min-height: 600px;
  padding: 16px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-primary);
  background-color: var(--hover-bg);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  resize: vertical;
  outline: none;
  white-space: pre;
  overflow-wrap: normal;
  overflow-x: auto;
  tab-size: 4;
}

.code-textarea:focus {
  border-color: #409EFF;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

.code-preview {
  background-color: #282c34 !important;
  border-radius: 6px;
  overflow: auto;
  max-height: 600px;
  border: 1px solid #3e4451;
}

.code-preview pre {
  margin: 0;
  padding: 16px;
  background: #282c34 !important;
}

.code-preview code {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre;
  overflow-wrap: normal;
  display: block;
  color: #abb2bf !important;
}

/* 强制应用语法高亮颜色 */
:deep(.hljs-keyword) { color: #c678dd !important; }
:deep(.hljs-string) { color: #98c379 !important; }
:deep(.hljs-comment) { color: #5c6370 !important; font-style: italic !important; }
:deep(.hljs-function) { color: #61afef !important; }
:deep(.hljs-title) { color: #e5c07b !important; }
:deep(.hljs-number) { color: #d19a66 !important; }
:deep(.hljs-built_in) { color: #e06c75 !important; }
:deep(.hljs-literal) { color: #56b6c2 !important; }
:deep(.hljs-attr) { color: #d19a66 !important; }
:deep(.hljs-variable) { color: #e06c75 !important; }
:deep(.hljs-type) { color: #e5c07b !important; }
:deep(.hljs-params) { color: #abb2bf !important; }
:deep(.hljs-meta) { color: #61afef !important; }
:deep(.hljs-link) { color: #61afef !important; text-decoration: underline !important; }
:deep(.hljs-tag) { color: #e06c75 !important; }
:deep(.hljs-name) { color: #e06c75 !important; }
:deep(.hljs-attribute) { color: #d19a66 !important; }
:deep(.hljs-bullet) { color: #61afef !important; }
:deep(.hljs-code) { color: #abb2bf !important; }
:deep(.hljs-emphasis) { font-style: italic !important; }
:deep(.hljs-strong) { font-weight: bold !important; }
:deep(.hljs-formula) { color: #c678dd !important; }
:deep(.hljs-quote) { color: #5c6370 !important; font-style: italic !important; }
:deep(.hljs-doctag) { color: #c678dd !important; }
:deep(.hljs-selector-tag) { color: #e06c75 !important; }
:deep(.hljs-selector-id) { color: #61afef !important; }
:deep(.hljs-selector-class) { color: #d19a66 !important; }
:deep(.hljs-section) { color: #e06c75 !important; font-weight: bold !important; }
:deep(.hljs-regexp) { color: #56b6c2 !important; }
:deep(.hljs-symbol) { color: #61afef !important; }
:deep(.hljs-class) { color: #e5c07b !important; }

/* 自定义滚动条 */
.code-textarea::-webkit-scrollbar,
.code-preview::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.code-textarea::-webkit-scrollbar-track,
.code-preview::-webkit-scrollbar-track {
  background: var(--bg-color);
  border-radius: 4px;
}

.code-textarea::-webkit-scrollbar-thumb,
.code-preview::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
}

.code-textarea::-webkit-scrollbar-thumb:hover,
.code-preview::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>
