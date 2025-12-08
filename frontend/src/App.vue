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
          <p>个人工作站</p>
        </div>

        <el-menu
          :default-active="$route.path"
          router
          class="sidebar-menu"
        >
          <el-menu-item index="/dashboard">
            <el-icon><el-icon-house /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>

          <el-menu-item index="/files">
            <el-icon><el-icon-folder /></el-icon>
            <span>文件管理</span>
          </el-menu-item>

          <el-menu-item index="/code">
            <el-icon><el-icon-document-copy /></el-icon>
            <span>代码编辑</span>
          </el-menu-item>

          <el-menu-item index="/notes">
            <el-icon><el-icon-document /></el-icon>
            <span>笔记</span>
          </el-menu-item>

          <el-menu-item index="/data">
            <el-icon><el-icon-data-analysis /></el-icon>
            <span>数据分析</span>
          </el-menu-item>

          <el-menu-item index="/calculator">
            <el-icon><el-icon-promotion /></el-icon>
            <span>计算器</span>
          </el-menu-item>

          <el-menu-item index="/todos">
            <el-icon><el-icon-list /></el-icon>
            <span>TODO目标</span>
          </el-menu-item>

          <el-menu-item index="/chat">
            <el-icon><el-icon-chat-dot-round /></el-icon>
            <span>AI 聊天</span>
          </el-menu-item>

          <el-menu-item index="/settings">
            <el-icon><el-icon-setting /></el-icon>
            <span>设置</span>
          </el-menu-item>
        </el-menu>

        <!-- 背景音乐控制 -->
        <div class="music-controls">
          <MusicPlayer />
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
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useThemeStore } from './store/theme'
import MusicPlayer from './components/MusicPlayer.vue'

const route = useRoute()
const themeStore = useThemeStore()

// 判断是否为欢迎页面
const isWelcomePage = computed(() => {
  return route.path === '/'
})

const pageTitle = computed(() => {
  const titles = {
    '/dashboard': '仪表盘',
    '/files': '文件管理',
    '/code': '代码编辑器',
    '/notes': '笔记',
    '/data': '数据分析',
    '/calculator': '高级计算器',
    '/todos': 'TODO目标',
    '/chat': 'AI 聊天',
    '/settings': '设置'
  }
  return titles[route.path] || 'NexusHub'
})

const toggleTheme = () => {
  themeStore.toggleTheme()
}

onMounted(() => {
  themeStore.loadTheme()
})
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background-color: var(--bg-color);
}

.sidebar {
  background-color: var(--card-bg);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  height: 100vh;
  position: sticky;
  top: 0;
}

.logo {
  padding: 30px 20px;
  text-align: center;
  border-bottom: 1px solid var(--border-color);
}

.logo h2 {
  margin: 0;
  font-size: 24px;
  color: var(--text-primary);
  font-weight: 700;
}

.logo p {
  margin: 5px 0 0 0;
  font-size: 12px;
  color: var(--text-secondary);
}

.sidebar-menu {
  flex: 1;
  border: none;
  background-color: transparent;
}

.music-controls {
  padding: 15px;
  border-top: 1px solid var(--border-color);
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
