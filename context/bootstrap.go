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
	funs       []func()
	asyncFuns  []func()
	Cron       *cron.Cron
}

var b *Boot

func NewInstance() *Boot {
	cron := cron.New()
	cron.Start()
	boot := &Boot{
		Cron: cron,
	}
	boot.cmd = Cmd{
		Params: make(map[string]*string),
		Cmd:    make(map[string]func(c JewelProperties)),
	}
	boot.inject = inject.New()
	b = boot
	return boot
}
func GetBoot() *Boot {
	return b
}
func (b *Boot) GetInject() *inject.Injector {
	return b.inject
}
func (b *Boot) GetCmd() *Cmd {
	return &b.cmd
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
func (b *Boot) AddFun(fun func()) *Boot {
	b.funs = append(b.funs, fun)
	return b
}
func (b *Boot) AddAsyncFun(fun func()) *Boot {
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
	b.cmd.defaultCmd(func(c JewelProperties) {
		b.defaultService(c, c.Jewel.Profiles.Active)
	})
	b.cmd.StartAndDir(b, dir)
	return b
}
func (b *Boot) BindHttp(r func(engine *gin.Engine)) {

	b.cmd.httpCmd(func(c JewelProperties) {
		b.http(c, []func(engine *gin.Engine){r})
	})
	b.cmd.Http(b)
}

func (b *Boot) defaultService(c JewelProperties, env string) {
	db := NewDb()
	err := db.Open(c)
	if err != nil {
		seelog.Error(err)
		seelog.Flush()
		os.Exit(-1)
	}
	if db.RedisDb != nil {
		for name, client := range db.RedisDb {
			b.inject.ApplyWithName(name, client)
		}
	}
	if db.MysqlDb != nil {
		for name, mysql := range db.MysqlDb {
			b.inject.ApplyWithName(name, mysql)
		}
	}
	if db.PostDb != nil {
		for name, postgres := range db.PostDb {
			b.inject.ApplyWithName(name, postgres)
		}
	}
	if db.MgoDb != nil {
		for name, postgres := range db.PostDb {
			b.inject.ApplyWithName(name, postgres)
		}
	}
	if db.RestTemplate != nil {
		for name, restTemplate := range db.RestTemplate {
			b.inject.ApplyWithName(name, restTemplate)
		}
	}
	Services.ServiceMap[DB] = db
}
func (b *Boot) http(c JewelProperties, fs []func(engine *gin.Engine)) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	b.defaultRouter(engine, c.Jewel.Profiles.Active, c.Jewel.Port, time.Now().String(), c.Jewel.Name)
	registeries(fs)
	load(engine)
	engine.Run(fmt.Sprintf(":%d", + c.Jewel.Port))
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
}

/*
获取程序运行路径
*/
func getCurrentDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pwd + "/config"
}
