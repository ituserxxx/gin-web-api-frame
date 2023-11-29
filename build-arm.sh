export GOOS=linux
export GOARCH=arm
export CGO_ENABLED=1
export CC=arm-linux-gnueabihf-gcc
export GOARM=7
go build -o zhihui-server-arm