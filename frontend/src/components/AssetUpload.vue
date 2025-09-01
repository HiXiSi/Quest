<template>
  <div class="asset-upload">
    <el-upload
      ref="uploadRef"
      :action="uploadUrl"
      :headers="uploadHeaders"
      :before-upload="beforeUpload"
      :on-success="handleSuccess"
      :on-error="handleError"
      :on-progress="handleProgress"
      :show-file-list="false"
      :auto-upload="true"
      accept="*/*"
    >
      <template #trigger>
        <el-button :icon="Upload" :disabled="disabled">
          选择文件
        </el-button>
      </template>
    </el-upload>
    
    <!-- 文件信息显示 -->
    <div v-if="selectedFile" class="file-info">
      <div class="file-item">
        <!-- 文件预览 -->
        <div v-if="isImageFile(selectedFile)" class="file-preview image-preview">
          <img :src="previewUrl" alt="Preview" />
        </div>
        <div v-else-if="isVideoFile(selectedFile)" class="file-preview video-preview">
          <video :src="previewUrl" muted></video>
          <div class="play-overlay">
            <el-icon class="play-icon"><VideoPlay /></el-icon>
          </div>
        </div>
        <el-icon v-else class="file-icon"><Document /></el-icon>
        
        <div class="file-details">
          <span class="file-name">{{ selectedFile.name }}</span>
          <span class="file-size">({{ formatFileSize(selectedFile.size) }})</span>
        </div>
        <el-button
          v-if="!uploading"
          @click="clearFile"
          :icon="Delete"
          size="small"
          type="danger"
          text
        />
      </div>
      
      <!-- 上传进度 -->
      <div v-if="uploading" class="upload-progress">
        <el-progress
          :percentage="uploadProgress"
          :status="uploadStatus"
        />
      </div>
      
      <!-- 上传按钮 -->
      <!-- 移除了手动上传按钮，因为现在是自动上传 -->
    </div>
    
    <!-- 已上传文件显示 -->
    <div v-if="fileUrl && !uploading" class="uploaded-file">
      <div class="success-info">
        <el-icon class="success-icon"><Check /></el-icon>
        <span>文件上传成功</span>
        <el-button @click="clearAll" size="small" text>重新选择</el-button>
      </div>
      
      <!-- 已上传文件预览 -->
      <div v-if="isImageUrl(fileUrl)" class="uploaded-preview image-preview">
        <img :src="absoluteFileUrl" alt="Uploaded file" />
      </div>
      <div v-else-if="isVideoUrl(fileUrl)" class="uploaded-preview video-preview">
        <video :src="absoluteFileUrl" muted></video>
        <div class="play-overlay">
          <el-icon class="play-icon"><VideoPlay /></el-icon>
        </div>
      </div>
      
      <div class="file-url">
        <el-input
          :model-value="absoluteFileUrl"
          readonly
          size="small"
        >
          <template #append>
            <el-button @click="copyUrl" size="small">复制</el-button>
          </template>
        </el-input>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload, Document, Delete, Check, VideoPlay } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { toAbsoluteUrl } from '@/utils/urlHelper'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const userStore = useUserStore()

