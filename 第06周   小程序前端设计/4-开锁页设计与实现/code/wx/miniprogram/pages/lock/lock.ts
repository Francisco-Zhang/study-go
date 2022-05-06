import { getUserProfile } from "../../utils/util"
// 获取应用实例
const app = getApp<IAppOption>()
const shareLocationKey = "share_location"
Page({

    /**
     * 页面的初始数据
     */
    data: {
        shareLocation:false,
        avatarURL: ''
    },

    /**
     * 生命周期函数--监听页面加载
     */
    async onLoad() {
        const userInfo = await app.globalData.userInfo
        this.setData({
            avatarURL: userInfo.avatarUrl,
            shareLocation:wx.getStorageSync(shareLocationKey) || false
        })
    },

    /**
     * 生命周期函数--监听页面初次渲染完成
     */
    onReady() {

    },

    /**
     * 生命周期函数--监听页面显示
     */
    onShow() {

    },

    /**
     * 生命周期函数--监听页面隐藏
     */
    onHide() {

    },

    /**
     * 生命周期函数--监听页面卸载
     */
    onUnload() {

    },

    /**
     * 页面相关事件处理函数--监听用户下拉动作
     */
    onPullDownRefresh() {

    },

    /**
     * 页面上拉触底事件的处理函数
     */
    onReachBottom() {

    },

    /**
     * 用户点击右上角分享
     */
    onShareAppMessage() {

    },
    // 使用语法糖，进一步简化写法
    async getUserProfile() {
        try {
            const userProfile = await getUserProfile()
            console.log(userProfile)
            app.resolveUserInfo(userProfile.userInfo)
        } catch (error) {
            console.log("用户拒绝授权")
        }
    },
    onShareLocation(e: any) {
        const shareLocation:boolean = e.detail.value
        wx.setStorageSync(shareLocationKey,shareLocation)
    }

})