<template>
  <div class="collection-page">
    <div class="page-header">
      <h2>我的收藏</h2>
      <el-button type="primary" @click="showAddDialog = true">
        <el-icon><Plus /></el-icon>
        添加收藏
      </el-button>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="collection-filters">
      <div class="search-box">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索收藏..."
          clearable
          prefix-icon="el-icon-Search"
          class="search-input"
        >
          <template #append>
             <el-button @click="showUrlDialog = true">
                <el-icon><Link /></el-icon> 爬取
             </el-button>
          </template>
        </el-input>
      </div>
      <div class="filter-actions">
        <el-select
          v-model="selectedCategory"
          placeholder="选择分类"
          clearable
          class="category-select"
        >
          <el-option label="所有分类" value="" />
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>
        <el-dropdown @command="handleBulkAction">
          <el-button type="default">
            批量操作 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="batchDelete" :disabled="selectedItems.length === 0">
                <el-icon><Delete /></el-icon> 批量删除
              </el-dropdown-item>
              <el-dropdown-item command="batchMove" :disabled="selectedItems.length === 0">
                <el-icon><SwitchButton /></el-icon> 批量移动
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button type="default" @click="showCategoryDialog = true">
          <el-icon><Plus /></el-icon> 管理分类
        </el-button>
      </div>
    </div>

    <!-- 收藏类型切换 -->
    <div class="collection-tabs">
      <el-tabs v-model="activeTab" @tab-click="handleTabClick">
        <el-tab-pane label="全部" name="all">
          <el-icon><Collection /></el-icon>
        </el-tab-pane>
        <el-tab-pane label="图片" name="image">
          <el-icon><PictureFilled /></el-icon>
        </el-tab-pane>
        <el-tab-pane label="视频" name="video">
          <el-icon><VideoPlay /></el-icon>
        </el-tab-pane>
        <el-tab-pane label="文章" name="article">
          <el-icon><Document /></el-icon>
        </el-tab-pane>
        <el-tab-pane label="链接" name="link">
          <el-icon><Link /></el-icon>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 收藏列表 -->
    <div class="collection-grid">
      <div 
        v-for="item in filteredCollections" 
        :key="item.id" 
        class="collection-item"
        :class="`type-${item.type} ${selectedItems.includes(item.id) ? 'selected' : ''}`"
        @click="viewItem(item)"
      >
        <div class="item-checkbox">
          <el-checkbox v-model="selectedItems" :label="item.id" @click.stop />
        </div>
        <div class="item-thumbnail">
          <el-image 
            v-if="item.type === 'image'" 
            :src="item.thumbnail || item.url" 
            :fit="'cover'"
            :preview-src-list="[item.url]"
          >
            <template #error>
              <div class="image-error">
                <el-icon><PictureFilled /></el-icon>
              </div>
            </template>
          </el-image>
          <div v-else-if="item.type === 'video'" class="video-thumbnail">
            <el-icon class="play-icon"><VideoPlay /></el-icon>
            <el-image 
              :src="item.thumbnail" 
              :fit="'cover'"
            >
              <template #error>
                <div class="video-error">
                  <el-icon><VideoPlay /></el-icon>
                </div>
              </template>
            </el-image>
          </div>
          <div v-else-if="item.type === 'article'" class="article-thumbnail">
            <el-icon class="article-icon"><Document /></el-icon>
            <div class="article-info">
              <div class="article-domain">{{ getDomain(item.url) }}</div>
            </div>
          </div>
          <div v-else class="link-thumbnail">
            <el-icon class="link-icon"><Link /></el-icon>
            <div class="link-info">
              <div class="link-domain">{{ getDomain(item.url) }}</div>
            </div>
          </div>
        </div>
        <div class="item-info">
          <h3 class="item-title">{{ item.title }}</h3>
          <p class="item-description">{{ item.description || '无描述' }}</p>
          <div class="item-meta">
            <span class="item-type-badge">{{ getItemTypeLabel(item.type) }}</span>
            <span class="item-date">{{ formatDate(item.created_at) }}</span>
          </div>
        </div>
        <div class="item-actions">
          <el-button 
            type="text" 
            @click.stop="editItem(item)"
            class="action-btn"
          >
            <el-icon><EditPen /></el-icon>
          </el-button>
          <el-button 
            type="text" 
            @click.stop="deleteItem(item)"
            class="action-btn delete-btn"
          >
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
      </div>
      <div v-if="filteredCollections.length === 0" class="empty-collection">
        <el-empty description="暂无收藏" />
      </div>
    </div>

    <!-- 添加/编辑收藏对话框 -->
    <el-dialog 
      v-model="showAddDialog" 
      :title="currentItem.id ? '编辑收藏' : '添加收藏'" 
      width="500px"
      :before-close="handleCloseDialog"
    >
      <el-form 
        ref="collectionForm" 
        :model="currentItem" 
        label-width="80px"
        :rules="formRules"
      >
        <el-form-item label="类型" prop="type">
          <el-select v-model="currentItem.type" placeholder="请选择类型">
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
            <el-option label="文章" value="article" />
            <el-option label="链接" value="link" />
          </el-select>
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="currentItem.category_id" placeholder="选择分类" clearable>
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="currentItem.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="URL" prop="url">
          <el-input v-model="currentItem.url" placeholder="请输入URL" />
        </el-form-item>
        <el-form-item label="缩略图" prop="thumbnail">
          <el-input v-model="currentItem.thumbnail" placeholder="请输入缩略图URL（可选）" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="currentItem.description" 
            type="textarea" 
            rows="3" 
            placeholder="请输入描述（可选）" 
          />
        </el-form-item>
        <el-form-item label="标签">
          <el-tag 
            v-for="tag in currentItem.tags" 
            :key="tag" 
            closable 
            @close="removeTag(tag)"
          >
            {{ tag }}
          </el-tag>
          <el-input 
            v-model="newTag" 
            class="tag-input" 
            placeholder="按Enter添加标签"
            @keyup.enter="addTag"
          />
        </el-form-item>

        <div v-if="currentItem.resources && currentItem.resources.length > 0" class="detected-resources">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
              <h4>检测到的资源 (点击选择)</h4>
              <div v-if="imageResources.length > 0">
                <el-checkbox v-model="selectAllImages" @change="handleSelectAllImages">全选图片</el-checkbox>
                <el-button type="primary" size="small" @click="batchAddImages" :disabled="selectedResourceIndices.length === 0">
                  批量添加 ({{ selectedResourceIndices.length }})
                </el-button>
              </div>
            </div>
            
            <!-- 图片资源网格 -->
            <div v-if="imageResources.length > 0" class="resource-section">
              <div class="resource-grid">
                <div 
                  v-for="(res, idx) in imageResources" 
                  :key="idx" 
                  class="resource-grid-item"
                  :class="{ selected: selectedResourceIndices.includes(idx) }"
                  @click="toggleResourceSelection(idx)"
                >
                  <el-image 
                    :src="res.url" 
                    fit="cover" 
                    loading="lazy"
                    draggable="false"
                    style="pointer-events: none;"
                  >
                    <template #error>
                      <div class="image-slot">
                        <el-icon><PictureFilled /></el-icon>
                      </div>
                    </template>
                  </el-image>
                  <div class="resource-overlay">
                    <el-icon v-if="selectedResourceIndices.includes(idx)"><Check /></el-icon>
                    <el-icon v-else><Plus /></el-icon>
                  </div>
                </div>
              </div>
            </div>

          <!-- 其他资源列表 -->
          <div v-if="otherResources.length > 0" class="resource-list">
            <div 
              v-for="(res, idx) in otherResources" 
              :key="idx" 
              class="resource-item"
              @click="selectResource(res)"
            >
              <el-tag size="small" :type="getResourceTypeTag(res.type)">{{ res.type }}</el-tag>
              <span class="resource-url" :title="res.url">{{ res.url }}</span>
            </div>
          </div>
        </div>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleCloseDialog">取消</el-button>
          <el-button type="primary" @click="saveItem">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 分类管理对话框 -->
    <el-dialog 
      v-model="showCategoryDialog" 
      title="分类管理" 
      width="400px"
    >
      <div class="category-list">
        <div 
          v-for="category in categories" 
          :key="category.id"
          class="category-item"
        >
          <div class="category-info">
            <el-input
              v-model="category.name"
              size="small"
              placeholder="分类名称"
              @change="updateCategory(category)"
              class="category-name-input"
            />
          </div>
          <div class="category-actions">
            <el-button
              type="text"
              size="small"
              @click="deleteCategory(category)"
              class="delete-category-btn"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
        <div class="add-category">
          <el-input
            v-model="newCategoryName"
            size="small"
            placeholder="新分类名称"
            suffix-icon="el-icon-Plus"
            @keyup.enter="addCategory"
            @click.suffix="addCategory"
          />
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCategoryDialog = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 批量移动对话框 -->
    <el-dialog 
      v-model="showBatchMoveDialog" 
      title="批量移动收藏" 
      width="400px"
    >
      <el-form :model="batchMoveForm">
        <el-form-item label="目标分类">
          <el-select
            v-model="batchMoveForm.targetCategoryId"
            placeholder="选择目标分类"
            style="width: 100%"
          >
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showBatchMoveDialog = false">取消</el-button>
          <el-button type="primary" @click="executeBatchMove">确定移动</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- URL 爬取对话框 -->
    <el-dialog
      v-model="showUrlDialog"
      title="从 URL 添加"
      width="400px"
    >
      <el-form :model="urlForm">
        <el-form-item label="网址 URL">
          <el-input v-model="urlForm.url" placeholder="https://example.com" @keyup.enter="parseUrl" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showUrlDialog = false">取消</el-button>
          <el-button type="primary" :loading="urlLoading" @click="parseUrl">获取信息</el-button>
        </span>
      </template>
    </el-dialog>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { 
  Plus, Collection, PictureFilled, VideoPlay, Document, Link, 
  EditPen, Delete, ArrowDown, SwitchButton, Search, Check
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import api from '../api'

const router = useRouter()

// 收藏数据
const collections = ref([])
const activeTab = ref('all')
const showAddDialog = ref(false)
const showCategoryDialog = ref(false)
const showBatchMoveDialog = ref(false)
const showUrlDialog = ref(false)
const urlLoading = ref(false)
const urlForm = ref({ url: '' })

const newTag = ref('')
const newCategoryName = ref('')
const searchKeyword = ref('')
const selectedCategory = ref('')
const selectedItems = ref([])
const selectAllImages = ref(false)
const selectedResourceIndices = ref([])

// 分类数据
const categories = ref([
  { id: 1, name: '默认分类' },
  { id: 2, name: '工作资料' },
  { id: 3, name: '学习资源' }
])

// 表单数据
const currentItem = ref({
  id: null,
  title: '',
  url: '',
  type: 'image',
  thumbnail: '',
  description: '',
  tags: [],
  category_id: '',
  resources: []
})

// 批量移动表单
const batchMoveForm = ref({
  targetCategoryId: ''
})

// 表单验证规则
const formRules = {
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  url: [{ required: true, message: '请输入URL', trigger: 'blur' }]
}

const imageResources = computed(() => {
  return currentItem.value.resources?.filter(r => r.type === 'image' || r.type === 'img') || []
})

const otherResources = computed(() => {
  return currentItem.value.resources?.filter(r => r.type !== 'image' && r.type !== 'img') || []
})

// 过滤后的收藏列表
const filteredCollections = computed(() => {
  let filtered = collections.value
  
  // 按类型过滤
  if (activeTab.value !== 'all') {
    filtered = filtered.filter(item => item.type === activeTab.value)
  }
  
  // 按分类过滤
  if (selectedCategory.value) {
    filtered = filtered.filter(item => item.category_id === selectedCategory.value)
  }
  
  // 按关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(item => 
      item.title.toLowerCase().includes(keyword) ||
      item.description.toLowerCase().includes(keyword) ||
      item.tags.some(tag => tag.toLowerCase().includes(keyword)) ||
      item.url.toLowerCase().includes(keyword)
    )
  }
  
  return filtered
})

