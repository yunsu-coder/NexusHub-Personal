<template>
  <el-container class="app-container" :class="{ 'zen-mode': isZenMode }">
    <!-- Static Background (Performance Optimized) -->
    <div class="bg-layer"></div>
    
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
      v-if="!isWelcomePage && !isZenMode"
      @toggle-zen="toggleZenMode"
      @open-palette="showCommandPalette = true"
    />

    <!-- Zen Mode Exit Button -->
    <div
      v-if="!isWelcomePage && isZenMode"
      class="zen-exit-button"
      @click="toggleZenMode"
    >
      退出专注模式
    </div>

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

const route = useRoute()
const isWelcomePage = computed(() => route.path === '/')
const isZenMode = ref(false)
const showCommandPalette = ref(false)
const globalLoading = ref(false)
const currentTimeStr = ref('')
const latency = ref(0)
const isOnline = ref(true)

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

const toggleZenMode = () => {
  isZenMode.value = !isZenMode.value
  ElMessage.info(isZenMode.value ? 'Zen Mode On' : 'Zen Mode Off')
}
</script>

<style scoped>
.app-container {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  color: #e0e0e0;
}

/* Background - Static High Quality Gradient/Image */
.bg-layer {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: -1;
  background: radial-gradient(circle at center, #2b32b2, #1488cc);
  background-image: url('https://images.unsplash.com/photo-1451187580459-43490279c0fa?q=80&w=2072&auto=format&fit=crop');
  background-size: cover;
  background-position: center;
}
.bg-layer::after {
  content: '';
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.5); /* Overlay */
  backdrop-filter: blur(2px);
}

/* Header */
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  z-index: 100;
  border-bottom: 1px solid rgba(255,255,255,0.05);
}

.glass-panel {
  background: rgba(20, 20, 20, 0.6);
  backdrop-filter: blur(12px);
}

.logo-text {
  font-weight: 700;
  font-size: 14px;
  letter-spacing: 0.5px;
  color: #fff;
}
.version {
  font-weight: 400;
  opacity: 0.6;
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
}
.status-item:hover { opacity: 1; }

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #666;
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
  scrollbar-color: rgba(255,255,255,0.2) transparent;
}

.view-wrapper {
  min-height: 100%;
  padding: 20px;
  /* Make space for Dock */
  padding-bottom: 100px; 
}

/* Zen Mode */
.zen-mode .app-header {
  transform: translateY(-100%);
  transition: transform 0.3s;
}
.zen-mode .view-wrapper {
  padding-bottom: 20px;
}

/* Zen Mode Exit Button */
.zen-exit-button {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 120;
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  font-size: 12px;
  cursor: pointer;
  border: 1px solid rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(12px);
  transition: background 0.2s, transform 0.1s;
}

.zen-exit-button:hover {
  background: rgba(0, 0, 0, 0.8);
  transform: translateY(-1px);
}

/* Transitions */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
