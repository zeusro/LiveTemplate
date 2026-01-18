#!/bin/bash

# 加载公共配置和工具函数
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "${SCRIPT_DIR}/config.sh"
source "${SCRIPT_DIR}/utils.sh"

# 使用参数或默认值
URL="${1:-$DEFAULT_URL}"
OUTPUT_FILE="${2:-$DEFAULT_OUTPUT_FILE}"

# 检查Chrome
CHROME=$(check_chrome) || exit 1

echo "正在使用Chrome截屏..."
echo "URL: $URL"
echo "输出文件: $OUTPUT_FILE"
echo "窗口大小: ${WINDOW_WIDTH}x${WINDOW_HEIGHT}"

# 使用Chrome headless模式截屏
"$CHROME" \
    --headless=new \
    --disable-gpu \
    --window-size=${WINDOW_WIDTH},${WINDOW_HEIGHT} \
    --hide-scrollbars \
    --run-all-compositor-stages-before-draw \
    --virtual-time-budget=${VIRTUAL_TIME_BUDGET} \
    --screenshot="$OUTPUT_FILE" \
    "$URL" 2>&1 | grep -v "DevTools listening"

# 检查结果
if verify_screenshot "$OUTPUT_FILE" "$WINDOW_WIDTH"; then
    echo ""
    echo "✓ 截屏成功！"
    echo "文件已保存到: $OUTPUT_FILE"
    show_image_info "$OUTPUT_FILE"
else
    echo ""
    echo "✗ 截屏失败"
    echo "可能的原因："
    echo "1. Chrome未正确安装"
    echo "2. 网络连接问题，无法访问目标URL"
    echo "3. 页面加载时间过长"
    exit 1
fi
