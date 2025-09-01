<template>
  <div class="page-container">
    <div class="designer-container">
      <!-- 左侧字段面板 -->
      <div class="fields-panel">
        <h3>字段类型</h3>
        <div class="field-categories">
          <div class="field-items">
            <div
              v-for="field in fieldTypes"
              :key="field.type"
              class="field-item"
              draggable="true"
              @dragstart="handleDragStart($event, field)"
            >
              <el-icon>
                <component :is="field.icon" />
              </el-icon>
              <span>{{ field.label }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 中间设计区域 -->
      <div class="design-area">
        <div class="design-header">
          <div class="design-title">
            <el-input
              v-model="formSchema.name"
              placeholder="请输入数据表名称"
              size="large"
              style="font-size: 18px; font-weight: bold;"
            />
          </div>
          <div class="design-description">
            <el-input
              v-model="formSchema.description"
              placeholder="请输入数据表描述（可选）"
              type="textarea"
              :rows="2"
            />
          </div>
        </div>
        
        <!-- 设计操作按钮 -->
        <div class="design-actions">
          <el-button @click="previewForm">
            <el-icon><View /></el-icon>
            预览表单
          </el-button>
          <el-button @click="clearForm" type="warning">
            <el-icon><Delete /></el-icon>
            清空字段
          </el-button>
          <el-button @click="saveForm" type="primary" :loading="saving">
            <el-icon><Check /></el-icon>
            保存表结构
          </el-button>
        </div>
        
        <div class="fields-canvas">
          <div
            class="drop-zone"
            @dragover.prevent
            @drop="handleDrop"
            :class="{ 'drag-over': isDragOver }"
            @dragenter.prevent="isDragOver = true"
            @dragleave.prevent="isDragOver = false"
          >
            <div v-if="!formSchema.fields || formSchema.fields.length === 0" class="empty-hint">
              <el-icon size="48" color="#c0c4cc">
                <Plus />
              </el-icon>
              <p>从左侧拖拽字段类型到此处来创建字段</p>
            </div>
            
            <!-- 字段列表 -->
            <div class="fields-list">
              <div
                v-for="(field, index) in formSchema.fields || []"
                :key="field.id"
                class="field-row"
                :class="{ active: selectedFieldIndex === index }"
                @click="selectField(index)"
              >
                <div class="field-info">
                  <div class="field-name">{{ field.name }}</div>
                  <div class="field-type">{{ getFieldTypeLabel(field.type) }}</div>
                </div>
                <div class="field-controls">
                  <el-button-group size="small">
                    <el-button @click.stop="moveFieldUp(index)" :disabled="index === 0">
                      <el-icon><ArrowUp /></el-icon>
                    </el-button>
                    <el-button @click.stop="moveFieldDown(index)" :disabled="index === (formSchema.fields || []).length - 1">
                      <el-icon><ArrowDown /></el-icon>
                    </el-button>
                    <el-button @click.stop="duplicateField(index)">
                      <el-icon><CopyDocument /></el-icon>
                    </el-button>
                    <el-button @click.stop="removeField(index)" type="danger">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </el-button-group>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 右侧属性面板 -->
      <div class="properties-panel">
        <h3>字段属性</h3>
        <div v-if="selectedField" class="property-form">
          <el-form label-width="100px" size="small">
            <!-- 1. 基础属性 -->
            <h4>基础属性</h4>
            <el-form-item label="字段名">
              <el-input v-model="selectedField.name" placeholder="字段名（英文）" />
            </el-form-item>
            <el-form-item label="字段标签">
              <el-input v-model="selectedField.label" placeholder="界面显示名称（中文）" />
            </el-form-item>
            <el-form-item label="默认值">
              <el-input v-model="selectedField.default_value" placeholder="默认值" />
            </el-form-item>
            <el-form-item>
              <el-checkbox v-model="selectedField.required">必填字段</el-checkbox>
            </el-form-item>
            
            <!-- 2. 格式属性 -->
            <h4>格式属性</h4>
            
            <!-- 唯一ID类型的格式 -->
            <template v-if="selectedField.type === 'unique_id'">
              <el-form-item label="ID类型">
                <el-select v-model="selectedField.id_type">
                  <el-option label="自增整数" value="auto_increment" />
                  <el-option label="UUID字符串" value="uuid" />
                </el-select>
              </el-form-item>
            </template>
            
            <!-- 字符串类型的格式 -->
            <template v-if="selectedField.type === 'string'">
              <el-form-item label="字符串格式">
                <el-select v-model="selectedField.format">
                  <el-option label="普通文本" value="text" />
                  <el-option label="邮箱地址" value="email" />
                  <el-option label="手机号码" value="phone" />
                  <el-option label="网址URL" value="url" />
                </el-select>
              </el-form-item>
              <el-form-item label="最小长度">
                <el-input-number v-model="selectedField.min_length" :min="0" :controls="false" />
              </el-form-item>
              <el-form-item label="最大长度">
                <el-input-number v-model="selectedField.max_length" :min="1" :controls="false" />
              </el-form-item>
            </template>
            
            <!-- 数值类型的格式 -->
            <template v-if="selectedField.type === 'integer' || selectedField.type === 'float'">
              <el-form-item v-if="selectedField.type === 'float'" label="小数位数">
                <el-input-number v-model="selectedField.precision" :min="0" :max="10" :controls="false" />
              </el-form-item>
              <el-form-item label="最小值">
                <el-input-number v-model="selectedField.min_value" :controls="false" />
              </el-form-item>
              <el-form-item label="最大值">
                <el-input-number v-model="selectedField.max_value" :controls="false" />
              </el-form-item>
            </template>
            
            <!-- 时间类型的格式 -->
            <template v-if="selectedField.type === 'datetime'">
              <el-form-item label="格式化">
                <el-select v-model="selectedField.time_format">
                  <el-option label="Date对象" value="date_object" />
                  <el-option label="YYYY-MM-DD HH:mm:ss" value="datetime" />
                  <el-option label="YYYY-MM-DD" value="date" />
                  <el-option label="HH:mm:ss" value="time" />
                </el-select>
              </el-form-item>
            </template>
            
            <!-- 枚举类型的选项 -->
            <template v-if="selectedField.type === 'single_enum' || selectedField.type === 'multi_enum'">
              <el-form-item label="枚举选项">
                <div class="enum-options-editor">
                  <div
                    v-for="(option, index) in selectedField.enum_options || []"
                    :key="index"
                    class="enum-option-item"
                  >
                    <el-input
                      v-model="option.value"
                      placeholder="英文值"
                      size="small"
                      style="width: 120px;"
                    />
                    <el-input
                      v-model="option.label"
                      placeholder="中文显示"
                      size="small"
                      style="width: 120px;"
                    />
                    <el-button
                      @click="removeEnumOption(index)"
                      size="small"
                      type="danger"
                      text
                    >
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </div>
                  <el-button @click="addEnumOption" size="small" type="primary" text>
                    <el-icon><Plus /></el-icon>
                    添加选项
                  </el-button>
                </div>
              </el-form-item>
            </template>
            
            <!-- 3. 表单形式 -->
            <h4>表单形式</h4>
            
            <!-- 唯一ID的表单形式 -->
            <template v-if="selectedField.type === 'unique_id'">
              <el-form-item label="表单展示">
                <el-select v-model="selectedField.input_type">
                  <el-option label="隐藏字段" value="hidden" />
                  <el-option label="只读展示" value="readonly" />
                </el-select>
              </el-form-item>
            </template>
            
            <!-- 字符串的表单形式 -->
            <template v-if="selectedField.type === 'string'">
              <el-form-item label="输入方式">
                <el-select v-model="selectedField.input_type">
                  <el-option label="单行输入" value="input" />
                  <el-option label="密码输入" value="password" />
                  <el-option label="多行输入" value="textarea" />
                  <el-option label="文件上传" value="file" />
                </el-select>
              </el-form-item>
              <el-form-item v-if="selectedField.input_type === 'textarea'" label="文本行数">
                <el-input-number v-model="selectedField.textarea_rows" :min="2" :max="10" :controls="false" />
              </el-form-item>
            </template>
            
            <!-- 时间的表单形式 -->
            <template v-if="selectedField.type === 'datetime'">
              <el-form-item v-if="selectedField.update_mode === 'user_input'" label="选择方式">
                <el-select v-model="selectedField.input_type">
                  <el-option label="日期时间" value="datetime" />
                  <el-option label="仅日期" value="date" />
                  <el-option label="仅时间" value="time" />
                </el-select>
              </el-form-item>
              <el-form-item v-else label="显示方式">
                <el-select v-model="selectedField.input_type">
                  <el-option label="日期时间" value="datetime" />
                  <el-option label="仅日期" value="date" />
                  <el-option label="仅时间" value="time" />
                </el-select>
              </el-form-item>
            </template>
            
            <!-- 枚举的表单形式 -->
            <template v-if="selectedField.type === 'single_enum'">
              <el-form-item label="选择方式">
                <el-select v-model="selectedField.input_type">
                  <el-option label="下拉选择" value="select" />
                  <el-option label="单选按钮" value="radio" />
                </el-select>
              </el-form-item>
            </template>
            
            <template v-if="selectedField.type === 'multi_enum'">
              <el-form-item label="选择方式">
                <el-select v-model="selectedField.input_type">
                  <el-option label="多选框" value="checkbox" />
                  <el-option label="多选下拉" value="multi-select" />
                </el-select>
              </el-form-item>
            </template>
            
            <!-- 布尔的表单形式 -->
            <template v-if="selectedField.type === 'boolean'">
              <el-form-item label="选择方式">
                <el-select v-model="selectedField.input_type">
                  <el-option label="开关" value="switch" />
                  <el-option label="单选按钮" value="radio" />
                  <el-option label="复选框" value="checkbox" />
                </el-select>
              </el-form-item>
            </template>
            
            <el-form-item label="占位符">
              <el-input v-model="selectedField.placeholder" placeholder="输入提示文本" />
            </el-form-item>
          </el-form>
        </div>
        <div v-else class="no-selection">
          <p>请选择一个字段来编辑属性</p>
        </div>
      </div>
    </div>
    
    <!-- 表单预览对话框 -->
    <el-dialog
      v-model="previewVisible"
      title="表单预览"
      width="800px"
      :close-on-click-modal="false"
    >
      <div class="form-preview">
        <h2>{{ formSchema.name || '未命名表单' }}</h2>
        <p v-if="formSchema.description">{{ formSchema.description }}</p>
        
        <el-form label-width="120px">
          <div
            v-for="field in formSchema.fields"
            :key="field.id"
          >
            <FormFieldRenderer :field="field" :preview="false" />
          </div>
        </el-form>
      </div>
      
      <template #footer>
        <el-button @click="previewVisible = false">关闭预览</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, ArrowUp, ArrowDown, CopyDocument, Delete, View, Check,
  EditPen, Calendar, Select, Document, List, Upload,
  Switch, StarFilled, Timer, Key
} from '@element-plus/icons-vue'
import api from '@/utils/api'
import FormFieldRenderer from '@/components/FormFieldRenderer.vue'

