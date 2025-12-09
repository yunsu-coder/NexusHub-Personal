<template>
  <div class="ai-chat-page">
    <el-container class="chat-container">
      <!-- 左侧会话列表 -->
      <el-aside width="280px" class="session-list">
        <div class="session-header">
          <h3>AI 会话</h3>
          <el-button type="primary" circle @click="createNewSession" :icon="Plus" size="small" />
        </div>
        <div class="sessions">
          <div 
            v-for="session in sessions" 
            :key="session.id"
            class="session-item"
            :class="{ active: session.id === currentSessionId }"
            @click="switchSession(session.id)"
          >
            <div class="session-info">
              <div class="session-title">{{ session.title }}</div>
              <div class="session-time">{{ formatSessionTime(session.updated_at) }}</div>
            </div>
            <el-button 
              type="danger" 
              text 
              circle 
              size="small"
              :icon="Delete"
              @click.stop="deleteSession(session.id)"
            />
          </div>
        </div>
      </el-aside>

      <!-- 右侧聊天区域 -->
      <el-main class="chat-main">
        <div class="chat-header">
          <div class="session-actions">
            <h3>{{ currentSession.title }}</h3>
            <el-button 
              type="primary" 
              text 
              size="small"
              @click="renameSession"
            >
              重命名
            </el-button>
            <el-button 
              type="warning" 
              text 
              size="small"
              @click="clearCurrentSession"
            >
              清空会话
            </el-button>
            <el-button 
              type="success" 
              text 
              size="small"
              @click="showSettings = true"
            >
              设置
            </el-button>
          </div>
        </div>

        <!-- 聊天消息区域 -->
        <div class="chat-messages" ref="messagesRef">
          <div v-for="(message, index) in currentSession.messages" :key="message.id || index" class="message">
            <div :class="['message-content', message.role]">
              <div class="message-header">
                <span class="message-role">{{ message.role === 'assistant' ? 'AI' : '我' }}</span>
                <span class="message-time">{{ formatMessageTime(message.created_at) }}</span>
              </div>
              <div 
                class="markdown-content" 
                v-html="renderMarkdown(message.content)"
              ></div>
              <div class="message-actions">
                <el-button 
                  type="text" 
                  size="small"
                  @click="copyMessage(message.content)"
                >
                  复制
                </el-button>
                <el-button 
                  type="text" 
                  size="small"
                  @click="replyToMessage(message)"
                >
                  回复
                </el-button>
              </div>
            </div>
          </div>
          <div v-if="isTyping" class="typing-indicator">
            <el-skeleton :rows="3" animated />
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="input-container">
          <el-input
            v-model="inputMessage"
            type="textarea"
            :rows="2"
            placeholder="输入消息... (Enter 发送, Ctrl+Enter 换行)"
            resize="none"
            @keydown="handleKeyDown"
          />
          <div class="input-actions">
            <el-button type="primary" :loading="loading" @click="sendMessage">
              发送
            </el-button>
          </div>
        </div>
      </el-main>
    </el-container>

    <!-- 设置面板 -->
    <el-drawer
      v-model="showSettings"
      title="AI 设置"
      size="300px"
      direction="rtl"
    >
      <el-form :model="aiSettings" label-width="100px">
        <el-form-item label="API URL">
          <el-input v-model="aiSettings.apiUrl" placeholder="例如: https://api.openai.com/v1/chat/completions" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="aiSettings.apiKey" type="password" show-password />
        </el-form-item>
        <el-form-item label="模型">
          <el-select v-model="aiSettings.model" placeholder="选择模型">
            <el-option label="gpt-4o" value="gpt-4o" />
            <el-option label="gpt-3.5-turbo" value="gpt-3.5-turbo" />
            <el-option label="claude-3-opus-20240229" value="claude-3-opus-20240229" />
            <el-option label="gemini-pro" value="gemini-pro" />
          </el-select>
        </el-form-item>
        <el-form-item label="温度">
          <el-slider v-model="aiSettings.temperature" :min="0" :max="2" :step="0.1" />
        </el-form-item>
        <el-form-item label="最大响应长度">
          <el-input-number v-model="aiSettings.maxResponseLength" :min="100" :max="2000" :step="100" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveAISettings">保存设置</el-button>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick, computed } from 'vue'
