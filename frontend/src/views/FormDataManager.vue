<template>
  <div class="page-container">
    <div class="card-container">
      <!-- 页面头部 -->
      <div class="header-section">
        <div class="header-left">
          <el-button @click="goBack" size="large">
            <el-icon><ArrowLeft /></el-icon>
            返回表单列表
          </el-button>
          <div class="form-info">
            <h2>{{ formSchema?.name || '加载中...' }}</h2>
            <p v-if="formSchema?.description">{{ formSchema.description }}</p>
          </div>
        </div>
        <div class="header-actions">
          <el-button @click="showAddDialog" type="primary" size="large">
            <el-icon><Plus /></el-icon>
            添加记录
          </el-button>
        </div>
      </div>
      
      <!-- 搜索和筛选 -->
      <div class="filter-section">
        <div class="filter-left">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索记录..."
            clearable
            @keyup.enter="handleSearch"
            @clear="handleSearch"
            style="width: 300px;"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
        <div class="filter-right">
          <el-button @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="exportData" :loading="exporting">
            <el-icon><Download /></el-icon>
            导出数据
          </el-button>
        </div>
      </div>
      
      <!-- 数据表格 -->
      <div class="table-section">
        <el-table
          v-loading="loading"
          :data="records"
          stripe
          style="width: 100%"
          :empty-text="records.length === 0 ? '暂无数据记录' : ''"
        >
          <!-- 序号列 -->
          <el-table-column type="index" width="60" align="center" label="#" />
          
          <!-- 记录ID列 -->
          <el-table-column label="ID" width="80" align="center">
            <template #default="{ row }">
              <el-tag type="info" size="small">{{ row.id }}</el-tag>
            </template>
          </el-table-column>
          
          <!-- 动态字段列 -->
          <el-table-column
            v-for="field in visibleFields"
            :key="field.name"
            :prop="field.name"
            :label="field.label"
            :min-width="getColumnWidth(field)"
            show-overflow-tooltip
          >
            <template #default="{ row }">
              <!-- 布尔类型 -->
              <span v-if="field.type === 'boolean'">
                {{ formatFieldValue(row.data[field.name], field) }}
              </span>
              <!-- 单选枚举 -->
              <el-tag v-else-if="field.type === 'single_enum'" size="small">
                {{ formatFieldValue(row.data[field.name], field) }}
              </el-tag>
              <!-- 多选枚举 -->
              <div v-else-if="field.type === 'multi_enum'" class="checkbox-values">
                <el-tag
                  v-for="value in (row.data[field.name] || [])"
                  :key="value"
                  size="small"
                  style="margin-right: 4px;"
                >
                  {{ formatEnumOptionLabel(value, field) }}
                </el-tag>
              </div>
              <!-- 唯一ID -->
              <el-tag v-else-if="field.type === 'unique_id'" type="info" size="small">
                {{ row.data[field.name] || '-' }}
              </el-tag>
              <!-- 时间类型 -->
              <span v-else-if="field.type === 'datetime'">
                {{ formatDateTimeValue(row.data[field.name], field) }}
              </span>
              <!-- 数值类型 -->
              <span v-else-if="field.type === 'integer' || field.type === 'float'">
                {{ formatNumberValue(row.data[field.name], field) }}
              </span>
              <!-- 字符串和其他类型 -->
              <span v-else>
                <!-- 文件上传类型显示为链接或预览 -->
                <div v-if="field.input_type === 'file' && row.data[field.name]">
                  <div 
                    v-if="isImageUrl(row.data[field.name])" 
                    class="file-preview image-preview small"
                    @click="previewFile(row.data[field.name])"
                  >
                    <img :src="toAbsoluteUrl(row.data[field.name])" alt="Preview" />
                  </div>
                  <div 
                    v-else-if="isVideoUrl(row.data[field.name])"
                    class="file-preview video-preview small"
                    @click="previewFile(row.data[field.name])"
                  >
                    <video :src="toAbsoluteUrl(row.data[field.name])" muted></video>
                    <div class="play-overlay">
                      <el-icon class="play-icon"><VideoPlay /></el-icon>
                    </div>
                  </div>
                  <a 
                    v-else
                    :href="toAbsoluteUrl(row.data[field.name])"
                    target="_blank"
                    class="file-link"
                  >
                    查看文件
                  </a>
                </div>
                <span v-else>
                  {{ formatFieldValue(row.data[field.name], field) }}
                </span>
              </span>
            </template>
          </el-table-column>
          
          <!-- 创建时间列 -->
          <el-table-column label="创建时间" min-width="150" align="center">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          
          <!-- 更新时间列 -->
          <el-table-column label="更新时间" min-width="150" align="center">
            <template #default="{ row }">
              {{ formatDate(row.updated_at) }}
            </template>
          </el-table-column>
          
          <!-- 操作列 -->
          <el-table-column label="操作" width="240" align="center" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="viewRecord(row)">
                <el-icon><View /></el-icon>
                查看
              </el-button>
              <el-button size="small" type="primary" @click="editRecord(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button size="small" type="danger" @click="deleteRecord(row)">
                <el-icon><Delete /></el-icon>
                删除
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
          style="margin-top: 20px; text-align: right;"
        />
      </div>
    </div>
    
    <!-- 添加/编辑记录对话框 -->
    <el-dialog
      v-model="formDialogVisible"
      :title="isEditing ? '编辑记录' : '添加记录'"
      width="900px"
      :close-on-click-modal="false"
      @close="handleDialogClose"
    >
      <el-form
        ref="recordFormRef"
        :model="recordForm"
        label-width="120px"
        :rules="formRules"
      >
        <el-row :gutter="24">
          <el-col
            v-for="field in editableFields"
            :key="field.name"
            :span="getFieldSpan(field)"
          >
            <FormFieldRenderer
              :field="field"
              v-model="recordForm[field.name]"
              :preview="false"
            />
          </el-col>
        </el-row>
      </el-form>
      
      <template #footer>
        <el-button @click="formDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveRecord" :loading="saving">
          {{ isEditing ? '更新' : '保存' }}
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 查看记录对话框 -->
    <el-dialog
      v-model="viewDialogVisible"
      title="查看记录详情"
      width="700px"
    >
      <div v-if="currentRecord" class="record-detail">
        <div
          v-for="field in viewFields"
          :key="field.name"
          class="detail-item"
        >
          <label class="detail-label">{{ field.label }}：</label>
          <div class="detail-value">
            <!-- 布尔类型 -->
            <span v-if="field.type === 'boolean'">
              {{ formatFieldValue(currentRecord.data[field.name], field) }}
            </span>
            <!-- 单选枚举 -->
            <el-tag v-else-if="field.type === 'single_enum'" size="small">
              {{ formatFieldValue(currentRecord.data[field.name], field) }}
            </el-tag>
            <!-- 多选枚举 -->
            <div v-else-if="field.type === 'multi_enum'" class="checkbox-values">
              <el-tag
                v-for="value in (currentRecord.data[field.name] || [])"
                :key="value"
                size="small"
                style="margin-right: 4px;"
              >
                {{ formatEnumOptionLabel(value, field) }}
              </el-tag>
            </div>
            <!-- 唯一ID -->
            <el-tag v-else-if="field.type === 'unique_id'" type="info" size="small">
              {{ currentRecord.id || '-' }}
            </el-tag>
            <!-- 时间类型 -->
            <span v-else-if="field.type === 'datetime'">
              {{ formatDateTimeValue(currentRecord.data[field.name], field) }}
            </span>
            <!-- 数值类型 -->
            <span v-else-if="field.type === 'integer' || field.type === 'float'">
              {{ formatNumberValue(currentRecord.data[field.name], field) }}
            </span>
            <!-- 其他类型 -->
            <span v-else>
              <!-- 文件上传类型显示为链接或预览 -->
              <div v-if="field.input_type === 'file' && currentRecord.data[field.name]">
                <div 
                  v-if="isImageUrl(currentRecord.data[field.name])"
                  class="file-preview image-preview"
                  @click="previewFile(currentRecord.data[field.name])"
                >
                  <img :src="toAbsoluteUrl(currentRecord.data[field.name])" alt="Preview" />
                </div>
                <div 
                  v-else-if="isVideoUrl(currentRecord.data[field.name])"
                  class="file-preview video-preview"
                  @click="previewFile(currentRecord.data[field.name])"
                >
                  <video :src="toAbsoluteUrl(currentRecord.data[field.name])" muted></video>
                  <div class="play-overlay">
                    <el-icon class="play-icon"><VideoPlay /></el-icon>
                  </div>
                </div>
                <a 
                  v-else
                  :href="toAbsoluteUrl(currentRecord.data[field.name])"
                  target="_blank"
                  class="file-link"
                >
                  查看文件
                </a>
              </div>
              <span v-else>
                {{ formatFieldValue(currentRecord.data[field.name], field) || '-' }}
              </span>
            </span>
          </div>
        </div>
        
        <div class="detail-item">
          <label class="detail-label">创建时间：</label>
          <div class="detail-value">{{ formatDate(currentRecord.created_at) }}</div>
        </div>
        
        <div class="detail-item">
          <label class="detail-label">更新时间：</label>
          <div class="detail-value">{{ formatDate(currentRecord.updated_at) }}</div>
        </div>
      </div>
      
      <template #footer>
        <el-button @click="viewDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="editRecord(currentRecord)">编辑此记录</el-button>
      </template>
    </el-dialog>
    
    <!-- 文件预览对话框 -->
    <el-dialog
      v-model="previewDialogVisible"
      :title="previewTitle"
      width="800px"
      center
    >
      <div class="preview-container">
        <img v-if="isImageUrl(currentPreviewUrl)" :src="currentPreviewUrl" alt="Preview" class="preview-image" />
        <video v-else-if="isVideoUrl(currentPreviewUrl)" :src="currentPreviewUrl" controls class="preview-video"></video>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Edit, Delete, VideoPlay } from '@element-plus/icons-vue'
