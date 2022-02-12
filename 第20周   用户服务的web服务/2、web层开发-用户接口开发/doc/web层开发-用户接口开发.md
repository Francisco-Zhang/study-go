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

