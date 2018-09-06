package context

import (
	"flag"
	"fmt"
	"github.com/SunMaybo/jewel-inject/inject"
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
	fmt.Println("                                                                          /$$")
	fmt.Println("                                                                          | $$")
	fmt.Println("          /$$$$$$   /$$$$$$          /$$  /$$$$$$  /$$  /$$  /$$  /$$$$$$ | $$")
	fmt.Println("         /$$__  $$ /$$__  $$ /$$$$$$|__/ /$$__  $$| $$ | $$ | $$ /$$__  $$| $$")
	fmt.Println("	| $$  \\ $$| $$  \\ $$|______/ /$$| $$$$$$$$| $$ | $$ | $$| $$$$$$$$| $$")
	fmt.Println("	| $$  | $$| $$  | $$        | $$| $$_____/| $$ | $$ | $$| $$_____/| $$")
	fmt.Println("	|  $$$$$$$|  $$$$$$/        | $$|  $$$$$$$|  $$$$$/$$$$/|  $$$$$$$| $$")
	fmt.Println("	\\____  $$ \\______/         | $$ \\_______/ \\_____/\\___/  \\_______/|__/")
	fmt.Println("	/$$  \\ $$             /$$  | $$")
	fmt.Println("	|  $$$$$$/            |  $$$$$$/")
	fmt.Println("	\\______/              \\______/")
	fmt.Println("    ::  go-jewel  ::  (V2.0.1)")
	properties := Properties{}
	dir = GetCurrentDirectory(dir)
	fileName := LoadFileName(dir)
	jewel := &JewelProperties{}
	properties.Load(fileName, jewel)
	if env != "" {
		jewel.Jewel.Profiles.Active = env
	}
	NewLogger(dir + "/log.xml")
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

	fmt.Println("=============================================================")
	fmt.Printf("             project:  %s                        \n", jewel.Jewel.Name)
	fmt.Printf("         environment:  %s                    \n", jewel.Jewel.Profiles.Active)
	fmt.Printf("                port:  %d                           \n", jewel.Jewel.Port)
	fmt.Println("=============================================================")

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
