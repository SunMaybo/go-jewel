package main

import (
	"github.com/SunMaybo/go-jewel/context"
	"github.com/gin-gonic/gin"
)

func main() {
	boot := context.NewInstance()
	boot.Run(func(engine *gin.Engine) {

	})
}
