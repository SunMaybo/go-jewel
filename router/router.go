package router

import (
	"github.com/gin-gonic/gin"
	"container/list"
)

var array = list.New()

func Load(engine *gin.Engine) {

	for e := array.Front(); e != nil; e = e.Next() {
		e.Value.(func(engine *gin.Engine))(engine)
	}
}

func Registeries(fs []func(engine *gin.Engine)) {
	for _, v := range fs {
		if v != nil {
			array.PushBack(v)
		}

	}
}
func Register(fun func(engine *gin.Engine)) {
	array.PushBack(fun)
}
