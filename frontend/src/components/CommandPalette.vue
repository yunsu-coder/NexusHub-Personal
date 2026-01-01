<template>
  <el-dialog
    v-model="visible"
    :show-close="false"
    :modal="true"
    :lock-scroll="true"
    class="command-palette-dialog"
    width="600px"
    align-center
    destroy-on-close
  >
    <div class="palette-container">
      <div class="search-input-wrapper">
        <el-icon class="search-icon"><Search /></el-icon>
        <input 
          ref="inputRef"
          v-model="query"
          class="palette-input"
          placeholder="搜索页面、命令或文件..."
          @keydown.down.prevent="selectNext"
          @keydown.up.prevent="selectPrev"
          @keydown.enter="execute"
          @keydown.esc="close"
        />
        <span class="esc-hint">ESC</span>
      </div>

      <div class="results-list" v-if="filteredCommands.length > 0">
        <div 
          v-for="(cmd, index) in filteredCommands" 
          :key="index"
          class="result-item"
          :class="{ selected: selectedIndex === index }"
          @click="executeCommand(cmd)"
          @mousemove="selectedIndex = index"
        >
          <div class="cmd-icon">
            <el-icon><component :is="cmd.icon" /></el-icon>
          </div>
          <div class="cmd-info">
            <div class="cmd-title">{{ cmd.title }}</div>
            <div class="cmd-desc" v-if="cmd.desc">{{ cmd.desc }}</div>
          </div>
          <div class="cmd-shortcut" v-if="cmd.shortcut">{{ cmd.shortcut }}</div>
        </div>
      </div>
      <div class="empty-state" v-else>
        无相关结果
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useThemeStore } from '../store/theme'
import { 
  Search, Odometer, Folder, Notebook, List, Collection, 
  ChatLineRound, Monitor, Reading, Moon, Sunny, VideoPlay,
  FullScreen
} from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: Boolean
})

const emit = defineEmits(['update:modelValue'])

const router = useRouter()
const themeStore = useThemeStore()
const inputRef = ref(null)
const query = ref('')
const selectedIndex = ref(0)

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const commands = [
  // 导航
  { title: '前往 仪表盘', desc: 'Go to Dashboard', icon: Odometer, action: () => router.push('/dashboard') },
  { title: '前往 文件管理', desc: 'Go to Files', icon: Folder, action: () => router.push('/files') },
  { title: '前往 笔记', desc: 'Go to Notes', icon: Notebook, action: () => router.push('/notes') },
  { title: '前往 TODO', desc: 'Go to Todos', icon: List, action: () => router.push('/todos') },
  { title: '前往 收藏', desc: 'Go to Collection', icon: Collection, action: () => router.push('/collection') },
  
  { title: '前往 代码竞技场', desc: 'Go to Code Arena', icon: Monitor, action: () => router.push('/code') },
  { title: '前往 博客', desc: 'Go to Blog', icon: Reading, action: () => router.push('/blog') },
  
  // 系统动作
  { title: '切换主题', desc: 'Toggle Dark/Light Mode', icon: Moon, action: () => themeStore.toggleTheme() },
]

const filteredCommands = computed(() => {
  if (!query.value) return commands
  const lowerQuery = query.value.toLowerCase()
  return commands.filter(cmd => 
    cmd.title.toLowerCase().includes(lowerQuery) || 
    (cmd.desc && cmd.desc.toLowerCase().includes(lowerQuery))
  )
})

const selectNext = () => {
  selectedIndex.value = (selectedIndex.value + 1) % filteredCommands.value.length
}

const selectPrev = () => {
  selectedIndex.value = (selectedIndex.value - 1 + filteredCommands.value.length) % filteredCommands.value.length
}

const execute = () => {
  if (filteredCommands.value[selectedIndex.value]) {
    executeCommand(filteredCommands.value[selectedIndex.value])
  }
}

const executeCommand = (cmd) => {
  cmd.action()
  close()
}

const close = () => {
  visible.value = false
  query.value = ''
  selectedIndex.value = 0
}

watch(visible, (val) => {
  if (val) {
    nextTick(() => {
      inputRef.value?.focus()
    })
  }
})
</script>

<style scoped>
.palette-container {
  background: #1e1e1e;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 20px 50px rgba(0,0,0,0.5);
}

.search-input-wrapper {
  display: flex;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #333;
  gap: 12px;
}

.search-icon {
  font-size: 20px;
  color: #888;
}

.palette-input {
  flex: 1;
  background: transparent;
  border: none;
  font-size: 18px;
  color: #fff;
  outline: none;
}

.esc-hint {
  font-size: 12px;
  color: #666;
  border: 1px solid #444;
  padding: 2px 6px;
  border-radius: 4px;
}

.results-list {
  max-height: 400px;
  overflow-y: auto;
  padding: 8px;
}

.result-item {
  display: flex;
  align-items: center;
  padding: 12px;
  gap: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.1s;
  color: #ccc;
}

.result-item.selected {
  background: #2d2d2d;
  color: #fff;
}

.cmd-icon {
  display: flex;
  align-items: center;
  font-size: 18px;
}

.cmd-info {
  flex: 1;
}

.cmd-title {
  font-size: 14px;
  font-weight: 500;
}

.cmd-desc {
  font-size: 12px;
  color: #888;
  margin-top: 2px;
}

.empty-state {
  padding: 30px;
  text-align: center;
  color: #666;
}
</style>

<style>
/* 覆盖 Element Plus Dialog 默认样式以实现更纯粹的面板 */
.command-palette-dialog .el-dialog {
  background: transparent;
  box-shadow: none;
  border-radius: 12px;
}
.command-palette-dialog .el-dialog__header,
.command-palette-dialog .el-dialog__body {
  padding: 0;
  margin: 0;
}
</style>