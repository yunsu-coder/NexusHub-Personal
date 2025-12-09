<template>
  <div class="file-manager">
    <!-- 上传区域 -->
    <el-card class="upload-section">
      <el-upload
        ref="uploadRef"
        class="upload-dragger"
        drag
        :action="uploadUrl"
        :headers="uploadHeaders"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        :show-file-list="false"
        multiple
        name="file"
      >
        <div class="upload-content">
          <el-icon class="upload-icon"><el-icon-upload-filled /></el-icon>
          <div class="el-upload__text">
            拖拽文件到此处或 <em>点击上传</em>
          </div>
          <div class="el-upload__tip">
            支持所有文件类型，单个文件最大 500MB
          </div>
        </div>
      </el-upload>
    </el-card>

    <!-- 文件分类标签 -->
    <el-card class="category-section">
      <el-radio-group v-model="selectedCategory" @change="loadFiles">
        <el-radio-button label="">全部</el-radio-button>
        <el-radio-button label="media">媒体</el-radio-button>
        <el-radio-button label="document">文档</el-radio-button>
        <el-radio-button label="code">代码</el-radio-button>
        <el-radio-button label="archive">压缩包</el-radio-button>
        <el-radio-button label="other">其他</el-radio-button>
      </el-radio-group>
    </el-card>

    <!-- 文件列表 -->
    <el-card class="files-section">
      <template #header>
        <div class="card-header">
          <div class="search-box">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索文件名"
              prefix-icon="Search"
              clearable
              style="width: 300px"
            />
          </div>
          <el-button type="primary" :icon="Refresh" @click="loadFiles">刷新</el-button>
        </div>
      </template>

      <el-table :data="filteredFiles" style="width: 100%" v-loading="loading">
        <el-table-column label="文件名" min-width="200">
          <template #default="{ row }">
            <div class="file-name">
              <el-icon :size="20">
                <component :is="getFileIcon(row.extension)" />
              </el-icon>
              <span>{{ row.file_name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="category" label="分类" width="100">
          <template #default="{ row }">
            <el-tag :type="getCategoryType(row.category)">
              {{ getCategoryName(row.category) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="file_size" label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.file_size) }}
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="上传时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              size="small"
              type="primary"
              :icon="View"
              @click="previewFile(row)"
            >
              预览
            </el-button>
            <el-button
              size="small"
              type="success"
              :icon="Download"
              @click="downloadFile(row)"
            >
              下载
            </el-button>
            <el-button
              size="small"
              type="danger"
              :icon="Delete"
              @click="deleteFile(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 文件预览对话框 -->
    <el-dialog
      v-model="previewVisible"
      :title="currentPreviewFile?.file_name"
      width="80%"
      @close="closePreview"
    >
      <div class="preview-content">
        <!-- 图片预览 -->
        <el-image
          v-if="isImage(currentPreviewFile)"
          :src="getFileUrl(currentPreviewFile)"
          fit="contain"
          style="max-width: 100%; max-height: 70vh"
        />

        <!-- 视频预览 -->
        <video
          v-else-if="isVideo(currentPreviewFile)"
          :src="getFileUrl(currentPreviewFile)"
          controls
          style="max-width: 100%; max-height: 70vh"
        />

        <!-- 音频预览 -->
        <audio
          v-else-if="isAudio(currentPreviewFile)"
          :src="getFileUrl(currentPreviewFile)"
          controls
          style="max-width: 100%; margin: 20px 0"
        />

        <!-- 代码预览（带语法高亮） -->
        <pre v-else-if="isCode(currentPreviewFile)" class="code-preview">
          <code :class="`language-${getCodeLanguage(currentPreviewFile.extension)}`">{{ previewContent }}</code>
        </pre>

        <!-- 文本预览 -->
        <pre v-else-if="isText(currentPreviewFile)" class="text-preview">{{ previewContent }}</pre>

        <!-- 不支持预览 -->
        <el-empty v-else description="此文件类型不支持预览" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  UploadFilled,
  Refresh,
  View,
  Download,
  Delete,
  Document,
  Picture,
  VideoPlay,
  Files,
  Folder,
  Search
} from '@element-plus/icons-vue'
import api from '../api'
import config from '../config'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

