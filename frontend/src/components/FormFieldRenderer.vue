<template>
  <div class="form-field" :style="{ width: field.width || '100%' }">
    <el-form-item
      :label="field.label"
      :required="field.required"
      :prop="preview ? undefined : field.name"
    >
      <!-- 唯一ID字段 -->
      <el-input
        v-if="field.type === 'unique_id'"
        v-model="fieldValue"
        :placeholder="field.placeholder"
        :disabled="true"
        readonly
      />
      
      <!-- 整数字段 -->
      <el-input-number
        v-else-if="field.type === 'integer'"
        v-model="fieldValue"
        :placeholder="field.placeholder"
        :disabled="preview"
        :min="field.min_value"
        :max="field.max_value"
        :step="1"
        style="width: 100%"
      />
      
      <!-- 浮点数字段 -->
      <el-input-number
        v-else-if="field.type === 'float'"
        v-model="fieldValue"
        :placeholder="field.placeholder"
        :disabled="preview"
        :min="field.min_value"
        :max="field.max_value"
        :precision="field.precision || 2"
        :step="0.01"
        style="width: 100%"
      />
      
      <!-- 字符串字段 -->
      <el-input
        v-else-if="field.type === 'string' && field.input_type === 'input'"
        v-model="fieldValue"
        :type="getStringInputType(field)"
        :placeholder="field.placeholder"
        :disabled="preview"
        :maxlength="field.max_length"
        :minlength="field.min_length"
        show-word-limit
      />
      
      <!-- 密码字段 -->
      <el-input
        v-else-if="field.type === 'string' && field.input_type === 'password'"
        v-model="fieldValue"
        type="password"
        :placeholder="field.placeholder"
        :disabled="preview"
        :maxlength="field.max_length"
        :minlength="field.min_length"
        show-password
        show-word-limit
      />
      
      <!-- 多行文本输入 -->
      <el-input
        v-else-if="field.type === 'string' && field.input_type === 'textarea'"
        v-model="fieldValue"
        type="textarea"
        :placeholder="field.placeholder"
        :disabled="preview"
        :maxlength="field.max_length"
        :minlength="field.min_length"
        :rows="field.textarea_rows || 4"
        show-word-limit
      />
      
      <!-- 文件上传 -->
      <AssetUpload
        v-else-if="field.type === 'string' && field.input_type === 'file'"
        v-model="fieldValue"
        :disabled="preview"
      />
      
      <!-- 布尔字段 -->
      <el-switch
        v-else-if="field.type === 'boolean' && field.input_type === 'switch'"
        v-model="fieldValue"
        :disabled="preview"
      />
      
      <!-- 布尔单选按钮 -->
      <el-radio-group
        v-else-if="field.type === 'boolean' && field.input_type === 'radio'"
        v-model="fieldValue"
        :disabled="preview"
      >
        <el-radio :label="true">是</el-radio>
        <el-radio :label="false">否</el-radio>
      </el-radio-group>
      
      <!-- 布尔复选框 -->
      <el-checkbox
        v-else-if="field.type === 'boolean' && field.input_type === 'checkbox'"
        v-model="fieldValue"
        :disabled="preview"
      >
        {{ field.label }}
      </el-checkbox>
      
      <!-- 时间字段 - 日期时间 -->
      <el-date-picker
        v-else-if="field.type === 'datetime' && field.input_type === 'datetime'"
        v-model="fieldValue"
        type="datetime"
        :placeholder="field.placeholder"
        :disabled="preview"
        style="width: 100%"
      />
      
      <!-- 时间字段 - 仅日期 -->
      <el-date-picker
        v-else-if="field.type === 'datetime' && field.input_type === 'date'"
        v-model="fieldValue"
        type="date"
        :placeholder="field.placeholder"
        :disabled="preview"
        style="width: 100%"
      />
      
      <!-- 时间字段 - 仅时间 -->
      <el-time-picker
        v-else-if="field.type === 'datetime' && field.input_type === 'time'"
        v-model="fieldValue"
        :placeholder="field.placeholder"
        :disabled="preview"
        style="width: 100%"
      />
      
      <!-- 单选枚举 - 下拉选择 -->
      <el-select
        v-else-if="field.type === 'single_enum' && field.input_type === 'select'"
        v-model="fieldValue"
        :placeholder="field.placeholder"
        :disabled="preview"
        style="width: 100%"
      >
        <el-option
          v-for="option in field.enum_options || []"
          :key="option.value"
          :label="option.label"
          :value="option.value"
        />
      </el-select>
      
      <!-- 单选枚举 - 单选按钮 -->
      <el-radio-group
        v-else-if="field.type === 'single_enum' && field.input_type === 'radio'"
        v-model="fieldValue"
        :disabled="preview"
      >
        <el-radio
          v-for="option in field.enum_options || []"
          :key="option.value"
          :label="option.value"
        >
          {{ option.label }}
        </el-radio>
      </el-radio-group>
      
      <!-- 多选枚举 - 多选框 -->
      <el-checkbox-group
        v-else-if="field.type === 'multi_enum' && field.input_type === 'checkbox'"
        v-model="fieldValue"
        :disabled="preview"
      >
        <el-checkbox
          v-for="option in field.enum_options || []"
          :key="option.value"
          :label="option.value"
        >
          {{ option.label }}
        </el-checkbox>
      </el-checkbox-group>
      
      <!-- 多选枚举 - 多选下拉 -->
      <el-select
        v-else-if="field.type === 'multi_enum' && field.input_type === 'multi-select'"
        v-model="fieldValue"
        :placeholder="field.placeholder"
        :disabled="preview"
        multiple
        style="width: 100%"
      >
        <el-option
          v-for="option in field.enum_options || []"
          :key="option.value"
          :label="option.label"
          :value="option.value"
        />
      </el-select>
      
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
import AssetUpload from './AssetUpload.vue'

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

// 获取字符串输入类型
const getStringInputType = (field) => {
  switch (field.format) {
    case 'email':
      return 'email'
    case 'url':
      return 'url'
    case 'phone':
      return 'tel'
    default:
      return 'text'
  }
}

// 字段值
const fieldValue = computed({
  get() {
    if (props.modelValue !== null && props.modelValue !== undefined) {
      return props.modelValue
    }
    
    // 根据字段类型返回默认值
    switch (props.field.type) {
      case 'multi_enum':
        return []
      case 'integer':
      case 'float':
        return props.field.default_value ? Number(props.field.default_value) : null
      case 'boolean':
        return props.field.default_value === 'true' || props.field.default_value === true
      case 'unique_id':
        return ''
      case 'datetime':
        return props.field.default_value || null
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
  font-size: 14px;
}

:deep(.el-form-item) {
  margin-bottom: 16px;
}

:deep(.el-radio) {
  margin-right: 12px;
}

:deep(.el-checkbox) {
  margin-right: 12px;
}

/* 紧凑布局优化 */
:deep(.el-input),
:deep(.el-select),
:deep(.el-date-picker),
:deep(.el-time-picker) {
  width: 100%;
}

:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-radio-group) {
  width: 100%;
}

:deep(.el-checkbox-group) {
  width: 100%;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
</style>