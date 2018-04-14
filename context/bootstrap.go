package context

import (
	"github.com/gin-gonic/gin"
	"fmt"
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
		b.defaultService(c, []func(engine *gin.Engine){r})
	})
	b.cmd.Start()

}
func (b *boot) Run2(dir string, env string, r func(engine *gin.Engine), fun func(cfgMap ConfigMap)) {
	b.cmd.putExtend(fun)
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, []func(engine *gin.Engine){r})
	})
	b.cmd.StartConfigDir(dir, env)
}
func (b *boot) Run3(dir string, env string, fun func(cfgMap ConfigMap), fs ...func(engine *gin.Engine)) {
	b.cmd.putExtend(fun)
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, fs)
	})
	b.cmd.StartConfigDir(dir, env)
}
func (b *boot) RunWithExtend(r func(engine *gin.Engine), fun func(cfgMap ConfigMap)) {
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, []func(engine *gin.Engine){r})
	})
	b.cmd.putExtend(fun)

	b.cmd.Start()
}
func (b *boot) RunWithConfig(c Config) {
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, nil)
	})
	b.cmd.StartConfig(c)
}

func (b *boot) RunWithConfigDir(dir string, env string) {
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, nil)
	})
	b.cmd.StartConfigDir(dir, env)
}
func (b *boot) RunWithConfigDirAndExtend(dir string, env string, fun func(cfgMap ConfigMap)) {
	b.cmd.putExtend(fun)
	b.cmd.defaultCmd(func(c Config) {
		b.defaultService(c, nil)
	})
	b.cmd.StartConfigDir(dir, env)
}

func (b *boot) defaultService(c Config, fs []func(engine *gin.Engine)) {
	NewLogger(c.Jewel.Log)
	db := Db{}
	db.Open(c)
	Services.ServiceMap[DB] = db
	if c.Jewel.Port > 0 {
		engine := gin.Default()
		registeries(fs)
		load(engine)
		engine.Run(fmt.Sprintf(":%d", + c.Jewel.Port))
	}
}
