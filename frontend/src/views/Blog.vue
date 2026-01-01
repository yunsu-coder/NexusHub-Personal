<template>
  <div class="blog-page">
    <div class="blog-header">
      <div class="header-left">
        <h2>我的博客</h2>
        <span class="subtitle">记录生活，分享技术</span>
      </div>
      <el-tooltip content="创建新的博客文章" placement="bottom">
        <el-button type="primary" :icon="EditPen" @click="openEditor()">写文章</el-button>
      </el-tooltip>
    </div>

    <!-- 搜索和分类 -->
    <div class="blog-filters">
      <div class="filter-left">
        <el-input
          v-model="searchQuery"
          placeholder="搜索文章..."
          :prefix-icon="Search"
          clearable
          @input="onSearch"
          style="width: 300px"
        />
        <el-select
          v-model="selectedCategory"
          placeholder="选择分类"
          clearable
          @change="onFilter"
          style="width: 150px; margin-left: 10px"
        >
          <el-option label="全部" value="" />
          <el-option
            v-for="category in categories"
            :key="category"
            :label="category"
            :value="category"
          />
        </el-select>
      </div>
      <div class="filter-right">
        <el-button
          size="small"
          :type="sortOrder === 'desc' ? 'primary' : 'default'"
          @click="toggleSort"
        >
          {{ sortOrder === 'desc' ? '最新优先' : '最早优先' }}
        </el-button>
      </div>
    </div>

    <div class="blog-content">
      <div v-if="loading" class="loading-state">
        <el-skeleton :rows="3" animated />
      </div>
      
      <div v-else class="blog-grid">
        <div 
          v-for="post in filteredPosts" 
          :key="post.id" 
          class="blog-card"
          @click="viewPost(post)"
        >
          <div class="post-cover" v-if="post.cover || getRandomCover(post.id)">
            <img :src="post.cover || getRandomCover(post.id)" alt="cover" loading="lazy">
            <div class="post-overlay">
              <el-button circle :icon="View" />
            </div>
          </div>
          <div class="post-body">
            <div class="post-meta">
              <el-tag size="small" :type="getRandomType(post.id)">{{ post.category || '默认' }}</el-tag>
              <span class="time">{{ formatDate(post.created_at) }}</span>
            </div>
            <h3 class="post-title">{{ post.title }}</h3>
            <p class="post-excerpt">{{ getExcerpt(post.content) }}</p>
            <div class="post-footer">
              <div class="author">
                <el-avatar :size="24" :src="post.authorAvatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'" />
                <span>{{ post.author || '我' }}</span>
              </div>
              <div class="actions">
                <el-tooltip content="查看文章" placement="top">
                  <el-button circle :icon="View" />
                </el-tooltip>
                <el-tooltip content="编辑文章" placement="top">
                  <el-button text circle :icon="Edit" @click.stop="openEditor(post)" />
                </el-tooltip>
                <el-tooltip content="删除文章" placement="top">
                  <el-button text circle :icon="Delete" type="danger" @click.stop="deletePost(post)" />
                </el-tooltip>
              </div>
            </div>
          </div>
        </div>
        
        <el-empty v-if="!loading && posts.length === 0" description="暂无文章，开始写作吧！" />
      </div>
    </div>

    <!-- 编辑/阅读 对话框 (全屏) -->
    <el-dialog 
      v-model="showEditDialog" 
      :title="isEditing ? '编辑文章' : '阅读文章'" 
      fullscreen
      class="blog-dialog"
    >
      <div class="editor-layout" v-if="isEditing">
        <div class="editor-header">
          <el-input v-model="currentPost.title" placeholder="文章标题" class="title-input" />
          <el-input v-model="currentPost.category" placeholder="分类" style="width: 150px" />
          <el-button type="primary" @click="savePost">发布</el-button>
        </div>
        <div class="editor-meta">
          <div class="meta-row">
            <div class="meta-item">
              <span class="meta-label">标签：</span>
              <el-input v-model="currentPost.tags" placeholder="请输入标签，多个标签用逗号分隔" style="width: 300px" />
            </div>
            <div class="meta-item">
              <span class="meta-label">封面：</span>
              <el-input v-model="currentPost.cover" placeholder="封面图片URL" style="width: 400px" />
            </div>
          </div>
        </div>
        <div class="editor-main">
          <v-md-editor 
            v-model="currentPost.content" 
            height="100%"
            placeholder="开始写作..."
            :toolbar="customToolbar"
            preview-theme="auto"
          ></v-md-editor>
        </div>
      </div>
      
      <div class="reader-layout" v-else>
         <div class="reader-container">
           <h1 class="reader-title">{{ currentPost.title }}</h1>
           <div class="reader-meta">
             <span>{{ formatDate(currentPost.created_at) }}</span>
             <span>{{ currentPost.category }}</span>
             <span v-if="currentPost.tags">{{ currentPost.tags.split(',').join(' · ') }}</span>
           </div>
           <el-divider />
           <v-md-editor 
             :model-value="currentPost.content" 
             mode="preview"
             preview-theme="auto"
           ></v-md-editor>
         </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { EditPen, View, Edit, Delete, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

const showEditDialog = ref(false)
const isEditing = ref(false)
const currentPost = ref({ title: '', content: '', category: '', tags: '', cover: '' })
const posts = ref([])
const loading = ref(false)
const searchQuery = ref('')
const selectedCategory = ref('')
const sortOrder = ref('desc')
const allCategories = ref([])

// 自定义编辑器工具栏
const customToolbar = [
  'bold',
  'italic',
  'strike',
  '|',
  'title',
  'underline',
  'quote',
  '|',
  'list',
  'ordered-list',
  'check',
  '|',
  'code',
  'inline-code',
  'table',
  '|',
  'link',
  'image',
  '|',
  'undo',
  'redo',
  '|',
  'save',
  'preview',
  'fullscreen'
]

const fetchPosts = async () => {
  loading.value = true
  try {
    const data = await api.get('/blog')
    posts.value = data
    // 提取所有分类
    const categories = [...new Set(data.map(post => post.category || '默认'))]
    allCategories.value = categories
  } catch (e) {
    ElMessage.error('加载文章失败')
  } finally {
    loading.value = false
  }
}

// 计算属性：过滤和排序后的帖子
const filteredPosts = computed(() => {
  let result = [...posts.value]
  
  // 分类过滤
  if (selectedCategory.value) {
    result = result.filter(post => post.category === selectedCategory.value)
  }
  
  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(post => 
      post.title.toLowerCase().includes(query) || 
      post.content.toLowerCase().includes(query) ||
      (post.tags && post.tags.toLowerCase().includes(query))
    )
  }
  
  // 排序
  result.sort((a, b) => {
    const dateA = new Date(a.created_at).getTime()
    const dateB = new Date(b.created_at).getTime()
    return sortOrder.value === 'desc' ? dateB - dateA : dateA - dateB
  })
  
  return result
})

