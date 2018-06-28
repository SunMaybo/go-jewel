package controller

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
	"os"
	"strings"
	"fmt"
	"log"
)

type TestController struct {

}
func (t TestController)Test(c *gin.Context)  {
 fmt.Println(getCurrentDirectory())
}
/*
获取程序运行路径
*/
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