const route = useRoute()
const router = useRouter()

// 当前编辑的表单结构
const formSchema = ref({
  id: null,
  name: '',
  description: '',
  fields: []
})

// 界面状态
const selectedFieldIndex = ref(-1)
const isDragOver = ref(false)
const previewVisible = ref(false)
const saving = ref(false)

// 可拖拽的字段类型（基于数据库字段）
const fieldTypes = ref([
  { type: 'unique_id', label: '唯一id', icon: 'Key', dbType: 'PRIMARY KEY' },
  { type: 'integer', label: '整数', icon: 'EditPen', dbType: 'INTEGER' },
  { type: 'float', label: '浮点数', icon: 'EditPen', dbType: 'FLOAT' },
  { type: 'string', label: '字符串', icon: 'Document', dbType: 'VARCHAR' },
  { type: 'boolean', label: '布尔', icon: 'Switch', dbType: 'BOOLEAN' },
  { type: 'datetime', label: '时间', icon: 'Timer', dbType: 'DATETIME' },
  { type: 'single_enum', label: '单选枚举', icon: 'Select', dbType: 'VARCHAR' },
  { type: 'multi_enum', label: '多选枚举', icon: 'List', dbType: 'JSON' }
])

// 当前选中的字段
const selectedField = computed(() => {
  const fields = formSchema.value.fields || []
  if (selectedFieldIndex.value >= 0 && selectedFieldIndex.value < fields.length) {
    return fields[selectedFieldIndex.value]
  }
  return null
})

