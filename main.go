package main

import (
	"go-jewel/jewel"
	"github.com/gin-gonic/gin"
	"github.com/SunMaybo/jewel-inject/inject"
)

func main() {
  jewel:=jewel.NewHttp()
	jewel.HttpStart(func(router *gin.RouterGroup, injector *inject.Injector) {
		
	})
}
