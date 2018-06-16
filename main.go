package main

import (
	"github.com/SunMaybo/go-jewel/context"
)

func main() {
	boot := context.NewInstance()
	boot.RunWithConfigDirAndExtend("./config", "www", func(cfgMap context.ConfigMap) {

	})
}
