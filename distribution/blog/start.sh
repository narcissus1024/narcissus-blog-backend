#!/bin/bash
set -euo pipefail

# 启动blog-backend容器脚本
# eg: ./start.sh -i /app/data/nginx/html/img
# [提示] 若需要默认的配置文件，可以通过ocker run -ti crpi-dvcqq1n8apk7iww5.cn-beijing.personal.cr.aliyuncs.com/narcissus1024/dev/blog-backend:latest  bash
#        命令启动容器，然后通过docker cp

# 显示帮助信息
display_help() {
    echo "Usage: $0 [options...]"
    echo
    echo "   -t, --tag            Image tag, default is 'latest'."
    echo "   -p, --port           The port is map to the local port, default is 9002."
    echo "   -i, --img-dir        The img dir in container, default is '/app/data/nginx/html/img'."
    echo "   -h, --help           Display this help message."
    echo
    exit 1
}

ROOT_DIR=${HOME}
BLOG_DIR="${ROOT_DIR}/narcissus-blog"
BLOG_IMG_DIR="${BLOG_DIR}/img" # 将容器内图片目录挂载到主机目录
BLOG_CONF_DIR="${BLOG_DIR}/conf" # 将容器内配置目录挂载到主机目录

TAG="latest" # 镜像标签
PORT=9002 # 映射本地服务端口

# 容器配置
TARGET_PORT=9002 # 容器内服务端口
TARGET_HOME_DIR="/app"
TARGET_IMG_DIR="${TARGET_HOME_DIR}/data/nginx/html/img" # 必须指定为配置文件中imgDataDir的值

while getopts "t:p:i:h" flag
do
    case "${flag}" in
        t) TAG=${OPTARG};;
        p) PORT=${OPTARG};;
        i) TARGET_IMG_DIR=${OPTARG};;
        h) display_help;;
        \?) display_help;;
    esac
done

# 验证必需参数是否为空
if [ -z "$TAG" ] || [ -z "$PORT" ] || [ -z "$TARGET_IMG_DIR" ]; then
    echo "Error: Missing required arguments."
    display_help
fi

# 如果BLOG_DIR不存在，则创建
if [ ! -d "$BLOG_DIR" ]; then
    mkdir -p "$BLOG_DIR"
fi

# 如果BLOG_IMG_DIR不存在，则创建
if [ ! -d "$BLOG_IMG_DIR" ]; then
    mkdir -p "$BLOG_IMG_DIR"
fi

# 如果BLOG_CONF_DIR不存在，则创建
if [ ! -d "$BLOG_CONF_DIR" ]; then
    mkdir -p "$BLOG_CONF_DIR"
    # 默认配置文件
    cat <<EOF > "$BLOG_CONF_DIR/conf.yaml"
app:
  name: blog_narcissus
  port: 9002
  domain: 'localhost'
  imgDataDir: /app/data/nginx/html/img
  imgProxyURL: http://127.0.0.1:9001
  privateKeyDir: /app/conf
  publicKeyDir: /app/conf
mysql:
  user: root
  password: admin
  host: 127.0.0.1
  port: 3306
  dbname: blog_narcissus
redis:
  host: 127.0.0.1
  port: 6379
  password: admin
  db: 0
logger:
  logLevel: info
  logFormat: logfmt
EOF

fi

IMAGE="crpi-dvcqq1n8apk7iww5.cn-beijing.personal.cr.aliyuncs.com/narcissus1024/dev/blog-backend:$TAG"

# 打印配置参数
echo "--- Docker run configuration ---"
echo "  TAG: $TAG"
echo "  PORT: $PORT"
echo "  TARGET_IMG_DIR: $TARGET_IMG_DIR"
echo "  IMAGE: $IMAGE"
echo "-------------------------"

# 启动容器
read -p "Do you want to start the container? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Container start cancelled."
    exit 1
fi
# 启动容器
docker run --name blog-backend -p $TARGET_PORT:$PORT \
-v "$BLOG_CONF_DIR":"${TARGET_HOME_DIR}/conf" \
-v "$BLOG_IMG_DIR":"$TARGET_IMG_DIR" \
-d $IMAGE
