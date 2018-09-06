package main

import (
	"github.com/SunMaybo/go-jewel/jewel"
)

func main() {
	jewel := jewel.New()
	jewel.Cmd("redis_start", "", func() {

	})
	jewel.Cmd("start", "", func() {

	})
	jewel.Start()
}
