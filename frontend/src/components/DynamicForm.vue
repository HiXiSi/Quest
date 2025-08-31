<template>
  <div class="dynamic-form">
    <div v-if="schema" class="form-container">
      <!-- 表单头部信息 -->
      <div class="form-header">
        <h2>{{ schema.name }}</h2>
        <p v-if="schema.description" class="form-description">{{ schema.description }}</p>
      </div>
      
      <!-- 表单内容 -->
      <el-form
        ref="dynamicFormRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
        @submit.prevent="handleSubmit"
      >
        <div
          v-for="field in schema.fields"
          :key="field.name"
          class="form-field-wrapper"
        >
          <FormFieldRenderer
            :field="field"
            v-model="formData[field.name]"
            :preview="false"
          />
        </div>
        
        <!-- 表单操作按钮 -->
        <div class="form-actions">
          <el-button
            v-if="showResetButton"
            @click="resetForm"
            :disabled="submitting"
          >
            重置
          </el-button>
          <el-button
            type="primary"
            @click="handleSubmit"
            :loading="submitting"
            :disabled="!isFormValid"
          >
            {{ submitButtonText }}
          </el-button>
        </div>
      </el-form>
    </div>
    
    <div v-else class="no-schema">
      <el-empty description="表单配置加载失败" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import FormFieldRenderer from './FormFieldRenderer.vue'

