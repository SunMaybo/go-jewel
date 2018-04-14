package context

import (
	"os"
	"errors"
	"flag"
	"strconv"
	"time"
	"fmt"
)

type Cmd struct {
	Params  map[string]string
	Cmd     map[string]func(c Config)
	Extends func(cfgMap ConfigMap)
}

func (c *Cmd) Load() error {
	if len(os.Args) < 2 {
		return errors.New("too less cmd")
	}
	return nil
}
func (c *Cmd) PutFlagString(name string, value string, usage string) {
	e := flag.String(name, value, usage)
	c.Params[name] = *e
}
func (c *Cmd) PutFlagInt(name string, value int, usage string) {
	e := flag.Int(name, value, usage)
	c.Params[name] = string(*e)
}
func (c *Cmd) PutFlagDuration(name string, value time.Duration, usage string) {
	e := flag.Duration(name, value, usage)
	c.Params[name] = string(*e)
}
func (c *Cmd) PutFlagFloat64(name string, value float64, usage string) {
	e := flag.Float64(name, value, usage)
	c.Params[name] = strconv.FormatFloat(*e, 'E', -1, 64)
}
func (c *Cmd) PutCmd(name string, fun func(c Config)) {
	c.Cmd[name] = fun
}
func (c *Cmd) putExtend(fun func(cfgMap ConfigMap)) {
	c.Extends = fun
}

func (c *Cmd) defaultCmd(fun func(c Config)) {
	c.Cmd["default"] = fun
}
func (c *Cmd) Start() {
	c.PutFlagString("e", "www", "startup environment")
	flag.Parse()
	cmd := flag.Arg(0)
	if cmd != "" {
		fmt.Printf("action: %s\n", cmd)
		fmt.Printf("env: %s\n", c.Params["e"])
	}

	fmt.Printf("-------------------------------------------------------\n")

	cfg := Load("./config", c.Params["env"])
	if fun, ok := c.Cmd["extend"]; ok {
		fun(cfg)
	}
	cfgMap := LoadMap("./config", c.Params["env"])
	if c.Extends != nil {
		c.Extends(cfgMap)
	}
	c.Cmd["default"](cfg) //默认的方法
	if fun, ok := c.Cmd[cmd]; ok {
		fun(cfg)
	} else {
		fmt.Println("cmd no found")
	}

}
func (c *Cmd) StartConfig(cfg Config) {
	if fun, ok := c.Cmd["extend"]; ok {
		fun(cfg)
	}
	c.Cmd["default"](cfg) //默认的方法

}
func (c *Cmd) StartConfigDir(dir string, env string) {
	cfg := Load(dir, env)
	if fun, ok := c.Cmd["extend"]; ok {
		fun(cfg)
	}
	cfgMap := LoadMap(dir, env)
	if c.Extends != nil {
		c.Extends(cfgMap)
	}
	c.Cmd["default"](cfg) //默认的方法
}
