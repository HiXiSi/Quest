<template>
  <div class="page-container">
    <div class="card-container">
      <div class="header-actions">
        <h2>文件管理</h2>
        <el-button type="primary" @click="$router.push('/upload')">
          <el-icon><Upload /></el-icon>
          上传文件
        </el-button>
      </div>
      
      <!-- 搜索和筛选 -->
      <div class="filter-bar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索文件名或描述"
          clearable
          @change="handleSearch"
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <el-select
          v-model="selectedCategory"
          placeholder="选择分类"
          clearable
          @change="handleFilterChange"
          class="filter-select"
        >
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>
        
        <el-select
          v-model="selectedFileType"
          placeholder="文件类型"
          clearable
          @change="handleFilterChange"
          class="filter-select"
        >
          <el-option label="图片" value="image" />
          <el-option label="文档" value="document" />
          <el-option label="PDF" value="pdf" />
          <el-option label="文本" value="text" />
          <el-option label="其他" value="other" />
        </el-select>
      </div>
      
      <!-- 文件列表 -->
      <div v-loading="loading" class="file-list">
        <div v-if="files.length === 0" class="no-data">
          <el-empty description="暂无文件" />
        </div>
        
        <div v-else class="file-grid">
          <div
            v-for="file in files"
            :key="file.id"
            class="file-card"
          >
            <div class="file-card-header">
              <div class="file-icon-wrapper">
                <el-icon size="24" :color="getFileTypeColor(file.file_type)">
                  <component :is="getFileTypeIcon(file.file_type)" />
                </el-icon>
              </div>
              <div class="file-actions">
                <el-dropdown trigger="click" placement="bottom-end">
                  <el-button size="small" text>
                    <el-icon><MoreFilled /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="previewFile(file)">
                        <el-icon><View /></el-icon> 预览
                      </el-dropdown-item>
                      <el-dropdown-item @click="editFile(file)">
                        <el-icon><Edit /></el-icon> 编辑
                      </el-dropdown-item>
                      <el-dropdown-item @click="downloadFile(file)">
                        <el-icon><Download /></el-icon> 下载
                      </el-dropdown-item>
                      <el-dropdown-item divided @click="deleteFile(file)">
                        <el-icon><Delete /></el-icon> 删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
            
            <div class="file-content" @click="previewFile(file)">
              <div class="file-name" :title="file.original_name">
                {{ file.original_name }}
              </div>
              
              <div class="file-description" v-if="file.description">
                {{ file.description }}
              </div>
              
              <div class="file-tags" v-if="file.tags && file.tags.length > 0">
                <el-tag
                  v-for="tag in file.tags.slice(0, 3)"
                  :key="tag.id"
                  size="small"
                  :color="tag.color || '#409EFF'"
                  style="border: none; color: white; margin-right: 4px; margin-bottom: 4px;"
                >
                  {{ tag.name }}
                </el-tag>
                <span v-if="file.tags.length > 3" class="more-tags">
                  +{{ file.tags.length - 3 }}
                </span>
              </div>
              
              <div class="file-meta">
                <span class="file-size">{{ formatFileSize(file.file_size) }}</span>
                <span class="file-date">{{ formatDate(file.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 分页 -->
      <el-pagination
        v-if="total > 0"
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
    
    <!-- 编辑文件对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑文件信息"
      width="600px"
    >
      <el-form
        ref="editFormRef"
        :model="editForm"
        label-width="80px"
      >
        <el-form-item label="文件名">
          <el-input v-model="editForm.original_name" readonly />
        </el-form-item>
        
        <el-form-item label="分类">
          <el-select v-model="editForm.category_id" placeholder="选择分类" clearable>
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="标签">
          <el-select
            v-model="editForm.tag_ids"
            multiple
            placeholder="选择标签"
            style="width: 100%"
          >
            <el-option
              v-for="tag in allTags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            >
              <span style="float: left">{{ tag.name }}</span>
              <span style="float: right; color: var(--el-text-color-secondary); font-size: 13px">
                <el-tag size="small" :color="tag.color || '#409EFF'" style="border: none; color: white;">
                  {{ tag.name }}
                </el-tag>
              </span>
            </el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="描述">
          <el-input
            v-model="editForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入文件描述"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveFileEdit" :loading="editLoading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Upload, Search, Download, Delete, Edit, View, MoreFilled,
  Picture, Document, EditPen, VideoPlay, Headset, Folder
} from '@element-plus/icons-vue'
import api from '@/utils/api'

const loading = ref(false)
const files = ref([])
const categories = ref([])
const allTags = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchKeyword = ref('')
const selectedCategory = ref('')
const selectedFileType = ref('')

// 编辑文件相关状态
const editDialogVisible = ref(false)
const editLoading = ref(false)
const editFormRef = ref()
const editForm = ref({
  id: null,
  original_name: '',
  category_id: null,
  tag_ids: [],
  description: ''
})

const fetchFiles = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (searchKeyword.value) params.keyword = searchKeyword.value
    if (selectedCategory.value) params.category_id = selectedCategory.value
    if (selectedFileType.value) params.file_type = selectedFileType.value
    
    const response = await api.get('/files', { params })
    files.value = response.data.data.list || []
    total.value = response.data.data.total || 0
  } catch (error) {
    console.error('获取文件列表失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const response = await api.get('/categories')
    categories.value = response.data.data || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

const fetchTags = async () => {
  try {
    const response = await api.get('/tags')
    allTags.value = response.data.data.list || []
  } catch (error) {
    console.error('获取标签列表失败:', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchFiles()
}

const handleFilterChange = () => {
  currentPage.value = 1
  fetchFiles()
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchFiles()
}

const handleCurrentChange = () => {
  fetchFiles()
}

const previewFile = async (file) => {
  // 预览文件
  if (file.file_type === 'image' || file.file_type === 'text' || file.file_type === 'pdf') {
    try {
      // 对于文本文件，获取内容并在新窗口显示
      if (file.file_type === 'text') {
        const response = await api.get(`/files/${file.id}/preview`)
        const newWindow = window.open('', '_blank')
        newWindow.document.write(`
          <html>
            <head><title>${file.original_name}</title></head>
            <body style="font-family: monospace; white-space: pre-wrap; padding: 20px;">
              ${response.data.content}
            </body>
          </html>
        `)
      } else {
        // 对于图片和PDF，获取blob并创建预览链接
        const response = await api.get(`/files/${file.id}/preview`, {
          responseType: 'blob'
        })
        const blob = new Blob([response.data], { type: file.mime_type })
        const url = window.URL.createObjectURL(blob)
        window.open(url, '_blank')
        // 延迟释放URL以确保文件能正常加载
        setTimeout(() => window.URL.revokeObjectURL(url), 10000)
      }
    } catch (error) {
      console.error('预览文件失败:', error)
      ElMessage.error('预览失败')
    }
  } else {
    downloadFile(file)
  }
}

const editFile = (file) => {
  editForm.value = {
    id: file.id,
    original_name: file.original_name,
    category_id: file.category_id || null,
    tag_ids: file.tags ? file.tags.map(tag => tag.id) : [],
    description: file.description || ''
  }
  editDialogVisible.value = true
}

const saveFileEdit = async () => {
  editLoading.value = true
  try {
    await api.put(`/files/${editForm.value.id}`, {
      category_id: editForm.value.category_id,
      tag_ids: editForm.value.tag_ids,
      description: editForm.value.description
    })
    
    ElMessage.success('文件信息更新成功')
    editDialogVisible.value = false
    fetchFiles()
  } catch (error) {
    console.error('更新文件信息失败:', error)
    ElMessage.error('更新失败')
  } finally {
    editLoading.value = false
  }
}

const downloadFile = async (file) => {
  try {
    const response = await api.get(`/files/${file.id}/download`, {
      responseType: 'blob'
    })
    
    // 创建下载链接
    const blob = new Blob([response.data])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = file.original_name
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('下载成功')
  } catch (error) {
    console.error('下载文件失败:', error)
    ElMessage.error('下载失败')
  }
}

const deleteFile = async (file) => {
  try {
    await ElMessageBox.confirm(`确定要删除文件 "${file.original_name}" 吗？`, '确认删除', {
      type: 'warning'
    })
    
    await api.delete(`/files/${file.id}`)
    ElMessage.success('文件已移到回收站')
    fetchFiles()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除文件失败:', error)
    }
  }
}

const getFileTypeIcon = (fileType) => {
  const iconMap = {
    image: 'Picture',
    document: 'Document',
    pdf: 'Document',
    text: 'EditPen',
    video: 'VideoPlay',
    audio: 'Headset',
    archive: 'Folder',
    other: 'Document'
  }
  return iconMap[fileType] || 'Document'
}

const getFileTypeColor = (fileType) => {
  const colorMap = {
    image: '#67c23a',
    document: '#409eff',
    pdf: '#f56c6c',
    text: '#e6a23c',
    video: '#909399',
    audio: '#909399',
    archive: '#909399',
    other: '#909399'
  }
  return colorMap[fileType] || '#909399'
}

const formatFileSize = (size) => {
  const units = ['B', 'KB', 'MB', 'GB']
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  return `${size.toFixed(1)} ${units[unitIndex]}`
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString()
}

onMounted(() => {
  fetchFiles()
  fetchCategories()
  fetchTags()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-container {
  background: white;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filter-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  align-items: center;
}

.search-input {
  width: 280px;
}

.filter-select {
  width: 160px;
}

.file-list {
  min-height: 400px;
}

.file-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.file-card {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: white;
  transition: all 0.3s ease;
  overflow: hidden;
}

.file-card:hover {
  border-color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
  transform: translateY(-2px);
}

.file-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px 8px;
  background: #f8f9fa;
}

.file-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.file-content {
  padding: 12px 16px 16px;
  cursor: pointer;
}

.file-name {
  font-weight: 500;
  font-size: 14px;
  color: #303133;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.4;
}

.file-description {
  font-size: 12px;
  color: #606266;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-tags {
  margin-bottom: 8px;
  min-height: 20px;
}

.more-tags {
  font-size: 12px;
  color: #909399;
}

.file-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #909399;
}

.file-size {
  font-weight: 500;
}

.file-date {
  opacity: 0.8;
}

.no-data {
  text-align: center;
  padding: 60px 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-input,
  .filter-select {
    width: 100%;
  }
  
  .file-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  }
}
</style>