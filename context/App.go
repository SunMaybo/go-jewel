package context

import (
	"io/ioutil"
	"strings"
	"log"
	"os"
	"github.com/cihub/seelog"
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
	if name != "" {
		app.Load(dir + "/" + name)
	}
	return app
}
func GetCurrentDirectory(dir string) string {
	pwd, err := os.Getwd()
	if err != nil {
		seelog.Error(err)
		os.Exit(1)
	}
	if dir == "" {
		return pwd
	}
	if strings.HasPrefix(dir, ".") {
		return pwd + "/" + dir
	}
	return dir
}

func LoadFileName(dir string) string {
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
	return dir + "/" + name
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
	if name != "" {
		app := ConfigStruct{}
		app.Load(dir+"/"+name, inter)
	}
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

	if name != "" {
		app := ConfigStruct{}
		app.Load(dir+"/"+name, inter)
	}
}
func LoadEnvFileName(dir, env string) string {
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
	if name != "" {
		return dir + "/" + name
	} else {
		return ""
	}
}
