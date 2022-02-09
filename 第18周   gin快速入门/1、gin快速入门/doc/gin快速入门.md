## 1、 gin的helloworld体验

### 官方文档地址

 https://gin-gonic.com/docs/quickstart/

Gin is a web framework written in Go (Golang),比较主流的go web 开发框架。

Beego是一个相对大而全的开发框架。

### 安装

go get -u github.com/gin-gonic/gin

go mod 模式下会自动根据import引用。

### 代码

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
    r.Run() // listen and serve on 0.0.0.0:8080，自定义端口（":8083"）
}
```

r.Get函数的第二个参数 就是一个   func(*Context) 的函数。上面和下面的这种写法是等效的。

```go

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	r := gin.Default()
	r.GET("/ping", pong)
	r.Run() // listen and serve on 0.0.0.0:8080
}
```

```go
c.JSON(200, map[string]interface{}{   //也可以这样写，gin.H就是一个map类型的别名。用于缩短代码长度。
		"message": "pong",
	})
```

### 使用get、post、put等http方法

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	// 使用默认中间件创建一个gin路由器
	// logger and recovery (crash-free) 中间件
	router := gin.Default()

	//restful 的开发中
	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.PATCH("/somePatch", patching)
	router.HEAD("/someHead", head)
	router.OPTIONS("/someOptions", options)

	// 默认启动的是 8080端口，也可以自己定义启动端口
	router.Run()
	// router.Run(":3000") for a hard coded port
}
```



## 2、 使用New和Default初始化路由器的区别

new 和 default 都可以创建路由。

使用default 可以自动创建 logger and recovery (crash-free) 中间件，服务端收到请求后logger 会打印日志。

recovery 可以处理异常 panic，会自动返回错误状态码。

new 创建的路由在方法抛出异常后，没有返回响应。

## 3、 gin的路由分组

路由分组是为了区分不同的服务，由于相同的服务路由前缀相同，为了减少路径重复，使用路由。路由分组内的路径会自动再前面加上分组路径。

```go

func main() {
	router := gin.Default()
	goodsGroup := router.Group("/goods")
	{//大括号是为了增强可读性，表示代码块，没有逻辑意义，也可以去掉。
		goodsGroup.GET("", goodsList)
		goodsGroup.GET("/:id/:action/add", goodsDetail) //获取商品id为1的详细信息 模式
		goodsGroup.POST("", createGoods)
	}

	router.Run(":8083")
}
```

## 4、 获取url中的变量

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/user/:name/:action/", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		c.String(http.StatusOK, "%s is %s", name, action)
	})
    
    // *号这种模式要慎用，平时业务场景很少用到
    //action 会取到name后面的全部路径作为参数
    //例如 user/tom/delete/a/b   ,action 的值为路径 /tom/delete/a/b，最开始有个`/`
    //一般操作文件路径有可能用到，会匹配到很多路径。后面不固定
    r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		c.String(http.StatusOK, "%s is %s", name, action)
	})


	r.Run(":8082") 
}
```

约束路径参数的类型,比如 id 只能是int 类型。

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Person struct {
	ID   int    `uri:"id" binding:"required"`
	Name string `uri:"name" binding:"required"`
}

func main() {
	router := gin.Default()
	router.GET("/:name/:id", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.Status(404)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"name": person.Name,
			"id":   person.ID,
		})
	})
	router.Run(":8083")
}
```

## 5、 获取get和post表单信息

### 获取Get参数

```go
func main() {
	router := gin.Default()

	// 匹配的url格式:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")  //默认值 Guest
		lastname := c.Query("lastname") // 是 c.Request.URL.Query().Get("lastname") 的简写

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.Run(":8080")
}
```

### 获取post参数

```go
func main() {
	router := gin.Default()

	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous") // 此方法可以设置默认值

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	router.Run(":8080")
}
```



### get、post混合

```go
POST /post?id=1234&page=1 HTTP/1.1
Content-Type: application/x-www-form-urlencoded

name=manu&message=this_is_great
func main() {
	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {

		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")

		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})
	router.Run(":8080")
}
```



## 6、 gin返回protobuf