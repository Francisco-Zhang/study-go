## 1、关系型数据库 vs 非关系型数据库

### sql

优点

- 成熟、开发人员熟悉
- 丰富的生态
- 一致性保证（最大优点）

缺点

- Object-Relational Mapping （为了描述一个对象，需要建立很多张表，关系表）
- 性能

用途

- 遗留的系统
- ToB的系统

### no-sql

种类：Redis、MongoDB、ElasticSearch、HBase.....

MongoDB 优点

- 保存的JSON文档为一条记录，一条记录包含数据量大，结构复杂
- 丰富的查询功能
- 性能比关系型数据库好（牺牲了关系型db的功能，事务支持不是特别高）



MongoDB 缺点

- 事务支持不是特别高
- 不支持Join



用途

- 快速开发
- ToC的系统
- Serverless 云开发的宠儿，背后都是MongoDB
  - Firebase
  - LeanCloud
  - 腾讯云开发



## 2、用docker来启动MongoDB

### 安装容器服务

```shell
docker run -d -p 27017:27017 --name car-mongo mongo:latest
```

### vscode安装插件

安装插件：MongoDB for VS Code

点击vscode左侧叶子图标——点击连接——通过 localhost:27017 建立连接

点击 *Create New Playground*  编写脚本,然后点击 右上角的 播放 按钮运行脚本。

### 设置默认数据库

这样就不需要每次都在 Playground 中，使用 use('coolcar'); 来指定数据库了。

现有连接右键复制连接字符串——删除现有连接——新建连接——使用字符串连接——原有字符串中加入默认库 coolcar

mongodb://localhost:27017/coolcar?readPreference=primary&ssl=false

## 3、MongoDB的CRUD操作

collection 就是通常意义上的表的概念

db.sales.drop(); 删除表

在配置中可以关掉 确认 弹窗。

### create

选中脚本，点击 黄色💡 图标，运行选择行。

```sql
db.acount.insert({
  open_id:"123",
  login_count:0,
});
```

插入多条

```sql
db.acount.insertMany([{
  _id: "user123",
  open_id:"123",
  login_count:0,
},{
  open_id:"456",
  login_count:0,
},
]);
```

### Retrieve

```sql
-- 查找所有
db.acount.find();

-- 根据条件查找,返回数组
db.acount.find({
    _id:ObjectId("628d8a753fce4d042ad61ed5")
});

db.acount.findOne({
    login_count:0
});


-- and 关系
db.acount.find({
     login_count:{$gt:3},
     open_id:"ssss",
})

-- or
db.acount.find({
    $or:[
        {
            login_count:{$gt:3},
            open_id: "1234"
        },{
            login_count:0,
        }
    ] 
})

```



### Update

```sql
-- 默认更新一个
db.acount.update({
    _id:ObjectId("628d8a753fce4d042ad61ed5")
},{
    $set: {
      login_count: 1
    }
})
```

配合find使用

```sql
acount=db.acount.find({
    login_count:1
});

db.acount.update({
    _id:ObjectId("628d8a753fce4d042ad61ed5")
},{
    $set: {
      login_count: acount.login_count+1
    }
})
```

单条语句保证原子性

```sql
db.acount.update({
    _id:ObjectId("628d8a753fce4d042ad61ed5")
},{
    $inc: {
      login_count: 1
    }
})
```

```sql
-- 更新多个字段，一条语句是原子性的，但是多条不保证。
db.acount.update({
    _id:ObjectId("628d8a753fce4d042ad61ed5")
},{
    $inc: {
      login_count: 1
    },
    $set: {
      open_id: "1234"
    }
})
```

```sql
db.acount.update({
    _id:ObjectId("628d8a753fce4d042ad61ed5")
},{
    $set: {
      profile:{
          name:"abc",
          age:28
      }
    }
})

db.acount.find({
    "profile.age":{$gt:3},
});
```



### Delete

```sql
db.acount.deleteOne({
    _id:"user123"
})
```



### 建索引

```sql
db.acount.createIndex({
    "profile.age":1,  --  1表示正序，-1表示倒序
});
```



## 4、用MongoDB Playground模拟用户登陆

```javascript
db.account.drop()

db.account.insertMany([
    {open_id:'123'},
    {open_id:'456'},
])

db.account.createIndex({
    open_id:1,
},{
    unique:true
});

function resolveOpenID_(open_id){
  return db.account.updateOne({
        open_id:open_id
    },{
        $set: {
          open_id:open_id
        }
    },{
        upsert:true
    });
}

// 相比updateOne，功能相同，但是返回值不同。find 会返回查找到的记录
// 通过建索引，实现唯一，不然无法保证
function resolveOpenID(open_id){
    return db.account.findAndModify({
        query:{
        open_id:open_id
        },
        update:{
            $set: {
                open_id:open_id
            }
        },
        upsert:true,
        new:true,  //返回更新后的数据
    });
}

resolveOpenID('789')

db.account.find()
```



## 5、通过go语言来操作MongoDB

server项目下面新建目录cmd--mongo,新建文件 main.go

```go
package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&ssl=false"))
	if err != nil {
		panic(err)
	}
	col := mc.Database("coolcar").Collection("account")
	findRows(c, col)
}

func findRows(c context.Context, col *mongo.Collection) {
	cur, err := col.Find(c, bson.M{})
	if err != nil {
		panic(err)
	}
	for cur.Next(c) {
		var row struct {
			ID     primitive.ObjectID `bson:"_id"`
			OpenID string             `bson:"open_id"`
		}
		err := cur.Decode(&row)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", row)
	}
}

func findRowsOne(c context.Context, col *mongo.Collection) {
	res := col.FindOne(c, bson.M{
		"open_id": "1231",
	})
	var row struct {
		ID     primitive.ObjectID `bson:"_id"`
		OpenID string             `bson:"open_id"`
	}
	err := res.Decode(&row)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", row)
}

func insertRows(c context.Context, col *mongo.Collection) {
	res, err := col.InsertMany(c, []interface{}{
		bson.M{
			"open_id": "1231",
		},
		bson.M{
			"open_id": "4561",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", res)
}
```

