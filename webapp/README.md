# md2html Webapp

Vue 3 前端应用，为在线工具箱平台提供用户界面，首个工具为 **Markdown → HTML 转换器**。

## 技术栈

- **Vue 3** + **TypeScript** + **Vite** 8
- **Pinia** 状态管理
- **Vue Router** 4 路由
- **VueUse** 工具函数库
- **marked.js** (Web Worker) — Markdown 解析
- **highlight.js** — 代码高亮
- **Mermaid** 10 — 流程图渲染
- **KaTeX** — 数学公式渲染

## 目录结构

```
webapp/
├── index.html                    # HTML 入口（引入 CDN 资源）
├── vite.config.ts                # Vite 配置
├── .env / .env.production        # 环境变量
├── src/
│   ├── main.ts                   # 应用入口
│   ├── App.vue                   # 根组件
│   ├── router/index.ts           # 路由配置
│   ├── stores/                   # Pinia Store
│   │   ├── app.ts                # 全局状态（主题、导航）
│   │   ├── auth.ts               # 认证状态
│   │   ├── editor.ts             # 编辑器状态（内容、行号、统计）
│   │   ├── preview.ts            # 预览状态（渲染结果、视图模式）
│   │   └── history.ts            # 历史记录状态
│   ├── views/                    # 页面视图
│   │   ├── HomeView.vue          # 工具箱首页
│   │   ├── Md2HtmlView.vue       # Markdown 转换器主页面
│   │   ├── HistoryView.vue       # 历史记录
│   │   ├── LoginView.vue         # 登录
│   │   └── SettingsView.vue      # 设置
│   ├── components/
│   │   └── md2html/              # md2html 工具组件
│   │       ├── ToolbarBar.vue    # 顶部工具栏
│   │       ├── FindBox.vue       # 查找替换（备用独立组件）
│   │       └── DropZone.vue      # 拖拽导入（备用独立组件）
│   ├── composables/              # Vue Composables（核心逻辑）
│   │   ├── index.ts              # 统一导出
│   │   ├── useMarkdownRender.ts  # 渲染管线编排
│   │   ├── useMarkdownWorker.ts  # Web Worker 管理
│   │   ├── useMermaid.ts         # Mermaid 渲染 + 下载
│   │   ├── useKaTeX.ts           # KaTeX 数学公式
│   │   ├── useHighlight.ts       # 代码高亮 + 复制
│   │   ├── useScrollSync.ts      # 编辑器/预览滚动同步
│   │   ├── usePanelResize.ts     # 面板拖拽分栏
│   │   ├── useFileIO.ts          # 文件导入导出
│   │   ├── useFindReplace.ts     # 查找替换
│   │   ├── useToast.ts           # Toast 通知
│   │   └── useExportHtml.ts      # 导出完整 HTML
│   ├── services/
│   │   └── api/                  # API 客户端层
│   │       ├── client.ts         # fetch 封装（拦截器、baseURL）
│   │       ├── auth.ts           # 认证 API
│   │       └── md2html.ts        # md2html 工具 API
│   ├── assets/styles/
│   │   ├── variables.css         # CSS 变量（主题色、间距、字体）
│   │   └── base.css              # 全局基础样式
│   ├── utils/constants.ts        # 常量定义
│   └── types/                    # TypeScript 类型（预留）
```

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

启动后访问 `http://localhost:5173`。

### 生产构建

```bash
npm run build
```

输出到 `dist/` 目录，可直接部署到 Nginx 或静态文件服务器。

### 预览构建结果

```bash
npm run preview
```

## 路由设计

| 路径 | 页面 | 认证 | 说明 |
|------|------|------|------|
| `/` | HomeView | - | 工具箱首页 |
| `/md2html` | Md2HtmlView | - | Markdown 转换器 |
| `/json-formatter` | JsonFormatterView | - | JSON 格式化 |
| `/history` | HistoryView | JWT | 转换历史 |
| `/login` | LoginView | - | 登录 |
| `/settings` | SettingsView | - | 设置 |

> 未来新增工具时，在 `router/index.ts` 添加路由并在 `views/` 下创建对应页面即可。

## 核心功能

### Markdown 渲染管线

```
用户输入 → 120ms 防抖 → Web Worker (marked.js)
    → postMessage → 更新 Store
    → 后处理：Heading ID → Mermaid → KaTeX → 代码高亮/复制
    → localStorage 缓存
```

**设计原则**：本地渲染优先，后端 API 作为兼容能力。

### Composable 职责

| Composable | 职责 |
|------------|------|
| `useMarkdownRender` | 编排整个渲染管线，串联 Worker → 后处理 |
| `useMarkdownWorker` | Web Worker 生命周期、防抖、渲染令牌 |
| `useMermaid` | 检测 mermaid 代码块、渲染 SVG、添加下载按钮 |
| `useKaTeX` | 检测数学公式、调用 KaTeX 渲染 |
| `useHighlight` | 代码高亮、语言标签、复制按钮、Heading ID |
| `useScrollSync` | 编辑器与预览的双向滚动同步 |
| `usePanelResize` | 拖拽分栏、比例持久化到 localStorage |
| `useFileIO` | 文件导入（.md/.txt）、拖拽导入 |
| `useFindReplace` | 查找下一个、替换、全部替换 |
| `useExportHtml` | 导出完整 HTML（含目录、Mermaid、Lightbox） |
| `useToast` | Toast 通知 |

### 主题

通过 CSS 变量实现 Light/Dark 主题切换，变量定义在 `assets/styles/variables.css`。主题状态持久化到 localStorage。

## 环境变量

| 变量 | 开发默认值 | 生产值 | 说明 |
|------|-----------|--------|------|
| `VITE_API_BASE` | `http://localhost:8000` | `/` | 后端 API 地址 |

开发环境配置在 `.env`，生产环境配置在 `.env.production`。

## API 降级策略

前端优先使用本地 Web Worker 渲染。当需要历史记录、主题同步等后端功能时，调用 API；若 API 不可用，自动降级为纯本地模式。

## Nginx 部署参考

```nginx
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;

    # SPA history 模式
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://backend:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## License

MIT
