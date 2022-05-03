## 1、 WXML简介

### wx:if 标签控制隐藏

```xml
<text wx:if="{{showPath}}">pages/learncss/learncss.wxml</text>
```

```typescript
Page({
    data: {
        showPath : true,
    },
})

// true 显示，false不显示，整个元素就没了，输出给浏览器的标签为 <page></page>
// true 输出标签：<page> <text>pages/learncss/learncss.wxml</text></page>
```

### else

有两种方式,只修改wxml文件不需要重新编译

```xml
<text wx:if="{{showPath}}">pages/learncss/learncss.wxml</text>
<text wx:if="{{!showPath}}">learn css</text>
```

```xml
<text wx:if="{{showPath}}">pages/learncss/learncss.wxml</text>
<block wx:else>
    <text>learn css</text>
</block>
```

### for

```xml
<!-- 类似于div,wx:key具有优化作用，不写会有warning, `*this`表示key就是value本身  -->
<view wx:for="{{values}}" wx:for-item="val" wx:key="*this">
    value is {{val}}
</view>
```

### wx:key的优化

```typescript
data: {
        showPath : true,
        values:[
            {
                name:'john',
                id:1,
            },
            {
                name:'mary',
                id:2,
            },
            {
                name:'tom',
                id:3,
            },
        ]
    },
```

```xml
<view wx:for="{{values}}" wx:for-item="val" wx:key="id">
    value is {{val}}
</view>
```

优化的地方在于，当我们添加 id =4 的员工时，wx 可以根据id,判断之前三个已经渲染过了，不需要再进行重新渲染

## 2、CSS选择器

样式直接设置在标签上，量非常大，css 选择器能解决这个问题。

在chrome开发者工具，command+f 打开 css选择器 调试工具。在查找框内输入：.item.blue，就能找到该类选择器选中的元素。

div#title：表示class含有title的div。

.item div: 表示class含有item的子元素div。空格代表父子关系，没有空格代表同时选择。

.trips>div : 表示trips 下面的一级子元素 div。

```html
<div>
    <div id="title">租辆车</div>
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
    </div>
    <div class="item blue">
        <div>天河体育中心到广州塔</div>
        <div>价格：80元</div>
    </div>
</div>

<style>
    div {
        font-size: xx-large;
    }
    #title {
        font-size: 60px;
    }
</style>
```

## 3、 CSS相关问题的提问方法

属性设置了并没有生效，有可能是组合错了，并不少属性写错了。

提问方式：

1. 推荐网站：codepen.io。可以把链接发给其他人，让别人帮忙查找问题。
2. 微信开发者工具 -- 项目 -- 新建代码片段 -- 分享 -- 生成分享链接