// 获取所有收藏
async function fetchCollections() {
  try {
    const data = await api.get('/collections')
    // 转换tags字符串为数组
    collections.value = data.map(item => ({
      ...item,
      tags: item.tags ? item.tags.split(',').filter(tag => tag) : [],
      category_id: item.category_id || ''
    }))
  } catch (error) {
    console.error('获取收藏失败:', error)
    ElMessage.error('获取收藏失败')
    // 如果API不可用，使用本地模拟数据
    collections.value = getMockCollections()
  }
}

// 获取所有分类
async function fetchCategories() {
  try {
    const data = await api.get('/categories')
    categories.value = data
  } catch (error) {
    console.error('获取分类失败:', error)
    // 使用本地存储的分类或默认分类
    const savedCategories = localStorage.getItem('collection_categories')
    if (savedCategories) {
      categories.value = JSON.parse(savedCategories)
    } else {
      // 使用默认分类
      categories.value = [
        { id: 1, name: '默认分类' },
        { id: 2, name: '工作资料' },
        { id: 3, name: '学习资源' }
      ]
      saveCategoriesToLocalStorage()
    }
  }
}

// 爬取 URL
async function parseUrl() {
  if (!urlForm.value.url) return
  urlLoading.value = true
  try {
    const res = await api.post('/collections/parse', { url: urlForm.value.url })
    currentItem.value = {
      ...currentItem.value,
      title: res.title || '无标题',
      description: res.description || '',
      url: urlForm.value.url,
      thumbnail: res.image || '',
      type: 'link', // 默认为链接
      resources: res.resources || []
    }
    showUrlDialog.value = false
    showAddDialog.value = true
    urlForm.value.url = ''
    ElMessage.success('解析成功')
  } catch (e) {
    ElMessage.error('解析失败: ' + (e.response?.data?.error || e.message))
  } finally {
    urlLoading.value = false
  }
}

