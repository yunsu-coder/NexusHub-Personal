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
            <label>主题模板</label>
            <el-radio-group v-model="theme.theme_template" @change="applyThemeTemplate">
              <el-radio-button label="default">默认主题</el-radio-button>
              <el-radio-button label="neon">霓虹主题</el-radio-button>
              <el-radio-button label="forest">森林主题</el-radio-button>
              <el-radio-button label="ocean">海洋主题</el-radio-button>
              <el-radio-button label="sunset">日落主题</el-radio-button>
            </el-radio-group>
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



          <div class="setting-item" style="margin-top: 30px">
            <el-alert title="提示" type="info" :closable="false">
              <p>• 主题更改会立即生效</p>
  
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
            <el-tag>v1.0.3</el-tag>
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
  theme_template: 'default',

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

// 应用主题模板
const applyThemeTemplate = () => {
  const templates = {
    default: {
      primary_color: '#000000',
      secondary_color: '#ffffff'
    },
    neon: {
      primary_color: '#00ff00',
      secondary_color: '#ff00ff'
    },
    forest: {
      primary_color: '#006400',
      secondary_color: '#228B22'
    },
    ocean: {
      primary_color: '#000080',
      secondary_color: '#00BFFF'
    },
    sunset: {
      primary_color: '#FF6347',
      secondary_color: '#FFD700'
    }
  }

  const selectedTemplate = templates[theme.value.theme_template]
  if (selectedTemplate) {
    theme.value.primary_color = selectedTemplate.primary_color
    theme.value.secondary_color = selectedTemplate.secondary_color
    saveTheme()
    ElMessage.success('主题模板已应用')
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
