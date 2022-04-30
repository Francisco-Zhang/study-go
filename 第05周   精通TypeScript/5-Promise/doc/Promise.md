## 1、 回调函数的缺点

### 前端异步运行机制

- 发起请求
- 事件处理函数结束（不能一直等待请求响应，页面会卡死）
- 请求结束
- 进入回调函数

### 回调地狱

实际场景中，一般需要很多次接口调用，过多的回调函数造成回调地狱。

```typescript
//如果不是回调函数
const setting = wx.getSetting()
if (setting.authSetting['scope.userInfo']){
  const userInfo = wx.getUserInfo()
  this.globalData.userInfo = userInfo
}

// 查看是否授权,获取用户信息，会写很多层
// 造成回调地狱，这就是 Promise 要解决的问题
wx.getSetting({
  success (res){
    if (res.authSetting['scope.userInfo']) {
      // 已经授权，可以直接调用 getUserInfo 获取头像昵称
      wx.getUserInfo({
        success: res => {
          this.globalData.userInfo = userInfo
        }
      })
    }
  }
})
```

## 2、 Promise的创建和使用

### 网络请求模拟

因为发送请求后不能阻塞UI响应，所以需要通过回调函数处理服务器返回的结果。下面是模拟网络请求的处理过程

```typescript
function add(a:number,b:number,callback:(res:number) => void): void{
  setTimeout(() =>{
    callback(a+b)
  },2000)
   
}

//有可能因为add过快执行结束，导致没有输出结果。可以通过chrome本身的控制台查看到结果
//实际运行中，在 playgroud 也是可以看到输出结果的
add(2,3,res => {
  console.log('2+3',res)
})
```

### 回调套回调

再次相加，只能回调套回调

```typescript
function add(a:number,b:number,callback:(res:number) => void): void{
  setTimeout(() =>{
    callback(a+b)
  },2000)
   
}

add(2,3,res => {
  console.log('2+3',res)
  add(res,4,res2 => {
   	 console.log('2+3+4',res2)
  })
})
```

### 使用Promise改造

```typescript
function add(a:number,b:number): Promise<number>{
  return new Promise((resolve,reject)=>{
    setTimeout(() =>{
    resolve(a+b)
    },2000)
  })
}

add(2,3).then(res=>{
    console.log('2+3',res)
  //如果不返回 promise 对象，则 resolve 函数会执行所有的then，参数为 undefined,否则交由下一个promise处理所有的then。
  //一定要写 return，then 才会链式执行
    return add(res,4)  
  }).then(res=>{
    console.log('2+3+4',res)
    return add(res,6)
  }).then(res=>{
    console.log('2+3+4+6',res)
  })

//不打log的写法
add(2,3)
.then(res=>add(res,4))
.then(res=>add(res,6))
.then(res=>{
  	console.log('final result is',res)
})
```

### 需要注意的用法

```typescript
let p = add(2,3)
p.then(res=>{
    console.log('2+3',res)
  }).then(res=>{
  //需要注意的是，此处会打印 undefined 因为此处接收的 res 为第一个then的返回结果，或第一个 then 返回的Promis的resolve结果。
    console.log('2222',res)    
  })

//这种写法适用于微信多个page切换获取用户信息。不必每次切换页面都重新发送请求。可以共用一个global Promise 对象。
p.then(res=>{
    console.log('3333',res)
})
```



### 异常处理

```typescript
function add(a:number,b:number): Promise<number>{
  return new Promise((resolve,reject)=>{
    if(b % 17 ===0 ){
        reject(`bad number ${b}`)
    }
    setTimeout(() =>{
    resolve(a+b)
    },2000)
  })
}

add(2,17).then(res=>{
  console.log('2+17',res)
  return add(res,4) 
}).then(res=>{
  console.log('2+3+4',res)
  return add(res,6)
}).then(res=>{
  console.log('2+3+4+6',res)
}).catch(err =>{  //任意一步出错，都会到这里，但是只能执行到一次
  console.log('caught error',err)
})
```

### (2+3)*4+5

```typescript
function add(a:number,b:number): Promise<number>{
  return new Promise((resolve,reject)=>{
    if(b % 17 ===0 ){
        reject(`bad number ${b}`)
    }
    setTimeout(() =>{
    resolve(a+b)
    },2000)
  })
}

function mul(a:number,b:number): Promise<number>{
  return new Promise((resolve,reject)=>{
    setTimeout(() =>{
    resolve(a*b)
    },3000)
  })
}

//(2+3)*4+5
add(2,3).then(res=>{
    console.log('2+3',res)
    return mul(res,4)
  }).then(res=>{
    console.log('(2+3)*4',res)
    return add(res,5)
  }).then(res=>{
    console.log('(2+3)*4+5',res)  //最后一步，没有then要执行了，就不用return了。
  }).catch(err =>{
    console.log('caught error',err)
  })
```



