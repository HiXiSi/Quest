<template>
  <div class="page-container">
    <div class="card-container">
      <div class="header-actions">
        <h2>表单管理</h2>
        <el-button type="primary" @click="showCreateDialog = true">
          <el-icon><Plus /></el-icon>
          创建表单
        </el-button>
      </div>
      
      <!-- 搜索筛选 -->
      <div class="filter-bar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索表单名称或描述"
          clearable
          @change="handleSearch"
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
      
      <!-- 表单列表 -->
      <div v-loading="loading" class="form-list">
        <div v-if="forms.length === 0" class="no-data">
          <el-empty description="暂无表单" />
        </div>
        
        <div v-else class="form-grid">
          <div
            v-for="form in forms"
            :key="form.id"
            class="form-card"
            @click="viewFormData(form)"
          >
            <!-- 表单卡片头部 -->
            <div class="form-card-header">
              <div class="form-info-line">
                <el-icon size="16" color="#409eff" class="form-type-icon">
                  <DocumentCopy />
                </el-icon>
                <span class="form-name" :title="form.name">{{ form.name }}</span>
              </div>
              <div class="form-actions" @click.stop>
                <el-dropdown trigger="click" placement="bottom-end">
                  <el-button size="small" text>
                    <el-icon><MoreFilled /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="viewFormData(form)">
                        <el-icon><View /></el-icon> 查看数据
                      </el-dropdown-item>
                      <el-dropdown-item @click="editForm(form)">
                        <el-icon><Edit /></el-icon> 编辑结构
                      </el-dropdown-item>
                      <el-dropdown-item @click="addRecord(form)">
                        <el-icon><Plus /></el-icon> 添加记录
                      </el-dropdown-item>
                      <el-dropdown-item divided @click="deleteForm(form)">
                        <el-icon><Delete /></el-icon> 删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
            
            <!-- 表单预览区域 -->
            <div class="form-preview">
              <div class="form-icon">
                <el-icon size="48" color="#409eff">
                  <DocumentCopy />
                </el-icon>
              </div>
              <div class="form-summary">
                <div class="field-count">{{ getFieldCount(form) }} 个字段</div>
                <div class="record-count" v-if="form.record_count !== undefined">
                  {{ form.record_count }} 条记录
                </div>
              </div>
            </div>
            
            <div class="form-content">
              <!-- 描述 -->
              <div class="form-description">
                <span v-if="form.description">{{ form.description }}</span>
                <span v-else class="no-description">无描述</span>
              </div>
              
              <!-- 创建信息 -->
              <div class="form-meta">
                <span class="form-creator">{{ form.user?.username || '未知用户' }}</span>
                <span class="form-date">{{ formatDate(form.created_at) }}</span>
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
    
    <!-- 创建表单对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      title="创建新表单"
      width="500px"
    >
      <el-form :model="createForm" label-width="80px" ref="createFormRef">
        <el-form-item label="表单名称" prop="name" :rules="[{ required: true, message: '请输入表单名称' }]">
          <el-input v-model="createForm.name" placeholder="请输入表单名称" />
        </el-form-item>
        <el-form-item label="表单描述">
          <el-input
            v-model="createForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入表单描述（可选）"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="createFormAndDesign" :loading="createLoading">
          创建并设计
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, Search, DocumentCopy, View, Edit, Delete, MoreFilled
} from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'

const router = useRouter()

const loading = ref(false)
const forms = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchKeyword = ref('')

// 创建表单相关
const showCreateDialog = ref(false)
const createLoading = ref(false)
const createFormRef = ref()
const createForm = ref({
  name: '',
  description: ''
})

