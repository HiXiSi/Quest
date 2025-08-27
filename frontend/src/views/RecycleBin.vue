<!-- RecycleBin.vue -->
<template>
  <div class="page-container">
    <div class="card-container">
      <div class="header-actions">
        <h2>回收站</h2>
        <div class="actions">
          <el-button
            v-if="selectedFiles.length > 0"
            type="success"
            @click="batchRestore"
          >
            <el-icon><RefreshRight /></el-icon>
            批量恢复 ({{ selectedFiles.length }})
          </el-button>
          <el-button
            v-if="selectedFiles.length > 0"
            type="danger"
            @click="batchPermanentDelete"
          >
            <el-icon><Delete /></el-icon>
            批量彻底删除 ({{ selectedFiles.length }})
          </el-button>
          <el-button type="danger" @click="emptyRecycleBin">
            <el-icon><Delete /></el-icon>
            清空回收站
          </el-button>
        </div>
      </div>
      
      <!-- 文件列表 -->
      <div v-loading="loading" class="file-list">
        <div v-if="files.length === 0" class="no-data">
          <el-empty description="回收站为空" />
        </div>
        
        <div v-else>
          <el-table
            :data="files"
            @selection-change="handleSelectionChange"
            style="width: 100%"
          >
            <el-table-column type="selection" width="55" />
            <el-table-column label="文件名称" min-width="200">
              <template #default="{ row }">
                <div class="file-item">
                  <el-icon :size="20" :color="getFileTypeColor(row.file_type)">
                    <component :is="getFileTypeIcon(row.file_type)" />
                  </el-icon>
                  <span class="file-name">{{ row.original_name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="file_size" label="文件大小" width="120">
              <template #default="{ row }">
                {{ formatFileSize(row.file_size) }}
              </template>
            </el-table-column>
            <el-table-column prop="file_type" label="文件类型" width="100" />
            <el-table-column prop="deleted_at" label="删除时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.deleted_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button size="small" type="success" @click="restoreFile(row)">
                  <el-icon><RefreshRight /></el-icon>
                  恢复
                </el-button>
                <el-button size="small" type="danger" @click="permanentDelete(row)">
                  <el-icon><Delete /></el-icon>
                  彻底删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          
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
            style="margin-top: 20px"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  RefreshRight, Delete, Picture, VideoCamera, Headset, Document, FolderOpened
} from '@element-plus/icons-vue'
import api from '@/utils/api'

const loading = ref(false)
const files = ref([])
const selectedFiles = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const fetchDeletedFiles = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    const response = await api.get('/recycle', { params })
    files.value = response.data.data.list || []
    total.value = response.data.data.total || 0
  } catch (error) {
    console.error('获取回收站文件失败:', error)
    ElMessage.error('获取回收站文件失败')
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (selection) => {
  selectedFiles.value = selection
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchDeletedFiles()
}

const handleCurrentChange = () => {
  fetchDeletedFiles()
}

const restoreFile = async (file) => {
  try {
    await ElMessageBox.confirm(`确定要恢复文件 "${file.original_name}" 吗？`, '确认恢复', {
      type: 'success'
    })
    
    await api.post(`/files/${file.id}/restore`)
    ElMessage.success('文件恢复成功')
    fetchDeletedFiles()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('恢复文件失败:', error)
      ElMessage.error('恢复文件失败')
    }
  }
}

const batchRestore = async () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning('请选择要恢复的文件')
    return
  }
  
  try {
    await ElMessageBox.confirm(`确定要恢复 ${selectedFiles.value.length} 个文件吗？`, '确认批量恢复', {
      type: 'success'
    })
    
    const fileIds = selectedFiles.value.map(file => file.id)
    await api.post('/files/batch-restore', { file_ids: fileIds })
    ElMessage.success('批量恢复成功')
    fetchDeletedFiles()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量恢复失败:', error)
      ElMessage.error('批量恢复失败')
    }
  }
}

const permanentDelete = async (file) => {
  try {
    await ElMessageBox.confirm(
      `确定要彻底删除文件 "${file.original_name}" 吗？\n\n此操作不可恢复！`,
      '确认彻底删除',
      {
        type: 'error',
        confirmButtonText: '彻底删除',
        confirmButtonClass: 'el-button--danger'
      }
    )
    
    await api.delete(`/recycle/${file.id}`)
    ElMessage.success('文件已彻底删除')
    fetchDeletedFiles()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('彻底删除失败:', error)
      ElMessage.error('彻底删除失败')
    }
  }
}

const batchPermanentDelete = async () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning('请选择要彻底删除的文件')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要彻底删除 ${selectedFiles.value.length} 个文件吗？\n\n此操作不可恢复！`,
      '确认批量彻底删除',
      {
        type: 'error',
        confirmButtonText: '批量彻底删除',
        confirmButtonClass: 'el-button--danger'
      }
    )
    
    const fileIds = selectedFiles.value.map(file => file.id)
    await api.post('/recycle/batch-delete', { file_ids: fileIds })
    ElMessage.success('批量彻底删除成功')
    fetchDeletedFiles()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量彻底删除失败:', error)
      ElMessage.error('批量彻底删除失败')
    }
  }
}

const emptyRecycleBin = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清空回收站吗？\n\n此操作将彻底删除所有文件，不可恢复！',
      '确认清空回收站',
      {
        type: 'error',
        confirmButtonText: '清空回收站',
        confirmButtonClass: 'el-button--danger'
      }
    )
    
    await api.delete('/recycle/empty')
    ElMessage.success('回收站已清空')
    fetchDeletedFiles()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('清空回收站失败:', error)
      ElMessage.error('清空回收站失败')
    }
  }
}

const getFileTypeIcon = (fileType) => {
  const iconMap = {
    image: 'Picture',
    video: 'VideoCamera',
    audio: 'Headset',
    document: 'Document',
    pdf: 'Document',
    text: 'Document',
    archive: 'FolderOpened',
    other: 'Document'
  }
  return iconMap[fileType] || 'Document'
}

const getFileTypeColor = (fileType) => {
  const colorMap = {
    image: '#67C23A',
    video: '#E6A23C',
    audio: '#409EFF',
    document: '#909399',
    pdf: '#F56C6C',
    text: '#909399',
    archive: '#E6A23C',
    other: '#C0C4CC'
  }
  return colorMap[fileType] || '#C0C4CC'
}

const formatFileSize = (size) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / 1024 / 1024).toFixed(1) + ' MB'
  return (size / 1024 / 1024 / 1024).toFixed(1) + ' GB'
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchDeletedFiles()
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

.actions {
  display: flex;
  gap: 10px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.no-data {
  text-align: center;
  padding: 40px 0;
}
</style>