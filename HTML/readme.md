# HTML 工具集

## IE强制最新，多核强制极速

```HTML
<meta http-equiv="X-UA-Compatible" content="IE=edge" />
<meta name="renderer" content="webkit" /> 
```

## Chrome 截屏工具

本目录提供多个Chrome截屏脚本，支持使用Chrome headless模式对网页进行截屏。

### 配置文件

#### `config.sh` / `config.js`
公共配置文件，包含默认URL、窗口尺寸、输出路径等配置。可以通过环境变量覆盖：

**环境变量：**
- `SCREENSHOT_URL` - 要截屏的URL（默认：https://www.sysu.edu.cn/news/info/1881/1394311.htm）
- `SCREENSHOT_WIDTH` - 窗口宽度（默认：1920）
- `SCREENSHOT_HEIGHT` - 窗口高度（默认：4000）
- `SCREENSHOT_TIMEOUT` - 等待时间，毫秒（默认：30000）

#### `utils.sh` / `utils.js`
公共工具函数库，提供：
- `get_chrome_path()` / `getChromePath()` - 获取Chrome路径
- `check_chrome()` / `checkChrome()` - 检查Chrome是否安装
- `show_image_info()` / `showImageInfo()` - 显示图片信息
- `verify_screenshot()` / `verifyScreenshot()` - 验证截屏结果

### Shell 脚本

#### `screenshot.sh`
基础截屏脚本，使用Chrome headless模式截取网页。

**用法：**
```bash
# 使用默认URL和输出文件
./screenshot.sh

# 指定URL
./screenshot.sh "https://example.com"

# 指定URL和输出文件
./screenshot.sh "https://example.com" "output.png"

# 使用环境变量
SCREENSHOT_URL="https://example.com" SCREENSHOT_WIDTH=1920 ./screenshot.sh
```

**参数：**
1. URL（可选）- 要截屏的网页URL，默认使用config.sh中的DEFAULT_URL
2. 输出文件（可选）- 输出图片路径，默认为当前目录下的screenshot.png

**特点：**
- 窗口大小：1920x4000（可通过环境变量配置）
- 自动等待页面加载完成
- 支持显示图片尺寸信息

### Node.js 脚本

#### `screenshot-cdp.js`
使用Chrome DevTools Protocol (CDP) 进行截屏的Node.js脚本，支持真正的滚动到底部并截取完整页面。

**用法：**
```bash
# 使用默认URL和输出文件
node screenshot-cdp.js

# 指定URL
node screenshot-cdp.js "https://example.com"

# 指定URL和输出文件
node screenshot-cdp.js "https://example.com" "output.png"

# 使用环境变量
SCREENSHOT_URL="https://example.com" SCREENSHOT_WIDTH=1920 node screenshot-cdp.js
```

**参数：**
1. URL（可选）- 要截屏的网页URL，默认使用config.js中的defaultUrl
2. 输出文件（可选）- 输出图片文件名（不含路径），默认为screenshot.png，保存在脚本目录

**特点：**
- 使用Chrome DevTools Protocol进行精确控制
- 自动滚动到页面底部，确保所有内容加载
- 支持动态内容加载
- 使用WebSocket与Chrome通信

**依赖：**
- Node.js (v18+)
- Chrome浏览器

**工作原理：**
1. 启动Chrome headless模式，启用远程调试端口
2. 通过WebSocket连接到Chrome DevTools Protocol
3. 执行JavaScript代码滚动到页面底部
4. 等待所有动态内容加载完成
5. 使用Page.captureScreenshot API截取完整页面

**代码示例：**
```javascript
const { checkChrome, getRandomDebugPort } = require('./utils');
const config = require('./config');

// 检查Chrome
const chromePath = checkChrome();
if (!chromePath) {
    console.error('错误: 未找到Chrome浏览器');
    process.exit(1);
}

// 获取随机调试端口
const debugPort = getRandomDebugPort();
```

### 使用示例

#### 基本截屏
```bash
# 截取默认URL
./screenshot.sh

# 截取指定URL
./screenshot.sh "https://www.example.com"
```

#### 使用Node.js脚本
```bash
# 使用CDP方法截屏（支持真正的滚动）
node screenshot-cdp.js "https://www.example.com" "output.png"
```

#### 自定义配置
```bash
# 设置自定义窗口大小
SCREENSHOT_WIDTH=2560 SCREENSHOT_HEIGHT=2000 ./screenshot.sh

# 设置自定义URL和超时时间
SCREENSHOT_URL="https://www.sysu.edu.cn/news/info/1881/1394311.htm" SCREENSHOT_TIMEOUT=3000  SCREENSHOT_HEIGHT=2000 ./screenshot.sh
```

### 文件说明

| 文件 | 类型 | 窗口高度 | 特点 | 适用场景 |
|------|------|---------|------|---------|
| `config.sh` | Shell配置 | - | Bash脚本的公共配置文件 | - |
| `config.js` | JS配置 | - | Node.js脚本的公共配置文件 | - |
| `utils.sh` | Shell工具 | - | Bash脚本的公共工具函数 | - |
| `utils.js` | JS工具 | - | Node.js脚本的公共工具函数 | - |
| `screenshot.sh` | Shell脚本 | 4000px | 基础截屏，快速简单 | 短页面、快速截屏 |
| `screenshot-cdp.js` | Node.js脚本 | 4000px | 使用CDP真正滚动，支持动态内容 | 需要滚动加载的页面 |

**功能对比：**

1. **screenshot.sh** - 基础截屏
   - 窗口高度：4000px（可配置）
   - 特点：简单快速，适合大多数场景
   - 使用：`./screenshot.sh [URL] [输出文件]`

2. **screenshot-cdp.js** - 高级截屏（Node.js）
   - 窗口高度：4000px（可配置）
   - 特点：真正滚动到底部，支持动态内容加载
   - 使用：`node screenshot-cdp.js [URL] [输出文件]`
   - 优势：可以处理需要滚动才能加载的内容

### 常见问题

**Q: 为什么截屏高度只有4000px？**
A: `screenshot.sh`默认窗口高度为4000px。如需截取完整页面或处理动态内容，请使用`screenshot-cdp.js`。

**Q: 页面被截断了怎么办？**
A: Chrome headless的最大高度限制为16384px。如果页面超过此高度，可以考虑：
1. 使用Puppeteer等工具
2. 分块截屏后拼接
3. 使用`screenshot-cdp.js`尝试滚动截屏

**Q: 如何修改默认URL？**
A: 可以通过以下方式：
1. 修改`config.sh`或`config.js`中的默认URL
2. 使用环境变量`SCREENSHOT_URL`
3. 在命令行中传递URL参数

**Q: Node.js脚本报错"Cannot find module 'ws'"？**
A: 需要安装ws模块：
```bash
cd HTML
npm install ws
```

### 技术说明

- **Chrome headless模式**：无界面运行Chrome，适合自动化任务
- **Chrome DevTools Protocol (CDP)**：Chrome的调试协议，支持精确控制浏览器行为
- **最大高度限制**：Chrome的`--screenshot`参数最大高度为16384px
- **虚拟时间预算**：`--virtual-time-budget`参数用于等待页面加载和渲染
