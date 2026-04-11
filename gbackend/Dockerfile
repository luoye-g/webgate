# 使用最小的alpine镜像作为运行环境
FROM alpine:latest

# 设置工作目录
WORKDIR /data/app


# 复制配置文件
COPY ./blog_backend ./

# 暴露端口
EXPOSE 9001

# 运行应用
CMD ["./blog_backend"]