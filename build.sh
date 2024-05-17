rm -rf ./output

mkdir output
cp -r ./ssh ./output/
go build -o ./output/webgate