// 获取所有分类
const categories = computed(() => {
  return allCategories.value
})

// 搜索处理
const onSearch = () => {
  // 搜索自动触发过滤
}

// 分类过滤处理
const onFilter = () => {
  // 分类变化自动触发过滤
}

// 排序切换
const toggleSort = () => {
  sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc'
}

// 更新编辑器，添加标签和封面支持
const openEditor = (post = null) => {
  isEditing.value = true
  if (post) {
    currentPost.value = { ...post }
  } else {
    currentPost.value = { title: '', content: '', category: '', tags: '', cover: '' }
  }
  showEditDialog.value = true
}

// 更新保存方法，添加标签和封面支持
const savePost = async () => {
  if (!currentPost.value.title || !currentPost.value.content) return ElMessage.warning('标题和内容不能为空')
  
  try {
    if (currentPost.value.id) {
        await api.put(`/blog/${currentPost.value.id}`, currentPost.value)
    } else {
        await api.post('/blog', currentPost.value)
    }
    ElMessage.success('发布成功')
    showEditDialog.value = false
    fetchPosts()
  } catch (e) {
    ElMessage.error('发布失败')
  }
}

const viewPost = (post) => {
  isEditing.value = false
  currentPost.value = { ...post }
  showEditDialog.value = true
}

const deletePost = async (post) => {
  try {
    await ElMessageBox.confirm('确定删除这篇文章吗?', '提示', { type: 'warning' })
    await api.delete(`/blog/${post.id}`)
    ElMessage.success('已删除')
    fetchPosts()
  } catch (e) {}
}

const formatDate = (date) => new Date(date).toLocaleDateString()

