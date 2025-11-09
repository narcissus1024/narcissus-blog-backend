#!/bin/bash
set -euo pipefail

# --- 默认配置 ---
PROJECT_NAME="nblog"
# 镜像配置
REPO="crpi-dvcqq1n8apk7iww5.cn-beijing.personal.cr.aliyuncs.com/narcissus1024"

# 后端镜像
BACKEND_IMAGE_NAME="blog-backend"
BACKEND_IMAGE_TAG="latest"
BACKEND_PORT=9002

# 前端镜像
FRONTEND_IMAGE_NAME="blog-frontend"
FRONTEND_IMAGE_TAG="latest"
FRONTEND_PORT=80
FRONTEND_MAP_PORT=80

# mysql
MYSQL_PASSWORD="admin"
# redis
REDIS_PASSWORD="admin"

# 本地挂载目录配置
ROOT_DIR="${HOME}/narcissus-blog"
CONF_DIR="${ROOT_DIR}/conf" # 配置目录
SECRETS_DIR="${ROOT_DIR}/secrets" # 密钥目录
# 配置目录
BACKEND_CONF_DIR="${CONF_DIR}/backend"
FRONTEND_CONF_DIR="${CONF_DIR}/frontend"
NGINX_CONF_DIR="${CONF_DIR}/nginx"
MYSQL_CONF_DIR="${CONF_DIR}/mysql"
REDIS_CONF_DIR="${CONF_DIR}/redis"

# Nginx 配置
NGINX_CONF_FILE="${NGINX_CONF_DIR}/nginx.conf"
NGINX_CONF_D="${NGINX_CONF_DIR}/conf.d"
NGINX_LOG_DIR="${ROOT_DIR}/logs/nginx"
# ssl文件名：cert.pem和key.pem
SSL_CERT_DIR="${NGINX_CONF_DIR}/ssl" # 如需修改，请改为自定义目录

# 域名
DOMAIN="blog.cn"  # 后端使用cookie，主域名
IMG_DOMAIN="img.${DOMAIN}"  # 后端使用，图片域名
NGINX_SERVER_NAME="${DOMAIN} www.${DOMAIN}"
NGINX_IMG_SERVER_NAME="img.${DOMAIN}"

# Nginx 模板变量
LISTEN="80"  # 监听端口 443 ssl
SSL_OPEN="off"  # 是否启用SSL，可选值：on, off
SSL_CERTIFICATE="/etc/nginx/ssl/cert.pem"  # SSL证书路径（容器内路径）
SSL_CERTIFICATE_KEY="/etc/nginx/ssl/key.pem"  # SSL证书密钥路径（容器内路径）
BACKEND_DOMAIN="${PROJECT_NAME}-backend-1"  # 后端服务名
BACKEND_SERVER="${BACKEND_DOMAIN}:${BACKEND_PORT}"  # 后端服务地址

# 后端配置模板变量
IMG_PROXY_URL="${IMG_DOMAIN}"  # 图片代理URL
if [ "${SSL_OPEN}" = "on" ]; then
    IMG_PROXY_URL="https://${IMG_DOMAIN}"
else
    IMG_PROXY_URL="http://${IMG_DOMAIN}"
fi
MYSQL_HOST="${PROJECT_NAME}-mysql-1"  # MySQL主机
REDIS_HOST="${PROJECT_NAME}-redis-1"  # Redis主机

# === 脚本参数 START ===
# 服务启动配置
SERVER="all"  # 指定要启动的服务，默认为all

# 步骤控制标志
DO_CONFIG=false      # 步骤1：生成配置
DO_START=false       # 步骤2：生成 .env 并启动 docker compose
DO_CLEANUP=false     # 步骤3：清理 docker compose
# === 脚本参数 END ===

# 确认提示函数
confirm_prompt() {
    local prompt_message="$1"
    read -p "${prompt_message} (y/N) " REPLY
    if [[ ! $REPLY =~ ^[Yy]([Ee][Ss])?$ ]]; then
        return 1
    fi
    return 0
}

# 显示帮助信息
display_help() {
    echo "Usage: $0 [options...]"
    echo
    echo "Steps (if no step option is given, default: -c -s):"
    echo "   -c                    Step 1: prepare/generate configuration files and directories (idempotent)."
    echo "   -s                    Step 2: generate .env and start Docker Compose services."
    echo "   -d, --down            Step 3: clean Docker Compose services and configurations."
    echo
    echo "Other options:"
    echo "   --service             Service to start (e.g., 'nblog-backend', 'all'). Default is 'all'."
    echo "   -h, --help             Display this help message."
    echo
    exit 1
}

