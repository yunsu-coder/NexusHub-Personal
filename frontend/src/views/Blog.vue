<template>
  <div class="blog-page">
    <div class="blog-header glass-panel">
      <div class="header-left">
        <h2>我的博客</h2>
        <span class="subtitle">记录生活，分享技术</span>
      </div>
      <el-button type="primary" :icon="EditPen" @click="openEditor()">写文章</el-button>
    </div>

    <div class="blog-content">
      <div v-if="loading" class="loading-state">
        <el-skeleton :rows="3" animated />
      </div>
      
      <div v-else class="blog-grid">
        <div 
          v-for="post in posts" 
          :key="post.id" 
          class="blog-card glass-panel"
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
              <el-tag size="small" effect="dark" :type="getRandomType(post.id)">{{ post.category || '默认' }}</el-tag>
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
                <el-button text circle :icon="Edit" @click.stop="openEditor(post)" />
                <el-button text circle :icon="Delete" type="danger" @click.stop="deletePost(post)" />
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
          <el-input v-model="currentPost.category" placeholder="分类" style="width: 200px" />
          <el-button type="primary" @click="savePost">发布</el-button>
        </div>
        <div class="editor-main">
          <v-md-editor 
            v-model="currentPost.content" 
            height="100%"
            placeholder="开始写作..."
          ></v-md-editor>
        </div>
      </div>
      
      <div class="reader-layout" v-else>
         <div class="reader-container">
           <h1 class="reader-title">{{ currentPost.title }}</h1>
           <div class="reader-meta">
             <span>{{ formatDate(currentPost.created_at) }}</span>
             <span>{{ currentPost.category }}</span>
           </div>
           <el-divider />
           <v-md-editor 
             :model-value="currentPost.content" 
             mode="preview"
           ></v-md-editor>
         </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { EditPen, View, Edit, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

const showEditDialog = ref(false)
const isEditing = ref(false)
const currentPost = ref({ title: '', content: '', category: '' })
const posts = ref([])
const loading = ref(false)

const fetchPosts = async () => {
  loading.value = true
  try {
    const data = await api.get('/blog')
    posts.value = data
  } catch (e) {
    ElMessage.error('加载文章失败')
  } finally {
    loading.value = false
  }
}

const openEditor = (post = null) => {
  isEditing.value = true
  if (post) {
    currentPost.value = { ...post }
  } else {
    currentPost.value = { title: '', content: '', category: '' }
  }
  showEditDialog.value = true
}

const viewPost = (post) => {
  isEditing.value = false
  currentPost.value = { ...post }
  showEditDialog.value = true
}

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
  padding: 20px;
  gap: 20px;
}

.glass-panel {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
}

.blog-header {
  padding: 20px 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #fff;
}
.header-left h2 { margin: 0; font-size: 24px; }
.subtitle { opacity: 0.7; font-size: 14px; margin-top: 5px; display: block; }

.blog-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 20px;
}

.blog-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.blog-card {
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  display: flex;
  flex-direction: column;
  height: 380px;
}
.blog-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  background: rgba(255,255,255,0.1);
}

.post-cover {
  height: 180px;
  position: relative;
  overflow: hidden;
}
.post-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}
.blog-card:hover .post-cover img {
  transform: scale(1.05);
}

.post-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
}
.blog-card:hover .post-overlay { opacity: 1; }

.post-body {
  padding: 20px;
  flex: 1;
  display: flex;
  flex-direction: column;
  color: #eee;
}

.post-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}
.time { font-size: 12px; opacity: 0.6; }

.post-title {
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 10px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-excerpt {
  font-size: 13px;
  opacity: 0.7;
  line-height: 1.6;
  margin: 0 0 20px 0;
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
  border-top: 1px solid rgba(255,255,255,0.1);
  padding-top: 15px;
}
.author { display: flex; align-items: center; gap: 8px; font-size: 13px; opacity: 0.9; }

/* Dialog Styles */
.editor-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f9f9f9;
}
.editor-header {
  padding: 15px 30px;
  background: #fff;
  border-bottom: 1px solid #eee;
  display: flex;
  gap: 20px;
  align-items: center;
}
.title-input { flex: 1; font-size: 18px; font-weight: bold; }
.editor-main { flex: 1; overflow: hidden; }

.reader-layout {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
}
.reader-title { font-size: 32px; font-weight: 800; color: #333; margin-bottom: 20px; }
.reader-meta { color: #999; margin-bottom: 30px; display: flex; gap: 20px; }
</style>