## 3、 同时等待多个Promise

```typescript
//(2+3)*(4+5)
//等两个Promise都返回结果再执行下一步操作，实际场景中经常多个接口并行
// then(res=>{}) 可以直接替换成 then(([a,b])=>{})
Promise.all([add(2,3),add(4,5)]).then(res=>{
  const [a,b] = res
  console.log('result ',a,b)
  return mul(a,b)
}).then(res=>{
  console.log('(2+3)*(4+5)',res)
})

//两个中任意一个返回，则执行then
Promise.race([add(2,3),add(4,5)]).then(res=>{
  console.log(res)
})
```

## 4、将小程序API改写成Promise

### 改造前

```typescript
wx.getSetting({
  success (res){
    if (res.authSetting['scope.userInfo']) {
      // 已经授权，可以直接调用 getUserInfo 获取头像昵称
      wx.getUserInfo({
        success: res => {
          this.globalData.userInfo = userInfo
        }
      })
    }
  }
})
```



### 改造后

```typescript
//util.ts 文件
export function getSetting(): Promise<WechatMiniprogram.GetSettingSuccessCallbackResult> {
  return new Promise((resolve, reject) => {
    wx.getSetting({
      success: resolve,
      fail: reject,
    })
  })
}

export function getUserInfo(): Promise<WechatMiniprogram.GetUserInfoSuccessCallbackResult> {
  return new Promise((resolve, reject) => {
    wx.getUserInfo({
      success: resolve,
      fail: reject,
    })
  })
}

//app.ts 文件
getSetting().then(res => {
  if (res.authSetting['scope.userInfo']) {
      return getUserInfo()
  }
  return undefined //返回对象不是Promise类型，会直接传给下一个then
}).then(res => {
  if(!res){
    return
  }
  this.globalData.userInfo = res.userInfo  //上面已经有if判断，可以不用加问号了， res?.userInfo
  if (this.userInfoReadyCallback) {
    this.userInfoReadyCallback(res)
  }
})    
```

## 5、获取用户头像

### 跳转页面后如何通知用户信息已经获取

- 回调函数
- EventEmitter
- Promise

微信的做法是设置回调函数

```typescript
//app.ts
if (this.userInfoReadyCallback) {
    this.userInfoReadyCallback(res)
}

//index page 
onLoad() {
  if(app.globalData.userInfo){
   this.setData({
      userInfo: app.globalData.userInfo,
      hasUserInfo: true
    })
  }else{
    //设置app.ts的回调函数，这种方式特别不稳定，if else 很难选择
    //这种方式没法解决多页面问题，如果是多个页面，需要多个回调函数，app.ts里需要维护一个回调函数数组，
    //如果有的page已经unload调了，那回调函数还要不要调用呢
    app.userInfoReadyCallback = res => {  
      this.setData({
      userInfo: app.globalData.userInfo,
      hasUserInfo: true
    })
    }
  } 
}
```

EventEmitter,是js一个标准，可以管理回调函数数组。在js中很好用，但是在ts中很难定义事件，类型很难定义。

所以最为推荐 Promise 的形式。

改造 app.globalData.userInfo,原来的代码：

```typescript
interface IAppOption {
   globalData: {
     userInfo?: WechatMiniprogram.UserInfo,
   }
   userInfoReadyCallback?: WechatMiniprogram.GetUserInfoSuccessCallback,
 }
```

改造后

```typescript
interface IAppOption {
  globalData: {
    userInfo: Promise<WechatMiniprogram.UserInfo> ,
  }
}

//index.ts，如果app.ts很快拿到userInfo，则直接执行then
//否则，可以等待app.ts收到响应后执行then
//可以很方便的在多个页面使用
onLoad() {
  app.globalData.userInfo.then(res => {
     this.setData({
      userInfo: res,
      hasUserInfo: true
    })
  })
}
//语法上的小技巧
app.globalData.userInfo.then(userInfo => {
     this.setData({
      userInfo,
      hasUserInfo: true
    })
 })

//app.ts 设置全局的 userInfo
App<IAppOption>({
  globalData: {
      userInfo:new Promise((resolve,reject)=>{
        getSetting().then(res => {
        if (res.authSetting['scope.userInfo']) {
            return getUserInfo()
        }
        return undefined //返回对象不是Promise类型，会直接传给下一个then
      	}).then(res => {
        if(!res){
          return
        }
   			//通知
        resolve(res)
      	}).catch(reject) 
      }),
  },
)
```

### 控制台查看 app

在控制台通过 getApp() 函数可以获取全局的 app

```typescript
//调试的时候获取全局的 userInfo 查看 Promise 的状态
getApp().globalData.userInfo
//可以看到 状态是 resolved 或者 pending 等。
```

