package context

import (
	"os"
	"errors"
	"flag"
	"fmt"
	"github.com/robfig/cron"
	"strconv"
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
func (c *Cmd) StartAndDir(b *Boot, dir string) {
	c.PutFlagString("e", "", "startup environment")
	flag.Parse()
	cmd := flag.Arg(0)
	if cmd != "" {
		fmt.Printf("action: %s\n", cmd)
		fmt.Printf("env: %s\n", *c.Params["e"])
	}

	fmt.Printf("-------------------------------------------------------\n")
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
	b.inject.Apply(Services.Db().MysqlDb, Services.Db().RedisDb)
	b.inject.Inject() //依赖扫描于加载
	fmt.Println(len(b.asyncFuns))
	for e, _ := range b.asyncFuns {
		go func(fun func()) {
			fun()
		}(b.asyncFuns[e])
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

}
func (c *Cmd) Start(b *Boot) {
	c.PutFlagString("e", "", "startup environment")
	c.PutFlagString("p", "", "startup port")
	flag.Parse()
	cmd := flag.Arg(0)
	if cmd != "" {
		fmt.Printf("action: %s\n", cmd)
		fmt.Printf("env: %s\n", *c.Params["e"])
	}
	fmt.Printf("-------------------------------------------------------\n")
	b.inject.Apply(b.cfgPointer ...)
	//注册依赖
	b.inject.Apply(b.injector...)
	dir := getCurrentDirectory()
	cfg := Load(dir)

	var env string
	if *c.Params["e"] != "" {
		env = *c.Params["e"]
	} else if cfg.Jewel.Profiles.Active != "" {
		env = *c.Params["e"]
	} else {
		env = "default"
	}

	for _, v := range b.cfgPointer {
		LoadCfg(dir, v)
		LoadEnvCfg(dir, env, v)
	}
	LoadEnvCfg(dir, env, &cfg)

	cfg.Jewel.Profiles.Active = env
	if *c.Params["p"] != "" {
		cfg.Jewel.Port, _ = strconv.Atoi(*c.Params["p"])
	} else if cfg.Jewel.Port <= 0 {
		cfg.Jewel.Port = 8080
	}
	// 注册配置
	b.inject.Apply(&cfg)

	c.Cmd["default"](cfg) //默认的方法
	if fun, ok := c.Cmd[cmd]; ok {
		fun(cfg)
	}
	b.inject.Apply(Services.Db().MysqlDb, Services.Db().RedisDb)
	b.inject.Inject() //依赖扫描于加载
	for e, _ := range b.asyncFuns {
		go func(fun func()) {
			fun()
		}(b.asyncFuns[e])
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

}
func (c *Cmd) Http(b *Boot) () {
	cfg := b.inject.Service(&Config{}).(Config)
	c.Cmd["http"](cfg) //默认的方法
}
