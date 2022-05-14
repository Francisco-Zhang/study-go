import { routing } from "../../utils/routing"

// pages/mytrips/mytrips.ts
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
    },

    /**
     * 生命周期函数--监听页面加载
     */
    onLoad() {

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
    onSwiperChange(e:any){
        //caused by program,eg:通过程序修改 current 值引起的轮播图变更事件。
        if(!e.detail.source){
            return
        }

    },
    //data-promotion-id 数据传参方式，微信里所有的元素都支持这种方法。
    //currentTarget bindtap的元素
    onPromotionItemTap(e:any){
        console.log(e)
        const promotionID=e.currentTarget.dataset.promotionId
        if(promotionID){

        }
    },
    onRegisterTap(){
        wx.navigateTo({
            url:routing.register()
        })
    }
})