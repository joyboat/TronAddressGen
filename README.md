# TronAddressGen
波场的靓号生成器

## 如何自己编译

1. 安装go环境
2. 项目下运行 `go mod tidy`
3. 编译 `go build -o TronAddressGen.exe`
4. 项目目录下会有一个 `TronAddressGen.exe`
5. 在命令行下，执行 `TronAddressGen.exe -endTimes=6 -numAddr=20 -numWorker=20`
6. 跑吧