<template>
  <div class="file-manager" @dragover.prevent="handleDragOver" @drop.prevent="handleDrop" @dragleave="handleDragLeave">
    <!-- 拖拽上传覆盖层 -->
    <div v-if="isDragging" class="drag-overlay">
      <el-icon :size="60"><UploadFilled /></el-icon>
      <div class="drag-text">释放以上传文件</div>
    </div>

    <!-- 顶部工具栏 -->
    <div class="file-toolbar glass-panel">
      <div class="toolbar-left">
        <el-upload
          class="upload-btn"
          :action="uploadUrl"
          :headers="uploadHeaders"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          :show-file-list="false"
          :before-upload="beforeUpload"
          multiple
        >
          <el-button type="primary" :icon="UploadFilled">上传</el-button>
        </el-upload>
        <el-button :icon="FolderAdd" @click="createFolder">新建文件夹</el-button>
      </div>
      <div class="toolbar-right">
        <el-select v-model="typeFilter" size="small" class="type-select" placeholder="类型">
          <el-option label="全部" value="all" />
          <el-option label="图片" value="image" />
          <el-option label="文档" value="document" />
          <el-option label="视频" value="video" />
          <el-option label="音频" value="audio" />
          <el-option label="其他" value="other" />
        </el-select>
        <el-input
          v-model="searchKeyword"
          placeholder="搜索文件..."
          prefix-icon="Search"
          class="search-input"
          clearable
        />
        <el-divider direction="vertical" />
        <el-tooltip content="网格视图">
          <el-button :type="viewMode === 'grid' ? 'primary' : ''" :icon="Menu" circle @click="viewMode = 'grid'" />
        </el-tooltip>
        <el-tooltip content="列表视图">
          <el-button :type="viewMode === 'list' ? 'primary' : ''" :icon="List" circle @click="viewMode = 'list'" />
        </el-tooltip>
      </div>
    </div>

    <!-- 面包屑导航 -->
    <div class="breadcrumb-bar glass-panel">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item><el-icon><HomeFilled /></el-icon></el-breadcrumb-item>
        <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index">{{ item }}</el-breadcrumb-item>
      </el-breadcrumb>
      <span class="file-count">{{ filteredFiles.length }} 个项目</span>
    </div>

    <!-- 内容区域 -->
    <div class="file-content" v-loading="loading" element-loading-background="rgba(0, 0, 0, 0.3)">
      <!-- 骨架屏加载 -->
      <div v-if="loading && files.length === 0" class="skeleton-wrapper">
        <el-skeleton :rows="3" animated />
      </div>

      <!-- 空状态 -->
      <el-empty v-if="!loading && filteredFiles.length === 0" description="暂无文件，拖拽或点击上传" />

      <!-- 网格视图 -->
      <div v-if="viewMode === 'grid' && filteredFiles.length > 0" class="file-grid">
        <div 
          v-for="file in filteredFiles" 
          :key="file.id" 
          class="file-card"
          :class="{ 'is-selected': selectedFile?.id === file.id }"
          @click="selectFile(file)"
          @dblclick="handleFileDblClick(file)"
          @contextmenu.prevent="showContextMenu($event, file)"
        >
          <div class="file-preview">
            <!-- 图片缩略图 -->
            <img v-if="isImage(file)" :src="getFileUrl(file)" loading="lazy" class="thumbnail-img" />
            <!-- 通用图标 -->
            <el-icon v-else :size="48" :color="getFileColor(file.extension)">
              <component :is="getFileIcon(file.extension)" />
            </el-icon>
          </div>
          <div class="file-name" :title="file.file_name">{{ file.file_name }}</div>
          <div class="file-meta">{{ formatFileSize(file.file_size) }}</div>
        </div>
      </div>

      <!-- 列表视图 -->
      <div v-if="viewMode === 'list' && filteredFiles.length > 0" class="file-list-wrapper">
        <el-table 
          :data="filteredFiles" 
          style="width: 100%; background: transparent;" 
          @row-click="selectFile"
          @row-dblclick="handleFileDblClick"
          @row-contextmenu="showContextMenu"
          highlight-current-row
        >
          <el-table-column label="名称" min-width="300">
            <template #default="{ row }">
              <div class="list-name">
                <el-icon :size="20" :color="getFileColor(row.extension)" style="margin-right: 8px">
                  <component :is="getFileIcon(row.extension)" />
                </el-icon>
                <span>{{ row.file_name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="file_size" label="大小" width="120">
            <template #default="{ row }">{{ formatFileSize(row.file_size) }}</template>
          </el-table-column>
          <el-table-column prop="updated_at" label="修改日期" width="180">
            <template #default="{ row }">{{ formatTime(row.updated_at || row.created_at) }}</template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- 文件详情抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      title="文件详情"
      direction="rtl"
      size="350px"
    >
      <div v-if="selectedFile" class="file-detail">
        <div class="detail-preview">
          <img v-if="isImage(selectedFile)" :src="getFileUrl(selectedFile)" class="detail-img" />
          <el-icon v-else :size="80" :color="getFileColor(selectedFile.extension)">
            <component :is="getFileIcon(selectedFile.extension)" />
          </el-icon>
        </div>
        <div class="detail-info">
          <h3>{{ selectedFile.file_name }}</h3>
          <p><strong>类型:</strong> {{ selectedFile.extension || '未知' }}</p>
          <p><strong>大小:</strong> {{ formatFileSize(selectedFile.file_size) }}</p>
          <p><strong>上传时间:</strong> {{ formatTime(selectedFile.created_at) }}</p>
          <p><strong>路径:</strong> {{ selectedFile.file_path }}</p>
        </div>
        <div class="detail-actions">
          <el-button type="primary" @click="downloadFile(selectedFile)" block>下载</el-button>
          <el-button type="danger" @click="deleteFile(selectedFile)" block>删除</el-button>
        </div>
      </div>
    </el-drawer>

    <!-- 右键菜单 -->
    <div 
      v-show="contextMenu.visible" 
      class="context-menu glass-panel" 
      :style="{ top: contextMenu.top + 'px', left: contextMenu.left + 'px' }"
    >
      <div class="menu-item" @click="handleFileDblClick(contextMenu.file)">
        <el-icon><View /></el-icon> 打开
      </div>
      <div class="menu-item" @click="downloadFile(contextMenu.file)">
        <el-icon><Download /></el-icon> 下载
      </div>
      <div class="menu-item" @click="renameFile(contextMenu.file)">
        <el-icon><Edit /></el-icon> 重命名
      </div>
      <div class="menu-divider"></div>
      <div class="menu-item danger" @click="deleteFile(contextMenu.file)">
        <el-icon><Delete /></el-icon> 删除
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  UploadFilled, FolderAdd, Search, Menu, List, 
  Document, Picture, VideoPlay, Files, Folder,
  HomeFilled, View, Download, Edit, Delete, Headset
} from '@element-plus/icons-vue'
import api from '../api'
import config from '../config'

// State
const loading = ref(false)
const viewMode = ref('grid')
const searchKeyword = ref('')
const typeFilter = ref('all')
const files = ref([])
const breadcrumbs = ref(['Home'])
const contextMenu = ref({ visible: false, top: 0, left: 0, file: null })
const isDragging = ref(false)
const selectedFile = ref(null)
const drawerVisible = ref(false)

// Computed
const uploadUrl = computed(() => `${config.api.baseURL}/files/upload`)
// 上传请求头，携带本地 token（如果存在）
const uploadHeaders = computed(() => {
  const token = localStorage.getItem('token')
  return token
    ? { Authorization: `Bearer ${token}` }
    : {}
}) 

const filteredFiles = computed(() => {
  let list = [...files.value]

  // 类型筛选
  if (typeFilter.value !== 'all') {
    list = list.filter(f => {
      const ext = (f.extension || '').toLowerCase()
      if (!ext) return typeFilter.value === 'other'
      const imgExt = ['.jpg', '.jpeg', '.png', '.gif', '.webp']
      const docExt = ['.pdf', '.doc', '.docx', '.ppt', '.pptx', '.xls', '.xlsx', '.txt']
      const videoExt = ['.mp4', '.mov', '.avi']
      const audioExt = ['.mp3', '.wav']
      if (typeFilter.value === 'image') return imgExt.includes(ext)
      if (typeFilter.value === 'document') return docExt.includes(ext)
      if (typeFilter.value === 'video') return videoExt.includes(ext)
      if (typeFilter.value === 'audio') return audioExt.includes(ext)
      return ![...imgExt, ...docExt, ...videoExt, ...audioExt].includes(ext)
    })
  }

  // 关键字搜索
  if (searchKeyword.value) {
    const q = searchKeyword.value.toLowerCase()
    list = list.filter(f => f.file_name.toLowerCase().includes(q))
  }

  return list
})

// Methods
const loadFiles = async () => {
  loading.value = true
  try {
    files.value = await api.get('/files')
  } catch (e) {
    ElMessage.error('加载文件失败')
  } finally {
    loading.value = false
  }
}

// Drag & Drop
const handleDragOver = () => isDragging.value = true
const handleDragLeave = (e) => {
  if (e.relatedTarget === null) isDragging.value = false
}
const handleDrop = async (e) => {
  isDragging.value = false
  const droppedFiles = e.dataTransfer.files
  if (droppedFiles.length === 0) return

  const formData = new FormData()
  // Currently only support single file upload via this simple handler for demo
  // Production should loop and upload all
  for (let i = 0; i < droppedFiles.length; i++) {
    formData.append('file', droppedFiles[i])
    try {
        await api.post('/files/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        })
        ElMessage.success(`已上传: ${droppedFiles[i].name}`)
    } catch (err) {
        ElMessage.error(`上传失败: ${droppedFiles[i].name}`)
    }
  }
  loadFiles()
}

const selectFile = (file) => {
  selectedFile.value = file
  drawerVisible.value = true
}

const handleFileDblClick = (file) => {
  if (isImage(file)) {
    // Preview image
    window.open(getFileUrl(file), '_blank')
  } else {
    downloadFile(file)
  }
}

const showContextMenu = (event, file) => {
  const e = event.event || event
  contextMenu.value = {
    visible: true,
    top: e.clientY,
    left: e.clientX,
    file: file
  }
}

const closeContextMenu = () => contextMenu.value.visible = false

const renameFile = async (file) => {
  closeContextMenu()
  try {
    const { value } = await ElMessageBox.prompt('请输入新文件名', '重命名', {
      inputValue: file.file_name,
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })
    if (value && value !== file.file_name) {
      await api.put(`/files/${file.id}/rename`, { new_name: value })
      ElMessage.success('重命名成功')
      loadFiles()
    }
  } catch (e) {}
}

const downloadFile = (file) => {
  closeContextMenu()
  window.open(`${config.api.baseURL}/files/download/${file.id}`, '_blank')
}

const deleteFile = async (file) => {
  closeContextMenu()
  try {
    await ElMessageBox.confirm(`确定删除 "${file.file_name}" 吗?`, '警告', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await api.delete(`/files/${file.id}`)
    ElMessage.success('已删除')
    loadFiles()
    if (selectedFile.value?.id === file.id) {
        selectedFile.value = null
        drawerVisible.value = false
    }
  } catch (e) {}
}

const handleUploadSuccess = () => {
  ElMessage.success('上传成功')
  loadFiles()
}
const handleUploadError = (error) => {
  try {
    const msg = error?.response?.data?.message || error?.response?.data?.error
    if (msg) {
      ElMessage.error(`上传失败: ${msg}`)
      return
    }
  } catch {}
  ElMessage.error('上传失败')
}
const beforeUpload = (file) => {
  const maxSizeMB = 50
  const isLtMax = file.size / 1024 / 1024 < maxSizeMB
  if (!isLtMax) {
    ElMessage.error(`文件过大，最大支持 ${maxSizeMB}MB`)
    return false
  }

  const ext = '.' + file.name.split('.').pop().toLowerCase()
  const allowed = ['.jpg','.jpeg','.png','.gif','.webp','.pdf','.doc','.docx','.ppt','.pptx','.xls','.xlsx','.txt','.mp4','.mov','.avi','.mp3','.wav','.zip','.rar','.7z']
  if (!allowed.includes(ext)) {
    ElMessage.error('不支持的文件类型')
    return false
  }
  return true
}

const createFolder = () => {
  ElMessage.info('文件夹功能正在开发中 (Coming Soon)')
}

// Utils
const isImage = (file) => ['.jpg', '.jpeg', '.png', '.gif', '.webp'].includes(file.extension?.toLowerCase())
const getFileUrl = (file) => `${config.api.baseURL}/files/download/${file.id}`

const getFileIcon = (ext) => {
  if (!ext) return Folder
  const e = ext.toLowerCase()
  if (['.jpg','.jpeg','.png','.gif'].includes(e)) return Picture
  if (['.mp4','.mov','.avi'].includes(e)) return VideoPlay
  if (['.mp3','.wav'].includes(e)) return Headset
  if (['.zip','.rar','.7z'].includes(e)) return Files
  return Document
}

const getFileColor = (ext) => {
  if (!ext) return '#E6A23C'
  const e = ext.toLowerCase()
  if (['.jpg','.jpeg','.png'].includes(e)) return '#67C23A'
  if (['.pdf','.ppt'].includes(e)) return '#F56C6C'
  if (['.doc','.docx'].includes(e)) return '#409EFF'
  if (['.zip','.rar'].includes(e)) return '#909399'
  return '#A0CFFF'
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
}

const formatTime = (t) => {
    if(!t) return '-'
    return new Date(t).toLocaleString()
}

onMounted(() => {
  loadFiles()
  document.addEventListener('click', closeContextMenu)
})
onUnmounted(() => document.removeEventListener('click', closeContextMenu))
</script>

<style scoped>
.file-manager {
  height: 100%;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
  gap: 16px;
}

/* Glass Panel Utility */
.glass-panel {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
}

/* Toolbar */
.file-toolbar {
  padding: 12px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
}
.toolbar-left, .toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}
.type-select { width: 130px; }
.search-input { width: 240px; }

/* Breadcrumb */
.breadcrumb-bar {
  padding: 12px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
}
.file-count { font-size: 12px; opacity: 0.6; }

/* Content Area */
.file-content {
  flex: 1;
  overflow-y: auto;
  position: relative;
  padding: 0 10px; /* space for scrollbar */
}

/* Grid View */
.file-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 16px;
  padding-bottom: 40px;
}

.file-card {
  background: rgba(255,255,255,0.05);
  border: 1px solid transparent;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s ease;
}
.file-card:hover {
  background: rgba(255,255,255,0.1);
  transform: translateY(-2px);
}
.file-card.is-selected {
  background: rgba(64, 158, 255, 0.15);
  border-color: #409EFF;
}

.file-preview {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
}
.thumbnail-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
}

