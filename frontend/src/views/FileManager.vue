<template>
  <div 
    class="file-manager" 
    @dragover.prevent="handleDragOver" 
    @drop.prevent="handleDrop" 
    @dragleave="handleDragLeave"
    @mousedown="handleDragStart"
    @mousemove="handleDragMove"
    @mouseup="handleDragEnd"
    @mouseleave="handleDragEnd"
  >
    <!-- 拖拽上传覆盖层 -->
    <div v-if="isDragging" class="drag-overlay">
      <el-icon :size="60"><UploadFilled /></el-icon>
      <div class="drag-text">释放以上传文件</div>
    </div>
    
    <!-- 拖拽选择区域 -->
    <div 
      v-if="dragSelectionVisible" 
      class="drag-selection"
      :style="{
        left: `${Math.min(dragStart.x, dragEnd.x)}px`,
        top: `${Math.min(dragStart.y, dragEnd.y)}px`,
        width: `${Math.abs(dragEnd.x - dragStart.x)}px`,
        height: `${Math.abs(dragEnd.y - dragStart.y)}px`
      }"
    ></div>

    <!-- 顶部工具栏 -->
    <div class="file-toolbar glass-panel">
      <div class="toolbar-left">
        <el-upload
          class="upload-btn"
          :auto-upload="false"
          :multiple="true"
          :show-file-list="false"
          :on-change="handleFileChange"
          :before-upload="beforeUpload"
          :disabled="isUploading"
        >
          <el-button type="primary" :icon="UploadFilled" :disabled="isUploading">上传</el-button>
        </el-upload>
        <el-button :icon="FolderAdd" @click="createFolder">新建文件夹</el-button>
        <el-select 
          v-model="selectedStorage" 
          placeholder="存储方式" 
          size="small"
          class="storage-select"
        >
          <el-option 
            v-for="option in storageOptions" 
            :key="option.value" 
            :label="option.label" 
            :value="option.value" 
          />
        </el-select>
      </div>
      <div class="toolbar-right">
        <el-button :icon="Setting" @click="storageConfigVisible = true">存储配置</el-button>
        <el-button 
          :icon="RefreshRight"
          :type="'primary'"
          :disabled="isSyncing"
          @click="startSync"
        >
          同步文件
          <el-progress 
            v-if="isSyncing" 
            type="circle" 
            :percentage="syncProgress" 
            :size="18"
            :stroke-width="3"
            :show-text="false"
            style="margin-left: 8px;"
          />
        </el-button>
        <!-- 批量操作按钮 -->
        <div v-if="selectedFiles.length > 0" class="batch-actions">
          <el-button 
            type="danger" 
            :icon="Delete" 
            @click="batchDeleteFiles"
            :disabled="selectedFiles.length === 0"
          >
            批量删除 ({{ selectedFiles.length }})
          </el-button>
          <el-button 
            type="warning" 
            @click="clearSelection"
          >
            取消选择
          </el-button>
        </div>
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
            :data-file-id="file.id"
            :class="{ 'is-selected': selectedFile?.id === file.id, 'is-batch-selected': selectedFiles.includes(file.id) }"
            @click.stop="selectFile(file)"
            @dblclick="handleFileDblClick(file)"
            @contextmenu.prevent="showContextMenu($event, file)"
          >
            <!-- 批量选择复选框 -->
            <div class="batch-select-checkbox">
              <el-checkbox 
                v-model="selectedFiles" 
                :label="file.id" 
                @change="handleBatchSelectChange(file)"
                @click.stop
              />
            </div>
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
          @selection-change="handleSelectionChange"
          highlight-current-row
        >
          <el-table-column type="selection" width="55" />
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

    <!-- 批量上传对话框 -->
    <el-dialog
      v-model="isUploading"
      title="批量上传"
      width="600px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
    >
      <div class="upload-progress">
        <div class="progress-header">
          <div class="progress-info">
            <span>上传进度</span>
            <span>{{ uploadProgress }}%</span>
          </div>
          <el-progress :percentage="uploadProgress" :stroke-width="8" :show-text="false"></el-progress>
        </div>
        <div class="file-list">
          <div v-for="(file, index) in uploadingFiles" :key="index" class="file-item">
            <div class="file-info">
              <el-icon v-if="file.status === 'success'" color="#67C23A"><CircleCheck /></el-icon>
              <el-icon v-else-if="file.status === 'error'" color="#F56C6C"><WarningFilled /></el-icon>
              <el-icon v-else-if="file.status === 'uploading'" color="#409EFF"><Loading /></el-icon>
              <el-icon v-else color="#909399"><Clock /></el-icon>
              <span class="file-name">{{ file.name }}</span>
            </div>
            <div class="file-size">{{ formatFileSize(file.size) }}</div>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="upload-footer">
          <el-button 
            v-if="isUploading" 
            type="danger" 
            @click="cancelUpload" 
            :disabled="uploadingFiles.every(f => f.status !== 'uploading')"
          >
            取消上传
          </el-button>
          <el-button 
            v-else 
            type="primary" 
            @click="uploadProgress = 0; uploadingFiles = []"
          >
            完成
          </el-button>
        </div>
      </template>
    </el-dialog>

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
    
    <!-- 存储配置对话框 -->
    <el-dialog
      v-model="storageConfigVisible"
      title="对象存储配置"
      width="800px"
      center
    >
      <div class="storage-config">
        <div class="setting-item">
          <label>默认存储方式</label>
          <el-select v-model="storageSettings.defaultStorage" placeholder="选择默认存储方式">
            <el-option label="本地存储" value="local" />
            <el-option label="阿里云OSS" value="oss" />
            <el-option label="腾讯云COS" value="cos" />
            <el-option label="AWS S3" value="s3" />
            <el-option label="七牛云Kodo" value="qiniu" />
          </el-select>
        </div>

        <div class="tip-text">
          <el-icon color="#409EFF"><WarningFilled /></el-icon>
          <span>所有对象存储类型都将统一使用S3存储服务</span>
        </div>

        <!-- 存储配置折叠面板 -->
        <el-collapse v-model="activeStoragePanel">
          <!-- AWS S3配置 -->
          <el-collapse-item title="S3存储配置" name="s3">
            <div class="setting-item">
              <label>Endpoint</label>
              <el-input v-model="storageSettings.providers.s3.endpoint" placeholder="S3服务端点 (如: https://s3.amazonaws.com)" />
            </div>
            <div class="setting-item">
              <label>Access Key ID</label>
              <el-input v-model="storageSettings.providers.s3.accessKeyId" placeholder="访问密钥ID" />
            </div>
            <div class="setting-item">
              <label>Secret Access Key</label>
              <el-input v-model="storageSettings.providers.s3.secretAccessKey" type="password" placeholder="访问密钥Secret" show-password />
            </div>
            <div class="setting-item">
              <label>Bucket</label>
              <el-input v-model="storageSettings.providers.s3.bucket" placeholder="存储桶名称" />
            </div>
            <div class="setting-item">
              <label>Region</label>
              <el-input v-model="storageSettings.providers.s3.region" placeholder="存储区域 (如: us-east-1)" />
            </div>
            <div class="setting-item">
              <label>启用S3存储</label>
              <el-switch v-model="storageSettings.providers.s3.enabled" />
            </div>
          </el-collapse-item>
          
          <!-- 七牛云配置 -->
          <el-collapse-item title="七牛云存储配置" name="qiniu">
            <div class="setting-item">
              <label>Access Key</label>
              <el-input v-model="storageSettings.providers.qiniu.accessKey" placeholder="七牛云Access Key" />
            </div>
            <div class="setting-item">
              <label>Secret Key</label>
              <el-input v-model="storageSettings.providers.qiniu.secretKey" type="password" placeholder="七牛云Secret Key" show-password />
            </div>
            <div class="setting-item">
              <label>Bucket</label>
              <el-input v-model="storageSettings.providers.qiniu.bucket" placeholder="存储桶名称" />
            </div>
            <div class="setting-item">
              <label>Region</label>
              <el-input v-model="storageSettings.providers.qiniu.region" placeholder="存储区域 (如: cn-east-1)" />
            </div>
            <div class="setting-item">
              <label>Domain</label>
              <el-input v-model="storageSettings.providers.qiniu.domain" placeholder="自定义域名" />
            </div>
            <div class="setting-item">
              <label>启用七牛云存储</label>
              <el-switch v-model="storageSettings.providers.qiniu.enabled" />
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="storageConfigVisible = false">取消</el-button>
          <el-button type="primary" @click="saveStorageSettings">保存配置</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  UploadFilled, FolderAdd, Search, Menu, List, 
  Document, Picture, VideoPlay, Files, Folder,
  HomeFilled, View, Download, Edit, Delete, Headset,
  CircleCheck, WarningFilled, Loading, Clock, Setting,
  RefreshRight
} from '@element-plus/icons-vue'
import api from '../api'
import config from '../config'
import StorageServiceFactory from '../services/storageService'

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

