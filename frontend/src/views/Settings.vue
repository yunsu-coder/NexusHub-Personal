<template>
  <div class="settings-page">
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>🎨 主题设置</span>
          </template>

          <div class="setting-item">
            <label>主题模式</label>
            <el-radio-group v-model="theme.theme_name" @change="saveTheme">
              <el-radio-button label="dark">深色</el-radio-button>
              <el-radio-button label="light">浅色</el-radio-button>
            </el-radio-group>
          </div>

          <div class="setting-item">
            <label>主色调</label>
            <el-color-picker v-model="theme.primary_color" @change="saveTheme" show-alpha />
          </div>

          <div class="setting-item">
            <label>辅助色</label>
            <el-color-picker v-model="theme.secondary_color" @change="saveTheme" show-alpha />
          </div>

          <div class="setting-item">
            <label>背景图片</label>
            <div style="display: flex; gap: 10px; margin-bottom: 10px">
              <el-input v-model="theme.background_image" placeholder="输入图片 URL..." @change="saveTheme" style="flex: 1" />
              <el-button @click="previewBackgroundImage" :disabled="!theme.background_image">预览</el-button>
            </div>
            <el-upload
              :action="`http://localhost:8080/api/v1/files/upload`"
              :headers="uploadHeaders"
              :show-file-list="false"
              accept="image/*"
              :on-success="handleImageUploadSuccess"
              :before-upload="beforeImageUpload"
            >
              <el-button type="primary" size="small">本地上传图片</el-button>
            </el-upload>
          </div>

          <div class="setting-item">
            <label>背景图片透明度: {{ Math.round(theme.background_opacity * 100) }}%</label>
            <el-slider v-model="theme.background_opacity" :min="0" :max="1" :step="0.05" @change="saveTheme" />
          </div>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>
            <span>🤖 AI 聊天设置</span>
          </template>

          <div class="setting-item">
            <label>AI API Key</label>
            <el-input
              v-model="aiSettings.apiKey"
              type="password"
              placeholder="输入 API Key..."
              show-password
              @change="saveAISettings"
            />
          </div>

          <div class="setting-item">
            <label>API URL</label>
            <el-input
              v-model="aiSettings.apiUrl"
              placeholder="例如: https://api.openai.com/v1/chat/completions"
              @change="saveAISettings"
            />
          </div>

          <div class="setting-item">
            <label>AI 模型</label>
            <el-select v-model="aiSettings.model" placeholder="选择模型" @change="handleModelChange" style="width: 100%">
              <el-option label="GPT-3.5 Turbo" value="gpt-3.5-turbo" />
              <el-option label="GPT-4" value="gpt-4" />
              <el-option label="GPT-4 Turbo" value="gpt-4-turbo-preview" />
              <el-option label="Claude 3 Sonnet" value="claude-3-sonnet" />
              <el-option label="Claude 3 Opus" value="claude-3-opus" />
              <el-option label="DeepSeek Chat" value="deepseek-chat" />
              <el-option label="自定义" value="custom" />
            </el-select>
          </div>

          <div class="setting-item">
            <el-alert title="AI 配置说明" type="warning" :closable="false">
              <p>• API Key 存储在本地浏览器，不会上传到服务器</p>
              <p>• 支持 OpenAI、Claude、Gemini 等多种 API</p>
              <p>• 测试连接前请确保 API Key 有效</p>
            </el-alert>
          </div>

          <div class="setting-item">
            <el-button type="primary" @click="testAIConnection" :loading="testing">测试 AI 连接</el-button>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>
            <span>🎵 音乐设置</span>
          </template>

          <div class="setting-item">
            <label>背景音乐</label>
            <div style="display: flex; gap: 10px; margin-bottom: 10px">
              <el-input v-model="theme.background_music" placeholder="输入音乐 URL..." @change="saveTheme" style="flex: 1" />
              <el-button @click="previewBackgroundMusic" :disabled="!theme.background_music">预览</el-button>
            </div>
            <el-upload
              :action="`http://localhost:8080/api/v1/files/upload`"
              :headers="uploadHeaders"
              :show-file-list="false"
              accept="audio/*"
              :on-success="handleMusicUploadSuccess"
              :before-upload="beforeMusicUpload"
            >
              <el-button type="primary" size="small">本地上传音乐</el-button>
            </el-upload>
          </div>

          <div class="setting-item">
            <label>音量: {{ Math.round(theme.music_volume * 100) }}%</label>
            <el-slider v-model="theme.music_volume" :min="0" :max="1" :step="0.1" @change="saveTheme" />
          </div>

          <div class="setting-item" style="margin-top: 30px">
            <el-alert title="提示" type="info" :closable="false">
              <p>• 主题更改会立即生效</p>
              <p>• 背景音乐支持 MP3、WAV 等格式</p>
              <p>• 所有设置自动保存</p>
            </el-alert>
          </div>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>
            <span>ℹ️ 系统信息</span>
          </template>

          <div class="info-item">
            <span>版本</span>
            <el-tag>v1.0.0</el-tag>
          </div>

          <div class="info-item">
            <span>前端</span>
            <el-tag type="success">Vue 3 + Element Plus</el-tag>
          </div>

          <div class="info-item">
            <span>后端</span>
            <el-tag type="warning">Go + Gin</el-tag>
          </div>

          <div class="info-item">
            <span>数据库</span>
            <el-tag type="danger">MySQL</el-tag>
          </div>

          <div class="info-item">
            <span>AI 状态</span>
            <el-tag :type="aiSettings.apiKey ? 'success' : 'info'">
              {{ aiSettings.apiKey ? '已配置' : '未配置' }}
            </el-tag>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useThemeStore } from '../store/theme'
