package context

import "github.com/SunMaybo/go-jewel/context"

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

func (c *Context) Db() context.Db {
	return c.Service(DB).(context.Db)
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
