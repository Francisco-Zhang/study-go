## 1、 回调函数的缺点

### 前端异步运行机制

- 发起请求
- 事件处理函数结束（不能一直等待请求响应，页面会卡死）
- 请求结束
- 进入回调函数

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

