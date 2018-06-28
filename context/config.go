package context

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"encoding/json"
	"encoding/xml"
	"strings"
)

type Config struct {
	Jewel struct {
		Log            string `json:"log" yaml:"log"`
		Max_Open_Conns int    `json:"max-open-conns" yaml:"max-open-conns"`
		Max_Idle_Conns int    `json:"max-idle-conns" yaml:"max-idle-conns" xml:"max-idle-conns"`
		SqlShow        bool   `json:"sql_show" yaml:"sql_show" xml:"sql_show"`
		Mysql          string `json:"mysql" yaml:"mysql" xml:"mysql"`
		Name           string `json:"name" yaml:"name" xml:"name"`
		Port           int    `json:"port" yaml:"port" xml:"port"`
		Postgres       string `json:"postgres" yaml:"postgres" xml:"postgres"`
		Profiles struct {
			Active string `json:"active"`
		} `json:"profiles" yaml:"profiles" xml:"profiles"`
		Redis struct {
			Db       int    `json:"db"`
			Host     string `json:"host"`
			Password string `json:"password"`
		} `json:"redis" yaml:"redis" xml:"redis"`
		JsonRpc struct {
			Enabled  *bool   `json:"enabled"`
			UserName string `json:"username"`
			Password string `json:"password"`
		} `json:"jsonrpc" yaml:"jsonrpc" xml:"jsonrpc"`
		Sqlite3 string `json:"sqlite3" yaml:"sqlite3" xml:"sqlite3"`
	} `json:"jewel" yaml:"jewel" xml:"jewel"`
}

func (config *Config) Load(fileName string) {
	if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		config.loadYml(fileName)
	} else if strings.HasSuffix(fileName, ".xml") {
		config.loadXml(fileName)
	} else {
		config.loadJson(fileName)
	}
}

func (config *Config) loadYml(fileName string) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}

}
func (config *Config) loadXml(fileName string) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	xml.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}
}
func (config *Config) loadJson(fileName string) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}
}

type ConfigStruct struct {
}

func (config *ConfigStruct) Load(fileName string, inter interface{}) {
	if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		config.loadYml(fileName, inter)
	} else if strings.HasSuffix(fileName, ".xml") {
		config.loadXml(fileName, inter)
	} else {
		config.loadJson(fileName, inter)
	}
}

func (config *ConfigStruct) loadYml(fileName string, inter interface{}) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(buff, inter)
	if err != nil {
		log.Fatalln(err)
	}

}
func (config *ConfigStruct) loadXml(fileName string, inter interface{}) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = xml.Unmarshal(buff, &inter)
	if err != nil {
		log.Fatalln(err)
	}
}
func (config *ConfigStruct) loadJson(fileName string, inter interface{}) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(buff, &inter)
	if err != nil {
		log.Fatalln(err)
	}
}