// 添加/更新收藏
async function saveItem() {
  try {
    // 转换tags数组为逗号分隔字符串
    const collectionData = {
      ...currentItem.value,
      tags: currentItem.value.tags.join(',')
    }
    
    if (currentItem.value.id) {
      // 更新收藏
      const data = await api.put(`/collections/${currentItem.value.id}`, collectionData)
      // 转换tags字符串为数组
      data.tags = data.tags.split(',').filter(tag => tag)
      const index = collections.value.findIndex(item => item.id === currentItem.value.id)
      if (index !== -1) {
        collections.value[index] = data
      }
      ElMessage.success('更新收藏成功')
    } else {
      // 添加收藏
      const data = await api.post('/collections', collectionData)
      // 转换tags字符串为数组
      data.tags = data.tags.split(',').filter(tag => tag)
      collections.value.unshift(data)
      ElMessage.success('添加收藏成功')
    }
    showAddDialog.value = false
    resetForm()
  } catch (error) {
    console.error('保存收藏失败:', error)
    ElMessage.error('保存收藏失败')
    // 如果API不可用，使用本地模拟数据
    handleMockSave()
  }
}

// 添加分类
function addCategory() {
  if (!newCategoryName.value.trim()) return
  
  const newCategory = {
    id: Date.now(),
    name: newCategoryName.value.trim()
  }
  
  categories.value.push(newCategory)
  newCategoryName.value = ''
  saveCategoriesToLocalStorage()
  ElMessage.success('分类添加成功')
}