import { ElMessage, ElMessageBox, ElDrawer, ElForm, ElFormItem, ElInput, ElSelect, ElOption, ElSlider, ElInputNumber, ElSkeleton } from 'element-plus'
import { Delete, Plus } from '@element-plus/icons-vue'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'
import { marked } from 'marked'
import api from '../api'

marked.setOptions({
  highlight: function (code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      return hljs.highlight(code, { language: lang }).value
    }
    return code
  }
})

// 引用
const messagesRef = ref(null)

// 响应式数据
const inputMessage = ref('')
const loading = ref(false)
const isTyping = ref(false)
const showSettings = ref(false)

// 会话管理
const sessions = ref([])
const currentSessionId = ref(null)
const aiSettings = ref({
  apiUrl: '',
  apiKey: '',
  model: 'gpt-3.5-turbo',
  temperature: 0.7,
  maxResponseLength: 1000
})

// 计算属性
const currentSession = computed(() => {
  return sessions.value.find(s => s.id === currentSessionId.value) || {
    title: '新会话',
    messages: []
  }
})

// 页面挂载
onMounted(() => {
  loadHistory()
  loadAISettings()
})

// 加载历史会话
const loadHistory = async () => {
  // 尝试从旧格式加载历史记录
  const oldHistory = localStorage.getItem('chat_history')
  if (oldHistory) {
    try {
      const oldMessages = JSON.parse(oldHistory)
      if (oldMessages.length > 0) {
        // 创建一个新会话并导入旧消息
        const importSession = {
          id: 'imported_' + Date.now().toString(),
          title: '导入的历史记录',
          messages: oldMessages,
          created_at: oldMessages[0]?.created_at || new Date().toISOString(),
          updated_at: oldMessages[oldMessages.length - 1]?.created_at || new Date().toISOString()
        }
        sessions.value.unshift(importSession)
        currentSessionId.value = importSession.id
        saveSessionsToLocalStorage()
        
        // 删除旧格式数据
        localStorage.removeItem('chat_history')
        return
      }
    } catch (e) {
      console.error('Failed to load old chat history', e)
    }
  }
  
  // 正常加载会话
  loadSessionsFromLocalStorage()
}

// 加载AI设置
const loadAISettings = () => {
  try {
    const saved = localStorage.getItem('ai_settings')
    if (saved) {
      aiSettings.value = JSON.parse(saved)
    } else {
      // 默认设置
      aiSettings.value = {
        apiUrl: 'https://api.openai.com/v1/chat/completions',
        apiKey: '',
        model: 'gpt-3.5-turbo',
        temperature: 0.7,
        maxResponseLength: 1000
      }
      saveAISettings()
    }
  } catch (e) {
    console.error('Failed to load AI settings', e)
  }
}

// 保存AI设置
const saveAISettings = () => {
  try {
    localStorage.setItem('ai_settings', JSON.stringify(aiSettings.value))
    ElMessage.success('设置已保存')
    showSettings.value = false
  } catch (e) {
    console.error('Failed to save AI settings', e)
    ElMessage.error('保存设置失败')
  }
}

// 会话管理函数
const createNewSession = () => {
  const newSession = {
    id: Date.now().toString(),
    title: `新会话 ${sessions.value.length + 1}`,
    messages: [],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString()
  }
  sessions.value.unshift(newSession)
  currentSessionId.value = newSession.id
  saveSessionsToLocalStorage()
}

const switchSession = (sessionId) => {
  currentSessionId.value = sessionId
  nextTick(() => {
    scrollToBottom()
  })
}

const deleteSession = async (sessionId) => {
  try {
    await ElMessageBox.confirm('确定要删除这个会话吗？', '提示', { type: 'warning' })
    const index = sessions.value.findIndex(s => s.id === sessionId)
    if (index !== -1) {
      sessions.value.splice(index, 1)
      if (sessionId === currentSessionId.value) {
        if (sessions.value.length > 0) {
          currentSessionId.value = sessions.value[0].id
        } else {
          createNewSession()
        }
      }
      saveSessionsToLocalStorage()
      ElMessage.success('会话已删除')
    }
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('删除会话失败')
    }
  }
}