import api from '@/utils/api'
import FormFieldRenderer from '@/components/FormFieldRenderer.vue'
import { toAbsoluteUrl } from '@/utils/urlHelper'

const route = useRoute()
const router = useRouter()

// 表单结构数据
const formSchema = ref(null)
const records = ref([])

// 分页和搜索
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchKeyword = ref('')

// 对话框状态
const formDialogVisible = ref(false)
const viewDialogVisible = ref(false)
const previewDialogVisible = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const exporting = ref(false)

// 表单数据
const recordFormRef = ref()
const recordForm = ref({})
const currentRecord = ref(null)
const editingRecordId = ref(null)

// 预览相关
const currentPreviewUrl = ref('')
const previewTitle = ref('')

// 可见字段（用于表格显示）
const visibleFields = computed(() => {
  if (!formSchema.value?.fields) return []
  // 过滤掉系统字段，这些字段在表格中单独显示
  return formSchema.value.fields.filter(field => {
    // 过滤掉系统管理的字段
    if (field.type === 'unique_id') return false
    return true
  }).slice(0, 6) // 限制显示字段数量，避免表格过宽
})

// 可编辑字段（用于表单显示）
const editableFields = computed(() => {
  if (!formSchema.value?.fields) return []
  return formSchema.value.fields.filter(field => {
    // 过滤掉不允许用户输入的字段
    if (field.type === 'unique_id') return false
    return true
  })
})

