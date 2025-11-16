# Narcissus Blog

## 项目简介

Narcissus Blog 是一个基于 GO和VUE 开发的博客系统，包含前端、后端、数据库和缓存等组件。

## 安装说明

### 前置条件

- Docker 和 Docker Compose
- Bash 环境
- OpenSSL（用于生成密钥，脚本会尝试自动安装）
- wget（用于下载安装包）

### 从 GitHub 下载并安装

您可以直接从 GitHub Releases 页面下载安装包，然后解压并使用安装脚本：

```bash
# 下载最新版本安装包（替换 v0.0.0 为实际版本号，如 v1.0.0）
wget https://github.com/narcissus1024/narcissus-blog-backend/releases/download/v0.0.0/narcissus-blog-v0.0.0.tar.gz


# 解压安装包
tar -xzf narcissus-blog-v0.0.0.tar.gz

# 进入解压后的目录
cd narcissus-blog-v0.0.0

# 运行安装脚本
./install/install.sh
```

### 安装脚本使用方法

安装脚本 `install.sh` 用于简化部署过程。该脚本分为三个步骤：

1. **生成配置文件**：为各服务创建必要的配置文件和目录
2. **生成 .env 并启动服务**：创建 Docker Compose 环境变量文件并启动服务
3. **清理服务**：停止并删除所有服务和数据

#### 基本用法

```bash
# 查看帮助信息
./install/install.sh -h

# 默认执行步骤1和2（生成配置并启动服务）
./install/install.sh

# 只生成配置文件，不启动服务
./install/install.sh -c

# 只启动服务（假设配置已存在）
./install/install.sh -s

# 清理所有服务和数据
./install/install.sh -d
```

### 配置说明

安装脚本会在用户的 HOME 目录下创建配置目录：`~/narcissus-blog/`，包含以下子目录：

```
~/narcissus-blog/
├── conf/            # 配置目录
│   ├── backend/     # 后端配置
│   ├── nginx/       # Nginx配置
│   ├── mysql/       # MySQL配置
│   └── redis/       # Redis配置
├── secrets/         # 密码和密钥
└── logs/            # 日志文件
    └── nginx/       # Nginx日志
```

#### 重要配置说明

##### 后端密钥配置

后端服务需要 RSA 密钥对进行传输加密。安装脚本会自动使用 OpenSSL 生成这些密钥：

- 私钥：`~/narcissus-blog/conf/backend/private.pem`（PKCS1 格式）
- 公钥：`~/narcissus-blog/conf/backend/public.pem`（PKCS1 格式）

如果需要使用自己的密钥，可以在运行脚本前手动将密钥文件放置在相应目录。

##### SSL 配置

如果需要启用 HTTPS（端口 443），请注意：

1. 修改脚本中的 `SSL_OPEN` 变量为 `on`
2. 将 SSL 证书文件放置在正确位置：
   - 证书文件：`~/narcissus-blog/conf/nginx/ssl/cert.pem`
   - 密钥文件：`~/narcissus-blog/conf/nginx/ssl/key.pem`

##### 自定义配置

如需修改默认配置，请编辑 `install/install.sh` 脚本中相关变量。

