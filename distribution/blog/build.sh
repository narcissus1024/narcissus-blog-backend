#!/bin/bash

set -euo pipefail

# Help function
display_help() {
    echo "Usage: $0 [-t tag] [-r repository] [-p] [-h]"
    echo
    echo "Build and optionally push the Docker image for the blog backend."
    echo
    echo "Options:"
    echo "  -t TAG         Set the image tag. Defaults to 'latest'."
    echo "  -r REPO        Set the repository. Defaults to 'crpi-dvcqq1n8apk7iww5.cn-beijing.personal.cr.aliyuncs.com/narcissus1024/dev'."
    echo "  -p             Push the image to the repository after a successful build."
    echo "  -h             Display this help message and exit."
    exit 0
}

# 默认值
TAG="latest"
REPO="crpi-dvcqq1n8apk7iww5.cn-beijing.personal.cr.aliyuncs.com/narcissus1024"
SHOULD_PUSH=false
PROJECT_ROOT=$(dirname "$(dirname "$(dirname "$(realpath "$0")")")") # 获取项目根目录

# 解析命令行参数
while getopts ":t:r:ph" flag
do
    case "${flag}" in
        t) TAG=${OPTARG};;
        r) REPO=${OPTARG};;
        p) SHOULD_PUSH=true;;
        h) display_help;;
        \?) echo "Invalid option: -$OPTARG" >&2; display_help; exit 1;;
        :) echo "Option -$OPTARG requires an argument." >&2; display_help; exit 1;;
    esac
done

# 镜像名称
IMAGE_NAME="blog-backend"
FULL_IMAGE_NAME="$REPO/$IMAGE_NAME:$TAG"

# --- Build Confirmation ---
echo "--- Build Information ---"
echo "Image:      $FULL_IMAGE_NAME"
echo "Platform:   linux/amd64"
echo "Dockerfile: $PROJECT_ROOT/distribution/blog/Dockerfile"
echo "Context:    $PROJECT_ROOT"
echo "-------------------------"

# Confirmation prompt
read -p "Do you want to proceed with the build? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Build cancelled."
    exit 1
fi

# 构建 Docker 镜像
echo "Starting build..."
docker build --platform linux/amd64 -t "$FULL_IMAGE_NAME" -f "$PROJECT_ROOT/distribution/blog/Dockerfile" "$PROJECT_ROOT"

# 如果构建失败，则退出
if [ $? -ne 0 ]; then
    echo "Docker build failed."
    exit 1
fi

# 如果指定了 -p 参数，则推送镜像
if [ "$SHOULD_PUSH" = true ]; then
    # --- Push Confirmation ---
    echo "--- Push Information ---"
    echo "Image to push: $FULL_IMAGE_NAME"
    echo "------------------------"
    read -p "Do you want to proceed with the push? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]
    then
        echo "Push cancelled."
        exit 0
    fi

    docker push "$FULL_IMAGE_NAME"
fi

echo "Done."