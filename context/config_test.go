package context

import (
	"testing"
	"fmt"
)

func TestConfig(t *testing.T) {
	app := ConfigMap{}
	app.Load("./app-www.yml")
	for _, v := range app {
		fmt.Println(v)
	}

}
func TestApp(t *testing.T) {
	app := Load("./")
	fmt.Println(app)
}

func TestConfigDir(t *testing.T) {
	hash := make(map[interface{}]interface{})
	hash["1234"] = 34
	if value, ok := hash["1234"]; ok {
		fmt.Println(value)
	}

}
