##go-jewel文档描述
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
### 配置文件秒睡
  1. 默认在启动目录下创建config文件夹用于存放文件
  2. 配置文件通过默认加载app.yml,通过指定环境来加载app-env.yml文件

## 结构提体依赖管理
## 扩展方法
##  命令行工具
## 定时任务
## HttpClient
## 日志
## 关系数据库操作
## redis数据库操作
## http-server
## 例子
