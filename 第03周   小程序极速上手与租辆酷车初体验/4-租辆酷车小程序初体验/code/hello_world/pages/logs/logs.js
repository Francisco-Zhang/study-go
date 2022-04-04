// logs.js
const util = require('../../utils/util.js')

Page({
    data: {
        logs: []
    },



    onShow() {
        console.log("lifecycle: logs onShow");
    },
    onHide() {
        console.log("lifecycle: logs onHide");
    },
    onReady() {
        console.log("lifecycle: logs onReady");
    },
    onUnload() {
        console.log("lifecycle: logs onUnload");
    },

    onLoad(opt) {
        console.log("lifecycle: logs onLoad");
        this.setData({
            logs: (wx.getStorageSync('logs') || []).map(log => {
                return {
                    date: util.formatTime(new Date(log)),
                    timeStamp: log
                }
            }),
            logColor: opt.color,
        })
    },
    onLogTap() {
        // wx.navigateTo({
        //     // url: '../test2/test2'
        //     url: '/pages/test2/test2' //这里改为绝对路径比较好
        // })

        //销毁当前页面(unLoad)，然后重新打开url页面
        wx.redirectTo({
            url: '/pages/test2/test2' //这里改为绝对路径比较好
        })
    }
})