## go-jewel文档描述

整合优秀的开源项目，实现快速启动golangweb项目，较为轻量。
### 整合的开源项目
1. 基于gorm实现关系型数据库操作
2. 基于go-redis实现数据库redis操作
3. 基于gin-gonic实现HttpServer功能
4. 基于seelog实现日志管理
5. 基于jewel-template实现HttpClient功能
6. 基于jewel-inject实现结构体依赖管理
7. 基于robfig定时任务管理

### 特性
1. 通过app.yml(yaml,json,xml) and app-env.yml 来管理配置文件
2. 提供方便的应用程序入口
3. 提供初始化加载服务的注册
4. 提供命令行注册服务，用于个性化启动参数配置
5. 提供全局结构体注入inject
6. 提供RestTemplate 用于restful接口调用，并扩展jsonrpc支持
7. 简单的基于cron 的定时任务支持

## 安装和入门
### 安装
```
go get -u github.com/SunMaybo/go-jewel
```
### 快速启动
```
boot := context.NewInstance().
start().BindHttp(func(engine *gin.Engine) {
 })

```

## 配置文件详解
### 配置文件描述
  1. 默认在启动目录下创建config文件夹用于存放文件
  2. 配置文件通过默认加载app.yml,通过指定环境来加载app-env.yml文件

### 数据库启动配置
1. 项目命名

 ```
 jewel：
     name: test
 ```

2.  环境配置

    ```
    jewel:
       profiles:
           active: www
    ```

3. mysql 配置

  ```
jewel:
   mysql: admin:mypass@tcp(127.0.0.1:3308)/tokenup-btcd? charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci
  ```
4. mysql全局配置配置

 	```
 	jewel:
   		 max-idle-conns: 10
   		 max-Open-conns: 100
   		 sql_show: true
	```

5. redis配置

  ```
  jewel:
     	redis:
     		 host: ""
      		 password: ""
      		 db: 0
  ```


## 结构体依赖管理(inject)
实现IOC依赖注入，仅仅支持全局依赖注入切为单例模式，注入的均采用指针注入

1. 注入采用tag标记，且空字符串代表通过type的name注入

```
type Stu struct{
 Person *Person `inject:""`
 Name string `yaml:"name"`
 Age int `yaml:"age"`
}
```
2. 添加到容器中,切必须在Start()函数执行前完成工作

```
boot.AddApply(&Stu{})

```

3. 扩展配置文件通过注入配置方式

```app.yml
name: xiaowang
 age: 34
```
```golang
 boot.AddApplyCfg(&Stu{})
```
## 扩展方法
扩展方法执行，必须在Start()函数执行前实现，切在Start()函数执行后执行

```
	boot.AddFun(func() {
		fmt.Println("Hello World !")
	})
```
```
	boot.AddAsyncFun(func() {
		fmt.Println("异步方法执行")
	})
```
##  命令行工具
## 定时任务
定时任务必须在Start()方法执行前实现，切在Start()函数执行后被启动

```
	boot.AddTask("task","*/1 * * * * ?", func() {
		fmt.Println("a task...")
	})
```
## HttpClient

整合[jewel-template](https://github.com/SunMaybo/jewel-template)提供restful，jsonrpc接口调用,你需要通过注入*rest.RestTemplate模版

```
	boot.AddApply(rest.Default())
```

## 日志
 整合[seelog](https://github.com/cihub/seelog)进行日志输出，你可以通过家在日志配置文件

```
jewel:
   log: ./config/log.xml
```
## 关系数据库操作
整合[gorm](https://github.com/jinzhu/gorm)，通过context获取*grom.DB

```golang
db:=context.Services.Db().MysqlDb
```
## redis数据库操作
整合[go-redis](https://github.com/go-redis/redis)，通过context获取*redis.Client

```golang
db:=context.Services.Db().RedisDb
```
## http-server
整合[gin](https://github.com/gin-gonic/gin)，包装jsonrpc,启动服务

```
boot.Start().BindHttp(func(engine *gin.Engine) {
	})
```
## 例子