const getExcerpt = (content) => {
    if (!content) return ''
    return content.replace(/[#*`]/g, '').substring(0, 100) + '...'
}

const getRandomCover = (id) => {
    const images = [
        'https://images.unsplash.com/photo-1499750310159-529801977349?w=500&auto=format&fit=crop&q=60',
        'https://images.unsplash.com/photo-1486312338219-ce68d2c6f44d?w=500&auto=format&fit=crop&q=60',
        'https://images.unsplash.com/photo-1519389950473-47ba0277781c?w=500&auto=format&fit=crop&q=60',
        'https://images.unsplash.com/photo-1455390582262-044cdead277a?w=500&auto=format&fit=crop&q=60'
    ]
    return images[(id || 0) % images.length]
}

const getRandomType = (id) => {
    const types = ['', 'success', 'warning', 'danger', 'info']
    return types[(id || 0) % types.length]
}

onMounted(fetchPosts)
</script>

<style scoped>
.blog-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 24px;
  gap: 24px;
  color: var(--text-primary);
  background: var(--bg-color);
}

.blog-header {
  padding: 24px 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
}

.blog-header:hover {
  box-shadow: var(--shadow-md);
}

.header-left h2 { 
  margin: 0; 
  font-size: 28px; 
  color: var(--text-primary);
}

.subtitle { 
  opacity: 0.7; 
  font-size: 14px; 
  margin-top: 5px; 
  display: block; 
  color: var(--text-secondary); 
}

.blog-filters {
  padding: 16px 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
  margin-bottom: 24px;
  color: var(--text-primary);
}

.blog-filters:hover {
  box-shadow: var(--shadow-md);
}

.filter-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-right {
  display: flex;
  gap: 10px;
}

.blog-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 24px;
}

/* Editor Meta Section */
.editor-meta {
  padding: 16px 32px;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  gap: 32px;
  align-items: center;
  flex-wrap: wrap;
  color: var(--text-primary);
}

.meta-row {
  display: flex;
  gap: 32px;
  flex-wrap: wrap;
  width: 100%;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.meta-label {
  font-weight: 500;
  color: var(--text-secondary);
  white-space: nowrap;
  font-size: 14px;
}

.blog-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 28px;
}

.blog-card {
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  height: 400px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  box-shadow: var(--shadow-sm);
}

.blog-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
  border-color: var(--primary-color);
}

.post-cover {
  height: 180px;
  position: relative;
  overflow: hidden;
  border-bottom: 1px solid var(--border-color);
}

.post-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.blog-card:hover .post-cover img {
  transform: scale(1.08);
}

.post-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.blog-card:hover .post-overlay {
  opacity: 1;
}

.post-body {
  padding: 24px;
  flex: 1;
  display: flex;
  flex-direction: column;
  color: var(--text-primary);
}

.post-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.time { 
  font-size: 13px; 
  color: var(--text-secondary);
}

.post-title {
  font-size: 19px;
  font-weight: 700;
  margin: 0 0 12px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  color: var(--text-primary);
  transition: color 0.3s ease;
}

.blog-card:hover .post-title {
  color: var(--primary-color);
}

.post-excerpt {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0 0 24px 0;
  flex: 1;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid var(--border-color);
  padding-top: 16px;
}

.author { 
  display: flex; 
  align-items: center; 
  gap: 10px; 
  font-size: 14px; 
  color: var(--text-secondary);
}

.actions {
  display: flex;
  gap: 8px;
}

/* Dialog Styles */
.editor-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-color);
}

.editor-header {
  padding: 16px 32px;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  gap: 20px;
  align-items: center;
  color: var(--text-primary);
}

.title-input {
  flex: 1; 
  font-size: 20px; 
  font-weight: bold; 
  color: var(--text-primary);
}

.editor-main {
  flex: 1; 
  overflow: hidden;
}

/* 自定义 markdown 编辑器样式 */
:deep(.v-md-editor) {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
  background: var(--card-bg);
}

:deep(.v-md-editor-preview) {
  background: var(--card-bg);
  color: var(--text-primary);
}

:deep(.v-md-editor-toolbar) {
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
}

.reader-layout {
  max-width: 800px;
  margin: 0 auto;
  padding: 48px 24px;
  background: var(--card-bg);
  min-height: 100vh;
}

.reader-title {
  font-size: 36px; 
  font-weight: 800; 
  color: var(--text-primary); 
  margin-bottom: 24px; 
  line-height: 1.2;
}

.reader-meta {
  color: var(--text-secondary); 
  margin-bottom: 32px; 
  display: flex; 
  gap: 24px; 
  font-size: 14px;
}

.reader-container {
  background: var(--card-bg);
  padding: 0;
  border-radius: 12px;
}
</style>
