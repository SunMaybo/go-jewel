package main

import (
	"github.com/SunMaybo/go-jewel/context"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/SunMaybo/jewel-template/template/rest"
)

func main() {
	boot := context.NewInstance()
	boot.Start().BindHttp(func(engine *gin.Engine) {
	})
}
