
// index.js

import { routing } from "../../utils/routing"

const initialLat = 23.099994
const initialLng = 113.324520
const app = getApp<IAppOption>()


Page({
  isPageShowing: false,
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
    scale: 10,
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
      latitude: 29.756825521115363,
      longitude: 121.8722114786053,
      width: 20,
      height: 20
    }],
    showModal:true,
    showCancel:true
  },
  pathIndex: 0,
  translateMarker(ctx: any) {
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
  async onLoad() {
    const userInfo = await app.globalData.userInfo
    this.setData({
      avatarURL: userInfo.avatarUrl,
    })
  },
  onMyLocatioTap() {
    wx.getLocation({
      type: 'gcj02',
      success: res => {
        this.setData({
          location: {
            longitude: res.longitude,
            latitude: res.latitude
          }
        })
      },
      fail: () => {
        wx.showToast({
          icon: 'none',
          title: '请前往设置页授权'  //用户点击了取消授权后，可以在设置里重新修改
        })
      }
    })
  },
  moveCars() {
    const map = wx.createMapContext('map')
    const dest = {
      latitude: initialLat,
      longitude: initialLng,
    }
    const moveCar = () => {
      dest.latitude += 0.1
      dest.longitude += 0.1
      map.translateMarker({
        destination: {
          latitude: dest.latitude,
          longitude: dest.longitude,
        },
        autoRotate: false,
        markerId: 0,
        rotate: 0,
        duration: 5000,
        animationEnd: () => {
          if (this.isPageShowing) {
            moveCar()
          }
        }
      })
    }
    moveCar()
  },

  onShow() {
    this.isPageShowing = true
  },
  onHide() {
    this.isPageShowing = false
  },
  onScanTap() {
    wx.scanCode({
      success: async () => {
        await this.selectComponent('#licModal').showModal()
         //TODO get carid from scan
         const carID = 'car123'
         const redirectURL = routing.lock({
           car_id: carID
         })
         wx.navigateTo({
           //redirectURL 有特殊字符，需要转义
           url: routing.register({
             redirectURL: redirectURL
           })
         })
      },
      fail: console.error,
    })
  },
  onMyTripsTap() {
    wx.navigateTo({
      url: routing.mytrips()
    })
  },


})


