package main

import (
	"github.com/SunMaybo/go-jewel/context"
	"github.com/gin-gonic/gin"
)

func main() {
	var stu Stu
	boot := context.NewInstance()
	boot.AddApplyCfg(&stu)
	boot.Run(func(engine *gin.Engine) {

	})
}

type Stu struct {
	Name   string `yml:"name"`
	Age    string `yml:"age"`
	Gendar string `yml:"gendar"`
}
