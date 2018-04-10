package context

import (
	"io/ioutil"
	"strings"
	"log"
)

func Load(dir string, env string) Config {
	ifs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	var name = ""
	for _, v := range ifs {
		if !v.IsDir() && strings.HasPrefix(v.Name(), "app"+"-"+env) {
			name = v.Name()
			break
		}
	}
	app := Config{}
	app.Load(dir + "/" + name)
	return app
}

func LoadMap(dir string, env string) map[interface{}]interface{} {
	ifs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	var name = ""
	for _, v := range ifs {
		if !v.IsDir() && strings.HasPrefix(v.Name(), "app"+"-"+env) {
			name = v.Name()
			break
		}
	}
	app := ConfigMap{}
	app.Load(dir + "/" + name)
	return app
}
