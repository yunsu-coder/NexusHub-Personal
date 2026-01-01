import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'
import config from '../config'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref({
    theme_name: 'light',
    primary_color: '#409eff',
    secondary_color: '#f56c6c',
    theme_template: 'sunset',
  })

  const isDark = ref(false)

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
    
    // 深色/浅色主题切换
    if (isDark.value) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }

    // 应用主题模板的完整颜色方案 - 确保所有主题下文字对比度符合WCAG标准
    // 所有主题使用单一纯色背景，确保视觉一致性
    const themeConfigs = {
      default: {
        primaryColor: isDark.value ? '#667eea' : '#409eff',
        secondaryColor: isDark.value ? '#f56565' : '#f56c6c',
        textPrimary: isDark.value ? '#ffffff' : '#111827',
        textSecondary: isDark.value ? '#a1a1aa' : '#4b5563',
        bgColor: isDark.value ? '#18181b' : '#f3f4f6',
        cardBg: isDark.value ? '#27272a' : '#ffffff',
        sidebarBg: isDark.value ? '#27272a' : '#ffffff',
        borderColor: isDark.value ? '#3f3f46' : '#e5e7eb'
      },
      dark: {
        primaryColor: '#667eea',
        secondaryColor: '#f56565',
        textPrimary: '#ffffff',
        textSecondary: '#a1a1aa',
        bgColor: '#18181b',
        cardBg: '#27272a',
        sidebarBg: '#27272a',
        borderColor: '#3f3f46'
      },
      light: {
        primaryColor: '#409eff',
        secondaryColor: '#f56c6c',
        textPrimary: '#111827',
        textSecondary: '#4b5563',
        bgColor: '#f3f4f6',
        cardBg: '#ffffff',
        sidebarBg: '#ffffff',
        borderColor: '#e5e7eb'
      },
      highContrast: {
        primaryColor: '#0066cc',
        secondaryColor: '#cc0000',
        textPrimary: '#000000',
        textSecondary: '#333333',
        bgColor: '#ffffff',
        cardBg: '#f0f0f0',
        sidebarBg: '#f0f0f0',
        borderColor: '#000000'
      },
      neon: {
        primaryColor: '#00ffff',
        secondaryColor: '#ff00ff',
        textPrimary: '#ffffff',
        textSecondary: '#e5e7eb',
        bgColor: '#000000',
        cardBg: '#171717',
        sidebarBg: '#171717',
        borderColor: '#00ffff'
      },
      forest: {
        primaryColor: '#22c55e',
        secondaryColor: '#16a34a',
        textPrimary: '#ffffff',
        textSecondary: '#dcfce7',
        bgColor: '#062e17',
        cardBg: '#065f46',
        sidebarBg: '#047857',
        borderColor: '#16a34a'
      },
      ocean: {
        primaryColor: '#0093e9',
        secondaryColor: '#80d0c7',
        textPrimary: '#ffffff',
        textSecondary: '#bae6fd',
        bgColor: '#0c102a',
        cardBg: '#0f172a',
        sidebarBg: '#1e293b',
        borderColor: '#0093e9'
      },
      sunset: {
        primaryColor: '#f97316',
        secondaryColor: '#ef4444',
        textPrimary: '#1f2937',
        textSecondary: '#6b7280',
        bgColor: '#fef3c7',
        cardBg: '#ffffff',
        sidebarBg: '#fef3c7',
        borderColor: '#f59e0b'
      }
    }

    // 确定当前主题配置
    let config
    if (themeConfigs[theme.value.theme_name]) {
      config = themeConfigs[theme.value.theme_name]
    } else {
      config = themeConfigs[theme.value.theme_template] || themeConfigs.default
    }
    
    // 应用所有CSS变量到文档根元素
    root.style.setProperty('--primary-color', config.primaryColor)
    root.style.setProperty('--secondary-color', config.secondaryColor)
    root.style.setProperty('--text-primary', config.textPrimary)
    root.style.setProperty('--text-secondary', config.textSecondary)
    root.style.setProperty('--bg-color', config.bgColor)
    root.style.setProperty('--card-bg', config.cardBg)
    root.style.setProperty('--sidebar-bg', config.sidebarBg)
    root.style.setProperty('--sidebar-text', config.textPrimary)
    root.style.setProperty('--sidebar-hover', isDark.value ? 'rgba(102, 126, 234, 0.1)' : 'rgba(64, 158, 255, 0.1)')
    root.style.setProperty('--border-color', config.borderColor)
    
    // Element Plus 组件主题变量
    root.style.setProperty('--el-color-primary', config.primaryColor)
    root.style.setProperty('--el-text-color-primary', config.textPrimary)
    root.style.setProperty('--el-text-color-regular', config.textPrimary)
    root.style.setProperty('--el-text-color-secondary', config.textSecondary)
    root.style.setProperty('--el-bg-color', config.bgColor)
    root.style.setProperty('--el-bg-color-page', config.bgColor)
    root.style.setProperty('--el-bg-color-overlay', config.cardBg)
    root.style.setProperty('--el-border-color', config.borderColor)
    root.style.setProperty('--el-card-bg-color', config.cardBg)
    root.style.setProperty('--el-card-border-color', config.borderColor)
    root.style.setProperty('--el-table-bg-color', config.cardBg)
    root.style.setProperty('--el-table-text-color', config.textPrimary)
    root.style.setProperty('--el-table-header-bg-color', isDark.value ? '#27272a' : '#f5f7fa')
    root.style.setProperty('--el-table-header-text-color', config.textPrimary)
    root.style.setProperty('--el-input-bg-color', isDark.value ? '#27272a' : '#ffffff')
    root.style.setProperty('--el-input-text-color', config.textPrimary)
    root.style.setProperty('--el-input-placeholder-color', config.textSecondary)
    root.style.setProperty('--el-input-border-color', config.borderColor)
    root.style.setProperty('--el-dialog-bg-color', config.cardBg)
    root.style.setProperty('--el-drawer-bg-color', config.cardBg)
    
    // 移除固定背景图，确保单一纯色背景
    root.style.setProperty('--app-bg-image', 'none')
    
    // 移除过渡效果类，以便下次切换时重新应用
    setTimeout(() => {
      root.classList.remove('theme-transition')
    }, 300) // 与过渡时间匹配
  }

  function toggleTheme() {
    isDark.value = !isDark.value
    theme.value.theme_name = isDark.value ? 'dark' : 'light'
    saveTheme(theme.value)
  }

  // API接口：获取当前背景色
  function getCurrentBgColor() {
    return document.documentElement.style.getPropertyValue('--bg-color') || getComputedStyle(document.documentElement).getPropertyValue('--bg-color')
  }

  // API接口：获取当前主题配置
  function getCurrentThemeConfig() {
    const root = document.documentElement
    return {
      primaryColor: getComputedStyle(root).getPropertyValue('--primary-color').trim(),
      secondaryColor: getComputedStyle(root).getPropertyValue('--secondary-color').trim(),
      textPrimary: getComputedStyle(root).getPropertyValue('--text-primary').trim(),
      textSecondary: getComputedStyle(root).getPropertyValue('--text-secondary').trim(),
      bgColor: getComputedStyle(root).getPropertyValue('--bg-color').trim(),
      cardBg: getComputedStyle(root).getPropertyValue('--card-bg').trim(),
      sidebarBg: getComputedStyle(root).getPropertyValue('--sidebar-bg').trim(),
      borderColor: getComputedStyle(root).getPropertyValue('--border-color').trim()
    }
  }

  // API接口：切换到指定主题
  function switchToTheme(themeName) {
    if (themeName) {
      theme.value.theme_name = themeName
      isDark.value = themeName === 'dark'
      applyTheme()
      saveTheme(theme.value)
    }
  }

  // API接口：更新背景色
  function updateBgColor(color) {
    if (color) {
      const root = document.documentElement
      root.classList.add('theme-transition')
      root.style.setProperty('--bg-color', color)
      root.style.setProperty('--el-bg-color', color)
      root.style.setProperty('--el-bg-color-page', color)
      
      // 移除过渡效果类
      setTimeout(() => {
        root.classList.remove('theme-transition')
      }, 300)
    }
  }

  return {
    theme,
    isDark,
    loadTheme,
    saveTheme,
    toggleTheme,
    applyTheme,
    // API接口方法
    getCurrentBgColor,
    getCurrentThemeConfig,
    switchToTheme,
    updateBgColor
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
