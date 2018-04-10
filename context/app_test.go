package context

import (
	"testing"
	"fmt"
)

func TestRunWithConfigDir(t *testing.T) {
 boot:=	BootStrap{}
	boot.RunWithConfigDirAndExtend("./","www", func(cfgMap ConfigMap) {
		fmt.Println(cfgMap)
	})

}
