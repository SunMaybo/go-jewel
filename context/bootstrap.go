package context

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/SunMaybo/go-jewel/router"
	"strconv"
)

func Run(r func(engine *gin.Engine)) {
	flagService := FlagService{
		Params: make(map[string]string),
		Cmd:    make(map[string]func(c Config))}
	flagService.Default(func(c Config) {
		defaultService(c)
	})
	flagService.PutCmd("start", func(c Config) {
		port := c.Jewel.Port
		if port <= 0 {
			port = 8080
		}
		if p, ok := flagService.Params["port"]; ok {

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
	flagService.Start()
}
func RunWithExtend(r func(engine *gin.Engine), fun func(c Config)) {
	flagService := FlagService{
		Params: make(map[string]string),
		Cmd:    make(map[string]func(c Config))}
	flagService.Default(func(c Config) {
		defaultService(c)
	})
	flagService.PutExtend(fun)
	flagService.PutCmd("start", func(c Config) {
		port := c.Jewel.Port
		if port <= 0 {
			port = 8080
		}
		if p, ok := flagService.Params["port"]; ok {

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
	flagService.Start()
}
func RunWithConfig(c Config) {
	flagService := FlagService{
		Params: make(map[string]string),
		Cmd:    make(map[string]func(c Config))}
	flagService.Default(func(c Config) {
		defaultService(c)
	})
	flagService.StartConfig(c)
}

func RunWithConfigDir(dir string, env string) {
	flagService := FlagService{
		Params: make(map[string]string),
		Cmd:    make(map[string]func(c Config)),
		Extend: make(map[string]func(c Config))}
	flagService.Default(func(c Config) {
		defaultService(c)
	})
	flagService.StartConfigDir(dir, env)
}
func RunWithConfigDirAndExtend(dir string, env string, fun func(c Config)) {
	flagService := FlagService{
		Params: make(map[string]string),
		Cmd:    make(map[string]func(c Config))}
	flagService.Default(func(c Config) {
		defaultService(c)
	})
	flagService.PutExtend(fun)
	flagService.StartConfigDir(dir, env)
}

func defaultService(c Config) {
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
