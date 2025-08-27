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
      
      <!-- 上传区域 -->
      <div class="upload-section">
        <el-upload
          ref="uploadRef"
          class="upload-demo"
          drag
          :action="uploadUrl"
          :headers="uploadHeaders"
          :data="uploadData"
          :before-upload="beforeUpload"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          :on-progress="handleUploadProgress"
          :file-list="fileList"
          multiple
          :limit="10"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">
            将文件拖到此处，或<em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              支持多文件上传，单个文件大小不超过100MB
            </div>
          </template>
        </el-upload>
      </div>
      
      <!-- 文件信息设置 -->
      <div class="file-settings">
        <h3 style="margin-bottom: 16px; color: #303133;">
          <el-icon style="margin-right: 8px;"><Setting /></el-icon>
          文件信息设置
        </h3>
        <el-alert
          title="这些设置将应用于新上传的文件"
          type="info"
          :closable="false"
          style="margin-bottom: 16px;"
        />
        <el-form :model="fileForm" label-width="100px">
          <el-form-item label="分类">
            <el-select
              v-model="fileForm.category_id"
              placeholder="选择分类（可选）"
              clearable
              style="width: 200px"
            >
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
              v-model="fileForm.tag_ids"
              placeholder="选择标签（可选）"
              multiple
              clearable
              style="width: 300px"
            >
              <el-option
                v-for="tag in tags"
                :key="tag.id"
                :label="tag.name"
                :value="tag.id"
              />
            </el-select>
          </el-form-item>
          
          <el-form-item label="描述">
            <el-input
              v-model="fileForm.description"
              type="textarea"
              :rows="3"
              placeholder="文件描述（可选）"
              style="width: 400px"
            />
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              @click="clearForm"
              :icon="RefreshRight"
            >
              清空设置
            </el-button>
            <el-text type="info" style="margin-left: 16px;">
              提示：拖拽或选择文件上传时，会自动应用这些设置
            </el-text>
          </el-form-item>
        </el-form>
      </div>
      
      <!-- 上传进度 -->
      <div v-if="uploadProgress.show" class="upload-progress">
        <el-progress
          :percentage="uploadProgress.percent"
          :status="uploadProgress.status"
        >
          <span>{{ uploadProgress.text }}</span>
        </el-progress>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { RefreshRight, ArrowLeft, UploadFilled, Setting } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const uploadRef = ref()
const fileList = ref([])
const categories = ref([])
const tags = ref([])

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

const uploadUrl = '/api/files/upload'
const uploadHeaders = {
  'Authorization': `Bearer ${userStore.token}`
}

const uploadData = ref({})

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await api.get('/categories')
    // 分类API直接返回数组，不是分页数据
    categories.value = response.data.data || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 获取标签列表
const fetchTags = async () => {
  try {
    const response = await api.get('/tags')
    // 后端返回的是分页数据，需要取 list 字段
    tags.value = response.data.data.list || []
  } catch (error) {
    console.error('获取标签列表失败:', error)
  }
}

// 上传前检查
const beforeUpload = (file) => {
  const isLt100M = file.size / 1024 / 1024 < 100
  if (!isLt100M) {
    ElMessage.error('文件大小不能超过100MB!')
    return false
  }
  
  // 设置上传数据
  uploadData.value = {
    category_id: fileForm.value.category_id || '',
    tag_ids: fileForm.value.tag_ids.join(','),
    description: fileForm.value.description || ''
  }
  
  return true
}

// 上传成功
const handleUploadSuccess = (response, file) => {
  ElMessage.success(`${file.name} 上传成功！`)
  uploadProgress.value.show = false
  
  // 上传成功后显示提示，问是否清空设置
  ElMessage({
    message: '文件上传成功！是否继续使用当前设置上传更多文件？',
    type: 'success',
    duration: 3000
  })
}

// 上传失败
const handleUploadError = (err, file) => {
  console.error('上传失败:', err)
  ElMessage.error(`${file.name} 上传失败！`)
  uploadProgress.value.show = false
}

// 上传进度
const handleUploadProgress = (event, file) => {
  uploadProgress.value = {
    show: true,
    percent: Math.round(event.percent),
    status: '', // 使用空字符串表示默认状态
    text: `正在上传 ${file.name}...`
  }
}

// 清空表单
const clearForm = () => {
  fileForm.value = {
    category_id: '',
    tag_ids: [],
    description: ''
  }
  ElMessage.success('设置已清空')
}

onMounted(() => {
  fetchCategories()
  fetchTags()
})
</script>

<style scoped>
.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.upload-section {
  margin-bottom: 30px;
}

.file-settings {
  margin-bottom: 20px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.upload-progress {
  margin-top: 20px;
}

.page-container {
  padding: 20px;
}

.card-container {
  background: white;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>