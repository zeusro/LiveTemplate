#!/bin/bash

# Go 升级脚本
# 用于将 Go 升级到最新版本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查当前版本
echo -e "${GREEN}检查当前 Go 版本...${NC}"
CURRENT_VERSION=$(go version | awk '{print $3}')
echo -e "当前版本: ${CURRENT_VERSION}"

# 获取最新版本号
echo -e "${GREEN}获取最新 Go 版本...${NC}"
LATEST_VERSION=$(curl -s https://go.dev/VERSION?m=text | head -n 1)
echo -e "最新版本: ${LATEST_VERSION}"

# 检查是否需要升级
if [ "$CURRENT_VERSION" == "$LATEST_VERSION" ]; then
    echo -e "${GREEN}您的 Go 已经是最新版本！${NC}"
    exit 0
fi

# 检测系统架构
ARCH=$(uname -m)
if [ "$ARCH" == "arm64" ]; then
    GO_ARCH="arm64"
else
    GO_ARCH="amd64"
fi

# 设置下载 URL
GO_VERSION=${LATEST_VERSION#go}  # 移除 "go" 前缀
DOWNLOAD_URL="https://go.dev/dl/${LATEST_VERSION}.darwin-${GO_ARCH}.tar.gz"
TEMP_DIR=$(mktemp -d)
DOWNLOAD_FILE="${TEMP_DIR}/go.tar.gz"

echo -e "${YELLOW}准备升级到 ${LATEST_VERSION}...${NC}"

# 备份旧版本（可选）
BACKUP_DIR="/usr/local/go-backup-$(date +%Y%m%d-%H%M%S)"
if [ -d "/usr/local/go" ]; then
    echo -e "${YELLOW}备份旧版本到 ${BACKUP_DIR}...${NC}"
    sudo mv /usr/local/go "$BACKUP_DIR"
    echo -e "${GREEN}备份完成${NC}"
fi

# 下载新版本
echo -e "${GREEN}下载 ${LATEST_VERSION}...${NC}"
curl -L -o "$DOWNLOAD_FILE" "$DOWNLOAD_URL"

# 验证下载
if [ ! -f "$DOWNLOAD_FILE" ]; then
    echo -e "${RED}下载失败！${NC}"
    exit 1
fi

# 安装新版本
echo -e "${GREEN}安装新版本...${NC}"
sudo tar -C /usr/local -xzf "$DOWNLOAD_FILE"

# 清理临时文件
rm -rf "$TEMP_DIR"

# 验证安装
echo -e "${GREEN}验证安装...${NC}"
NEW_VERSION=$(/usr/local/go/bin/go version)
echo -e "${GREEN}新版本: ${NEW_VERSION}${NC}"

# 检查 PATH
if ! echo "$PATH" | grep -q "/usr/local/go/bin"; then
    echo -e "${YELLOW}警告: /usr/local/go/bin 不在 PATH 中${NC}"
    echo -e "${YELLOW}请确保在 ~/.zshrc 或 ~/.bash_profile 中添加:${NC}"
    echo -e "export PATH=\$PATH:/usr/local/go/bin"
fi

echo -e "${GREEN}升级完成！${NC}"
echo -e "${YELLOW}如果旧版本备份在 ${BACKUP_DIR}，您可以稍后删除它${NC}"
