rm blog_backend
GOOS=linux GOARCH=amd64  go build -o blog_backend main.go

tag=v0.0.8
docker.exe build -f Dockerfile -t blog_backend:${tag} .

docker.exe tag  blog_backend:${tag} ccr.ccs.tencentyun.com/boxgeng/blog_backend:${tag}
docker.exe push ccr.ccs.tencentyun.com/boxgeng/blog_backend:${tag}
