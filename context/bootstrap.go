package context

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"encoding/json"
	"time"
	"github.com/cihub/seelog"
	"net/http"
	"os"
	"github.com/SunMaybo/jewel-inject/inject"
	"reflect"
	"log"
	"github.com/robfig/cron"
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
		log.Fatal("param must be pointer type")
	}
}
func (b *Boot) AddApplyCfg(pointers ... interface{}) *Boot {
	//映射配置文件内容
	for e := range pointers {
		checkPointer(pointers[e])
	}
	b.cfgPointer = append(b.cfgPointer, pointers...)
	return b
}
func (b *Boot) StartAndDir(dir string) (*Boot) {
	b.cmd.defaultCmd(func() {
		b.basePluginService()

		b.pluginService()

	})
	b.cmd.StartAndDir(b, dir)
	return b
}
func (b *Boot) BindHttp(r func(engine *gin.Engine)) {

	b.cmd.httpCmd(func() {
		b.http([]func(engine *gin.Engine){r})
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
		b.inject.ApplyWithName("plugin:"+name, inter)
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
			b.inject.ApplyWithName("plugin:redis."+name, client)
		}
	}
	if base.MysqlDb != nil {
		for name, mysql := range base.MysqlDb {
			b.inject.ApplyWithName("plugin:mysql."+name, mysql)
		}
	}
	if base.PostDb != nil {
		for name, postgres := range base.PostDb {
			b.inject.ApplyWithName("plugin:postgres."+name, postgres)
		}
	}
	if base.MgoDb != nil {
		for name, postgres := range base.PostDb {
			b.inject.ApplyWithName("plugin:mgo."+name, postgres)
		}
	}
	if base.RestTemplate != nil {
		for name, restTemplate := range base.RestTemplate {
			b.inject.ApplyWithName("plugin:rest."+name, restTemplate)
		}
	}
	b.inject.Apply(&base)

}
func (b *Boot) http(fs []func(engine *gin.Engine)) {
	var jewel JewelProperties
	jewel = b.GetInject().Service(&jewel).(JewelProperties)
	if jewel.Jewel.GinMode != nil {
		gin.SetMode(*jewel.Jewel.GinMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	b.defaultRouter(engine, jewel.Jewel.Profiles.Active, jewel.Jewel.Port, time.Now().String(), jewel.Jewel.Name)
	registeries(fs)
	load(engine)
	err := engine.Run(fmt.Sprintf(":%d", + jewel.Jewel.Port))
	if err != nil {
		log.Fatal(err)
	}
}

type Info struct {
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Env      string `json:"env"`
	BootTime string `json:"boot_time"`
	Db       DbTx   `json:"db"`
}

type DbTx struct {
	MysqlDb     string `json:"mysql_db"`
	PostDb      string `json:"postgresql_db"`
	SqlServerDb string `json:"sqlserver_db"`
	Sqlite3Db   string `json:"sqlite3_db"`
	MongodDb    string `json:"mongod_db"`
	RedisDb     string `json:"redis_db"`
}

func (b *Boot) defaultRouter(engine *gin.Engine, env string, port int, bootTime string, name string) {
	engine.GET("/info", func(context *gin.Context) {
		info := Info{
			Port:     port,
			Env:      env,
			BootTime: bootTime,
			Name:     name,
		}
		buff, err := json.Marshal(info)
		if err != nil {
			seelog.Error(err)
		}
		context.String(http.StatusOK, "%v", string(buff))
	})
	engine.GET("/healths", func(context *gin.Context) {
		services := b.GetInject().ServiceByPrefixName("plugin:")
		if services == nil {
			context.JSON(http.StatusOK, gin.H{"status": "UP"})
		}
		for _, service := range services {
			plugin := service.(Plugin)
			err := plugin.Health()
			if err != nil {
				context.JSON(http.StatusOK, gin.H{"status": "DOWN"})
			}
		}
		context.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
}
