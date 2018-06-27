package main

import (
	"github.com/SunMaybo/go-jewel/context"
	"github.com/gin-gonic/gin"
)

func main() {
	boot := context.NewInstance()
	boot.Start().BindHttp(func(engine *gin.Engine) {
	})
}
