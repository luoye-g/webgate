# 使用 Node.js 20 Alpine 作为基础镜像（轻量级）
FROM node:20-alpine

# 设置工作目录
WORKDIR /app

# 复制应用程序源代码
COPY . .

# 删除可能被复制进来的本地 node_modules，然后安装生产依赖
# 禁用 strict-ssl 以规避 SSL 证书验证问题（如证书过期等）
RUN rm -rf node_modules \
    && npm config set strict-ssl false \
    && npm config set registry https://registry.npmjs.org/ \
    && npm install --omit=dev --no-audit --no-fund

# 暴露端口（默认 3001，可通过环境变量覆盖）
EXPOSE 3001

# 设置环境变量
ENV NODE_ENV=production
ENV PORT=3001

# 启动应用程序
CMD ["npm", "start"]