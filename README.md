#go-jewel
go-jewel 是一个集成gin，gorm，seelog的框架，通过配置文件可以快的搭载框架，方便开发。
##特征
 1. 通过加载app_env.yml or app_env.json 来加载gorm，gin以及通过指定日志配置文件log.xml做指定日志搜集。
 2. 提供方便的应用程序入口。
 3. 将提供多数据源支持。
 4. 提供初始化加载服务的注册。
 5. 提供命令行注册服务，用于个性化启动参数配置。
 
 ## 快速开始
 ```
 boot := context.NewInstance()
 boot.Run(func(engine *gin.Engine) {
		engine.POST("/health", func(i *gin.Context) {
		})
	})
 ```
 ## 安装
 
 你可以通过一下命令安装
 
 ```
    go get -u github.com/SunMaybo/go-jewel
 ```
 ## 文档
 ## 例子
 ## 问题