// 查看对话框字段（用于查看记录详情）
const viewFields = computed(() => {
  if (!formSchema.value?.fields) return []
  // 查看对话框显示所有字段，包括唯一ID
  return formSchema.value.fields
})

// 表单验证规则
const formRules = computed(() => {
  const rules = {}
  if (formSchema.value?.fields) {
    formSchema.value.fields.forEach(field => {
      if (field.required) {
        rules[field.name] = [
          { required: true, message: `请输入${field.label}`, trigger: 'blur' }
        ]
      }
    })
  }
  return rules
})

// 判断是否为图片URL
const isImageUrl = (url) => {
  if (!url) return false
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp']
  return imageExtensions.some(ext => url.toLowerCase().endsWith(ext))
}

// 判断是否为视频URL
const isVideoUrl = (url) => {
  if (!url) return false
  const videoExtensions = ['.mp4', '.avi', '.mov', '.wmv', '.flv', '.webm']
  return videoExtensions.some(ext => url.toLowerCase().endsWith(ext))
}

// 文件预览
const previewFile = (url) => {
  currentPreviewUrl.value = toAbsoluteUrl(url)
  previewTitle.value = isImageUrl(url) ? '图片预览' : '视频预览'
  previewDialogVisible.value = true
}

// 获取表单结构
const fetchFormSchema = async () => {
  try {
    const response = await api.get(`/forms/${route.params.id}`)
    const formData = response.data.data
    
    // 解析schema字段，从 API 返回的数据中提取字段信息
    let fields = []
    if (formData.schema && formData.schema.fields) {
      fields = formData.schema.fields
    }
    
    // 构建完整的表单结构
    formSchema.value = {
      id: formData.id,
      name: formData.name,
      description: formData.description,
      fields: fields
    }
    
    console.log('表单结构:', formSchema.value) // 调试信息
  } catch (error) {
    console.error('获取表单结构失败:', error)
    ElMessage.error('获取表单结构失败')
    goBack()
  }
}

