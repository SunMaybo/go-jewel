package context

import (
	"github.com/gin-gonic/gin"
	"container/list"
	"github.com/SunMaybo/jewel-inject/inject"
)

var array = list.New()

func load(router *gin.RouterGroup) {

	for e := array.Front(); e != nil; e = e.Next() {
		e.Value.(func(engine *gin.RouterGroup, injector *inject.Injector))(router, nil)
	}
}
func registeries(fs []func(router *gin.RouterGroup, injector *inject.Injector)) {
	if fs == nil {
		return
	}
	for _, v := range fs {
		if v != nil {
			array.PushBack(v)
		}

	}
}
func register(fun func(engine *gin.RouterGroup, injector *inject.Injector)) {
	if fun == nil {
		return
	}
	array.PushBack(fun)
}
