<template>
  <el-container class="app-container">
    
    <el-header height="40px" class="app-header glass-panel" v-if="!isWelcomePage">
      <div class="header-left">
        <div class="logo-text">NexusHub <span class="version">v3.0.1 (Phoenix)</span></div>
      </div>
      <div class="header-right">
        <div class="status-item" @click="showCommandPalette = true">
          <el-tooltip content="Command Palette (Ctrl+K)" placement="bottom">
            <el-icon><Search /></el-icon>
          </el-tooltip>
        </div>
        <div class="status-item" :class="{ 'online': isOnline, 'offline': !isOnline }">
          <div class="status-dot"></div>
          <span class="status-text">{{ latency }}ms</span>
        </div>
        <div class="status-item time">{{ currentTimeStr }}</div>
      </div>
    </el-header>

    <el-main class="app-main">
       <LoadingOverlay :loading="globalLoading" />
       
       <!-- Content Area with Transition -->
       <router-view v-slot="{ Component }">
         <transition name="fade" mode="out-in">
           <div class="view-wrapper">
             <component :is="Component" />
           </div>
         </transition>
       </router-view>
    </el-main>

    <!-- Dock is fixed at bottom -->
    <TheDock 
      v-if="!isWelcomePage"
      @open-palette="showCommandPalette = true"
    />

    <CommandPalette v-model="showCommandPalette" />
  </el-container>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import LoadingOverlay from './components/LoadingOverlay.vue'
import TheDock from './components/TheDock.vue'
import CommandPalette from './components/CommandPalette.vue'
import api from './api'
import { ElMessage } from 'element-plus'
import { useThemeStore } from './store/theme'

const route = useRoute()
const isWelcomePage = computed(() => route.path === '/')
const showCommandPalette = ref(false)
const globalLoading = ref(false)
const currentTimeStr = ref('')
const latency = ref(0)
const isOnline = ref(true)
const themeStore = useThemeStore()

// Time
const updateTime = () => {
  const now = new Date()
  currentTimeStr.value = now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

// Health Check (Optimized: 30s interval)
const checkHealth = async () => {
  const start = Date.now()
  try {
    await api.get('/health')
    latency.value = Date.now() - start
    isOnline.value = true
  } catch (e) {
    latency.value = 0
    isOnline.value = false
  }
}

let timer, healthTimer

onMounted(() => {
  updateTime()
  checkHealth()
  timer = setInterval(updateTime, 1000)
  healthTimer = setInterval(checkHealth, 30000)
  
  // 加载主题设置
  themeStore.loadTheme()
  
  window.addEventListener('keydown', handleKeydown)
  window.addEventListener('api:loading', handleApiLoading)
})

onUnmounted(() => {
  clearInterval(timer)
  clearInterval(healthTimer)
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('api:loading', handleApiLoading)
})

const handleKeydown = (e) => {
  if (e.ctrlKey && e.key === 'k') {
    e.preventDefault()
    showCommandPalette.value = !showCommandPalette.value
  }
}

const handleApiLoading = (e) => {
  globalLoading.value = !!e.detail
}


</script>

<style scoped>
.app-container {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  color: var(--text-primary);
  transition: color 0.3s ease;
}



/* Header */
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  z-index: 100;
  border-bottom: 1px solid var(--border-color);
  transition: border-color 0.3s ease;
}

.glass-panel {
  background: var(--card-bg);
  backdrop-filter: blur(12px);
  transition: background 0.3s ease;
}

.logo-text {
  font-weight: 700;
  font-size: 14px;
  letter-spacing: 0.5px;
  color: var(--text-primary);
  transition: color 0.3s ease;
}
.version {
  font-weight: 400;
  opacity: 0.9;
  font-size: 11px;
  margin-left: 6px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  cursor: pointer;
  opacity: 0.8;
  transition: opacity 0.2s;
  color: var(--text-secondary);
}
.status-item:hover { opacity: 1; }

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--text-secondary);
  transition: background 0.3s ease;
}
.online .status-dot { background: #67c23a; box-shadow: 0 0 4px #67c23a; }
.offline .status-dot { background: #f56c6c; }

/* Main Content */
.app-main {
  position: relative;
  padding: 0; /* Let views handle padding */
  overflow-y: auto;
  /* Fix scrollbar */
  scrollbar-width: thin;
  scrollbar-color: var(--border-color) transparent;
}

.view-wrapper {
  min-height: 100%;
  padding: 20px;
  /* Make space for Dock */
  padding-bottom: 100px; 
}



/* Transitions */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
