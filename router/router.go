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

func Register(funs ... func(engine *gin.Engine)) {
	for _, v := range funs {
		if v != nil {
			array.PushBack(v)
		}

	}
}