// 更新分类
function updateCategory(category) {
  if (!category.name.trim()) {
    ElMessage.error('分类名称不能为空')
    return
  }
  saveCategoriesToLocalStorage()
  ElMessage.success('分类更新成功')
}

// 删除分类
function deleteCategory(category) {
  // 检查是否有收藏使用该分类
  const hasItems = collections.value.some(item => item.category_id === category.id)
  if (hasItems) {
    ElMessage.error('该分类下有收藏，无法删除')
    return
  }
  
  const index = categories.value.findIndex(c => c.id === category.id)
  if (index !== -1) {
    categories.value.splice(index, 1)
    saveCategoriesToLocalStorage()
    ElMessage.success('分类删除成功')
  }
}

// 批量操作处理
function handleBulkAction(command) {
  switch (command) {
    case 'batchDelete':
      batchDeleteItems()
      break
    case 'batchMove':
      showBatchMoveDialog.value = true
      break
    default:
      break
  }
}

// 批量删除
async function batchDeleteItems() {
  if (selectedItems.value.length === 0) return
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedItems.value.length} 个收藏吗？`, 
      '确认删除', 
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 批量删除
    for (const id of selectedItems.value) {
      await api.delete(`/collections/${id}`)
    }
    
    // 更新本地数据
    collections.value = collections.value.filter(item => !selectedItems.value.includes(item.id))
    selectedItems.value = []
    ElMessage.success('批量删除成功')
  } catch (error) {
    if (error === 'cancel') return
    console.error('批量删除失败:', error)
    ElMessage.error('批量删除失败')
    
    // 模拟批量删除
    collections.value = collections.value.filter(item => !selectedItems.value.includes(item.id))
    selectedItems.value = []
    ElMessage.success('批量删除成功（模拟）')
  }
}

// 执行批量移动
async function executeBatchMove() {
  if (!batchMoveForm.value.targetCategoryId) {
    ElMessage.error('请选择目标分类')
    return
  }
  
  try {
    // 批量移动
    for (const id of selectedItems.value) {
      const item = collections.value.find(item => item.id === id)
      if (item) {
        await api.put(`/collections/${id}`, {
          ...item,
          category_id: batchMoveForm.value.targetCategoryId,
          tags: item.tags.join(',')
        })
      }
    }
    
    // 更新本地数据
    collections.value.forEach(item => {
      if (selectedItems.value.includes(item.id)) {
        item.category_id = batchMoveForm.value.targetCategoryId
      }
    })
    
    selectedItems.value = []
    showBatchMoveDialog.value = false
    batchMoveForm.value.targetCategoryId = ''
    ElMessage.success('批量移动成功')
  } catch (error) {
    console.error('批量移动失败:', error)
    ElMessage.error('批量移动失败')
    
    // 模拟批量移动
    collections.value.forEach(item => {
      if (selectedItems.value.includes(item.id)) {
        item.category_id = batchMoveForm.value.targetCategoryId
      }
    })
    
    selectedItems.value = []
    showBatchMoveDialog.value = false
    batchMoveForm.value.targetCategoryId = ''
    ElMessage.success('批量移动成功（模拟）')
  }
}

// 删除收藏
async function deleteItem(item) {
  try {
    await ElMessageBox.confirm('确定要删除这个收藏吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await api.delete(`/collections/${item.id}`)
    const index = collections.value.findIndex(i => i.id === item.id)
    if (index !== -1) {
      collections.value.splice(index, 1)
    }
    ElMessage.success('删除收藏成功')
  } catch (error) {
    if (error === 'cancel') return
    console.error('删除收藏失败:', error)
    ElMessage.error('删除收藏失败')
    // 如果API不可用，使用本地模拟删除
    const index = collections.value.findIndex(i => i.id === item.id)
    if (index !== -1) {
      collections.value.splice(index, 1)
    }
    ElMessage.success('删除收藏成功')
  }
}

// 编辑收藏
function editItem(item) {
  // 确保tags是数组格式
  currentItem.value = {
    ...item,
    tags: Array.isArray(item.tags) ? item.tags : (item.tags ? item.tags.split(',').filter(tag => tag) : []),
    category_id: item.category_id || ''
  }
  showAddDialog.value = true
}

// 查看收藏详情
function viewItem(item) {
  if (item.type === 'video') {
    // 打开视频播放
    window.open(item.url, '_blank')
  } else {
    // 打开链接
    window.open(item.url, '_blank')
  }
}

// 添加标签
function addTag() {
  if (newTag.value && !currentItem.value.tags.includes(newTag.value)) {
    currentItem.value.tags.push(newTag.value)
    newTag.value = ''
  }
}

// 移除标签
function removeTag(tag) {
  const index = currentItem.value.tags.indexOf(tag)
  if (index !== -1) {
    currentItem.value.tags.splice(index, 1)
  }
}

// 关闭对话框
function handleCloseDialog() {
  showAddDialog.value = false
  resetForm()
}

// 重置表单
function resetForm() {
  currentItem.value = {
    id: null,
    title: '',
    url: '',
    type: 'image',
    thumbnail: '',
    description: '',
    tags: []
  }
  newTag.value = ''
}

// 标签切换
function handleTabClick(tab) {
  activeTab.value = tab.props.name
}

function selectResource(res) {
  // 单选逻辑保持不变，用于“其他资源列表”
  currentItem.value.url = res.url
  currentItem.value.type = res.type === 'img' ? 'image' : (res.type === 'audio' ? 'video' : res.type)
  
  if (res.type === 'image' || res.type === 'img') {
     currentItem.value.thumbnail = res.url
     currentItem.value.type = 'image'
  } else if (res.type === 'video') {
     currentItem.value.type = 'video'
  } else if (res.type === 'audio') {
     currentItem.value.type = 'video' 
  }
  ElMessage.success('已选择资源')
}

// 切换资源选中状态（用于批量）
function toggleResourceSelection(idx) {
  const index = selectedResourceIndices.value.indexOf(idx)
  if (index === -1) {
    selectedResourceIndices.value.push(idx)
  } else {
    selectedResourceIndices.value.splice(index, 1)
  }
  // 如果单选，也更新主表单
  if (selectedResourceIndices.value.length === 1) {
    selectResource(imageResources.value[selectedResourceIndices.value[0]])
  }
  selectAllImages.value = selectedResourceIndices.value.length === imageResources.value.length
}

function handleSelectAllImages(val) {
  if (val) {
    selectedResourceIndices.value = imageResources.value.map((_, idx) => idx)
  } else {
    selectedResourceIndices.value = []
  }
}

async function batchAddImages() {
  if (selectedResourceIndices.value.length === 0) return
  
  const selectedRes = selectedResourceIndices.value.map(idx => imageResources.value[idx])
  
  let successCount = 0
  for (const res of selectedRes) {
    // 关键修复：深拷贝 currentItem 并清空 id，确保是新记录
    const newItem = JSON.parse(JSON.stringify(currentItem.value))
    newItem.title = newItem.title.replace(/\s\(\d+\)$/, '') + ` (${successCount + 1})`
    newItem.url = res.url
    newItem.thumbnail = res.url
    newItem.type = 'image'
    newItem.id = 0 // 明确置为 0 或 null

    try {
      // 构造请求数据
      const collectionData = { ...newItem, tags: Array.isArray(newItem.tags) ? newItem.tags.join(',') : newItem.tags }
      const data = await api.post('/collections', collectionData)
      
      // 修正返回数据的 tags 格式
      if (typeof data.tags === 'string') {
          data.tags = data.tags.split(',').filter(tag => tag)
      }
      
      collections.value.unshift(data)
      successCount++
    } catch (e) {
      console.error('批量添加失败', e)
    }
  }
  
  ElMessage.success(`成功添加 ${successCount} 张图片`)
  showAddDialog.value = false
  resetForm()
}

function getResourceTypeTag(type) {
  const map = {
    image: 'success',
    img: 'success',
    video: 'danger',
    audio: 'warning'
  }
  return map[type] || 'info'
}

// 工具函数
function getItemTypeLabel(type) {
  const labels = {
    image: '图片',
    video: '视频',
    article: '文章',
    link: '链接'
  }
  return labels[type] || '未知'
}

function formatDate(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

function getDomain(url) {
  try {
    const urlObj = new URL(url)
    return urlObj.hostname.replace('www.', '')
  } catch (error) {
    return '未知'
  }
}

// 保存分类到本地存储
function saveCategoriesToLocalStorage() {
  localStorage.setItem('collection_categories', JSON.stringify(categories.value))
}

// 从本地存储加载分类
function loadCategoriesFromLocalStorage() {
  const saved = localStorage.getItem('collection_categories')
  if (saved) {
    categories.value = JSON.parse(saved)
  }
}

// 模拟数据
function getMockCollections() {
  return [
    {
      id: 1,
      title: '美丽的风景图片',
      url: 'https://picsum.photos/id/1015/800/600',
      type: 'image',
      thumbnail: 'https://picsum.photos/id/1015/200/150',
      description: '这是一张美丽的风景图片，展示了山川河流的壮丽景色',
      tags: ['风景', '自然', '摄影'],
      category_id: 1,
      created_at: '2024-01-15T10:30:00Z'
    },
    {
      id: 2,
      title: '教程视频：Vue3基础入门',
      url: 'https://example.com/videos/vue3-tutorial.mp4',
      type: 'video',
      thumbnail: 'https://picsum.photos/id/1025/200/150',
      description: 'Vue3基础入门教程，适合初学者学习',
      tags: ['Vue3', '教程', '前端'],
      category_id: 3,
      created_at: '2024-01-14T15:20:00Z'
    },
    {
      id: 3,
      title: '现代JavaScript开发最佳实践',
      url: 'https://example.com/articles/javascript-best-practices',
      type: 'article',
      description: '探讨现代JavaScript开发的最佳实践和技巧',
      tags: ['JavaScript', '前端', '最佳实践'],
      category_id: 2,
      created_at: '2024-01-13T09:15:00Z'
    },
    {
      id: 4,
      title: 'GitHub - 世界最大的代码托管平台',
      url: 'https://github.com',
      type: 'link',
      description: 'GitHub是全球最大的代码托管平台，用于版本控制和协作',
      tags: ['GitHub', '代码', '开发'],
      category_id: 2,
      created_at: '2024-01-12T14:45:00Z'
    }
  ]
}

// 模拟保存
function handleMockSave() {
  if (currentItem.value.id) {
    // 模拟更新
    const index = collections.value.findIndex(item => item.id === currentItem.value.id)
    if (index !== -1) {
      collections.value[index] = { ...currentItem.value }
    }
  } else {
    // 模拟添加
    const newItem = {
      ...currentItem.value,
      id: Date.now(),
      created_at: new Date().toISOString()
    }
    collections.value.unshift(newItem)
  }
  showAddDialog.value = false
  resetForm()
  ElMessage.success('保存成功（模拟数据）')
}

// 页面加载时获取数据
onMounted(() => {
  loadCategoriesFromLocalStorage()
  fetchCollections()
})

// 监听分类变化，保存到本地存储
watch(categories, () => {
  saveCategoriesToLocalStorage()
}, { deep: true })
</script>

<style scoped>
.collection-page {
  padding: 24px;
  animation: fadeIn 0.5s ease;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-color);
}

.page-header h2 {
  color: var(--text-primary);
  font-size: 28px;
  font-weight: 700;
  margin: 0;
}

.collection-tabs {
  margin-bottom: 24px;
  background: var(--card-bg);
  border-radius: 12px;
  padding: 8px;
  box-shadow: var(--shadow-sm);
}

.collection-tabs :deep(.el-tabs__header) {
  margin: 0;
}

.collection-tabs :deep(.el-tabs__nav-wrap) {
  padding: 0;
}

.collection-tabs :deep(.el-tabs__nav) {
  display: flex;
  width: 100%;
  justify-content: space-around;
}

.collection-tabs :deep(.el-tabs__item) {
  flex: 1;
  text-align: center;
  color: var(--text-secondary);
  font-weight: 500;
  padding: 12px 0;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.collection-tabs :deep(.el-tabs__item.is-active) {
  color: var(--primary-color);
  background: var(--card-bg);
  box-shadow: var(--shadow-sm);
}

.collection-tabs :deep(.el-tabs__active-bar) {
  display: none;
}

.collection-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.collection-item {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-sm);
}

.collection-item:hover {
  transform: translateY(-8px);
  box-shadow: var(--shadow-lg);
  border-color: var(--primary-color);
}

.item-thumbnail {
  position: relative;
  width: 100%;
  height: 200px;
  overflow: hidden;
  background: var(--bg-color);
}

.item-thumbnail img {
  width: 100%;
  height: 100%;
  transition: transform 0.3s ease;
}

.collection-item:hover .item-thumbnail :deep(img),
.collection-item:hover .item-thumbnail :deep(.el-image__inner) {
  transform: scale(1.03);
}

/* 关闭 Element Plus Image 的悬停遮罩以避免闪屏 */
.item-thumbnail :deep(.el-image__mask) {
  opacity: 0 !important;
  pointer-events: none !important;
}

/* 保持 Image 内部容器稳定，避免尺寸抖动 */
.item-thumbnail :deep(.el-image__wrapper) {
  width: 100% !important;
  height: 100% !important;
}

.video-thumbnail {
  position: relative;
  width: 100%;
  height: 100%;
}

.play-icon {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 48px;
  color: white;
  opacity: 0.9;
  z-index: 1;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  pointer-events: none;
}

.video-thumbnail:hover .play-icon {
  transform: translate(-50%, -50%) scale(1.1);
  opacity: 1;
}

.article-thumbnail, .link-thumbnail {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
}

.article-icon, .link-icon {
  font-size: 64px;
  color: white;
  opacity: 0.8;
}

.article-info, .link-info {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent);
  color: white;
}

.article-domain, .link-domain {
  font-size: 14px;
  font-weight: 600;
  opacity: 0.9;
}

.image-error, .video-error {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-color);
}

.image-error .el-icon, .video-error .el-icon {
  font-size: 48px;
  color: var(--text-secondary);
}

.item-info {
  padding: 16px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.item-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.item-description {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0 0 12px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  flex: 1;
}

.item-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-secondary);
}

.item-type-badge {
  padding: 4px 8px;
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.1);
  font-weight: 500;
}

.type-image .item-type-badge {
  background: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.type-video .item-type-badge {
  background: rgba(250, 128, 114, 0.1);
  color: #fa8072;
}

.type-article .item-type-badge {
  background: rgba(139, 195, 74, 0.1);
  color: #8bc34a;
}

.type-link .item-type-badge {
  background: rgba(255, 193, 7, 0.1);
  color: #ffc107;
}

.item-actions {
  display: flex;
  justify-content: flex-end;
  padding: 0 16px 16px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.collection-item:hover .item-actions {
  opacity: 1;
}

.action-btn {
  padding: 4px 8px;
  font-size: 16px;
}

.delete-btn {
  color: #f56c6c;
}

.delete-btn:hover {
  color: #f56c6c;
  background: rgba(245, 108, 108, 0.1);
}

.empty-collection {
  grid-column: 1 / -1;
  padding: 80px 20px;
  text-align: center;
}

/* 搜索和筛选区域样式 */
.collection-filters {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 16px;
  background: var(--card-bg);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
}

.search-box {
  flex: 1;
  margin-right: 16px;
}

.search-input {
  max-width: 400px;
}

.filter-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.category-select {
  min-width: 180px;
}

/* 收藏项多选功能样式 */
.item-checkbox {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 10;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.collection-item:hover .item-checkbox {
  opacity: 1;
}

.collection-item.selected {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.collection-item.selected .item-checkbox {
  opacity: 1;
}

/* 分类管理对话框样式 */
.category-list {
  max-height: 300px;
  overflow-y: auto;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  margin-bottom: 8px;
  background: var(--card-bg);
  border-radius: 8px;
  transition: all 0.3s ease;
}

.category-item:hover {
  background: var(--bg-color);
}

.category-name-input {
  max-width: 200px;
}

.delete-category-btn {
  color: #f56c6c;
}

.delete-category-btn:hover {
  background: rgba(245, 108, 108, 0.1);
}

.add-category {
  margin-top: 16px;
  padding: 12px;
  background: rgba(64, 158, 255, 0.05);
  border-radius: 8px;
  border: 1px dashed var(--primary-color);
}

/* 批量移动对话框样式 */
.batch-move-form {
  padding: 16px 0;
}

.tag-input {
  width: 120px;
  margin-top: 8px;
}

.detected-resources {
  margin-top: 20px;
  border-top: 1px solid var(--border-color);
  padding-top: 10px;
}

.detected-resources h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: var(--text-secondary);
}

.resource-list {
  max-height: 150px;
  overflow-y: auto;
  border: 1px solid var(--border-color);
  border-radius: 4px;
}

.resource-section {
  margin-bottom: 10px;
}

.resource-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 8px;
  max-height: 200px;
  overflow-y: auto;
  padding: 4px;
}

.resource-grid-item {
  position: relative;
  aspect-ratio: 1;
  border-radius: 4px;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid var(--border-color);
}

.resource-grid-item.selected {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px var(--primary-color);
}

.resource-grid-item:hover {
  border-color: var(--primary-color);
}

.resource-grid-item .el-image {
  width: 100%;
  height: 100%;
}

.resource-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.resource-grid-item:hover .resource-overlay {
  opacity: 1;
}

.resource-overlay .el-icon {
  color: white;
  font-size: 20px;
}

.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  color: #909399;
}

.resource-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  cursor: pointer;
  border-bottom: 1px solid var(--border-color);
  transition: background-color 0.2s;
}

.resource-item:last-child {
  border-bottom: none;
}

.resource-item:hover {
  background-color: var(--hover-bg);
}

.resource-url {
  font-size: 12px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .collection-page {
    padding: 16px;
  }
  
  .collection-grid {
    grid-template-columns: 1fr;
  }
  
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .page-header h2 {
    font-size: 24px;
  }
}
</style>
