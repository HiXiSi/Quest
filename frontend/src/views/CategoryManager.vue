<template>
  <div class="page-container">
    <div class="card-container">
      <div class="header-actions">
        <h2>分类管理</h2>
        <el-button type="primary" @click="showAddDialog">
          <el-icon><Plus /></el-icon>
          新增分类
        </el-button>
      </div>
      
      <!-- 分类列表 -->
      <div v-loading="loading">
        <el-table :data="categories" style="width: 100%">
          <el-table-column prop="name" label="分类名称" width="200" />
          <el-table-column prop="description" label="描述" show-overflow-tooltip />
          <el-table-column prop="color" label="颜色" width="100">
            <template #default="{ row }">
              <div
                class="color-preview"
                :style="{ backgroundColor: row.color }"
              ></div>
            </template>
          </el-table-column>
          <el-table-column prop="icon" label="图标" width="100">
            <template #default="{ row }">
              <el-icon v-if="row.icon" :size="20">
                <component :is="row.icon" />
              </el-icon>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="{ row }">
              <el-button size="small" @click="editCategory(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button size="small" type="danger" @click="deleteCategory(row)">
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
        :title="isEdit ? '编辑分类' : '新增分类'"
        width="500px"
      >
        <el-form
          ref="formRef"
          :model="categoryForm"
          :rules="formRules"
          label-width="80px"
        >
          <el-form-item label="分类名称" prop="name">
            <el-input v-model="categoryForm.name" placeholder="请输入分类名称" />
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input
              v-model="categoryForm.description"
              type="textarea"
              :rows="3"
              placeholder="请输入分类描述"
            />
          </el-form-item>
          <el-form-item label="颜色" prop="color">
            <el-color-picker v-model="categoryForm.color" />
          </el-form-item>
          <el-form-item label="图标" prop="icon">
            <el-select v-model="categoryForm.icon" placeholder="选择图标" clearable>
              <el-option value="Folder" label="文件夹">
                <el-icon><Folder /></el-icon>
                文件夹
              </el-option>
              <el-option value="Document" label="文档">
                <el-icon><Document /></el-icon>
                文档
              </el-option>
              <el-option value="Picture" label="图片">
                <el-icon><Picture /></el-icon>
                图片
              </el-option>
              <el-option value="VideoCamera" label="视频">
                <el-icon><VideoCamera /></el-icon>
                视频
              </el-option>
            </el-select>
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
const categories = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

const categoryForm = ref({
  id: null,
  name: '',
  description: '',
  color: '#409EFF',
  icon: ''
})

const formRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 1, max: 50, message: '名称长度在1-50个字符', trigger: 'blur' }
  ]
}

const fetchCategories = async () => {
  loading.value = true
  try {
    const response = await api.get('/categories')
    categories.value = response.data.data || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
    ElMessage.error('获取分类列表失败')
  } finally {
    loading.value = false
  }
}

const showAddDialog = () => {
  isEdit.value = false
  categoryForm.value = {
    id: null,
    name: '',
    description: '',
    color: '#409EFF',
    icon: ''
  }
  dialogVisible.value = true
}

const editCategory = (category) => {
  isEdit.value = true
  categoryForm.value = { ...category }
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (isEdit.value) {
          await api.put(`/categories/${categoryForm.value.id}`, categoryForm.value)
          ElMessage.success('分类更新成功')
        } else {
          await api.post('/categories', categoryForm.value)
          ElMessage.success('分类创建成功')
        }
        dialogVisible.value = false
        fetchCategories()
      } catch (error) {
        console.error('保存分类失败:', error)
        ElMessage.error('保存分类失败')
      }
    }
  })
}

const deleteCategory = async (category) => {
  try {
    await ElMessageBox.confirm(`确定要删除分类 "${category.name}" 吗？`, '确认删除', {
      type: 'warning'
    })
    
    await api.delete(`/categories/${category.id}`)
    ElMessage.success('分类删除成功')
    fetchCategories()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除分类失败:', error)
      ElMessage.error('删除分类失败')
    }
  }
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchCategories()
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

.color-preview {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: 1px solid #ddd;
}
</style>