GO环境安装

1、安装编译环境，官网下载，直接安装即可

2、go 1.11开始支持模块化，通过命令打开

​		go env -w GO111MODUL=on

3、更改代理地址，下载包使用

​		go env -w GOPROXY=https://goproxy.cn,direct               direct表示前面地址无效的情况下，直接原地址拉取