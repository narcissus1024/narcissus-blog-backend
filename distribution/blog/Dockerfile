# Dockerfile for narcissus-blog

# ---- Build Stage ----
FROM golang:1.23-bullseye AS builder

# 更换 apt 镜像源为国内源（使用 HTTP 而非 HTTPS）
RUN cat <<EOF > /etc/apt/sources.list
deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye main contrib non-free
deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-updates main contrib non-free
deb http://mirrors.tuna.tsinghua.edu.cn/debian-security/ bullseye-security main
EOF

# 安装 CGO 构建工具
RUN apt-get update && apt-get install -y build-essential

# 设置工作目录
WORKDIR /app

# 设置 GOPROXY
ENV GOPROXY=https://goproxy.cn,direct

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码
COPY . .

# 编译应用
# CGO_ENABLED=1 启用 CGO
# -ldflags="-s -w" 去除调试信息，减小二进制文件大小
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags="-s -w" -o narcissus-blog ./cmd/blog/main.go

# ---- Release Stage ----
FROM debian:bullseye-slim

# 更换 apt 镜像源为国内源（使用 HTTP 而非 HTTPS）
RUN cat <<EOF > /etc/apt/sources.list
deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye main contrib non-free
deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-updates main contrib non-free
deb http://mirrors.tuna.tsinghua.edu.cn/debian-security/ bullseye-security main
EOF

# 先安装 ca-certificates 然后安装 curl 用于健康检查
RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译后的二进制文件
COPY --from=builder /app/narcissus-blog .

# 复制配置文件
COPY conf/ ./conf/

# 暴露端口
EXPOSE 9002

# 设置入口点
CMD ["./narcissus-blog"]