// 拖拽开始
const handleDragStart = (event, fieldType) => {
  event.dataTransfer.setData('fieldType', JSON.stringify(fieldType))
}

// 放置字段
const handleDrop = (event) => {
  event.preventDefault()
  isDragOver.value = false
  
  const fieldTypeData = event.dataTransfer.getData('fieldType')
  if (!fieldTypeData) return
  
  const fieldType = JSON.parse(fieldTypeData)
  const newField = createFieldFromType(fieldType)
  
  // 确保fields数组存在
  if (!formSchema.value.fields) {
    formSchema.value.fields = []
  }
  
  formSchema.value.fields.push(newField)
  selectedFieldIndex.value = formSchema.value.fields.length - 1
}

// 根据类型创建字段
const createFieldFromType = (fieldType) => {
  // 确保fields数组存在
  if (!formSchema.value.fields) {
    formSchema.value.fields = []
  }
  
  const baseField = {
    id: 'field_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9),
    type: fieldType.type,
    name: `field_${formSchema.value.fields.length + 1}`,
    label: fieldType.label,
    placeholder: `请输入${fieldType.label}`,
    required: false,
    default_value: '',
    dbType: fieldType.dbType
  }
  
  // 根据不同类型添加特定属性
  switch (fieldType.type) {
    case 'unique_id':
      return {
        ...baseField,
        id_type: 'auto_increment', // 默认为自增整数
        input_type: 'hidden' // 默认为隐藏字段
      }
    case 'integer':
      return {
        ...baseField,
        min_value: null,
        max_value: null,
        input_type: 'number'
      }
    case 'float':
      return {
        ...baseField,
        min_value: null,
        max_value: null,
        precision: 2,
        input_type: 'number'
      }
    case 'string':
      return {
        ...baseField,
        min_length: null,
        max_length: null,
        format: 'text',
        input_type: 'input',
        textarea_rows: 4
      }
    case 'boolean':
      return {
        ...baseField,
        input_type: 'switch'
      }
    case 'datetime':
      return {
        ...baseField,
        time_format: 'datetime', // 默认日期时间格式
        input_type: 'datetime'
      }
    case 'single_enum':
      return {
        ...baseField,
        enum_options: [
          { value: 'option1', label: '选项1' },
          { value: 'option2', label: '选项2' }
        ],
        input_type: 'select'
      }
    case 'multi_enum':
      return {
        ...baseField,
        enum_options: [
          { value: 'option1', label: '选项1' },
          { value: 'option2', label: '选项2' }
        ],
        input_type: 'checkbox'
      }
    default:
      return baseField
  }
}

