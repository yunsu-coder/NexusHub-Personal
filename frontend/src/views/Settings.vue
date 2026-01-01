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
              <el-radio-button label="highContrast">高对比度</el-radio-button>
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
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>
            <span>ℹ️ 系统信息</span>
          </template>

          <div class="info-item">
            <span>版本</span>
            <el-tag>v3.0.1</el-tag>
          </div>
          
          <div class="info-item">
            <span>技术栈</span>
            <el-tag type="success">Vue 3 + Go + Gin</el-tag>
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
      primary_color: '#409eff',
      secondary_color: '#f56c6c'
    },
    neon: {
      primary_color: '#00ffff',
      secondary_color: '#ff00ff'
    },
    forest: {
      primary_color: '#22c55e',
      secondary_color: '#16a34a'
    },
    ocean: {
      primary_color: '#0093e9',
      secondary_color: '#80d0c7'
    },
    sunset: {
      primary_color: '#f97316',
      secondary_color: '#ef4444'
    }
  }

  const selectedTemplate = templates[theme.value.theme_template]
  if (selectedTemplate) {
    theme.value.primary_color = selectedTemplate.primary_color
    theme.value.secondary_color = selectedTemplate.secondary_color
    
    // 立即应用到 Store 以便获得即时视觉反馈
    themeStore.theme.theme_template = theme.value.theme_template
    themeStore.theme.primary_color = theme.value.primary_color
    themeStore.theme.secondary_color = theme.value.secondary_color
    themeStore.applyTheme()
    
    // 保存到后端
    saveTheme()
    ElMessage.success('主题模板已应用')
  }
}

const loadTheme = async () => {
  try {
    theme.value = await api.get('/theme')
    // 更新到主题仓库并应用
    themeStore.theme = theme.value
    themeStore.isDark = theme.value.theme_name === 'dark'
    themeStore.applyTheme()
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