// 获取表单列表
const fetchForms = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    
    const response = await api.get('/forms', { params })
    forms.value = response.data.data.list || []
    total.value = response.data.data.total || 0
  } catch (error) {
    console.error('获取表单列表失败:', error)
    ElMessage.error('获取表单列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchForms()
}

// 分页处理
const handleSizeChange = () => {
  currentPage.value = 1
  fetchForms()
}

const handleCurrentChange = () => {
  fetchForms()
}

// 创建表单并跳转到设计器
const createFormAndDesign = async () => {
  if (!createFormRef.value) return
  
  const valid = await createFormRef.value.validate().catch(() => false)
  if (!valid) return
  
  createLoading.value = true
  try {
    const response = await api.post('/forms', {
      name: createForm.value.name,
      description: createForm.value.description,
      fields: [] // 不添加默认字段，让用户自己设计
    })
    
    ElMessage.success('表单创建成功')
    showCreateDialog.value = false
    
    // 重置表单
    createForm.value = { name: '', description: '' }
    
    // 跳转到表单设计器
    router.push(`/forms/${response.data.data.id}/design`)
  } catch (error) {
    console.error('创建表单失败:', error)
    ElMessage.error('创建表单失败')
  } finally {
    createLoading.value = false
  }
}

// 查看表单数据
const viewFormData = (form) => {
  router.push(`/forms/${form.id}/data`)
}

// 编辑表单结构
const editForm = (form) => {
  router.push(`/forms/${form.id}/design`)
}

// 添加记录
const addRecord = (form) => {
  router.push(`/forms/${form.id}/records/new`)
}

// 删除表单
const deleteForm = async (form) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除表单 "${form.name}" 吗？删除后无法恢复！`,
      '确认删除',
      { type: 'warning' }
    )
    
    await api.delete(`/forms/${form.id}`)
    ElMessage.success('表单删除成功')
    fetchForms()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除表单失败:', error)
      ElMessage.error('删除表单失败')
    }
  }
}

// 获取字段数量
const getFieldCount = (form) => {
  try {
    // 从 schema.fields 中获取字段数量
    if (form.schema && form.schema.fields) {
      return form.schema.fields.length || 0
    }
    // 兼容旧格式：如果 schema 是字符串，尝试解析
    if (typeof form.schema === 'string') {
      const schema = JSON.parse(form.schema)
      return schema.fields?.length || 0
    }
    return 0
  } catch {
    return 0
  }
}

// 格式化日期
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString()
}

onMounted(() => {
  fetchForms()
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
    }

    .form-list {
      min-height: 400px;

      .no-data {
        text-align: center;
        padding: 60px 0;
      }
    }

    .form-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
      gap: 16px;
      margin-bottom: 24px;

      .form-card {
        border: 1px solid #e4e7ed;
        border-radius: 8px;
        background: white;
        transition: all 0.3s ease;
        overflow: hidden;
        cursor: pointer;

        &:hover {
          border-color: #409eff;
          box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
          transform: translateY(-2px);
        }

        .form-card-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 12px 16px;
          background: #f8f9fa;
          border-bottom: 1px solid #e4e7ed;

          .form-info-line {
            display: flex;
            align-items: center;
            flex: 1;
            min-width: 0;

            .form-type-icon {
              margin-right: 8px;
              flex-shrink: 0;
            }

            .form-name {
              font-size: 14px;
              font-weight: 500;
              color: #303133;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
              flex: 1;
              min-width: 0;
            }
          }

          .form-actions {
            flex-shrink: 0;
          }
        }

        .form-preview {
          height: 120px;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          background: linear-gradient(135deg, #f0f7ff 0%, #e6f3ff 100%);
          position: relative;

          .form-icon {
            margin-bottom: 8px;
          }

          .form-summary {
            text-align: center;
            
            .field-count {
              font-size: 14px;
              color: #409eff;
              font-weight: 500;
            }

            .record-count {
              font-size: 12px;
              color: #909399;
              margin-top: 2px;
            }
          }
        }

        .form-content {
          padding: 12px 16px;

          .form-description {
            font-size: 12px;
            color: #606266;
            margin-bottom: 8px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            line-height: 1.4;

            .no-description {
              color: #c0c4cc;
            }
          }

          .form-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: 12px;
            color: #909399;

            .form-creator {
              font-weight: 500;
            }

            .form-date {
              opacity: 0.8;
            }
          }
        }
      }
    }
  }
}
</style>