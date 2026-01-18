// 公共工具函数

const fs = require('fs');
const path = require('path');

/**
 * 获取Chrome路径
 * @returns {string} Chrome可执行文件路径
 */
function getChromePath() {
    if (process.platform === 'darwin') {
        const chromePath = '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome';
        if (fs.existsSync(chromePath)) {
            return chromePath;
        }
    }
    
    // 尝试从PATH中查找
    return 'google-chrome';
}

/**
 * 检查Chrome是否安装
 * @returns {string|null} Chrome路径，如果未找到则返回null
 */
function checkChrome() {
    const chromePath = getChromePath();
    
    // 对于macOS，检查文件是否存在
    if (process.platform === 'darwin' && chromePath.startsWith('/Applications')) {
        if (!fs.existsSync(chromePath)) {
            return null;
        }
    }
    
    return chromePath;
}

/**
 * 获取随机调试端口
 * @param {number} min - 最小端口号
 * @param {number} max - 最大端口号
 * @returns {number} 随机端口号
 */
function getRandomDebugPort(min = 9222, max = 10222) {
    return min + Math.floor(Math.random() * (max - min));
}

/**
 * 显示图片信息
 * @param {string} imageFile - 图片文件路径
 */
function showImageInfo(imageFile) {
    if (!fs.existsSync(imageFile)) {
        console.error('错误: 图片文件不存在');
        return;
    }
    
    const stats = fs.statSync(imageFile);
    console.log(`文件大小: ${(stats.size / 1024 / 1024).toFixed(2)} MB`);
    
    // 如果系统有sips命令（macOS），显示图片尺寸
    if (process.platform === 'darwin') {
        try {
            const { execSync } = require('child_process');
            const dimensions = execSync(`sips -g pixelWidth -g pixelHeight "${imageFile}" 2>/dev/null`, { encoding: 'utf8' });
            console.log('图片尺寸:');
            console.log(dimensions);
        } catch (e) {
            // sips命令不可用，跳过
        }
    }
}

/**
 * 验证截屏结果
 * @param {string} imageFile - 图片文件路径
 * @param {number} expectedWidth - 期望的宽度
 * @returns {boolean} 验证是否通过
 */
function verifyScreenshot(imageFile, expectedWidth = 1920) {
    if (!fs.existsSync(imageFile) || fs.statSync(imageFile).size === 0) {
        console.error('✗ 截屏失败：文件不存在或为空');
        return false;
    }
    
    // 在macOS上可以使用sips检查尺寸
    if (process.platform === 'darwin') {
        try {
            const { execSync } = require('child_process');
            const dimensions = execSync(`sips -g pixelWidth -g pixelHeight "${imageFile}" 2>/dev/null`, { encoding: 'utf8' });
            const heightMatch = dimensions.match(/pixelHeight:\s*(\d+)/);
            const widthMatch = dimensions.match(/pixelWidth:\s*(\d+)/);
            
            if (heightMatch && widthMatch) {
                const height = parseInt(heightMatch[1]);
                const width = parseInt(widthMatch[1]);
                
                if (height === 16384) {
                    console.log('⚠️  注意: 图片高度达到16384px（Chrome的最大限制）');
                    console.log('   如果页面实际高度超过此限制，底部内容可能被截断');
                }
                
                if (width === expectedWidth) {
                    console.log(`✓ 图片宽度: ${width}px（符合要求）`);
                }
            }
        } catch (e) {
            // sips命令不可用，跳过验证
        }
    }
    
    return true;
}

module.exports = {
    getChromePath,
    checkChrome,
    getRandomDebugPort,
    showImageInfo,
    verifyScreenshot
};
