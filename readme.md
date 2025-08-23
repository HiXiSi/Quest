# 资料管理平台

一个功能完整的前后端分离Web应用，用于管理结构化数据和非结构化文件，支持文件分类、标签管理等功能。

## 项目概述

### 核心功能
- **用户管理**：用户注册、登录、权限控制
- **文件管理**：多文件上传、拖拽上传、文件预览与下载
- **分类管理**：支持树状结构的文件分类
- **标签管理**：多标签系统，支持标签筛选
- **回收站功能**：删除的文件可恢复
- **文件预览**：支持图片、PDF、文本文件在线预览

### 技术栈

#### 后端
- **语言**：Go (Golang)
- **Web框架**：Gin
- **ORM**：GORM
- **数据库**：SQLite（可扩展为PostgreSQL）
- **认证**：JWT

#### 前端
- **框架**：Vue 3 + Composition API
- **构建工具**：Vite
- **UI组件库**：Element Plus
- **状态管理**：Pinia
- **HTTP客户端**：Axios

## 项目结构

```
Quest1/
├── backend/                 # 后端代码
│   ├── config/             # 配置文件
│   ├── controllers/        # 控制器
│   ├── middlewares/        # 中间件
│   ├── models/             # 数据模型
│   ├── routes/             # 路由配置
│   ├── utils/              # 工具函数
│   ├── main.go             # 主入口文件
│   └── go.mod              # Go模块文件
├── frontend/               # 前端代码
│   ├── src/
│   │   ├── components/     # 组件
│   │   ├── views/          # 页面
│   │   ├── stores/         # 状态管理
│   │   ├── utils/          # 工具函数
│   │   └── router/         # 路由配置
│   ├── package.json        # 前端依赖
│   └── vite.config.js      # Vite配置
├── uploads/                # 文件上传目录
├── static/                 # 静态文件目录
└── README.md               # 项目说明
```

## 快速开始

### 前置要求
- Go 1.21+
- Node.js 18+
- npm 或 yarn

### 后端启动

1. 进入后端目录
```bash
cd backend
```

2. 安装依赖
```bash
go mod download
```

3. 启动服务
```bash
go run main.go
```

后端服务将在 `http://localhost:8080` 启动

### 前端启动

1. 进入前端目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
```

3. 启动开发服务器
```bash
npm run dev
```

前端应用将在 `http://localhost:3000` 启动

## API文档

### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

### 用户管理
- `GET /api/users/profile` - 获取用户信息
- `PUT /api/users/profile` - 更新用户信息

### 分类管理
- `GET /api/categories` - 获取分类列表
- `POST /api/categories` - 创建分类
- `PUT /api/categories/:id` - 更新分类
- `DELETE /api/categories/:id` - 删除分类

### 标签管理
- `GET /api/tags` - 获取标签列表
- `POST /api/tags` - 创建标签
- `PUT /api/tags/:id` - 更新标签
- `DELETE /api/tags/:id` - 删除标签

### 文件管理
- `GET /api/files` - 获取文件列表
- `POST /api/files/upload` - 文件上传
- `GET /api/files/:id` - 获取文件信息
- `PUT /api/files/:id` - 更新文件信息
- `DELETE /api/files/:id` - 删除文件（移到回收站）
- `POST /api/files/:id/restore` - 恢复文件
- `GET /api/files/:id/download` - 下载文件
- `GET /api/files/:id/preview` - 预览文件

### 回收站
- `GET /api/recycle` - 获取回收站文件
- `DELETE /api/recycle/empty` - 清空回收站

## 功能特性

### 文件上传
- 支持多文件选择上传
- 支持拖拽上传
- 文件大小限制（100MB）
- 文件类型检测
- MD5和SHA256哈希验证

### 文件管理
- 文件列表展示
- 按分类、标签筛选
- 关键词搜索
- 文件重命名
- 文件移动（更改分类）
- 批量操作

### 文件预览
- 图片在线预览
- PDF文件预览
- 文本文件在线查看和编辑
- 缩略图生成

### 用户权限
- 普通用户：管理自己的文件
- 管理员：管理所有用户和文件
- JWT Token认证

## 部署说明

### 后端部署
1. 编译二进制文件
```bash
go build -o material-platform main.go
```

2. 运行
```bash
./material-platform
```

### 前端部署
1. 构建生产版本
```bash
npm run build
```

2. 将 `dist` 目录部署到Web服务器

## 数据库设计

### 主要数据表
- `users` - 用户表
- `categories` - 分类表（支持层级结构）
- `tags` - 标签表
- `files` - 文件表
- `file_tags` - 文件标签关联表

## 扩展计划

1. **存储扩展**：集成MinIO对象存储
2. **数据库升级**：迁移到PostgreSQL
3. **文件处理**：图片压缩、视频转码
4. **搜索功能**：集成Elasticsearch全文搜索
5. **API扩展**：RESTful API完善
6. **移动端**：开发移动端应用

## 许可证

MIT License

## 贡献指南

欢迎提交Issue和Pull Request来改进项目。

## 联系方式

如有问题，请通过Issue联系我们。