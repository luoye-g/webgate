rm -rf ./webgate_product
rm webgate.zip

mkdir webgate_product
cp -r ./etc ./webgate_product/
cp -r ./frontend ./webgate_product/
cp run.sh ./webgate_product/
cd backend
GOOS=linux GOARCH=arm64 GOARCH=amd64 go build -o ../webgate_product/webgate
cd ..
zip -r webgate.zip ./webgate_product
rm -rf webgate_product