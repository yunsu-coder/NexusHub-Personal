<template>
  <div class="quick-links-widget">
    <div class="widget-header">
      <h3>快速链接</h3>
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
        :class="{ 'offline': link.status === 'offline' }"
      >
        <div class="link-card">
          <div class="status-dot" :class="link.status || 'unknown'"></div>
          <div class="link-icon">
            <img :src="getIconUrl(link)" @error="handleIconError($event, link)" />
            <span class="fallback-text">{{ link.title[0] || '?' }}</span>
          </div>
          <span class="link-title">{{ link.title }}</span>
          <div class="latency-tag" v-if="link.latency">{{ link.latency }}ms</div>
          <div class="delete-btn" @click.prevent="removeLink(index)">
            <el-icon><Close /></el-icon>
          </div>
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
  padding: 24px;
  background: var(--card-bg);
  border-radius: 12px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  height: 100%;
}

.quick-links-widget:hover {
  box-shadow: var(--shadow-md);
}

.widget-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.widget-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.links-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 16px;
}

.link-item {
  text-decoration: none;
  color: var(--text-primary);
  transition: all 0.3s ease;
  display: block;
}

.link-item.offline {
  opacity: 0.7;
}

.link-card {
  position: relative;
  background: var(--bg-light);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  height: 100%;
  min-height: 100px;
}

.link-item:hover .link-card {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
  border-color: var(--primary-color);
  background: var(--card-bg);
}

.link-icon {
  width: 56px;
  height: 56px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
  overflow: hidden;
  position: relative;
}

.link-icon img {
  width: 36px;
  height: 36px;
  object-fit: contain;
  transition: opacity 0.3s ease;
}

.link-icon img[src=""] {
  opacity: 0;
}

.fallback-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 24px;
  font-weight: 700;
  color: var(--primary-color);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.link-icon img[src=""] + .fallback-text,
.link-icon img.error + .fallback-text {
  opacity: 1;
}

.link-title {
  font-size: 13px;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  font-weight: 500;
  color: var(--text-primary);
  line-height: 1.3;
}

.delete-btn {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 22px;
  height: 22px;
  background: var(--danger-color);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  opacity: 0;
  transition: all 0.3s ease;
  cursor: pointer;
  z-index: 10;
}

.link-item:hover .delete-btn {
  opacity: 1;
  transform: scale(1);
}

.delete-btn:hover {
  background: #f56c6c;
  transform: scale(1.1);
}

/* 状态指示灯 */
.status-dot {
  position: absolute;
  top: 8px;
  left: 8px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: 2px solid var(--card-bg);
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
}

.status-dot.online { 
  background: var(--success-color); 
  box-shadow: 0 0 8px var(--success-color); 
}
.status-dot.offline { 
  background: var(--danger-color); 
}
.status-dot.unknown { 
  background: var(--text-secondary); 
}

.latency-tag {
  position: absolute;
  bottom: 8px;
  right: 8px;
  font-size: 10px;
  color: var(--text-secondary);
  background: var(--bg-light);
  padding: 2px 6px;
  border-radius: 10px;
  border: 1px solid var(--border-color);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .links-grid {
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
    gap: 12px;
  }
  
  .link-card {
    min-height: 80px;
    padding: 12px 8px;
  }
  
  .link-icon {
    width: 48px;
    height: 48px;
    margin-bottom: 8px;
  }
  
  .link-icon img {
    width: 32px;
    height: 32px;
  }
  
  .fallback-text {
    font-size: 20px;
  }
  
  .link-title {
    font-size: 12px;
  }
}
</style>