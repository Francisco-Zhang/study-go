## 1、集群的配置及版本

### 简介

kind用于本机进行k8s集群练习，k8s的版本由云厂商决定。![1](img/1.png)

### 腾讯云安装

云产品-->集群

运行时组件选 containerd，docker已经被弃用。

![2](img/2.png)



## 2、kubectl的安装

https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/    kubectl 安装包文档说明页面。

下载：https://dl.k8s.io/release/v1.18.0/bin/windows/amd64/kubectl.exe      版本号可以修改。

然后将可执行文件目录添加PATH目录。



## 3、用kind来配置本地集群

还有其他的工具，例如minikube

Kind 是 kubernetes in docker 的简写。kubernetes in docker  is  not using docker。

官网：https://kind.sigs.k8s.io/![3](img/3.png)

[release notes]   https://github.com/kubernetes-sigs/kind/releases  有kind 与 k8s版本对应关系。

If you have [go](https://golang.org/) ([1.17+](https://golang.org/doc/devel/release.html#policy)) and [docker](https://www.docker.com/) installed `go install sigs.k8s.io/kind@v0.11.1 && kind create cluster` is all you need!

安装 Kind 之前需要先有 Go 和 Docker环境。Kind 镜像是运行在docker服务之上的。

在项目的 go.mod 同级目录 使用 go get命令安装kind,只有是在mod项目中才能使用go get 命令

```shell
go get sigs.k8s.io/kind@v0.8.0
```

Kind 会被安装在gopath 的bin 目录下面，将bin目录添加到Path 环境变量，就可以直接在控制台使用 kind命令。

```shell
#创建集群
kind create cluster
```

整个结构是 kind 镜像里边 运行的 k8s 集群。![4](img/4.png)

为了系统运行的更流程，建议给Docker服务分配更多的CPU,内存空间。

## 4、用kubectl来连接集群

查看集群信息命令：kubectl  cluster-info  --context kind-kind  或者 没有context,直接使用  kubectl  cluster-info

查看版本命令：	kubectl   version

腾讯云 k8s集群访问

1、基本信息里边开通外网访问权限

2、下载配置文件 cls-f5ukhgwj-config

3、重命名为  qcloud-coolcar-config

4、 export KUBECONFIG = ~/Downloads/qcloud-coolcar-config  

​		windows 用 set %KUBECONFIG%=~/Downloads/qcloud-coolcar-config  



vscode 插件管理：添加config，可以查看k8s信息

![5](img/5.png)





查看 kind 集群

kind get kubeconfig  将内容保存到文件，然后使用插件访问。



设置环境变量后，kubctl访问腾讯云

export KUBECONFIG = ~/Downloads/qcloud-coolcar-config  

kubectl  cluster-info

清除后，访问本机kind

unset   KUBECONFIG 

kubectl  cluster-info