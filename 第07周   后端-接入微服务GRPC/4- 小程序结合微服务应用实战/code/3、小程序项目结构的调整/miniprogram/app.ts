import { IAppOption } from "./appoption"

let resolveUserInfo: (value: WechatMiniprogram.UserInfo) => void
let rejectUserInfo: (reason?: any) => void

// app.ts
App<IAppOption>({
  globalData: {
    userInfo:new Promise((resolve,reject)=>{
      resolveUserInfo=resolve
      rejectUserInfo=reject
    })
  },
  resolveUserInfo(userInfo: WechatMiniprogram.UserInfo){
    resolveUserInfo(userInfo)
  },
  rejectUserInfo(reason?: any){
    rejectUserInfo(reason)
  },
  onLaunch() {
    // 展示本地存储能力
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)

    // wx.request({
    //   url:'http://localhost:8080/trip/trip123',
    //   method:'GET',
    //   success:(res)=>{
    //     const getTripResp = res.data
    //     console.log(getTripResp)
    //   },
    //   fail:console.error
    // })

    // 登录
    wx.login({
      success: res => {
        console.log(res.code)
        // 发送 res.code 到后台换取 openId, sessionKey, unionId
      },
    })
  },
})