import api from '../api'

const themeStore = useThemeStore()
const testing = ref(false)

const theme = ref({
  theme_name: 'dark',
  primary_color: '#000000',
  secondary_color: '#ffffff',
  background_image: '',
  background_opacity: 1.0,
  background_music: '',
  music_volume: 0.5
})

const aiSettings = ref({
  apiKey: '',
  apiUrl: 'https://api.openai.com/v1/chat/completions',
  model: 'gpt-3.5-turbo'
})

// 上传请求头 - 包含JWT token
const uploadHeaders = computed(() => {
  const token = localStorage.getItem('token')
  return {
    Authorization: `Bearer ${token}`
  }
})

// 图片上传前验证
const beforeImageUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB!')
    return false
  }
  return true
}

// 音乐上传前验证
const beforeMusicUpload = (file) => {
  const isAudio = file.type.startsWith('audio/')
  const isLt20M = file.size / 1024 / 1024 < 20

  if (!isAudio) {
    ElMessage.error('只能上传音频文件!')
    return false
  }
  if (!isLt20M) {
    ElMessage.error('音频大小不能超过 20MB!')
    return false
  }
  return true
}

// 图片上传成功
const handleImageUploadSuccess = (response) => {
  if (response && response.url) {
    theme.value.background_image = `http://localhost:8080${response.url}`
    saveTheme()
    ElMessage.success('背景图片上传成功')
  }
}

// 音乐上传成功
const handleMusicUploadSuccess = (response) => {
  if (response && response.url) {
    theme.value.background_music = `http://localhost:8080${response.url}`
    saveTheme()
    ElMessage.success('背景音乐上传成功')
  }
}

const loadTheme = async () => {
  try {
    theme.value = await api.get('/theme')
  } catch {
    ElMessage.error('加载设置失败')
  }
}

const saveTheme = async () => {
  try {
    await themeStore.saveTheme(theme.value)
    ElMessage.success('设置已保存')
  } catch {
    ElMessage.error('保存失败')
  }
}

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

const saveAISettings = () => {
  try {
    localStorage.setItem('ai_settings', JSON.stringify(aiSettings.value))
    ElMessage.success('AI 设置已保存')
  } catch {
    ElMessage.error('保存 AI 设置失败')
  }
}

const handleModelChange = () => {
  // 自动配置 DeepSeek API
  if (aiSettings.value.model === 'deepseek-chat') {
    aiSettings.value.apiKey = 'sk-ba421705882c40bb87d018b3971faa38'
    aiSettings.value.apiUrl = 'https://api.deepseek.com/v1/chat/completions'
    ElMessage.success('DeepSeek API 已自动配置')
  }
  saveAISettings()
}

const previewBackgroundImage = () => {
  if (!theme.value.background_image) return

  const img = new Image()
  img.onload = () => {
    ElMessage.success('图片加载成功！')
    // 临时应用预览
    document.documentElement.style.setProperty('--bg-image', `url(${theme.value.background_image})`)
  }
  img.onerror = () => {
    ElMessage.error('图片加载失败，请检查 URL 是否正确')
  }
  img.src = theme.value.background_image
}

const previewBackgroundMusic = () => {
  if (!theme.value.background_music) return

  const audio = new Audio()
  audio.oncanplay = () => {
    ElMessage.success('音乐加载成功！可以在侧边栏播放')
  }
  audio.onerror = () => {
    ElMessage.error('音乐加载失败，请检查 URL 是否正确')
  }
  audio.src = theme.value.background_music
}

const testAIConnection = async () => {
  if (!aiSettings.value.apiKey) {
    ElMessage.warning('请先配置 API Key')
    return
  }

  testing.value = true
  try {
    const response = await fetch(aiSettings.value.apiUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${aiSettings.value.apiKey}`
      },
      body: JSON.stringify({
        model: aiSettings.value.model,
        messages: [{ role: 'user', content: 'Hello, this is a test message.' }],
        max_tokens: 10
      })
    })

    if (response.ok) {
      ElMessage.success('AI 连接测试成功！')
    } else {
      const error = await response.text()
      ElMessage.error(`连接失败: ${response.status} - ${error}`)
    }
  } catch (error) {
    ElMessage.error(`连接失败: ${error.message}`)
  } finally {
    testing.value = false
  }
}

onMounted(() => {
  loadTheme()
  loadAISettings()
})
</script>

<style scoped>
.settings-page {
  animation: fadeIn 0.5s;
}

.setting-item {
  margin-bottom: 25px;
}

.setting-item label {
  display: block;
  margin-bottom: 10px;
  font-weight: 500;
  color: var(--text-primary);
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.info-item:last-child {
  border-bottom: none;
}

.info-item span {
  color: var(--text-primary);
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
