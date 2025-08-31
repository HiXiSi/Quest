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
          border
          style="width: 100%"
          :empty-text="records.length === 0 ? '暂无数据记录' : ''"
        >
          <!-- 序号列 -->
          <el-table-column type="index" width="50" align="center" />
          
          <!-- 动态字段列 -->
          <el-table-column
            v-for="field in visibleFields"
            :key="field.name"
            :prop="field.name"
            :label="field.label"
            :width="getColumnWidth(field)"
            show-overflow-tooltip
          >
            <template #default="{ row }">
              <span v-if="field.type === 'switch'">
                {{ formatFieldValue(row.data[field.name], field) }}
              </span>
              <el-tag v-else-if="field.type === 'select'" size="small">
                {{ formatFieldValue(row.data[field.name], field) }}
              </el-tag>
              <div v-else-if="field.type === 'checkbox'" class="checkbox-values">
                <el-tag
                  v-for="value in (row.data[field.name] || [])"
                  :key="value"
                  size="small"
                  style="margin-right: 4px;"
                >
                  {{ formatOptionLabel(value, field) }}
                </el-tag>
              </div>
              <el-rate
                v-else-if="field.type === 'rate'"
                :model-value="row.data[field.name] || 0"
                disabled
                show-score
                text-color="#ff9900"
                score-template="{value}"
              />
              <span v-else>{{ formatFieldValue(row.data[field.name], field) }}</span>
            </template>
          </el-table-column>
          
          <!-- 创建时间列 -->
          <el-table-column label="创建时间" width="180" align="center">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          
          <!-- 操作列 -->
          <el-table-column label="操作" width="200" align="center" fixed="right">
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
      width="800px"
      :close-on-click-modal="false"
      @close="handleDialogClose"
    >
      <el-form
        ref="recordFormRef"
        :model="recordForm"
        label-width="120px"
        :rules="formRules"
      >
        <div
          v-for="field in formSchema?.fields || []"
          :key="field.name"
        >
          <FormFieldRenderer
            :field="field"
            v-model="recordForm[field.name]"
            :preview="false"
          />
        </div>
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
          v-for="field in formSchema?.fields || []"
          :key="field.name"
          class="detail-item"
        >
          <label class="detail-label">{{ field.label }}：</label>
          <div class="detail-value">
            <span v-if="field.type === 'switch'">
              {{ formatFieldValue(currentRecord.data[field.name], field) }}
            </span>
            <el-tag v-else-if="field.type === 'select'" size="small">
              {{ formatFieldValue(currentRecord.data[field.name], field) }}
            </el-tag>
            <div v-else-if="field.type === 'checkbox'" class="checkbox-values">
              <el-tag
                v-for="value in (currentRecord.data[field.name] || [])"
                :key="value"
                size="small"
                style="margin-right: 4px;"
              >
                {{ formatOptionLabel(value, field) }}
              </el-tag>
            </div>
            <el-rate
              v-else-if="field.type === 'rate'"
              :model-value="currentRecord.data[field.name] || 0"
              disabled
              show-score
              text-color="#ff9900"
              score-template="{value}"
            />
            <span v-else>{{ formatFieldValue(currentRecord.data[field.name], field) || '-' }}</span>
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
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft, Plus, Search, Download, View, Edit, Delete
} from '@element-plus/icons-vue'
import api from '@/utils/api'
import FormFieldRenderer from '@/components/FormFieldRenderer.vue'

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
const isEditing = ref(false)
const saving = ref(false)
const exporting = ref(false)

// 表单数据
const recordFormRef = ref()
const recordForm = ref({})
const currentRecord = ref(null)
const editingRecordId = ref(null)

// 可见字段（用于表格显示）
const visibleFields = computed(() => {
  if (!formSchema.value?.fields) return []
  // 限制显示的字段数量，避免表格过宽
  return formSchema.value.fields.slice(0, 6)
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

// 获取表单结构
const fetchFormSchema = async () => {
  try {
    const response = await api.get(`/forms/${route.params.id}`)
    formSchema.value = response.data.data
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
  if (formSchema.value?.fields) {
    formSchema.value.fields.forEach(field => {
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
    case 'switch':
      return value ? '是' : '否'
    case 'date':
      return new Date(value).toLocaleDateString()
    case 'datetime':
      return new Date(value).toLocaleString()
    case 'time':
      return value
    case 'select':
    case 'radio':
      return formatOptionLabel(value, field)
    default:
      return value
  }
}

// 格式化选项标签
const formatOptionLabel = (value, field) => {
  if (!field.options) return value
  const option = field.options.find(opt => opt.value === value)
  return option ? option.label : value
}

// 格式化日期
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString()
}

// 获取列宽度
const getColumnWidth = (field) => {
  switch (field.type) {
    case 'switch':
      return 80
    case 'rate':
      return 120
    case 'date':
      return 120
    case 'datetime':
      return 160
    case 'checkbox':
      return 200
    case 'text':
      // 根据inputType决定列宽
      return field.inputType === 'multi' ? 250 : 150
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
      }
    }
  }
}
</style>