#!/bin/bash
# 公共配置文件

# 默认URL（可以通过环境变量或参数覆盖）
DEFAULT_URL="${SCREENSHOT_URL:-https://www.sysu.edu.cn/news/info/1881/1394311.htm}"

# 输出目录（脚本所在目录）
OUTPUT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 默认输出文件名
DEFAULT_OUTPUT_FILE="${OUTPUT_DIR}/screenshot.png"

# 窗口尺寸
WINDOW_WIDTH="${SCREENSHOT_WIDTH:-1920}"
WINDOW_HEIGHT="${SCREENSHOT_HEIGHT:-4000}"

# 等待时间（毫秒）
VIRTUAL_TIME_BUDGET="${SCREENSHOT_TIMEOUT:-30000}"