// 存储方式选择
const selectedStorage = ref(config.storageProviders.default)
const storageOptions = ref(Object.entries(config.storageProviders.providers).map(([key, provider]) => ({
  label: provider.name,
  value: key
})))

// 存储配置相关状态
const storageConfigVisible = ref(false)
const storageSettings = ref({
  defaultStorage: 'local',
  providers: {
    oss: {
      endpoint: '',
      accessKeyId: '',
      accessKeySecret: '',
      bucket: '',
      region: ''
    },
    cos: {
      secretId: '',
      secretKey: '',
      region: '',
      bucket: ''
    },
    s3: {
      endpoint: '',
      accessKeyId: '',
      secretAccessKey: '',
      region: '',
      bucket: '',
      enabled: false
    },
    qiniu: {
      accessKey: '',
      secretKey: '',
      bucket: '',
      region: '',
      domain: '',
      enabled: false
    }
  }
})

// 折叠面板激活项
const activeStoragePanel = ref([])

// 同步相关状态
const isSyncing = ref(false)
const syncProgress = ref(0)
const syncStatus = ref('idle') // idle, syncing, success, error
const syncResult = ref(null)

// 批量操作相关状态
const selectedFiles = ref([])
const isSelectAll = ref(false)

// 拖拽选中相关状态
const isDragSelecting = ref(false)
const dragStart = ref({ x: 0, y: 0 })
const dragEnd = ref({ x: 0, y: 0 })
const dragSelectionVisible = ref(false)

