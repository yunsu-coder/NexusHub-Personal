<template>
  <div id="app" class="app-container">
    <!-- 欢迎页面 - 全屏显示,不带布局 -->
    <router-view v-if="isWelcomePage" />

    <!-- 其他页面 - 带侧边栏布局 -->
    <el-container v-else>
      <!-- 侧边栏 -->
      <el-aside width="250px" class="sidebar">
        <div class="logo">
          <h2>NexusHub</h2>
          <div class="logo-subtitle">个人工作站</div>
        </div>

        <div class="simple-menu">
          <div class="menu-item" :class="{ active: route.path === '/dashboard' }" @click="router.push('/dashboard')">
            <span>仪表盘</span>
          </div>
          <div class="menu-item" :class="{ active: route.path === '/files' }" @click="router.push('/files')">
            <span>文件管理</span>
          </div>
          <div class="menu-item" :class="{ active: route.path === '/notes' }" @click="router.push('/notes')">
            <span>笔记</span>
          </div>
          <div class="menu-item" :class="{ active: route.path === '/data' }" @click="router.push('/data')">
            <span>数据分析</span>
          </div>
          <div class="menu-item" :class="{ active: route.path === '/todos' }" @click="router.push('/todos')">
            <span>TODO目标</span>
          </div>
          <div class="menu-item" :class="{ active: route.path === '/collection' }" @click="router.push('/collection')">
            <span>我的收藏</span>
          </div>
          <div class="menu-item" :class="{ active: route.path === '/chat' }" @click="router.push('/chat')">
            <span>AI 聊天</span>
          </div>
          <div class="menu-item" :class="{ active: route.path === '/settings' }" @click="router.push('/settings')">
            <span>设置</span>
          </div>
        </div>


      </el-aside>

      <!-- 主内容区 -->
      <el-container>
        <el-header class="app-header">
          <div class="header-left">
            <h3>{{ pageTitle }}</h3>
          </div>
          <div class="header-right">
            <el-button
              circle
              @click="toggleTheme"
            >
              {{ themeStore.isDark ? '☀️' : '🌙' }}
            </el-button>
          </div>
        </el-header>

      <el-main class="app-main">
        <LoadingOverlay :loading="globalLoading" />
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useThemeStore } from './store/theme'
import LoadingOverlay from './components/LoadingOverlay.vue'
import { ElMessage } from 'element-plus'


const route = useRoute()
const router = useRouter()
const themeStore = useThemeStore()
const globalLoading = ref(false)

  // 菜单点击处理，用于调试
    const handleMenuClick = (index, indexPath) => {
      console.log('Menu clicked:', index, indexPath);
    };

    onMounted(() => {
      // Load theme from local storage
      const savedTheme = localStorage.getItem('theme');
      if (savedTheme) {
        document.documentElement.setAttribute('data-theme', savedTheme);
      }
      
      // Load theme from store
      themeStore.loadTheme();
    });

  // 判断是否为欢迎页面
const isWelcomePage = computed(() => {
  return route.path === '/'
})

const pageTitle = computed(() => {
    const titles = {
      '/dashboard': '仪表盘',
      '/files': '文件管理',
      '/notes': '笔记',
      '/data': '数据分析',
  
      '/todos': 'TODO目标',
      '/collection': '我的收藏',
      '/chat': 'AI 聊天',
      '/settings': '设置'
    }
    return titles[route.path] || 'NexusHub'
  })

const toggleTheme = () => {
  themeStore.toggleTheme()
  ElMessage.success('主题已切换')
}

// 监听全局 API 加载事件
const handleApiLoading = (e) => {
  globalLoading.value = !!(e?.detail)
}

onMounted(() => {
  window.addEventListener('api:loading', handleApiLoading)
})

onUnmounted(() => {
  window.removeEventListener('api:loading', handleApiLoading)
})


</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background-color: var(--bg-color);
  animation: fadeIn 0.5s ease;
}

.sidebar {
  width: 240px;
  height: 100vh;
  background-color: var(--sidebar-bg);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  padding: 20px;
  transition: transform 0.3s ease;
}

.logo {
  text-align: center;
  padding: 20px 0;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--border-color);
}

.logo h2 {
  font-size: 24px;
  font-weight: bold;
  color: var(--primary-color);
  margin: 0 0 5px 0;
}

.logo-subtitle {
  font-size: 14px;
  color: var(--text-secondary);
}

.simple-menu {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.menu-item {
  height: 50px;
  line-height: 50px;
  border-radius: 8px;
  padding: 0 20px;
  transition: all 0.3s ease;
  cursor: pointer;
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
  position: relative;
}

.menu-item:hover {
  background-color: rgba(64, 158, 255, 0.1);
  color: var(--primary-color);
}

.menu-item.active {
  background-color: rgba(64, 158, 255, 0.1);
  color: var(--primary-color);
}



.app-header {
  background-color: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 30px;
  height: 60px;
}

.header-left h3 {
  margin: 0;
  font-size: 20px;
  color: var(--text-primary);
}

.header-right {
  display: flex;
  gap: 10px;
}

.app-main {
  padding: 30px;
  background-color: var(--bg-color);
  min-height: calc(100vh - 60px);
  transition: all 0.3s ease;
  position: relative;
  z-index: 1;
}

/* 平滑滚动 */
.app-main::-webkit-scrollbar {
  width: 8px;
}

.app-main::-webkit-scrollbar-track {
  background: var(--bg-color);
}

.app-main::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
}

.app-main::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
