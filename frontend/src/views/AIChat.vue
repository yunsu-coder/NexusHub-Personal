<template>
  <div class="ai-chat-page">
    <div class="chat-container glass-panel">
      <!-- Chat Header -->
      <div class="chat-header">
        <div class="header-left">
          <h3>Nexus AI</h3>
          <el-tag size="small" effect="plain" round>{{ aiSettings.model }}</el-tag>
        </div>
        <div class="header-right">
          <el-tooltip content="清空会话">
            <el-button circle :icon="Delete" size="small" @click="clearCurrentSession" />
          </el-tooltip>
          <el-tooltip content="设置">
            <el-button circle :icon="Setting" size="small" @click="showSettingsDrawer = true" />
          </el-tooltip>
        </div>
      </div>

      <!-- Chat Messages -->
      <div class="chat-messages" ref="messagesRef">
        <div v-if="currentSession.messages.length === 0" class="empty-state">
          <div class="ai-avatar-large">
             <div class="pulse-ring"></div>
             <el-icon><Cpu /></el-icon>
          </div>
          <h2>How can I help you today?</h2>
          <div class="suggestion-chips">
            <div class="chip" @click="inputMessage = '帮我写一份周报'">
              <el-icon><Document /></el-icon> 写周报
            </div>
            <div class="chip" @click="inputMessage = '解释一下量子纠缠'">
              <el-icon><Reading /></el-icon> 解释概念
            </div>
            <div class="chip" @click="inputMessage = '用Python写一个HTTP服务器'">
              <el-icon><Monitor /></el-icon> 写代码
            </div>
          </div>
        </div>
        
        <div v-else class="message-list">
           <div 
             v-for="(message, index) in currentSession.messages" 
             :key="message.id || index" 
             class="message-row" 
             :class="message.role"
           >
             <div class="message-avatar">
               <el-avatar v-if="message.role === 'user'" :size="36" :src="aiSettings.userAvatar">User</el-avatar>
               <div v-else class="ai-avatar-small"><el-icon><Cpu /></el-icon></div>
             </div>
             
             <div class="message-content-wrapper">
               <div class="message-sender">{{ message.role === 'user' ? 'You' : 'Nexus AI' }}</div>
               <div class="message-bubble" :class="message.role">
                 <div class="markdown-body" v-html="renderMarkdown(message.content)"></div>
               </div>
               <div class="message-actions">
                  <el-button link size="small" @click="copyMessage(message.content)">复制</el-button>
               </div>
             </div>
           </div>
           
           <div v-if="isTyping" class="message-row assistant">
             <div class="message-avatar">
               <div class="ai-avatar-small"><el-icon><Cpu /></el-icon></div>
             </div>
             <div class="message-content-wrapper">
                <div class="typing-indicator">
                  <span></span><span></span><span></span>
                </div>
             </div>
           </div>
        </div>
      </div>

      <!-- Input Area -->
      <div class="input-area">
        <div class="input-box glass-panel">
          <el-input
            v-model="inputMessage"
            type="textarea"
            :autosize="{ minRows: 1, maxRows: 6 }"
            placeholder="Send a message..."
            resize="none"
            class="chat-input"
            @keydown.enter.exact.prevent="sendMessage"
          />
          <div class="input-actions">
             <el-button text circle :icon="Picture" @click="triggerUpload" />
             <el-button type="primary" circle :icon="Position" :loading="loading" @click="sendMessage" :disabled="!inputMessage.trim() && !loading" />
          </div>
        </div>
        <div class="input-footer">
           Nexus AI can make mistakes. Consider checking important information.
        </div>
      </div>
    </div>

    <!-- Settings Drawer -->
    <el-drawer v-model="showSettingsDrawer" title="AI Configuration" size="350px">
      <el-form :model="aiSettings" label-position="top">
        <el-form-item label="API Endpoint">
          <el-input v-model="aiSettings.apiUrl" placeholder="https://api.openai.com/v1/chat/completions" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="aiSettings.apiKey" type="password" show-password />
        </el-form-item>
        <el-form-item label="Model">
          <el-input v-model="aiSettings.model" placeholder="gpt-3.5-turbo" />
        </el-form-item>
        <el-form-item label="Temperature">
          <el-slider v-model="aiSettings.temperature" :min="0" :max="2" :step="0.1" />
        </el-form-item>
      </el-form>
      <template #footer>
         <el-button type="primary" @click="saveAISettings">Save Settings</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick, computed } from 'vue'
