package context

import (
	"flag"
	"github.com/SunMaybo/jewel-inject/inject"
	"log"
	"github.com/SunMaybo/go-jewel/logs"
)

type Cmd struct {
	Params map[string]*string
	Cmd    map[string]func()
}

func (c *Cmd) PutFlagString(name string, value string, usage string) {
	e := flag.String(name, value, usage)
	c.Params[name] = e
}

func (c *Cmd) PutCmd(name string, fun func()) {
	c.Cmd[name] = fun
}
func (c *Cmd) defaultCmd(fun func()) {
	c.Cmd["default"] = fun
}
func (c *Cmd) httpCmd(fun func()) {
	c.Cmd["http"] = fun
}
func (c *Cmd) Start(b *Boot, dir, env string) {
	properties := Properties{}
	dir, err := GetCurrentDirectory(dir)
	if err != nil {
		log.Fatal(err)
	}
	fileName := LoadFileName(dir)
	jewel := &JewelProperties{
		Jewel: Jewel{Name: "jewel-project", Profiles: Profiles{Active: "test"}, Log: Log{Level: "debug"}, Server: ServerProperties{Port: 8080}},
	}

	if fileName != "" {
		properties.Load(fileName, jewel)
	}
	if env == "" {
		env = jewel.Jewel.Profiles.Active
	} else {
		jewel.Jewel.Profiles.Active = env
	}
	logs.GetLog(jewel.Jewel.Log.Level)
	for _, v := range b.cfgPointer {
		LoadCfg(dir, v)
		LoadEnvCfg(dir, env, v)
	}
	envFileName := LoadEnvFileName(dir, env)
	if envFileName != "" {
		properties.Load(envFileName, jewel)
	}
	b.inject.Apply(b.cfgPointer ...)
	b.inject.Apply(jewel)
	c.Cmd["default"]() //默认的方法
	b.inject.Apply(jewel)
	b.inject.Apply(b.injector...)
	b.inject.Inject() //依赖扫描于加载
	for e := range b.asyncFuns {
		go func(fun func(inject *inject.Injector)) {
			fun(b.GetInject())
		}(b.asyncFuns[e])
	}
	for _, f := range b.funs {
		f(b.GetInject())
	}
}
func (c *Cmd) Http(b *Boot) () {
	c.Cmd["http"]() //默认的方法
}
