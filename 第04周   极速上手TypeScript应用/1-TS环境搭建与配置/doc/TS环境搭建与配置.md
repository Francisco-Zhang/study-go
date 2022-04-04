## 1、 node和typescript的安装

### node、npm安装

官网下载安装包安装，会自动添加环境变量



### typescript安装

一般不需要安装，可以使用使用 npx tsc 命令允许 ts

npx tsc: 如果本地有 tsc 包就运行，没有就下载，然后临时安装运行 tsc,运行结束后删除临时安装的tsc。

```sh
#查看仓库地址
npm get registry
#设置成国内地址
npm set registry  https://registry.npm.taobao.org
```



我自己尝试的时候，发现需要先安装 typescript

安装 typescript：

```shell
npm install -g typescript  #如果提示权限不足，加sudo
```

安装完成后我们可以使用 **tsc** 命令来执行 TypeScript 的相关代码，以下是查看版本号：

```sh
$ tsc -v
Version 4.6.3
```

之后可以正常使用 npx tsc 命令。



## 2、 typescript小程序代码的生成

新建小程序项目，语言选择 typescript

<img src="img/1.png"  />

新版本开发者工具针对ts编译进行了优化：

工具通过对内置的编译流程进行优化，以编译插件的方式，改进了对 typescript 项目支持。

1. 相比起之前 Typescript 项目中会同时存在 ts 文件和 js 文件，新的模板只需要创建 ts 文件即可，无需再生成同名的 js 文件。
2. 新的模板无需在每次编译前执行 npm run tsc 命令。
