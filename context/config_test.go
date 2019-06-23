package context

import (
	"testing"
	"fmt"
)

func TestConfigDir(t *testing.T) {
	hash := make(map[interface{}]interface{})
	hash["1234"] = 34
	if value, ok := hash["1234"]; ok {
		fmt.Println(value)
	}

}
