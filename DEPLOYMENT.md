# NexusHub-Personal 部署指南

## 📦 生产环境部署

### 系统要求
- Node.js 16+
- Go 1.21+
- MySQL 8.0+
- Nginx (推荐)

---

## 🔧 后端部署

### 1. 配置文件
创建生产环境配置文件 `backend/config/production.yaml`:

```yaml
server:
  port: 8080
  mode: release  # 生产模式

database:
  host: your-mysql-host
  port: 3306
  user: your-db-user
  password: your-db-password
  dbname: nexushub
  charset: utf8mb4

jwt:
  secret: your-secret-key-change-this  # 必须修改!
  expire_hours: 168  # 7天

upload:
  max_size: 524288000  # 500MB
  allowed_extensions:
    - .jpg
    - .png
    - .pdf
    - .xlsx
    - .mp4
    - .mp3
```

### 2. 编译后端
```bash
cd backend
go build -o nexushub-server cmd/server/main.go
```

### 3. 使用systemd管理服务
创建 `/etc/systemd/system/nexushub.service`:

```ini
[Unit]
Description=NexusHub Personal Server
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/nexushub/backend
ExecStart=/var/www/nexushub/backend/nexushub-server
Restart=always
RestartSec=5
StandardOutput=append:/var/log/nexushub/access.log
StandardError=append:/var/log/nexushub/error.log

[Install]
WantedBy=multi-user.target
```

启动服务:
```bash
sudo systemctl daemon-reload
sudo systemctl enable nexushub
sudo systemctl start nexushub
sudo systemctl status nexushub
```

---

## 🌐 前端部署

### 1. 环境变量配置
创建 `frontend/.env.production`:

```env
VITE_API_BASE_URL=https://yourdomain.com/api/v1
VITE_APP_TITLE=NexusHub Personal
```

### 2. 编译前端
```bash
cd frontend
npm install
npm run build
```

### 3. Nginx配置
创建 `/etc/nginx/sites-available/nexushub`:

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    # 前端静态文件
    location / {
        root /var/www/nexushub/frontend/dist;
        try_files $uri $uri/ /index.html;

        # 缓存策略
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # 后端API代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # 上传文件大小限制
        client_max_body_size 500M;
    }

    # 上传文件访问
    location /uploads/ {
        alias /var/www/nexushub/backend/uploads/;
        expires 1y;
        add_header Cache-Control "public";
    }
}
```

启用站点:
```bash
sudo ln -s /etc/nginx/sites-available/nexushub /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 🔒 SSL证书 (HTTPS)

使用Let's Encrypt免费证书:

```bash
sudo apt-get install certbot python3-certbot-nginx
sudo certbot --nginx -d yourdomain.com
sudo systemctl reload nginx
```

---

## 💾 数据库初始化

```sql
CREATE DATABASE nexushub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'nexushub_user'@'localhost' IDENTIFIED BY 'strong_password';
GRANT ALL PRIVILEGES ON nexushub.* TO 'nexushub_user'@'localhost';
FLUSH PRIVILEGES;
```

---

## 📝 Docker部署 (可选)

### Dockerfile (后端)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o nexushub-server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/nexushub-server .
COPY --from=builder /app/config ./config
EXPOSE 8080
CMD ["./nexushub-server"]
```

### docker-compose.yml
```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: nexushub
      MYSQL_USER: nexushub_user
      MYSQL_PASSWORD: nexushub_pass
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    depends_on:
      - mysql
    environment:
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: nexushub_user
      DB_PASSWORD: nexushub_pass
      DB_NAME: nexushub
    volumes:
      - ./backend/uploads:/root/uploads

  frontend:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./frontend/dist:/usr/share/nginx/html
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - backend

volumes:
  mysql_data:
```

启动:
```bash
docker-compose up -d
```

---

## 🔧 环境变量

### 后端环境变量
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=nexushub_user
export DB_PASSWORD=your_password
export DB_NAME=nexushub
export JWT_SECRET=your-secret-key
export GIN_MODE=release
export PORT=8080
```

---

## 📊 监控和日志

### 1. 日志管理
```bash
# 创建日志目录
sudo mkdir -p /var/log/nexushub
sudo chown www-data:www-data /var/log/nexushub

# 日志轮转配置
sudo nano /etc/logrotate.d/nexushub
```

```
/var/log/nexushub/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 www-data www-data
    sharedscripts
}
```

### 2. 监控
推荐使用:
- **Prometheus + Grafana**: 系统监控
- **ELK Stack**: 日志分析
- **Uptime Kuma**: 服务可用性监控

---

## 🔐 安全建议

1. **修改默认JWT密钥**
   - 在 `backend/internal/middleware/auth.go` 中修改 `jwtSecret`

2. **数据库安全**
   - 使用强密码
   - 限制远程访问
   - 定期备份

3. **文件上传安全**
   - 限制文件类型
   - 扫描病毒
   - 限制文件大小

4. **防火墙配置**
   ```bash
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw allow 22/tcp
   sudo ufw enable
   ```

5. **定期更新**
   ```bash
   sudo apt update && sudo apt upgrade
   ```

---

## 📦 备份策略

### 数据库备份
```bash
#!/bin/bash
# backup.sh
BACKUP_DIR="/var/backups/nexushub"
DATE=$(date +%Y%m%d_%H%M%S)

mysqldump -u nexushub_user -p nexushub > $BACKUP_DIR/db_$DATE.sql
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz /var/www/nexushub/backend/uploads

# 保留最近7天的备份
find $BACKUP_DIR -type f -mtime +7 -delete
```

添加到crontab:
```bash
0 2 * * * /var/www/nexushub/backup.sh
```

---

## 🚀 性能优化

1. **启用Gzip压缩** (Nginx)
```nginx
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;
```

2. **数据库索引优化**
```sql
CREATE INDEX idx_files_user_id ON files(user_id);
CREATE INDEX idx_notes_user_id ON notes(user_id);
CREATE INDEX idx_tasks_user_id_status ON tasks(user_id, status);
```

3. **CDN加速**
   - 使用CDN加速静态资源
   - 推荐: Cloudflare, AWS CloudFront

---

## ✅ 部署检查清单

- [ ] 数据库已创建并配置
- [ ] 后端配置文件已更新
- [ ] JWT密钥已修改
- [ ] 前端API地址已配置
- [ ] Nginx配置已完成
- [ ] SSL证书已安装
- [ ] 防火墙规则已设置
- [ ] 日志目录已创建
- [ ] 备份脚本已配置
- [ ] 服务自启动已启用
- [ ] 初始管理员账号已创建

---

## 📞 故障排查

### 后端无法启动
```bash
# 检查日志
sudo journalctl -u nexushub -f

# 检查端口占用
sudo netstat -tulpn | grep 8080

# 检查数据库连接
mysql -u nexushub_user -p -h localhost nexushub
```

### 前端页面404
```bash
# 检查Nginx配置
sudo nginx -t

# 检查文件权限
ls -la /var/www/nexushub/frontend/dist

# 查看Nginx日志
sudo tail -f /var/log/nginx/error.log
```

### 文件上传失败
- 检查Nginx `client_max_body_size`
- 检查uploads目录权限
- 检查磁盘空间

---

## 🎉 完成!

部署完成后,访问 `https://yourdomain.com` 开始使用 NexusHub Personal!

首次使用请注册管理员账号。

---

**技术支持**: GitHub Issues
**文档版本**: v1.0.0
**更新日期**: 2024-12-08
