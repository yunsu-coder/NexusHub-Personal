# 云存储集成配置指南

## 概述

NexusHub-Personal 现在支持集成第三方对象存储服务（当前支持 MinIO）。当配置了云存储后，文件上传将自动使用云存储服务，未配置时自动回退到本地存储。

## 支持的云存储提供商

- **MinIO**: 兼容 S3 的对象存储服务
- **AWS S3**: 亚马逊云存储（计划支持）
- **阿里云 OSS**: 阿里云对象存储（计划支持）
- **腾讯云 COS**: 腾讯云对象存储（计划支持）

## MinIO 配置

### 1. 环境变量配置

在 `backend/.env` 文件中添加以下配置：

```env
# MinIO 云存储配置
CLOUD_STORAGE_PROVIDER=minio
CLOUD_STORAGE_ACCESS_KEY=your_minio_access_key
CLOUD_STORAGE_SECRET_KEY=your_minio_secret_key
CLOUD_STORAGE_BUCKET=your_bucket_name
CLOUD_STORAGE_ENDPOINT=http://localhost:9000
CLOUD_STORAGE_REGION=us-east-1
```

### 2. 配置说明

- **CLOUD_STORAGE_PROVIDER**: 设置为 `minio`
- **CLOUD_STORAGE_ACCESS_KEY**: MinIO 的访问密钥
- **CLOUD_STORAGE_SECRET_KEY**: MinIO 的密钥
- **CLOUD_STORAGE_BUCKET**: 存储桶名称（将自动创建）
- **CLOUD_STORAGE_ENDPOINT**: MinIO 服务地址
- **CLOUD_STORAGE_REGION**: 区域（可选）

### 3. MinIO 服务部署

#### 使用 Docker 部署 MinIO：

```bash
# 拉取 MinIO 镜像
docker pull minio/minio

# 运行 MinIO 容器
docker run -d \
  --name minio \
  -p 9000:9000 \
  -p 9001:9001 \
  -e MINIO_ROOT_USER=admin \
  -e MINIO_ROOT_PASSWORD=password123 \
  -v /path/to/data:/data \
  minio/minio server /data --console-address ":9001"
```

#### 使用 Docker Compose 部署：

创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  minio:
    image: minio/minio:latest
    container_name: minio
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: admin
      MINIO_ROOT_PASSWORD: password123
    volumes:
      - ./minio-data:/data
    command: server /data --console-address ":9001"
    restart: unless-stopped
```

启动服务：
```bash
docker-compose up -d
```

## 功能特性

### 1. 自动存储选择
- **启用云存储**: 文件自动上传到配置的云存储服务
- **本地存储**: 未配置云存储时自动使用本地存储
- **失败回退**: 云存储初始化失败时自动回退到本地存储

### 2. 文件管理
- **上传**: 支持批量上传，显示上传进度
- **下载**: 自动选择云存储或本地存储下载
- **删除**: 根据文件存储位置自动选择删除方式
- **预览**: 支持图片、文档等文件的云存储预览

### 3. 安全特性
- **访问控制**: 文件按用户ID组织存储
- **路径隔离**: 用户文件通过路径隔离
- **清理机制**: 数据库操作失败时自动清理云存储文件

## 存储路径结构

### 本地存储
```
storage/
├── uploads/
│   ├── media/
│   ├── document/
│   ├── code/
│   ├── archive/
│   └── other/
```

### 云存储（MinIO）
```
bucket/
├── {userID}_{category}_{timestamp}/
│   └── {filename}
```

## 监控和日志

系统会记录以下日志：
- 云存储初始化状态
- 文件上传/下载/删除操作
- 存储失败和回退情况
- 错误和警告信息

## 故障排除

### 1. 云存储连接失败
- 检查网络连接
- 验证访问密钥和密钥
- 确认 MinIO 服务运行状态

### 2. 存储桶不存在
- 系统会自动创建存储桶
- 检查 MinIO 权限配置

### 3. 文件上传失败
- 检查文件大小限制
- 验证网络连接
- 查看日志获取详细错误信息

## 配置检查

可以通过以下方式检查云存储配置：

1. **应用启动日志**: 查看初始化信息
2. **环境变量**: 检查 `.env` 文件配置
3. **连接测试**: 尝试上传文件验证配置

## 性能优化建议

1. **网络配置**: 使用内网地址减少延迟
2. **并发控制**: 合理设置上传并发数
3. **缓存策略**: 对频繁访问的文件启用缓存
4. **CDN集成**: 对于公开文件可考虑使用 CDN

## 注意事项

- 云存储配置会覆盖本地存储设置
- 现有本地文件不受影响
- 建议在生产环境中使用 HTTPS
- 定期备份重要数据
- 监控存储使用量和成本