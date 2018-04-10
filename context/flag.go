package context

import (
	"os"
	"errors"
	"flag"
	"strconv"
	"time"
	"fmt"
)

type FlagService struct {
	Params map[string]string
	Cmd    map[string]func(c Config)
}

func (f *FlagService) Load() error {
	if len(os.Args) < 2 {
		return errors.New("too less cmd")
	}
	return nil
}
func (f *FlagService) PutFlagString(name string, value string, usage string) {
	e := flag.String(name, value, usage)
	f.Params[name] = *e
}
func (f *FlagService) PutFlagInt(name string, value int, usage string) {
	e := flag.Int(name, value, usage)
	f.Params[name] = string(*e)
}
func (f *FlagService) PutFlagDuration(name string, value time.Duration, usage string) {
	e := flag.Duration(name, value, usage)
	f.Params[name] = string(*e)
}
func (f *FlagService) PutFlagFloat64(name string, value float64, usage string) {
	e := flag.Float64(name, value, usage)
	f.Params[name] = strconv.FormatFloat(*e, 'E', -1, 64)
}
func (f *FlagService) PutCmd(name string, fun func(c Config)) {
	f.Cmd[name] = fun
}
func (f *FlagService) Default(fun func(c Config)) {
	f.Cmd["default"] = fun
}
func (f *FlagService) Start() {
	f.PutFlagString("env", "www", "startup environment")
	cmd := flag.Arg(0)
	cfg := Load("./config", f.Params["env"])
	f.Cmd["default"](cfg) //默认的方法
	fun := f.Cmd[cmd]
	if fun != nil {
		fun(cfg)
	} else {
		fmt.Println("cmd no found")
	}
}
func (f *FlagService) StartConfig(cfg Config) {
	f.Cmd["default"](cfg) //默认的方法
}
func (f *FlagService) StartConfigDir(dir string, env string) {
	cfg := Load(dir, f.Params["env"])
	f.Cmd["default"](cfg) //默认的方法
}