# 解析命令行参数
parse_args() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --service)
                SERVER="$2"
                shift 2
                ;;
            -c)
                DO_CONFIG=true
                shift
                ;;
            -s)
                DO_START=true
                shift
                ;;
            -d|--down)
                DO_CLEANUP=true
                shift
                ;;
            -h|--help)
                display_help
                ;;
            *)
                echo "Unknown option: $1"
                display_help
                ;;
        esac
    done
}

# 检查并安装OpenSSL
check_install_openssl() {
    echo "--- Checking OpenSSL installation... ---"
    if command -v openssl &> /dev/null; then
        echo "OpenSSL is already installed."
        return 0
    fi
    
    echo "OpenSSL not found. Attempting to install..."
    
    # 检测操作系统类型
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$NAME
    elif type lsb_release >/dev/null 2>&1; then
        OS=$(lsb_release -si)
    elif [ -f /etc/lsb-release ]; then
        . /etc/lsb-release
        OS=$DISTRIB_ID
    elif [ -f /etc/debian_version ]; then
        OS="Debian"
    else
        OS=$(uname -s)
    fi
    
    # 根据不同的操作系统安装OpenSSL
    case "$OS" in
        *Ubuntu*|*Debian*)
            echo "Detected Debian/Ubuntu system"
            sudo apt-get update && sudo apt-get install -y openssl
            ;;
        *CentOS*|*Red*|*Fedora*|*Amazon*)
            echo "Detected RHEL/CentOS/Fedora system"
            sudo yum install -y openssl
            ;;
        *SUSE*)
            echo "Detected SUSE system"
            sudo zypper install -y openssl
            ;;
        *Alpine*)
            echo "Detected Alpine system"
            apk add --no-cache openssl
            ;;
        *Darwin*)
            echo "Detected macOS system"
            if command -v brew &> /dev/null; then
                brew install openssl
            else
                echo "Homebrew not found. Please install Homebrew first or OpenSSL manually."
                return 1
            fi
            ;;
        *MINGW*|*MSYS*|*CYGWIN*)
            echo "Detected Windows system"
            echo "Please install OpenSSL manually on Windows."
            return 1
            ;;
        *)
            echo "Unsupported operating system: $OS"
            echo "Please install OpenSSL manually."
            return 1
            ;;
    esac
    
    # 验证安装
    if command -v openssl &> /dev/null; then
        echo "OpenSSL installed successfully."
        return 0
    else
        echo "Failed to install OpenSSL. Please install it manually."
        return 1
    fi
}

# 生成RSA密钥对
generate_rsa_keys() {
    echo "--- Generating RSA key pair... ---"
    local keys_dir="$BACKEND_CONF_DIR"
    local private_key="$keys_dir/private.pem"
    local public_key="$keys_dir/public.pem"
    local temp_public_key="$keys_dir/public_key_pkcs8.pem"
    
    # 如果密钥已存在，跳过生成
    if [ -f "$private_key" ] && [ -f "$public_key" ]; then
        echo "RSA keys already exist. Skipping generation."
        return 0
    fi
    
    # 确保目录存在
    mkdir -p "$keys_dir"
    
    # 检查OpenSSL是否安装
    if ! command -v openssl &> /dev/null; then
        echo "OpenSSL not found. Please install it first."
        return 1
    fi
    
    # 生成私钥 (PKCS1格式)
    echo "Generating private key (PKCS1 format)..."
    # 检查OpenSSL版本，如果是3.0+版本需要使用-traditional参数
    openssl_version=$(openssl version | awk '{print $2}')
    if [[ "$openssl_version" =~ ^3\. ]]; then
        openssl genrsa -traditional -out "$private_key" 2048
    else
        openssl genrsa -out "$private_key" 2048
    fi
    
    # 提取公钥 (PKCS8格式)
    echo "Extracting public key (PKCS8 format)..."
    openssl rsa -in "$private_key" -pubout -out "$temp_public_key"
    
    # 将公钥转换为PKCS1格式
    echo "Converting public key to PKCS1 format..."
    openssl rsa -pubin -in "$temp_public_key" -RSAPublicKey_out -out "$public_key"
    
    # 删除临时文件
    rm -f "$temp_public_key"
    
    # 设置权限
    chmod 600 "$private_key"
    chmod 644 "$public_key"
    
    echo "RSA key pair generated successfully:"
    echo "  Private key: $private_key"
    echo "  Public key: $public_key"
    return 0
}

