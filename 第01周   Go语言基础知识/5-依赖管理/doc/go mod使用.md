```shell
#首先开启 	 go mod: go env -w GO111MODULE=on
#新建文件夹  gomodtest
#初始化项目  go mod init gomodtest
#安装依赖 	 go get -u go.uber.org/zap   安装在gopath 的bin 目录下

# go.mod 文件会新增依赖版本配置

module gomodtest

go 1.17

require (
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
)

#go.sum 文件会记录详细hash，确保下载的包文件没有被篡改
```

