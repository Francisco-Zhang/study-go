## 1、安装Go插件

插件中搜索go——安装 Go Team at Google

## 2、安装工具

view（查看）—— Command Palette（命令面板）—— 输入”go: install“——选择 go: install/update tools

全选安装这些工具（先配置国内镜像地址，安装会很快）。

## 3、添加server

因为vscode的插件上认项目的根目录的，所以需要添加文件夹server到工作区

这样server就是一个单独的项目，go 进行编译的时候只编译 server

使用 go mod init coolcar 进行项目初始化

新建文件，会有语言提示。如果弹窗 安装 gopls，则进行安装

首选项——设置——输入 go language —— 勾选 Use Language Server

Go: Format Tool 选择 goimports

新建文件，文件名随便写 hello.go

然后根据提示就可以 定义package , main函数。

## 4、运行

运行——以非调试模式运行——环境选择 go ——运行结果在 debug console 里

或者

点击底部 “选择并启动调试配置”——添加配置 server—— 选择第一个 Launch Package
