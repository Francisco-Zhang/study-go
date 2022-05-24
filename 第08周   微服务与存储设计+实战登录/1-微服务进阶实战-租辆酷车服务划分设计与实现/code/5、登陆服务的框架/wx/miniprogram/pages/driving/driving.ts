import { formatDuration, formatFee } from "../../utils/format"
import { routing } from "../../utils/routing"

const centPerSec = 0.7

function durationStr(sec: number) {
    const dur = formatDuration(sec)
    return `${dur.hh}:${dur.mm}:${dur.ss}`
}

Page({

    timer: undefined as number | undefined,
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


    onLoad(opt: Record<'trip_id', string>) {
        const o: routing.DrivingOpts = opt
        console.log('current trip', o.trip_id)
        this.setupLocationUpdator()
        this.setupTimer()
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
        if (this.timer) {
            clearInterval(this.timer)
        }
    },

    setupTimer() {
        let elapsedSec = 0
        let cents = 0
        this.timer = setInterval(() => {
            elapsedSec++
            cents += centPerSec
            this.setData({
                elapsed: durationStr(elapsedSec),
                fee: formatFee(cents)
            })
        }, 1000)
    },
    onEndTripTap() {
        wx.redirectTo({
            url:routing.mytrips()
        })
    }
})