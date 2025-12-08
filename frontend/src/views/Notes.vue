<template>
  <div class="notes-page">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="notes-list-card">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>笔记列表</span>
              <el-button type="primary" size="small" @click="createNote">新建</el-button>
            </div>
          </template>
          <el-input v-model="search" placeholder="搜索..." clearable style="margin-bottom: 15px" />
          <div class="notes-list">
            <div v-for="note in filteredNotes" :key="note.id" class="note-item" :class="{active: current?.id === note.id}" @click="selectNote(note)">
              <h4>{{ note.title || '无标题' }}</h4>
              <p>{{ note.content?.substring(0,50) }}...</p>
              <small>{{ new Date(note.updated_at).toLocaleString('zh-CN') }}</small>
            </div>
            <el-empty v-if="notes.length === 0" />
          </div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card v-if="current">
          <template #header>
            <div style="display: flex; gap: 10px">
              <el-input v-model="current.title" placeholder="标题..." style="flex: 1" />
              <el-button type="primary" @click="save">保存</el-button>
              <el-button type="danger" @click="del">删除</el-button>
            </div>
          </template>
          <el-input v-model="current.content" type="textarea" :rows="20" placeholder="内容..." />
          <el-input v-model="current.tags" placeholder="标签" style="margin-top: 10px" />
        </el-card>
        <el-empty v-else description="选择笔记" />
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'
const notes = ref([])
const current = ref(null)
const search = ref('')
const filteredNotes = computed(() => notes.value.filter(n => !search.value || n.title?.includes(search.value) || n.content?.includes(search.value)))
const load = async () => { try { notes.value = await api.get('/notes') } catch { ElMessage.error('加载失败') } }
const createNote = () => current.value = { title: '', content: '', tags: '' }
const selectNote = (n) => current.value = {...n}
const save = async () => { try { if(current.value.id) await api.put(`/notes/${current.value.id}`, current.value); else current.value = await api.post('/notes', current.value); ElMessage.success('保存成功'); load() } catch { ElMessage.error('保存失败') } }
const del = async () => { try { await ElMessageBox.confirm('确定删除?'); await api.delete(`/notes/${current.value.id}`); ElMessage.success('删除成功'); current.value = null; load() } catch(e) { if(e!=='cancel') ElMessage.error('删除失败') } }
onMounted(load)
</script>

<style scoped>
.notes-page { min-height: 70vh; }
.notes-list-card :deep(.el-card__body) { max-height: 70vh; overflow-y: auto; }
.notes-list { display: flex; flex-direction: column; gap: 10px; }
.note-item { padding: 12px; border: 1px solid var(--border-color); border-radius: 6px; cursor: pointer; transition: all 0.3s; }
.note-item:hover, .note-item.active { background: var(--hover-bg); border-color: #409EFF; }
.note-item h4 { margin: 0 0 8px 0; color: var(--text-primary); }
.note-item p { margin: 0 0 8px 0; font-size: 13px; color: var(--text-secondary); }
.note-item small { font-size: 12px; color: var(--text-secondary); }
</style>
