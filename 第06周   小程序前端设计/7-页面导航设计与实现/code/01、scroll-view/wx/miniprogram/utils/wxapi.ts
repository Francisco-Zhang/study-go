export const formatTime = (date: Date) => {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours()
  const minute = date.getMinutes()
  const second = date.getSeconds()

  return (
    [year, month, day].map(formatNumber).join('/') +
    ' ' +
    [hour, minute, second].map(formatNumber).join(':')
  )
}

const formatNumber = (n: number) => {
  const s = n.toString()
  return s[1] ? s : '0' + s
}

//为了减少回调，采用Promise改造
export function getUserProfile(): Promise<WechatMiniprogram.GetUserProfileSuccessCallbackResult>{
  return new Promise((resolve,reject)=>{
      wx.getUserProfile({
      desc: '需要你的个人信息来为你提供服务', // 声明获取用户个人信息后的用途，后续会展示在弹窗中，请谨慎填写
      success: resolve,   //等同于 res => resolve(res),
      fail:reject,
    })
  })
}