import camelcaseKeys from "camelcase-keys"
import { IAppOption } from "./appoption"
import { auth } from "./service/proto_gen/auth/auth_pb"
import { coolcar } from "./service/proto_gen/trip_pb"

let resolveUserInfo: (value: WechatMiniprogram.UserInfo) => void
let rejectUserInfo: (reason?: any) => void

// app.ts
App<IAppOption>({
  globalData: {
    userInfo: new Promise((resolve, reject) => {
      resolveUserInfo = resolve
      rejectUserInfo = reject
    })
  },
  resolveUserInfo(userInfo: WechatMiniprogram.UserInfo) {
    resolveUserInfo(userInfo)
  },
  rejectUserInfo(reason?: any) {
    rejectUserInfo(reason)
  },
  onLaunch() {
    // 展示本地存储能力
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)



    // 登录
    wx.login({
      success: res => {
        console.log(res.code)
        // 发送 res.code 到后台换取 openId, sessionKey, unionId
        wx.request({
          url: 'http://localhost:8080/v1/auth/login',
          method: 'POST',
          data: {
            code: res.code
          } as auth.v1.ILoginRequest,
          success: res=>{
            const loginResp :auth.v1.ILoginResponse = auth.v1.LoginResponse.fromObject(
              camelcaseKeys(res.data as object)
            )
            console.log(loginResp)
          },
          fail: console.error
        })
      },
    })
  },
})