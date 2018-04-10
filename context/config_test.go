package context

import (
	"testing"
	"fmt"
)

func TestConfig(t *testing.T) {
	app := Config{}
	app.Load("./app-www.yml")
	fmt.Println(app)
}
func TestApp(t *testing.T) {
	app := Load("./", "www")
	fmt.Println(app)
}

func TestConfigDir(t *testing.T) {
	hash:=make(map[interface{}]interface{})
	hash["1234"]=34
	if value,ok:=hash["1234"];ok {
		fmt.Println(value)
	}

}
