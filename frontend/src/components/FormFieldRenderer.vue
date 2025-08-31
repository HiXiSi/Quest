<template>
  <div class="form-field" :style="{ width: field.width || '100%' }">
    <el-form-item
      :label="field.label"
      :required="field.required"
      :prop="preview ? undefined : field.name"
    >
      <!-- 文本输入 -->
      <el-input
        v-if="field.type === 'text' && (!field.inputType || field.inputType === 'single')"
        v-model="fieldValue"
        :placeholder="field.placeholder"
        :disabled="preview"
        :maxlength="field.max_length"
        :minlength="field.min_length"
        show-word-limit
      />
      
      <!-- 多行文本输入 -->
      <el-input
        v-else-if="field.type === 'text' && field.inputType === 'multi'"
        v-model="fieldValue"
        type="textarea"
        :placeholder="field.placeholder"
        :disabled="preview"
        :maxlength="field.max_length"
        :minlength="field.min_length"
        :rows="field.rows || 4"
        show-word-limit
      />
      
      <!-- 数字 -->
      <el-input-number
        v-else-if="field.type === 'number'"
        v-model="fieldValue"
        :placeholder="field.placeholder"
        :disabled="preview"
        :min="field.min"
        :max="field.max"
        :step="field.step || 1"
        style="width: 100%"
      />
      
      <!-- 邮箱 -->
      <el-input
        v-else-if="field.type === 'email'"
        v-model="fieldValue"
        type="email"
        :placeholder="field.placeholder"
        :disabled="preview"
        :maxlength="field.max_length"
        :minlength="field.min_length"
      />
      
      <!-- 网址 -->
      <el-input
        v-else-if="field.type === 'url'"
        v-model="fieldValue"
        type="url"
        :placeholder="field.placeholder"
        :disabled="preview"
        :maxlength="field.max_length"
        :minlength="field.min_length"
      />
      
      <!-- 密码 -->
      <el-input
        v-else-if="field.type === 'password'"
        v-model="fieldValue"
        type="password"
        :placeholder="field.placeholder"
        :disabled="preview"
        :maxlength="field.max_length"
        :minlength="field.min_length"
        show-password
      />
      
      <!-- 下拉选择 -->
      <el-select
        v-else-if="field.type === 'select'"
        v-model="fieldValue"
        :placeholder="field.placeholder"
        :disabled="preview"
        style="width: 100%"
      >
        <el-option
          v-for="option in field.options"
          :key="option.value"
          :label="option.label"
          :value="option.value"
        />
      </el-select>
      
      <!-- 单选按钮 -->
      <el-radio-group
        v-else-if="field.type === 'radio'"
        v-model="fieldValue"
        :disabled="preview"
      >
        <el-radio
          v-for="option in field.options"
          :key="option.value"
          :label="option.value"
        >
          {{ option.label }}
        </el-radio>
      </el-radio-group>
      
      <!-- 多选框 -->
      <el-checkbox-group
        v-else-if="field.type === 'checkbox'"
        v-model="fieldValue"
        :disabled="preview"
      >
        <el-checkbox
          v-for="option in field.options"
          :key="option.value"
          :label="option.value"
        >
          {{ option.label }}
        </el-checkbox>
      </el-checkbox-group>
      
      <!-- 日期 -->
      <el-date-picker
        v-else-if="field.type === 'date'"
        v-model="fieldValue"
        type="date"
        :placeholder="field.placeholder"
        :disabled="preview"
        style="width: 100%"
      />
      
      <!-- 日期时间 -->
      <el-date-picker
        v-else-if="field.type === 'datetime'"
        v-model="fieldValue"
        type="datetime"
        :placeholder="field.placeholder"
        :disabled="preview"
        style="width: 100%"
      />
      
      <!-- 时间 -->
      <el-time-picker
        v-else-if="field.type === 'time'"
        v-model="fieldValue"
        :placeholder="field.placeholder"
        :disabled="preview"
        style="width: 100%"
      />
      
      <!-- 文件上传 -->
      <el-upload
        v-else-if="field.type === 'file'"
        class="upload-demo"
        :disabled="preview"
        :auto-upload="false"
        :show-file-list="true"
        :limit="1"
      >
        <el-button type="primary" :disabled="preview">
          <el-icon><Upload /></el-icon>
          选择文件
        </el-button>
        <template #tip>
          <div class="el-upload__tip">{{ field.placeholder || '请选择文件' }}</div>
        </template>
      </el-upload>
      
      <!-- 开关 -->
      <el-switch
        v-else-if="field.type === 'switch'"
        v-model="fieldValue"
        :disabled="preview"
      />
      
      <!-- 评分 -->
      <el-rate
        v-else-if="field.type === 'rate'"
        v-model="fieldValue"
        :disabled="preview"
      />
      
      <!-- 未知类型 -->
      <div v-else class="unknown-field">
        <el-alert
          title="未知字段类型"
          :description="`字段类型 '${field.type}' 暂不支持`"
          type="warning"
          show-icon
          :closable="false"
        />
      </div>
    </el-form-item>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { Upload } from '@element-plus/icons-vue'

const props = defineProps({
  field: {
    type: Object,
    required: true
  },
  preview: {
    type: Boolean,
    default: false
  },
  modelValue: {
    type: [String, Number, Boolean, Array, Date],
    default: null
  }
})

const emit = defineEmits(['update:modelValue'])

// 字段值
const fieldValue = computed({
  get() {
    if (props.modelValue !== null && props.modelValue !== undefined) {
      return props.modelValue
    }
    
    // 根据字段类型返回默认值
    switch (props.field.type) {
      case 'checkbox':
        return []
      case 'number':
        return props.field.default_value ? Number(props.field.default_value) : null
      case 'switch':
        return props.field.default_value === 'true' || props.field.default_value === true
      case 'rate':
        return props.field.default_value ? Number(props.field.default_value) : 0
      default:
        return props.field.default_value || ''
    }
  },
  set(value) {
    emit('update:modelValue', value)
  }
})

// 监听字段配置变化，更新默认值
watch(
  () => props.field.default_value,
  (newValue) => {
    if (!props.modelValue && newValue) {
      fieldValue.value = newValue
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.form-field {
  margin-bottom: 16px;

  .unknown-field {
    margin: 8px 0;
  }

  .upload-demo {
    .el-upload__tip {
      font-size: 12px;
      color: #606266;
      margin-top: 4px;
    }
  }
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-form-item) {
  margin-bottom: 0;
}

:deep(.el-radio) {
  margin-right: 16px;
}

:deep(.el-checkbox) {
  margin-right: 16px;
}
</style>