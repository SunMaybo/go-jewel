package context

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/SunMaybo/go-jewel/router"
	"strconv"
)

type BootStrap struct {
	flagService FlagService
}

func (b *BootStrap) Run(r func(engine *gin.Engine)) {
	b.flagService = FlagService{
		Params:  make(map[string]string),
		Cmd:     make(map[string]func(c Config)),
		Extends: nil}
	b.flagService.Default(func(c Config) {
		b.defaultService(c)
	})
	b.flagService.PutCmd("start", func(c Config) {
		port := c.Jewel.Port
		if port <= 0 {
			port = 8080
		}
		if p, ok := b.flagService.Params["port"]; ok {

			cmdPort, _ := strconv.Atoi(p)
			if cmdPort > 0 {
				port = cmdPort
			}

		}
		engine := gin.Default()
		router.Register(r)
		router.Load(engine)
		engine.Run(fmt.Sprintf(":%d", + port))
	})
	b.flagService.Start()
}
func (b *BootStrap) RunWithExtend(r func(engine *gin.Engine), fun func(cfgMap ConfigMap)) {
	b.flagService = FlagService{
		Params:  make(map[string]string),
		Cmd:     make(map[string]func(c Config)),
		Extends: nil}
	b.flagService.Default(func(c Config) {
		b.defaultService(c)
	})
	b.flagService.PutExtend(fun)
	b.flagService.PutCmd("start", func(c Config) {
		port := c.Jewel.Port
		if port <= 0 {
			port = 8080
		}
		if p, ok := b.flagService.Params["port"]; ok {

			cmdPort, _ := strconv.Atoi(p)
			if cmdPort > 0 {
				port = cmdPort
			}

		}
		engine := gin.Default()
		router.Register(r)
		router.Load(engine)
		engine.Run(fmt.Sprintf(":%d", + port))
	})
	b.flagService.Start()
}
func (b *BootStrap) RunWithConfig(c Config) {
	b.flagService = FlagService{
		Params:  make(map[string]string),
		Cmd:     make(map[string]func(c Config)),
		Extends: nil}
	b.flagService.Default(func(c Config) {
		b.defaultService(c)
	})
	b.flagService.StartConfig(c)
}

func (b *BootStrap) RunWithConfigDir(dir string, env string) {
	b.flagService = FlagService{
		Params:  make(map[string]string),
		Cmd:     make(map[string]func(c Config)),
		Extends: nil}
	b.flagService.Default(func(c Config) {
		b.defaultService(c)
	})
	b.flagService.StartConfigDir(dir, env)
}
func (b *BootStrap) RunWithConfigDirAndExtend(dir string, env string, fun func(cfgMap ConfigMap)) {
	b.flagService = FlagService{
		Params:  make(map[string]string),
		Cmd:     make(map[string]func(c Config)),
		Extends: nil}
	b.flagService.PutExtend(fun)
	b.flagService.Default(func(c Config) {
		b.defaultService(c)
	})
	b.flagService.StartConfigDir(dir, env)
}

func (b *BootStrap) defaultService(c Config) {
	//1. 日志
	//log := Logger{}
	//see := log.GetLogger(c.Jewel.Log)
	//Services.ServiceMap[LOG] = see
	//2. 数据库
	NewLogger(c.Jewel.Log)
	db := Db{}
	db.Open(c)
	Services.ServiceMap[DB] = db
}
