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

