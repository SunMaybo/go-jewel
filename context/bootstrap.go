package context

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"encoding/json"
	"time"
	"github.com/cihub/seelog"
	"net/http"
)

type boot struct {
	cmd Cmd
}

func NewInstance() boot {
	boot := boot{}
	boot.cmd = Cmd{
		Params:  make(map[string]string),
		Cmd:     make(map[string]func(c Config)),
		Extends: nil}
	return boot
}
func (b *boot) GetCmd() *Cmd {
	return &b.cmd
}

func (b *boot) Run(r func(engine *gin.Engine)) {
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, []func(engine *gin.Engine){r}, c.Jewel.Profiles.Active, c.Jewel.Port)
	})
	b.cmd.Start()

}
func (b *boot) Run2(dir string, env string, r func(engine *gin.Engine), fun func(cfgMap ConfigMap)) {
	b.cmd.putExtend(fun)
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, []func(engine *gin.Engine){r}, env, c.Jewel.Port)
	})
	b.cmd.StartConfigDir(dir, env)
}
func (b *boot) Run3(dir string, env string, fun func(cfgMap ConfigMap), fs ...func(engine *gin.Engine)) {
	b.cmd.putExtend(fun)
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, fs, c.Jewel.Profiles.Active, c.Jewel.Port)
	})
	b.cmd.StartConfigDir(dir, env)
}
func (b *boot) RunWithExtend(r func(engine *gin.Engine), fun func(cfgMap ConfigMap)) {
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, []func(engine *gin.Engine){r}, c.Jewel.Profiles.Active, c.Jewel.Port)
	})
	b.cmd.putExtend(fun)

	b.cmd.Start()
}
func (b *boot) RunWithConfig(c Config) {
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, nil, c.Jewel.Profiles.Active, c.Jewel.Port)
	})
	b.cmd.StartConfig(c)
}

func (b *boot) RunWithConfigDir(dir string, env string) {
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, nil, env, c.Jewel.Port)
	})
	b.cmd.StartConfigDir(dir, env)
}
func (b *boot) RunWithConfigDirAndExtend(dir string, env string, fun func(cfgMap ConfigMap)) {
	b.cmd.putExtend(fun)
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, nil, env, c.Jewel.Port)
	})
	b.cmd.StartConfigDir(dir, env)
}

func (b *boot) defaultService(c Config, fs []func(engine *gin.Engine), env string, port int) {
	NewLogger(c.Jewel.Log)
	db := Db{}
	db.Open(c)
	Services.ServiceMap[DB] = db
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