import { Delete, Setting, Position, Picture, Cpu, Document, Reading, Monitor } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'

marked.setOptions({
  highlight: function (code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      return hljs.highlight(code, { language: lang }).value
    }
    return hljs.highlightAuto(code).value
  }
})

const messagesRef = ref(null)
const inputMessage = ref('')
const loading = ref(false)
const isTyping = ref(false)
const showSettingsDrawer = ref(false)
const sessions = ref([])
const currentSessionId = ref(null)

const aiSettings = ref({
  apiUrl: 'https://api.openai.com/v1/chat/completions',
  apiKey: '',
  model: 'gpt-3.5-turbo',
  temperature: 0.7,
  userAvatar: ''
})

const currentSession = computed(() => {
  return sessions.value.find(s => s.id === currentSessionId.value) || { messages: [] }
})

const loadHistory = () => {
  const saved = localStorage.getItem('chat_sessions')
  if (saved) {
    sessions.value = JSON.parse(saved)
    if (sessions.value.length > 0) currentSessionId.value = sessions.value[0].id
    else createNewSession()
  } else {
    createNewSession()
  }
}

const createNewSession = () => {
  const newSession = {
    id: Date.now().toString(),
    messages: [],
    created_at: new Date().toISOString()
  }
  sessions.value.unshift(newSession)
  currentSessionId.value = newSession.id
}

const loadSettings = () => {
  const saved = localStorage.getItem('ai_settings')
  if (saved) aiSettings.value = JSON.parse(saved)
}

const saveAISettings = () => {
  localStorage.setItem('ai_settings', JSON.stringify(aiSettings.value))
  ElMessage.success('Settings saved')
  showSettingsDrawer.value = false
}

