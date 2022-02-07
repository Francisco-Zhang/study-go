## 1、什么是orm？如何正确看待orm？

### 什么是orm

ORM全称是：Object Relational Mapping（对象关系映射），其主要作用是在编程中，把面向对象的概念跟数据库中的表概念对应起来。举例来说就是，我定义一个对象，那就对应着一张表，这个对象的实例，就对应着表中的一条记录。

对于数据来说，最重要最常用的是表：表中有列，orm就是将一张表映射成一个类，表中的列映射成类中的一个类。java、python是这样的。但是对于go语言而言，表映射成struct。列如何映射？列可以映射成struct中的类型，int->int，但是有另外一个问题，就是数据库中的列具备很好的描述性，但是struct 有tag。 所以表和struct映射是非常合理的。

### 常用orm

https://github.com/go-gorm/gorm

https://github.com/facebook/ent

https://github.com/jmoiron/sqlx

https://gitea.com/xorm/xorm/src/branch/master/README_CN.md

https://github.com/didi/gendry/blob/master/translation/zhcn/README.md

gorm是关注读最高的。各框架差异不会很大。

### orm的优缺点

**优点**

1. 提高了开发效率
2. 屏蔽sql细节。可以自动对实体Entity对象与数据库只的Table进行字段与属性的映射；不用直接sql编码。
3. 屏蔽各种数据库之间的差异

**缺点**

1. orm会牺牲程序的执行效率，以及会固定思维模式。
2. 太过依赖orm导致sql理解不够
3. 对于固定的orm过重，导致切换到其他orm代价高

### 如何正确看到orm和sql之间的关系

1. sql为主，orm为辅
2. orm主要目的是为了增加代码可维护性和开发效率。