const renameSession = async () => {
  try {
    const { value } = await ElMessageBox.prompt('请输入新的会话名称', '重命名会话', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputValue: currentSession.value.title
    })
    if (value.trim()) {
      const session = sessions.value.find(s => s.id === currentSessionId.value)
      if (session) {
        session.title = value.trim()
        saveSessionsToLocalStorage()
        ElMessage.success('会话已重命名')
      }
    }
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('重命名失败')
    }
  }
}

const clearCurrentSession = async () => {
  try {
    await ElMessageBox.confirm('确定要清空当前会话的所有消息吗？', '提示', { type: 'warning' })
    const session = sessions.value.find(s => s.id === currentSessionId.value)
    if (session) {
      session.messages = []
      session.updated_at = new Date().toISOString()
      saveSessionsToLocalStorage()
      ElMessage.success('会话已清空')
    }
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('清空会话失败')
    }
  }
}

// 格式化时间
const formatSessionTime = (timeStr) => {
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now - date
  
  if (diff < 60000) { // 1分钟内
    return '刚刚'
  } else if (diff < 3600000) { // 1小时内
    return `${Math.floor(diff / 60000)}分钟前`
  } else if (diff < 86400000) { // 1天内
    return `${Math.floor(diff / 3600000)}小时前`
  } else {
    return date.toLocaleDateString()
  }
}