// 获取记录列表
const fetchRecords = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    
    const response = await api.get(`/forms/${route.params.id}/records`, { params })
    records.value = response.data.data.list || []
    total.value = response.data.data.total || 0
  } catch (error) {
    console.error('获取记录列表失败:', error)
    ElMessage.error('获取记录列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchRecords()
}

// 分页处理
const handleSizeChange = () => {
  currentPage.value = 1
  fetchRecords()
}

const handleCurrentChange = () => {
  fetchRecords()
}

// 返回表单列表
const goBack = () => {
  router.push('/forms')
}

// 显示添加对话框
const showAddDialog = () => {
  isEditing.value = false
  editingRecordId.value = null
  resetForm()
  formDialogVisible.value = true
}

// 查看记录
const viewRecord = (record) => {
  currentRecord.value = record
  console.log('查看记录:', record)
  console.log('表单结构字段:', formSchema.value?.fields)
  viewDialogVisible.value = true
}

// 编辑记录
const editRecord = (record) => {
  isEditing.value = true
  editingRecordId.value = record.id
  currentRecord.value = record
  
  // 填充表单数据
  recordForm.value = { ...record.data }
  
  viewDialogVisible.value = false
  formDialogVisible.value = true
}

// 删除记录
const deleteRecord = async (record) => {
  try {
    await ElMessageBox.confirm('确定要删除这条记录吗？删除后无法恢复。', '确认删除', {
      type: 'warning'
    })
    
    await api.delete(`/forms/records/${record.id}`)
    ElMessage.success('记录删除成功')
    fetchRecords()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除记录失败:', error)
      ElMessage.error('删除记录失败')
    }
  }
}

