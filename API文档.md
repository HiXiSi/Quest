# 资料管理平台 API 接口文档

## 1. 概述

本文档描述了资料管理平台的后端 API 接口，包括用户管理、文件管理、分类管理、标签管理等功能模块的接口说明。

**基础信息：**
- 基础URL: `http://localhost:8080/api`
- 认证方式: JWT Token (Bearer Token)
- 响应格式: JSON

## 2. 通用响应格式

所有API接口都采用统一的响应格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

**字段说明：**
- `code`: 响应状态码，0表示成功，其他值表示错误
- `message`: 响应消息
- `data`: 响应数据，可能为对象、数组或null

**分页响应格式：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "pages": 5
  }
}
```

## 3. 用户认证接口

### 3.1 用户注册

**接口地址：** `POST /auth/register`

**请求参数：**
```json
{
  "username": "admin",
  "email": "admin@example.com",
  "password": "password123"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "user",
      "avatar": "",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  }
}
```

### 3.2 用户登录

**接口地址：** `POST /auth/login`

**请求参数：**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**响应示例：** 同注册接口

## 4. 用户管理接口

### 4.1 获取用户个人信息

**接口地址：** `GET /users/profile`

**请求头：** `Authorization: Bearer {token}`

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "user",
    "avatar": "",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 4.2 更新用户个人信息

**接口地址：** `PUT /users/profile`

**请求头：** `Authorization: Bearer {token}`

**请求参数：**
```json
{
  "email": "newemail@example.com",
  "avatar": "http://example.com/avatar.jpg"
}
```

## 5. 文件管理接口

### 5.1 文件上传

**接口地址：** `POST /files/upload`

**请求头：** `Authorization: Bearer {token}`

**请求类型：** `multipart/form-data`

**请求参数：**
- `file`: 文件对象（必需）
- `category_id`: 分类ID（可选）
- `description`: 文件描述（可选）
- `tag_ids`: 标签ID列表，逗号分隔（可选）

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "file_id": 1,
    "file_name": "document_1234567890.pdf",
    "file_size": 1024000,
    "file_path": "../uploads/2023/01/01/document_1234567890.pdf"
  }
}
```

### 5.2 获取文件列表

**接口地址：** `GET /files`

**请求头：** `Authorization: Bearer {token}`

**查询参数：**
- `page`: 页码（默认1）
- `page_size`: 每页大小（默认20）
- `keyword`: 搜索关键词
- `category_id`: 分类ID
- `tag_id`: 标签ID
- `file_type`: 文件类型
- `sort_by`: 排序字段（默认created_at）
- `sort_order`: 排序方向（asc/desc，默认desc）

### 5.3 获取单个文件信息

**接口地址：** `GET /files/{id}`

**请求头：** `Authorization: Bearer {token}`

### 5.4 更新文件信息

**接口地址：** `PUT /files/{id}`

**请求头：** `Authorization: Bearer {token}`

**请求参数：**
```json
{
  "original_name": "新文件名.pdf",
  "description": "文件描述",
  "category_id": 1,
  "tag_ids": [1, 2, 3],
  "is_public": true
}
```

### 5.5 删除文件（移到回收站）

**接口地址：** `DELETE /files/{id}`

**请求头：** `Authorization: Bearer {token}`

### 5.6 恢复文件

**接口地址：** `POST /files/{id}/restore`

**请求头：** `Authorization: Bearer {token}`

### 5.7 下载文件

**接口地址：** `GET /files/{id}/download`

**请求头：** `Authorization: Bearer {token}`

### 5.8 预览文件

**接口地址：** `GET /files/{id}/preview`

**请求头：** `Authorization: Bearer {token}`

### 5.9 批量删除文件

**接口地址：** `POST /files/batch-delete`

**请求头：** `Authorization: Bearer {token}`

**请求参数：**
```json
{
  "file_ids": [1, 2, 3]
}
```

### 5.10 批量恢复文件

**接口地址：** `POST /files/batch-restore`

**请求头：** `Authorization: Bearer {token}`

**请求参数：**
```json
{
  "file_ids": [1, 2, 3]
}
```

