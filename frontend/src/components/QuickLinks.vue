<template>
  <div class="quick-links-widget glass-card">
    <div class="widget-header">
      <h3>Quick Links</h3>
      <div class="header-actions">
        <el-button link :icon="Refresh" size="small" @click="checkAllStatus" :loading="checking">刷新状态</el-button>
        <el-button link :icon="Plus" size="small" @click="showAddDialog = true"></el-button>
      </div>
    </div>
    <div class="links-grid">
      <a 
        v-for="(link, index) in links" 
        :key="index" 
        :href="link.url" 
        target="_blank" 
        class="link-item"
      >
        <div class="status-dot" :class="link.status || 'unknown'"></div>
        <div class="link-icon">
          <img :src="getIconUrl(link)" @error="handleIconError($event, link)" />
        </div>
        <span class="link-title">{{ link.title }}</span>
        <div class="latency-tag" v-if="link.latency">{{ link.latency }}ms</div>
        <div class="delete-btn" @click.prevent="removeLink(index)">
          <el-icon><Close /></el-icon>
        </div>
      </a>
    </div>

    <!-- 添加链接对话框 -->
    <el-dialog v-model="showAddDialog" title="添加快捷方式" width="300px" append-to-body>
      <el-form :model="form" label-position="top">
        <el-form-item label="名称">
          <el-input v-model="form.title" placeholder="GitHub" />
        </el-form-item>
        <el-form-item label="URL">
          <el-input v-model="form.url" placeholder="https://github.com" />
        </el-form-item>
        <el-form-item label="图标 URL (可选)">
          <el-input v-model="form.icon" placeholder="https://..." />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="addLink">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Plus, Close, Refresh } from '@element-plus/icons-vue'
import api from '../api'

const showAddDialog = ref(false)
const checking = ref(false)
const form = ref({ title: '', url: '', icon: '' })

const links = ref([
  { title: 'GitHub', url: 'https://github.com', icon: '' },
  { title: 'ChatGPT', url: 'https://chat.openai.com', icon: '' },
  { title: 'Bilibili', url: 'https://www.bilibili.com', icon: '' },
  { title: 'Localhost', url: 'http://localhost:8080', icon: '' },
])

onMounted(() => {
  const saved = localStorage.getItem('quick_links')
  if (saved) {
    links.value = JSON.parse(saved)
  }
  checkAllStatus()
})

const getIconUrl = (link) => {
  if (link.icon) return link.icon
  try {
    const domain = new URL(link.url).hostname
    return `https://www.google.com/s2/favicons?domain=${domain}&sz=64`
  } catch (e) {
    return ''
  }
}

const handleIconError = (e, link) => {
  // 如果 Google 服务也失败了，显示默认占位符（这里使用透明像素或默认图，或者利用 CSS 显示首字母）
  e.target.style.display = 'none'
  if (!e.target.parentElement.querySelector('.fallback-text')) {
    const span = document.createElement('span')
    span.className = 'fallback-text'
    span.innerText = link.title[0] || '?'
    e.target.parentElement.appendChild(span)
  }
}

const save = () => {
  localStorage.setItem('quick_links', JSON.stringify(links.value))
}

const addLink = () => {
  if (!form.value.title || !form.value.url) return
  links.value.push({ ...form.value, status: 'unknown' })
  save()
  form.value = { title: '', url: '', icon: '' }
  showAddDialog.value = false
  checkAllStatus()
}

const removeLink = (index) => {
  links.value.splice(index, 1)
  save()
}

const checkAllStatus = async () => {
  checking.value = true
  for (const link of links.value) {
    try {
      const res = await api.post('/monitor/check', { url: link.url })
      link.status = res.status
      link.latency = res.latency
    } catch (e) {
      link.status = 'offline'
      link.latency = 0
    }
  }
  save()
  checking.value = false
}
</script>

<style scoped>
.quick-links-widget {
  padding: 20px;
  height: 100%;
}

.widget-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.widget-header h3 {
  margin: 0;
  font-size: 16px;
  color: var(--text-primary);
}

.links-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 15px;
}

.link-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  color: var(--text-primary);
  position: relative;
  transition: transform 0.2s;
  padding: 10px;
  border-radius: 8px;
}

.link-item:hover {
  background: rgba(255,255,255,0.1);
  transform: translateY(-2px);
}

.link-icon {
  width: 48px;
  height: 48px;
  background: rgba(255,255,255,0.1);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
  overflow: hidden;
}

.link-icon img {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.fallback-icon, .fallback-text {
  font-size: 20px;
  font-weight: bold;
  color: #fff;
}

.link-title {
  font-size: 12px;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
}

.delete-btn {
  position: absolute;
  top: 0;
  right: 0;
  width: 20px;
  height: 20px;
  background: #f56c6c;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.2s;
}

.link-item:hover .delete-btn {
  opacity: 1;
}

/* 状态指示灯 */
.status-dot {
  position: absolute;
  top: 8px;
  left: 8px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  border: 1px solid rgba(0,0,0,0.2);
}

.status-dot.online { background: #67c23a; box-shadow: 0 0 5px #67c23a; }
.status-dot.offline { background: #f56c6c; }
.status-dot.unknown { background: #909399; }

.latency-tag {
  position: absolute;
  bottom: 2px;
  right: 5px;
  font-size: 9px;
  color: rgba(255,255,255,0.6);
}
</style>