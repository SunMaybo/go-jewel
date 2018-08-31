package context

import (
	"flag"
	"fmt"
)

type Cmd struct {
	Params map[string]*string
	Cmd    map[string]func(c JewelProperties)
}

func (c *Cmd) PutFlagString(name string, value string, usage string) {
	e := flag.String(name, value, usage)
	c.Params[name] = e
}

func (c *Cmd) PutCmd(name string, fun func(c JewelProperties)) {
	c.Cmd[name] = fun
}
func (c *Cmd) defaultCmd(fun func(c JewelProperties)) {
	c.Cmd["default"] = fun
}
func (c *Cmd) httpCmd(fun func(c JewelProperties)) {
	c.Cmd["http"] = fun
}
func (c *Cmd) StartAndDir(b *Boot, dir string) {
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
	fmt.Println("    ::  go-jewel  ::  (V2.0.0)")
	properties := Properties{}
	fileName := LoadFileName(dir)
	jewel := &JewelProperties{}
	properties.Load(fileName, jewel)
	env := jewel.Jewel.Profiles.Active
	NewLogger(dir + "/log.xml")
	for _, v := range b.cfgPointer {
		LoadCfg(dir, v)
		LoadEnvCfg(dir, env, v)
	}
	envFileName := LoadEnvFileName(dir, env)
	properties.Load(envFileName, jewel)

	fmt.Println("=============================================================")
	fmt.Printf("             project:  %s                        \n", jewel.Jewel.Name)
	fmt.Printf("         environment:  %s                    \n", jewel.Jewel.Profiles.Active)
	fmt.Printf("                port:  %d                           \n", jewel.Jewel.Port)
	fmt.Println("=============================================================")
	c.Cmd["default"](*jewel) //默认的方法
	b.inject.Apply(b.cfgPointer ...)
	b.inject.Apply(jewel)
	b.inject.Apply(b.injector...)
	b.inject.Inject() //依赖扫描于加载
	fmt.Println(len(b.asyncFuns))
	for e := range b.asyncFuns {
		go func(fun func()) {
			fun()
		}(b.asyncFuns[e])
	}
	for _, f := range b.funs {
		f()
	}
}
func (c *Cmd) Http(b *Boot) () {
	cfg := b.inject.Service(&JewelProperties{}).(JewelProperties)
	c.Cmd["http"](cfg) //默认的方法
}
