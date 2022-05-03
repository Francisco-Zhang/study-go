## 1、 控制元素的位置

### 大小

绝对值、相对值（rem、rpx）、百分比

小程序多用相对值 rpx。rpx:不管屏幕大小，宽度都是750rpx。只要高度设置750rpx，宽度100%，就是一个正方形。



### 设置大小的分类

box-sizing:	border-box ,高度设置多少，就是高多少

box-sizing:	content-box ,默认的，高度不包含 border

微信小程序的大小显示方式更像是 border-box



### 位置设置

top、left、 right、 bottom

这四个属性单独设置都是没有反应的，一定要配上另一个值 position: relative,fixed,absolute.

fixed,absolute 不占位置，会单独给设置了位置的元素分配一个层进行渲染。

一般是父元素设置relative，子元素设置absolute，这样子元素可以根据滚动条滚动,并且以父元素为基准，新分配一个图层进行显示。

如果父元素不设置relative，则子元素的absolute位置是相对整个page。

如果父元素不设置relative，则子元素的relative位置是相对元素原来的位置偏移，并且占用原来图层的元素位置。

父元素设置了relative，但是没有设置具体的top、left等，相当于偏移原来的位置为0，对父元素没有任何影响。

```html
<div>
    <div id="title" class="blue">租辆车</div>
    <div>css 教学</div>
</div>

<div class="trip">
    <div class="item blue">
        <div>中关村到天安门</div>
        <div>价格：120元</div>
    </div>
    <div class="item red">
        <div>陆家嘴到迪士尼</div>
        <div>价格：50元</div>
        <div class="star"></div>
    </div>
    <div class="item blue">
        <div>天河体育中心到广州塔</div>
        <div>价格：80元</div>
    </div>
</div>


<style>
    div {
        font-size: x-large;
    }
    
    #title {
        font-size: 60px;
    }
    
    .red {
        color: red;
    }
    
    .blue {
        color: blue;
    }
    
    .item.blue {
        color: cornflowerblue;
    }
    /* 为了方便显示大小，显示边框 */
    
    .item {
        border: 2px solid;
    }
    
    .item.red {
        position: relative;
    }
    
    .star {
        width: 20px;
        height: 20px;
        background-color: red;
        position: absolute;
        right: 0;
        top: 0;
    }
</style>
```

## 2、 文本样式

折行：white-space,work-break(设置英文字符用),text-overflow



## 3、flex布局

```css
display:flex
flex-direction:column
align-items:center
```

## 4、 在小程序中使用css

### text 和 view区别

在小程序中 text 和 view 内部都可以放置文本，区别是 text 有个特殊的属性user-select，可以让文本能够被选择

```xml
<text user-select wx:if="{{showPath}}">pages/learncss/learncss.wxml</text>
<view wx:for="{{values}}" wx:for-item="val" wx:key="*this">
    employee is {{val.name}}
</view>
```

### 公用样式container

```xml
<!-- app.wxss中有container样式，app.wxss中定义的样式是公共的，所有页面都生效，自己页面的样式只有自己页面能够使用 -->
<view class="container">
<text user-select wx:if="{{showPath}}">pages/learncss/learncss.wxml</text>
<block wx:else>
    <text>learn css</text>
</block>


<view wx:for="{{values}}" wx:for-item="val" wx:key="*this">
    employee is {{val.name}}
</view>
</view>
```

### rpx

在iPhone 6，7，8上面 正好屏幕宽度 375px,宽度 750rpx。是一个2倍的关系。

```xml
<!-- pages/learncss/learncss.wxml -->
<view class="container">
    <text user-select wx:if="{{showPath}}">pages/learncss/learncss.wxml</text>
    <block wx:else>
        <text>learn css</text>
    </block>
    <!-- 类似于div,wx:key具有优化作用，不写会有warning, `*this`表示key就是value本身 -->
    <view class="list">
        <view class="item {{val.id%2?'blue':'red'}}" wx:for="{{values}}" wx:for-item="val" wx:key="*this">
            employee is {{val.name}}
        </view>
    </view>
</view>

/* pages/learncss/learncss.wxss */

page {
    height: 100%;
}

.item {
    font-size: xx-large;
}

.list {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;
    height: 750rpx; 			/* 设置高度=宽度 */
    align-self: stretch; 	/* 此处 align-self: stretch 等价于 width: 100%; */
}

.blue {
    color: blue;
}

.red {
    color: red;
}
```

