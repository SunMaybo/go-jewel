package main

import (
	"github.com/gin-gonic/gin"
	"github.com/SunMaybo/jewel-inject/inject"
	"go-jewel/jewel"
)

func main() {
	jewel := jewel.NewHttp()
	jewel.HttpStart(func(router *gin.RouterGroup, injector *inject.Injector) {

	})
}
