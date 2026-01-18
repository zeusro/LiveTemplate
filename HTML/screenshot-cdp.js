#!/usr/bin/env node

const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');
const config = require('./config');
const { checkChrome, getRandomDebugPort, showImageInfo, verifyScreenshot } = require('./utils');

// 使用参数或默认值
const URL = process.argv[2] || config.defaultUrl;
const OUTPUT_FILE = path.join(config.outputDir, process.argv[3] || config.defaultOutputFile);

// 检查Chrome
const chromePath = checkChrome();
if (!chromePath) {
    console.error('错误: 未找到Chrome浏览器');
    console.error('请确保已安装Google Chrome或Chromium');
    process.exit(1);
}

const DEBUG_PORT = getRandomDebugPort(config.debugPortMin, config.debugPortMax);

console.log('正在使用Chrome DevTools Protocol截屏...');
console.log(`URL: ${URL}`);
console.log(`输出文件: ${OUTPUT_FILE}`);
console.log(`窗口大小: ${config.windowWidth}x${config.windowHeight}\n`);

// 启动Chrome
const chrome = spawn(chromePath, [
    '--headless=new',
    '--disable-gpu',
    `--remote-debugging-port=${DEBUG_PORT}`,
    `--window-size=${config.windowWidth},${config.windowHeight}`,
    '--hide-scrollbars',
    URL
], {
    stdio: 'ignore'
});

(async function() {
    try {
        // 等待Chrome启动
        await new Promise(resolve => setTimeout(resolve, 3000));

        // 获取页面标签
        const tabsResponse = await fetch(`http://localhost:${DEBUG_PORT}/json`);
        const tabs = await tabsResponse.json();
        const tabId = tabs[0]?.id;

        if (!tabId) {
            console.error('错误: 无法获取页面标签');
            chrome.kill();
            process.exit(1);
        }

        console.log('正在加载页面并滚动到底部...');

        // 执行滚动到底部的脚本
        const scrollScript = `
        (async function() {
            // 等待页面加载
            await new Promise(resolve => {
                if (document.readyState === 'complete') {
                    resolve();
                } else {
                    window.addEventListener('load', resolve);
                }
            });
            
            // 等待动态内容加载
            await new Promise(resolve => setTimeout(resolve, 2000));
            
            // 滚动到底部
            let lastHeight = 0;
            let currentHeight = document.body.scrollHeight;
            let attempts = 0;
            const maxAttempts = 50;
            
            while (attempts < maxAttempts) {
                window.scrollTo(0, document.body.scrollHeight);
                await new Promise(resolve => setTimeout(resolve, 300));
                
                currentHeight = document.body.scrollHeight;
                if (currentHeight === lastHeight) {
                    break;
                }
                lastHeight = currentHeight;
                attempts++;
            }
            
            // 最后滚动到最底部
            window.scrollTo(0, document.body.scrollHeight);
            await new Promise(resolve => setTimeout(resolve, 1000));
            
            return document.body.scrollHeight;
        })();
        `;

        // 执行滚动
        const evaluateResponse = await fetch(`http://localhost:${DEBUG_PORT}/json/runtime/evaluate`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ expression: scrollScript })
        });

        await evaluateResponse.json();
        await new Promise(resolve => setTimeout(resolve, 5000));

        console.log('正在截屏...');

        // 使用Page.captureScreenshot
        const screenshotResponse = await fetch(`http://localhost:${DEBUG_PORT}/json/page/captureScreenshot`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                format: 'png',
                fromSurface: true,
                captureBeyondViewport: true
            })
        });

        const screenshotData = await screenshotResponse.json();

        if (screenshotData.data) {
            // 保存图片
            const imageBuffer = Buffer.from(screenshotData.data, 'base64');
            fs.writeFileSync(OUTPUT_FILE, imageBuffer);
            
            console.log('\n✓ 截屏成功！');
            console.log(`文件已保存到: ${OUTPUT_FILE}`);
            showImageInfo(OUTPUT_FILE);
            
            if (verifyScreenshot(OUTPUT_FILE, config.windowWidth)) {
                console.log('');
            }
        } else {
            console.error('错误: 无法获取截屏数据');
            console.error('响应:', screenshotData);
            chrome.kill();
            process.exit(1);
        }

        // 关闭Chrome
        chrome.kill();
    } catch (error) {
        console.error('发生错误:', error);
        chrome.kill();
        process.exit(1);
    }
})();
