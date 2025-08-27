<template>
  <div class="page-container">
    <div class="card-container">
      <div class="header-actions">
        <h2>文件上传</h2>
        <el-button @click="$router.back()">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
      </div>
      
      <!-- 主体内容：左右布局 -->
      <div class="upload-content">
        <!-- 左侧：文件选择区域 -->
        <div class="upload-left">
          <h3 class="section-title">
            <el-icon><UploadFilled /></el-icon>
            选择文件
          </h3>
          
          <el-upload
            ref="uploadRef"
            class="upload-demo"
            drag
            :auto-upload="false"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            :file-list="fileList"
            multiple
            :limit="10"
            :show-file-list="true"
          >
            <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击选择</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持多文件选择，单个文件大小不超过100MB
              </div>
            </template>
          </el-upload>
          
          <!-- 文件列表 -->
          <div v-if="selectedFiles.length > 0" class="file-list">
            <h4>已选择文件 ({{ selectedFiles.length }})</h4>
            <div 
              class="file-item" 
              v-for="(file, index) in selectedFiles" 
              :key="`${file.name}-${file.size}-${file.lastModified || index}`"
            >
              <div class="file-info">
                <el-icon><Document /></el-icon>
                <span class="file-name">{{ file.name }}</span>
                <span class="file-size">({{ formatFileSize(file.size) }})</span>
              </div>
              <el-button 
                size="small" 
                type="danger" 
                text 
                @click="removeFile(index)"
                :title="`删除 ${file.name}`"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
        
        <!-- 右侧：文件信息表单 -->
        <div class="upload-right">
          <h3 class="section-title">
            <el-icon><Setting /></el-icon>
            文件信息
          </h3>
          
          <el-form :model="fileForm" label-width="70px" label-position="top" class="upload-form">
            <!-- 分类和标签同行 -->
            <div class="form-row">
              <el-form-item label="分类" class="form-item-half">
                <el-select
                  v-model="fileForm.category_id"
                  placeholder="选择分类（可选）"
                  clearable
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
              
              <el-form-item label="标签" class="form-item-half">
                <el-select
                  v-model="fileForm.tag_ids"
                  placeholder="选择标签（可选）"
                  multiple
                  clearable
                  style="width: 100%"
                >
                  <el-option
                    v-for="tag in tags"
                    :key="tag.id"
                    :label="tag.name"
                    :value="tag.id"
                  />
                </el-select>
              </el-form-item>
            </div>
            
            <el-form-item label="描述">
              <el-input
                v-model="fileForm.description"
                type="textarea"
                :rows="4"
                placeholder="文件描述（可选）"
              />
            </el-form-item>
            
            <el-form-item class="form-actions-item">
              <div class="form-actions">
                <el-button 
                  type="primary" 
                  @click="uploadFiles"
                  :disabled="selectedFiles.length === 0 || uploading"
                  :loading="uploading"
                  class="upload-btn"
                >
                  <el-icon v-if="!uploading"><Upload /></el-icon>
                  {{ uploading ? '正在上传...' : `上传文件 (${selectedFiles.length})` }}
                </el-button>
                
                <div class="secondary-actions">
                  <el-button 
                    @click="clearForm"
                    :icon="RefreshRight"
                    class="action-btn"
                  >
                    清空表单
                  </el-button>
                  
                  <el-button 
                    @click="clearFiles"
                    :icon="Delete"
                    :disabled="selectedFiles.length === 0"
                    class="action-btn"
                  >
                    清空文件
                  </el-button>
                </div>
              </div>
            </el-form-item>
          </el-form>
          
          <!-- 上传进度 -->
          <div v-if="uploadProgress.show" class="upload-progress">
            <div class="progress-title">上传进度</div>
            <el-progress
              :percentage="uploadProgress.percent"
              :status="uploadProgress.status"
            >
              <span>{{ uploadProgress.text }}</span>
            </el-progress>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  RefreshRight, ArrowLeft, UploadFilled, Setting, Upload, Document, Delete 
} from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const uploadRef = ref()
const fileList = ref([])
const selectedFiles = ref([]) // 存储选中的文件
const categories = ref([])
const tags = ref([])
const uploading = ref(false)

