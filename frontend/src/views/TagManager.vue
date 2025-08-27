<!-- TagManager.vue -->
<template>
  <div class="page-container">
    <div class="card-container">
      <div class="header-actions">
        <h2>标签管理</h2>
        <el-button type="primary" @click="showAddDialog">
          <el-icon><Plus /></el-icon>
          新增标签
        </el-button>
      </div>
      
      <!-- 标签列表 -->
      <div v-loading="loading">
        <el-table :data="tags" style="width: 100%">
          <el-table-column prop="name" label="标签名称" width="200" />
          <el-table-column prop="description" label="描述" show-overflow-tooltip />
          <el-table-column prop="color" label="颜色" width="100">
            <template #default="{ row }">
              <el-tag :color="row.color || '#409EFF'" style="border: none; color: white;">
                {{ row.name }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="usage_count" label="使用次数" width="120">
            <template #default="{ row }">
              <el-badge :value="row.usage_count" :max="999" />
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="{ row }">
              <el-button size="small" @click="editTag(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button size="small" type="danger" @click="deleteTag(row)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <!-- 新增/编辑对话框 -->
      <el-dialog
        v-model="dialogVisible"
        :title="isEdit ? '编辑标签' : '新增标签'"
        width="500px"
      >
        <el-form
          ref="formRef"
          :model="tagForm"
          :rules="formRules"
          label-width="80px"
        >
          <el-form-item label="标签名称" prop="name">
            <el-input v-model="tagForm.name" placeholder="请输入标签名称" />
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input
              v-model="tagForm.description"
              type="textarea"
              :rows="3"
              placeholder="请输入标签描述"
            />
          </el-form-item>
          <el-form-item label="颜色" prop="color">
            <el-color-picker v-model="tagForm.color" />
          </el-form-item>
          <el-form-item label="预览">
            <el-tag :color="tagForm.color || '#409EFF'" style="border: none; color: white;">
              {{ tagForm.name || '标签名称' }}
            </el-tag>
          </el-form-item>
        </el-form>
        
        <template #footer>
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/utils/api'

const loading = ref(false)
const tags = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

const tagForm = ref({
  id: null,
  name: '',
  description: '',
  color: '#409EFF'
})

const formRules = {
  name: [
    { required: true, message: '请输入标签名称', trigger: 'blur' },
    { min: 1, max: 30, message: '名称长度在1-30个字符', trigger: 'blur' }
  ]
}

const fetchTags = async () => {
  loading.value = true
  try {
    const response = await api.get('/tags')
    // 后端返回的是分页数据，需要取 list 字段
    tags.value = response.data.data.list || []
  } catch (error) {
    console.error('获取标签列表失败:', error)
    ElMessage.error('获取标签列表失败')
  } finally {
    loading.value = false
  }
}

const showAddDialog = () => {
  isEdit.value = false
  tagForm.value = {
    id: null,
    name: '',
    description: '',
    color: '#409EFF'
  }
  dialogVisible.value = true
}

const editTag = (tag) => {
  isEdit.value = true
  tagForm.value = { 
    ...tag,
    // 确保 color字段有默认值
    color: tag.color || '#409EFF'
  }
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (isEdit.value) {
          await api.put(`/tags/${tagForm.value.id}`, tagForm.value)
          ElMessage.success('标签更新成功')
        } else {
          await api.post('/tags', tagForm.value)
          ElMessage.success('标签创建成功')
        }
        dialogVisible.value = false
        fetchTags()
      } catch (error) {
        console.error('保存标签失败:', error)
        ElMessage.error('保存标签失败')
      }
    }
  })
}

const deleteTag = async (tag) => {
  try {
    await ElMessageBox.confirm(`确定要删除标签 "${tag.name}" 吗？`, '确认删除', {
      type: 'warning'
    })
    
    await api.delete(`/tags/${tag.id}`)
    ElMessage.success('标签删除成功')
    fetchTags()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除标签失败:', error)
      ElMessage.error('删除标签失败')
    }
  }
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
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
</style>