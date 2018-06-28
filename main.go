package main

import (
	"github.com/SunMaybo/go-jewel/context"
	"github.com/gin-gonic/gin"
	"fmt"
)

func main() {
	type Stu struct {
		Name string
		Age  string
	}
	boot := context.NewInstance()
	boot.JsonRpc().RegisterJsonRpc("test", func(name string, age float64, stu map[string]interface{}) string {
		return "ok"
	})
	boot.GetCmd().PutFlagString("f", "./config", "open a file")
	boot.GetCmd().PutCmd("start", func(c context.Config) {
		fmt.Println("start ....")
		f:=boot.GetCmd().Params["f"]
		fmt.Println(*f)
	})

	boot.Start().BindJsonRpc("/test", func(engine *gin.Engine) {
		engine.GET("/students", func(c *gin.Context) {

		})
	})
}
