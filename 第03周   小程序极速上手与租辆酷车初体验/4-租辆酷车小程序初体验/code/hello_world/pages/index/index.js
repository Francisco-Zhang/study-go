// index.js
// 获取应用实例

import { raceData } from "./raceData"

// 默认值
const setting = {
    skew: 0,
    rotate: 0,
    showLocation: false,
    showScale: false,
    subKey: '',
    layerStyle: 1,
    enableZoom: true,
    enableScroll: true,
    enableRotate: false,
    showCompass: false,
    enable3D: false,
    enableOverlooking: false,
    enableSatellite: false,
    enableTraffic: false,
}
const initialLat = 29.761267625855936
const initialLng = 121.87264654736123



Page({
    data: {
        avatarURL: '',
        setting: {
            skew: 0,
            rotate: 0,
            showLocation: true,
            showScale: true,
            subKey: '',
            layerStyle: 1,
            enableZoom: true,
            enableScroll: true,
            enableRotate: false,
            showCompass: false,
            enable3D: false,
            enableOverlooking: false,
            enableSatellite: false,
            enableTraffic: false,
        },
        location: {
            latitude: initialLat,
            longitude: initialLng,
        },
        scale: 16,
        markers: [{
            iconPath: "/resources/car.png",
            id: 0,
            latitude: 23.099994,
            longitude: 113.324520,
            width: 50,
            height: 50
        }, {
            iconPath: "/resources/car.png",
            id: 1,
            latitude: 23.099994,
            longitude: 113.324520,
            width: 50,
            height: 50
        }, {
            iconPath: "/resources/car.png",
            id: 2,
            latitude: 29.756825521115363,
            longitude: 121.8722114786053,
            width: 20,
            height: 20
        }]
    },
    pathIndex: 0,
    translateMarker(ctx) {
        console.log("pathIndex: " + this.pathIndex)
        this.pathIndex++
            if (this.pathIndex >= raceData.path.length) {
                return
            }
        ctx.translateMarker({
            markerId: 2,
            destination: {
                latitude: raceData.path[this.pathIndex].lat,
                longitude: raceData.path[this.pathIndex].lng,
            },
            duration: 200,
            success: () => this.translateMarker(ctx)
        })
    },
    onReady() {
        console.log("onready")
            //第一个参数上map的id,使用wx的方法，不必渲染整个map
        const ctx = wx.createMapContext('map', this)
        this.translateMarker(ctx)
    },
    onLoad() {
        console.log(raceData)
    }

})