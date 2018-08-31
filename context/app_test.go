package context

import (
	"testing"
	"fmt"
)

func TestRunWithConfigDir(t *testing.T) {
	pro := NewProperties()
	jewel := JewelProperties{}
	pro.Load("./../config/app-new.yml", &jewel)
	fmt.Println(*jewel.Jewel.MySql["primary"].URL)
	fmt.Printf("%+v", jewel)
}
