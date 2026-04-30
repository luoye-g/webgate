const express = require('express');
const path = require('path');
const { createProxyMiddleware } = require('http-proxy-middleware');
const cookieParser = require('cookie-parser');
const morgan = require('morgan');

const app = express();
const PORT = (() => {
  const envPort = process.env.PORT;
  if (envPort) {
    const port = parseInt(envPort, 10);
    if (isNaN(port) || port < 1 || port > 65535) {
      console.warn(`警告: 环境变量 PORT=${envPort} 无效，使用默认端口 3001`);
      return 3001;
    }
    return port;
  }
  return 3001; // 使用3001避免与常见开发服务器冲突
})();

// 中间件
app.use(morgan('combined')); // 使用combined格式记录完整访问日志
app.use(cookieParser());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// 静态文件服务
app.use('/assets', express.static(path.join(__dirname, 'templates/assets')));
app.use('/js', express.static(path.join(__dirname, 'js')));
app.use('/node_modules', express.static(path.join(__dirname, 'node_modules')));


// 路由 - HTML页面
app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, 'templates/html/index.html'));
});

app.get('/login', (req, res) => {
  res.sendFile(path.join(__dirname, 'templates/html/login.html'));
});

app.get('/sign-in', (req, res) => {
  res.sendFile(path.join(__dirname, 'templates/html/sign-in.html'));
});

app.get('/profile', (req, res) => {
  res.sendFile(path.join(__dirname, 'templates/html/profile.html'));
});

app.get('/test', (req, res) => {
  res.sendFile(path.join(__dirname, 'templates/html/test.html'));
});

// 博客管理路由
app.get('/blog-list', (req, res) => {
  res.sendFile(path.join(__dirname, 'templates/html/blog-list.html'));
});

app.get('/blog-create', (req, res) => {
  res.redirect(301, '/blog-editor');
});

app.get('/blog-edit/:id', (req, res) => {
  res.redirect(301, `/blog-editor/${req.params.id}`);
});

// 统一的博客创建/编辑页（合并页面）
// /blog-editor           -> create 模式
// /blog-editor/new       -> create 模式
// /blog-editor/:id (数字) -> edit 模式
app.get('/blog-editor', (req, res) => {
  res.sendFile(path.join(__dirname, 'templates/html/blog-editor.html'));
});

app.get('/blog-editor/:id', (req, res) => {
  res.sendFile(path.join(__dirname, 'templates/html/blog-editor.html'));
});

app.get('/blog-view/:id', (req, res) => {
  res.sendFile(path.join(__dirname, 'templates/html/blog-view.html'));
});

// 兼容老路径：永久重定向到新的统一路径 /blog-view/:id
app.get('/blog-public-view/:id', (req, res) => {
  res.redirect(301, `/blog-view/${req.params.id}`);
});

// 启动服务器
const server = app.listen(PORT, () => {
  console.log(`服务器运行在 http://localhost:${PORT}`);
}).on('error', (err) => {
  if (err.code === 'EADDRINUSE') {
    console.error(`端口 ${PORT} 已被占用，请使用其他端口`);
    process.exit(1);
  } else {
    console.error('服务器启动失败:', err);
    process.exit(1);
  }
});

// 优雅关闭处理
process.on('SIGTERM', () => {
  console.log('收到 SIGTERM 信号，正在优雅关闭服务器...');
  server.close(() => {
    console.log('服务器已关闭');
    process.exit(0);
  });
});

process.on('SIGINT', () => {
  console.log('收到 SIGINT 信号，正在优雅关闭服务器...');
  server.close(() => {
    console.log('服务器已关闭');
    process.exit(0);
  });
});

module.exports = app;