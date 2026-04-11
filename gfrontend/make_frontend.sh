tag=v0.0.4
docker.exe build -f Dockerfile -t blog_frontend:${tag} .

docker.exe tag  blog_frontend:${tag} ccr.ccs.tencentyun.com/boxgeng/blog_frontend:${tag}
docker.exe push ccr.ccs.tencentyun.com/boxgeng/blog_frontend:${tag}