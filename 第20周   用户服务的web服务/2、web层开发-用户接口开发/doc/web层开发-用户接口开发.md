## 1 、表单验证的初始化

## 2、 自定义mobile验证器

```go
type PassWordLoginForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required"` //手机号码格式有规范可寻， 自定义validator
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
}


func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	//使用正则表达式判断是否合法
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok{
		return false
	}
	return true
}

//注册验证器
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    _ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
    _ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
        return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
    }, func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("mobile", fe.Field())
        return t
    })
}
```

## 3 、登录逻辑完善

```go
func PassWordLogin(c *gin.Context) {
	//表单验证
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	//拨号连接用户grpc服务器
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务】失败",
			"msg", err.Error())
	}
	//登录的逻辑
	userSrvClient := proto.NewUserClient(conn)
	if rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		//只是查询到用户了而已，并没有检查密码
		if passRsp, pasErr := userSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); pasErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败",
			})
		} else {
			if passRsp.Success {
				c.JSON(http.StatusOK, map[string]string{
					"msg": "登录成功",
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登录失败",
				})
			}

		}
	}
}
```

## 4、 session机制在微服务下的问题



### ![2](img/2.PNG)



微服务数据库隔离，session机制失效，解决办法是 session 放到 redis 中。

jwt,  json web tocken 可以做到不存储 json ，仍然可以验证。



## 5、 json web token的认证机制

密钥一定不能泄露，密钥是服务端做验证用的，密钥只能是服务器知道。



## 6、 集成jwt到gin中

key 随机生成 网址：https://suijimimashengcheng.bmcx.com/



## 7、 给url添加登录权限验证

```go
func InitUserRouter(Router *gin.RouterGroup) {
    //如果是user这一组url都加权限，则是 Router.Group("user").Use(middlewares.JWTAuth())
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	{	//单个接口加权限
		UserRouter.GET("list", middlewares.JWTAuth(),middlewares.IsAdminAuth(), api.GetUserList) 
		UserRouter.POST("pwd_login", api.PassWordLogin)
	}
}


func IsAdminAuth() gin.HandlerFunc{
	return func(ctx *gin.Context){
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg":"无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
```



## 8、 如何解决前后端的跨域问题

跨域的问题 - 可以后端解决, 也可以前端来解决(比如vue 的proxy)，前端解决本质上是前端自己充当服务端然后转发请求。

### 跨域演示

简单请求，不会出现跨域提示。去掉 dataType、beforeSend就会变成简单请求。

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <script src="http://libs.baidu.com/jquery/2.0.0/jquery.min.js"></script>
</head>
<body>
    <button type="button" id="query">请求数据</button>
    <div id="content" style="background-color: aquamarine; width: 300px;height: 500px"></div>
</body>
<script type="text/javascript">
    $("#query").click(function () {
        $.ajax(
            {
                url:"http://127.0.0.1:8021/u/v1/user/list",
                dataType: "json",  //去掉才能变成简单请求
                type: "get",
                //header中添加数据就属于非简单请求，这样浏览器就会自动发送OPTIONS请求，然后出现权限校验
                beforeSend: function(request) {
                  request.setRequestHeader("x-token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiTmlja05hbWUiOiJib2JieTAiLCJBdXRob3JpdHlJZCI6MSwiZXhwIjoxNjQ3MjQ5ODI4LCJpc3MiOiJpbW9vYyIsIm5iZiI6MTY0NDY1NzgyOH0.b58yDz3o3WaoF4Dff0FaZyU4MqbEQrya4nnKsmU4-qU")
                },
                success: function (result) {
                    console.log(result.data);
                    $("#content").text(result.data)
                },
                error: function (data) {
                    alert("请求出错")
                }
            }
        );
    });
</script>
</html>
```

服务端跨域提示：

本地浏览器打开html，点击按钮，出现跨域请求。

```http
request URL: http://127.0.0.1:8021/u/v1/user/list
Referrer Policy: strict-origin-when-cross-origin
```



### 什么情况下会发送 OPTIONS 请求?

当一个请求跨域且不是简单请求时就会发送 OPTIONS 请求

这里要注意两个问题：

1. 请求跨域

- 一般常见形式是 CORS 跨域；
- nginx 代理转发对浏览器而言没有跨域；

2. 简单请求：满足下列条件的请求可以视为简单请求

- 使用下列方法之一：
  GET
  HEAD
  POST

- Content-Type 的值仅限于下列三者之一：
  text/plain
  multipart/form-data
  application/x-www-form-urlencoded

上面的条件，只要有一点不满足就不是简单请求，例如常见的 application/json，或者携带自定义请求头，或请求方法为 PUT、DELETE。

当请求跨域且不是简单请求时，浏览器首先使用 OPTIONS 方法发起一个预检请求到服务器，以获知服务器是否允许该实际请求。"预检请求“的使用，可以避免跨域请求对服务器的用户数据产生未预期的影响。


### 添加OPTIONS请求处理中间件

cors.go

```go
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		//允许请求头 携带 x-token
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
```

router.go

```go

func Routers() *gin.Engine {
	Router := gin.Default()
	//配置跨域
	Router.Use(middlewares.Cors())
	//生成全局所有分组的顶层分组
	ApiGroup := Router.Group("/u/v1")
	router2.InitUserRouter(ApiGroup)

	return Router
}
```



## 9、 获取图片验证码

第三方库：https://mojotv.cn/go/refactor-base64-captcha

chaptcha.go

```go

var store = base64Captcha.DefaultMemStore


func GetCaptcha(ctx *gin.Context){
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		zap.S().Errorf("生成验证码错误,: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":"生成验证码错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath": b64s,
	})
}
```

base.go

```go
func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
	}
}
```

router.go

```go
ApiGroup := Router.Group("/u/v1")
router.InitUserRouter(ApiGroup)
router.InitBaseRouter(ApiGroup)
```

user.go

```go
//同一个包的变量可以直接访问
if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false) {
    c.JSON(http.StatusBadRequest, gin.H{
        "captcha": "验证码错误",
    })
    return
}
```



## 10、 阿里云发送短信

### 代码

```go
package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

func main(){
	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", "xxxx", "xxx")
	if err != nil {
		panic(err)
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = "xxx"                         //手机号
	request.QueryParams["SignName"] = "xxx"                               //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "xxx"       //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + "777777" + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	fmt.Print( client.DoAction(request, response))
	//  fmt.Print(response)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
	//json数据解析
}
```

### redis的安装

docker 安装

```shell
docker run -p 6379:6379 -d redis:latest redis-server  
docker container update --restart=always 容器名字 #设置容器自动启动，比如服务器重启后，容器会退出。设置后自动启动
```

redis的go驱动

```http
https://github.com/go-redis/redis
```



## 11、 redis保存验证码

```go
//将验证码保存起来 - redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)
```



## 12、 用户注册接口