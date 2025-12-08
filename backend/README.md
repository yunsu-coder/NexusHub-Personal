# NexusHub-Personal Backend

## 快速运行

### 1. 确保 MySQL 已启动并创建数据库

```bash
mysql -u root -p
CREATE DATABASE mywite CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 配置环境变量

复制 `.env.example` 到 `.env`:
```bash
cp .env.example .env
```

确认配置正确（默认已配置）:
```
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=Gzh85218
DB_NAME=mywite
```

### 3. 安装依赖并运行

```bash
# 安装依赖
go mod tidy

# 运行服务器
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动

### 4. 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 获取主题配置
curl http://localhost:8080/api/v1/theme

# 创建笔记
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"测试笔记","content":"这是内容","tags":"测试"}'

# 上传文件
curl -X POST http://localhost:8080/api/v1/files/upload \
  -F "file=@/path/to/file.jpg"
```

## 项目结构

```
backend/
├── cmd/server/          # 应用入口
├── internal/
│   ├── config/          # 配置管理
│   ├── database/        # 数据库初始化
│   ├── handler/         # HTTP 请求处理
│   ├── middleware/      # 中间件（CORS、用户上下文）
│   ├── model/           # 数据模型
│   ├── router/          # 路由注册
│   └── service/         # 业务逻辑
└── storage/             # 文件存储目录
```

## 数据库迁移

首次运行时会自动创建以下表:
- users (用户)
- notes (笔记)
- files (文件)
- tasks (任务)
- bookmarks (书签)
- code_snippets (代码片段)
- themes (主题配置)
- chat_messages (聊天记录)

并自动创建默认用户 (ID=1) 和默认主题。

## 开发建议

### 添加新功能

1. 在 `internal/model/` 添加数据模型
2. 在 `internal/service/` 添加业务逻辑
3. 在 `internal/handler/` 添加 HTTP 处理器
4. 在 `internal/router/router.go` 注册路由
5. 在 `internal/database/database.go` 的 `autoMigrate()` 中添加模型迁移

### 生产环境部署

1. 修改 `.env` 中 `GIN_MODE=release`
2. 使用 `go build` 编译二进制文件:
   ```bash
   go build -o nexushub cmd/server/main.go
   ```
3. 使用系统服务管理工具（如 systemd）运行
