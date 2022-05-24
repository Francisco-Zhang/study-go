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

