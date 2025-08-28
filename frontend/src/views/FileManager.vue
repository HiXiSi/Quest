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
          <el-option label="视频" value="video" />
          <el-option label="三维模型" value="3d" />
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
            <!-- 第一行：文件类型图标、文件名称、操作按钮 -->
            <div class="file-card-header">
              <div class="file-info-line">
                <el-icon size="16" :color="getFileTypeColor(file.file_type, file)" class="file-type-icon">
                  <component :is="getFileTypeIcon(file.file_type, file)" />
                </el-icon>
                <span class="file-name" :title="file.original_name">{{ file.original_name }}</span>
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
            
            <!-- 预览区域 -->
            <div class="file-preview" @click="previewFile(file)">
              <!-- 图片预览 -->
              <div v-if="isImageFile(file)" class="preview-image">
                <img 
                  :src="getPreviewUrl(file)"
                  :alt="file.original_name"
                  @error="handleImageError($event, file)"
                  loading="lazy"
                />
                <div class="preview-fallback">
                  <el-icon size="48" :color="getFileTypeColor(file.file_type, file)">
                    <Picture />
                  </el-icon>
                </div>
              </div>
              
              <!-- 视频预览 -->
              <div v-else-if="isVideoFile(file)" class="preview-video">
                <video 
                  :src="getPreviewUrl(file)"
                  muted
                  preload="metadata"
                  @error="handleVideoError($event, file)"
                >
                  <source :src="getPreviewUrl(file)" :type="file.mime_type">
                </video>
                <div class="preview-fallback">
                  <el-icon size="48" :color="getFileTypeColor(file.file_type, file)">
                    <VideoPlay />
                  </el-icon>
                </div>
                <div class="video-overlay">
                  <el-icon size="32" color="white"><VideoPlay /></el-icon>
                </div>
              </div>
              
              <!-- 三维模型预览 -->
              <div v-else-if="is3DModelFile(file)" class="preview-model">
                <div class="model-icon">
                  <el-icon size="48" color="#9c27b0">
                    <Box />
                  </el-icon>
                </div>
                <div class="model-overlay">
                  <el-icon size="24" color="white"><View /></el-icon>
                  <span class="overlay-text">点击预览</span>
                </div>
                <div class="model-type-badge">
                  {{ getFileExtension(file).toUpperCase() }}
                </div>
              </div>
              
              <!-- 默认显示区域 -->
              <div v-else class="preview-default">
                <el-icon size="48" :color="getFileTypeColor(file.file_type, file)">
                  <component :is="getFileTypeIcon(file.file_type, file)" />
                </el-icon>
              </div>
            </div>
            
            <div class="file-content">              
              <!-- 标签一行 -->
              <div class="file-tags">
                <template v-if="file.tags && file.tags.length > 0">
                  <el-tag
                    v-for="tag in file.tags.slice(0, 3)"
                    :key="tag.id"
                    size="small"
                    :color="tag.color || '#409EFF'"
                    style="border: none; color: white; margin-right: 4px;"
                  >
                    {{ tag.name }}
                  </el-tag>
                  <span v-if="file.tags.length > 3" class="more-tags">
                    +{{ file.tags.length - 3 }}
                  </span>
                </template>
                <span v-else class="no-tags">无标签</span>
              </div>
              
              <!-- 描述一行 -->
              <div class="file-description">
                <span v-if="file.description">{{ file.description }}</span>
                <span v-else class="no-description">无描述</span>
              </div>
              
              <!-- 大小和日期 -->
              <div class="file-meta">
                <div class="file-category" v-if="file.category">
                  <el-icon size="12"><Folder /></el-icon>
                  <span>{{ file.category.name }}</span>
                </div>
                <div class="file-category" v-else>
                  <el-icon size="12"><Folder /></el-icon>
                  <span class="no-category">未分类</span>
                </div>
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
    
    <!-- 三维模型预览对话框 -->
    <el-dialog
      v-model="modelPreviewVisible"
      :title="'三维模型预览 - ' + (currentPreviewFile?.original_name || '')"
      width="80%"
      :close-on-click-modal="false"
      @close="handleCloseModelPreview"
    >
      <div class="model-preview-container">
        <ModelViewer
          v-if="modelPreviewVisible && currentPreviewFile"
          :model-url="getModelUrl(currentPreviewFile)"
          :file-extension="getFileExtension(currentPreviewFile)"
        />
      </div>
      
      <template #footer>
        <el-button @click="downloadFile(currentPreviewFile)" type="primary">
          <el-icon><Download /></el-icon>
          下载模型
        </el-button>
        <el-button @click="modelPreviewVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Upload, Search, Download, Delete, Edit, View, MoreFilled,
  Picture, Document, EditPen, VideoPlay, Headset, Folder, Box
} from '@element-plus/icons-vue'
import api from '@/utils/api'
import ModelViewer from '@/components/ModelViewer.vue'

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