// 选择字段
const selectField = (index) => {
  selectedFieldIndex.value = index
}

// 获取字段类型显示名称
const getFieldTypeLabel = (type) => {
  const typeMap = {
    unique_id: '唯一id',
    integer: '整数',
    float: '浮点数',
    string: '字符串',
    boolean: '布尔',
    datetime: '时间',
    single_enum: '单选枚举',
    multi_enum: '多选枚举'
  }
  return typeMap[type] || type
}

// 获取数据库类型
const getDbType = (type) => {
  const dbTypeMap = {
    unique_id: 'PRIMARY KEY',
    integer: 'INTEGER',
    float: 'FLOAT',
    string: 'VARCHAR',
    boolean: 'BOOLEAN',
    datetime: 'DATETIME',
    single_enum: 'VARCHAR',
    multi_enum: 'JSON'
  }
  return dbTypeMap[type] || 'VARCHAR'
}

// 添加枚举选项
const addEnumOption = () => {
  if (!selectedField.value.enum_options) {
    selectedField.value.enum_options = []
  }
  selectedField.value.enum_options.push({
    value: `option${selectedField.value.enum_options.length + 1}`,
    label: `选项${selectedField.value.enum_options.length + 1}`
  })
}

// 删除枚举选项
const removeEnumOption = (index) => {
  if (selectedField.value.enum_options) {
    selectedField.value.enum_options.splice(index, 1)
  }
}

// 移动字段
const moveFieldUp = (index) => {
  if (index > 0) {
    const fields = formSchema.value.fields
    const temp = fields[index]
    fields[index] = fields[index - 1]
    fields[index - 1] = temp
    selectedFieldIndex.value = index - 1
  }
}

const moveFieldDown = (index) => {
  const fields = formSchema.value.fields
  if (index < fields.length - 1) {
    const temp = fields[index]
    fields[index] = fields[index + 1]
    fields[index + 1] = temp
    selectedFieldIndex.value = index + 1
  }
}

