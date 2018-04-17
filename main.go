package main

import (
	"go-jewel/context"
	"go-jewel/example"
)

func main() {
	/*boot := context.NewInstance()

	boot.GetCmd().PutCmd("start", func(c context.Config) {
		fmt.Println(c)
	})*/
	/*boot.RunWithExtend(func(engine *gin.Engine) {
		engine.GET("/info", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello World!")
		})
	}, func(cfgMap context.ConfigMap) {
		fmt.Println(cfgMap)
	})*/

	/*boot.Run3("./config", "www", nil, func(engine *gin.Engine) {
		engine.POST("/health", func(i *gin.Context) {
		})
	})*/
	boot := context.NewInstance()
	boot.JsonRpc().RegisterJsonRpc("test", example.Test)
	boot.RunJsonRpc2("/", "./config", "www", nil, nil)
}
