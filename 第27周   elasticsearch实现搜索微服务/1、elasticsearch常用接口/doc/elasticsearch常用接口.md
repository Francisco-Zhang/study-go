## 映射（mapping）1、 我们为什么要用elasticsearch进行搜索

### mysql搜索面临的问题

- 性能低下
- 没有相关性排名
- 无法全文搜索
- 搜索不准确，没有分词



### 什么是全文搜索

我们生活中的数据总体分为两种：结构化数据和非结构化数据。

- 结构化数据：指具有固定格式或有限长度的数据，如数据库，元数据等。
- 非结构化数据：指不定长或无固定格式的数据，如邮件，word文档等。
- 半结构化数据，如XML，HTML，JSON 等，当根据需要可按结构化数据来处理，也可抽取出纯文本按非结构化数据来处理。

**非结构化数据又一种叫法叫全文数据。**

按照数据的分类，搜索也分为两种：

- 对结构化数据的搜索：如对数据库的搜索，用 SQL语句。再如对元数据的搜索，如利用windows搜索对文件名，类型，修改时间进行搜索等。 

- 对非结构化数据的搜索：如利用 windows 的搜索也可以搜索文件内容，Linux 下的 grep命令，再如用Google 和百度可以搜索大量内容数据

对非结构化数据也即对全文数据的搜索主要有两种方法：

一种是**顺序扫描法(Serial Scanning)：**所谓顺序扫描，比如要找内容包含某一个字符串的文件，就是一个文档一个文档的看，对于每一个文档，从头看到尾，如果此文档包含此字符串，则此文档为我们要找的文件，接着看下一个文件，直到扫描完所有的文件。如利用windows的搜索也可以搜索文件内容，只是相当的慢。如果你有一个80G硬盘，如果想在上面找到一个内容包含某字符串的文件，不花他几个小时，怕是做不到。Linux下的grep命令也是这一种方式。大家可能觉得这种方法比较原始，但对于小数据量的文件，这种方法还是最直接，最方便的。但是对于大量的文件，这种方法就很慢了。

另一种是**全文检索(Full-text Search)：**即先建立索引，再对索引进行搜索。索引是从非结构化数据中提取出之后重新组织的信息。

### 什么是elasticsearch

Elasticsearch是一个基于Apache Lucene(TM)的开源的高扩展的分布式搜索引擎 。当然Elasticsearch并不仅仅是Lucene那么简单，他不仅包含了全文搜索功能，还可以进行以下工作：

- 分布式实时文件存储，并将每一个字段都编入索引，使其可以被搜索。
- 实时分析的分布式搜索引擎。
- 可以扩展到上百台服务器，处理PB级别的结构化或非结构化数据。

**ES的适用场景**

- 维基百科 全文检索、高亮、搜索推荐
- The Guardian(国外新闻网站) 用户行为日志（点击，浏览，收藏，评论）+社交网络数据（对某某新闻的相关看法），数据分析，给到每篇新闻文章的作者，让他们知道他的文章的公众反馈（好、坏、热门。。。）
- Stack Overflow(国外程序异常讨论论坛)，全文检索，搜索到相关问题和答案，如果程序报错了，就会将报错信息粘贴到里面去，搜索有没有对应的答案
- github，搜索上千亿行的代码
- 电商网站，检索商品
- 日志数据的分析 elk技术（使用最多，将日志以可视化的方式呈现）
- 商品价格监控网站，用户设定某商品的价格阈值，当低于该阈值的时候，发送通知消息给用户
- BI系统，商业智能Business Intelligence。比如有个大型商场集团，BI，分析一下某某地区最近3年的用户消费金额的趋势以及用户群体的组成构成，产出相关的数张报表。



**ELASTICSEARCH的特点**

1. 可以作为大型分布式集群（数百台服务器）技术，处理PB级的数据，服务大公司；也可以运行在单机上服务于小公司
2. Elasticsearch不是什么新技术，主要是将全文检索、数据分析以及分布式技术，合并在了一起，才形成了独一无二的ES：lucene(全文检索)，商用的数据分析软件，分布式数据库
3. 对用户而言，是开箱即用的，非常简单，作为中小型应用，直接3分钟部署一下ES，就可以作为生产环境的系统来使用了，此时的场景是数据量不大，操作不是太复杂
4. （数据库的功能面对很多领域是不够用的（事务，还有各种联机事务型的操作）；特殊的功能，比如**全文检索**，**同义词处理**，**相关度排名**，**复杂数据分析**，**海量数据的近实时处理**，Elasticsearch作为传统数据库的一个补充，提供了数据库所不能提供的很多功能



