package jewel

import (
	"testing"
	"path/filepath"
	"os"
	"strings"
	"log"
	"fmt"
)

func TestApplicationBoot(t *testing.T) {
	fmt.Println(GetCurrentDirectory("./config"))
}
func GetCurrentDirectory(dir string) string {
	abs, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	root := strings.Replace(abs, "\\", "/", -1)
	if strings.HasPrefix(dir, ".") {
		dir = root +"/"+ dir
	}
	return dir
}