const uploadRef = ref()
const loading = ref(false)
const selectedCategory = ref('')
const files = ref([])
const previewVisible = ref(false)
const currentPreviewFile = ref(null)
const previewContent = ref('')
const searchKeyword = ref('')
const filteredFiles = computed(() => {
  if (!searchKeyword.value) return files.value
  return files.value.filter(file => 
    file.file_name.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

// 上传请求头 - 包含JWT token
const uploadHeaders = computed(() => ({ }))

const uploadUrl = computed(() => `${config.api.baseURL}/files/upload`)

const loadFiles = async () => {
  loading.value = true
  try {
    const url = selectedCategory.value
      ? `/files/category/${selectedCategory.value}`
      : '/files'
    files.value = await api.get(url)
  } catch (error) {
    ElMessage.error('加载文件列表失败')
  } finally {
    loading.value = false
  }
}

const handleUploadSuccess = () => {
  ElMessage.success('上传成功')
  loadFiles()
}

const handleUploadError = () => {
  ElMessage.error('上传失败')
}

const previewFile = async (file) => {
  currentPreviewFile.value = file
  previewVisible.value = true

  if (isText(file) || isCode(file)) {
    try {
      const response = await fetch(getFileUrl(file))
      previewContent.value = await response.text()
      // 延迟执行语法高亮，确保DOM已更新
      setTimeout(() => {
        if (isCode(file)) {
          hljs.highlightAll()
        }
      }, 0)
    } catch (error) {
      ElMessage.error('加载文件内容失败')
    }
  }
}

const closePreview = () => {
  previewVisible.value = false
  currentPreviewFile.value = null
  previewContent.value = ''
}

const downloadFile = (file) => {
  window.open(`${config.api.baseURL}/files/download/${file.id}`, '_blank')
}

const deleteFile = async (file) => {
  try {
    await ElMessageBox.confirm('确定要删除此文件吗？', '提示', {
      type: 'warning'
    })
    await api.delete(`/files/${file.id}`)
    ElMessage.success('删除成功')
    loadFiles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const getFileIcon = (ext) => {
  const imageExts = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg']
  const videoExts = ['.mp4', '.avi', '.mov', '.mkv', '.webm']
  const audioExts = ['.mp3', '.wav', '.ogg', '.flac', '.aac']
  const archiveExts = ['.zip', '.rar', '.7z', '.tar', '.gz']

  if (imageExts.includes(ext?.toLowerCase())) return Picture
    if (videoExts.includes(ext?.toLowerCase())) return VideoPlay
    if (audioExts.includes(ext?.toLowerCase())) return Document
    if (archiveExts.includes(ext?.toLowerCase())) return Files
    return Document
}

const getCategoryType = (category) => {
  const types = {
    media: 'success',
    document: 'primary',
    code: 'warning',
    archive: 'info',
    other: 'default'
  }
  return types[category] || 'default'
}

const getCategoryName = (category) => {
  const names = {
    media: '媒体',
    document: '文档',
    code: '代码',
    archive: '压缩包',
    other: '其他'
  }
  return names[category] || category
}

const formatFileSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
}

const formatTime = (time) => {
  return new Date(time).toLocaleString('zh-CN')
}

const getFileUrl = (file) => {
  return `${config.api.baseURL}/files/download/${file.id}`
}

const isImage = (file) => {
  return ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg'].includes(file.extension?.toLowerCase())
}

const isVideo = (file) => {
  return ['.mp4', '.avi', '.mov', '.mkv', '.webm'].includes(file.extension?.toLowerCase())
}

const isAudio = (file) => {
  return ['.mp3', '.wav', '.ogg', '.flac', '.aac'].includes(file.extension?.toLowerCase())
}

const isText = (file) => {
  return ['.txt', '.md', '.json', '.xml', '.csv', '.log', '.ini', '.yaml', '.yml', '.toml'].includes(file.extension?.toLowerCase())
}

const isCode = (file) => {
  return ['.html', '.htm', '.css', '.scss', '.less', '.js', '.jsx', '.ts', '.tsx', 
          '.java', '.c', '.cpp', '.h', '.hpp', '.go', '.py', '.php', '.rb', '.swift', 
          '.kt', '.rs', '.vue', '.svelte', '.lua', '.sql', '.sh', '.bash'].includes(file.extension?.toLowerCase())
}

const getCodeLanguage = (ext) => {
  const languageMap = {
    '.html': 'html',
    '.htm': 'html',
    '.css': 'css',
    '.scss': 'scss',
    '.less': 'less',
    '.js': 'javascript',
    '.jsx': 'javascript',
    '.ts': 'typescript',
    '.tsx': 'typescript',
    '.java': 'java',
    '.c': 'c',
    '.cpp': 'cpp',
    '.h': 'c',
    '.hpp': 'cpp',
    '.go': 'go',
    '.py': 'python',
    '.php': 'php',
    '.rb': 'ruby',
    '.swift': 'swift',
    '.kt': 'kotlin',
    '.rs': 'rust',
    '.vue': 'vue',
    '.svelte': 'svelte',
    '.lua': 'lua',
    '.sql': 'sql',
    '.sh': 'bash',
    '.bash': 'bash',
    '.json': 'json',
    '.xml': 'xml',
    '.yaml': 'yaml',
    '.yml': 'yaml',
    '.md': 'markdown'
  }
  return languageMap[ext?.toLowerCase()] || 'plaintext'
}

onMounted(() => {
  loadFiles()
})
</script>

<style scoped>
.file-manager {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.upload-section {
  background: linear-gradient(135deg, var(--card-bg) 0%, var(--hover-bg) 100%);
}

.upload-dragger {
  width: 100%;
}

.upload-content {
  padding: 40px;
  text-align: center;
}

.upload-icon {
  font-size: 67px;
  color: #409EFF;
  margin-bottom: 16px;
}

.el-upload__text {
  font-size: 16px;
  color: var(--text-primary);
}

.el-upload__tip {
  margin-top: 10px;
  font-size: 12px;
  color: var(--text-secondary);
}

.category-section :deep(.el-card__body) {
  padding: 15px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 10px;
}

.text-preview {
  max-height: 70vh;
  overflow: auto;
  padding: 20px;
  background: var(--hover-bg);
  border-radius: 4px;
  color: var(--text-primary);
  font-family: 'Consolas', 'Monaco', monospace;
}

.code-preview {
  max-height: 70vh;
  overflow: auto;
  padding: 20px;
  background: var(--hover-bg);
  border-radius: 4px;
  color: var(--text-primary);
}

.search-box {
  margin-bottom: 15px;
}

:deep(.el-table__header-wrapper th) {
  background-color: var(--hover-bg);
  color: var(--text-primary);
  font-weight: bold;
}

:deep(.el-table__body-wrapper tr:hover td) {
  background-color: var(--hover-bg) !important;
}

.preview-content {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}
</style>
