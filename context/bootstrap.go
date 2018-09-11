package context

import (
	"github.com/gin-gonic/gin"
	"time"
	"github.com/cihub/seelog"
	"net/http"
	"os"
	"github.com/SunMaybo/jewel-inject/inject"
	"reflect"
	"log"
	"github.com/robfig/cron"
	"html/template"
	"github.com/SunMaybo/go-jewel/prometheus"
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
		log.Fatalf("param must be pointer type")
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
			seelog.Error(err)
			seelog.Flush()
			os.Exit(-1)
		}
		name, inter := plugin.Interface()
		b.GetInject().ApplyWithName("plugin:"+name, inter)
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
		seelog.Error(err)
		seelog.Flush()
		os.Exit(-1)
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
		for name, postgres := range base.PostDb {
			b.inject.ApplyWithName("mgo."+name, postgres)
		}
	}
	if base.RestTemplate != nil {
		for name, restTemplate := range base.RestTemplate {
			b.inject.ApplyWithName("rest."+name, restTemplate)
		}
	}
	name, inter := base.Interface()
	b.GetInject().ApplyWithName("plugin:"+name, inter)
}
func (b *Boot) http(fs []func(router *gin.RouterGroup, injector *inject.Injector)) {
	var jewel JewelProperties
	jewel = b.GetInject().Service(&jewel).(JewelProperties)
	serverProperties := jewel.Jewel.Server
	server, err := serverProperties.Create()
	if err != nil {
		log.Fatal(err)
	}
	engine := gin.Default()
	if serverProperties.GinMode != nil {
		gin.SetMode(*serverProperties.GinMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
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
	b.defaultRouter(router, jewel.Jewel.Profiles.Active, *serverProperties.Port, time.Now().String(), jewel.Jewel.Name)
	registeries(fs)
	load(router, b.GetInject())
	server.Handler = engine
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
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
		}
		for _, service := range services {
			plugin := service.(Plugin)
			err := plugin.Health()
			if err != nil {
				context.JSON(http.StatusOK, gin.H{"status": "DOWN", "message": err.Error()})
			}
		}
		context.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
}