const props = defineProps({
  // 表单结构配置
  schema: {
    type: Object,
    default: null
  },
  // 初始数据
  initialData: {
    type: Object,
    default: () => ({})
  },
  // 是否显示重置按钮
  showResetButton: {
    type: Boolean,
    default: true
  },
  // 提交按钮文本
  submitButtonText: {
    type: String,
    default: '提交'
  },
  // 是否自动验证
  autoValidate: {
    type: Boolean,
    default: true
  },
  // 是否禁用表单
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit', 'reset', 'change', 'validate'])

const dynamicFormRef = ref()
const formData = ref({})
const submitting = ref(false)

// 表单验证规则
const formRules = computed(() => {
  const rules = {}
  if (props.schema?.fields) {
    props.schema.fields.forEach(field => {
      const fieldRules = []
      
      // 必填验证
      if (field.required) {
        fieldRules.push({
          required: true,
          message: `请输入${field.label}`,
          trigger: ['blur', 'change']
        })
      }
      
      // 字符长度验证
      if (field.min_length || field.max_length) {
        if (field.type === 'text' || field.type === 'email' || field.type === 'url') {
          fieldRules.push({
            min: field.min_length || 0,
            max: field.max_length || Infinity,
            message: `长度应在 ${field.min_length || 0} 到 ${field.max_length || '无限'} 个字符之间`,
            trigger: ['blur', 'change']
          })
        }
      }
      
      // 数字范围验证
      if (field.type === 'number' && (field.min !== null || field.max !== null)) {
        fieldRules.push({
          type: 'number',
          min: field.min !== null ? field.min : -Infinity,
          max: field.max !== null ? field.max : Infinity,
          message: `数值应在 ${field.min !== null ? field.min : '无限小'} 到 ${field.max !== null ? field.max : '无限大'} 之间`,
          trigger: ['blur', 'change']
        })
      }
      
      // 邮箱格式验证
      if (field.type === 'email') {
        fieldRules.push({
          type: 'email',
          message: '请输入正确的邮箱格式',
          trigger: ['blur', 'change']
        })
      }
      
      // URL格式验证
      if (field.type === 'url') {
        fieldRules.push({
          type: 'url',
          message: '请输入正确的网址格式',
          trigger: ['blur', 'change']
        })
      }
      
      if (fieldRules.length > 0) {
        rules[field.name] = fieldRules
      }
    })
  }
  return rules
})

// 表单是否有效
const isFormValid = computed(() => {
  if (!props.autoValidate || !props.schema?.fields) return true
  
  return props.schema.fields.every(field => {
    if (!field.required) return true
    
    const value = formData.value[field.name]
    
    // 检查必填字段
    if (field.required) {
      if (value === null || value === undefined || value === '') {
        return false
      }
      
      // 数组类型（如多选框）需要检查长度
      if (Array.isArray(value) && value.length === 0) {
        return false
      }
    }
    
    return true
  })
})

// 初始化表单数据
const initFormData = () => {
  const data = {}
  
  if (props.schema?.fields) {
    props.schema.fields.forEach(field => {
      // 优先使用传入的初始数据
      if (props.initialData && props.initialData[field.name] !== undefined) {
        data[field.name] = props.initialData[field.name]
      } else if (field.default_value !== undefined && field.default_value !== '') {
        // 使用字段配置的默认值
        switch (field.type) {
          case 'number':
            data[field.name] = Number(field.default_value) || 0
            break
          case 'switch':
            data[field.name] = field.default_value === 'true' || field.default_value === true
            break
          case 'checkbox':
            data[field.name] = Array.isArray(field.default_value) ? field.default_value : []
            break
          case 'rate':
            data[field.name] = Number(field.default_value) || 0
            break
          default:
            data[field.name] = field.default_value
        }
      } else {
        // 使用类型默认值
        switch (field.type) {
          case 'checkbox':
            data[field.name] = []
            break
          case 'number':
          case 'rate':
            data[field.name] = 0
            break
          case 'switch':
            data[field.name] = false
            break
          default:
            data[field.name] = ''
        }
      }
    })
  }
  
  formData.value = data
}

// 重置表单
const resetForm = () => {
  if (dynamicFormRef.value) {
    dynamicFormRef.value.resetFields()
  }
  initFormData()
  emit('reset', { ...formData.value })
  ElMessage.success('表单已重置')
}

// 验证表单
const validateForm = async () => {
  if (!dynamicFormRef.value) return false
  
  try {
    await dynamicFormRef.value.validate()
    return true
  } catch (error) {
    return false
  }
}

// 提交表单
const handleSubmit = async () => {
  if (props.disabled) return
  
  const isValid = await validateForm()
  if (!isValid) {
    ElMessage.warning('请完善表单信息')
    return
  }
  
  submitting.value = true
  
  try {
    // 准备提交数据
    const submitData = { ...formData.value }
    
    // 处理特殊字段类型的数据格式
    if (props.schema?.fields) {
      props.schema.fields.forEach(field => {
        const value = submitData[field.name]
        
        // 处理日期类型
        if ((field.type === 'date' || field.type === 'datetime' || field.type === 'time') && value) {
          submitData[field.name] = value.toISOString ? value.toISOString() : value
        }
        
        // 处理文件上传类型（这里简化处理，实际项目中需要处理文件上传）
        if (field.type === 'file' && value) {
          // TODO: 处理文件上传逻辑
        }
      })
    }
    
    emit('submit', {
      data: submitData,
      schema: props.schema
    })
  } catch (error) {
    console.error('表单提交失败:', error)
    ElMessage.error('表单提交失败')
  } finally {
    submitting.value = false
  }
}

// 监听表单数据变化
watch(
  formData,
  (newData) => {
    emit('change', { ...newData })
    
    // 自动验证
    if (props.autoValidate && dynamicFormRef.value) {
      emit('validate', isFormValid.value)
    }
  },
  { deep: true }
)

// 监听schema变化，重新初始化表单
watch(
  () => props.schema,
  () => {
    initFormData()
  },
  { immediate: true }
)

// 监听初始数据变化
watch(
  () => props.initialData,
  () => {
    initFormData()
  },
  { deep: true }
)

// 暴露方法给父组件
defineExpose({
  validateForm,
  resetForm,
  getFormData: () => ({ ...formData.value }),
  setFormData: (data) => {
    formData.value = { ...formData.value, ...data }
  }
})

onMounted(() => {
  initFormData()
})
</script>

<style scoped>
.dynamic-form {
  .form-container {
    max-width: 800px;
    margin: 0 auto;

    .form-header {
      text-align: center;
      margin-bottom: 32px;
      padding-bottom: 16px;
      border-bottom: 1px solid #e4e7ed;

      h2 {
        margin: 0 0 8px 0;
        color: #303133;
        font-size: 24px;
      }

      .form-description {
        margin: 0;
        color: #606266;
        font-size: 14px;
        line-height: 1.5;
      }
    }

    .form-field-wrapper {
      margin-bottom: 24px;
    }

    .form-actions {
      text-align: center;
      margin-top: 32px;
      padding-top: 24px;
      border-top: 1px solid #e4e7ed;

      .el-button {
        margin: 0 8px;
        min-width: 100px;
      }
    }
  }

  .no-schema {
    text-align: center;
    padding: 60px 0;
  }
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #606266;
}

:deep(.el-form-item__error) {
  font-size: 12px;
}
</style>