import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref({
    theme_name: 'dark',
    primary_color: '#000000',
    secondary_color: '#ffffff',
    background_image: '',
    background_opacity: 1.0,
    background_music: '',
    music_volume: 0.5
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
    root.style.setProperty('--primary-color', theme.value.primary_color)
    root.style.setProperty('--secondary-color', theme.value.secondary_color)

    if (isDark.value) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }

    if (theme.value.background_image) {
      root.style.setProperty('--bg-image', `url(${theme.value.background_image})`)
      const opacity = theme.value.background_opacity !== undefined ? theme.value.background_opacity : 1.0
      root.style.setProperty('--bg-opacity', opacity)
    }
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