## 6. 分类管理接口

### 6.1 获取分类列表

**接口地址：** `GET /categories`

**请求头：** `Authorization: Bearer {token}`

**查询参数：**
- `tree`: 是否返回树形结构（true/false）

### 6.2 创建分类

**接口地址：** `POST /categories`

**请求头：** `Authorization: Bearer {token}`

**请求参数：**
```json
{
  "name": "文档分类",
  "description": "分类描述",
  "parent_id": null,
  "icon": "document",
  "color": "#409eff",
  "sort": 0,
  "is_active": true
}
```

### 6.3 更新分类

**接口地址：** `PUT /categories/{id}`

**请求头：** `Authorization: Bearer {token}`

### 6.4 删除分类

**接口地址：** `DELETE /categories/{id}`

**请求头：** `Authorization: Bearer {token}`

## 7. 标签管理接口

### 7.1 获取标签列表

**接口地址：** `GET /tags`

**请求头：** `Authorization: Bearer {token}`

**查询参数：**
- `page`: 页码（默认1）
- `page_size`: 每页大小（默认20）
- `keyword`: 搜索关键词

### 7.2 创建标签

**接口地址：** `POST /tags`

**请求头：** `Authorization: Bearer {token}`

**请求参数：**
```json
{
  "name": "重要",
  "color": "#f56c6c",
  "description": "重要文件标签"
}
```

### 7.3 更新标签

**接口地址：** `PUT /tags/{id}`

**请求头：** `Authorization: Bearer {token}`

### 7.4 删除标签

**接口地址：** `DELETE /tags/{id}`

**请求头：** `Authorization: Bearer {token}`

## 8. 回收站接口

### 8.1 获取回收站文件列表

**接口地址：** `GET /recycle`

**请求头：** `Authorization: Bearer {token}`

**查询参数：**
- `page`: 页码（默认1）
- `page_size`: 每页大小（默认20）

### 8.2 清空回收站

**接口地址：** `DELETE /recycle/empty`

**请求头：** `Authorization: Bearer {token}`

## 9. 管理员接口

### 9.1 获取所有用户

**接口地址：** `GET /admin/users`

**请求头：** `Authorization: Bearer {token}`

**权限要求：** 管理员

**查询参数：**
- `page`: 页码（默认1）
- `page_size`: 每页大小（默认10）

### 9.2 更新用户信息

**接口地址：** `PUT /admin/users/{id}`

**请求头：** `Authorization: Bearer {token}`

**权限要求：** 管理员

**请求参数：**
```json
{
  "username": "newusername",
  "email": "newemail@example.com",
  "role": "admin",
  "avatar": "http://example.com/avatar.jpg"
}
```

### 9.3 删除用户

**接口地址：** `DELETE /admin/users/{id}`

**请求头：** `Authorization: Bearer {token}`

**权限要求：** 管理员

### 9.4 获取系统统计信息

**接口地址：** `GET /admin/stats`

**请求头：** `Authorization: Bearer {token}`

**权限要求：** 管理员

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_count": 100,
    "file_count": 1000,
    "category_count": 20,
    "tag_count": 50,
    "total_size": 1073741824,
    "total_size_formatted": "1.0 GB"
  }
}
```

## 10. 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权，需要登录 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 11. 使用示例

### 用户登录并上传文件示例

1. **用户登录：**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

2. **上传文件：**
```bash
curl -X POST http://localhost:8080/api/files/upload \
  -H "Authorization: Bearer {获取到的token}" \
  -F "file=@/path/to/your/file.pdf" \
  -F "description=测试文件"
```

3. **获取文件列表：**
```bash
curl -X GET "http://localhost:8080/api/files?page=1&page_size=10" \
  -H "Authorization: Bearer {token}"
```

## 12. 注意事项

1. 所有需要认证的接口都必须在请求头中包含有效的JWT token
2. 文件上传限制大小为100MB
3. 管理员权限接口需要用户角色为"admin"
4. 删除操作通常是软删除，文件会被移到回收站
5. 清空回收站会永久删除文件，请谨慎操作