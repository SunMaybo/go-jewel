package context

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/SunMaybo/go-jewel/router"
	"strconv"
)

type BootStrap struct {
	cmd Cmd
}

func (b *BootStrap) Run(r func(engine *gin.Engine)) {
	b.cmd = Cmd{
		Params:  make(map[string]string),
		Cmd:     make(map[string]func(c Config)),
		Extends: nil}
	b.cmd.Default(func(c Config) {
		b.defaultService(c)
	})
	b.cmd.Start()

}
func (b *BootStrap) RunWithExtend(r func(engine *gin.Engine), fun func(cfgMap ConfigMap)) {
	b.cmd = Cmd{
		Params:  make(map[string]string),
		Cmd:     make(map[string]func(c Config)),
		Extends: nil}
	b.cmd.Default(func(c Config) {
		b.defaultService(c)
	})
	b.cmd.PutExtend(fun)
	b.cmd.PutCmd("start", func(c Config) {
		port := c.Jewel.Port
		if port <= 0 {
			port = 8080
		}
		if p, ok := b.cmd.Params["port"]; ok {

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
	b.cmd.Start()
}
func (b *BootStrap) RunWithConfig(c Config) {
	b.cmd = Cmd{
		Params:  make(map[string]string),
		Cmd:     make(map[string]func(c Config)),
		Extends: nil}
	b.cmd.Default(func(c Config) {
		b.defaultService(c)
	})
	b.cmd.StartConfig(c)
}

func (b *BootStrap) RunWithConfigDir(dir string, env string) {
	b.cmd = Cmd{
		Params:  make(map[string]string),
		Cmd:     make(map[string]func(c Config)),
		Extends: nil}
	b.cmd.Default(func(c Config) {
		b.defaultService(c)
	})
	b.cmd.StartConfigDir(dir, env)
}
func (b *BootStrap) RunWithConfigDirAndExtend(dir string, env string, fun func(cfgMap ConfigMap)) {
	b.cmd = Cmd{
		Params:  make(map[string]string),
		Cmd:     make(map[string]func(c Config)),
		Extends: nil}
	b.cmd.PutExtend(fun)
	b.cmd.Default(func(c Config) {
		b.defaultService(c)
	})
	b.cmd.StartConfigDir(dir, env)
}

func (b *BootStrap) defaultService(c Config, fs []func(engine *gin.Engine)) {
	//1. 日志
	//log := Logger{}
	//see := log.GetLogger(c.Jewel.Log)
	//Services.ServiceMap[LOG] = see
	//2. 数据库
	NewLogger(c.Jewel.Log)
	db := Db{}
	db.Open(c)
	Services.ServiceMap[DB] = db
	if c.Jewel.Port > 0 {
		engine := gin.Default()
		router.Registeries(fs)
		router.Load(engine)
		engine.Run(fmt.Sprintf(":%d", + c.Jewel.Port))
	}
}
