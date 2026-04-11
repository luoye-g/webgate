# 使用 Node.js 18 Alpine 作为基础镜像（轻量级）
FROM node:18-alpine

# 设置工作目录
WORKDIR /app

# 复制 package.json 和 package-lock.json
COPY package*.json ./

# 安装生产依赖（不安装 devDependencies）
RUN npm ci --only=production

# 复制应用程序源代码
COPY . .

# 暴露端口（默认 3001，可通过环境变量覆盖）
EXPOSE 3001

# 设置环境变量
ENV NODE_ENV=production
ENV PORT=3001

# 启动应用程序
CMD ["npm", "start"]