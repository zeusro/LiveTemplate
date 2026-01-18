// 公共配置文件

const path = require('path');

// 默认配置
const config = {
    // 默认URL（可以通过环境变量覆盖）
    defaultUrl: process.env.SCREENSHOT_URL || 'https://www.sysu.edu.cn/news/info/1881/1394311.htm',
    
    // 输出目录（脚本所在目录）
    outputDir: __dirname,
    
    // 默认输出文件名
    defaultOutputFile: 'screenshot.png',
    
    // 窗口尺寸
    windowWidth: parseInt(process.env.SCREENSHOT_WIDTH) || 1920,
    windowHeight: parseInt(process.env.SCREENSHOT_HEIGHT) || 4000,
    
    // 等待时间（毫秒）
    virtualTimeBudget: parseInt(process.env.SCREENSHOT_TIMEOUT) || 30000,
    
    // 调试端口范围
    debugPortMin: 9222,
    debugPortMax: 10222
};

module.exports = config;
