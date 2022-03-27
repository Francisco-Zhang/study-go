## 1、 需求分析-数据库实体分析



## ![1](img/1.PNG)





## 2、 需求分析-商品微服务接口分析



## 3、 商品分类表结构设计应该注意什么？

```
类型， 这个字段是否能为null， 这个字段应该设置为可以为null还是设置为空， 0
实际开发过程中 尽量设置为不为null，而是设置默认值
https://zhuanlan.zhihu.com/p/73997266
这些类型我们使用int32还是int，proto中没有int类型，为了减少类型转换，使用int32
```

## 4、 品牌、轮播图表结构设计

**CategoryID 和 BrandsID构成 唯一联合索引，防止插入重复数据。**

采用自己定义一张表，而不是gorm帮我们生成的方式。

```go
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"` //如果不写int32，容易被对应成bigint类型。
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands   Brands
}

func (GoodsCategoryBrand) TableName() string {	//自定义表名
	return "goodscategorybrand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url string `gorm:"type:varchar(200);not null"`
	Index int32 `gorm:"type:int;default:1;not null"`
}
```

## 5、 商品表结构设计

自定义数据类型

```go
type GormList []string

func (g GormList) Value() (driver.Value, error){
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type Goods struct {
	BaseModel
	.
  .
  .
  //因为很少会根据图片的id查询商品，所以放弃使用关联表的做法，采用自定义类型。
	Images GormList `gorm:"type:varchar(1000);not null"`  
	DescImages GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string `gorm:"type:varchar(200);not null"`
}
```

