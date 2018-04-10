package go_jewel

import (
	"github.com/SunMaybo/go-jewel/context"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/SunMaybo/go-jewel/router"
	"strconv"
)

func Run(r func(engine *gin.Engine)) {
	flagService := context.FlagService{
		Params: make(map[string]string),
		Cmd:    make(map[string]func(c context.Config))}
	flagService.Default(func(c context.Config) {
		defaultService(c)
	})
	flagService.PutCmd("start", func(c context.Config) {
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

func RunWithConfig(c context.Config) {
	flagService := context.FlagService{
		Params: make(map[string]string),
		Cmd:    make(map[string]func(c context.Config))}
	flagService.Default(func(c context.Config) {
		defaultService(c)
	})
	flagService.StartConfig(c)
}

func RunWithConfigDir(dir string, env string) {
	flagService := context.FlagService{
		Params: make(map[string]string),
		Cmd:    make(map[string]func(c context.Config))}
	flagService.Default(func(c context.Config) {
		defaultService(c)
	})
	flagService.StartConfigDir(dir, env)
}

func defaultService(c context.Config) {
	//1. 日志
	//log := Logger{}
	//see := log.GetLogger(c.Jewel.Log)
	//Services.ServiceMap[LOG] = see
	//2. 数据库
	NewLogger(c.Jewel.Log)
	db := context.Db{}
	db.Open(c)
	context.Services.ServiceMap[context.DB] = db
}
