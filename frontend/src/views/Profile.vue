<template>
  <div class="page-container">
    <div class="card-container">
      <h2>个人中心</h2>
      
      <div v-loading="loading" class="profile-content">
        <el-form
          ref="formRef"
          :model="userForm"
          :rules="formRules"
          label-width="100px"
          style="max-width: 600px"
        >
          <el-form-item label="用户名">
            <el-input v-model="userForm.username" disabled />
          </el-form-item>
          
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="userForm.email" />
          </el-form-item>
          
          <el-form-item label="角色">
            <el-tag :type="userForm.role === 'admin' ? 'danger' : 'primary'">
              {{ userForm.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </el-form-item>
          
          <el-form-item label="头像">
            <div class="avatar-section">
              <el-avatar
                :size="80"
                :src="userForm.avatar"
                class="avatar-display"
              >
                <el-icon><User /></el-icon>
              </el-avatar>
              <div class="avatar-actions">
                <el-input
                  v-model="userForm.avatar"
                  placeholder="请输入头像 URL"
                  style="width: 300px"
                />
                <el-button @click="clearAvatar" style="margin-left: 10px">
                  清除
                </el-button>
              </div>
            </div>
          </el-form-item>
          
          <el-form-item label="注册时间">
            <span>{{ formatDate(userForm.created_at) }}</span>
          </el-form-item>
          
          <el-form-item label="更新时间">
            <span>{{ formatDate(userForm.updated_at) }}</span>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="updateProfile" :loading="updating">
              保存修改
            </el-button>
            <el-button @click="resetForm">
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const loading = ref(false)
const updating = ref(false)
const formRef = ref()

const userForm = ref({
  id: null,
  username: '',
  email: '',
  role: '',
  avatar: '',
  created_at: '',
  updated_at: ''
})

const formRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const fetchUserProfile = async () => {
  loading.value = true
  try {
    const response = await api.get('/users/profile')
    userForm.value = { ...response.data.data }
  } catch (error) {
    console.error('获取用户信息失败:', error)
    ElMessage.error('获取用户信息失败')
  } finally {
    loading.value = false
  }
}

const updateProfile = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      updating.value = true
      try {
        const updateData = {
          email: userForm.value.email,
          avatar: userForm.value.avatar
        }
        
        const response = await api.put('/users/profile', updateData)
        userForm.value = { ...response.data.data }
        
        // 更新全局用户信息
        if (userStore.updateUser) {
          userStore.updateUser(response.data.data)
        }
        
        ElMessage.success('个人信息更新成功')
      } catch (error) {
        console.error('更新个人信息失败:', error)
        ElMessage.error('更新个人信息失败')
      } finally {
        updating.value = false
      }
    }
  })
}

const resetForm = () => {
  fetchUserProfile()
}

const clearAvatar = () => {
  userForm.value.avatar = ''
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchUserProfile()
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

.profile-content {
  margin-top: 20px;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.avatar-display {
  flex-shrink: 0;
}

.avatar-actions {
  flex: 1;
  display: flex;
  align-items: center;
}
</style>