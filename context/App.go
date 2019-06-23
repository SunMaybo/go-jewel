package context

import (
	"io/ioutil"
	"strings"
	"os"
)

var suffixCfgName = []string{"yml", "yaml", "xml", "json"}

func GetCurrentDirectory(dir string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if dir == "" {
		return pwd, nil
	}
	if strings.HasPrefix(dir, ".") {
		return pwd + "/" + dir, nil
	}
	return dir, nil
}

func LoadFileName(dir string) string {
	ifs, err := ioutil.ReadDir(dir)
	if err != nil {
		return ""
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
	if name == "" {
		return ""
	} else {
		return dir + "/" + name
	}
}
func LoadCfg(dir string, inter interface{}) {
	ifs, err := ioutil.ReadDir(dir)
	if err != nil {
		return
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
		return
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
		return ""
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