// 批量上传相关状态
const uploadProgress = ref(0)
const uploadingFiles = ref([])
const isUploading = ref(false)

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

// 文件选择处理
const handleFileChange = (file, fileList) => {
  // 处理从el-upload组件选择的文件
  if (file.raw) {
    // 如果是单个文件，先验证再上传
    if (beforeUpload(file)) {
      startBatchUpload([file.raw])
    }
  } else if (fileList && fileList.length > 0) {
    // 如果是多个文件，先过滤通过验证的文件
    const files = fileList.map(f => f.raw).filter(Boolean).filter(f => {
      // 为原始文件创建一个类似el-upload的文件对象，以便beforeUpload处理
      const uploadFile = { ...f, name: f.name, size: f.size }
      return beforeUpload(uploadFile)
    })
    if (files.length > 0) {
      startBatchUpload(files)
    }
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

  // 过滤通过验证的文件
  const validFiles = Array.from(droppedFiles).filter(file => {
    // 为原始文件创建一个类似el-upload的文件对象，以便beforeUpload处理
    const uploadFile = { ...file, name: file.name, size: file.size }
    return beforeUpload(uploadFile)
  })
  
  if (validFiles.length > 0) {
    startBatchUpload(validFiles)
  }
}

const uploadSingleFile = async (file) => {
  try {
    // 调试信息：检查文件信息和存储配置
    console.log('Uploading file:', file.name)
    console.log('Selected storage:', selectedStorage.value)
    console.log('Storage settings:', storageSettings.value)
    
    // 检查是否选择了S3存储但未启用
    if (selectedStorage.value === 's3' && !storageSettings.value.providers.s3.enabled) {
      ElMessage.error(`上传失败: ${file.name}，S3存储未启用，请先在存储配置中启用S3存储`)
      return false
    }
    
    // 获取存储服务实例
    const storageService = StorageServiceFactory.getStorageService(selectedStorage.value)
    
    // 上传文件
    await storageService.upload(file, {
      path: breadcrumbs.value.join('/')
    })
    
    return true
  } catch (err) {
    console.error('Upload Error:', err)
    // 提供更详细的错误信息
    let errorMsg = `上传失败: ${file.name}`
    if (err.response) {
      const { status, data } = err.response
      console.log('Response status:', status)
      console.log('Response data:', data)
      if (status === 403) {
        errorMsg = `上传失败: ${file.name}，权限不足`
      } else if (status === 413) {
        errorMsg = `上传失败: ${file.name}，文件过大`
      } else if (status === 404) {
        errorMsg = `上传失败: ${file.name}，后端API不存在，请检查服务是否正常`
      } else if (status === 500) {
        errorMsg = `上传失败: ${file.name}，后端服务错误，请检查服务器日志`
      } else {
        errorMsg = `上传失败: ${file.name}，错误信息: ${data?.message || data?.error || '未知错误'}`
      }
    } else if (err.request) {
      console.log('No response received:', err.request)
      errorMsg = `上传失败: ${file.name}，网络连接失败，请检查后端服务是否正常运行`
    } else {
      console.log('Request error:', err.message)
      errorMsg = `上传失败: ${file.name}，错误: ${err.message}`
    }
    ElMessage.error(errorMsg)
    return false
  }
}

// 批量上传处理
const startBatchUpload = async (files) => {
  isUploading.value = true
  uploadProgress.value = 0
  uploadingFiles.value = files.map(file => ({
    name: file.name,
    size: file.size,
    status: 'pending', // pending, uploading, success, error
    progress: 0,
    error: '',
    fileRef: file // 保存文件引用以供取消使用
  }))

  // 异步上传文件
  const totalFiles = files.length
  let uploadedCount = 0
  let failedCount = 0
  
  for (let i = 0; i < files.length; i++) {
    const file = files[i]
    const fileIndex = i
    
    // 更新状态为正在上传
    uploadingFiles.value[fileIndex].status = 'uploading'
    
    try {
      // 模拟进度更新（实际项目中可能需要从服务器获取真实进度）
      const progressInterval = setInterval(() => {
        if (uploadingFiles.value[fileIndex].progress < 90) {
          uploadingFiles.value[fileIndex].progress += 10
        }
      }, 200)
      
      // 上传文件
      const success = await uploadSingleFile(file)
      
      // 完成上传，更新进度
      clearInterval(progressInterval)
      if (success) {
        uploadingFiles.value[fileIndex].status = 'success'
        uploadingFiles.value[fileIndex].progress = 100
        uploadedCount++
      } else {
        uploadingFiles.value[fileIndex].status = 'error'
        uploadingFiles.value[fileIndex].progress = 100
        failedCount++
      }
      
      uploadProgress.value = Math.round(((uploadedCount + failedCount) / totalFiles) * 100)
      
    } catch (err) {
      uploadingFiles.value[fileIndex].status = 'error'
      uploadingFiles.value[fileIndex].progress = 100
      uploadingFiles.value[fileIndex].error = err.message || '上传失败'
      failedCount++
      uploadProgress.value = Math.round(((uploadedCount + failedCount) / totalFiles) * 100)
    }
  }
  
  // 上传完成
  isUploading.value = false
  
  // 显示上传结果
  if (uploadedCount > 0) {
    ElMessage.success(`成功上传 ${uploadedCount} 个文件`)
  }
  if (failedCount > 0) {
    ElMessage.error(`上传失败 ${failedCount} 个文件`)
  }
  
  // 重新加载文件列表
  loadFiles()
}

// 取消上传
const cancelUpload = () => {
  // 重置上传状态
  isUploading.value = false
  uploadProgress.value = 0
  uploadingFiles.value = []
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
  const maxSizeMB = config.upload.maxSize / 1024 / 1024
  const isLtMax = file.size / 1024 / 1024 < maxSizeMB
  if (!isLtMax) {
    ElMessage.error(`文件过大，最大支持 ${maxSizeMB}MB`)
    return false
  }

  // 获取文件扩展名，带点号和不带点号两种格式
  const fileName = file.name
  const extWithoutDot = fileName.split('.').pop().toLowerCase()
  const extWithDot = `.${extWithoutDot}`
  
  // 检查是否是思源笔记配置中排除的扩展名
  if (config.upload.思源笔记?.excludeExtensions?.includes(extWithDot)) {
    ElMessage.error(`不支持的文件类型: ${extWithoutDot}，已被配置排除`)
    return false
  }

  // 检查是否允许所有文件扩展名
  if (config.upload.allowedExtensions) {
    return true
  }
  
  // 使用配置文件中的允许类型
  const allAllowedTypes = [
    ...config.upload.acceptedTypes.image,
    ...config.upload.acceptedTypes.document,
    ...config.upload.acceptedTypes.code,
    ...config.upload.acceptedTypes.audio,
    ...config.upload.acceptedTypes.video,
    ...config.upload.acceptedTypes.archive,
    ...(config.upload.acceptedTypes.note || [])
  ]
  
  // 检查文件扩展名是否在允许列表中（同时检查带点号和不带点号的格式）
  const isAllowed = allAllowedTypes.some(type => {
    const typeWithoutDot = type.replace('.', '').toLowerCase()
    return extWithoutDot === typeWithoutDot || extWithDot === type.toLowerCase()
  })
  
  if (!isAllowed) {
    ElMessage.error(`不支持的文件类型: ${extWithoutDot}，请上传支持的文件格式`)
    return false
  }
  return true
}

const createFolder = () => {
  ElMessage.info('文件夹功能正在开发中 (Coming Soon)')
}

/**
 * 开始文件同步
 */
const startSync = async () => {
  try {
    // 检查七牛云配置是否已启用
    if (!storageSettings.value.providers.qiniu.enabled) {
      ElMessage.error('七牛云存储未启用，请先在存储配置中启用')
      return
    }
    
    isSyncing.value = true
    syncProgress.value = 0
    syncStatus.value = 'syncing'
    
    // 获取存储服务实例
    const storageService = StorageServiceFactory.getStorageService('qiniu')
    
    // 执行同步操作
    syncResult.value = await storageService.sync({
      path: breadcrumbs.value.join('/')
    })
    
    // 同步完成
    syncStatus.value = 'success'
    syncProgress.value = 100
    
    // 重新加载文件列表
    loadFiles()
  } catch (err) {
    console.error('Sync Error:', err)
    syncStatus.value = 'error'
    ElMessage.error('同步失败，请检查控制台日志获取详细信息')
  } finally {
    isSyncing.value = false
  }
}

/**
 * 监听同步状态变化
 */
const handleSyncStatusChange = (event) => {
  const { storageType, status } = event.detail
  if (storageType === 'qiniu') {
    isSyncing.value = status.isSyncing
    syncProgress.value = status.progress
    syncStatus.value = status.status
  }
}

/**
 * 处理列表视图选择变化
 */
const handleSelectionChange = (selection) => {
  selectedFiles.value = selection.map(file => file.id)
  isSelectAll.value = selection.length === filteredFiles.value.length
}

/**
 * 处理网格视图选择变化
 */
const handleBatchSelectChange = (file) => {
  // 处理单选逻辑
  if (selectedFiles.value.includes(file.id)) {
    // 取消选择
    selectedFiles.value = selectedFiles.value.filter(id => id !== file.id)
  } else {
    // 添加选择
    selectedFiles.value.push(file.id)
  }
  
  // 更新全选状态
  isSelectAll.value = selectedFiles.value.length === filteredFiles.value.length
}

/**
 * 批量删除文件
 */
const batchDeleteFiles = async () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning('请先选择要删除的文件')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedFiles.value.length} 个文件吗？此操作不可恢复。`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 执行批量删除操作
    const deletePromises = selectedFiles.value.map(fileId => api.delete(`/files/${fileId}`))
    await Promise.all(deletePromises)
    
    ElMessage.success(`成功删除 ${selectedFiles.value.length} 个文件`)
    
    // 清除选择
    selectedFiles.value = []
    isSelectAll.value = false
    
    // 重新加载文件列表
    loadFiles()
  } catch (err) {
    if (err !== 'cancel') {
      console.error('Batch Delete Error:', err)
      ElMessage.error('批量删除失败，请检查控制台日志获取详细信息')
    }
  }
}

/**
 * 清除选择
 */
const clearSelection = () => {
  selectedFiles.value = []
  isSelectAll.value = false
}

/**
 * 切换全选状态
 */
const toggleSelectAll = () => {
  if (isSelectAll.value) {
    // 取消全选
    selectedFiles.value = []
  } else {
    // 全选
    selectedFiles.value = filteredFiles.value.map(file => file.id)
  }
  isSelectAll.value = !isSelectAll.value
}

/**
 * 处理拖拽开始
 */
const handleDragStart = (event) => {
  // 只处理左键拖拽
  if (event.button !== 0) return
  
  // 如果点击的是文件卡片或复选框，不启动拖拽选择
  if (event.target.closest('.file-card') || event.target.closest('.el-checkbox')) {
    return
  }
  
  // 记录拖拽开始位置
  dragStart.value = { x: event.clientX, y: event.clientY }
  dragEnd.value = { ...dragStart.value }
  isDragSelecting.value = true
  dragSelectionVisible.value = true
}

/**
 * 处理拖拽移动
 */
const handleDragMove = (event) => {
  if (!isDragSelecting.value) return
  
  // 更新拖拽结束位置
  dragEnd.value = { x: event.clientX, y: event.clientY }
}

/**
 * 处理拖拽结束
 */
const handleDragEnd = () => {
  if (!isDragSelecting.value) return
  
  // 计算拖拽区域
  const rect = {
    left: Math.min(dragStart.value.x, dragEnd.value.x),
    top: Math.min(dragStart.value.y, dragEnd.value.y),
    right: Math.max(dragStart.value.x, dragEnd.value.x),
    bottom: Math.max(dragStart.value.y, dragEnd.value.y)
  }
  
  // 获取所有文件卡片元素
  const fileCards = document.querySelectorAll('.file-card')
  
  // 检查每个文件卡片是否在拖拽区域内
  fileCards.forEach(card => {
    const cardRect = card.getBoundingClientRect()
    
    // 检查卡片是否与拖拽区域相交
    if (cardRect.left < rect.right && 
        cardRect.right > rect.left && 
        cardRect.top < rect.bottom && 
        cardRect.bottom > rect.top) {
      // 获取文件ID
      const fileId = card.dataset.fileId
      if (fileId && !selectedFiles.value.includes(fileId)) {
        selectedFiles.value.push(fileId)
      }
    }
  })
  
  // 更新全选状态
  isSelectAll.value = selectedFiles.value.length === filteredFiles.value.length
  
  // 重置拖拽状态
  isDragSelecting.value = false
  dragSelectionVisible.value = false
}

/**
 * 处理文件卡片渲染，添加data-file-id属性
 */
const handleFileCardRender = (file) => {
  // 这个方法用于在模板中添加data-file-id属性
  return file.id
}

// Utils
const getFileExtension = (file) => {
  if (file.extension) {
    let ext = (file.extension || '').toLowerCase()
    // 确保扩展名带有点号，与config.upload.acceptedTypes格式一致
    return ext.startsWith('.') ? ext : '.' + ext
  }
  if (file.file_name) {
    return '.' + file.file_name.split('.').pop().toLowerCase()
  }
  return ''
}

// 检查文件类型的通用函数，支持带点号和不带点号的扩展名
const checkFileType = (ext, typeList) => {
  const normalizedExt = ext.toLowerCase()
  // 检查带点号和不带点号两种格式
  return typeList.some(type => {
    const normalizedType = type.toLowerCase()
    return normalizedExt === normalizedType || normalizedExt === normalizedType.replace('.', '')
  })
}

const isImage = (file) => {
  const ext = getFileExtension(file)
  return checkFileType(ext, config.upload.acceptedTypes.image)
}

const isAudio = (file) => {
  const ext = getFileExtension(file)
  return checkFileType(ext, config.upload.acceptedTypes.audio)
}

const isVideo = (file) => {
  const ext = getFileExtension(file)
  return checkFileType(ext, config.upload.acceptedTypes.video)
}

const isArchive = (file) => {
  const ext = getFileExtension(file)
  return checkFileType(ext, config.upload.acceptedTypes.archive)
}

const isDocument = (file) => {
  const ext = getFileExtension(file)
  return checkFileType(ext, config.upload.acceptedTypes.document)
}

const isCode = (file) => {
  const ext = getFileExtension(file)
  return checkFileType(ext, config.upload.acceptedTypes.code)
}

const getFileUrl = (file) => `${config.api.baseURL}/files/download/${file.id}`

const getFileIcon = (ext) => {
  if (!ext) return Folder
  
  // 图片
  if (checkFileType(ext, config.upload.acceptedTypes.image)) return Picture
  
  // 视频
  if (checkFileType(ext, config.upload.acceptedTypes.video)) return VideoPlay
  
  // 音频
  if (checkFileType(ext, config.upload.acceptedTypes.audio)) return Headset
  
  // 文档
  if (checkFileType(ext, config.upload.acceptedTypes.document)) return Document
  
  // 代码
  if (checkFileType(ext, config.upload.acceptedTypes.code)) return Document
  
  // 压缩文件
  if (checkFileType(ext, config.upload.acceptedTypes.archive)) return Files
  
  // 默认文档图标
  return Document
}

const getFileColor = (ext) => {
  if (!ext) return '#E6A23C'
  const normalizedExt = ext.toLowerCase()
  
  // 图片
  if (checkFileType(ext, config.upload.acceptedTypes.image)) return '#67C23A'
  
  // 文档
  if (checkFileType(ext, config.upload.acceptedTypes.document)) {
    if (normalizedExt === '.pdf' || normalizedExt === 'pdf') return '#F56C6C' // PDF红色
    if (normalizedExt === '.doc' || normalizedExt === 'doc' || normalizedExt === '.docx' || normalizedExt === 'docx') return '#409EFF' // Word蓝色
    if (normalizedExt === '.xls' || normalizedExt === 'xls' || normalizedExt === '.xlsx' || normalizedExt === 'xlsx') return '#67C23A' // Excel绿色
    if (normalizedExt === '.ppt' || normalizedExt === 'ppt' || normalizedExt === '.pptx' || normalizedExt === 'pptx') return '#E6A23C' // PowerPoint橙色
    return '#909399' // 其他文档
  }
  
  // 代码
  if (checkFileType(ext, config.upload.acceptedTypes.code)) return '#A0CFFF'
  
  // 视频
  if (checkFileType(ext, config.upload.acceptedTypes.video)) return '#909399'
  
  // 音频
  if (checkFileType(ext, config.upload.acceptedTypes.audio)) return '#909399'
  
  // 压缩文件
  if (checkFileType(ext, config.upload.acceptedTypes.archive)) return '#909399'
  
  // 默认颜色
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

// 加载存储设置
const loadStorageSettings = () => {
  const saved = localStorage.getItem('storage_settings')
  if (saved) {
    try {
      storageSettings.value = JSON.parse(saved)
      // 更新当前选择的存储方式
      selectedStorage.value = storageSettings.value.defaultStorage
    } catch {
      console.error('Failed to parse storage settings')
    }
  }
}

// 保存存储设置
const saveStorageSettings = () => {
  try {
    localStorage.setItem('storage_settings', JSON.stringify(storageSettings.value))
    // 更新当前选择的存储方式
    selectedStorage.value = storageSettings.value.defaultStorage
    ElMessage.success('存储设置已保存')
    // 关闭对话框
    storageConfigVisible.value = false
  } catch {
    ElMessage.error('保存存储设置失败')
  }
}

onMounted(() => {
  loadFiles()
  loadStorageSettings()
  document.addEventListener('click', closeContextMenu)
  // 添加同步状态监听器
  window.addEventListener('storage:syncStatus', handleSyncStatusChange)
})
onUnmounted(() => {
  document.removeEventListener('click', closeContextMenu)
  // 移除同步状态监听器
  window.removeEventListener('storage:syncStatus', handleSyncStatusChange)
})
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

/* Glass Panel Utility - 使用主题变量，移除背景图片 */
.glass-panel {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
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
.storage-select { width: 150px; }
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
  position: relative;
}
.file-card:hover {
  background: rgba(255,255,255,0.1);
  transform: translateY(-2px);
}
.file-card.is-selected {
  background: rgba(64, 158, 255, 0.15);
  border-color: #409EFF;
}
.file-card.is-batch-selected {
  background: rgba(245, 108, 108, 0.15);
  border-color: #F56C6C;
}

/* 批量选择复选框样式 */
.batch-select-checkbox {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 4px;
  padding: 4px;
  transform: scale(1.2); /* 增大复选框大小 */
}

/* 调整Element Plus复选框样式，增大点击区域 */
:deep(.el-checkbox) {
  --el-checkbox-size: 20px; /* 调整复选框大小 */
}

:deep(.el-checkbox__input) {
  width: 20px;
  height: 20px;
}

:deep(.el-checkbox__inner) {
  width: 20px;
  height: 20px;
}

:deep(.el-checkbox__input.is-checked .el-checkbox__inner::after) {
  left: 7px;
  top: 3px;
  width: 6px;
  height: 12px;
}

/* 批量操作按钮样式 */
.batch-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

/* 工具栏样式调整 */
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

/* 拖拽选择区域样式 */
.drag-selection {
  position: fixed;
  border: 2px solid #409EFF;
  background-color: rgba(64, 158, 255, 0.2);
  pointer-events: none;
  z-index: 9999;
  box-shadow: 0 0 10px rgba(64, 158, 255, 0.3);
  border-radius: 4px;
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

/* 存储配置样式 */
.storage-config {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.tip-text {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: rgba(64, 158, 255, 0.1);
  border: 1px solid rgba(64, 158, 255, 0.3);
  border-radius: 8px;
  font-size: 14px;
  color: #409EFF;
  margin: 12px 0;
}

.setting-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.setting-item label {
  font-weight: 500;
  color: var(--text-primary);
  font-size: 14px;
}

/* 批量上传对话框样式 */
.upload-progress {
  max-height: 400px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.progress-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.file-list {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 8px;
  padding: 8px;
}

.file-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  border-bottom: 1px solid rgba(255,255,255,0.05);
}

.file-item:last-child {
  border-bottom: none;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 8px;
  overflow: hidden;
}

.file-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 350px;
}

.file-size {
  font-size: 12px;
  opacity: 0.7;
}

.upload-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}
</style>
