package jewel

import (
	"github.com/alecthomas/kingpin"
	"os"
	"github.com/SunMaybo/go-jewel/context"
	"github.com/gin-gonic/gin"
	"github.com/SunMaybo/jewel-inject/inject"
	"github.com/SunMaybo/go-jewel/registry"
)

type Jewel struct {
	boot    *context.Boot
	app     *kingpin.Application
	cmd     []*kingpin.CmdClause
	cmdFunc map[string]func()
}

type CmdParam struct {
	Param map[string]interface{}
	Name  string
}

func NewHttp() *Jewel {
	jewel := &Jewel{
		boot:    context.NewInstance(),
		app:     kingpin.UsageTemplate(kingpin.ManPageTemplate),
		cmdFunc: make(map[string]func()),
	}
	jewel.boot.AddApplyCfg(&registry.JewelRegisterProperties{})
	var plugin context.Plugin
	plugin = &registry.EtcRegisterPlugin{}
	jewel.boot.AddPlugins(plugin)
	jewel.cmd = append(jewel.cmd, jewel.app.Command("server", "Start a http server"))
	return jewel
}

func New() *Jewel {
	jewel := &Jewel{
		boot:    context.NewInstance(),
		app:     kingpin.UsageTemplate(kingpin.DefaultUsageTemplate),
		cmdFunc: make(map[string]func()),
	}
	return jewel
}

func (jewel *Jewel) Cmd(name, help string, fun func()) *kingpin.CmdClause {
	c := jewel.app.Command(name, help)
	jewel.cmd = append(jewel.cmd, c)
	jewel.cmdFunc[name] = fun
	return c
}
func (jewel *Jewel) Boot() *context.Boot {
	return jewel.boot
}
func (jewel *Jewel) KingpinApp() *kingpin.Application {
	return jewel.app
}
func (jewel *Jewel) AddBeanProperties(properties ... interface{}) *Jewel {
	for _, p := range properties {
		jewel.boot = jewel.boot.AddApplyCfg(p)
	}
	return jewel
}

func (jewel *Jewel) AddBean(beans ...  interface{}) *Jewel {
	for _, bean := range beans {
		jewel.boot = jewel.boot.AddApply(bean)
	}
	return jewel
}

func (jewel *Jewel) AddAsyncFun(fun func(injector *inject.Injector)) *Jewel {
	jewel.boot = jewel.boot.AddAsyncFun(fun)
	return jewel
}

func (jewel *Jewel) AddSyncFun(fun func(injector *inject.Injector)) *Jewel {
	jewel.boot = jewel.boot.AddFun(fun)
	return jewel
}
func (jewel *Jewel) AddTask(name, cron string, fun func()) *Jewel {
	jewel.boot = jewel.boot.AddTask(name, cron, fun)
	return jewel
}

func (jewel *Jewel) AddPlugins(plugins ... context.Plugin) {
	jewel.boot.AddPlugins(plugins...)
}
func (jewel *Jewel) HttpStart(httpFun func(router *gin.RouterGroup, injector *inject.Injector)) {
	cmdParams := jewel.getCmdParam()
	c := kingpin.MustParse(jewel.app.Parse(os.Args[1:]))
	for _, cmd := range cmdParams {
		if cmd.Name == c && c == "server" {
			if cmd.Param["env"] == nil {
				jewel.boot = jewel.boot.Start(*cmd.Param["dir"].(*string), "")
			} else {
				jewel.boot = jewel.boot.Start(*cmd.Param["dir"].(*string), *cmd.Param["env"].(*string))

			}
			etcRegister := jewel.boot.GetInject().ServicePtrByName("plugin:etcd_register")
			if etcRegister != nil {
				reg := etcRegister.(*registry.EtcRegistry)
				registerOperation := registry.EtcRegisterOperation{
					Registry: reg,
				}
				jewel.boot.BindHttp(httpFun, registerOperation.HttpBindOp)
			} else {
				jewel.boot.BindHttp(httpFun)
			}
			jewel.boot.Close()
			return
		} else if cmd.Name == c {
			if fun, ok := jewel.cmdFunc[c]; ok {
				if cmd.Param["env"] == nil {
					jewel.boot = jewel.boot.Start(*cmd.Param["dir"].(*string), "")
				} else {
					jewel.boot = jewel.boot.Start(*cmd.Param["dir"].(*string), *cmd.Param["env"].(*string))

				}
				fun()
				jewel.boot.Close()
				return
			}
		}

	}
}

func (jewel *Jewel) getCmdParam() []CmdParam {
	var cmdParams []CmdParam
	for _, cmd := range jewel.cmd {
		target := cmd.Flag("config", "The directory where the configuration files are located").Default("./config").String()
		env := cmd.Flag("jewel.profiles.active", "The env where the configuration files are located").String()
		cmdParam := CmdParam{
			Param: make(map[string]interface{}),
		}
		cmdParam.Name = cmd.FullCommand()
		cmdParam.Param["dir"] = target
		cmdParam.Param["env"] = env
		cmdParams = append(cmdParams, cmdParam)
	}
	return cmdParams
}

func (jewel *Jewel) Start() {
	cmdParams := jewel.getCmdParam()
	c := kingpin.MustParse(jewel.app.Parse(os.Args[1:]))
	for _, cmd := range cmdParams {
		if cmd.Name == c {
			if fun, ok := jewel.cmdFunc[c]; ok {
				if cmd.Param["env"] == nil {
					jewel.boot = jewel.boot.Start(*cmd.Param["dir"].(*string), "")
				} else {
					jewel.boot = jewel.boot.Start(*cmd.Param["dir"].(*string), *cmd.Param["env"].(*string))

				}
				fun()
				jewel.boot.Close()
				return
			}

		}
	}
}
