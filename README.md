## go-jewel
go-jewel 是一个集成gin，gorm，seelog的框架，通过配置文件可以快的搭载框架，方便开发。
### 特征
 1. 通过加载app.yml or app_env.yml 来加载gorm，gin以及通过指定日志配置文件log.xml做指定日志搜集。
 2. 提供方便的应用程序入口。
 3. 将提供多数据源支持。
 4. 提供初始化加载服务的注册。
 5. 提供命令行注册服务，用于个性化启动参数配置。
 6. 提供全局结构体注入inject
 7. 提供RestTemplate 用于restful接口调用，并扩展jsonrpc支持.
 8. 简单的基于cron 的定时任务支持
### 快速开始

 ```
 boot := context.NewInstance().start()

 ```
 ### 你想启动http-server,你需要这样做
 ```
 boot := context.NewInstance().
 start().
 BindHttp(func(engine *gin.Engine) {
 })
 ```
 ### 提供全局同步，异步方法支持
 ```
  boot.AddAsyncFun(func() {
 		
 	})
 	boot.AddFun(func() {
 	})
 ```
 ### 提供定时任务支持
 ```
 	boot.AddTask("task_name","*/1 * * * * ?", func() {
 		
 	})
 ```
 ### inject依赖注入支持
 #### mapping 配置文件
 1. golang实现
 ```
	type SystemConfig struct {
		Name    string `yaml:"name"`
		Ip      string `yaml:"ip"`
		Port    int    `yaml:"port"`
		Version string `yaml:"version"`
	}
 boot.AddApplyCfg(&SystemConfig{})
 ```
 2. 配置文件app.yml
 ```
     name: test_project
       ip: 127.0.0.1
     port: 8080
  version: 1.0.0  
 ```
 #### 结构体inject管理
 
 1. golang实现申请管理
 
 ```
 boot.AddApply(
 		&worker.TaskWorker{},
 		rest.Default(),
 	)
 	
 ```
 
 2. golang 获取服务
 
  ```
  rest:=boot.GetInject().Service(rest.Default()).(rest.RestTemplate)
  ```
 
 3. 结构体之间依赖，通过tag实现
 其中""代表，使用type的name依赖，也可以手动命名，推荐使用type的name自动注入
 
 ```
 type UserService struct {
 	UserSafeDao *dao.UserSafeDb `inject:""`
 }
 ```
### 全局配置文件app.yml
根据配置自动启动，http-server,mysql 以及redis
```
jewel:
   profiles:
        active: dev
   port: 7089
   log: ./config/log.xml
   mysql: root:mypass@tcp(127.0.0.1:3306)/tokenup-risk-controller?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci
   sql_show: true
   redis:
      host: 54.255.249.22:6379
      db: 0
   max-idle-conns: 10
   max-open-conns: 100
```
通过 jewel.profiles.active指定的环境名来扩展配置文件 列子:app-dev.yml
```
     name: test_project
       ip: 127.0.0.1
     port: 8080
  version: 1.0.0 
```
 
 ### 安装
 
 你可以通过一下命令安装
 
 ```
    go get -u github.com/SunMaybo/go-jewel
 ```