## 2、gorm连接数据库

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:admin123@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
```



## 3、快速体验auto migrate功能

可以自动建表

```go
import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Product struct {
	gorm.Model
	Code  sql.NullString
	Price uint
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:admin123@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	//设置全局的logger，这个logger在我们执行每个sql语句的时候会打印每一行sql
	//sql才是最重要的，本着这个原则我尽量的给大家看到每个api背后的sql语句是什么

	//定义一个表结构， 将表结构直接生成对应的表(会自动生成数据库表) - migrations
	// 迁移 schema
	_ = db.AutoMigrate(&Product{}) //此处应该有sql语句
}
```



## 4、gorm的Model的逻辑删除

### 代码

```go
import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Product struct {
	gorm.Model
	Code  sql.NullString
	Price uint
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:admin123@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	//设置全局的logger，这个logger在我们执行每个sql语句的时候会打印每一行sql
	//sql才是最重要的，本着这个原则我尽量的给大家看到每个api背后的sql语句是什么

	//定义一个表结构， 将表结构直接生成对应的表(会自动生成数据库表) - migrations
	// 迁移 schema
	_ = db.AutoMigrate(&Product{}) //此处应该有sql语句

	// 新增
	db.Create(&Product{Code: sql.NullString{"D42", true}, Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // 根据整形主键查找
	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	db.Model(&product).Updates(Product{Price: 200, Code: sql.NullString{"", true}}) // 仅更新非零值字段
	//如果我们去更新一个product 只设置了price：200
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product， 并没有执行delete语句，逻辑删除
	//db.Delete(&product, 1)
}
```



### 对应sql语句

会自动操作一些公共字段，删除是逻辑删除，会更新删除时间字段，查询的时候默认查询更新时间为null的数据。

```shell
[0.000ms] [rows:-] SELECT DATABASE()

[1.012ms] [rows:1] SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE 'gorm_test%' ORDER BY SCHEMA_NAME='gorm_test' DESC limit 1

[3.671ms] [rows:-] SELECT count(*) FROM information_schema.tables WHERE table_schema = 'gorm_test' AND table_name = 'products' AND table_type = 'BASE TABLE'

[0.506ms] [rows:-] SELECT DATABASE()

[0.613ms] [rows:1] SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE 'gorm_test%' ORDER BY SCHEMA_NAME='gorm_test' DESC limit 1

[0.885ms] [rows:-] SELECT column_name, is_nullable, data_type, character_maximum_length, numeric_precision, numeric_scale , datetime_precision FROM information_schema.columns WHERE 
table_schema = 'gorm_test' AND table_name = 'products' ORDER BY ORDINAL_POSITION

[0.000ms] [rows:-] SELECT DATABASE()

[0.552ms] [rows:1] SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE 'gorm_test%' ORDER BY SCHEMA_NAME='gorm_test' DESC limit 1

[1.205ms] [rows:-] SELECT count(*) FROM information_schema.statistics WHERE table_schema = 'gorm_test' AND table_name = 'products' AND index_name = 'idx_products_deleted_at'        

[15.866ms] [rows:1] INSERT INTO `products` (`created_at`,`updated_at`,`deleted_at`,`code`,`price`) VALUES ('2022-02-06 21:14:24.061','2022-02-06 21:14:24.061',NULL,'D42',100)       

[0.767ms] [rows:1] SELECT * FROM `products` WHERE `products`.`id` = 1 AND `products`.`deleted_at` IS NULL ORDER BY `products`.`id` LIMIT 1

[0.506ms] [rows:0] SELECT * FROM `products` WHERE code = 'D42' AND `products`.`deleted_at` IS NULL AND `products`.`id` = 1 ORDER BY `products`.`id` LIMIT 1

[6.370ms] [rows:1] UPDATE `products` SET `price`=200,`updated_at`='2022-02-06 21:14:24.08' WHERE `id` = 1 AND `products`.`deleted_at` IS NULL

[6.730ms] [rows:1] UPDATE `products` SET `updated_at`='2022-02-06 21:14:24.086',`code`='',`price`=200 WHERE `id` = 1 AND `products`.`deleted_at` IS NULL
```



## 5、通过NullString解决不能更新零值的问题

```go
type Product struct {
	gorm.Model
	Code  string
	Price uint
}
db.Model(&product).Updates(Product{Price: 200, Code: ''}) // 仅更新非零值字段,code是无法更新的。

//要解决这个问题，需要更改模型struct的数据类型
db.Model(&product).Updates(Product{Price: 200, Code: sql.NullString{"", true}}) 
//Price是0也是无法更新的，解决这类问题，需要更改数据类型为相应的 NullBool,NullInt32,NullFloat64等。
```



## 6、表结构定义细节

文档地址：https://learnku.com/docs/gorm/v2/models/9729

约定
GORM 倾向于约定，而不是配置。默认情况下，GORM 使用 ID 作为主键，使用结构体名的 **蛇形复数** 作为表名(自动加s)，字段名的 蛇形 作为列名，并使用 CreatedAt、UpdatedAt 字段追踪创建、更新时间

遵循 GORM 已有的约定，可以减少您的配置和代码量。如果约定不符合您的需求，GORM 允许您自定义配置它们



### 字段标签

声明 model 时，tag 是可选的，GORM 支持以下 tag： tag 名大小写不敏感，但建议使用 `camelCase` 风格

| 标签名         | **说明**                                                     |
| -------------- | ------------------------------------------------------------ |
| column         | 指定 db 列名                                                 |
| type           | 列数据类型，推荐使用兼容性好的通用类型，例如：所有数据库都支持 bool、int、uint、float、string、time、bytes 并且可以和其他标签一起使用，例如：not null、size, autoIncrement… 像 varbinary(8) 这样指定数据库数据类型也是支持的。在使用指定数据库数据类型时，它需要是完整的数据库数据类型，如：MEDIUMINT UNSIGNED not NULL AUTO_INSTREMENT |
| size           | 指定列大小，例如：`size:256`                                 |
| primaryKey     | 指定列为主键                                                 |
| unique         | 指定列为唯一                                                 |
| default        | 指定列的默认值                                               |
| precision      | 指定列的精度                                                 |
| scale          | 指定列大小                                                   |
| not null       | 指定列为 NOT NULL                                            |
| autoIncrement  | 指定列为自动增长                                             |
| embedded       | 嵌套字段                                                     |
| embeddedPrefix | 嵌入字段的列名前缀                                           |

### 代码举例

```go
package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	UserID uint   `gorm:"primarykey"`
	Name   string `gorm:"column:user_name;type:varchar(50);index:idx_user_name;unique;default:'bobby'"`
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:admin123@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&User{}) //此处应该有sql语句

	db.Create(&User{})

}
```