// 复制字段
const duplicateField = (index) => {
  const originalField = formSchema.value.fields[index]
  const duplicatedField = {
    ...JSON.parse(JSON.stringify(originalField)),
    id: 'field_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9),
    name: originalField.name + '_copy'
  }
  formSchema.value.fields.splice(index + 1, 0, duplicatedField)
  selectedFieldIndex.value = index + 1
}

// 删除字段
const removeField = (index) => {
  if (!formSchema.value.fields) {
    formSchema.value.fields = []
    return
  }
  
  formSchema.value.fields.splice(index, 1)
  if (selectedFieldIndex.value >= formSchema.value.fields.length) {
    selectedFieldIndex.value = formSchema.value.fields.length - 1
  }
}

// 添加选项
const addOption = () => {
  if (selectedField.value && selectedField.value.options) {
    selectedField.value.options.push({
      label: `选项${selectedField.value.options.length + 1}`,
      value: `option${selectedField.value.options.length + 1}`
    })
  }
}

// 删除选项
const removeOption = (index) => {
  if (selectedField.value && selectedField.value.options) {
    selectedField.value.options.splice(index, 1)
  }
}

// 预览表单
const previewForm = () => {
  if (!formSchema.value.fields || formSchema.value.fields.length === 0) {
    ElMessage.warning('请先添加字段')
    return
  }
  previewVisible.value = true
}

// 清空表单
const clearForm = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有字段吗？', '确认清空', {
      type: 'warning'
    })
    
    formSchema.value.fields = []
    selectedFieldIndex.value = -1
    ElMessage.success('表单已清空')
  } catch (error) {
    // 用户取消
  }
}