const raceData = {
  "car_id": "222222",
  "path": [{
    "elapsed_ms": 0,
    "lat": 29.756825521115363,
    "lng": 121.87221147860531
  }, {
    "elapsed_ms": 100,
    "lat": 29.75684382898254,
    "lng": 121.87224327516209
  }, {
    "elapsed_ms": 200,
    "lat": 29.756862303343542,
    "lng": 121.87226556885331
  }, {
    "elapsed_ms": 300,
    "lat": 29.756880443491006,
    "lng": 121.87228852797627
  }, {
    "elapsed_ms": 400,
    "lat": 29.756898917036104,
    "lng": 121.87231148714312
  }, {
    "elapsed_ms": 500,
    "lat": 29.75691722369424,
    "lng": 121.87233461267527
  }, {
    "elapsed_ms": 600,
    "lat": 29.75693530156913,
    "lng": 121.8723579045862
  },
  {
    "elapsed_ms": 700,
    "lat": 29.75695400311557,
    "lng": 121.8723813628892
  },
  {
    "elapsed_ms": 800,
    "lat": 29.7569726427814,
    "lng": 121.8724042122407
  }, {
    "elapsed_ms": 900,
    "lat": 29.756991448943772,
    "lng": 121.8724244595153
  }, {

    "elapsed_ms": 1000,
    "lat": 29.75701025491105,
    "lng": 121.8724523705914
  }, {
    "elapsed_ms": 110,
    "lat": 29.75702939427861,
    "lng": 121.87247602821199
  }, {
    "elapsed_ms": 120,
    "lat": 29.757048533451606,
    "lng": 121.872499857452
  }, {
    "elapsed_ms": 130,
    "lat": 29.757067839332514,
    "lng": 121.8725239431159
  }, {
    "elapsed_ms": 1400,
    "lat": 29.75708714502009,
    "lng": 121.87254806725889
  }, {

    "elapsed_ms": 1500,
    "lat": 29.75710678389689,
    "lng": 121.87257235761386
  }, {
    "elapsed_ms ": 1600,
    "lat": 29.75712608982862,
    "lng": 121.87259631523894
  }, {

    "elapsed ns": 1700,
    "lat": 29.757205979105887,
    "lng": 121.87262060564686
  }, {
    "elapsed_ms ": 1800,
    "lat": 29.757165867951205,
    "lng": 121.87264472972564
  }, {
    "elapsed_ms": 1900,
    "lat": 29.757186006760595,
    "lng": 121.87266918656208
  }, {
    "elapsed_ms": 2000,
    "lat": 29.757186006760595,
    "lng": 121.87269347704337
  }, {
    "elapsed_ms ": 210,
    "lat": 29.757226284432072,
    "lng": 121.87271810029618
  }, {
    "elapsed_ms ": 220,
    "lat": 29.757246923367386,
    "lng": 121.87274255723293
  }, {
    "elapsed_ms": 2300,
    "lat": 29.757267562320276,
    "lng": 121.87276701419012
  },
  {
    "elapsed_ms ": 2400,
    "lat": 29.757288534463942,
    "lng": 121.87279163755689
  },
  {
    "elapsed_ms": 2500,
    "lat": 29.757309506834236,
    "lng": 121.87281609458115
  }, {
    "elapsed_ms": 2600,
    "lat": 29.75733081218686,
    "lng": 121.87284088437868
  }, {
    "elapsed_ms ": 270,
    "lat": 29.757352117557563,
    "lng": 121.87286567419723
  }, {
    "elapsed_ms": 2800,
    "lat": 29.757373589637396,
    "lng": 121.87289046404985
  }, {
    "elapsed_ms": 290,
    "lat": 29.757395394700794,
    "lng": 121.8729158667717
  }, {

    "elapsed_ms": 3000,
    "lat": 29.75741686608864,
    "lng": 121.87294054293616
  }, {
    "elapmed_ms ": 310,
    "lat": 29.7574383819297,
    "lng": 121.87296583198372
  }, {
    "elapmed_ms": 3200,
    "lat": 29.757468095865,
    "lng": 121.87299128741746
  }, {
    "elapmed_ms": 330,
    "lat": 29.75748261430908,
    "lng": 121.87301674286044
  },
  {

    "elapmed_ms": 3400,
    "lat": 29.75750475225028,
    "lng": 121.87304236471624
  },
  {
    "elapmed_ms": 300,
    "lat": 29.757527056643845,
    "lng": 121.8730681529723
  },
  {
    "elapmed_ms": 360,
    "lat": 29.7575493610823,
    "lng": 121.87309394125124
  },
  {
    "elapmed_ms": 370,
    "lat": 29.757571832024162,
    "lng": 121.87311989593103
  }, {
    "elapsed_ms ": 380,
    "lat": 29.75759430298594,
    "lng": 121.87314585063403
  },
  {
    "elapsed_ms": 3900,
    "lat": 29.757617107142654,
    "lng": 121.87317197175153
  },
  {
    "elapsed_ms": 4000,
    "lat": 29.75763974462855,
    "lng": 121.87319809287956
  }, {
    "elapsed_ms": 4100,
    "lat": 29.757662381721136,
    "lng": 121.873245467623
  },
  {
    "elapsed_ms": 4200,
    "lat": 29.75768535221648,
    "lng": 121.8732510069509
  },
  {
    "elapsed_ms": 4300,
    "lat": 29.75770832252608,
    "lng": 121.8732762101798
  }, {
    "elapsed_ms": 4400,
    "lat": 29.757731292856715,
    "lng": 121.8730424136533
  },
  ]
}