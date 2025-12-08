# NexusHub-Personal 启动指南

## ⚠️ 数据库连接错误

当前遇到数据库连接失败：`Access denied for user 'root'@'localhost'`

### 解决步骤

#### 方案 1：确认 MySQL 密码（推荐）

1. 打开命令行，测试 MySQL 连接：
```bash
mysql -u root -p
# 输入密码：Gzh85218
```

2. 如果密码不正确，请修改 `backend/.env` 文件中的 `DB_PASSWORD`

#### 方案 2：创建新的 MySQL 用户

```sql
-- 登录 MySQL
mysql -u root -p

-- 创建数据库
CREATE DATABASE mywite CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户并授权
CREATE USER 'nexushub'@'localhost' IDENTIFIED BY 'nexushub123';
GRANT ALL PRIVILEGES ON mywite.* TO 'nexushub'@'localhost';
FLUSH PRIVILEGES;
```

然后修改 `.env`:
```
DB_USER=nexushub
DB_PASSWORD=nexushub123
```

#### 方案 3：使用空密码（仅本地开发）

如果您的 MySQL root 用户没有密码，修改 `.env`:
```
DB_PASSWORD=
```

### 完整启动流程

1. **确保 MySQL 服务运行**
```bash
# Windows
net start MySQL80

# 或通过服务管理器启动
```

2. **创建数据库**
```sql
mysql -u root -p
CREATE DATABASE IF NOT EXISTS mywite CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

3. **修改配置文件**
编辑 `backend/.env`，确认数据库配置正确

4. **启动后端服务**
```bash
cd NexusHub-Personal/backend
go run cmd/server/main.go
```

5. **验证服务**
```bash
# 新开一个终端
curl http://localhost:8080/health
# 应返回: {"status":"ok"}
```

## 项目功能

✅ 已完成的后端功能：
- 笔记管理 API
- 文件上传/下载/分类
- 代码片段管理
- 任务管理
- 书签管理
- 主题配置系统
- AI 聊天接口（待接入 AI API）

📋 待开发：
- Vue 前端界面
- 代码编辑器
- 文件预览
- 背景音乐播放器

## 测试 API

服务启动成功后，可以测试以下接口：

```bash
# 1. 健康检查
curl http://localhost:8080/health

# 2. 获取主题配置
curl http://localhost:8080/api/v1/theme

# 3. 创建笔记
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"测试笔记","content":"# Hello\n这是测试内容","tags":"测试"}'

# 4. 获取所有笔记
curl http://localhost:8080/api/v1/notes

# 5. 上传文件（替换为实际文件路径）
curl -X POST http://localhost:8080/api/v1/files/upload \
  -F "file=@test.jpg"

# 6. 创建代码片段
curl -X POST http://localhost:8080/api/v1/code \
  -H "Content-Type: application/json" \
  -d '{"title":"Hello Go","language":"go","code":"package main\n\nfunc main() {\n    println(\"Hello\")\n}"}'
```

完整 API 文档请查看 [README.md](../README.md)
