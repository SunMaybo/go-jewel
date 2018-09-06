package main

import (

	"github.com/gin-gonic/gin"
	"go-jewel/jewel"
)

func main() {
	jewel := jewel.NewHttp()
	jewel.HttpStart(func(engine *gin.Engine) {
	})

}