// 保存记录
const saveRecord = async () => {
  try {
    await recordFormRef.value.validate()
  } catch (error) {
    ElMessage.warning('请完善必填字段')
    return
  }
  
  saving.value = true
  try {
    const data = {
      schema_id: parseInt(route.params.id),
      data: recordForm.value
    }
    
    if (isEditing.value) {
      await api.put(`/forms/records/${editingRecordId.value}`, data)
      ElMessage.success('记录更新成功')
    } else {
      await api.post('/forms/records', data)
      ElMessage.success('记录添加成功')
    }
    
    formDialogVisible.value = false
    fetchRecords()
  } catch (error) {
    console.error('保存记录失败:', error)
    ElMessage.error('保存记录失败')
  } finally {
    saving.value = false
  }
}

// 重置表单
const resetForm = () => {
  recordForm.value = {}
  // 只为可编辑字段设置默认值
  if (editableFields.value) {
    editableFields.value.forEach(field => {
      if (field.default_value) {
        recordForm.value[field.name] = field.default_value
      }
    })
  }
}

// 对话框关闭处理
const handleDialogClose = () => {
  if (recordFormRef.value) {
    recordFormRef.value.resetFields()
  }
  resetForm()
}

// 导出数据
const exportData = async () => {
  exporting.value = true
  try {
    // 这里可以实现数据导出功能
    // 简单示例：导出为JSON格式
    const allRecords = []
    let page = 1
    let hasMore = true
    
    while (hasMore) {
      const response = await api.get(`/forms/${route.params.id}/records`, {
        params: { page, page_size: 100 }
      })
      const pageRecords = response.data.data.list || []
      allRecords.push(...pageRecords)
      
      hasMore = pageRecords.length === 100
      page++
    }
    
    // 创建下载链接
    const dataStr = JSON.stringify(allRecords, null, 2)
    const dataBlob = new Blob([dataStr], { type: 'application/json' })
    const url = URL.createObjectURL(dataBlob)
    
    const link = document.createElement('a')
    link.href = url
    link.download = `${formSchema.value.name}_数据导出_${new Date().toLocaleDateString()}.json`
    link.click()
    
    URL.revokeObjectURL(url)
    ElMessage.success('数据导出成功')
  } catch (error) {
    console.error('导出数据失败:', error)
    ElMessage.error('导出数据失败')
  } finally {
    exporting.value = false
  }
}

// 格式化字段值
const formatFieldValue = (value, field) => {
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  
  switch (field.type) {
    case 'boolean':
      return value ? '是' : '否'
    case 'single_enum':
      return formatEnumOptionLabel(value, field)
    case 'multi_enum':
      if (Array.isArray(value)) {
        return value.map(v => formatEnumOptionLabel(v, field)).join(', ')
      }
      return value
    case 'unique_id':
      return value
    case 'datetime':
      return formatDateTimeValue(value, field)
    case 'integer':
    case 'float':
      return formatNumberValue(value, field)
    default:
      return value
  }
}

// 格式化枚举选项标签
const formatEnumOptionLabel = (value, field) => {
  if (!field.enum_options) return value
  const option = field.enum_options.find(opt => opt.value === value)
  return option ? option.label : value
}

// 格式化时间值
const formatDateTimeValue = (value, field) => {
  if (!value) return '-'
  
  const date = new Date(value)
  
  switch (field.time_format) {
    case 'date':
      return date.toLocaleDateString()
    case 'time':
      return date.toLocaleTimeString()
    case 'datetime':
    default:
      return date.toLocaleString()
  }
}

// 格式化数值
const formatNumberValue = (value, field) => {
  if (value === null || value === undefined) return '-'
  
  if (field.type === 'float' && field.precision) {
    return Number(value).toFixed(field.precision)
  }
  
  return value.toString()
}

// 获取字段在表单中的占位
const getFieldSpan = (field) => {
  // 多行文本、多选枚举和文件上传独占一行
  if (field.type === 'string' && field.input_type === 'textarea') {
    return 24
  }
  if (field.type === 'string' && field.input_type === 'file') {
    return 24
  }
  if (field.type === 'multi_enum') {
    return 24
  }
  // 其他字段一行放两个
  return 12
}

// 格式化日期
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString()
}

