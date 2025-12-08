<template>
  <div class="ai-chat-page">
    <el-card class="chat-container">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>🤖 AI 助手</span>
          <el-button size="small" @click="clearHistory">清空历史</el-button>
        </div>
      </template>

      <div class="messages-container" ref="messagesRef">
        <div v-for="(msg, index) in messages" :key="index" class="message-item" :class="`message-${msg.role}`">
          <div class="message-avatar">
            {{ msg.role === 'user' ? '👤' : '🤖' }}
          </div>
          <div class="message-content">
            <div class="message-text">{{ msg.content }}</div>
            <div class="message-time">{{ new Date(msg.created_at).toLocaleTimeString('zh-CN') }}</div>
          </div>
        </div>
        <el-empty v-if="messages.length === 0" description="开始对话吧！" />
      </div>

      <div class="input-container">
        <el-input
          v-model="inputMessage"
          type="textarea"
          :rows="3"
          placeholder="输入消息... (Enter发送, Ctrl+Enter换行)"
          @keydown="handleKeyDown"
        />
        <el-button type="primary" :loading="loading" @click="sendMessage" style="margin-top: 10px">
          发送 (Enter)
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

const messages = ref([])
const inputMessage = ref('')
const loading = ref(false)
const messagesRef = ref(null)
const aiSettings = ref(null)

const loadAISettings = () => {
  const saved = localStorage.getItem('ai_settings')
  if (saved) {
    try {
      aiSettings.value = JSON.parse(saved)
    } catch {
      console.error('Failed to parse AI settings')
    }
  }
}

const loadHistory = async () => {
  try {
    const history = await api.get('/chat/history?limit=50')
    messages.value = history.reverse()
    scrollToBottom()
  } catch {
    // 如果后端获取失败，从localStorage加载
    const saved = localStorage.getItem('chat_history')
    if (saved) {
      try {
        messages.value = JSON.parse(saved)
      } catch (e) {
        console.error('Failed to load chat history', e)
      }
    }
  }
}

const saveToLocalStorage = () => {
  try {
    localStorage.setItem('chat_history', JSON.stringify(messages.value))
  } catch (e) {
    console.error('Failed to save chat history', e)
  }
}

const sendMessage = async () => {
  if (!inputMessage.value.trim()) return

  // 检查AI配置
  if (!aiSettings.value?.apiKey) {
    ElMessage.warning('请先在设置中配置 AI API Key')
    return
  }

  loading.value = true
  const userMsg = inputMessage.value
  inputMessage.value = ''

  // 添加用户消息
  const userMessage = {
    role: 'user',
    content: userMsg,
    created_at: new Date().toISOString()
  }
  messages.value.push(userMessage)
  scrollToBottom()

  try {
    // 准备对话历史（最近10条）
    const conversationHistory = messages.value.slice(-10).map(m => ({
      role: m.role === 'assistant' ? 'assistant' : 'user',
      content: m.content
    }))

    // 调用AI API
    const response = await fetch(aiSettings.value.apiUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${aiSettings.value.apiKey}`
      },
      body: JSON.stringify({
        model: aiSettings.value.model,
        messages: conversationHistory,
        max_tokens: 1000,
        temperature: 0.7
      })
    })

    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`)
    }

    const data = await response.json()
    const aiResponse = data.choices?.[0]?.message?.content || '抱歉，我无法回答这个问题。'

    // 添加AI回复
    const assistantMessage = {
      role: 'assistant',
      content: aiResponse,
      created_at: new Date().toISOString()
    }
    messages.value.push(assistantMessage)

    // 保存到后端和localStorage
    try {
      await api.post('/chat/message', { content: userMsg })
    } catch {
      // 后端保存失败也没关系，继续使用localStorage
    }

    saveToLocalStorage()
    scrollToBottom()
    ElMessage.success('发送成功')
  } catch (error) {
    ElMessage.error(`发送失败: ${error.message}`)
    // 移除失败的用户消息
    messages.value.pop()
  } finally {
    loading.value = false
  }
}

const clearHistory = async () => {
  try {
    await ElMessageBox.confirm('确定清空所有对话历史吗？', '提示', { type: 'warning' })

    // 清空后端
    try {
      await api.delete('/chat/history')
    } catch {
      // 忽略后端错误
    }

    // 清空localStorage
    localStorage.removeItem('chat_history')
    messages.value = []
    ElMessage.success('清空成功')
  } catch(e) {
    if(e !== 'cancel') ElMessage.error('清空失败')
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  })
}

const handleKeyDown = (event) => {
  // Ctrl+Enter 换行（默认行为，不需要处理）
  // Enter 发送消息
  if (event.key === 'Enter' && !event.ctrlKey) {
    event.preventDefault() // 阻止默认换行
    sendMessage()
  }
}

onMounted(() => {
  loadAISettings()
  loadHistory()
})
</script>

<style scoped>
.ai-chat-page {
  height: calc(100vh - 120px);
}

.chat-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.chat-container :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.message-item {
  display: flex;
  gap: 12px;
  animation: fadeIn 0.3s;
}

.message-user {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  background: var(--hover-bg);
  flex-shrink: 0;
}

.message-content {
  max-width: 70%;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.message-user .message-content {
  align-items: flex-end;
}

.message-text {
  padding: 12px 16px;
  border-radius: 12px;
  background: var(--card-bg);
  color: var(--text-primary);
  line-height: 1.6;
  word-wrap: break-word;
  border: 1px solid var(--border-color);
}

.message-user .message-text {
  background: #409EFF;
  color: white;
  border-color: #409EFF;
}

.message-time {
  font-size: 12px;
  color: var(--text-secondary);
  padding: 0 4px;
}

.input-container {
  padding: 20px;
  border-top: 1px solid var(--border-color);
  background: var(--card-bg);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