prepare_passwords() {
    echo "--- Preparing password... ---"
    # 创建secrets目录
    if [ ! -d "$SECRETS_DIR" ]; then
        mkdir -p "$SECRETS_DIR"
        chmod 700 "$SECRETS_DIR"
    fi
    # 生成密码文件（每行一个密码）
    echo "$MYSQL_PASSWORD" > "$SECRETS_DIR/mysql_root_password"
    echo "$REDIS_PASSWORD" > "$SECRETS_DIR/redis_password"

    # 设置文件权限
    chmod 600 "$SECRETS_DIR"/*
}

# 准备前端配置
prepare_frontend() {
    echo "--- Preparing frontend configuration... ---"
    if [ ! -d "$NGINX_LOG_DIR" ]; then
        mkdir -p "$NGINX_LOG_DIR"
    fi
    if [ ! -d "$NGINX_CONF_DIR" ]; then
        mkdir -p "$NGINX_CONF_DIR"
        mkdir -p "$NGINX_CONF_D"
        
        # 复制原始模板文件
        cp "./nginx/server.conf" "$NGINX_CONF_D/server.conf.template"
        cp "./nginx/nginx.conf" "$NGINX_CONF_DIR/nginx.conf.template"
        
        # 生成 SSL 配置
        if [ "$SSL_OPEN" = "on" ]; then
            SSL_CONFIG="ssl on;\n    ssl_certificate      $SSL_CERTIFICATE;\n    ssl_certificate_key  $SSL_CERTIFICATE_KEY;"
        else
            SSL_CONFIG="# SSL is disabled"
        fi
        
        # 生成 server_name 配置
        if [ -n "$NGINX_SERVER_NAME" ]; then
            SERVER_NAME="server_name  $NGINX_SERVER_NAME;"
        else
            SERVER_NAME="# No server_name specified"
        fi

        # 生成 img server_name 配置
        if [ -n "$NGINX_IMG_SERVER_NAME" ]; then
            IMG_SERVER_NAME="server_name  $NGINX_IMG_SERVER_NAME;"
        else
            IMG_SERVER_NAME="# No server_name specified"
        fi
        
        # 替换 server.conf 模板变量
        sed -e "s|{{LISTEN}}|$LISTEN|g" \
            -e "s|{{SERVER_NAME}}|$SERVER_NAME|g" \
            -e "s|{{SSL_CONFIG}}|$SSL_CONFIG|g" \
            -e "s|{{BACKEND_DOMAIN}}|$BACKEND_DOMAIN|g" \
            -e "s|{{IMG_SERVER_NAME}}|$IMG_SERVER_NAME|g" \
            "$NGINX_CONF_D/server.conf.template" > "$NGINX_CONF_D/server.conf"
        
        # 替换 nginx.conf 模板变量
        sed -e "s|{{BACKEND_DOMAIN}}|$BACKEND_DOMAIN|g" \
            -e "s|{{BACKEND_SERVER}}|$BACKEND_SERVER|g" \
            "$NGINX_CONF_DIR/nginx.conf.template" > "$NGINX_CONF_DIR/nginx.conf"
    fi
}

# 准备后端配置
prepare_backend() {
    echo "--- Preparing backend configuration... ---"
    # BACKEND_CONF_DIR
    if [ ! -d "$BACKEND_CONF_DIR" ]; then
        mkdir -p "$BACKEND_CONF_DIR"
        
        # 复制原始模板文件
        cp "./backend/conf.yaml" "$BACKEND_CONF_DIR/conf.yaml.template"
        
        # 替换配置文件中的变量
        sed -e "s|{{BACKEND_PORT}}|$BACKEND_PORT|g" \
            -e "s|{{DOMAIN}}|$DOMAIN|g" \
            -e "s|{{IMG_PROXY_URL}}|$IMG_PROXY_URL|g" \
            -e "s|{{MYSQL_PASSWORD}}|$MYSQL_PASSWORD|g" \
            -e "s|{{MYSQL_HOST}}|$MYSQL_HOST|g" \
            -e "s|{{REDIS_HOST}}|$REDIS_HOST|g" \
            -e "s|{{REDIS_PASSWORD}}|$REDIS_PASSWORD|g" \
            "$BACKEND_CONF_DIR/conf.yaml.template" > "$BACKEND_CONF_DIR/conf.yaml"
        
        echo "Backend configuration file generated with the following settings:"
        echo "  BACKEND_PORT: $BACKEND_PORT"
        echo "  DOMAIN: $DOMAIN"
        echo "  MySQL Host: $MYSQL_HOST"
        echo "  Redis Host: $REDIS_HOST"
        
        # 检查OpenSSL是否安装
        if ! check_install_openssl; then
            echo "OpenSSL not found. Please install it first."
            exit 1
        fi
        # 生成RSA密钥对
        if ! generate_rsa_keys; then
            echo "Error: Failed to generate RSA keys. Please check the logs above."
            exit 1
        fi
    fi
}

# 准备 MySQL 配置
prepare_mysql() {
    echo "--- Preparing MySQL configuration... ---"
    if [ ! -d "$MYSQL_CONF_DIR" ]; then
        mkdir -p "$MYSQL_CONF_DIR"
        cp "./mysql/my.cnf" "$MYSQL_CONF_DIR/"
        cp "./mysql/sql.sql" "$MYSQL_CONF_DIR/"
    fi
}

# 准备 Redis 配置
prepare_redis() {
    echo "--- Preparing Redis configuration... ---"
    if [ ! -d "$REDIS_CONF_DIR" ]; then
        mkdir -p "$REDIS_CONF_DIR"
        cp "./redis/redis.conf" "$REDIS_CONF_DIR/"
    fi
}

# 生成 .env 文件
generate_env_file() {
    echo "--- Generating .env file... ---"
    cat <<EOF > .env
# Docker Compose Environment Variables

# Backend
BACKEND_PORT=${BACKEND_PORT}

# Images
REPO=${REPO}
BACKEND_IMAGE_NAME=${BACKEND_IMAGE_NAME}
BACKEND_IMAGE_TAG=${BACKEND_IMAGE_TAG}
FRONTEND_IMAGE_NAME=${FRONTEND_IMAGE_NAME}
FRONTEND_IMAGE_TAG=${FRONTEND_IMAGE_TAG}
FRONTEND_PORT=${FRONTEND_PORT}
FRONTEND_MAP_PORT=${FRONTEND_MAP_PORT}

# Volumes
BACKEND_CONF_DIR=${BACKEND_CONF_DIR}
FRONTEND_CONF_DIR=${FRONTEND_CONF_DIR}
MYSQL_CONF_DIR=${MYSQL_CONF_DIR}
REDIS_CONF_DIR=${REDIS_CONF_DIR}

# Nginx
NGINX_CONF_FILE=${NGINX_CONF_FILE}
NGINX_CONF_D=${NGINX_CONF_D}
NGINX_LOG_DIR=${NGINX_LOG_DIR}
SSL_CERT_DIR=${SSL_CERT_DIR}

# Secrets
SECRETS_DIR=${SECRETS_DIR}
EOF
}

# 清理服务和配置
cleanup_services() {
    echo "--- Cleaning up Docker Compose services and configurations ---"

    # 确认是否继续
    if ! confirm_prompt "This will stop all services with data and delete configurations. Continue?"; then
        echo "Cleanup cancelled."
        exit 1
    fi

    # 停止并删除容器
    echo "Stopping and removing Docker Compose services..."
    docker compose -p $PROJECT_NAME down -v
    
    # 删除配置目录
    echo "Removing configuration directories..."
    if [ -d "$ROOT_DIR" ]; then
        rm -rf "$ROOT_DIR"
        echo "Removed $ROOT_DIR"
    fi
    
    # 删除.env文件
    if [ -f ".env" ]; then
        rm -f ".env"
        echo "Removed .env file"
    fi

    echo "Cleanup completed successfully!"
}

# 主函数
main() {
    parse_args "$@"

    # 如果没有指定任何步骤参数，默认执行步骤1和步骤2
    if [ "$DO_CONFIG" = false ] && [ "$DO_START" = false ] && [ "$DO_CLEANUP" = false ]; then
        DO_CONFIG=true
        DO_START=true
    fi

    # 步骤1：生成配置
    if [ "$DO_CONFIG" = true ]; then
        prepare_passwords
        prepare_frontend
        if ! prepare_backend; then
            echo "Error: Failed to prepare backend configuration. Please check the logs above."
            exit 1
        fi
        prepare_mysql
        prepare_redis
    fi


    # 步骤2：生成 .env 并启动 docker compose
    if [ "$DO_START" = true ]; then
        if ! confirm_prompt "Please check the configuration of each service in ${CONF_DIR} is correct."; then
            echo "Deployment cancelled."
            exit 1
        fi
        generate_env_file

        echo "--- Docker Compose Configuration ---"
        echo "Project Name: $PROJECT_NAME"
        echo "Server: $SERVER"
        echo ".ENV: "
        cat .env
        echo "----------------------------------"

        if ! confirm_prompt "Do you want to start the services?"; then
            echo "Deployment cancelled."
            exit 1
        fi

        local services_to_start=""
        if [[ "$SERVER" != "all" ]]; then
            services_to_start=$SERVER
        fi

        echo "Starting services with docker-compose... ($SERVER)"
        docker compose -p $PROJECT_NAME up -d $services_to_start

        echo "Deployment completed successfully!"
    fi

    # 步骤3：清理 docker compose
    if [ "$DO_CLEANUP" = true ]; then
        cleanup_services
    fi
}

main "$@"
