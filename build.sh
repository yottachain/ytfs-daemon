out_file_name=ytfs-daemon

echo build linux
GOOS=linux GOARCH=amd64 go build -o ${out_file_name}-linux ./main.go
echo build darwin
GOOS=darwin GOARCH=amd64 go build -o ${out_file_name}-darwin ./main.go
echo build windows
GOOS=windows GOARCH=amd64 go build -o ${out_file_name}-windows ./main.go

rm -rf out*

mkdir out

mv ${out_file_name}* ./out

cp *lua ./out

tar -czvf ${out_file_name}.tar.gz ./out