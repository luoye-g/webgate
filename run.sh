# step1 拷贝新配置到nginx
cp ./etc/nginx/nginx.conf /home/luoye/nginx/conf/nginx.conf
# step2 重新启动nginx
sudo /home/luoye/nginx/sbin/nginx -s reload
# step3 启动go服务
nohup ./webgate > ./log.txt &

/bin/bash 
cd /home/ubuntu/product/
./backend-linux