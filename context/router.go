package context

import (
	"github.com/gin-gonic/gin"
	"container/list"
)

var array = list.New()

func load(engine *gin.Engine) {

	for e := array.Front(); e != nil; e = e.Next() {
		e.Value.(func(engine *gin.Engine))(engine)
	}
}
func registeries(fs []func(engine *gin.Engine)) {
	if fs == nil {
		return
	}
	for _, v := range fs {
		if v != nil {
			array.PushBack(v)
		}

	}
}
func register(fun func(engine *gin.Engine)) {
	if fun == nil {
		return
	}
	array.PushBack(fun)
}
