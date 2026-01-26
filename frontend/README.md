# Athens 依赖管理前端

基于 Vite + React + TypeScript + Tailwind CSS + ShadCN UI 构建的 Athens Go 模块代理管理界面。

## 技术栈

- **框架**: React 19 + TypeScript
- **构建工具**: Vite 7
- **样式**: Tailwind CSS 4
- **UI 组件**: ShadCN UI (基于 Radix UI)
- **图标**: Lucide React

## 快速开始

### 安装依赖
```bash
npm install
```

### 开发模式
```bash
npm run dev
```
访问 `http://localhost:5173`

### 构建生产版本
```bash
npm run build
```

### 预览生产构建
```bash
npm run preview
```

## 项目结构

```
frontend/
├── src/
│   ├── components/
│   │   └── ui/              # ShadCN UI 组件
│   │       ├── button.tsx
│   │       └── card.tsx
│   ├── lib/
│   │   └── utils.ts         # 工具函数
│   ├── App.tsx              # 主应用组件
│   ├── main.tsx             # 入口文件
│   └── index.css            # 全局样式
├── public/                  # 静态资源
├── tailwind.config.js       # Tailwind 配置
├── vite.config.ts           # Vite 配置
└── package.json
```

## 功能特性

- ✅ 模块统计展示
- ✅ 模块列表查看
- ✅ 响应式设计
- ✅ 暗色模式支持
- 🚧 实时数据同步 (开发中)
- 🚧 模块搜索 (开发中)
- 🚧 模块详情页 (开发中)

## 与后端集成

前端需要连接到 Athens 后端 API。在 `src` 目录下创建 API 客户端：

```typescript
// src/api/client.ts
const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:3000'

export async function fetchModules() {
  const response = await fetch(`${API_BASE}/catalog`)
  return response.json()
}
```

在 `.env` 文件中配置后端地址：
```
VITE_API_BASE=http://localhost:3000
```

## 开发指南

### 添加新的 UI 组件

1. 在 `src/components/ui/` 下创建新组件
2. 使用 ShadCN UI 的设计模式
3. 导出并在页面中使用

### 样式定制

在 `tailwind.config.js` 中修改主题配置，或在 `src/index.css` 中调整 CSS 变量。

## License

与 Athens 项目保持一致
