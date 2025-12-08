# NexusHub-Personal

> 一个功能强大的个人工作站系统 - 集成文件管理、笔记、代码编辑、数据分析、AI聊天等多种功能

![Version](https://img.shields.io/badge/version-1.0.0-blue) ![License](https://img.shields.io/badge/license-MIT-green)

## ✨ 核心功能

- 🔐 **用户认证**: JWT token认证,安全可靠
- 📁 **文件管理**: 多格式文件上传下载,分类管理
- 📝 **笔记系统**: Markdown支持,富文本编辑
- 💻 **代码编辑器**: 语法高亮,多语言支持
- 📊 **数据分析**: CSV/Excel/PDF导入,5种可视化图表
- 🧮 **科学计算器**: 完整的科学计算功能
- ✅ **TODO管理**: 任务追踪,优先级管理
- 🤖 **AI聊天**: 支持OpenAI/DeepSeek/Claude等多种模型
- 🎨 **主题定制**: 拟物风格UI,背景图片和音乐

## 🚀 快速开始

### 前置要求
- Node.js 16+
- Go 1.21+
- MySQL 8.0+

### 1. 数据库配置
```sql
CREATE DATABASE nexushub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 启动后端
```bash
cd backend
go mod download
go run cmd/server/main.go
```
后端运行在 `http://localhost:8080`

### 3. 启动前端
```bash
cd frontend
npm install
npm run dev
```
前端运行在 `http://localhost:5173`

### 4. 开始使用
访问 `http://localhost:5173` 并注册账号!

## 📖 技术栈

**前端**: Vue 3, Element Plus, ECharts, Math.js, XLSX.js, PDF.js
**后端**: Go, Gin, GORM, JWT
**数据库**: MySQL 8.0

## 📂 项目结构

```
NexusHub-Personal/
├── frontend/          # Vue 3前端
├── backend/           # Go后端
├── DEPLOYMENT.md      # 部署文档
└── README.md          # 本文件
```

## 🚢 生产部署

详细部署指南请查看 [DEPLOYMENT.md](./DEPLOYMENT.md)

## 📱 功能页面

- `/auth` - 登录注册
- `/dashboard` - 仪表盘
- `/files` - 文件管理
- `/notes` - 笔记
- `/code` - 代码编辑
- `/data` - 数据分析
- `/calculator` - 计算器
- `/todos` - TODO管理
- `/chat` - AI聊天
- `/settings` - 系统设置

## 🎯 主要特色

### 1. 拟物化设计
精美的拟物风格UI,平滑动画,现代化视觉体验

### 2. 多格式支持
文档(PDF, Excel), 图片(JPG, PNG), 音视频(MP3, MP4), 数据(CSV, XLSX)

### 3. AI集成
支持多种AI模型,本地配置,数据安全

### 4. 智能数据分析
自动识别数据类型,5种可视化图表,AI辅助分析

## 📝 更新日志

### v1.0.0 (2024-12-08)
- ✅ 完整的用户认证系统
- ✅ 文件管理和预览
- ✅ 笔记和代码编辑
- ✅ 数据分析和可视化
- ✅ AI聊天集成
- ✅ 拟物化UI设计
- ✅ 完整部署文档

## 📄 许可证

MIT License

## 🙏 致谢

感谢Vue.js, Element Plus, Go Gin, GORM, ECharts及所有开源项目!

---

**⭐ 如果这个项目对你有帮助,请给个Star!**
