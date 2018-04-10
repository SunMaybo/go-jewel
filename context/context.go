package context

import (
	//"github.com/cihub/seelog"
	"github.com/regcostajr/go-web3/db"
)

const (
	DB  = "DB"
	WEB = "WEB"
	LOG = "LOG"
)

var Services = Context{ServiceMap: make(map[string]interface{})}

type IContext interface {
	Register(name string, inter interface{})
	Service(name string) interface{}
}

type Context struct {
	ServiceMap map[string]interface{}
}

func (c *Context) Register(name string, inter interface{}) {
	c.ServiceMap[name] = inter
}

func (c *Context) Service(name string) interface{} {
	return c.ServiceMap[name]
}

func (c *Context) Db() db.DB {
	return c.Service(DB).(db.DB)
}
/*func (c *Context) Log() seelog.LoggerInterface {
	if log, ok := c.Service(LOG).(seelog.LoggerInterface); ok {
		return log
	}
	return nil
}*/
func (c *Context) web() interface{} {
	return c.Service(WEB)
}
