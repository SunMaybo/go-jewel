package router

import (
	"github.com/gin-gonic/gin"
	"go-jewel/controller"
)

func Router(engine *gin.Engine) {
	t := controller.TestController{}
	engine.GET("", t.Test)
}
