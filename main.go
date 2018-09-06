package main

import (

	"github.com/gin-gonic/gin"
	"github.com/SunMaybo/go-jewel/jewel"
)

func main() {
	jewel := jewel.NewHttp()
	jewel.HttpStart(func(engine *gin.Engine) {
	})

}