.file-name {
  font-size: 13px;
  font-weight: 500;
  text-align: center;
  width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}
.file-meta { font-size: 11px; opacity: 0.5; }

/* List View */
.file-list-wrapper :deep(.el-table) {
  --el-table-bg-color: transparent;
  --el-table-tr-bg-color: transparent;
  --el-table-header-bg-color: rgba(255,255,255,0.05);
  --el-table-text-color: #eee;
  --el-table-header-text-color: #ddd;
  --el-table-row-hover-bg-color: rgba(255,255,255,0.1);
  --el-table-border-color: rgba(255,255,255,0.1);
}
.list-name { display: flex; align-items: center; cursor: pointer; }

/* Drag Overlay */
.drag-overlay {
  position: absolute;
  inset: 0;
  background: rgba(64, 158, 255, 0.2);
  backdrop-filter: blur(4px);
  border: 2px dashed #409EFF;
  border-radius: 16px;
  z-index: 999;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  pointer-events: none; /* Let events pass through to container */
}
.drag-text { font-size: 18px; font-weight: bold; margin-top: 16px; }

/* Context Menu */
.context-menu {
  position: fixed;
  background: #1e1e1e;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 6px 0;
  min-width: 160px;
  z-index: 9999;
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
}
.menu-item {
  padding: 10px 16px;
  font-size: 13px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #eee;
}
.menu-item:hover { background: #409EFF; color: #fff; }
.menu-item.danger { color: #F56C6C; }
.menu-item.danger:hover { background: #F56C6C; color: #fff; }
.menu-divider { height: 1px; background: #333; margin: 4px 0; }

/* Drawer */
.file-detail { padding: 20px; display: flex; flex-direction: column; height: 100%; }
.detail-preview { 
  height: 200px; 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  background: rgba(0,0,0,0.2); 
  border-radius: 8px; 
  margin-bottom: 20px; 
}
.detail-img { max-width: 100%; max-height: 100%; border-radius: 4px; }
.detail-info p { margin: 10px 0; font-size: 14px; color: #ccc; }
.detail-actions { margin-top: auto; display: flex; flex-direction: column; gap: 10px; }
</style>
