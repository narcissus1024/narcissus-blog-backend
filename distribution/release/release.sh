#!/bin/bash
set -euo pipefail

# 脚本说明：将install目录和README.md打包为一个.tar.gz文件
# 用法：./distribution/release/release.sh [版本号]

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# 项目根目录
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
# 发布目录
RELEASE_DIR="${PROJECT_ROOT}/distribution/release"
# 临时目录
TMP_DIR="${RELEASE_DIR}/tmp"

# 默认版本号
VERSION=${1:-"v1.0.0"}
# 发布包名称
RELEASE_NAME="narcissus-blog-${VERSION}"
# 发布包路径
RELEASE_PACKAGE="${RELEASE_DIR}/${RELEASE_NAME}.tar.gz"

# 显示帮助信息
display_help() {
    echo "用法: $0 [版本号]"
    echo
    echo "将install目录和README.md打包为一个.tar.gz文件"
    echo "参数:"
    echo "  版本号    可选，指定发布版本号，默认为v1.0.0"
    echo
    echo "示例:"
    echo "  $0 v1.0.0"
    echo
    exit 1
}

# 如果第一个参数是-h或--help，显示帮助信息
if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    display_help
fi

# 清理函数，用于清理临时目录
cleanup() {
    echo "清理临时目录..."
    rm -rf "${TMP_DIR}"
}

# 注册退出时的清理函数
trap cleanup EXIT

# 创建临时目录
echo "创建临时目录..."
mkdir -p "${TMP_DIR}/${RELEASE_NAME}"

# 复制文件到临时目录
echo "复制文件到临时目录..."
cp -r "${PROJECT_ROOT}/distribution/install" "${TMP_DIR}/${RELEASE_NAME}/"
cp "${PROJECT_ROOT}/README.md" "${TMP_DIR}/${RELEASE_NAME}/"

# 确保install目录中的脚本有执行权限
echo "设置脚本执行权限..."
chmod +x "${TMP_DIR}/${RELEASE_NAME}/install/install.sh"

# 创建发布目录（如果不存在）
mkdir -p "${RELEASE_DIR}"

# 打包文件
echo "打包文件为 ${RELEASE_PACKAGE}..."
tar -czf "${RELEASE_PACKAGE}" -C "${TMP_DIR}" "${RELEASE_NAME}"

echo "打包完成: ${RELEASE_PACKAGE}"
echo "包含文件:"
tar -tvf "${RELEASE_PACKAGE}"

# 清理临时目录（由trap自动调用）