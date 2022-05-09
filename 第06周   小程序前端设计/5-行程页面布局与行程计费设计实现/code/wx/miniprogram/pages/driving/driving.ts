
Page({

    data: {
        location: {
            latitude: 30,
            longitude: 120,
        },
        scale: 12,
        elapsed: '00:00:00',
        fee: '0.00',
        markers: [
            {
                iconPath: "/resources/car.png",
                id: 0,
                latitude: 30,
                longitude: 120,
                width: 20,
                height: 20,
            },
        ],
    },
    setupLocationUpdator() {
        //微信可以主动监测位置变动，而不是我们不断的查询。
        wx.startLocationUpdate({
            fail: console.error
        })
        wx.onLocationChange(loc => {
            console.log(loc)
            this.setData({
                location: {
                    latitude: loc.latitude,
                    longitude: loc.latitude
                }
            })
        })
    },

    /**
     * 生命周期函数--监听页面加载
     */
    onLoad() {
        this.setupLocationUpdator()
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
        wx.stopLocationUpdate()
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

    }
})