const sendMessage = async () => {
  if (!inputMessage.value.trim() || loading.value) return
  
  const userMsg = inputMessage.value
  inputMessage.value = ''
  
  // 先追加用户消息
  currentSession.value.messages.push({
      role: 'user',
      content: userMsg,
      created_at: new Date().toISOString()
  })
  scrollToBottom()

  // 如果没有配置 API Key，走本地假回复模式，保证可用性
  if (!aiSettings.value.apiKey) {
      const assistantMsg = {
        role: 'assistant',
        content: `（本地离线模式）\n\n你说的是：\n${userMsg}\n\n当前未配置 AI Key，请在右上角“设置”中填写后即可接入真实模型。`,
        created_at: new Date().toISOString()
      }
      currentSession.value.messages.push(assistantMsg)
      localStorage.setItem('chat_sessions', JSON.stringify(sessions.value))
      ElMessage.info('当前处于本地离线演示模式，未调用真实 AI 接口')
      scrollToBottom()
      return
  }

  loading.value = true
  isTyping.value = true
  
  // 准备 Assistant 消息（用于流式追加内容）
  const assistantMsg = {
      role: 'assistant',
      content: '',
      created_at: new Date().toISOString()
  }
  currentSession.value.messages.push(assistantMsg)
  
  try {
    const history = currentSession.value.messages.slice(0, -1).map(m => ({
        role: m.role,
        content: m.content
    }))
    
    const response = await fetch(aiSettings.value.apiUrl, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${aiSettings.value.apiKey}`
        },
        body: JSON.stringify({
            model: aiSettings.value.model,
            messages: history,
            temperature: aiSettings.value.temperature,
            stream: true
        })
    })

    if (!response.ok) {
        const errorText = await response.text().catch(() => '')
        throw new Error(`API Error: ${response.status}${errorText ? ' - ' + errorText : ''}`)
    }

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    
    while(true) {
        const { done, value } = await reader.read()
        if (done) break
        
        const chunk = decoder.decode(value)
        const lines = chunk.split('\n')
        
        for (const line of lines) {
            if (line.startsWith('data: ')) {
                const data = line.slice(6)
                if (data === '[DONE]') continue
                try {
                    const json = JSON.parse(data)
                    const delta = json.choices[0]?.delta?.content || ''
                    assistantMsg.content += delta
                    scrollToBottom()
                } catch(e) {}
            }
        }
    }
  } catch (e) {
      const msg = e?.message || 'Unknown error'
      assistantMsg.content += `\n[Error: ${msg}]`
      ElMessage.error(`AI 请求失败：${msg}`)
  } finally {
      // 统一持久化会话
      localStorage.setItem('chat_sessions', JSON.stringify(sessions.value))
      loading.value = false
      isTyping.value = false
  }
}

const scrollToBottom = () => {
    nextTick(() => {
        if (messagesRef.value) {
            messagesRef.value.scrollTop = messagesRef.value.scrollHeight
        }
    })
}

const clearCurrentSession = () => {
    currentSession.value.messages = []
    localStorage.setItem('chat_sessions', JSON.stringify(sessions.value))
}

const triggerUpload = () => {
    ElMessage.info('File upload is a demo feature')
}

const copyMessage = (text) => {
    navigator.clipboard.writeText(text)
    ElMessage.success('Copied')
}

const renderMarkdown = (text) => marked(text)

onMounted(() => {
    loadSettings()
    loadHistory()
})
</script>

<style scoped>
.ai-chat-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 20px;
}

.chat-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #fff; /* Fallback */
}

.glass-panel {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

@media (prefers-color-scheme: dark) {
    .glass-panel {
        background: rgba(30, 30, 30, 0.95);
        color: #eee;
    }
}

.chat-header {
  padding: 15px 20px;
  border-bottom: 1px solid rgba(0,0,0,0.05);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-left { display: flex; align-items: center; gap: 10px; }
.header-left h3 { margin: 0; font-size: 16px; }

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  opacity: 0.8;
}

.ai-avatar-large {
  width: 80px;
  height: 80px;
  background: #000;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40px;
  color: #fff;
  margin-bottom: 20px;
  position: relative;
}

.pulse-ring {
    position: absolute;
    width: 100%;
    height: 100%;
    border-radius: 50%;
    border: 2px solid #000;
    animation: pulse 2s infinite;
}

@keyframes pulse {
    0% { transform: scale(1); opacity: 0.5; }
    100% { transform: scale(1.5); opacity: 0; }
}

.suggestion-chips {
  display: flex;
  gap: 10px;
  margin-top: 30px;
}
.chip {
  padding: 10px 20px;
  border: 1px solid #ddd;
  border-radius: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  transition: all 0.2s;
}
.chip:hover {
  background: #f0f0f0;
  transform: translateY(-2px);
}

.message-row {
  display: flex;
  gap: 15px;
  margin-bottom: 25px;
  max-width: 800px;
  margin-left: auto;
  margin-right: auto;
  width: 100%;
}

.message-avatar {
    flex-shrink: 0;
    margin-top: 5px;
}
.ai-avatar-small {
    width: 36px;
    height: 36px;
    background: #000;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
}

.message-content-wrapper {
    flex: 1;
    min-width: 0;
}
.message-sender { font-size: 12px; margin-bottom: 5px; opacity: 0.6; font-weight: 600; }

.message-bubble {
    padding: 0 15px;
    line-height: 1.6;
}
.message-row.user .message-bubble {
    /* No bubble for user, plain text or slight bg? */
    /* Let's keep it simple like ChatGPT */
}

.markdown-body :deep(pre) {
    background: #1e1e1e;
    padding: 15px;
    border-radius: 8px;
    overflow-x: auto;
}
.markdown-body :deep(code) {
    font-family: 'Consolas', monospace;
}

.message-actions {
    opacity: 0;
    transition: opacity 0.2s;
}
.message-row:hover .message-actions { opacity: 1; }

.input-area {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
}

.input-box {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 16px;
  background: #fff;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
}

.chat-input :deep(.el-textarea__inner) {
    box-shadow: none;
    border: none;
    background: transparent;
    padding: 5px 10px;
}

.input-actions {
    display: flex;
    justify-content: space-between;
    padding: 5px 10px 0;
}

.input-footer {
    text-align: center;
    font-size: 11px;
    opacity: 0.5;
    margin-top: 10px;
}

.typing-indicator span {
    display: inline-block;
    width: 6px;
    height: 6px;
    background: #aaa;
    border-radius: 50%;
    margin-right: 4px;
    animation: typing 1s infinite;
}
.typing-indicator span:nth-child(2) { animation-delay: 0.2s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.4s; }

@keyframes typing {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-5px); }
}
</style>
