package context

import (
	"io/ioutil"
	"strings"
	"log"
)

var suffixCfgName = []string{"yml", "yaml", "xml", "json"}

func Load(dir string) Config {
	ifs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	var name = ""
loop:
	for _, v := range ifs {
		for e := range suffixCfgName {
			if !v.IsDir() && strings.HasPrefix(v.Name(), "app."+suffixCfgName[e]) {
				name = v.Name()
				break loop
			}
		}

	}
	app := Config{}
	app.Load(dir + "/" + name)
	return app
}
func LoadCfg(dir string, inter interface{}) {
	ifs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	var name = ""
loop:
	for _, v := range ifs {
		for e := range suffixCfgName {
			if !v.IsDir() && strings.HasPrefix(v.Name(), "app."+suffixCfgName[e]) {
				name = v.Name()
				break loop
			}
		}
	}
	app := ConfigStruct{}
	app.Load(dir+"/"+name, inter)
}
func LoadEnvCfg(dir, env string, inter interface{}) {
	ifs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	var name = ""
loop:
	for _, v := range ifs {
		for e := range suffixCfgName {
			if !v.IsDir() && strings.HasPrefix(v.Name(), "app-"+env+"."+suffixCfgName[e]) {
				name = v.Name()
				break loop
			}
		}
	}
	app := ConfigStruct{}
	app.Load(dir+"/"+name, inter)
}
