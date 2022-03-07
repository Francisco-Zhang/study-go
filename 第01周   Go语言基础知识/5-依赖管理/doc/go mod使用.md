### 基本使用

```shell
#首先开启 	 go mod: go env -w GO111MODULE=on
#新建文件夹  gomodtest
#初始化项目  go mod init gomodtest
#安装依赖 	 go get -u go.uber.org/zap   安装在gopath 的bin 目录下
#开发项目安装  go install ./...  会将本项目安装在 bin 目录下，go build 只编译不安装

# go.mod 文件会新增依赖版本配置

module gomodtest

go 1.17

require (
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
)

#有的包后面是一长串字符串，是因为go mod 没有发现确定的版本，按照最新提交的下载。字符串可能是提交记录id
#go.sum 文件会记录详细hash，确保下载的包文件没有被篡改
# indirect,当包被项目 import之后，go mod 文件中对应的包就不再显示 indirect 标识。

```



### 安装特定的版本

```go
go get -u go.uber.org/zap@v1.11
//go 程序不是直接根据import的路径去加载包，而是先到 go.mod 文件内查找版本，再去gopath目录对应版本目录加载文件。
```

### go mod tidy

使用该命令清洁后， go.sum 文件 会清除掉 go mod 中不再使用的包的信息。例如 安装了 v1.11的zap后，go mod文件依然有 原版本信息，使用该命令可以去除。

如果我们想升级zap版本， 使用 go get -u go.uber.org/zap 就是安装的最新版本。 



### 增加依赖的两种方法

1、使用 go get 命令   2、直接在代码中import,然后通过IDE 的 Sync 自动安装。



### go build

使用 go build ./...   命令会自动对项目下所有文件进行编译，可以借助该命令 完成项目所有文件 import 包的自动安装。

实际项目使用中  go get ./... 可以安装依赖