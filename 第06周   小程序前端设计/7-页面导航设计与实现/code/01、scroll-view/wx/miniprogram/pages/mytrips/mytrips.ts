import { routing } from "../../utils/routing"

interface Trip {
    id: string
    start: string
    end: string
    duration: string
    fee: string
    distance: string
    status:string
}
const app = getApp<IAppOption>()
Page({

    data: {
        promotionItems: [
            {
                img: 'https://img.mukewang.com/5f7301d80001fdee18720764.jpg',
                promotionID: 1,
            },
            {
                img: 'https://img.mukewang.com/5f6805710001326c18720764.jpg',
                promotionID: 2,
            },
            {
                img: 'https://img.mukewang.com/5f6173b400013d4718720764.jpg',
                promotionID: 3,
            },
            {
                img: 'https://img.mukewang.com/5f7141ad0001b36418720764.jpg',
                promotionID: 4,
            },
        ],
        //licStatus: licStatusMap.get(rental.v1.IdentityStatus.UNSUBMITTED),
        avatarURL: '',
        tripsHeight: 0,
        navCount: 0,
        //mainItems: [] as MainItem[],
        mainScroll: '',
        // navItems: [] as NavItem[],
        navSel: '',
        navScroll: '',
        trips:[] as Trip[],
        scrollTop:0,
        scrollView:'',
    },

    /**
     * 生命周期函数--监听页面加载
     */
    async onLoad() {
        this.populateTrips()
        const userInfo = await app.globalData.userInfo
        this.setData({
            avatarURL: userInfo.avatarUrl,
        })
    },

    populateTrips(){
        const trips:Trip[] = []
        for(let i=0;i<100;i++){
            trips.push({
                id:(1001+i).toString(),
                start:'东方明珠',
                end:'迪士尼',
                distance:'12公里',
                fee:'128.00元',
                duration:'0时44分',
                status:'已完成'
            })
        }
        this.setData({
            trips:trips
        })
    },


    /**
     * 生命周期函数--监听页面初次渲染完成,渲染完成才能知道窗口高度
     */
    onReady() {
        wx.createSelectorQuery().select('#heading')
        .boundingClientRect(rect=>{
            this.setData({
                tripsHeight: wx.getSystemInfoSync().windowHeight-rect.height
            })
        }).exec()
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
    onSwiperChange(e: any) {
        //caused by program,eg:通过程序修改 current 值引起的轮播图变更事件。
        if (!e.detail.source) {
            return
        }

    },
    //data-promotion-id 数据传参方式，微信里所有的元素都支持这种方法。
    //currentTarget bindtap的元素
    onPromotionItemTap(e: any) {
        console.log(e)
        const promotionID = e.currentTarget.dataset.promotionId
        if (promotionID) {

        }
    },
    onRegisterTap() {
        wx.navigateTo({
            url: routing.register()
        })
    }
})