// 模型预览相关状态
const modelPreviewVisible = ref(false)
const currentPreviewFile = ref(null)

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
  // 检查文件是否为三维模型
  if (is3DModelFile(file)) {
    currentPreviewFile.value = file
    modelPreviewVisible.value = true
    return
  }
  
  // 检查文件是否可预览
  if (isImageFile(file) || isVideoFile(file) || file.file_type === 'pdf' || file.file_type === 'text') {
    // 可预览的文件，在新窗口打开预览
    const previewUrl = getPreviewUrl(file)
    window.open(previewUrl, '_blank')
  } else {
    // 不可预览的文件，直接下载
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
  // 直接使用下载链接
  const downloadUrl = `/api/files/${file.id}/download`
  window.open(downloadUrl, '_blank')
  ElMessage.success('开始下载')
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

const getFileTypeIcon = (fileType, file = null) => {
  // 如果提供了file对象，检查是否为三维模型文件
  if (file && is3DModelFile(file)) {
    return 'Box'
  }
  
  // 检查是否为三维模型文件
  if (fileType === '3d' || fileType === 'model') {
    return 'Box'
  }
  
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

const getFileTypeColor = (fileType, file = null) => {
  // 如果提供了file对象，检查是否为三维模型文件
  if (file && is3DModelFile(file)) {
    return '#9c27b0'
  }
  
  // 检查是否为三维模型文件
  if (fileType === '3d' || fileType === 'model') {
    return '#9c27b0'
  }
  
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

// 判断是否为图片文件
const isImageFile = (file) => {
  const imageExtensions = ['png', 'jpg', 'jpeg', 'gif', 'webp', 'bmp', 'svg']
  const extension = file.original_name.split('.').pop()?.toLowerCase()
  return imageExtensions.includes(extension) || file.file_type === 'image'
}

// 判断是否为视频文件
const isVideoFile = (file) => {
  const videoExtensions = ['mp4', 'avi', 'mov', 'wmv', 'flv', 'webm', 'mkv']
  const extension = file.original_name.split('.').pop()?.toLowerCase()
  return videoExtensions.includes(extension) || file.file_type === 'video'
}

// 判断是否为三维模型文件
const is3DModelFile = (file) => {
  const modelExtensions = ['gltf', 'glb', 'fbx']
  const extension = file.original_name.split('.').pop()?.toLowerCase()
  return modelExtensions.includes(extension)
}

// 获取文件扩展名
const getFileExtension = (file) => {
  return file.original_name.split('.').pop()?.toLowerCase() || ''
}

// 获取模型文件URL
const getModelUrl = (file) => {
  return `/api/files/${file.id}/preview`
}

// 关闭模型预览对话框
const handleCloseModelPreview = () => {
  modelPreviewVisible.value = false
  currentPreviewFile.value = null
}

// 生成预览URL（简化版本）
const getPreviewUrl = (file) => {
  return `/api/files/${file.id}/preview`
}

// 图片加载错误处理
const handleImageError = (event, file) => {
  console.log('图片加载失败，显示默认图标')
  event.target.style.display = 'none'
  // 显示fallback图标
  const fallback = event.target.nextElementSibling
  if (fallback && fallback.classList.contains('preview-fallback')) {
    fallback.style.display = 'flex'
  }
}

// 视频加载错误处理
const handleVideoError = (event, file) => {
  console.log('视频加载失败，显示默认图标')
  event.target.style.display = 'none'
  // 显示fallback图标
  const fallback = event.target.parentElement.querySelector('.preview-fallback')
  if (fallback) {
    fallback.style.display = 'flex'
  }
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

  .card-container {
    background: white;
    padding: 24px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

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

      .search-input {
        width: 280px;
      }

      .filter-select {
        width: 160px;
      }
    }

    .file-list {
      min-height: 400px;

      .no-data {
        text-align: center;
        padding: 60px 0;
      }
    }

    .file-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
      gap: 16px;
      margin-bottom: 24px;

      .file-card {
        border: 1px solid #e4e7ed;
        border-radius: 8px;
        background: white;
        transition: all 0.3s ease;
        overflow: hidden;

        &:hover {
          border-color: #409eff;
          box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
          transform: translateY(-2px);

          .preview-image img {
            transform: scale(1.05);
          }

          .video-overlay {
            background: rgba(0, 0, 0, 0.8);
            transform: translate(-50%, -50%) scale(1.1);
          }
        }

        .file-card-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 8px 0 8px 8px;
          background: #f8f9fa;
          border-bottom: 1px solid #e4e7ed;

          .file-info-line {
            display: flex;
            align-items: center;
            flex: 1;
            min-width: 0;

            .file-type-icon {
              margin-right: 6px;
              flex-shrink: 0;
            }

            .file-name {
              font-size: 13px;
              color: #303133;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
              flex: 1;
              min-width: 0;
            }
          }
        }

        /* 文件预览区域样式 */
        .file-preview {
          height: 100px;
          position: relative;
          cursor: pointer;
          background: #f8f9fa;
          display: flex;
          align-items: center;
          justify-content: center;
          overflow: hidden;

          .preview-image {
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;

            img {
              max-width: 100%;
              max-height: 100%;
              object-fit: cover;
              border-radius: 0;
              transition: transform 0.3s ease;
            }
          }

          .preview-video {
            width: 100%;
            height: 100%;
            position: relative;
            display: flex;
            align-items: center;
            justify-content: center;

            video {
              width: 100%;
              height: 100%;
              object-fit: cover;
              border-radius: 0;
            }
          }

          .video-overlay {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            background: rgba(0, 0, 0, 0.6);
            border-radius: 50%;
            width: 48px;
            height: 48px;
            display: flex;
            align-items: center;
            justify-content: center;
            transition: all 0.3s ease;
            pointer-events: none;
          }

          .preview-default {
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
          }
          
          .preview-model {
            width: 100%;
            height: 100%;
            position: relative;
            display: flex;
            align-items: center;
            justify-content: center;
            background: linear-gradient(135deg, #e1bee7 0%, #f3e5f5 100%);
            cursor: pointer;
            transition: all 0.3s ease;
          
            &:hover {
              background: linear-gradient(135deg, #ce93d8 0%, #e1bee7 100%);
              transform: scale(1.02);
          
              .model-overlay {
                opacity: 1;
                transform: translate(-50%, -50%) scale(1);
              }
          
              .model-type-badge {
                transform: scale(1.1);
              }
            }
          
            .model-icon {
              opacity: 0.7;
              transition: opacity 0.3s ease;
            }
          
            .model-overlay {
              position: absolute;
              top: 50%;
              left: 50%;
              transform: translate(-50%, -50%) scale(0.8);
              background: rgba(156, 39, 176, 0.8);
              border-radius: 50%;
              width: 60px;
              height: 60px;
              display: flex;
              flex-direction: column;
              align-items: center;
              justify-content: center;
              opacity: 0;
              transition: all 0.3s ease;
              pointer-events: none;
          
              .overlay-text {
                font-size: 10px;
                color: white;
                margin-top: 2px;
                font-weight: 500;
              }
            }
          
            .model-type-badge {
              position: absolute;
              top: 8px;
              right: 8px;
              background: rgba(156, 39, 176, 0.9);
              color: white;
              padding: 2px 6px;
              border-radius: 4px;
              font-size: 10px;
              font-weight: bold;
              transition: transform 0.3s ease;
            }
          }

          .preview-fallback {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            display: none;
            align-items: center;
            justify-content: center;
            background: #f8f9fa;
            opacity: 0.9;
          }
        }

        .file-content {
          padding: 8px;

          .file-category {
            display: flex;
            align-items: center;
            font-size: 12px;
            color: #606266;

            .el-icon {
              margin-right: 4px;
            }
          }

          .no-category {
            color: #c0c4cc;
          }

          .file-tags {
            margin-bottom: 6px;
            min-height: 20px;
            display: flex;
            align-items: center;
            flex-wrap: wrap;
            gap: 2px;
          }

          .no-tags {
            font-size: 12px;
            color: #c0c4cc;
          }

          .file-description {
            font-size: 12px;
            color: #606266;
            margin-bottom: 4px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            line-height: 1.4;
          }

          .no-description {
            color: #c0c4cc;
          }

          .file-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: 12px;
            color: #909399;

            .file-size {
              font-weight: 500;
            }

            .file-date {
              opacity: 0.8;
            }
          }

          .more-tags {
            font-size: 12px;
            color: #909399;
          }
        }
      }
    }
  }
}

/* 模型预览对话框样式 */
.model-preview-container {
  height: 60vh;
  min-height: 400px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e4e7ed;
}
</style>