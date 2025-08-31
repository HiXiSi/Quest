<template>
  <div class="page-container">
    <div class="form-container">
      <div v-if="loading" class="loading-container">
        <el-loading 
          v-loading="loading"
          element-loading-text="加载表单中..."
          element-loading-background="rgba(0, 0, 0, 0.8)"
        />
      </div>
      
      <div v-else-if="error" class="error-container">
        <el-result
          icon="error"
          title="表单加载失败"
          :sub-title="error"
        >
          <template #extra>
            <el-button @click="loadForm">重新加载</el-button>
          </template>
        </el-result>
      </div>
      
      <div v-else-if="submitted" class="success-container">
        <el-result
          icon="success"
          title="提交成功"
          sub-title="您的表单已成功提交，感谢您的参与！"
        >
          <template #extra>
            <el-button @click="resetAndResubmit" type="primary">再次填写</el-button>
            <el-button @click="goBack">返回</el-button>
          </template>
        </el-result>
      </div>
      
      <div v-else-if="formSchema" class="form-content">
        <DynamicForm
          ref="dynamicFormRef"
          :schema="formSchema"
          :initial-data="initialData"
          :show-reset-button="true"
          submit-button-text="提交表单"
          @submit="handleFormSubmit"
          @reset="handleFormReset"
        />
        
        <!-- 表单信息 -->
        <div class="form-info">
          <el-divider />
          <div class="info-section">
            <p class="info-text">
              <el-icon><InfoFilled /></el-icon>
              表单创建时间：{{ formatDate(formSchema.created_at) }}
            </p>
            <p v-if="formSchema.updated_at !== formSchema.created_at" class="info-text">
              <el-icon><Edit /></el-icon>
              最后更新时间：{{ formatDate(formSchema.updated_at) }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { InfoFilled, Edit } from '@element-plus/icons-vue'
import api from '@/utils/api'
import DynamicForm from '@/components/DynamicForm.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const error = ref('')
const submitted = ref(false)
const formSchema = ref(null)
const dynamicFormRef = ref()

// 初始数据（如果是编辑模式）
const initialData = ref({})

// 表单ID
const formId = computed(() => route.params.id)

// 是否是编辑模式（从查询参数获取记录ID）
const isEditMode = computed(() => !!route.query.recordId)
const recordId = computed(() => route.query.recordId)

// 加载表单结构
const loadForm = async () => {
  loading.value = true
  error.value = ''
  
  try {
    // 获取表单结构
    const schemaResponse = await api.get(`/forms/${formId.value}`)
    formSchema.value = schemaResponse.data.data
    
    // 如果是编辑模式，加载现有记录
    if (isEditMode.value) {
      try {
        const recordResponse = await api.get(`/forms/records/${recordId.value}`)
        initialData.value = recordResponse.data.data.data || {}
      } catch (recordError) {
        console.error('加载记录失败:', recordError)
        ElMessage.warning('记录加载失败，将创建新记录')
      }
    }
  } catch (err) {
    console.error('加载表单失败:', err)
    error.value = err.response?.data?.message || '表单加载失败'
  } finally {
    loading.value = false
  }
}

// 处理表单提交
const handleFormSubmit = async ({ data, schema }) => {
  try {
    const submitData = {
      schema_id: parseInt(formId.value),
      data: data
    }
    
    if (isEditMode.value) {
      // 更新现有记录
      await api.put(`/forms/records/${recordId.value}`, submitData)
      ElMessage.success('记录更新成功')
    } else {
      // 创建新记录
      await api.post('/forms/records', submitData)
      ElMessage.success('表单提交成功')
    }
    
    submitted.value = true
  } catch (error) {
    console.error('表单提交失败:', error)
    ElMessage.error(error.response?.data?.message || '提交失败，请重试')
  }
}

// 处理表单重置
const handleFormReset = (data) => {
  ElMessage.info('表单已重置')
}

// 重置并重新填写
const resetAndResubmit = () => {
  submitted.value = false
  if (dynamicFormRef.value) {
    dynamicFormRef.value.resetForm()
  }
  
  // 如果是编辑模式，清除查询参数变为新建模式
  if (isEditMode.value) {
    router.replace({ 
      path: route.path,
      query: {}
    })
    initialData.value = {}
  }
}

// 返回上一页
const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push('/')
  }
}

// 格式化日期
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString()
}

onMounted(() => {
  loadForm()
})
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;

  .form-container {
    max-width: 900px;
    margin: 0 auto;
    background: white;
    border-radius: 12px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
    overflow: hidden;

    .loading-container {
      height: 400px;
      position: relative;
    }

    .error-container,
    .success-container {
      padding: 40px;
    }

    .form-content {
      padding: 40px;

      .form-info {
        margin-top: 20px;

        .info-section {
          .info-text {
            display: flex;
            align-items: center;
            margin: 8px 0;
            color: #909399;
            font-size: 13px;

            .el-icon {
              margin-right: 6px;
            }
          }
        }
      }
    }
  }
}

:deep(.dynamic-form) {
  .form-header {
    h2 {
      color: #303133;
      font-weight: 600;
    }
  }
  
  .form-actions {
    .el-button {
      padding: 12px 32px;
      font-size: 16px;
    }
    
    .el-button--primary {
      background: linear-gradient(45deg, #667eea, #764ba2);
      border: none;
      
      &:hover {
        opacity: 0.9;
      }
    }
  }
}
</style>