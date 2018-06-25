package context

import (
	"os"
	"errors"
	"flag"
	"fmt"
	"github.com/SunMaybo/jewel-inject/inject"
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
func (c *Cmd) Start(inject inject.Injector, cfgPointer []interface{}, injectors []interface{}) {
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
	for _, v := range cfgPointer {
		LoadCfg(dir, v)
		LoadEnvCfg(dir, env, v)
	}
	LoadEnvCfg(dir, env, &cfg)
	// 注册配置
	inject.Apply(cfgPointer ...)
	inject.Apply(&cfg)
	//注册依赖
	inject.Apply(injectors...)

	c.Cmd["default"](cfg) //默认的方法
	if fun, ok := c.Cmd[cmd]; ok {
		fun(cfg)
	} else {
		fmt.Println("cmd no found")
	}
	inject.Inject()    //依赖扫描于加载
	c.Cmd["http"](cfg) //默认的方法

}