const fileForm = ref({
  category_id: '',
  tag_ids: [],
  description: ''
})

const uploadProgress = ref({
  show: false,
  percent: 0,
  status: '',
  text: ''
})

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await api.get('/categories')
    categories.value = response.data.data || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 获取标签列表
const fetchTags = async () => {
  try {
    const response = await api.get('/tags')
    tags.value = response.data.data.list || []
  } catch (error) {
    console.error('获取标签列表失败:', error)
  }
}

// 文件选择变化
const handleFileChange = (file, fileList) => {
  // 检查文件大小
  const isLt100M = file.size / 1024 / 1024 < 100
  if (!isLt100M) {
    ElMessage.error('文件大小不能超过100MB!')
    // 移除超大文件
    const index = fileList.findIndex(item => item.uid === file.uid)
    if (index > -1) {
      fileList.splice(index, 1)
    }
    return false
  }
  
  // 更新选中文件列表
  selectedFiles.value = fileList.map(item => item.raw || item)
  // 同步更新组件内部的fileList
  fileList.value = fileList
}

// 移除文件
const handleFileRemove = (file, fileList) => {
  selectedFiles.value = fileList.map(item => item.raw || item)
  fileList.value = fileList
}

// 手动移除文件
const removeFile = (index) => {
  if (index < 0 || index >= selectedFiles.value.length) {
    ElMessage.error('文件索引错误')
    return
  }
  
  // 获取要删除的文件
  const fileToRemove = selectedFiles.value[index]
  console.log('删除文件:', fileToRemove.name, '索引:', index)
  
  // 从selectedFiles中移除
  selectedFiles.value.splice(index, 1)
  
  // 重新构建el-upload的fileList
  const newFileList = selectedFiles.value.map((file, newIndex) => ({
    name: file.name,
    size: file.size,
    raw: file,
    uid: file.lastModified + newIndex, // 生成唯一ID
    status: 'ready'
  }))
  
  // 直接替换el-upload的fileList
  if (uploadRef.value) {
    uploadRef.value.fileList = newFileList
  }
  
  // 同步更新fileList
  fileList.value = newFileList
  
  console.log('删除后的文件数量:', selectedFiles.value.length)
  ElMessage.success(`已移除 ${fileToRemove.name}`)
}

// 上传文件
const uploadFiles = async () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning('请先选择文件')
    return
  }
  
  uploading.value = true
  uploadProgress.value.show = true
  
  let successCount = 0
  let failCount = 0
  
  try {
    for (let i = 0; i < selectedFiles.value.length; i++) {
      const file = selectedFiles.value[i]
      
      // 更新进度
      uploadProgress.value = {
        show: true,
        percent: Math.round((i / selectedFiles.value.length) * 100),
        status: '',
        text: `正在上传 ${file.name}... (${i + 1}/${selectedFiles.value.length})`
      }
      
      try {
        // 创建 FormData
        const formData = new FormData()
        formData.append('file', file)
        formData.append('category_id', fileForm.value.category_id || '')
        formData.append('tag_ids', fileForm.value.tag_ids.join(','))
        formData.append('description', fileForm.value.description || '')
        
        // 上传文件
        await api.post('/files/upload', formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        successCount++
      } catch (error) {
        console.error(`上传 ${file.name} 失败:`, error)
        failCount++
      }
    }
    
    // 更新最终进度
    uploadProgress.value = {
      show: true,
      percent: 100,
      status: failCount > 0 ? 'exception' : 'success',
      text: `上传完成：成功 ${successCount} 个，失败 ${failCount} 个`
    }
    
    // 显示结果消息
    if (failCount === 0) {
      ElMessage.success(`所有文件上传成功！`)
      // 清空文件列表
      clearFiles()
    } else if (successCount === 0) {
      ElMessage.error('所有文件上传失败')
    } else {
      ElMessage.warning(`部分文件上传失败：成功 ${successCount} 个，失败 ${failCount} 个`)
    }
    
    // 3秒后隐藏进度条
    setTimeout(() => {
      uploadProgress.value.show = false
    }, 3000)
    
  } finally {
    uploading.value = false
  }
}