const uploadRef = ref()
const selectedFile = ref(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadStatus = ref('')
const previewUrl = ref('')

// 计算属性
const fileUrl = computed({
  get() {
    return props.modelValue
  },
  set(value) {
    emit('update:modelValue', value)
  }
})

// 用于界面显示的绝对URL
const absoluteFileUrl = computed(() => {
  return toAbsoluteUrl(props.modelValue)
})

const uploadUrl = computed(() => {
  // 使用相对路径，让Vite代理处理
  return '/api/assets/upload'
})

const uploadHeaders = computed(() => {
  return {
    'Authorization': `Bearer ${userStore.token}`
  }
})

// 判断是否为图片文件
const isImageFile = (file) => {
  return file.type.startsWith('image/')
}

// 判断是否为视频文件
const isVideoFile = (file) => {
  return file.type.startsWith('video/')
}

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

// 文件上传前检查
const beforeUpload = (file) => {
  // 文件大小限制 50MB
  if (file.size > 50 * 1024 * 1024) {
    ElMessage.error('文件大小不能超过50MB')
    return false
  }
  
  selectedFile.value = file
  
  // 生成预览URL
  if (isImageFile(file) || isVideoFile(file)) {
    previewUrl.value = URL.createObjectURL(file)
  }
  
  // 返回true以允许自动上传
  return true
}

// 上传成功
const handleSuccess = (response) => {
  uploading.value = false
  uploadProgress.value = 100
  uploadStatus.value = 'success'
  
  if (response.code === 0) {
    // 保持原始URL（相对路径）
    fileUrl.value = response.data.url
    selectedFile.value = null
    previewUrl.value = ''
    ElMessage.success('文件上传成功')
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

// 上传失败
const handleError = (error) => {
  uploading.value = false
  uploadStatus.value = 'exception'
  ElMessage.error('文件上传失败')
  console.error('Upload error:', error)
}

// 上传进度
const handleProgress = (event) => {
  uploading.value = true
  uploadProgress.value = Math.round(event.percent)
}

// 格式化文件大小
const formatFileSize = (size) => {
  if (size < 1024) {
    return size + ' B'
  } else if (size < 1024 * 1024) {
    return (size / 1024).toFixed(1) + ' KB'
  } else {
    return (size / (1024 * 1024)).toFixed(1) + ' MB'
  }
}

// 清除选中的文件
const clearFile = () => {
  selectedFile.value = null
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = ''
  }
  uploadRef.value.clearFiles()
}

// 清除所有
const clearAll = () => {
  selectedFile.value = null
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = ''
  }
  fileUrl.value = ''
  uploadRef.value.clearFiles()
}

// 复制URL
const copyUrl = async () => {
  try {
    await navigator.clipboard.writeText(absoluteFileUrl.value)
    ElMessage.success('URL已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 组件卸载前清理
onBeforeUnmount(() => {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
  }
})
</script>

<style scoped>
.asset-upload {
  .file-info {
    margin-top: 12px;
    padding: 12px;
    border: 1px solid #e4e7ed;
    border-radius: 4px;
    background: #f9f9f9;

    .file-item {
      display: flex;
      align-items: center;
      gap: 8px;

      .file-preview {
        width: 40px;
        height: 40px;
        border-radius: 4px;
        overflow: hidden;
        position: relative;
        flex-shrink: 0;

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
            opacity: 0;
            transition: opacity 0.3s;

            .play-icon {
              font-size: 20px;
              color: white;
            }
          }

          &:hover .play-overlay {
            opacity: 1;
          }
        }
      }

      .file-icon {
        font-size: 24px;
        color: #606266;
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
      }

      .file-details {
        flex: 1;
        min-width: 0;

        .file-name {
          font-size: 14px;
          color: #303133;
          display: block;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .file-size {
          font-size: 12px;
          color: #909399;
        }
      }
    }

    .upload-progress {
      margin-top: 12px;
    }

    .upload-actions {
      margin-top: 12px;
      text-align: right;
    }
  }

  .uploaded-file {
    margin-top: 12px;
    padding: 12px;
    border: 1px solid #67c23a;
    border-radius: 4px;
    background: #f0f9ff;

    .success-info {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 8px;

      .success-icon {
        color: #67c23a;
      }

      span {
        font-size: 14px;
        color: #67c23a;
      }
    }

    .uploaded-preview {
      width: 100%;
      height: 200px;
      border-radius: 4px;
      overflow: hidden;
      margin-bottom: 12px;
      position: relative;

      &.image-preview {
        img {
          width: 100%;
          height: 100%;
          object-fit: contain;
        }
      }

      &.video-preview {
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
          opacity: 0;
          transition: opacity 0.3s;

          .play-icon {
            font-size: 40px;
            color: white;
          }
        }

        &:hover .play-overlay {
          opacity: 1;
        }
      }
    }

    .file-url {
      .el-input {
        font-size: 12px;
      }
    }
  }
}
</style>