## 2、 安装elasticsearch和kibana

### 关闭防火墙

```sh
#禁用防火墙
systemctl stop firewalld.service
#停止并禁用开机启动 
sytemctl disable firewalld.service
#查看防火状态
systemctl status firewalld.service
```

### 使用docker安装elasticsearch

使用docker就不需要配置jvm的运行环境了

```shell
#新建es的config配置文件夹
mkdir -p /data/elasticsearch/config
#新建es的data目录
mkdir -p /data/elasticsearch/data
#给目录设置权限
chmod 777 -R /data/elasticsearch

echo "http.host:0.0.0.0" >> /data/elasticsearch/config/elasticsearch.yml


docker run --name elasticsearch   -p 9200:9200 -p 9300:9300
-e "discovery.type=single-node"
-e  ES_JAVA_OPTS="-Xms128m -Xmx256m" 
-v /data/elasticsearch/config/elasticsearch.yml :/usr/share/elasticsearch/config/elasticsearch.yml
-v /data/elasticsearch/data:/usr/share/elasticsearch/data 
-v /data/elasticsearch/plugins:/usr/share/elasticsearch/plugins
-d elasticsearch:7.10.1
```

访问 http://192.168.0.104:9200 有json数据返回，说明启动成功了。



### 通过docker安装kibana

es相当于mysql ,kibana相当于navicate

```shell
docker run -d --name kibana -e ELASTICSEARCH_HOST="http://192.168.0.104:9200"  -p 5601:5601  kibana:7.10.1
```

访问 http://192.168.0.104:5601 能看到后台页面说明 kibana启动成功了。

Home/Dev tools  在Console 可以写 es 的查询命令。



## 3、 elasticsearch中的基本概念

### es中的type、index、mapping 和 dsl

| mysql    | elasticsearch                                 |
| -------- | --------------------------------------------- |
| database |                                               |
| table    | index (7.x开始理解为table，type为固定值_doc） |
| row      | document                                      |
| column   | field                                         |
| schema   | 映射（mapping）                               |
| sql      | DSL（Descriptor Struct Laguage）              |



### 索引

**有两个含义：动词(insert)，名词(表)**，插入一条数据，我们一般说 ”索引一条数据“。

Elasticsearch将它的数据存储在一个或多个索引（index）中。索引就像数据库，可以向索引写入文档或者从索引中读取文档。



## 4、 通过put和post方法添加数据

### 查看索引

kibana 控制台：GET  _cat/indice

Postman: 		  ` http://192.168.0.104:9200/_cat/indice `

### 通过PUT+id新建数据

保存id 为1 的数据，新版本固定为_doc，之前为type。

PUT请求，不存在的话新建，存在的话更新，返回 "result" : "update"

```http
PUT account/_doc/1 
{
  "name": "test1",
  "age": 12,
  "company": {
  		"name":"aaaa",
  		"address":"beijing"
  },
  "tags": ["aa", "bb", "cc"]
}
返回结果
{
  "_index" : "customer",
  "_type" : "_doc",
  "_id" : "1",
  "_version" : 1,
  "result" : "created",
  "_shards" : {   //分片信息
    "total" : 2,
    "successful" : 1,
    "failed" : 0
  },
  "_seq_no" : 0,   //乐观锁
  "_primary_term" : 1
}
```

### POST不带id新建数据

```http
POST /user/_doc
{
  "name": "John Doe"
}

POST /user/_doc/1
{
  "name": "John Doe"
}
```

post不带id,每请求一次，新建一次数据

post带id,就和put是一样的操作，put是不允许不带id的。



### post + _create

没有就创建，有就报错

```http
POST /user/_create/1
{
  "name": "John Doe"
}
```

### 查看某个索引的基本信息

GET  /account