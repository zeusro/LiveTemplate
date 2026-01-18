#!/bin/bash
# 公共工具函数

# 获取Chrome路径
get_chrome_path() {
    if [ -f "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" ]; then
        echo "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
    elif command -v google-chrome &> /dev/null; then
        echo "google-chrome"
    elif command -v chromium &> /dev/null; then
        echo "chromium"
    else
        echo ""
    fi
}

# 检查Chrome是否安装
check_chrome() {
    local chrome_path=$(get_chrome_path)
    if [ -z "$chrome_path" ]; then
        echo "错误: 未找到Chrome浏览器" >&2
        echo "请确保已安装Google Chrome或Chromium" >&2
        return 1
    fi
    echo "$chrome_path"
}

# 显示图片信息
show_image_info() {
    local image_file="$1"
    if [ ! -f "$image_file" ] || [ ! -s "$image_file" ]; then
        return 1
    fi
    
    echo "文件信息:"
    ls -lh "$image_file"
    
    if command -v sips &> /dev/null; then
        echo ""
        echo "图片尺寸:"
        sips -g pixelWidth -g pixelHeight "$image_file" 2>/dev/null | grep -E "pixelWidth|pixelHeight"
    fi
}

# 验证截屏结果
verify_screenshot() {
    local image_file="$1"
    local expected_width="${2:-1920}"
    
    if [ ! -f "$image_file" ] || [ ! -s "$image_file" ]; then
        echo "✗ 截屏失败：文件不存在或为空"
        return 1
    fi
    
    if command -v sips &> /dev/null; then
        local dimensions=$(sips -g pixelWidth -g pixelHeight "$image_file" 2>/dev/null | grep -E "pixelWidth|pixelHeight")
        local height=$(echo "$dimensions" | grep pixelHeight | awk '{print $2}')
        local width=$(echo "$dimensions" | grep pixelWidth | awk '{print $2}')
        
        if [ "$height" -eq 16384 ]; then
            echo "⚠️  注意: 图片高度达到16384px（Chrome的最大限制）"
            echo "   如果页面实际高度超过此限制，底部内容可能被截断"
        fi
        
        if [ "$width" -eq "$expected_width" ]; then
            echo "✓ 图片宽度: ${width}px（符合要求）"
        fi
    fi
    
    return 0
}
