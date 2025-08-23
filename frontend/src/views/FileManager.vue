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
          style="width: 300px"
          clearable
          @change="handleSearch"
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
            class="file-item"
            @click="handleFileClick(file)"
          >
            <div class="file-icon">
              <el-icon size="48" :color="getFileTypeColor(file.file_type)">
                <component :is="getFileTypeIcon(file.file_type)" />
              </el-icon>
            </div>
            <div class="file-info">
              <div class="file-name" :title="file.original_name">
                {{ file.original_name }}
              </div>
              <div class="file-meta">
                <span>{{ formatFileSize(file.file_size) }}</span>
                <span>{{ formatDate(file.created_at) }}</span>
              </div>
            </div>
            <div class="file-actions">
              <el-button size="small" @click.stop="downloadFile(file)">
                <el-icon><Download /></el-icon>
              </el-button>
              <el-button size="small" type="danger" @click.stop="deleteFile(file)">
                <el-icon><Delete /></el-icon>
              </el-button>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/utils/api'

const loading = ref(false)
const files = ref([])
const categories = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchKeyword = ref('')
const selectedCategory = ref('')
const selectedFileType = ref('')

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

const handleFileClick = (file) => {
  // 预览文件
  if (file.file_type === 'image' || file.file_type === 'text' || file.file_type === 'pdf') {
    window.open(`/api/files/${file.id}/preview`, '_blank')
  } else {
    downloadFile(file)
  }
}

const downloadFile = (file) => {
  window.open(`/api/files/${file.id}/download`)
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
    audio: 'Headphone',
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
})
</script>

<style scoped>
.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filter-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.file-list {
  min-height: 400px;
}

.file-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.file-item {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.3s;
  background: white;
}

.file-item:hover {
  border-color: #409eff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.file-icon {
  text-align: center;
  margin-bottom: 12px;
}

.file-name {
  font-weight: 500;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-meta {
  display: flex;
  justify-content: space-between;
  color: #909399;
  font-size: 12px;
  margin-bottom: 12px;
}

.file-actions {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.no-data {
  text-align: center;
  padding: 60px 0;
}
</style>