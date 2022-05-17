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
interface MainItem {
    id: string
    navId: string
    navScrollId: string
    data: Trip
}


interface NavItem {
    id: string
    mainId: string
    label: string
}

interface MainItemQueryResult {
    id: string
    top: number
    dataset: {
        navId: string
        navScrollId: string
    }
}

const app = getApp<IAppOption>()
Page({
    scrollStates: {
        mainItems: [] as MainItemQueryResult[],
    },
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
        mainItems: [] as MainItem[],
        mainScroll: '',
        navItems: [] as NavItem[],
        navSel: '',
        navScroll: '',
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
        const mainItems: MainItem[] = []
        const navItems: NavItem[] = []
        let navSel=''
        let prevNav=''
        for(let i=0;i<100;i++){
            const mainId = 'main-' + i
            const navId = 'nav-' + i
            const tripId = (1001+i).toString()
            if(!prevNav){
                prevNav = navId
            }
            mainItems.push({
                id: mainId,
                navId: navId,
                navScrollId: prevNav,  //真正选中的是上一个元素，达到留白效果
                data: {
                    id:tripId,
                    start:'东方明珠',
                    end:'迪士尼',
                    distance:'12公里',
                    fee:'128.00元',
                    duration:'0时44分',
                    status:'已完成'
                },
            })

            navItems.push({
                id: navId,
                mainId: mainId,
                label: tripId,
            })
            if(i===0){
                navSel=navId
            }
            prevNav = navId
        }
        this.setData({
            mainItems,
            navItems,
            navSel,
        },()=>{
            this.prepareScrollStates()
        })
    },

    prepareScrollStates(){
        wx.createSelectorQuery().selectAll('.main-item')
        .fields({   //返回所有的 main-item 元素的 id,dataset,rect
            id: true,
            dataset: true,
            rect: true,
        }).exec(res => {
            this.scrollStates.mainItems = res[0]  //将这100个元素的位置信息保存到mainItems
        })
    },

    /**
     * 生命周期函数--监听页面初次渲染完成,渲染完成才能知道窗口高度
     */
    onReady() {
        wx.createSelectorQuery().select('#heading')
        .boundingClientRect(rect=>{
            const height =  wx.getSystemInfoSync().windowHeight-rect.height
            this.setData({
                tripsHeight:height,
                navCount:height/50
                
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
    },
    onNavItemTap(e:any){
        const mainId:string = e.currentTarget?.dataset?.mainId
        const navId:string = e.currentTarget?.id
        if(mainId && navId){
            this.setData({
                mainScroll:mainId,
                navSel:navId
            })
        }
    },
    onMainScroll(e: any) {
        //scrollTop表示滑动的位置，offsetTop 表示当前item 顶部距离滑动窗口的位置。
        const top: number = e.currentTarget?.offsetTop + e.detail?.scrollTop
        if (top === undefined) {
            return
        }

        const selItem = this.scrollStates.mainItems.find(
            v => v.top >= top)
        if (!selItem) {
            return
        }
        this.setData({
            navSel: selItem.dataset.navId,
            navScroll: selItem.dataset.navScrollId,
        })
    },

    onMianItemTap(e: any) {
        if (!e.currentTarget.dataset.tripInProgress) {
            return
        }
        const tripId = e.currentTarget.dataset.tripId
        if (tripId) {
            wx.redirectTo({
                url: routing.drving({
                    trip_id: tripId,
                }),
            })
        }
    }
})