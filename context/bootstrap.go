package context

import (
	"time"
	"net/http"
	"github.com/SunMaybo/jewel-inject/inject"
	"reflect"
	"github.com/robfig/cron"
	"html/template"
	"github.com/SunMaybo/go-jewel/prometheus"
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
	"github.com/SunMaybo/go-jewel/logs"
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"strings"
)

type Boot struct {
	inject     *inject.Injector
	cfgPointer []interface{}
	injector   []interface{}
	cmd        Cmd
	funs       []func(injector *inject.Injector)
	asyncFuns  []func(injector *inject.Injector)
	Cron       *cron.Cron
	plugins    []Plugin
}

func NewInstance() *Boot {
	cron := cron.New()
	cron.Start()
	boot := &Boot{
		Cron: cron,
	}
	boot.cmd = Cmd{
		Params: make(map[string]*string),
		Cmd:    make(map[string]func()),
	}
	boot.inject = inject.New()
	return boot
}
func (b *Boot) GetInject() *inject.Injector {
	return b.inject
}
func (b *Boot) AddPlugins(plugins ... Plugin) {
	for e := range plugins {
		checkPointer(plugins[e])
	}
	b.plugins = append(b.plugins, plugins...)
}

func (b *Boot) AddApply(pointers ... interface{}) *Boot {
	for e := range pointers {
		checkPointer(pointers[e])
	}
	b.injector = append(b.injector, pointers...)
	return b
}
func (b *Boot) AddTask(name, cron string, fun func()) *Boot {
	b.Cron.AddFunc(cron, fun)
	return b
}
func (b *Boot) AddFun(fun func(injector *inject.Injector)) *Boot {
	b.funs = append(b.funs, fun)
	return b
}
func (b *Boot) AddAsyncFun(fun func(injector *inject.Injector)) *Boot {
	b.asyncFuns = append(b.asyncFuns, fun)
	return b
}
func checkPointer(pointer interface{}) {
	if reflect.TypeOf(pointer).Kind() != reflect.Ptr {
		zap.S().Fatal("param must be pointer type")
	}
}
func (b *Boot) AddApplyCfg(pointers ... interface{}) *Boot {
	//映射配置文件内容
	for _, ptr := range pointers {
		checkPointer(ptr)
	}
	b.cfgPointer = append(b.cfgPointer, pointers...)
	return b
}
func (b *Boot) StartAndDir(dir string) (*Boot) {
	b.cmd.defaultCmd(func() {
		b.basePluginService()
		b.pluginService()

	})
	b.cmd.Start(b, dir, "")
	return b
}
func (b *Boot) Start(dir, env string) (*Boot) {
	b.cmd.defaultCmd(func() {
		b.basePluginService()
		b.pluginService()

	})
	b.cmd.Start(b, dir, env)
	return b
}
func (b *Boot) BindHttp(r ... func(router *gin.RouterGroup, injector *inject.Injector)) {

	b.cmd.httpCmd(func() {
		b.http(r)
	})
	b.cmd.Http(b)
}
func (b *Boot) pluginService() {
	for _, plugin := range b.plugins {
		err := plugin.Open(b.GetInject())
		if err != nil {
			zap.S().Fatal(err)
		}
		name := plugin.InterfaceName()
		b.GetInject().ApplyWithName("plugin:"+name, &plugin)
	}
}
func (b *Boot) Close() {
	for _, plugin := range b.plugins {
		plugin.Close()
	}
}
func (b *Boot) basePluginService() {
	base := NewBasePlugin()
	err := base.Open(b.GetInject())
	if err != nil {
		zap.S().Fatal(err)
	}
	if base.RedisDb != nil {
		for name, client := range base.RedisDb {
			b.inject.ApplyWithName("redis."+name, client)
		}
	}
	if base.MysqlDb != nil {
		for name, mysql := range base.MysqlDb {
			b.inject.ApplyWithName("mysql."+name, mysql)
		}
	}
	if base.PostDb != nil {
		for name, postgres := range base.PostDb {
			b.inject.ApplyWithName("postgres."+name, postgres)
		}
	}
	if base.MgoDb != nil {
		for name, mgo := range base.MgoDb {
			b.inject.ApplyWithName("mgo."+name, mgo)
		}
	}
	if base.RestTemplate != nil {
		for name, restTemplate := range base.RestTemplate {
			b.inject.ApplyWithName("rest."+name, restTemplate)
		}
	}
	name := base.InterfaceName()
	var plugin Plugin
	plugin = base
	b.GetInject().ApplyWithName("plugin:"+name, &plugin)
}
func (b *Boot) http(fs []func(router *gin.RouterGroup, injector *inject.Injector)) {
	var jewel JewelProperties
	jewel = b.GetInject().Service(&jewel).(JewelProperties)
	serverProperties := jewel.Jewel.Server
	server, err := serverProperties.Create()
	if err != nil {
		zap.S().Fatal(err)
	}
	fmt.Printf("listen and serve on %s\n", server.Addr)
	gin.DisableConsoleColor()
	if serverProperties.GinMode != nil {
		gin.SetMode(*serverProperties.GinMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(Cors())
	engine.Use(logs.Logger(logs.LOGGER))
	engine.Use(gin.Recovery())
	if serverProperties.EnablePprof == nil || *serverProperties.EnablePprof {
		ginpprof.Wrap(engine)
	}
	var router *gin.RouterGroup
	if serverProperties.ContextPath != nil {
		router = engine.Group(*serverProperties.ContextPath)
	} else {
		router = engine.Group("/")
	}
	if serverProperties.Templates != nil {
		engine.SetHTMLTemplate(template.New(*serverProperties.Templates))
	}
	b.defaultRouter(router, jewel.Jewel.Profiles.Active, serverProperties.Port, time.Now().String(), jewel.Jewel.Name)
	registeries(fs)
	load(router, b.GetInject())
	server.Handler = engine
	err = server.ListenAndServe()
	if err != nil {
		zap.S().Fatal(err)
	}

}

//
////// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

type Info struct {
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Env      string `json:"env"`
	BootTime string `json:"boot_time"`
}

func (b *Boot) defaultRouter(engine *gin.RouterGroup, env string, port int64, bootTime string, name string) {
	engine.GET("/info", func(context *gin.Context) {
		info := Info{
			Port:     int(port),
			Env:      env,
			BootTime: bootTime,
			Name:     name,
		}
		context.JSON(http.StatusOK, info)
	})
	p := prometheus.NewPrometheus("gin")
	p.Use(engine)
	engine.GET("/healths", func(context *gin.Context) {
		services := b.GetInject().ServiceByPrefixName("plugin:")
		if services == nil {
			context.JSON(http.StatusOK, gin.H{"status": "UP"})
			return
		}
		for _, service := range services {
			plugin := service.(Plugin)
			err := plugin.Health()
			if err != nil {
				context.JSON(http.StatusOK, gin.H{"status": "DOWN", "message": err.Error()})
				return
			}
		}
		context.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
}