const formatMessageTime = (timeStr) => {
  return new Date(timeStr).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

// 存储管理
const saveSessionsToLocalStorage = () => {
  try {
    localStorage.setItem('chat_sessions', JSON.stringify(sessions.value))
  } catch (e) {
    console.error('Failed to save sessions', e)
  }
}

const loadSessionsFromLocalStorage = () => {
  const saved = localStorage.getItem('chat_sessions')
  if (saved) {
    try {
      sessions.value = JSON.parse(saved)
      if (sessions.value.length > 0) {
        currentSessionId.value = sessions.value[0].id
      } else {
        createNewSession()
      }
    } catch (e) {
      console.error('Failed to load sessions', e)
      createNewSession()
    }
  } else {
    createNewSession()
  }
}

// 发送消息
const sendMessage = async () => {
  if (!inputMessage.value.trim() || !currentSessionId.value) return

  // 检查AI配置
  if (!aiSettings.value?.apiKey) {
    ElMessage.warning('请先在设置中配置 AI API Key')
    return
  }

  loading.value = true
  isTyping.value = true
  const userMsg = inputMessage.value
  inputMessage.value = ''

  // 添加用户消息
  const userMessage = {
    id: Date.now().toString(),
    role: 'user',
    content: userMsg,
    created_at: new Date().toISOString()
  }
  currentSession.value.messages.push(userMessage)
  scrollToBottom()

  try {
    // 准备对话历史
    const conversationHistory = currentSession.value.messages.map(m => ({
      role: m.role === 'assistant' ? 'assistant' : 'user',
      content: m.content
    }))

    // 调用AI API（流式）
    const response = await fetch(aiSettings.value.apiUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${aiSettings.value.apiKey}`
      },
      body: JSON.stringify({
        model: aiSettings.value.model,
        messages: conversationHistory,
        max_tokens: aiSettings.value.maxResponseLength,
        temperature: aiSettings.value.temperature,
        stream: true
      })
    })

    if (!response.ok) {
      throw new Error(`API请求失败: ${response.status}`)
    }

    // 创建AI回复消息
    const assistantMessage = {
      id: Date.now().toString(),
      role: 'assistant',
      content: '',
      created_at: new Date().toISOString()
    }
    currentSession.value.messages.push(assistantMessage)
    scrollToBottom()

    // 处理流式响应
    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = line.slice(6)
          if (data === '[DONE]') continue
          
          try {
            const json = JSON.parse(data)
            const delta = json.choices?.[0]?.delta?.content || ''
            if (delta) {
              assistantMessage.content += delta
              scrollToBottom()
            }
          } catch (e) {
            console.error('Failed to parse stream data', e)
          }
        }
      }
    }

    // 更新会话时间并保存
    currentSession.value.updated_at = new Date().toISOString()
    saveSessionsToLocalStorage()
    
    // 保存到后端
    try {
      await api.post('/chat/message', { content: userMsg })
    } catch {
      // 忽略后端保存失败
    }
    
    ElMessage.success('发送成功')
  } catch (error) {
    ElMessage.error(`发送失败: ${error.message}`)
    // 移除失败的用户消息和空的AI回复
    if (currentSession.value.messages.length > 0 && currentSession.value.messages[currentSession.value.messages.length - 1].role === 'assistant') {
      currentSession.value.messages.pop()
    }
    if (currentSession.value.messages.length > 0 && currentSession.value.messages[currentSession.value.messages.length - 1].role === 'user') {
      currentSession.value.messages.pop()
    }
  } finally {
    loading.value = false
    isTyping.value = false
  }
}

// 回复消息
const replyToMessage = (message) => {
  inputMessage.value = `回复: ${message.content}\n\n`
  nextTick(() => {
    const inputElement = document.querySelector('.el-textarea__inner')
    if (inputElement) {
      inputElement.focus()
      inputElement.setSelectionRange(inputMessage.value.length, inputMessage.value.length)
    }
  })
}

// 复制消息内容
const copyMessage = async (content) => {
  try {
    await navigator.clipboard.writeText(content)
    ElMessage.success('复制成功')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  })
}

// 处理键盘事件
const handleKeyDown = (event) => {
  // Enter 发送消息
  if (event.key === 'Enter' && !event.ctrlKey && !loading.value) {
    event.preventDefault() // 阻止默认换行
    sendMessage()
  }
}

// 渲染Markdown
const renderMarkdown = (content) => {
  return marked(content)
}

// 监听消息变化，确保新消息也能正确高亮
watch(
  () => currentSession.value.messages,
  () => {
    nextTick(() => {
      // 重新应用语法高亮到所有代码块
      const codeBlocks = document.querySelectorAll('.markdown-content pre code')
      codeBlocks.forEach((block) => {
        if (!block.classList.contains('hljs')) {
          hljs.highlightElement(block)
        }
      })
    })
  },
  { deep: true }
)
</script>

<style scoped>
.ai-chat-page {
  height: calc(100vh - 120px);
  width: 100%;
  overflow: hidden;
}

.chat-container {
  height: 100%;
}

/* 左侧会话列表 */
.session-list {
  border-right: 1px solid #e0e0e0;
  background-color: #f5f7fa;
  padding: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.session-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #e0e0e0;
}

.session-header h3 {
  margin: 0;
  font-size: 16px;
}

.sessions {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.session-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  margin-bottom: 4px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  background-color: white;
}

.session-item:hover {
  background-color: #ecf5ff;
}

.session-item.active {
  background-color: #409eff;
  color: white;
}

.session-info {
  flex: 1;
  min-width: 0;
  margin-right: 8px;
}

.session-title {
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.session-time {
  font-size: 12px;
  opacity: 0.7;
}

/* 右侧聊天区域 */
.chat-main {
  padding: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.chat-header {
  padding: 12px 16px;
  border-bottom: 1px solid #e0e0e0;
  background-color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.session-actions h3 {
  margin: 0;
  font-size: 16px;
  display: inline-block;
  margin-right: 16px;
  vertical-align: middle;
}

/* 聊天消息区域 */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background-color: #f5f7fa;
}

.message {
  margin-bottom: 16px;
  display: flex;
}

.message-content {
  max-width: 80%;
  padding: 12px 16px;
  border-radius: 12px;
  position: relative;
}

.message-content.user {
  margin-left: auto;
  background-color: #409eff;
  color: white;
  border-bottom-right-radius: 4px;
}

.message-content.assistant {
  background-color: white;
  border-bottom-left-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.message-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 12px;
  opacity: 0.8;
}

.message-role {
  font-weight: 500;
}

.message-time {
  font-size: 11px;
}

.markdown-content {
  line-height: 1.6;
  word-break: break-word;
}

.markdown-content p {
  margin: 0 0 8px 0;
}

.markdown-content pre {
  margin: 8px 0;
  padding: 12px;
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 8px;
  overflow-x: auto;
}

.markdown-content code {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
}

.message-actions {
  margin-top: 8px;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.message-actions .el-button {
  padding: 2px 8px;
  font-size: 12px;
}

/* 输入区域 */
.input-container {
  padding: 16px;
  background-color: white;
  border-top: 1px solid #e0e0e0;
}

.input-actions {
  margin-top: 8px;
  display: flex;
  justify-content: flex-end;
}

/* 打字指示器 */
.typing-indicator {
  padding: 12px 16px;
  background-color: white;
  border-radius: 12px;
  border-bottom-left-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  max-width: 80%;
}
</style>