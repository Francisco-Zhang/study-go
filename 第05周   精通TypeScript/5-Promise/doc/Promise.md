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
},
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

