<template>
  <div class="dock-container" :class="{ 'hidden': hidden }">
    <div class="dock">
      <div 
        v-for="item in menuItems" 
        :key="item.path" 
        class="dock-item"
        :class="{ active: currentPath === item.path }"
        @click="navigate(item.path)"
      >
        <el-tooltip :content="item.label" placement="top" :offset="20">
          <div class="icon-wrapper">
            <el-icon :size="24"><component :is="item.icon" /></el-icon>
          </div>
        </el-tooltip>
        <div class="dot" v-if="currentPath === item.path"></div>
      </div>

      <div class="divider"></div>

      <!-- 功能按钮 -->

      
      <div class="dock-item" @click="$emit('open-palette')">
        <el-tooltip content="命令面板 (Ctrl+K)" placement="top" :offset="20">
          <div class="icon-wrapper search-btn">
            <el-icon :size="24"><Search /></el-icon>
          </div>
        </el-tooltip>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { 
  Odometer, Folder, Notebook, List, Collection, 
  ChatLineRound, Monitor, Reading, Search,
  Headset
} from '@element-plus/icons-vue'

const props = defineProps({
  hidden: Boolean
})

const router = useRouter()
const route = useRoute()

const currentPath = computed(() => route.path)

const menuItems = [
  { label: '仪表盘', path: '/dashboard', icon: Odometer },
  { label: '文件', path: '/files', icon: Folder },
  { label: '笔记', path: '/notes', icon: Notebook },
  { label: '目标', path: '/todos', icon: List },
  { label: '收藏', path: '/collection', icon: Collection },

  { label: '代码', path: '/code', icon: Monitor },
  { label: '博客', path: '/blog', icon: Reading },
]

const navigate = (path) => {
  router.push(path)
}
</script>

<style scoped>
.dock-container {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 2000;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.dock-container.hidden {
  transform: translate(-50%, 150%);
}

.dock {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: var(--card-bg);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-lg);
  gap: 12px;
  transition: all 0.3s ease;
}

.dock-item {
  position: relative;
  cursor: pointer;
  transition: all 0.2s ease;
}

.icon-wrapper {
  width: 50px;
  height: 50px;
  border-radius: 14px;
  background: var(--hover-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-primary);
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  border: 1px solid var(--border-color);
}

.dock-item:hover .icon-wrapper {
  transform: translateY(-10px) scale(1.2);
  background: var(--card-bg);
  box-shadow: var(--shadow-md);
  margin: 0 10px; /* 挤开周围图标 */
}

/* 激活状态 */
.dock-item.active .icon-wrapper {
  background: var(--card-bg);
  color: var(--primary-color);
  box-shadow: 0 0 15px rgba(64, 158, 255, 0.3);
  border-color: var(--primary-color);
}

.zen-btn {
  background: rgba(230, 162, 60, 0.2);
}

.search-btn {
  background: rgba(103, 194, 58, 0.2);
}

.divider {
  width: 1px;
  height: 30px;
  background: var(--border-color);
  margin: 0 5px;
  transition: background 0.3s ease;
}

.dot {
  position: absolute;
  bottom: -6px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--primary-color);
  opacity: 0.8;
  transition: background 0.3s ease;
}
</style>