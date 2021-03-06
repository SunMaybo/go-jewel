package context

import (
	"testing"
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/SunMaybo/jewel-inject/inject"
)

func TestBootStrap(t *testing.T) {
	boot := NewInstance()
	boot.AddFun(func(injector *inject.Injector) {
		db := injector.ServiceByName("mysql.default").(gorm.DB)
		fmt.Print(db)
	})
	boot.StartAndDir("./../config")
}
func TestPlugin(t *testing.T) {
	var plugin Plugin
	plugin = &BasePlugin{}
	fmt.Println(plugin)
}