// 清空表单
const clearForm = () => {
  fileForm.value = {
    category_id: '',
    tag_ids: [],
    description: ''
  }
  ElMessage.success('表单已清空')
}

// 清空文件
const clearFiles = () => {
  selectedFiles.value = []
  fileList.value = []
  // 使用el-upload的clearFiles方法清空所有文件
  if (uploadRef.value) {
    uploadRef.value.clearFiles()
  }
}

// 格式化文件大小
const formatFileSize = (size) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / 1024 / 1024).toFixed(1) + ' MB'
  return (size / 1024 / 1024 / 1024).toFixed(1) + ' GB'
}

onMounted(() => {
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
      margin-bottom: 24px;
    }

    /* 主体内容左右布局 */
    .upload-content {
      display: flex;
      gap: 32px;
      min-height: 500px;
      align-items: flex-start;

      /* 响应式设计 */
      @media (max-width: 768px) {
        flex-direction: column;
        gap: 24px;
      }

      .upload-left {
        flex: 1;
        min-width: 300px;
        max-width: 500px;

        /* 响应式设计 */
        @media (max-width: 768px) {
          max-width: none;
        }

        .section-title {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 16px;
          color: #303133;
          font-size: 16px;
          font-weight: 600;
        }

        /* 上传组件样式 */
        .upload-demo {
          margin-bottom: 20px;

          :deep(.el-upload-dragger) {
            height: 180px;
          }
        }

        /* 文件列表样式 */
        .file-list {
          background: #f8f9fa;
          border-radius: 8px;
          padding: 16px;
          border: 1px solid #e4e7ed;

          h4 {
            margin: 0 0 12px 0;
            color: #303133;
            font-size: 14px;
          }

          .file-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 8px 0;
            border-bottom: 1px solid #ebeef5;

            &:last-child {
              border-bottom: none;
            }

            .file-info {
              display: flex;
              align-items: center;
              gap: 8px;
              flex: 1;
              min-width: 0;

              .file-name {
                font-weight: 500;
                color: #303133;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
              }

              .file-size {
                color: #909399;
                font-size: 12px;
                flex-shrink: 0;
              }
            }
          }
        }
      }

      .upload-right {
        flex: 1;
        min-width: 0;
        max-width: 480px; /* 限制最大宽度 */

        .section-title {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 16px;
          color: #303133;
          font-size: 16px;
          font-weight: 600;
        }

        /* 表单样式 */
        .upload-form {
          .form-row {
            display: flex;
            gap: 16px;
            margin-bottom: 16px;

            .form-item-half {
              flex: 1;
              margin-bottom: 0;
            }
          }
        }

        .form-actions-item {
          margin-bottom: 0;
        }

        .form-actions {
          display: flex;
          flex-direction: column;
          gap: 16px;

          .upload-btn {
            width: 100%;
            height: 44px;
            font-size: 16px;
            font-weight: 500;
          }

          .secondary-actions {
            display: flex;
            gap: 12px;

            .action-btn {
              flex: 1;
              height: 36px;
              font-size: 14px;
            }
          }
        }

        /* 进度条样式 */
        .upload-progress {
          margin-top: 20px;
          padding: 16px;
          background: #f8f9fa;
          border-radius: 8px;
          border: 1px solid #e4e7ed;

          .progress-title {
            font-size: 14px;
            color: #303133;
            margin-bottom: 8px;
            font-weight: 500;
          }
        }
      }
    }

    /* 组件内部样式修正 */
    :deep(.el-upload-list) {
      display: none;
    }

    :deep(.el-form-item) {
      margin-bottom: 20px;
    }

    :deep(.el-form-item__label) {
      font-weight: 500;
      color: #303133;
      font-size: 14px;
      margin-bottom: 6px;
    }

    :deep(.el-textarea__inner) {
      min-height: 80px;
    }

    :deep(.el-select) {
      .el-input__wrapper {
        border-radius: 6px;
      }
    }
  }
}
</style>