// 获取列宽度
const getColumnWidth = (field) => {
  switch (field.type) {
    case 'boolean':
      return 80
    case 'unique_id':
      return 120
    case 'integer':
      return 100
    case 'float':
      return 120
    case 'datetime':
      // 根据时间格式决定宽度
      switch (field.time_format) {
        case 'date':
          return 120
        case 'time':
          return 100
        case 'datetime':
        default:
          return 160
      }
    case 'single_enum':
      return 120
    case 'multi_enum':
      return 200
    case 'string':
      // 根据input_type决定列宽
      if (field.input_type === 'textarea') {
        return 250
      }
      if (field.input_type === 'file') {
        return 120
      }
      return 150
    default:
      return 150
  }
}

onMounted(() => {
  fetchFormSchema()
  fetchRecords()
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

    .header-section {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 24px;
      padding-bottom: 16px;
      border-bottom: 1px solid #e4e7ed;

      .header-left {
        display: flex;
        align-items: flex-start;
        gap: 16px;

        .form-info {
          h2 {
            margin: 0 0 4px 0;
            font-size: 20px;
            color: #303133;
          }

          p {
            margin: 0;
            color: #606266;
            font-size: 14px;
          }
        }
      }
    }

    .filter-section {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;

      .filter-right {
        display: flex;
        gap: 12px;
      }
    }

    .table-section {
      .checkbox-values {
        .el-tag {
          margin-bottom: 4px;
        }
      }
      
      .file-link {
        color: #409eff;
        text-decoration: none;
        
        &:hover {
          text-decoration: underline;
        }
      }
      
      .file-preview {
        cursor: pointer;
        border-radius: 4px;
        overflow: hidden;
        position: relative;
        display: inline-block;
        
        &.small {
          width: 60px;
          height: 60px;
          
          &.image-preview {
            img {
              width: 100%;
              height: 100%;
              object-fit: cover;
            }
          }
          
          &.video-preview {
            background: #000;
            
            video {
              width: 100%;
              height: 100%;
              object-fit: cover;
            }
            
            .play-overlay {
              position: absolute;
              top: 0;
              left: 0;
              width: 100%;
              height: 100%;
              display: flex;
              align-items: center;
              justify-content: center;
              background: rgba(0, 0, 0, 0.5);
              
              .play-icon {
                font-size: 16px;
                color: white;
              }
            }
          }
        }
      }
    }
  }

  .record-detail {
    .detail-item {
      display: flex;
      margin-bottom: 16px;
      align-items: flex-start;

      .detail-label {
        width: 120px;
        font-weight: 500;
        color: #606266;
        flex-shrink: 0;
      }

      .detail-value {
        flex: 1;
        color: #303133;

        .checkbox-values {
          .el-tag {
            margin-bottom: 4px;
          }
        }
        
        .file-preview {
          cursor: pointer;
          border-radius: 4px;
          overflow: hidden;
          position: relative;
          
          &.image-preview {
            max-width: 300px;
            max-height: 300px;
            
            img {
              width: 100%;
              height: 100%;
              object-fit: contain;
            }
          }
          
          &.video-preview {
            width: 100%;
            max-width: 300px;
            height: 200px;
            background: #000;
            
            video {
              width: 100%;
              height: 100%;
              object-fit: contain;
            }
            
            .play-overlay {
              position: absolute;
              top: 0;
              left: 0;
              width: 100%;
              height: 100%;
              display: flex;
              align-items: center;
              justify-content: center;
              background: rgba(0, 0, 0, 0.5);
              
              .play-icon {
                font-size: 40px;
                color: white;
              }
            }
          }
        }
      }
    }
  }
  
  .preview-container {
    text-align: center;
    
    .preview-image {
      max-width: 100%;
      max-height: 600px;
    }
    
    .preview-video {
      width: 100%;
      max-height: 600px;
    }
  }
}
</style>