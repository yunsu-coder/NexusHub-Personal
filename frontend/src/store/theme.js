import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'
import config from '../config'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref({
    theme_name: 'dark',
    primary_color: '#000000',
    secondary_color: '#ffffff',
    theme_template: 'default',
  
  })

  const isDark = ref(true)

  async function loadTheme() {
    try {
      const data = await api.get('/theme')
      theme.value = data
      isDark.value = data.theme_name === 'dark'
      applyTheme()
    } catch (error) {
      console.error('Failed to load theme:', error)
    }
  }

  async function saveTheme(newTheme) {
    try {
      const data = await api.put('/theme', newTheme)
      theme.value = data
      isDark.value = data.theme_name === 'dark'
      applyTheme()
    } catch (error) {
      console.error('Failed to save theme:', error)
    }
  }

  function applyTheme() {
    const root = document.documentElement
    
    // 为主题切换添加平滑过渡效果
    root.classList.add('theme-transition')
    
    // 基础颜色设置
    root.style.setProperty('--primary-color', theme.value.primary_color)
    root.style.setProperty('--secondary-color', theme.value.secondary_color)

    // 深色/浅色主题切换
    if (isDark.value) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }

    // 应用主题模板的完整颜色方案
    const themeConfigs = {
      default: {
        primaryColor: isDark.value ? '#667eea' : '#409eff',
        secondaryColor: isDark.value ? '#f56565' : '#f56c6c',
        textPrimary: isDark.value ? '#ffffff' : '#303133',
        textSecondary: isDark.value ? '#a0a0a0' : '#606266',
        bgColor: isDark.value ? '#121212' : '#f5f7fa',
        cardBg: isDark.value ? '#1e1e1e' : '#ffffff',
        sidebarBg: isDark.value ? '#181818' : '#ffffff',
        borderColor: isDark.value ? '#333333' : '#e4e7ed',
        gradientPrimary: isDark.value ? 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' : 'linear-gradient(135deg, #409eff 0%, #67c23a 100%)',
        gradientCard: isDark.value ? 'linear-gradient(135deg, rgba(30,30,30,0.85) 0%, rgba(18,18,18,0.95) 100%)' : 'linear-gradient(135deg, rgba(255,255,255,0.85) 0%, rgba(245,247,250,0.95) 100%)'
      },
      neon: {
        primaryColor: '#00ffff',
        secondaryColor: '#ff00ff',
        textPrimary: '#ffffff',
        textSecondary: '#00ffff',
        bgColor: '#0a0a0a',
        cardBg: '#111111',
        sidebarBg: '#0d0d0d',
        borderColor: '#00ffff',
        gradientPrimary: 'linear-gradient(135deg, #00ffff 0%, #ff00ff 100%)',
        gradientCard: 'linear-gradient(135deg, rgba(17,17,17,0.95) 0%, rgba(10,10,10,0.98) 100%)'
      },
      forest: {
        primaryColor: '#22c55e',
        secondaryColor: '#16a34a',
        textPrimary: '#ffffff',
        textSecondary: '#86efac',
        bgColor: '#06391a',
        cardBg: '#0a4e23',
        sidebarBg: '#074520',
        borderColor: '#22c55e',
        gradientPrimary: 'linear-gradient(135deg, #22c55e 0%, #16a34a 100%)',
        gradientCard: 'linear-gradient(135deg, rgba(10,78,35,0.9) 0%, rgba(6,57,26,0.95) 100%)'
      },
      ocean: {
        primaryColor: '#0093e9',
        secondaryColor: '#80d0c7',
        textPrimary: '#ffffff',
        textSecondary: '#80d0c7',
        bgColor: '#030b22',
        cardBg: '#05103d',
        sidebarBg: '#040d33',
        borderColor: '#0093e9',
        gradientPrimary: 'linear-gradient(135deg, #0093e9 0%, #80d0c7 100%)',
        gradientCard: 'linear-gradient(135deg, rgba(5,16,61,0.9) 0%, rgba(3,11,34,0.95) 100%)'
      },
      sunset: {
        primaryColor: '#f97316',
        secondaryColor: '#ef4444',
        textPrimary: '#1f2937',
        textSecondary: '#6b7280',
        bgColor: '#ffedd5',
        cardBg: '#ffffff',
        sidebarBg: '#fef3c7',
        borderColor: '#f59e0b',
        gradientPrimary: 'linear-gradient(135deg, #f97316 0%, #ef4444 100%)',
        gradientCard: 'linear-gradient(135deg, rgba(255,255,255,0.95) 0%, rgba(254,243,199,0.9) 100%)'
      }
    }

    const config = themeConfigs[theme.value.theme_template] || themeConfigs.default
    
    // 应用所有CSS变量
    root.style.setProperty('--primary-color', config.primaryColor)
    root.style.setProperty('--secondary-color', config.secondaryColor)
    root.style.setProperty('--text-primary', config.textPrimary)
    root.style.setProperty('--text-secondary', config.textSecondary)
    root.style.setProperty('--bg-color', config.bgColor)
    root.style.setProperty('--card-bg', config.cardBg)
    root.style.setProperty('--sidebar-bg', config.sidebarBg)
    root.style.setProperty('--sidebar-hover', isDark.value ? 'rgba(102, 126, 234, 0.1)' : 'rgba(64, 158, 255, 0.1)')
    root.style.setProperty('--border-color', config.borderColor)
    root.style.setProperty('--gradient-primary', config.gradientPrimary)
    root.style.setProperty('--gradient-card', config.gradientCard)

    // 背景图（避免 http 资源导致移动端阻止加载）
    const bg = normalizeAssetUrl(theme.value.background_image)
    if (bg) {
      root.style.setProperty('--app-bg-image', `url('${bg}')`)
    }
    
    // 移除过渡效果类，以便下次切换时重新应用
    setTimeout(() => {
      root.classList.remove('theme-transition')
    }, 500)
  }

  function toggleTheme() {
    isDark.value = !isDark.value
    theme.value.theme_name = isDark.value ? 'dark' : 'light'
    saveTheme(theme.value)
  }

  return {
    theme,
    isDark,
    loadTheme,
    saveTheme,
    toggleTheme,
    applyTheme
  }
})
  function normalizeAssetUrl(url) {
    if (!url) return ''
    try {
      const u = new URL(url)
      return u.pathname // 转为相对路径，避免混合内容
    } catch {
      return url
    }
  }
