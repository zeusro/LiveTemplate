# Part 3: 简易服务器状态监控看板

Go 后端 + Vue 3 + TypeScript 前端，SSE 实时推送 CPU/内存模拟数据，支持动态修改推送频率。

## 功能

- **后端 (Go)**  
  - `GET /api/metrics`：单次获取随机 CPU/内存  
  - `GET /api/stream`：SSE 流，按当前间隔持续推送  
  - `PUT /api/interval?seconds=5`：修改推送间隔（1–60 秒）  
- **前端 (Vue 3 + TS)**  
  - 连接 SSE，接收实时数据  
  - ECharts 绘制 CPU/内存实时曲线  
  - 下拉框调节推送间隔并即时生效  

## 运行

### 1. 启动后端

```bash
cd backend
go run main.go
```

后端监听 `http://localhost:8080`。

### 2. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端默认 `http://localhost:5173`，Vite 会将 `/api` 代理到 `http://localhost:8080`。

### 3. 使用

浏览器打开 `http://localhost:5173`，即可看到实时曲线与间隔控制。

## 目录结构

```
monitor-dashboard/
├── backend/
│   └── main.go       # Go SSE + 间隔 API
├── frontend/
│   ├── src/
│   │   ├── App.vue    # 图表 + 间隔控制
│   │   └── main.ts
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
└── README.md
```

## AI 使用报告（模板）

- **使用的 AI 模型**：_____________  
- **关键 Prompt 示例**：_____________  
- **AI 生成代码中的 Bug 及修复**：_____________  

（作答时请按测评要求填写并附在提交文档中。）
