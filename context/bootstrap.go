package context

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"encoding/json"
	"time"
	"github.com/cihub/seelog"
	"net/http"
	"github.com/SunMaybo/go-jewel/jsonrpc"
	"os"
	"github.com/SunMaybo/jewel-inject/inject"
	"reflect"
	"log"
)

var methodMap jsonrpc.MethodMap

type boot struct {
	inject     inject.Injector
	cfgPointer []interface{}
	injector   []interface{}
	cmd        Cmd
	taskfun    []Cron
	funs       []func()
	asyncFuns  []func()
}

func NewInstance() boot {
	boot := boot{}
	boot.cmd = Cmd{
		Params: make(map[string]*string),
		Cmd:    make(map[string]func(c Config)),
	}
	boot.inject = inject.New()
	return boot
}
func (b *boot) GetInject() inject.Injector {
	return b.inject
}
func (b *boot) GetCmd() *Cmd {
	return &b.cmd
}

func (b *boot) AddApply(pointer ... interface{}) *boot {
	checkPointer(pointer)
	b.injector = append(b.injector, pointer...)
	return b
}
func (b *boot) AddTask(name, cron string, fun func()) *boot {
	b.taskfun = append(b.taskfun, Cron{Name: name, Cron: cron, Fun: fun})
	return b
}
func (b *boot) AddFun(fun func()) *boot {
	b.funs = append(b.funs, fun)
	return b
}
func (b *boot) AddAsyncFun(fun func()) *boot {
	b.asyncFuns = append(b.asyncFuns, fun)
	return b
}
func checkPointer(pointer interface{}) {
	if reflect.TypeOf(pointer).Kind() != reflect.Ptr {
		log.Fatal("param must be pointer type")
	}
}
func (b *boot) AddApplyCfg(pointers ... interface{}) *boot {
	//映射配置文件内容
	for e := range pointers {
		checkPointer(pointers[e])
	}
	b.cfgPointer = append(b.cfgPointer, pointers...)
	return b
}
func (b *boot) JsonRpc() *boot {
	methodMap = make(jsonrpc.MethodMap)
	return b
}

func (b *boot) RegisterJsonRpc(name string, method interface{}) {
	methodMap.Register(name, method)
}

func (b *boot) RunJsonRpc(relativePath string, r func(engine *gin.Engine)) {
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, c.Jewel.Profiles.Active, c.Jewel.Port)

	})
	b.cmd.httpCmd(func(c Config) {
		if c.Jewel.JsonRpc.Enabled {
			b.http(c, []func(engine *gin.Engine){b.jsonRpc(relativePath, c.Jewel.JsonRpc.UserName, c.Jewel.JsonRpc.Password), r}, c.Jewel.Profiles.Active, c.Jewel.Port)
		} else {
			b.http(c, []func(engine *gin.Engine){r}, c.Jewel.Profiles.Active, c.Jewel.Port)
		}
	})
	b.cmd.Start(b)

}
func (b *boot) jsonRpc(relativePath string, username string, password string) func(engine *gin.Engine) {
	return func(engine *gin.Engine) {
		engine.POST(relativePath, func(context *gin.Context) {
			auth := context.GetHeader("Authorization")
			basicAuth := jsonrpc.BaseAuth(username, password)
			if auth != basicAuth {
				context.String(http.StatusUnauthorized, "authorization error")
				return
			}
			request := jsonrpc.Request{}
			context.BindJSON(&request)
			if request.JsonRpc != "2.0" {
				context.String(http.StatusBadRequest, "param error")
				return
			}
			if request.Method == "" {
				context.String(http.StatusBadRequest, "method is not nil")
				return
			}
			resp := methodMap.Call(request.Method, request.Params)
			resp.Id = request.Id
			context.JSON(http.StatusOK, resp)
		})
	}
}
func (b *boot) Run(r func(engine *gin.Engine)) {
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, c.Jewel.Profiles.Active, c.Jewel.Port)
	})
	b.cmd.httpCmd(func(c Config) {
		b.http(c, []func(engine *gin.Engine){r}, c.Jewel.Profiles.Active, c.Jewel.Port)
	})
	b.cmd.Start(b)
}

func (b *boot) defaultService(c Config, env string, port int) {
	NewLogger(c.Jewel.Log)
	db := Db{}
	err := db.Open(c)
	if err != nil {
		seelog.Error(err)
		seelog.Flush()
		return
	}
	Services.ServiceMap[DB] = db
}
func (b *boot) http(c Config, fs []func(engine *gin.Engine), env string, port int) {
	if c.Jewel.Port > 0 {
		engine := gin.Default()
		b.defaultRouter(engine, env, port, time.Now().String(), c.Jewel.Name)
		registeries(fs)
		load(engine)
		engine.Run(fmt.Sprintf(":%d", + c.Jewel.Port))
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

func (b *boot) defaultRouter(engine *gin.Engine, env string, port int, bootTime string, name string) {
	engine.GET("/info", func(context *gin.Context) {
		db := DbTx{}
		if Services.Db().MysqlDb != nil && Services.Db().MysqlDb.Error == nil {
			db.MysqlDb = "UP"
		}
		if Services.Db().PostDb != nil && Services.Db().PostDb.Error == nil {
			db.PostDb = "UP"
		}
		if Services.Db().SqlServerDb != nil && Services.Db().SqlServerDb.Error == nil {
			db.SqlServerDb = "UP"
		}
		if Services.Db().Sqlite3Db != nil && Services.Db().Sqlite3Db.Error == nil {
			db.Sqlite3Db = "UP"
		}
		if Services.Db().RedisDb != nil {
			db.RedisDb = "UP"
		}

		info := Info{
			Db:       db,
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