// 保存表单
const saveForm = async () => {
  if (!formSchema.value.name.trim()) {
    ElMessage.warning('请输入表单名称')
    return
  }
  
  // 允许保存空字段的表单，用户可以后续添加字段
  // if (!formSchema.value.fields || formSchema.value.fields.length === 0) {
  //   ElMessage.warning('请至少添加一个字段')
  //   return
  // }
  
  // 验证字段名称唯一性
  const fieldNames = formSchema.value.fields.map(f => f.name)
  const uniqueNames = new Set(fieldNames)
  if (fieldNames.length !== uniqueNames.size) {
    ElMessage.error('字段名称不能重复')
    return
  }
  
  saving.value = true
  try {
    const data = {
      name: formSchema.value.name,
      description: formSchema.value.description,
      fields: formSchema.value.fields
    }
    
    if (formSchema.value.id) {
      // 更新现有表单
      await api.put(`/forms/${formSchema.value.id}`, data)
      ElMessage.success('表单更新成功')
    } else {
      // 创建新表单
      await api.post('/forms', data)
      ElMessage.success('表单创建成功')
    }
    
    // 返回表单列表
    router.push('/forms')
  } catch (error) {
    console.error('保存表单失败:', error)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 加载表单（编辑模式）
const loadForm = async (id) => {
  try {
    const response = await api.get(`/forms/${id}`)
    const formData = response.data.data
    
    // 解析schema字段，从 API 返回的数据中提取字段信息
    let fields = []
    if (formData.schema && formData.schema.fields) {
      fields = formData.schema.fields
    }
    
    // 确保数据结构完整
    formSchema.value = {
      id: formData.id || null,
      name: formData.name || '',
      description: formData.description || '',
      fields: fields
    }
  } catch (error) {
    console.error('加载表单失败:', error)
    ElMessage.error('加载表单失败')
    router.push('/forms')
  }
}

onMounted(() => {
  // 如果有ID参数，说明是编辑模式
  if (route.params.id) {
    loadForm(route.params.id)
  }
})
</script>

<style scoped>
.page-container {
  height: 100vh;
  overflow: hidden;

  .designer-container {
    display: flex;
    height: 100%;

    .fields-panel {
      width: 250px;
      background: #f8f9fa;
      border-right: 1px solid #e4e7ed;
      overflow-y: auto;
      padding: 16px;

      h3 {
        margin: 0 0 16px 0;
        font-size: 16px;
        color: #303133;
      }

      h4 {
        margin: 16px 0 8px 0;
        font-size: 14px;
        color: #606266;
        font-weight: normal;
      }

      .field-items {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 8px;
        margin-bottom: 16px;

        .field-item {
          padding: 12px 8px;
          border: 1px solid #dcdfe6;
          border-radius: 4px;
          background: white;
          cursor: grab;
          text-align: center;
          transition: all 0.3s ease;
          user-select: none;
          position: relative;

          &:hover {
            border-color: #409eff;
            background: #ecf5ff;
          }

          &:active {
            cursor: grabbing;
          }

          span {
            display: block;
            font-size: 13px;
            color: #303133;
            margin: 4px 0;
            font-weight: 500;
          }
        }
      }
    }

    .design-area {
      flex: 1;
      display: flex;
      flex-direction: column;
      background: #fff;

      .design-header {
        padding: 16px 24px;
        border-bottom: 1px solid #e4e7ed;

        .design-title {
          margin-bottom: 12px;
        }
      }

      .fields-canvas {
        flex: 1;
        padding: 24px;
        overflow-y: auto;

        .drop-zone {
          min-height: 400px;
          border: 2px dashed #dcdfe6;
          border-radius: 8px;
          padding: 24px;
          transition: all 0.3s ease;

          &.drag-over {
            border-color: #409eff;
            background: #ecf5ff;
          }

          .empty-hint {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            height: 300px;
            color: #c0c4cc;

            p {
              margin: 16px 0 0 0;
              font-size: 14px;
            }
          }
          
          .fields-list {
            .field-row {
              display: flex;
              align-items: center;
              justify-content: space-between;
              padding: 12px 16px;
              margin-bottom: 8px;
              border: 1px solid #e4e7ed;
              border-radius: 6px;
              background: #fafbfc;
              cursor: pointer;
              transition: all 0.3s ease;
              
              &:hover {
                border-color: #c0c4cc;
                background: #f5f7fa;
              }
              
              &.active {
                border-color: #409eff;
                background: #ecf5ff;
              }
              
              .field-info {
                flex: 1;
                
                .field-name {
                  font-size: 14px;
                  font-weight: 500;
                  color: #303133;
                  margin-bottom: 4px;
                }
                
                .field-type {
                  font-size: 12px;
                  color: #909399;
                }
              }
              
              .field-controls {
                opacity: 0;
                transition: opacity 0.3s ease;
              }
              
              &:hover .field-controls,
              &.active .field-controls {
                opacity: 1;
              }
            }
          }

          .field-wrapper {
            margin-bottom: 16px;
            border: 1px solid transparent;
            border-radius: 6px;
            padding: 12px;
            position: relative;
            transition: all 0.3s ease;

            &:hover {
              border-color: #c0c4cc;

              .field-controls {
                opacity: 1;
              }
            }

            &.active {
              border-color: #409eff;
              background: #ecf5ff;

              .field-controls {
                opacity: 1;
              }
            }

            .field-content {
              pointer-events: none;
            }

            .field-controls {
              position: absolute;
              top: -12px;
              right: 8px;
              opacity: 0;
              transition: opacity 0.3s ease;
            }
          }
        }
      }

      .design-actions {
        padding: 16px 24px;
        border-bottom: 1px solid #e4e7ed;
        background: #f8f9fa;
        display: flex;
        gap: 12px;
        justify-content: flex-end;
        margin-bottom: 16px;
      }
    }

    .properties-panel {
      width: 300px;
      background: #f8f9fa;
      border-left: 1px solid #e4e7ed;
      overflow-y: auto;
      padding: 16px;

      h3 {
        margin: 0 0 16px 0;
        font-size: 16px;
        color: #303133;
      }

      h4 {
        margin: 16px 0 8px 0;
        font-size: 14px;
        color: #606266;
        font-weight: normal;
        border-bottom: 1px solid #e4e7ed;
        padding-bottom: 4px;
      }

      .property-form {
        .el-form-item {
          margin-bottom: 12px;
        }

        .options-editor {
          .option-item {
            display: flex;
            gap: 8px;
            margin-bottom: 8px;
            align-items: center;

            .el-input {
              flex: 1;
            }
          }
        }
        
        .enum-options-editor {
          .enum-option-item {
            display: flex;
            gap: 8px;
            margin-bottom: 8px;
            align-items: center;
          }
        }
      }

      .no-selection {
        text-align: center;
        color: #c0c4cc;
        padding: 40px 0;

        p {
          margin: 0;
          font-size: 14px;
        }
      }
    }
  }

  .form-preview {
    h2 {
      margin: 0 0 8px 0;
      color: #303133;
    }

    p {
      margin: 0 0 24px 0;
      color: #606266;
    }
  }
}
</style>