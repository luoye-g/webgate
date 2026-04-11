# WebGate Frontend Server

这是一个基于Node.js和Express的前端服务器，用于提供静态页面服务并将API请求代理到后端服务器。

## 安装依赖

```bash
npm install
```

## 启动服务器

### 生产环境
```bash
npm start
```

### 开发环境（自动重启）
```bash
npm run dev
```

## 服务器配置

- 默认端口：3000（可通过环境变量PORT修改）
- API代理：将 `/api/*` 请求转发到 `http://localhost:8080`
- WebSocket代理：将 `/ws` 连接转发到 `ws://localhost:8080`

## 路由映射

- `/` → `templates/html/index.html`
- `/login` → `templates/html/login.html`
- `/sign-in` → `templates/html/sign-in.html`
- `/test` → `templates/html/test.html`
- `/assets/*` → `templates/assets/*`（静态文件）

## 环境变量

- `PORT`: 服务器监听端口（默认: 3000）

## 注意事项

确保后端服务器在 `localhost:8080` 运行，以便API和WebSocket代理正常工作。