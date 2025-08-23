<template>
  <div class="dashboard">
    <el-row :gutter="20" class="mb-20">
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-icon files">
            <el-icon size="40"><Document /></el-icon>
          </div>
          <div class="stat-content">
            <h3>{{ stats.fileCount || 0 }}</h3>
            <p>总文件数</p>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-icon categories">
            <el-icon size="40"><Folder /></el-icon>
          </div>
          <div class="stat-content">
            <h3>{{ stats.categoryCount || 0 }}</h3>
            <p>分类数</p>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-icon tags">
            <el-icon size="40"><CollectionTag /></el-icon>
          </div>
          <div class="stat-content">
            <h3>{{ stats.tagCount || 0 }}</h3>
            <p>标签数</p>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-icon storage">
            <el-icon size="40"><Histogram /></el-icon>
          </div>
          <div class="stat-content">
            <h3>{{ stats.totalSizeFormatted || '0 B' }}</h3>
            <p>总存储</p>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <div class="card-container">
          <h3>快速操作</h3>
          <div class="quick-actions">
            <el-button type="primary" @click="$router.push('/upload')">
              <el-icon><Upload /></el-icon>
              上传文件
            </el-button>
            <el-button @click="$router.push('/categories')">
              <el-icon><Folder /></el-icon>
              管理分类
            </el-button>
            <el-button @click="$router.push('/tags')">
              <el-icon><CollectionTag /></el-icon>
              管理标签
            </el-button>
          </div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="card-container">
          <h3>最近上传</h3>
          <div class="recent-files">
            <div v-if="recentFiles.length === 0" class="no-data">
              暂无文件
            </div>
            <div
              v-for="file in recentFiles"
              :key="file.id"
              class="file-item-mini"
            >
              <el-icon class="file-icon-mini"><Document /></el-icon>
              <span class="file-name-mini">{{ file.original_name }}</span>
              <span class="file-time">{{ formatTime(file.created_at) }}</span>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/api'

const stats = ref({})
const recentFiles = ref([])

const fetchStats = async () => {
  try {
    const response = await api.get('/admin/stats')
    stats.value = response.data.data
  } catch (error) {
    console.error('获取统计信息失败:', error)
  }
}

const fetchRecentFiles = async () => {
  try {
    const response = await api.get('/files?page=1&page_size=5')
    recentFiles.value = response.data.data.list || []
  } catch (error) {
    console.error('获取最近文件失败:', error)
  }
}

const formatTime = (time) => {
  return new Date(time).toLocaleDateString()
}

onMounted(() => {
  fetchStats()
  fetchRecentFiles()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
}

.stat-icon.files { background: #ecf5ff; color: #409eff; }
.stat-icon.categories { background: #f0f9ff; color: #67c23a; }
.stat-icon.tags { background: #fdf6ec; color: #e6a23c; }
.stat-icon.storage { background: #fef0f0; color: #f56c6c; }

.stat-content h3 {
  margin: 0 0 4px 0;
  font-size: 24px;
  font-weight: 600;
}

.stat-content p {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.quick-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.recent-files {
  max-height: 300px;
  overflow-y: auto;
}

.file-item-mini {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.file-item-mini:last-child {
  border-bottom: none;
}

.file-icon-mini {
  margin-right: 8px;
  color: #909399;
}

.file-name-mini {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-time {
  color: #909399;
  font-size: 12px;
}

.no-data {
  text-align: center;
  color: #909399;
  padding: 40px 0;
}
</style>