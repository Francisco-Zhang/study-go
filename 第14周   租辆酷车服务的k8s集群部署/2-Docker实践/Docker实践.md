## 1、构建Docker镜像

### 制作服务镜像

容器状态改变后，根据容器制作镜像命令：	docker commit    container_id	alpine_1

步骤：编译、选择基础镜像、拷贝编译后的可执行文件、设置环境

![1](img/1.png)

```shell
#生成打包文件
go install ./gateway/...
#进入 gopath目录 可以看到打包后的 gateway
ls ~/go/bin
#为了不污染工作机,并且防止本机打包的程序无法在容器中运行，所以一般在Docker镜像中编译
rm ~/go/bin/gateway
#运行带有go环境的镜像
docker run -it golang:1.15
#修改go evn 配置
go env
go env -w GO111MODDULE=on
go env -w GOPROXY=https://goproxy.cn,direct
#然后拷贝项目源码到src目录，然后go install


```

```shell
#下面是采用Dockerfile自动完成构建
#deployment目录新建gateway目录，然后新建文件Dockerfile

# 启动编译环境
FROM golang:1.15-alpine AS builder

# 配置编译环境
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 拷贝源代码到镜像中,当前目录为server
COPY . /go/src/coolcar/server

# 编译
WORKDIR /go/src/coolcar/server
RUN go install ./gateway/...

# 设置服务入口
ENTRYPOINT [ "/go/bin/gateway" ]
```

```shell
#在server目录下面, -t指定镜像名,最后面的.指定Dockerfile当前目录为具体哪个目录
docker build -t coolcar/gateway -f ../deployment/gateway/Dockerfile .
#服务不需要-it,服务是不会退出的
docker run coolcar/gateways
#在vscode使用插件 attach shell 或者 使用命令 docker exec -it  进入docker容器内部。
```

