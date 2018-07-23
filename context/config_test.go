package context

import (
	"testing"
	"fmt"
)

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
