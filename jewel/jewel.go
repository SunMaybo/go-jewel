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

func NewHttp() *Jewel {
	jewel := &Jewel{
		boot:    context.NewInstance(),
		app:     kingpin.UsageTemplate(kingpin.DefaultUsageTemplate),
		cmdFunc: make(map[string]func()),
	}
	jewel.boot.AddApplyCfg(&registry.JewelRegisterProperties{})
	var plugin context.Plugin
	plugin = registry.EtcRegisterPlugin{}
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
	jewel.boot = jewel.boot.AddApplyCfg(properties)
	return jewel
}

func (jewel *Jewel) AddBean(beans ... interface{}) *Jewel {
	jewel.boot = jewel.boot.AddApply(beans)
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
func (jewel *Jewel) HttpStart(httpFun func(engine *gin.Engine)) {
	for _, cmd := range jewel.cmd {
		target := cmd.Flag("config", "The directory where the configuration files are located").Default("./config").String()
		c := kingpin.MustParse(jewel.app.Parse(os.Args[1:]))
		if cmd.FullCommand() == c && c == "server" {
			jewel.boot = jewel.boot.StartAndDir(*target)
			etcRegister := jewel.boot.GetInject().ServiceByName("plugin:etcd_register")
			if etcRegister!=nil {
			}
			jewel.boot.BindHttp(httpFun)
			jewel.boot.Close()
			return
		} else if cmd.FullCommand() == c {
			if fun, ok := jewel.cmdFunc[c]; ok {
				jewel.boot = jewel.boot.StartAndDir(*target)
				fun()
				jewel.boot.Close()
				return
			}
		}

	}
}

func (jewel *Jewel) Start() {
	for _, cmd := range jewel.cmd {
		target := cmd.Flag("config", "The directory where the configuration files are located").Default("./config").String()
		c := kingpin.MustParse(jewel.app.Parse(os.Args[1:]))
		if cmd.FullCommand() == c {
			if fun, ok := jewel.cmdFunc[c]; ok {
				jewel.boot = jewel.boot.StartAndDir(*target)
				fun()
				jewel.boot.Close()
				return
			}

		}
	}
}
