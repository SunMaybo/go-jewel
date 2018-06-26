package context

import (
	"os"
	"errors"
	"flag"
	"fmt"
	"github.com/robfig/cron"
)

type Cmd struct {
	Params map[string]*string
	Cmd    map[string]func(c Config)
}

func (c *Cmd) Load() error {
	if len(os.Args) < 2 {
		return errors.New("too less cmd")
	}
	return nil
}
func (c *Cmd) PutFlagString(name string, value string, usage string) {
	e := flag.String(name, value, usage)
	c.Params[name] = e
}

func (c *Cmd) PutCmd(name string, fun func(c Config)) {
	c.Cmd[name] = fun
}
func (c *Cmd) defaultCmd(fun func(c Config)) {
	c.Cmd["default"] = fun
}
func (c *Cmd) httpCmd(fun func(c Config)) {
	c.Cmd["http"] = fun
}
func (c *Cmd) Start(b *boot) {
	c.PutFlagString("e", "", "startup environment")
	flag.Parse()
	cmd := flag.Arg(0)
	if cmd != "" {
		fmt.Printf("action: %s\n", cmd)
		fmt.Printf("env: %s\n", *c.Params["e"])
	}

	fmt.Printf("-------------------------------------------------------\n")
	dir := getCurrentDirectory()
	cfg := Load(dir)
	env := cfg.Jewel.Profiles.Active
	for _, v := range b.cfgPointer {
		LoadCfg(dir, v)
		LoadEnvCfg(dir, env, v)
	}
	LoadEnvCfg(dir, env, &cfg)
	// 注册配置
	b.inject.Apply(b.cfgPointer ...)
	b.inject.Apply(&cfg)
	//注册依赖
	b.inject.Apply(b.injector...)
	c.Cmd["default"](cfg) //默认的方法
	if fun, ok := c.Cmd[cmd]; ok {
		fun(cfg)
	} else {
		fmt.Println("cmd no found")
	}
	b.inject.Apply(Services.Db().MysqlDb,Services.Db().RedisDb)
	b.inject.Inject() //依赖扫描于加载
	for _, f := range b.asyncFuns {
		go func() {
			f()
		}()
	}
	for _, f := range b.funs {
		f()
	}
	for _, task := range b.taskfun {
		c := cron.New()
		c.AddFunc(task.Cron, func() {
			task.Fun()
		})
		c.Start()

	}
	c.Cmd["http"](cfg) //默认的方法

}
