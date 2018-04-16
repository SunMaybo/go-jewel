package main

import (
	"fmt"
	"reflect"
	"strconv"
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

	/*h1 := hello
	fv := reflect.ValueOf(h1)
	fmt.Println(fv.Kind() == reflect.Func)
	values := make([]reflect.Value, 1)
	values[0] = reflect.ValueOf("Hello World")
	fv.Call(values)*/

	myType := &MyType{22,"sunmaybo"}
	test(myType)
	//fmt.Println(myType)     //就是检查一下myType对象内容
	//println("---------------")

}
func test(ty interface{})  {
	mtV := reflect.ValueOf(&ty).Elem()
	fmt.Println("Before:",mtV.MethodByName("String").Call(nil)[0])
	params := make([]reflect.Value,1)
	params[0] = reflect.ValueOf(18)
	mtV.MethodByName("SetI").Call(params)
	params[0] = reflect.ValueOf("reflection test")
	mtV.MethodByName("SetName").Call(params)
	fmt.Println("After:",mtV.MethodByName("String").Call(nil)[0])
}

func hello(str string) {
	fmt.Println(str)
}

type MyType struct {
	i    int
	name string
}

func (mt *MyType) SetI(i int) {
	mt.i = i
}

func (mt *MyType) SetName(name string) {
	mt.name = name
}

func (mt *MyType) String() string {
	return fmt.Sprintf("%p", mt) + "--name:" + mt.name + " i:" + strconv.Itoa(mt.i)
}
