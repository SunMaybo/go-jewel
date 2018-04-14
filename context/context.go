package context


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

func (c *Context) Db() Db {
	return c.Service(DB).(Db)
}

func (c *Context) web() interface{} {
	return c.Service(WEB)
}
